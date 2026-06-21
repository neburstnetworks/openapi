# Authentication

**[English](authentication.md)**

Neburst OpenAPI 采用 **HMAC-SHA256** 请求签名机制。每个请求都必须携带四个自定义请求头。系统不使用 Bearer Token 或会话 Cookie，每个请求都会被独立验证。

## Required Headers

| Header | 格式 | 说明 |
|--------|------|------|
| `X-Nb-Key` | `nb_key_` + 24 位十六进制字符 | 你的 API Key ID |
| `X-Nb-Timestamp` | Unix 时间戳（秒） | 当前时间；必须与服务器时钟相差不超过 **5 分钟** |
| `X-Nb-Nonce` | UUID v4 | 每个请求唯一；服务器会在 **10 分钟**内拒绝重复的 Nonce |
| `X-Nb-Signature` | 64 位十六进制字符 | `Hex(HMAC-SHA256(api_secret, string_to_sign))` |

## Signing Algorithm

### Step 1 -- Build the String to Sign

将以下五个部分用换行符（`\n`）拼接：

```
StringToSign = timestamp + "\n"
             + method    + "\n"
             + path      + "\n"
             + sorted_query_string + "\n"
             + SHA256(body)
```

| 组成部分 | 规则 |
|----------|------|
| **timestamp** | 与 `X-Nb-Timestamp` 请求头中发送的值相同 |
| **method** | 大写的 HTTP 方法：`GET`、`POST`、`PUT`、`DELETE` |
| **path** | 以 `/open/v1/...` 开头的请求路径（不包含协议、域名或查询字符串） |
| **sorted_query_string** | 查询参数按键名字母顺序排列。每个键和值都需进行 URL 编码，键值对之间用 `&` 连接。无查询参数时为空字符串。 |
| **SHA256(body)** | 原始请求体的 SHA-256 哈希值（十六进制编码）。对于没有请求体的请求（GET、DELETE），对空字符串取哈希：`e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |

### Step 2 -- Compute the Signature

```
Signature = Hex(HMAC-SHA256(api_secret, StringToSign))
```

- `api_secret` 是创建 API Key 时返回的密钥（以 `nb_secret_` 开头）。
- 输出为小写的十六进制字符串（64 个字符）。

### Step 3 -- Send the Request

将四个请求头附加到你的 HTTP 请求中：

```
X-Nb-Key: nb_key_1a2b3c4d5e6f7a8b9c0d1e2f
X-Nb-Timestamp: 1750291200
X-Nb-Nonce: 550e8400-e29b-41d4-a716-446655440000
X-Nb-Signature: 7f83b1657ff1fc53b92dc18148a1d65dfc2d4b1fa3d677284addd200126d9069
```

## Worked Example

假设你想在 `2026-06-18T12:00:00Z`（Unix 时间戳 `1750248000`）列出云实例。

**请求参数：**

| 字段 | 值 |
|------|-----|
| Method | `GET` |
| URL | `https://api.neburst.com/open/v1/compute/instance/list` |
| API Key ID | `nb_key_a1b2c3d4e5f6a7b8c9d0e1f2` |
| API Secret | `nb_secret_aabbccdd11223344556677889900aabb` |
| Timestamp | `1750248000` |
| Nonce | `550e8400-e29b-41d4-a716-446655440000` |
| Query String | （无） |
| Body | （无） |

### Step 1: Build String to Sign

```
1750248000
GET
/open/v1/compute/instance/list

e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
```

> 第 4 行（排序后的查询字符串）为空。
> 第 5 行是 `SHA256("")` 的结果。

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

# 空请求体的 SHA-256
BODY_HASH=$(printf "%s" "$BODY" | sha256sum | awk '{print $1}')

# 构建待签名字符串
STRING_TO_SIGN="${TIMESTAMP}
${METHOD}
${PATH}

${BODY_HASH}"

# 计算 HMAC-SHA256
SIGNATURE=$(printf "%s" "$STRING_TO_SIGN" | openssl dgst -sha256 -hmac "$API_SECRET" | awk '{print $NF}')

curl -s "https://api.neburst.com${PATH}" \
  -H "X-Nb-Key: ${API_KEY}" \
  -H "X-Nb-Timestamp: ${TIMESTAMP}" \
  -H "X-Nb-Nonce: ${NONCE}" \
  -H "X-Nb-Signature: ${SIGNATURE}"
```

## POST Request Example

对裸金属实例执行电源重启：

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

## 组合密钥（Base64）

在 [Neburst 控制面板](https://dash.neburst.com) 创建 API 密钥时，系统会提供一个**组合密钥** —— 一个 Base64 编码的字符串，同时包含 Key ID 和 Secret：

```
eyJrZXlfaWQiOiJuYl9rZXlfYTFiMmMzZDRlNWY2YTdiOGM5ZDBlMWYyIiwic2VjcmV0IjoibmJfc2VjcmV0X2FhYmJjY2RkMTEyMjMzNDQ1NTY2Nzc4ODk5MDBhYWJiIn0=
```

解码后得到 JSON 对象：

```json
{
  "key_id": "nb_key_a1b2c3d4e5f6a7b8c9d0e1f2",
  "secret": "nb_secret_aabbccdd11223344556677889900aabb"
}
```

所有官方 SDK 均支持直接传入组合密钥，无需分别提供 `key_id` 和 `secret`：

```go
// 组合密钥（单参数）
client := neburst.NewClient("https://api.neburst.com", "eyJrZXlfaWQ...")

// 等价于分别传入 key_id + secret
client := neburst.NewClient("https://api.neburst.com", "nb_key_...", "nb_secret_...")
```

CLI 添加账户时同样支持组合密钥：

```bash
neburst-cli account add my-account
# 在 "API Key" 提示处粘贴组合密钥即可
```

## Security Notes

- **时钟同步** -- 请通过 NTP 保持服务器时钟同步。时间戳与服务器时间相差超过 5 分钟的请求将被拒绝，错误码为 `58201`。
- **密钥存储** -- API Secret 仅在创建密钥时显示**一次**。请妥善保管（例如使用环境变量或密钥管理服务）。密钥无法再次获取。如果遗失，请使用 Revoke（轮换）接口生成新密钥。
- **Nonce 唯一性** -- 请使用 UUID v4 作为 Nonce。在 10 分钟窗口内重复使用同一 Nonce 将被拒绝，错误码为 `58202`。
- **仅限 HTTPS** -- 所有请求必须通过 HTTPS 发送。纯 HTTP 请求将被拒绝。

## Common Errors

| Code | Name | 触发条件 |
|------|------|----------|
| 58200 | `INVALID_SIGNATURE` | 计算的签名与服务端不匹配 |
| 58201 | `EXPIRED_TIMESTAMP` | 时间戳超出 +/- 5 分钟窗口 |
| 58202 | `DUPLICATE_NONCE` | 10 分钟内使用了重复的 Nonce |
| 58205 | `KEY_DISABLED` | API Key 已在控制面板中被禁用 |
| 58206 | `KEY_EXPIRED` | API Key 已过期 |
| 58207 | `KEY_NOT_FOUND` | Key ID 不存在 |
| 58208 | `MISSING_HEADERS` | 缺少一个或多个必需的请求头 |
