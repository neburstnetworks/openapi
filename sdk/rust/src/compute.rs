use std::collections::BTreeMap;
use crate::client::NeburstClient;
use crate::types::{Instance, Metrics, OSProfile, PaginatedResult, PowerStatus, Result, Traffic};

impl NeburstClient {
    // ── Cloud Instance ──

    pub async fn list_instances(&self, page: i32, page_size: i32) -> Result<PaginatedResult<Instance>> {
        let mut q = BTreeMap::new();
        q.insert("page".into(), page.to_string());
        q.insert("page_size".into(), page_size.to_string());
        self.request("GET", "/open/v1/compute/instance/list", Some(&q), None::<&()>).await
    }

    pub async fn get_instance(&self, id: &str) -> Result<Instance> {
        self.request("GET", &format!("/open/v1/compute/instance/{}", id), None, None::<&()>).await
    }

    pub async fn get_instance_status(&self, id: &str) -> Result<PowerStatus> {
        self.request("GET", &format!("/open/v1/compute/instance/{}/status", id), None, None::<&()>).await
    }

    pub async fn get_instance_traffic(&self, id: &str) -> Result<Traffic> {
        self.request("GET", &format!("/open/v1/compute/instance/{}/traffic", id), None, None::<&()>).await
    }

    pub async fn cloud_power_action(&self, id: &str, action: &str) -> Result<()> {
        let body = serde_json::json!({ "action": action });
        self.request_void("POST", &format!("/open/v1/compute/instance/{}/power", id), None, Some(&body)).await
    }

    pub async fn get_cloud_metrics(&self, id: &str) -> Result<Metrics> {
        self.request("GET", &format!("/open/v1/compute/instance/{}/metrics", id), None, None::<&()>).await
    }

    // ── Bare Metal ──

    pub async fn list_bare_metal_instances(&self, page: i32, page_size: i32) -> Result<PaginatedResult<Instance>> {
        let mut q = BTreeMap::new();
        q.insert("page".into(), page.to_string());
        q.insert("page_size".into(), page_size.to_string());
        self.request("GET", "/open/v1/compute/bare-metal/list", Some(&q), None::<&()>).await
    }

    pub async fn get_bare_metal_instance(&self, id: &str) -> Result<Instance> {
        self.request("GET", &format!("/open/v1/compute/bare-metal/{}", id), None, None::<&()>).await
    }

    pub async fn get_bare_metal_status(&self, id: &str) -> Result<PowerStatus> {
        self.request("GET", &format!("/open/v1/compute/bare-metal/{}/status", id), None, None::<&()>).await
    }

    pub async fn get_bare_metal_traffic(&self, id: &str) -> Result<Traffic> {
        self.request("GET", &format!("/open/v1/compute/bare-metal/{}/traffic", id), None, None::<&()>).await
    }

    pub async fn bare_metal_power_action(&self, id: &str, action: &str) -> Result<()> {
        let body = serde_json::json!({ "action": action });
        self.request_void("POST", &format!("/open/v1/compute/bare-metal/{}/power", id), None, Some(&body)).await
    }

    pub async fn get_bare_metal_metrics(&self, id: &str) -> Result<Metrics> {
        self.request("GET", &format!("/open/v1/compute/bare-metal/{}/metrics", id), None, None::<&()>).await
    }

    pub async fn get_reinstall_profiles(&self, id: &str) -> Result<Vec<OSProfile>> {
        self.request("GET", &format!("/open/v1/compute/bare-metal/{}/profiles", id), None, None::<&()>).await
    }

    pub async fn get_rescue_profiles(&self, id: &str) -> Result<Vec<OSProfile>> {
        self.request("GET", &format!("/open/v1/compute/bare-metal/{}/rescue-profiles", id), None, None::<&()>).await
    }

    pub async fn rebuild_instance(&self, id: &str, profile_id: i64, hostname: Option<&str>, public_keys: Option<&[String]>) -> Result<()> {
        let mut body = serde_json::json!({ "profile_id": profile_id });
        if let Some(h) = hostname { body["hostname"] = serde_json::Value::String(h.to_owned()); }
        if let Some(keys) = public_keys { body["public_keys"] = serde_json::json!(keys); }
        self.request_void("POST", &format!("/open/v1/compute/bare-metal/{}/rebuild", id), None, Some(&body)).await
    }

    pub async fn rescue_instance(&self, id: &str, profile_id: i64) -> Result<()> {
        let body = serde_json::json!({ "profile_id": profile_id });
        self.request_void("POST", &format!("/open/v1/compute/bare-metal/{}/rescue", id), None, Some(&body)).await
    }
}
