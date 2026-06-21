# Rate Limiting

**[中文版](rate-limiting_zh.md)**

The Neburst OpenAPI enforces per-user rate limits to ensure fair resource usage and platform stability.

## Limits

Two layers of rate limiting are applied to every request:

| Layer | Window | Max Requests | Scope |
|-------|--------|-------------|-------|
| Global | 1 minute (sliding) | **60 per user** | All endpoints combined |
| Per-request | 5 seconds (sliding) | **1 per user** | Each individual request |

- **Global limit**: No more than 60 requests per minute across all endpoints.
- **Per-request limit**: At most 1 request every 5 seconds. Sending requests faster than this will be rejected.

Both limits are applied **per user**, not per API key. If a user has multiple API keys, all keys share the same budget.

## Response Headers

Every response includes three rate-limit headers:

| Header | Type | Description |
|--------|------|-------------|
| `X-RateLimit-Limit` | integer | Maximum number of requests allowed in the current window. Always `60`. |
| `X-RateLimit-Remaining` | integer | Number of requests remaining before the limit is reached. |
| `X-RateLimit-Reset` | integer | Unix timestamp (seconds) indicating when the oldest request in the window expires, freeing capacity. |

### Example headers

```http
HTTP/1.1 200 OK
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 42
X-RateLimit-Reset: 1750291260
Content-Type: application/json
```

## When the Limit Is Exceeded

If you exceed 60 requests in the sliding window, the API returns:

- **HTTP Status:** `429 Too Many Requests`
- **Error Code:** `58204`

```json
{
  "code": 58204,
  "msg": "Rate limit exceeded"
}
```

The `X-RateLimit-Reset` header tells you the earliest Unix timestamp when a slot will become available.

## Implementation Details

The rate limiter uses a **Redis Sorted Set (ZSET)** sliding window algorithm:

1. Each request is recorded as a member in a ZSET keyed by user ID, with the current timestamp (milliseconds) as the score.
2. On each request, entries older than the window (60 seconds ago) are removed via `ZREMRANGEBYSCORE`.
3. The current request is added via `ZADD`.
4. `ZCARD` returns the total count of requests in the window.
5. If the count exceeds 60, the request is rejected.
6. The ZSET key has a TTL slightly longer than the window to ensure automatic cleanup.

This approach provides accurate per-user tracking with sub-second precision, without the "boundary burst" problem of fixed windows.

## Best Practices

### Respect the headers

Always check `X-RateLimit-Remaining` before making additional requests. If it reaches `0`, wait until the time indicated by `X-RateLimit-Reset`.

### Implement exponential backoff

When you receive a `429` response, do not retry immediately. Use exponential backoff:

```
wait_time = min(base_delay * 2^attempt, max_delay)
```

A reasonable starting point is `base_delay = 1 second`, `max_delay = 30 seconds`.

### Batch where possible

If you need data for multiple instances, consider whether you can use the list endpoint (`GET /compute/instance/list` or `GET /compute/bare-metal/list`) instead of making individual requests for each instance.

### Cache responses

For data that does not change frequently (instance metadata, invoices), cache responses client-side and refresh on a reasonable interval rather than polling aggressively.

### Monitor your usage

Track the `X-RateLimit-Remaining` header in your logs or metrics to detect when your integration is approaching the limit, so you can optimize before hitting `429` errors in production.
