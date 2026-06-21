# 错误码

**[English](errors.md)**

所有 API 响应均使用标准信封格式：

```json
{
  "code": <int>,
  "msg": "<string>",
  "data": ...
}
```

当 `code` 为 `0` 时，表示请求成功。非零值表示发生错误——`msg` 字段包含可读的错误描述。

## 错误码参考

| 错误码 | 名称 | HTTP 状态码 | 描述 |
|--------|------|-------------|------|
| 58200 | `INVALID_SIGNATURE` | 401 | `X-Nb-Signature` 请求头与服务器计算的签名不匹配。请检查签名算法、密钥和待签名字符串的构造。 |
| 58201 | `EXPIRED_TIMESTAMP` | 401 | `X-Nb-Timestamp` 的值与服务器当前时间相差超过 5 分钟。请确保系统时钟已通过 NTP 同步。 |
| 58202 | `DUPLICATE_NONCE` | 401 | `X-Nb-Nonce` 的值在最近 10 分钟内已被使用过。每个请求必须使用唯一的随机数（建议使用 UUID v4）。 |
| 58203 | `INSUFFICIENT_SCOPE` | 403 | API 密钥不具有该端点所需的权限范围。请在控制面板中检查密钥已分配的 Scope。 |
| 58204 | `RATE_LIMITED` | 429 | 已超过每用户每分钟 60 次请求的速率限制。请查看 `X-RateLimit-Reset` 响应头以了解窗口重置时间。 |
| 58205 | `KEY_DISABLED` | 401 | API 密钥已被手动禁用。请在控制面板中重新启用，或创建新密钥。 |
| 58206 | `KEY_EXPIRED` | 401 | API 密钥已超过其设定的过期时间。请创建新密钥。 |
| 58207 | `KEY_NOT_FOUND` | 401 | `X-Nb-Key` 的值与任何现有 API 密钥均不匹配。请确认 Key ID 是否正确。 |
| 58208 | `MISSING_HEADERS` | 400 | 请求中缺少一个或多个必需的认证请求头（`X-Nb-Key`、`X-Nb-Timestamp`、`X-Nb-Nonce`、`X-Nb-Signature`）。 |

## 业务逻辑错误

除上述认证/授权错误外，各端点还可能返回应用层错误：

| 错误码 | HTTP 状态码 | 描述 |
|--------|-------------|------|
| 58210 | 200 | 发生业务逻辑错误（例如实例未找到、电源状态无效）。`msg` 字段包含具体错误信息。 |
| 58400 | 200 | 无法解析请求体。请确保 JSON 格式正确且字段类型无误。 |

## 错误响应示例

### 认证失败

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

### 签名无效

```json
{
  "code": 58200,
  "msg": "Signature verification failed"
}
```

### 权限不足

尝试使用仅拥有 `instance:read` 权限的密钥调用 `POST /compute/instance/{id}/power`：

```json
{
  "code": 58203,
  "msg": "Insufficient scope"
}
```

### 频率限制

```json
{
  "code": 58204,
  "msg": "Rate limit exceeded"
}
```

触发频率限制时的响应头：

```
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 0
X-RateLimit-Reset: 1750291260
```

### 请求体无效

```json
{
  "code": 58400,
  "msg": "invalid request body"
}
```

## 故障排查

| 现象 | 可能原因 | 解决方法 |
|------|----------|----------|
| 始终收到 `58200` | 签名算法不正确 | 逐一检查待签名字符串的各个组成部分。参见 [认证](authentication.md)。 |
| 每个请求都返回 `58201` | 时钟偏差 | 运行 `ntpdate` 或启用 NTP 时间同步。 |
| 间歇性出现 `58202` | 随机数冲突 | 使用 UUID v4；不要在不同请求间复用随机数。 |
| 某些端点返回 `58203` | 缺少 Scope 权限 | 通过控制面板为密钥添加所需的 Scope。参见 [Scopes](scopes.md)。 |
| 密钥轮换后出现 `58207` | 使用了旧的 Key ID | 撤销/轮换后仅密钥密文会变更。如果是删除后重新创建，Key ID 也会改变。 |
