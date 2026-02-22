# Rate Limiting

## Overview

The API enforces two layers of rate limiting to protect against abuse and ensure fair usage:

1. **Burst protection** — Handled by Cloudflare. Returns `429` if more than 5 requests are made within a 10-second window.
2. **Quota-based limits** — Enforced at the application level per client IP using memcached counters.

## Quota Limits

| Limit   | Threshold      | Window   |
|---------|---------------|----------|
| Daily   | 100 requests  | 24 hours |
| Monthly | 1,000 requests | 30 days  |

When a limit is exceeded, the API responds with `429 Too Many Requests`.

## Response Headers

Every response from `/bookcover` includes rate limit headers:

```
X-RateLimit-Limit-Daily: 100
X-RateLimit-Remaining-Daily: 95
X-RateLimit-Limit-Monthly: 1000
X-RateLimit-Remaining-Monthly: 980
```

## Client IP Detection

The middleware resolves the client IP in this order:

1. `CF-Connecting-IP` — Set by Cloudflare (primary)
2. `X-Forwarded-For` — First IP in the chain
3. `X-Real-Ip` — Fallback proxy header
4. `RemoteAddr` — Direct TCP connection

## How It Works

- Request counters are stored in memcached with keys `ratelimit:{ip}:daily` and `ratelimit:{ip}:monthly`.
- Each key has a TTL matching its window (24h or 30 days), so counters reset automatically.
- Counters are incremented atomically using memcached's `Increment` operation with an `Add`-based fallback for key initialization.
- If memcached is unavailable, requests are allowed through (fail-open).
