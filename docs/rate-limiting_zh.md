# 频率限制

**[English](rate-limiting.md)**

Neburst OpenAPI 对每个用户实施速率限制，以确保资源的公平使用和平台的稳定性。

## 限制参数

每个请求会经过两层限流：

| 层级 | 窗口 | 最大请求数 | 适用范围 |
|------|------|-----------|---------|
| 全局限流 | 1 分钟（滑动窗口） | **每用户 60 次** | 所有端点合并计算 |
| 单次请求限流 | 5 秒（滑动窗口） | **每用户 1 次** | 每个独立请求 |

- **全局限流**：所有端点合计每分钟不超过 60 次请求。
- **单次请求限流**：每 5 秒最多 1 次请求，发送频率高于此将被拒绝。

两层限流均按**用户**维度计算，而非 API 密钥。同一用户的多个 API 密钥共享同一额度。

## 响应头

每个响应都包含三个速率限制相关的请求头：

| 请求头 | 类型 | 描述 |
|--------|------|------|
| `X-RateLimit-Limit` | integer | 当前窗口内允许的最大请求数。固定为 `60`。 |
| `X-RateLimit-Remaining` | integer | 达到限制前剩余的可用请求数。 |
| `X-RateLimit-Reset` | integer | Unix 时间戳（秒），表示窗口中最早的请求何时过期并释放容量。 |

### 响应头示例

```http
HTTP/1.1 200 OK
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 42
X-RateLimit-Reset: 1750291260
Content-Type: application/json
```

## 超过限制时

如果在滑动窗口内超过 60 次请求，API 将返回：

- **HTTP 状态码：** `429 Too Many Requests`
- **错误码：** `58204`

```json
{
  "code": 58204,
  "msg": "Rate limit exceeded"
}
```

`X-RateLimit-Reset` 响应头会告知你最早可用时隙对应的 Unix 时间戳。

## 实现细节

速率限制器使用基于 **Redis Sorted Set (ZSET)** 的滑动窗口算法：

1. 每个请求作为一个成员记录在以用户 ID 为键的 ZSET 中，当前时间戳（毫秒）作为分数。
2. 每次请求时，通过 `ZREMRANGEBYSCORE` 移除超出窗口范围（60 秒前）的条目。
3. 通过 `ZADD` 添加当前请求。
4. `ZCARD` 返回窗口内的请求总数。
5. 如果总数超过 60，则拒绝该请求。
6. ZSET 键的 TTL 略长于窗口大小，以确保自动清理。

该方案提供精确到亚秒级的用户维度跟踪，且不存在固定窗口的"边界突发"问题。

## 最佳实践

### 遵循响应头指引

在发送更多请求之前，始终检查 `X-RateLimit-Remaining`。当其值为 `0` 时，等待至 `X-RateLimit-Reset` 指示的时间。

### 实现指数退避

当收到 `429` 响应时，切勿立即重试。请使用指数退避策略：

```
wait_time = min(base_delay * 2^attempt, max_delay)
```

合理的起始参数为 `base_delay = 1 second`、`max_delay = 30 seconds`。

### 尽可能批量请求

如果需要获取多个实例的数据，优先考虑使用列表端点（`GET /compute/instance/list` 或 `GET /compute/bare-metal/list`），而非逐一请求每个实例。

### 缓存响应

对于不经常变化的数据（实例元数据、发票等），建议在客户端缓存响应，并按合理间隔刷新，而非频繁轮询。

### 监控使用量

在日志或监控指标中跟踪 `X-RateLimit-Remaining` 响应头，以便在接近限制时及时发现，从而在生产环境中触发 `429` 错误之前进行优化。
