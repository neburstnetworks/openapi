# Neburst API Rust SDK

Official Rust SDK for the [Neburst OpenAPI](https://neburst.com).

## Installation

Add to your `Cargo.toml`:

```toml
[dependencies]
neburst-api = "0.1.0"
tokio = { version = "1", features = ["full"] }
```

## Quick Start

```rust
use neburst_api::{NeburstClient, POWER_ON};

#[tokio::main]
async fn main() -> neburst_api::Result<()> {
    let client = NeburstClient::new(
        "https://api.neburst.com",
        "your-key-id",
        "your-api-secret",
    );

    // List all instances
    let instances = client.list_instances().await?;
    for inst in &instances {
        println!("{} - {}", inst.name, inst.status);
    }

    // Get instance details
    let instance = client.get_instance("instance-uuid").await?;

    // Get power status
    let status = client.get_instance_status("instance-uuid").await?;
    println!("Power: {}", status.status);

    // Power actions
    client.power_action("instance-uuid", POWER_ON).await?;

    // Rebuild instance
    client
        .rebuild_instance("instance-uuid", 1, Some("my-host"), None)
        .await?;

    // Rescue mode
    client.rescue_instance("instance-uuid", 1).await?;

    // Get traffic
    let traffic = client.get_instance_traffic("instance-uuid").await?;

    // Billing operations
    let balance = client.get_balance().await?;
    println!("Balance: {:.2} {}", balance.available, balance.currency);

    let invoices = client.list_invoices().await?;
    let invoice = client.get_invoice("invoice-uuid").await?;

    Ok(())
}
```

## Error Handling

```rust
use neburst_api::{NeburstClient, NeburstError};

match client.get_instance("invalid-id").await {
    Ok(inst) => println!("Found: {}", inst.name),
    Err(NeburstError::Api { code, message }) => {
        eprintln!("API error (code={}): {}", code, message);
    }
    Err(NeburstError::Http(e)) => eprintln!("HTTP error: {}", e),
    Err(NeburstError::Json(e)) => eprintln!("JSON error: {}", e),
}
```

## Authentication

All requests are signed with HMAC-SHA256. The signing process:

1. Construct the string to sign:
   ```
   StringToSign = timestamp + "\n" + method + "\n" + path + "\n" + sorted_query + "\n" + SHA256(body)
   ```
2. Compute `HMAC-SHA256(api_secret, StringToSign)`
3. Send the hex-encoded signature in the `X-Nb-Signature` header

Required headers on every request:
- `X-Nb-Key` - Your API key ID
- `X-Nb-Timestamp` - Current Unix timestamp (seconds)
- `X-Nb-Nonce` - Unique UUID v4 per request
- `X-Nb-Signature` - HMAC-SHA256 signature (hex)

## License

MIT
