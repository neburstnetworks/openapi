# Scopes

**[中文版](scopes_zh.md)**

Scopes control what an API key is allowed to do. Each endpoint requires a specific scope. A key without the required scope receives a `403` response with error code `58203`.

## Available Scopes

| Scope | Description |
|-------|-------------|
| `instance:read` | Read cloud instance information: list instances, get details, power status, traffic usage, metrics |
| `instance:power` | Execute cloud instance power actions: start, stop, restart |
| `bare-metal:read` | Read bare-metal instance information: list instances, get details, power status, traffic usage, OS profiles, rescue profiles, metrics |
| `bare-metal:power` | Execute bare-metal power actions: power on/off/cycle/reset |
| `bare-metal:rebuild` | Rebuild (reinstall) a bare-metal instance's operating system |
| `bare-metal:rescue` | Boot a bare-metal instance into rescue/recovery mode |
| `billing:read` | Read billing information: account balance, invoices |
| `user:read` | Read user account information |

## Wildcards

Scopes support two levels of wildcards:

| Pattern | Matches |
|---------|---------|
| `*` | **All scopes** -- full access to every endpoint |
| `instance:*` | All cloud instance scopes: `instance:read`, `instance:power` |
| `bare-metal:*` | All bare-metal scopes: `bare-metal:read`, `bare-metal:power`, `bare-metal:rebuild`, `bare-metal:rescue` |
| `billing:*` | All billing scopes: `billing:read` |
| `user:*` | All user scopes: `user:read` |

Wildcards are evaluated using prefix matching. `instance:*` matches any scope that starts with `instance:`.

## Scope-to-Endpoint Mapping

| Endpoint | Method | Required Scope |
|----------|--------|---------------|
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

## Scope Assignment

Scopes are assigned when creating an API key and can be updated later. The key's scope list is stored as a JSON array:

```json
["instance:read", "instance:power", "billing:read", "user:read"]
```

### Recommended scope sets by use case

**Monitoring / Dashboard integration:**
```json
["instance:read", "bare-metal:read", "billing:read", "user:read"]
```
Read-only access to instance status, billing data, and user information. Suitable for monitoring dashboards and alerting systems.

**Automation / Infrastructure management:**
```json
["instance:read", "instance:power", "bare-metal:read", "bare-metal:power"]
```
Read instance information and control power state. Suitable for auto-scaling scripts or scheduled restart automation.

**Full compute management:**
```json
["instance:*", "bare-metal:*"]
```
Complete control over all compute resources including rebuild and rescue. Suitable for infrastructure-as-code tools.

**Full access:**
```json
["*"]
```
Unrestricted access to all endpoints. Use sparingly and only when necessary.

## Principle of Least Privilege

Always assign the minimum set of scopes required for your use case:

- Do not use `*` when specific scopes suffice.
- Separate read-only keys from keys with write/action permissions.
- Create dedicated keys for each integration or automation script.
- Review and prune unused scopes periodically.

## Scope Validation Behavior

1. When a request reaches the server, the middleware extracts the required scope from an internal route-to-scope mapping.
2. The API key's assigned scopes (JSON array) are compared against the required scope.
3. A match occurs if any of the following is true:
   - The key has the exact scope (e.g., `instance:read` matches `instance:read`)
   - The key has a category wildcard (e.g., `instance:*` matches `instance:read`)
   - The key has the global wildcard (`*`)
4. If no match is found, the request is rejected with error code `58203` (`INSUFFICIENT_SCOPE`).
