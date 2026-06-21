# 接口参考

**[English](endpoints.md)**

Base URL: `https://api.neburst.com/open/v1`

所有响应均使用统一的信封格式：

```json
{
  "code": 0,
  "msg": "",
  "data": ...
}
```

`code = 0` 表示成功。非零值请参阅 [错误码](errors.md)。

---

## 云实例

接口路径前缀：`/compute/instance/`

### 获取云实例列表

```
GET /compute/instance/list?page=1&page_size=20
```

**所需权限：** `instance:read`

**查询参数：**

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `page` | integer | `1` | 页码（从 1 开始） |
| `page_size` | integer | `20` | 每页条数 |

**响应：**

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

**实例字段：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `uuid` | string | 是 | 实例 UUID |
| `name` | string | 是 | 用户自定义名称 |
| `type` | string | 是 | `"bare_metal"` 或 `"cloud"` |
| `status` | string | 是 | 生命周期状态（见下方状态枚举） |
| `primary_ipv4` | string | 否 | 主 IPv4 地址 |
| `ipv4_list` | string[] | 否 | 所有已分配的 IPv4 地址 |
| `ipv6_list` | string[] | 否 | 所有已分配的 IPv6 地址 |
| `specs` | object | 否 | 硬件规格 |
| `os_name` | string | 否 | 已安装的操作系统 |
| `region` | string | 否 | 数据中心区域（显示名称） |
| `hostname` | string | 否 | 实例主机名 |
| `pay_cycle` | string | 否 | 计费周期（如 `monthly`、`quarterly`、`yearly`） |
| `auto_renew` | boolean | 是 | 是否启用自动续费 |
| `next_pay_at` | string | 否 | 下次计费日期（ISO 8601） |
| `created_at` | string | 是 | 创建时间（ISO 8601） |

**Specs 对象：**

| 字段 | 类型 | 说明 |
|------|------|------|
| `cpu_model` | string | CPU 型号 |
| `cpu_cores` | integer | CPU 核心数 |
| `memory_gb` | integer | 内存大小（GB） |
| `disks` | object[] | 磁盘列表（`size_gb`、`type`） |
| `network_speed_gbps` | float | 网络端口速率（Gbps） |

**分页信封：**

| 字段 | 类型 | 说明 |
|------|------|------|
| `items` | array | 当前页的数据列表 |
| `total` | integer | 所有页的总条数 |
| `page` | integer | 当前页码 |
| `page_size` | integer | 每页条数 |

**云实例状态枚举：** `running`、`stopped`、`starting`、`restarting`、`stopping`、`pending`、`provisioning`、`deleting`、`terminating`、`terminated`、`error`、`failed`、`maintenance`、`upgrading`

---

### 获取云实例详情

```
GET /compute/instance/{id}
```

**所需权限：** `instance:read`

**路径参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| `id` | string (UUID) | 实例 UUID |

**响应：**

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

字段与 [获取云实例列表](#获取云实例列表) 响应一致。

---

### 获取云实例电源状态

```
GET /compute/instance/{id}/status
```

**所需权限：** `instance:read`

**路径参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| `id` | string (UUID) | 实例 UUID |

**响应：**

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

**字段：**

| 字段 | 类型 | 说明 |
|------|------|------|
| `status` | string | `"on"` 或 `"off"` |
| `is_installing` | boolean | 若正在安装/重装操作系统则为 `true` |

---

### 获取云实例流量

```
GET /compute/instance/{id}/traffic
```

**所需权限：** `instance:read`

**路径参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| `id` | string (UUID) | 实例 UUID |

**响应：**

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

**TrafficPackageDTO 字段：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 是 | 流量包名称 |
| `capacity_gb` | integer | 是 | 总容量（GB） |
| `used_gb` | float | 是 | 已使用量（GB） |
| `reset_cycle` | string | 否 | 重置周期（如 `monthly`） |

---

### 电源操作（云实例）

```
POST /compute/instance/{id}/power
```

**所需权限：** `instance:power`

**路径参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| `id` | string (UUID) | 实例 UUID |

**请求体：**

```json
{
  "action": "restart"
}
```

| 字段 | 类型 | 必填 | 可选值 |
|------|------|------|--------|
| `action` | string | 是 | `start`、`stop`、`restart` |

**操作说明：**

| 操作 | 行为 |
|------|------|
| `start` | 启动实例 |
| `stop` | 停止实例 |
| `restart` | 重启实例 |

**响应：**

```json
{
  "code": 0,
  "msg": ""
}
```

---

### 云实例监控指标

```
GET /compute/instance/{id}/metrics
```

**所需权限：** `instance:read`

**路径参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| `id` | string (UUID) | 实例 UUID |

**响应：**

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

**Metrics 字段：**

| 字段 | 类型 | 说明 |
|------|------|------|
| `cpu.percentage` | float | CPU 利用率 (%) |
| `memory.limit` | float | 总内存 |
| `memory.usage` | float | 已用内存 |
| `memory.free` | float | 可用内存 |
| `memory.percentage` | float | 内存利用率 (%) |
| `memory.unit` | string | 单位 (`MB`) |
| `disk.limit` | float | 磁盘总容量 |
| `disk.usage` | float | 已用磁盘 |
| `disk.free` | float | 可用磁盘 |
| `disk.percentage` | float | 磁盘利用率 (%) |
| `disk.unit` | string | 单位 (`GB`) |
| `bandwidth.limit` | float | 当前周期流量限额 |
| `bandwidth.allowance` | float | 配置的带宽配额 |
| `bandwidth.usage` | float | 已用出站流量 |
| `bandwidth.inbound` | float | 入站流量 |
| `bandwidth.outbound` | float | 出站流量 |
| `bandwidth.free` | float | 剩余流量 |
| `bandwidth.percentage` | float | 带宽利用率 (%) |
| `bandwidth.usage_unit` | string | 流量单位 (`GB`) |
| `bandwidth.limit_unit` | string | 限额单位 (`GB` 或 `TB`) |
| `bandwidth.started_time` | string | 周期开始时间 (ISO 8601) |
| `bandwidth.end_time` | string | 周期结束时间 (ISO 8601) |
| `network.inbound` | float | 当前入站速率 |
| `network.outbound` | float | 当前出站速率 |
| `network.unit` | string | 速率单位 (`Mbps` 或 `Gbps`) |

---

## 独立服务器

接口路径前缀：`/compute/bare-metal/`

### 获取独立服务器列表

```
GET /compute/bare-metal/list?page=1&page_size=20
```

**所需权限：** `bare-metal:read`

**查询参数：**

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `page` | integer | `1` | 页码（从 1 开始） |
| `page_size` | integer | `20` | 每页条数 |

**响应：**

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

实例字段与分页信封与 [获取云实例列表](#获取云实例列表) 一致。

**独立服务器状态枚举：** `unallocated`、`provisioning`、`allocated`

---

### 获取独立服务器详情

```
GET /compute/bare-metal/{id}
```

**所需权限：** `bare-metal:read`

**路径参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| `id` | string (UUID) | 实例 UUID |

**响应：**

结构与 [获取独立服务器列表](#获取独立服务器列表) 中的单条记录一致。

---

### 获取独立服务器电源状态

```
GET /compute/bare-metal/{id}/status
```

**所需权限：** `bare-metal:read`

**路径参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| `id` | string (UUID) | 实例 UUID |

**响应：**

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

**字段：**

| 字段 | 类型 | 说明 |
|------|------|------|
| `status` | string | `"on"` 或 `"off"` |
| `is_installing` | boolean | 若正在安装/重装操作系统则为 `true` |

---

### 获取独立服务器流量

```
GET /compute/bare-metal/{id}/traffic
```

**所需权限：** `bare-metal:read`

**路径参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| `id` | string (UUID) | 实例 UUID |

**响应：**

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

TrafficPackageDTO 字段与 [获取云实例流量](#获取云实例流量) 一致。

---

### 电源操作（独立服务器）

```
POST /compute/bare-metal/{id}/power
```

**所需权限：** `bare-metal:power`

**路径参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| `id` | string (UUID) | 实例 UUID |

**请求体：**

```json
{
  "action": "power-cycle"
}
```

| 字段 | 类型 | 必填 | 可选值 |
|------|------|------|--------|
| `action` | string | 是 | `power-on`、`power-off`、`power-cycle`、`power-reset` |

**操作说明：**

| 操作 | 行为 |
|------|------|
| `power-on` | 开机 |
| `power-off` | 优雅关机 |
| `power-cycle` | 优雅重启（先关机再开机） |
| `power-reset` | 硬重启（立即执行，等同于拔电源） |

**响应：**

```json
{
  "code": 0,
  "msg": ""
}
```

---

### 独立服务器监控指标

```
GET /compute/bare-metal/{id}/metrics
```

**所需权限：** `bare-metal:read`

**路径参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| `id` | string (UUID) | 实例 UUID |

**响应：**

与[云实例监控指标](#云实例监控指标)结构相同。独立服务器的部分字段可能为空（如 `cpu`、`memory`、`disk`，取决于 Agent 是否上报数据）。

---

### 重装系统

```
POST /compute/bare-metal/{id}/rebuild
```

**所需权限：** `bare-metal:rebuild`

> **警告：** 此操作具有破坏性，实例磁盘上的所有数据将被清除。

**路径参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| `id` | string (UUID) | 实例 UUID |

**请求体：**

```json
{
  "profile_id": 123,
  "hostname": "bm-prod-01",
  "public_keys": [
    "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAI... user@host"
  ]
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `profile_id` | integer | 是 | 要安装的操作系统配置 ID |
| `hostname` | string | 否 | 新的主机名 |
| `public_keys` | string[] | 否 | 注入到新系统的 SSH 公钥 |

**响应：**

```json
{
  "code": 0,
  "msg": ""
}
```

---

### 救援模式

```
POST /compute/bare-metal/{id}/rescue
```

**所需权限：** `bare-metal:rescue`

将实例引导至救援/恢复环境。实例的原始磁盘会被挂载，但操作系统运行在临时的救援镜像上。

**路径参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| `id` | string (UUID) | 实例 UUID |

**请求体：**

```json
{
  "profile_id": 456
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `profile_id` | integer | 是 | 救援系统配置 ID |

**响应：**

```json
{
  "code": 0,
  "msg": ""
}
```

---

### 获取操作系统配置列表

```
GET /compute/bare-metal/{id}/profiles
```

**所需权限：** `bare-metal:read`

列出可用于重装实例的操作系统配置。

**路径参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| `id` | string (UUID) | 实例 UUID |

**响应：**

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

**操作系统配置字段：**

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | integer | 配置 ID（用于重装/救援请求） |
| `name` | string | 操作系统显示名称 |
| `category` | string | 配置分类（如 `linux`、`windows`） |
| `is_rescue` | boolean | 是否为救援配置 |
| `features` | object | 该配置支持的功能 |
| `features.allow_ssh_keys` | boolean | 是否支持注入 SSH 公钥 |
| `features.allow_set_hostname` | boolean | 是否支持自定义主机名 |

---

### 获取救援配置列表

```
GET /compute/bare-metal/{id}/rescue-profiles
```

**所需权限：** `bare-metal:read`

列出实例可用的救援/恢复配置。

**路径参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| `id` | string (UUID) | 实例 UUID |

**响应：**

结构与 [获取操作系统配置列表](#获取操作系统配置列表) 一致，但仅返回 `is_rescue` 为 `true` 的配置。

---

## 账单

### 获取余额

```
GET /billing/balance
```

**所需权限：** `billing:read`

**参数：** 无

**响应：**

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

**字段：**

| 字段 | 类型 | 说明 |
|------|------|------|
| `available` | float | 可用余额（USD） |
| `locked` | float | 待处理订单冻结的余额 |
| `currency` | string | 固定为 `"USD"` |

---

### 获取账单列表

```
GET /billing/invoices?page=1&page_size=20
```

**所需权限：** `billing:read`

**查询参数：**

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `page` | integer | `1` | 页码（从 1 开始） |
| `page_size` | integer | `20` | 每页条数 |

**响应：**

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

**账单字段：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `uuid` | string | 是 | 账单 UUID |
| `amount` | float | 是 | 金额（USD） |
| `status` | string | 是 | 账单状态（见下方枚举） |
| `category` | string | 是 | 账单分类（见下方枚举） |
| `created_at` | string | 是 | 创建时间（ISO 8601） |
| `due_at` | string | 否 | 到期付款日期（ISO 8601） |

**账单状态枚举：** `paid`、`partially-paid`、`unpaid`、`refunded`、`partially-refunded`、`refunding`、`cancelled`、`fallback-to-balance`

**账单分类枚举：** `compute`、`topup`、`traffic`、`penalty`、`custom`

---

### 获取账单详情

```
GET /billing/invoices/{id}
```

**所需权限：** `billing:read`

**路径参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| `id` | string (UUID) | 账单 UUID |

**响应：**

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

## 用户

### 获取用户信息

```
GET /user/info
```

**所需权限：** `user:read`

**参数：** 无

**响应：**

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

**字段：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `uuid` | string | 是 | 用户 UUID |
| `email` | string | 是 | 邮箱地址 |
| `nickname` | string | 否 | 昵称 |
| `created_at` | string | 否 | 账号创建时间（ISO 8601） |

---

## 速查表

| 方法 | 接口 | 权限 | 说明 |
|------|------|------|------|
| GET | `/compute/instance/list` | `instance:read` | 获取云实例列表（分页） |
| GET | `/compute/instance/{id}` | `instance:read` | 获取云实例详情 |
| GET | `/compute/instance/{id}/status` | `instance:read` | 获取云实例电源状态 |
| GET | `/compute/instance/{id}/traffic` | `instance:read` | 获取云实例流量 |
| POST | `/compute/instance/{id}/power` | `instance:power` | 启动/停止/重启（云实例） |
| GET | `/compute/instance/{id}/metrics` | `instance:read` | 云实例监控指标 |
| GET | `/compute/bare-metal/list` | `bare-metal:read` | 获取独立服务器列表（分页） |
| GET | `/compute/bare-metal/{id}` | `bare-metal:read` | 获取独立服务器详情 |
| GET | `/compute/bare-metal/{id}/status` | `bare-metal:read` | 获取独立服务器电源状态 |
| GET | `/compute/bare-metal/{id}/traffic` | `bare-metal:read` | 获取独立服务器流量 |
| POST | `/compute/bare-metal/{id}/power` | `bare-metal:power` | 开/关/重启（独立服务器） |
| GET | `/compute/bare-metal/{id}/metrics` | `bare-metal:read` | 独立服务器监控指标 |
| POST | `/compute/bare-metal/{id}/rebuild` | `bare-metal:rebuild` | 重装系统（破坏性操作） |
| POST | `/compute/bare-metal/{id}/rescue` | `bare-metal:rescue` | 进入救援模式 |
| GET | `/compute/bare-metal/{id}/profiles` | `bare-metal:read` | 获取操作系统配置列表 |
| GET | `/compute/bare-metal/{id}/rescue-profiles` | `bare-metal:read` | 获取救援配置列表 |
| GET | `/billing/balance` | `billing:read` | 获取账户余额 |
| GET | `/billing/invoices` | `billing:read` | 获取账单列表（分页） |
| GET | `/billing/invoices/{id}` | `billing:read` | 获取账单详情 |
| GET | `/user/info` | `user:read` | 获取用户信息 |
