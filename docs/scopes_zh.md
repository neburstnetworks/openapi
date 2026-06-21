# 权限范围（Scopes）

**[English](scopes.md)**

权限范围控制 API Key 可以执行的操作。每个接口都需要特定的权限。若 Key 缺少所需权限，将收到 `403` 响应，错误码为 `58203`。

## 可用权限

| 权限 | 说明 |
|------|------|
| `instance:read` | 读取云实例信息：列出实例、获取详情、电源状态、流量用量、监控指标 |
| `instance:power` | 执行云实例电源操作：启动、停止、重启 |
| `bare-metal:read` | 读取独立服务器信息：列出实例、获取详情、电源状态、流量用量、操作系统配置、救援配置、监控指标 |
| `bare-metal:power` | 执行独立服务器电源操作：开机/关机/重启/硬重启 |
| `bare-metal:rebuild` | 重装独立服务器的操作系统 |
| `bare-metal:rescue` | 将独立服务器引导至救援/恢复模式 |
| `billing:read` | 读取账单信息：账户余额、账单记录 |
| `user:read` | 读取用户账号信息 |

## 通配符

权限支持两级通配符匹配：

| 模式 | 匹配范围 |
|------|----------|
| `*` | **所有权限** -- 完全访问所有接口 |
| `instance:*` | 所有云实例权限：`instance:read`、`instance:power` |
| `bare-metal:*` | 所有独立服务器权限：`bare-metal:read`、`bare-metal:power`、`bare-metal:rebuild`、`bare-metal:rescue` |
| `billing:*` | 所有账单权限：`billing:read` |
| `user:*` | 所有用户权限：`user:read` |

通配符采用前缀匹配。`instance:*` 可匹配所有以 `instance:` 开头的权限。

## 权限与接口映射

| 接口 | 方法 | 所需权限 |
|------|------|----------|
| `/compute/instance/list` | GET | `instance:read` |
| `/compute/instance/{id}` | GET | `instance:read` |
| `/compute/instance/{id}/status` | GET | `instance:read` |
| `/compute/instance/{id}/traffic` | GET | `instance:read` |
| `/compute/instance/{id}/power` | POST | `instance:power` |
| `/compute/instance/{id}/metrics` | GET | `instance:read` |
| `/compute/bare-metal/list` | GET | `bare-metal:read` |
| `/compute/bare-metal/{id}` | GET | `bare-metal:read` |
| `/compute/bare-metal/{id}/status` | GET | `bare-metal:read` |
| `/compute/bare-metal/{id}/traffic` | GET | `bare-metal:read` |
| `/compute/bare-metal/{id}/power` | POST | `bare-metal:power` |
| `/compute/bare-metal/{id}/metrics` | GET | `bare-metal:read` |
| `/compute/bare-metal/{id}/rebuild` | POST | `bare-metal:rebuild` |
| `/compute/bare-metal/{id}/rescue` | POST | `bare-metal:rescue` |
| `/compute/bare-metal/{id}/profiles` | GET | `bare-metal:read` |
| `/compute/bare-metal/{id}/rescue-profiles` | GET | `bare-metal:read` |
| `/billing/balance` | GET | `billing:read` |
| `/billing/invoices` | GET | `billing:read` |
| `/billing/invoices/{id}` | GET | `billing:read` |
| `/user/info` | GET | `user:read` |

## 权限分配

权限在创建 API Key 时分配，之后可以修改。Key 的权限列表以 JSON 数组形式存储：

```json
["instance:read", "instance:power", "billing:read", "user:read"]
```

### 按使用场景推荐的权限组合

**监控 / 仪表盘集成：**
```json
["instance:read", "bare-metal:read", "billing:read", "user:read"]
```
对实例状态、账单数据和用户信息的只读访问。适用于监控仪表盘和告警系统。

**自动化 / 基础设施管理：**
```json
["instance:read", "instance:power", "bare-metal:read", "bare-metal:power"]
```
读取实例信息并控制电源状态。适用于弹性伸缩脚本或定时重启自动化。

**完整计算资源管理：**
```json
["instance:*", "bare-metal:*"]
```
对所有计算资源的完全控制，包括重装系统和救援模式。适用于基础设施即代码工具。

**完全访问：**
```json
["*"]
```
对所有接口的无限制访问。请谨慎使用，仅在确实需要时授予。

## 最小权限原则

请始终仅分配满足需求所需的最少权限：

- 有特定权限可用时，不要使用 `*` 通配符。
- 将只读 Key 与具有写入/操作权限的 Key 分开管理。
- 为每个集成或自动化脚本创建独立的 Key。
- 定期审查并清理不再使用的权限。

## 权限校验流程

1. 请求到达服务器时，中间件从内部的路由-权限映射中提取所需权限。
2. 将 API Key 已分配的权限（JSON 数组）与所需权限进行比对。
3. 以下任一条件满足即为匹配：
   - Key 拥有完全一致的权限（如 `instance:read` 匹配 `instance:read`）
   - Key 拥有分类通配符（如 `instance:*` 匹配 `instance:read`）
   - Key 拥有全局通配符（`*`）
4. 若无匹配，请求将被拒绝，返回错误码 `58203`（`INSUFFICIENT_SCOPE`）。
