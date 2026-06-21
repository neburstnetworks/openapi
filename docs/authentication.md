# Authentication

**[中文版](authentication_zh.md)**

The Neburst OpenAPI uses **HMAC-SHA256** request signing. Every request must include four custom headers. There are no bearer tokens or session cookies -- each request is independently authenticated and verified.

## Required Headers

| Header | Format | Description |
|--------|--------|-------------|
| `X-Nb-Key` | `nb_key_` + 24 hex chars | Your API Key ID |
| `X-Nb-Timestamp` | Unix epoch (seconds) | Current time; must be within **+/- 5 minutes** of the server clock |
| `X-Nb-Nonce` | UUID v4 | Unique per request; the server rejects duplicates within a **10-minute** window |
| `X-Nb-Signature` | 64 hex chars | `Hex(HMAC-SHA256(api_secret, string_to_sign))` |

## Signing Algorithm

### Step 1 -- Build the String to Sign

Concatenate five components separated by newline characters (`\n`):

```
StringToSign = timestamp + "\n"
             + method    + "\n"
             + path      + "\n"
             + sorted_query_string + "\n"
             + SHA256(body)
```

| Component | Rules |
|-----------|-------|
| **timestamp** | The same value sent in `X-Nb-Timestamp` |
| **method** | Uppercase HTTP method: `GET`, `POST`, `PUT`, `DELETE` |
| **path** | Request path starting from `/open/v1/...` (no scheme, host, or query string) |
| **sorted_query_string** | Query parameters sorted alphabetically by key. Each key and value is URL-encoded. Pairs joined with `&`. Empty string if no query parameters. |
| **SHA256(body)** | Hex-encoded SHA-256 hash of the raw request body. For requests with no body (GET, DELETE), hash the empty string: `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |

### Step 2 -- Compute the Signature

```
Signature = Hex(HMAC-SHA256(api_secret, StringToSign))
```

- `api_secret` is the secret returned when you created the API key (starts with `nb_secret_`).
- The output is a lowercase hex string (64 characters).

### Step 3 -- Send the Request

Attach all four headers to your HTTP request:

```
X-Nb-Key: nb_key_1a2b3c4d5e6f7a8b9c0d1e2f
X-Nb-Timestamp: 1750291200
X-Nb-Nonce: 550e8400-e29b-41d4-a716-446655440000
X-Nb-Signature: 7f83b1657ff1fc53b92dc18148a1d65dfc2d4b1fa3d677284addd200126d9069
```

## Worked Example

Suppose you want to list cloud instances at `2026-06-18T12:00:00Z` (Unix timestamp `1750248000`).

**Request parameters:**

| Field | Value |
|-------|-------|
| Method | `GET` |
| URL | `https://api.neburst.com/open/v1/compute/instance/list` |
| API Key ID | `nb_key_a1b2c3d4e5f6a7b8c9d0e1f2` |
| API Secret | `nb_secret_aabbccdd11223344556677889900aabb` |
| Timestamp | `1750248000` |
| Nonce | `550e8400-e29b-41d4-a716-446655440000` |
| Query String | (none) |
| Body | (none) |

### Step 1: Build String to Sign

```
1750248000
GET
/open/v1/compute/instance/list

e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
```

> The 4th line (sorted query string) is empty.
> The 5th line is `SHA256("")`.

### Step 2: Compute Signature

```python
import hashlib, hmac

secret = "nb_secret_aabbccdd11223344556677889900aabb"
string_to_sign = "1750248000\nGET\n/open/v1/compute/instance/list\n\ne3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

signature = hmac.new(
    secret.encode(),
    string_to_sign.encode(),
    hashlib.sha256
).hexdigest()
```

### Step 3: Send the Request

```bash
TIMESTAMP="1750248000"
NONCE="550e8400-e29b-41d4-a716-446655440000"
API_KEY="nb_key_a1b2c3d4e5f6a7b8c9d0e1f2"
API_SECRET="nb_secret_aabbccdd11223344556677889900aabb"
METHOD="GET"
PATH="/open/v1/compute/instance/list"
BODY=""

# SHA-256 of empty body
BODY_HASH=$(printf "%s" "$BODY" | sha256sum | awk '{print $1}')

# Build string to sign
STRING_TO_SIGN="${TIMESTAMP}
${METHOD}
${PATH}

${BODY_HASH}"

# Compute HMAC-SHA256
SIGNATURE=$(printf "%s" "$STRING_TO_SIGN" | openssl dgst -sha256 -hmac "$API_SECRET" | awk '{print $NF}')

curl -s "https://api.neburst.com${PATH}" \
  -H "X-Nb-Key: ${API_KEY}" \
  -H "X-Nb-Timestamp: ${TIMESTAMP}" \
  -H "X-Nb-Nonce: ${NONCE}" \
  -H "X-Nb-Signature: ${SIGNATURE}"
```

## POST Request Example

Power-cycle a bare-metal instance:

```bash
TIMESTAMP=$(date +%s)
NONCE=$(uuidgen | tr '[:upper:]' '[:lower:]')
API_KEY="nb_key_a1b2c3d4e5f6a7b8c9d0e1f2"
API_SECRET="nb_secret_aabbccdd11223344556677889900aabb"
METHOD="POST"
PATH="/open/v1/compute/bare-metal/a1b2c3d4-e5f6-7890-abcd-ef1234567890/power"
BODY='{"action":"power-cycle"}'

BODY_HASH=$(printf "%s" "$BODY" | sha256sum | awk '{print $1}')

STRING_TO_SIGN="${TIMESTAMP}
${METHOD}
${PATH}

${BODY_HASH}"

SIGNATURE=$(printf "%s" "$STRING_TO_SIGN" | openssl dgst -sha256 -hmac "$API_SECRET" | awk '{print $NF}')

curl -s "https://api.neburst.com${PATH}" \
  -X POST \
  -H "Content-Type: application/json" \
  -H "X-Nb-Key: ${API_KEY}" \
  -H "X-Nb-Timestamp: ${TIMESTAMP}" \
  -H "X-Nb-Nonce: ${NONCE}" \
  -H "X-Nb-Signature: ${SIGNATURE}" \
  -d "$BODY"
```

## Security Notes

- **Clock synchronization** -- Keep your server's clock synced via NTP. Requests with timestamps more than 5 minutes from the server's time will be rejected with error code `58201`.
- **Secret storage** -- The API secret is shown **only once** when the key is created. Store it securely (e.g., environment variable, secrets manager). It cannot be retrieved later. If lost, use the Revoke (rotate) endpoint to generate a new secret.
- **Nonce uniqueness** -- Use UUID v4 for nonces. Reusing a nonce within the 10-minute window will be rejected with error code `58202`.
- **HTTPS only** -- All requests must be made over HTTPS. Plain HTTP requests will be rejected.

## Common Errors

| Code | Name | When |
|------|------|------|
| 58200 | `INVALID_SIGNATURE` | The computed signature does not match |
| 58201 | `EXPIRED_TIMESTAMP` | Timestamp is outside the +/- 5 minute window |
| 58202 | `DUPLICATE_NONCE` | Same nonce used within 10 minutes |
| 58205 | `KEY_DISABLED` | The API key has been disabled in the dashboard |
| 58206 | `KEY_EXPIRED` | The API key has passed its expiration date |
| 58207 | `KEY_NOT_FOUND` | The Key ID does not exist |
| 58208 | `MISSING_HEADERS` | One or more required headers are missing |
