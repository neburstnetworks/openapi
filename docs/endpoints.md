# Endpoint Reference

**[中文版](endpoints_zh.md)**

Base URL: `https://api.neburst.com/open/v1`

All responses use the standard envelope:

```json
{
  "code": 0,
  "msg": "",
  "data": ...
}
```

`code = 0` indicates success. See [Error Codes](errors.md) for non-zero values.

---

## Cloud Instance

Endpoints under `/compute/instance/`.

### List Cloud Instances

```
GET /compute/instance/list?page=1&page_size=20
```

**Scope:** `instance:read`

**Query parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | integer | `1` | Page number (1-based) |
| `page_size` | integer | `20` | Items per page |

**Response:**

```json
{
  "code": 0,
  "msg": "",
  "data": {
    "items": [
      {
        "uuid": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
        "name": "web-prod-01",
        "type": "cloud",
        "status": "running",
        "primary_ipv4": "203.0.113.10",
        "ipv4_list": ["203.0.113.10"],
        "ipv6_list": ["2001:db8::1"],
        "specs": {
          "cpu_model": "Intel Xeon",
          "cpu_cores": 4,
          "memory_gb": 8,
          "disks": [
            { "size_gb": 80, "type": "SSD" }
          ],
          "network_speed_gbps": 1
        },
        "os_name": "Ubuntu 24.04 LTS",
        "region": "Los Angeles",
        "hostname": "web-prod-01.neburst.net",
        "pay_cycle": "monthly",
        "auto_renew": true,
        "next_pay_at": "2026-07-18T00:00:00Z",
        "created_at": "2025-12-01T08:30:00Z"
      }
    ],
    "total": 42,
    "page": 1,
    "page_size": 20
  }
}
```

**Instance fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `uuid` | string | yes | Instance UUID |
| `name` | string | yes | User-defined name |
| `type` | string | yes | `"bare_metal"` or `"cloud"` |
| `status` | string | yes | Lifecycle status (see status enums below) |
| `primary_ipv4` | string | no | Primary IPv4 address |
| `ipv4_list` | string[] | no | All assigned IPv4 addresses |
| `ipv6_list` | string[] | no | All assigned IPv6 addresses |
| `specs` | object | no | Hardware specifications |
| `os_name` | string | no | Installed operating system |
| `region` | string | no | Datacenter region (display name) |
| `hostname` | string | no | Instance hostname |
| `pay_cycle` | string | no | Billing cycle (e.g., `monthly`, `quarterly`, `yearly`) |
| `auto_renew` | boolean | yes | Whether auto-renewal is enabled |
| `next_pay_at` | string | no | Next billing date (ISO 8601) |
| `created_at` | string | yes | Creation timestamp (ISO 8601) |

**Specs object:**

| Field | Type | Description |
|-------|------|-------------|
| `cpu_model` | string | CPU model name |
| `cpu_cores` | integer | Number of CPU cores |
| `memory_gb` | integer | Memory in GB |
| `disks` | object[] | List of disks (`size_gb`, `type`) |
| `network_speed_gbps` | float | Network port speed in Gbps |

**Pagination envelope:**

| Field | Type | Description |
|-------|------|-------------|
| `items` | array | List of items for the current page |
| `total` | integer | Total number of items across all pages |
| `page` | integer | Current page number |
| `page_size` | integer | Items per page |

**Cloud instance status values:** `running`, `stopped`, `starting`, `restarting`, `stopping`, `pending`, `provisioning`, `deleting`, `terminating`, `terminated`, `error`, `failed`, `maintenance`, `upgrading`

---

### Get Cloud Instance

```
GET /compute/instance/{id}
```

**Scope:** `instance:read`

**Path parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string (UUID) | Instance UUID |

**Response:**

```json
{
  "code": 0,
  "msg": "",
  "data": {
    "uuid": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "name": "web-prod-01",
    "type": "cloud",
    "status": "running",
    "primary_ipv4": "203.0.113.10",
    "ipv4_list": ["203.0.113.10"],
    "ipv6_list": ["2001:db8::1"],
    "specs": {
      "cpu_model": "Intel Xeon",
      "cpu_cores": 4,
      "memory_gb": 8,
      "disks": [
        { "size_gb": 80, "type": "SSD" }
      ],
      "network_speed_gbps": 1
    },
    "os_name": "Ubuntu 24.04 LTS",
    "region": "Los Angeles",
    "hostname": "web-prod-01.neburst.net",
    "pay_cycle": "monthly",
    "auto_renew": true,
    "next_pay_at": "2026-07-18T00:00:00Z",
    "created_at": "2025-12-01T08:30:00Z"
  }
}
```

Fields are the same as in the [List Cloud Instances](#list-cloud-instances) response.

---

### Get Cloud Instance Status

```
GET /compute/instance/{id}/status
```

**Scope:** `instance:read`

**Path parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string (UUID) | Instance UUID |

**Response:**

```json
{
  "code": 0,
  "msg": "",
  "data": {
    "status": "on",
    "is_installing": false
  }
}
```

**Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `status` | string | `"on"` or `"off"` |
| `is_installing` | boolean | `true` if an OS installation/rebuild is in progress |

---

### Get Cloud Instance Traffic

```
GET /compute/instance/{id}/traffic
```

**Scope:** `instance:read`

**Path parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string (UUID) | Instance UUID |

**Response:**

```json
{
  "code": 0,
  "msg": "",
  "data": {
    "packages": [
      {
        "name": "Monthly Traffic",
        "capacity_gb": 2000,
        "used_gb": 843.27,
        "reset_cycle": "monthly"
      }
    ]
  }
}
```

**TrafficPackageDTO fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | yes | Package name |
| `capacity_gb` | integer | yes | Total capacity in GB |
| `used_gb` | float | yes | Current usage in GB |
| `reset_cycle` | string | no | Reset period (e.g., `monthly`) |

---

### Power Action (Cloud)

```
POST /compute/instance/{id}/power
```

**Scope:** `instance:power`

**Path parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string (UUID) | Instance UUID |

**Request body:**

```json
{
  "action": "restart"
}
```

| Field | Type | Required | Values |
|-------|------|----------|--------|
| `action` | string | yes | `start`, `stop`, `restart` |

**Action descriptions:**

| Action | Behavior |
|--------|----------|
| `start` | Start the instance |
| `stop` | Stop the instance |
| `restart` | Restart the instance |

**Response:**

```json
{
  "code": 0,
  "msg": ""
}
```

---

### Cloud Instance Metrics

```
GET /compute/instance/{id}/metrics
```

**Scope:** `instance:read`

**Path parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string (UUID) | Instance UUID |

**Response:**

```json
{
  "code": 0,
  "msg": "",
  "data": {
    "cpu": {
      "percentage": 23.5
    },
    "memory": {
      "limit": 8192,
      "usage": 4096,
      "free": 4096,
      "percentage": 50.0,
      "unit": "MB"
    },
    "disk": {
      "limit": 80.0,
      "usage": 32.5,
      "free": 47.5,
      "percentage": 40.6,
      "unit": "GB"
    },
    "bandwidth": {
      "limit": 2000.0,
      "allowance": 2000.0,
      "usage": 843.2,
      "inbound": 120.3,
      "outbound": 843.2,
      "free": 1156.8,
      "percentage": 42.2,
      "usage_unit": "GB",
      "limit_unit": "GB",
      "started_time": "2026-06-01T00:00:00Z",
      "end_time": "2026-07-01T00:00:00Z"
    },
    "network": {
      "inbound": 125.5,
      "outbound": 98.3,
      "unit": "Mbps"
    }
  }
}
```

**Metrics fields:**

| Field | Type | Description |
|-------|------|-------------|
| `cpu.percentage` | float | CPU utilization (%) |
| `memory.limit` | float | Total memory |
| `memory.usage` | float | Used memory |
| `memory.free` | float | Free memory |
| `memory.percentage` | float | Memory utilization (%) |
| `memory.unit` | string | Unit (`MB`) |
| `disk.limit` | float | Total disk capacity |
| `disk.usage` | float | Used disk space |
| `disk.free` | float | Free disk space |
| `disk.percentage` | float | Disk utilization (%) |
| `disk.unit` | string | Unit (`GB`) |
| `bandwidth.limit` | float | Traffic limit for current period |
| `bandwidth.allowance` | float | Configured bandwidth allowance |
| `bandwidth.usage` | float | Outbound traffic used |
| `bandwidth.inbound` | float | Inbound traffic |
| `bandwidth.outbound` | float | Outbound traffic |
| `bandwidth.free` | float | Remaining traffic |
| `bandwidth.percentage` | float | Bandwidth utilization (%) |
| `bandwidth.usage_unit` | string | Traffic unit (`GB`) |
| `bandwidth.limit_unit` | string | Limit unit (`GB` or `TB`) |
| `bandwidth.started_time` | string | Period start (ISO 8601) |
| `bandwidth.end_time` | string | Period end (ISO 8601) |
| `network.inbound` | float | Current inbound speed |
| `network.outbound` | float | Current outbound speed |
| `network.unit` | string | Speed unit (`Mbps` or `Gbps`) |

---

## Bare Metal

Endpoints under `/compute/bare-metal/`.

### List Bare-Metal Instances

```
GET /compute/bare-metal/list?page=1&page_size=20
```

**Scope:** `bare-metal:read`

**Query parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | integer | `1` | Page number (1-based) |
| `page_size` | integer | `20` | Items per page |

**Response:**

```json
{
  "code": 0,
  "msg": "",
  "data": {
    "items": [
      {
        "uuid": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
        "name": "bm-prod-01",
        "type": "bare_metal",
        "status": "allocated",
        "primary_ipv4": "203.0.113.10",
        "ipv4_list": ["203.0.113.10"],
        "ipv6_list": ["2001:db8::1"],
        "specs": {
          "cpu_model": "Intel Xeon E-2388G",
          "cpu_cores": 8,
          "memory_gb": 32,
          "disks": [
            { "size_gb": 512, "type": "NVMe SSD" }
          ],
          "network_speed_gbps": 1
        },
        "os_name": "Ubuntu 24.04 LTS",
        "region": "Los Angeles",
        "hostname": "bm-prod-01.neburst.net",
        "pay_cycle": "monthly",
        "auto_renew": true,
        "next_pay_at": "2026-07-18T00:00:00Z",
        "created_at": "2025-12-01T08:30:00Z"
      }
    ],
    "total": 5,
    "page": 1,
    "page_size": 20
  }
}
```

Instance fields and pagination envelope are the same as in [List Cloud Instances](#list-cloud-instances).

**Bare-metal instance status values:** `unallocated`, `provisioning`, `allocated`

---

### Get Bare-Metal Instance

```
GET /compute/bare-metal/{id}
```

**Scope:** `bare-metal:read`

**Path parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string (UUID) | Instance UUID |

**Response:**

Same structure as [List Bare-Metal Instances](#list-bare-metal-instances) item.

---

### Get Bare-Metal Instance Status

```
GET /compute/bare-metal/{id}/status
```

**Scope:** `bare-metal:read`

**Path parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string (UUID) | Instance UUID |

**Response:**

```json
{
  "code": 0,
  "msg": "",
  "data": {
    "status": "on",
    "is_installing": false
  }
}
```

**Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `status` | string | `"on"` or `"off"` |
| `is_installing` | boolean | `true` if an OS installation/rebuild is in progress |

---

### Get Bare-Metal Instance Traffic

```
GET /compute/bare-metal/{id}/traffic
```

**Scope:** `bare-metal:read`

**Path parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string (UUID) | Instance UUID |

**Response:**

```json
{
  "code": 0,
  "msg": "",
  "data": {
    "packages": [
      {
        "name": "Monthly Traffic",
        "capacity_gb": 2000,
        "used_gb": 843.27,
        "reset_cycle": "monthly"
      }
    ]
  }
}
```

TrafficPackageDTO fields are the same as in [Get Cloud Instance Traffic](#get-cloud-instance-traffic).

---

### Power Action (Bare-Metal)

```
POST /compute/bare-metal/{id}/power
```

**Scope:** `bare-metal:power`

**Path parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string (UUID) | Instance UUID |

**Request body:**

```json
{
  "action": "power-cycle"
}
```

| Field | Type | Required | Values |
|-------|------|----------|--------|
| `action` | string | yes | `power-on`, `power-off`, `power-cycle`, `power-reset` |

**Action descriptions:**

| Action | Behavior |
|--------|----------|
| `power-on` | Start the instance |
| `power-off` | Graceful shutdown |
| `power-cycle` | Graceful reboot (shutdown then start) |
| `power-reset` | Hard reboot (immediate, equivalent to pulling the power plug) |

**Response:**

```json
{
  "code": 0,
  "msg": ""
}
```

---

### Bare-Metal Instance Metrics

```
GET /compute/bare-metal/{id}/metrics
```

**Scope:** `bare-metal:read`

**Path parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string (UUID) | Instance UUID |

**Response:**

Same structure as [Cloud Instance Metrics](#cloud-instance-metrics). Bare-metal metrics may have fewer fields populated (e.g., `cpu`, `memory`, and `disk` may be empty if the agent is not reporting).

---

### Rebuild Instance

```
POST /compute/bare-metal/{id}/rebuild
```

**Scope:** `bare-metal:rebuild`

> **Warning:** This operation is destructive. All data on the instance's disks will be erased.

**Path parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string (UUID) | Instance UUID |

**Request body:**

```json
{
  "profile_id": 123,
  "hostname": "bm-prod-01",
  "public_keys": [
    "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAI... user@host"
  ]
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `profile_id` | integer | yes | OS profile ID to install |
| `hostname` | string | no | New hostname for the instance |
| `public_keys` | string[] | no | SSH public keys to inject into the new OS |

**Response:**

```json
{
  "code": 0,
  "msg": ""
}
```

---

### Rescue Mode

```
POST /compute/bare-metal/{id}/rescue
```

**Scope:** `bare-metal:rescue`

Boot an instance into a rescue/recovery environment. The instance's original disks will be mounted, but the OS runs from a temporary rescue image.

**Path parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string (UUID) | Instance UUID |

**Request body:**

```json
{
  "profile_id": 456
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `profile_id` | integer | yes | Rescue OS profile ID |

**Response:**

```json
{
  "code": 0,
  "msg": ""
}
```

---

### List OS Profiles

```
GET /compute/bare-metal/{id}/profiles
```

**Scope:** `bare-metal:read`

List available OS profiles for rebuilding an instance.

**Path parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string (UUID) | Instance UUID |

**Response:**

```json
{
  "code": 0,
  "msg": "",
  "data": [
    {
      "id": 123,
      "name": "Ubuntu 24.04 LTS",
      "category": "linux",
      "is_rescue": false,
      "features": {
        "allow_ssh_keys": true,
        "allow_set_hostname": true
      }
    }
  ]
}
```

**OS Profile fields:**

| Field | Type | Description |
|-------|------|-------------|
| `id` | integer | Profile ID (used in rebuild/rescue requests) |
| `name` | string | Display name of the OS |
| `category` | string | Profile category (e.g., `linux`, `windows`) |
| `is_rescue` | boolean | Whether this is a rescue profile |
| `features` | object | Supported features for this profile |
| `features.allow_ssh_keys` | boolean | Whether SSH key injection is supported |
| `features.allow_set_hostname` | boolean | Whether custom hostname is supported |

---

### List Rescue Profiles

```
GET /compute/bare-metal/{id}/rescue-profiles
```

**Scope:** `bare-metal:read`

List available rescue/recovery profiles for an instance.

**Path parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string (UUID) | Instance UUID |

**Response:**

Same structure as [List OS Profiles](#list-os-profiles), but only returns profiles where `is_rescue` is `true`.

---

## Billing

### Get Balance

```
GET /billing/balance
```

**Scope:** `billing:read`

**Parameters:** None

**Response:**

```json
{
  "code": 0,
  "msg": "",
  "data": {
    "available": 128.50,
    "locked": 12.00,
    "currency": "USD"
  }
}
```

**Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `available` | float | Spendable balance in USD |
| `locked` | float | Balance reserved for pending orders |
| `currency` | string | Always `"USD"` |

---

### List Invoices

```
GET /billing/invoices?page=1&page_size=20
```

**Scope:** `billing:read`

**Query parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | integer | `1` | Page number (1-based) |
| `page_size` | integer | `20` | Items per page |

**Response:**

```json
{
  "code": 0,
  "msg": "",
  "data": {
    "items": [
      {
        "uuid": "f9e8d7c6-b5a4-3210-fedc-ba9876543210",
        "amount": 29.99,
        "status": "paid",
        "category": "compute",
        "created_at": "2026-06-01T00:00:00Z",
        "due_at": "2026-06-15T00:00:00Z"
      }
    ],
    "total": 15,
    "page": 1,
    "page_size": 20
  }
}
```

**Invoice fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `uuid` | string | yes | Invoice UUID |
| `amount` | float | yes | Amount in USD |
| `status` | string | yes | Invoice status (see enum below) |
| `category` | string | yes | Invoice category (see enum below) |
| `created_at` | string | yes | Creation timestamp (ISO 8601) |
| `due_at` | string | no | Payment due date (ISO 8601) |

**Invoice status enum:** `paid`, `partially-paid`, `unpaid`, `refunded`, `partially-refunded`, `refunding`, `cancelled`, `fallback-to-balance`

**Invoice category enum:** `compute`, `topup`, `traffic`, `penalty`, `custom`

---

### Get Invoice

```
GET /billing/invoices/{id}
```

**Scope:** `billing:read`

**Path parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string (UUID) | Invoice UUID |

**Response:**

```json
{
  "code": 0,
  "msg": "",
  "data": {
    "uuid": "f9e8d7c6-b5a4-3210-fedc-ba9876543210",
    "amount": 29.99,
    "status": "paid",
    "category": "compute",
    "created_at": "2026-06-01T00:00:00Z",
    "due_at": "2026-06-15T00:00:00Z"
  }
}
```

---

## User

### Get User Info

```
GET /user/info
```

**Scope:** `user:read`

**Parameters:** None

**Response:**

```json
{
  "code": 0,
  "msg": "",
  "data": {
    "uuid": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "email": "user@example.com",
    "nickname": "johndoe",
    "created_at": "2025-01-15T10:30:00Z"
  }
}
```

**Fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `uuid` | string | yes | User UUID |
| `email` | string | yes | User email address |
| `nickname` | string | no | Nickname |
| `created_at` | string | no | Account creation timestamp (ISO 8601) |

---

## Quick Reference

| Method | Endpoint | Scope | Description |
|--------|----------|-------|-------------|
| GET | `/compute/instance/list` | `instance:read` | List cloud instances (paginated) |
| GET | `/compute/instance/{id}` | `instance:read` | Get cloud instance details |
| GET | `/compute/instance/{id}/status` | `instance:read` | Get cloud instance power status |
| GET | `/compute/instance/{id}/traffic` | `instance:read` | Get cloud instance traffic usage |
| POST | `/compute/instance/{id}/power` | `instance:power` | Start/stop/restart (cloud) |
| GET | `/compute/instance/{id}/metrics` | `instance:read` | Cloud instance metrics |
| GET | `/compute/bare-metal/list` | `bare-metal:read` | List bare-metal instances (paginated) |
| GET | `/compute/bare-metal/{id}` | `bare-metal:read` | Get bare-metal instance details |
| GET | `/compute/bare-metal/{id}/status` | `bare-metal:read` | Get bare-metal instance power status |
| GET | `/compute/bare-metal/{id}/traffic` | `bare-metal:read` | Get bare-metal instance traffic usage |
| POST | `/compute/bare-metal/{id}/power` | `bare-metal:power` | Power on/off/cycle/reset (bare-metal) |
| GET | `/compute/bare-metal/{id}/metrics` | `bare-metal:read` | Bare-metal instance metrics |
| POST | `/compute/bare-metal/{id}/rebuild` | `bare-metal:rebuild` | Reinstall OS (destructive) |
| POST | `/compute/bare-metal/{id}/rescue` | `bare-metal:rescue` | Boot into rescue mode |
| GET | `/compute/bare-metal/{id}/profiles` | `bare-metal:read` | List OS profiles for rebuild |
| GET | `/compute/bare-metal/{id}/rescue-profiles` | `bare-metal:read` | List rescue profiles |
| GET | `/billing/balance` | `billing:read` | Get account balance |
| GET | `/billing/invoices` | `billing:read` | List invoices (paginated) |
| GET | `/billing/invoices/{id}` | `billing:read` | Get invoice details |
| GET | `/user/info` | `user:read` | Get user info |
