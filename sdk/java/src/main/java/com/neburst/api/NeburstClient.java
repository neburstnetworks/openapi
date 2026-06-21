package com.neburst.api;

import com.google.gson.Gson;
import com.google.gson.JsonElement;
import com.google.gson.reflect.TypeToken;
import com.neburst.api.model.ApiResponse;

import javax.crypto.Mac;
import javax.crypto.spec.SecretKeySpec;
import java.io.IOException;
import java.lang.reflect.Type;
import java.net.URI;
import java.net.URLEncoder;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.nio.charset.StandardCharsets;
import java.security.InvalidKeyException;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.time.Duration;
import java.time.Instant;
import java.util.Map;
import java.util.TreeMap;
import java.util.UUID;
import java.util.stream.Collectors;

/**
 * Core HTTP client for the Neburst OpenAPI.
 * Handles request signing (HMAC-SHA256) and response parsing.
 *
 * <pre>{@code
 * NeburstClient client = new NeburstClient("https://api.neburst.com", "your-key-id", "your-secret");
 * ComputeApi compute = new ComputeApi(client);
 * List<Instance> instances = compute.listInstances();
 * }</pre>
 */
public class NeburstClient {

    private static final String EMPTY_BODY_HASH = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855";

    private final String baseUrl;
    private final String keyId;
    private final String secret;
    private final HttpClient httpClient;
    private final Gson gson;

    /**
     * Creates a new NeburstClient.
     *
     * @param baseUrl the API base URL (e.g., "https://api.neburst.com")
     * @param keyId   the API key ID
     * @param secret  the API secret used for HMAC signing
     */
    public NeburstClient(String baseUrl, String keyId, String secret) {
        this.baseUrl = baseUrl.endsWith("/") ? baseUrl.substring(0, baseUrl.length() - 1) : baseUrl;
        this.keyId = keyId;
        this.secret = secret;
        this.httpClient = HttpClient.newBuilder()
                .connectTimeout(Duration.ofSeconds(30))
                .build();
        this.gson = new Gson();
    }

    /**
     * Executes an API request and deserializes the response data to the given type.
     *
     * @param method       HTTP method (GET, POST, PUT, DELETE)
     * @param path         API path (e.g., "/open/v1/compute/instances")
     * @param query        query parameters (nullable)
     * @param body         request body object (nullable, will be serialized to JSON)
     * @param responseType the type of the data field in the response (use TypeToken for generics)
     * @param <T>          the expected data type
     * @return the deserialized data field
     */
    <T> T doRequest(String method, String path, Map<String, String> query, Object body, Type responseType)
            throws NeburstApiException, IOException, InterruptedException {

        String bodyJson = (body != null) ? gson.toJson(body) : "";
        byte[] bodyBytes = bodyJson.isEmpty() ? new byte[0] : bodyJson.getBytes(StandardCharsets.UTF_8);

        String timestamp = String.valueOf(Instant.now().getEpochSecond());
        String nonce = UUID.randomUUID().toString();
        String sortedQuery = buildSortedQueryString(query);
        String bodyHash = bodyBytes.length == 0 ? EMPTY_BODY_HASH : sha256Hex(bodyBytes);

        String signature = sign(timestamp, method, path, sortedQuery, bodyHash);

        // Build URI
        String fullUrl = baseUrl + path;
        if (!sortedQuery.isEmpty()) {
            fullUrl += "?" + sortedQuery;
        }

        HttpRequest.Builder requestBuilder = HttpRequest.newBuilder()
                .uri(URI.create(fullUrl))
                .timeout(Duration.ofSeconds(60))
                .header("X-Nb-Key", keyId)
                .header("X-Nb-Timestamp", timestamp)
                .header("X-Nb-Nonce", nonce)
                .header("X-Nb-Signature", signature);

        if (body != null) {
            requestBuilder.header("Content-Type", "application/json");
            requestBuilder.method(method, HttpRequest.BodyPublishers.ofByteArray(bodyBytes));
        } else if ("GET".equalsIgnoreCase(method)) {
            requestBuilder.GET();
        } else {
            requestBuilder.method(method, HttpRequest.BodyPublishers.noBody());
        }

        HttpResponse<String> response = httpClient.send(requestBuilder.build(), HttpResponse.BodyHandlers.ofString());

        Type envelopeType = TypeToken.getParameterized(ApiResponse.class, JsonElement.class).getType();
        ApiResponse<JsonElement> apiResponse = gson.fromJson(response.body(), envelopeType);

        if (apiResponse == null) {
            throw new NeburstApiException(-1, "Failed to parse API response");
        }

        if (!apiResponse.isSuccess()) {
            throw new NeburstApiException(apiResponse.getCode(), apiResponse.getMsg());
        }

        if (responseType == null) {
            return null;
        }

        return gson.fromJson(apiResponse.getData(), responseType);
    }

    /**
     * Executes an API request that does not return data (void).
     * Still validates the response envelope for errors.
     */
    void doRequest(String method, String path, Map<String, String> query, Object body)
            throws NeburstApiException, IOException, InterruptedException {
        doRequest(method, path, query, body, null);
    }

    /**
     * Computes the HMAC-SHA256 signature for a request.
     *
     * <pre>
     * StringToSign = timestamp + "\n" + method + "\n" + path + "\n" + sortedQuery + "\n" + bodyHash
     * signature = Hex(HMAC-SHA256(secret, StringToSign))
     * </pre>
     */
    private String sign(String timestamp, String method, String path, String sortedQuery, String bodyHash) {
        String stringToSign = timestamp + "\n" + method + "\n" + path + "\n" + sortedQuery + "\n" + bodyHash;
        return hmacSha256Hex(secret, stringToSign);
    }

    /**
     * Builds a sorted query string from the given parameters.
     * Keys are sorted alphabetically; each key and value is URL-encoded individually.
     *
     * @param params query parameters (nullable)
     * @return the sorted query string, or empty string if params is null/empty
     */
    private String buildSortedQueryString(Map<String, String> params) {
        if (params == null || params.isEmpty()) {
            return "";
        }
        TreeMap<String, String> sorted = new TreeMap<>(params);
        return sorted.entrySet().stream()
                .map(e -> URLEncoder.encode(e.getKey(), StandardCharsets.UTF_8)
                        + "=" + URLEncoder.encode(e.getValue(), StandardCharsets.UTF_8))
                .collect(Collectors.joining("&"));
    }

    /**
     * Computes SHA-256 hash and returns the hex string.
     */
    private String sha256Hex(byte[] data) {
        try {
            MessageDigest digest = MessageDigest.getInstance("SHA-256");
            byte[] hash = digest.digest(data);
            return bytesToHex(hash);
        } catch (NoSuchAlgorithmException e) {
            throw new RuntimeException("SHA-256 not available", e);
        }
    }

    /**
     * Computes HMAC-SHA256 and returns the hex string.
     */
    private String hmacSha256Hex(String secret, String data) {
        try {
            Mac mac = Mac.getInstance("HmacSHA256");
            SecretKeySpec keySpec = new SecretKeySpec(secret.getBytes(StandardCharsets.UTF_8), "HmacSHA256");
            mac.init(keySpec);
            byte[] hash = mac.doFinal(data.getBytes(StandardCharsets.UTF_8));
            return bytesToHex(hash);
        } catch (NoSuchAlgorithmException | InvalidKeyException e) {
            throw new RuntimeException("HMAC-SHA256 computation failed", e);
        }
    }

    /**
     * Converts a byte array to a lowercase hex string.
     */
    private static String bytesToHex(byte[] bytes) {
        StringBuilder sb = new StringBuilder(bytes.length * 2);
        for (byte b : bytes) {
            sb.append(String.format("%02x", b));
        }
        return sb.toString();
    }
}
