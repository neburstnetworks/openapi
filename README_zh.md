# Neburst OpenAPI

**[English](README.md)**

Neburst OpenAPI 提供编程接口，用于管理云实例、账单等资源。

## 认证方式

所有请求均通过 **HMAC-SHA256** 签名进行身份验证。在 [Neburst 控制面板](https://dash.neburst.com) 中创建 API Key 后，你会获得一个 **Base64 组合密钥**，其中同时包含 `key_id` 和 `secret`。所有 SDK 和 CLI 均支持直接使用组合密钥。

每个请求必须携带以下 4 个 Header：

| Header | 说明 |
|--------|------|
| `X-Nb-Key` | 你的 API Key ID（`nb_key_...`） |
| `X-Nb-Timestamp` | Unix 时间戳（秒） |
| `X-Nb-Nonce` | UUID v4（一次性使用） |
| `X-Nb-Signature` | HMAC-SHA256 十六进制签名 |

完整的签名算法、组合密钥格式和示例请参阅 [docs/authentication_zh.md](docs/authentication_zh.md)。

## 基础 URL

```
https://api.neburst.com/open/v1/
```

## 接口列表

### 云实例

| 方法 | 路径 | 权限范围 | 说明 |
|------|------|----------|------|
| GET | `/compute/instance/list` | `instance:read` | 获取云实例列表（分页） |
| GET | `/compute/instance/{id}` | `instance:read` | 获取云实例详情 |
| GET | `/compute/instance/{id}/status` | `instance:read` | 获取电源状态 |
| GET | `/compute/instance/{id}/traffic` | `instance:read` | 获取流量使用情况 |
| POST | `/compute/instance/{id}/power` | `instance:power` | 开机 / 关机 / 重启 |
| GET | `/compute/instance/{id}/metrics` | `instance:read` | 获取实例监控指标 |

### 独立服务器

| 方法 | 路径 | 权限范围 | 说明 |
|------|------|----------|------|
| GET | `/compute/bare-metal/list` | `bare-metal:read` | 获取独服列表（分页） |
| GET | `/compute/bare-metal/{id}` | `bare-metal:read` | 获取独服详情 |
| GET | `/compute/bare-metal/{id}/status` | `bare-metal:read` | 获取电源状态 |
| GET | `/compute/bare-metal/{id}/traffic` | `bare-metal:read` | 获取流量使用情况 |
| POST | `/compute/bare-metal/{id}/power` | `bare-metal:power` | 开机 / 关机 / 电源循环 / 重置 |
| GET | `/compute/bare-metal/{id}/metrics` | `bare-metal:read` | 获取实例监控指标 |
| POST | `/compute/bare-metal/{id}/rebuild` | `bare-metal:rebuild` | 重装操作系统 |
| POST | `/compute/bare-metal/{id}/rescue` | `bare-metal:rescue` | 进入救援模式 |
| GET | `/compute/bare-metal/{id}/profiles` | `bare-metal:read` | 获取可用操作系统列表 |
| GET | `/compute/bare-metal/{id}/rescue-profiles` | `bare-metal:read` | 获取救援系统列表 |

### 账单

| 方法 | 路径 | 权限范围 | 说明 |
|------|------|----------|------|
| GET | `/billing/balance` | `billing:read` | 获取账户余额 |
| GET | `/billing/invoices` | `billing:read` | 获取账单列表（分页） |
| GET | `/billing/invoices/{id}` | `billing:read` | 获取账单详情 |

### 用户

| 方法 | 路径 | 权限范围 | 说明 |
|------|------|----------|------|
| GET | `/user/info` | `user:read` | 获取用户信息 |

完整的请求与响应文档请参阅 [docs/endpoints_zh.md](docs/endpoints_zh.md)。

## 速率限制

每用户每分钟 60 次请求，单次请求间隔不低于 5 秒。详见 [docs/rate-limiting_zh.md](docs/rate-limiting_zh.md)。

## SDK

| 语言 | 路径 | 依赖 |
|------|------|------|
| **Go** | [sdk/go/](sdk/go/) | 仅标准库 |
| **Python** | [sdk/python/](sdk/python/) | `requests` |
| **Java** | [sdk/java/](sdk/java/) | `gson`、Java 11+ |
| **Rust** | [sdk/rust/](sdk/rust/) | `reqwest`、`hmac`、`sha2`、`serde` |

### 快速上手

**Go**
```go
import "github.com/neburstnetworks/openapi/sdk/go/neburst"

// 使用 Base64 组合密钥
client := neburst.NewClient("https://api.neburst.com", "eyJrZXlfaWQ...")
// 或分别传入 key_id + secret
client := neburst.NewClient("https://api.neburst.com", "nb_key_...", "nb_secret_...")

result, err := client.ListInstances()
for _, inst := range result.Items {
    fmt.Println(inst.UUID, inst.Name, inst.PrimaryIPv4)
}
```

**Python**
```python
from neburst import NeburstClient

client = NeburstClient("https://api.neburst.com", "eyJrZXlfaWQ...")
result = client.list_instances()
for inst in result.items:
    print(inst.uuid, inst.name, inst.primary_ipv4)
```

**Java**
```java
import com.neburst.api.NeburstClient;

var client = new NeburstClient("https://api.neburst.com", "eyJrZXlfaWQ...");
var instances = client.compute().listInstances();
```

**Rust**
```rust
use neburst_api::NeburstClient;

let client = NeburstClient::new("https://api.neburst.com", "eyJrZXlfaWQ...");
let instances = client.list_instances().await?;
```

## 文档

- [认证方式](docs/authentication_zh.md) — HMAC-SHA256 签名流程与组合密钥
- [接口文档](docs/endpoints_zh.md) — 完整 API 参考
- [权限范围](docs/scopes_zh.md) — 权限系统说明
- [速率限制](docs/rate-limiting_zh.md) — 请求频率限制
- [错误码](docs/errors_zh.md) — 错误码参考
- [OpenAPI 规范](openapi.yaml) — 机器可读的 API 规范
