# Neburst OpenAPI

**[中文版](README_zh.md)**

Neburst OpenAPI provides programmatic access to manage compute instances, billing, and more.

## Authentication

All requests are authenticated using **HMAC-SHA256** signatures. Create an API Key in the [Neburst Dashboard](https://dash.neburst.com) — you will receive a **Base64 combined key** that contains both `key_id` and `secret`. All SDKs and the CLI accept this combined key directly.

Every request must include 4 headers:

| Header | Description |
|--------|-------------|
| `X-Nb-Key` | Your API Key ID (`nb_key_...`) |
| `X-Nb-Timestamp` | Unix epoch seconds |
| `X-Nb-Nonce` | UUID v4 (single-use) |
| `X-Nb-Signature` | HMAC-SHA256 hex signature |

See [docs/authentication.md](docs/authentication.md) for the full signing algorithm, combined key format, and examples.

## Base URL

```
https://api.neburst.com/open/v1/
```

## Endpoints

### Cloud Instance

| Method | Path | Scope | Description |
|--------|------|-------|-------------|
| GET | `/compute/instance/list` | `instance:read` | List cloud instances (paginated) |
| GET | `/compute/instance/{id}` | `instance:read` | Get cloud instance details |
| GET | `/compute/instance/{id}/status` | `instance:read` | Get power status |
| GET | `/compute/instance/{id}/traffic` | `instance:read` | Get traffic usage |
| POST | `/compute/instance/{id}/power` | `instance:power` | Start/stop/restart |
| GET | `/compute/instance/{id}/metrics` | `instance:read` | Instance metrics |

### Bare Metal

| Method | Path | Scope | Description |
|--------|------|-------|-------------|
| GET | `/compute/bare-metal/list` | `bare-metal:read` | List bare-metal instances (paginated) |
| GET | `/compute/bare-metal/{id}` | `bare-metal:read` | Get bare-metal instance details |
| GET | `/compute/bare-metal/{id}/status` | `bare-metal:read` | Get power status |
| GET | `/compute/bare-metal/{id}/traffic` | `bare-metal:read` | Get traffic usage |
| POST | `/compute/bare-metal/{id}/power` | `bare-metal:power` | Power on/off/cycle/reset |
| GET | `/compute/bare-metal/{id}/metrics` | `bare-metal:read` | Instance metrics |
| POST | `/compute/bare-metal/{id}/rebuild` | `bare-metal:rebuild` | Reinstall OS |
| POST | `/compute/bare-metal/{id}/rescue` | `bare-metal:rescue` | Enter rescue mode |
| GET | `/compute/bare-metal/{id}/profiles` | `bare-metal:read` | List OS profiles |
| GET | `/compute/bare-metal/{id}/rescue-profiles` | `bare-metal:read` | List rescue profiles |

### Billing

| Method | Path | Scope | Description |
|--------|------|-------|-------------|
| GET | `/billing/balance` | `billing:read` | Get account balance |
| GET | `/billing/invoices` | `billing:read` | List invoices (paginated) |
| GET | `/billing/invoices/{id}` | `billing:read` | Get invoice details |

### User

| Method | Path | Scope | Description |
|--------|------|-------|-------------|
| GET | `/user/info` | `user:read` | Get user info |

See [docs/endpoints.md](docs/endpoints.md) for full request/response documentation.

## Rate Limiting

60 requests per minute per user, with a minimum interval of 5 seconds between requests. See [docs/rate-limiting.md](docs/rate-limiting.md).

## SDKs

| Language | Path | Dependencies |
|----------|------|--------------|
| **Go** | [sdk/go/](sdk/go/) | Standard library only |
| **Python** | [sdk/python/](sdk/python/) | `requests` |
| **Java** | [sdk/java/](sdk/java/) | `gson`, Java 11+ |
| **Rust** | [sdk/rust/](sdk/rust/) | `reqwest`, `hmac`, `sha2`, `serde` |

### Quick Start

**Go**
```go
import "github.com/neburstnetworks/openapi/sdk/go/neburst"

// Combined Base64 key (single argument)
client := neburst.NewClient("https://api.neburst.com", "eyJrZXlfaWQ...")
// Or separate key_id + secret
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

## Documentation

- [Authentication](docs/authentication.md) — HMAC-SHA256 signing flow & combined key
- [Endpoints](docs/endpoints.md) — Full API reference
- [Scopes](docs/scopes.md) — Permission system
- [Rate Limiting](docs/rate-limiting.md) — Request limits
- [Errors](docs/errors.md) — Error codes reference
- [OpenAPI Spec](openapi.yaml) — Machine-readable spec
