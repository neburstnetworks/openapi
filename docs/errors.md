# Error Codes

**[中文版](errors_zh.md)**

All API responses use the standard envelope:

```json
{
  "code": <int>,
  "msg": "<string>",
  "data": ...
}
```

When `code` is `0`, the request succeeded. Non-zero codes indicate an error -- the `msg` field contains a human-readable description.

## Error Code Reference

| Code | Name | HTTP Status | Description |
|------|------|-------------|-------------|
| 58200 | `INVALID_SIGNATURE` | 401 | The `X-Nb-Signature` header does not match the server's computed signature. Double-check your signing algorithm, secret, and string-to-sign construction. |
| 58201 | `EXPIRED_TIMESTAMP` | 401 | The `X-Nb-Timestamp` value is more than 5 minutes away from the server's current time. Ensure your system clock is synchronized via NTP. |
| 58202 | `DUPLICATE_NONCE` | 401 | The `X-Nb-Nonce` value was already used within the last 10 minutes. Each request must use a unique nonce (UUID v4 recommended). |
| 58203 | `INSUFFICIENT_SCOPE` | 403 | The API key does not have the scope required for this endpoint. Check the key's assigned scopes in the dashboard. |
| 58204 | `RATE_LIMITED` | 429 | The per-user rate limit of 60 requests/minute has been exceeded. Inspect the `X-RateLimit-Reset` header for when the window resets. |
| 58205 | `KEY_DISABLED` | 401 | The API key has been manually disabled. Re-enable it in the dashboard or create a new key. |
| 58206 | `KEY_EXPIRED` | 401 | The API key has passed its configured expiration date. Create a new key. |
| 58207 | `KEY_NOT_FOUND` | 401 | The `X-Nb-Key` value does not match any existing API key. Verify the Key ID is correct. |
| 58208 | `MISSING_HEADERS` | 400 | One or more of the required authentication headers (`X-Nb-Key`, `X-Nb-Timestamp`, `X-Nb-Nonce`, `X-Nb-Signature`) is missing from the request. |

## Business Logic Errors

In addition to the authentication/authorization errors above, endpoints may return application-level errors:

| Code | HTTP Status | Description |
|------|-------------|-------------|
| 58210 | 200 | A business logic error occurred (e.g., instance not found, invalid power state). The `msg` field contains the specific error message. |
| 58400 | 200 | The request body could not be parsed. Ensure valid JSON and correct field types. |

## Error Response Examples

### Authentication failure

```bash
curl -s https://api.neburst.com/open/v1/compute/instance/list \
  -H "X-Nb-Key: nb_key_invalid"
```

```json
{
  "code": 58208,
  "msg": "Missing required authentication headers"
}
```

### Invalid signature

```json
{
  "code": 58200,
  "msg": "Signature verification failed"
}
```

### Insufficient scope

Attempting to call `POST /compute/instance/{id}/power` with a key that only has `instance:read`:

```json
{
  "code": 58203,
  "msg": "Insufficient scope"
}
```

### Rate limited

```json
{
  "code": 58204,
  "msg": "Rate limit exceeded"
}
```

Response headers when rate limited:

```
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 0
X-RateLimit-Reset: 1750291260
```

### Invalid request body

```json
{
  "code": 58400,
  "msg": "invalid request body"
}
```

## Troubleshooting

| Symptom | Likely Cause | Fix |
|---------|-------------|-----|
| Always getting `58200` | Incorrect signing algorithm | Verify each component of the string-to-sign. See [Authentication](authentication.md). |
| `58201` on every request | Clock drift | Run `ntpdate` or enable NTP synchronization. |
| `58202` intermittently | Nonce collision | Use UUID v4; do not reuse nonces across requests. |
| `58203` on certain endpoints | Missing scope | Add the required scope to the key via the dashboard. See [Scopes](scopes.md). |
| `58207` after key rotation | Using old Key ID | After revoking/rotating, only the secret changes. If you deleted and recreated, the Key ID changes too. |
