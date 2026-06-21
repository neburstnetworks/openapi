use base64::Engine;
use hmac::{Hmac, Mac};
use reqwest::Client;
use serde::de::DeserializeOwned;
use serde::Serialize;
use sha2::{Digest, Sha256};
use std::collections::BTreeMap;
use std::time::{SystemTime, UNIX_EPOCH};
use uuid::Uuid;

use crate::types::{ApiResponse, NeburstError, Result};

const EMPTY_BODY_SHA256: &str = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855";

/// Core HTTP client for the Neburst OpenAPI.
///
/// Handles HMAC-SHA256 request signing, JSON serialisation/deserialisation,
/// and the generic `{code, msg, data}` response envelope.
///
/// # Example
///
/// ```no_run
/// use neburst_api::NeburstClient;
///
/// #[tokio::main]
/// async fn main() -> neburst_api::Result<()> {
///     let client = NeburstClient::new(
///         "https://api.neburst.com",
///         "your-key-id",
///         "your-api-secret",
///     );
///     let instances = client.list_instances().await?;
///     Ok(())
/// }
/// ```
pub struct NeburstClient {
    base_url: String,
    key_id: String,
    secret: String,
    http: Client,
}

impl NeburstClient {
    /// Creates a new Neburst API client with key ID and secret.
    pub fn new(base_url: &str, key_id: &str, secret: &str) -> Self {
        Self {
            base_url: base_url.trim_end_matches('/').to_owned(),
            key_id: key_id.to_owned(),
            secret: secret.to_owned(),
            http: Client::builder()
                .timeout(std::time::Duration::from_secs(30))
                .build()
                .expect("failed to build reqwest client"),
        }
    }

    /// Creates a new client from a base64-encoded combined key.
    pub fn from_combined_key(base_url: &str, combined: &str) -> std::result::Result<Self, String> {
        let decoded = base64::engine::general_purpose::STANDARD
            .decode(combined)
            .map_err(|e| format!("base64 decode: {}", e))?;
        let parsed: serde_json::Value = serde_json::from_slice(&decoded)
            .map_err(|e| format!("json parse: {}", e))?;
        let key_id = parsed["key_id"].as_str().ok_or("missing key_id")?;
        let secret = parsed["secret"].as_str().ok_or("missing secret")?;
        Ok(Self::new(base_url, key_id, secret))
    }

    // ------------------------------------------------------------------
    // Internal request helpers
    // ------------------------------------------------------------------

    /// Executes a signed API request and deserialises the `data` field.
    pub(crate) async fn request<T: DeserializeOwned>(
        &self,
        method: &str,
        path: &str,
        query: Option<&BTreeMap<String, String>>,
        body: Option<&impl Serialize>,
    ) -> Result<T> {
        let resp = self.do_request(method, path, query, body).await?;

        let envelope: ApiResponse<serde_json::Value> = resp.json().await?;
        if envelope.code != 0 {
            return Err(NeburstError::Api {
                code: envelope.code,
                message: envelope.msg,
            });
        }

        let data = envelope.data.ok_or_else(|| NeburstError::Api {
            code: -1,
            message: "response data is null".into(),
        })?;

        serde_json::from_value(data).map_err(NeburstError::Json)
    }

    /// Executes a signed API request that returns no data (e.g. power action).
    pub(crate) async fn request_void(
        &self,
        method: &str,
        path: &str,
        query: Option<&BTreeMap<String, String>>,
        body: Option<&impl Serialize>,
    ) -> Result<()> {
        let resp = self.do_request(method, path, query, body).await?;

        let envelope: ApiResponse<serde_json::Value> = resp.json().await?;
        if envelope.code != 0 {
            return Err(NeburstError::Api {
                code: envelope.code,
                message: envelope.msg,
            });
        }
        Ok(())
    }

    /// Low-level request: signs, sends, and returns the raw `reqwest::Response`.
    async fn do_request(
        &self,
        method: &str,
        path: &str,
        query: Option<&BTreeMap<String, String>>,
        body: Option<&impl Serialize>,
    ) -> Result<reqwest::Response> {
        // Serialise body
        let body_bytes: Vec<u8> = match body {
            Some(b) => serde_json::to_vec(b)?,
            None => Vec::new(),
        };

        // Signing material
        let timestamp = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .expect("system clock before UNIX epoch")
            .as_secs()
            .to_string();
        let nonce = Uuid::new_v4().to_string();
        let sorted_query = Self::build_sorted_query(query);
        let body_hash = if body_bytes.is_empty() {
            EMPTY_BODY_SHA256.to_owned()
        } else {
            Self::sha256_hex(&body_bytes)
        };

        let signature = self.sign(&timestamp, method, path, &sorted_query, &body_hash);

        // Build URL
        let mut url = format!("{}{}", self.base_url, path);
        if !sorted_query.is_empty() {
            url.push('?');
            url.push_str(&sorted_query);
        }

        // Build request
        let mut req = self
            .http
            .request(method.parse().unwrap(), &url)
            .header("X-Nb-Key", &self.key_id)
            .header("X-Nb-Timestamp", &timestamp)
            .header("X-Nb-Nonce", &nonce)
            .header("X-Nb-Signature", &signature);

        if !body_bytes.is_empty() {
            req = req
                .header("Content-Type", "application/json")
                .body(body_bytes);
        }

        let resp = req.send().await?;
        Ok(resp)
    }

    // ------------------------------------------------------------------
    // Signing utilities
    // ------------------------------------------------------------------

    /// Computes the HMAC-SHA256 signature.
    ///
    /// ```text
    /// StringToSign = timestamp + "\n" + method + "\n" + path + "\n" + sorted_query + "\n" + body_hash
    /// signature = hex(HMAC-SHA256(secret, StringToSign))
    /// ```
    fn sign(
        &self,
        timestamp: &str,
        method: &str,
        path: &str,
        sorted_query: &str,
        body_hash: &str,
    ) -> String {
        let string_to_sign = format!(
            "{}\n{}\n{}\n{}\n{}",
            timestamp, method, path, sorted_query, body_hash
        );

        type HmacSha256 = Hmac<Sha256>;
        let mut mac =
            HmacSha256::new_from_slice(self.secret.as_bytes()).expect("HMAC key length error");
        mac.update(string_to_sign.as_bytes());
        hex::encode(mac.finalize().into_bytes())
    }

    /// SHA-256 hex digest.
    fn sha256_hex(data: &[u8]) -> String {
        let mut hasher = Sha256::new();
        hasher.update(data);
        hex::encode(hasher.finalize())
    }

    /// Builds a sorted, percent-encoded query string from a `BTreeMap`.
    ///
    /// Keys are already sorted (BTreeMap guarantees this). Each key and value
    /// is individually percent-encoded. Pairs are joined with `&`.
    /// Returns an empty string if the map is `None` or empty.
    fn build_sorted_query(params: Option<&BTreeMap<String, String>>) -> String {
        match params {
            None => String::new(),
            Some(map) if map.is_empty() => String::new(),
            Some(map) => map
                .iter()
                .map(|(k, v)| {
                    format!(
                        "{}={}",
                        urlencoding::encode(k),
                        urlencoding::encode(v)
                    )
                })
                .collect::<Vec<_>>()
                .join("&"),
        }
    }
}
