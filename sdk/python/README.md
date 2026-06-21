# Neburst Python SDK

Python SDK for the Neburst OpenAPI. Uses HMAC-SHA256 request signing.

## Installation

```bash
pip install neburst
```

Or install from source:

```bash
pip install -e /path/to/neburst-openapi/sdk/python
```

## Quick Start

```python
from neburst import NeburstClient

client = NeburstClient(
    base_url="https://api.neburst.com",
    key_id="your-key-id",
    secret="your-api-secret",
)

# List all instances
instances = client.list_instances()
for inst in instances:
    print(f"{inst.name} ({inst.type}) - {inst.status}")
```

## Authentication

All requests are signed with HMAC-SHA256. You need an API key pair (key ID + secret) which can be generated from the Neburst dashboard.

The SDK handles signing automatically — just pass your `key_id` and `secret` when creating the client.

## API Reference

### Compute

| Method | Description |
|---|---|
| `list_instances()` | List all compute instances |
| `get_instance(id)` | Get instance details |
| `get_instance_status(id)` | Get power status |
| `get_instance_traffic(id)` | Get traffic usage |
| `power_action(id, action)` | Power on/off/cycle/reset |
| `rebuild_instance(id, profile_id, hostname?, public_keys?)` | Reinstall OS |
| `rescue_instance(id, profile_id)` | Boot into rescue mode |

**Power actions:** `"power-on"`, `"power-off"`, `"power-cycle"`, `"power-reset"`

### Billing

| Method | Description |
|---|---|
| `get_balance()` | Get account balance |
| `list_invoices()` | List all invoices |
| `get_invoice(id)` | Get invoice details |

## Error Handling

```python
from neburst import NeburstClient, APIError

client = NeburstClient(base_url="https://api.neburst.com", key_id="...", secret="...")

try:
    instance = client.get_instance("non-existent-uuid")
except APIError as e:
    print(f"API error {e.code}: {e.msg}")
```

## Data Types

- **Instance** — uuid, name, type, status, region, hostname, pay_cycle, auto_renew, next_pay_at, created_at
- **PowerStatus** — status, is_installing
- **Traffic** — packages (list of TrafficPackage)
- **TrafficPackage** — name, capacity_gb, used_gb, reset_cycle
- **Balance** — available, locked, currency
- **Invoice** — uuid, amount, status, category, created_at, due_at

## Requirements

- Python >= 3.8
- requests >= 2.20.0
