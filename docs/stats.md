# Cache Statistics

The API provides real-time cache performance metrics to help you understand usage patterns and optimize your integration.

## Endpoint

```
GET /debug/cache-stats
```

Returns JSON with cache performance data:

```json
{
  "total_requests": 1720000,
  "cache_hits": 1400000,
  "cache_misses": 320000,
  "new_books_cached": 310000,
  "scraping_errors": 10000,
  "hit_ratio": 81.4,
  "miss_ratio": 18.6,
  "new_book_ratio": 18.0
}
```

## Metrics Explained

| Metric | Description |
|--------|-------------|
| `total_requests` | Total API requests since server start |
| `cache_hits` | Requests served from cache (fast, no external calls) |
| `cache_misses` | Requests requiring a fresh lookup |
| `new_books_cached` | Unique books added to cache |
| `scraping_errors` | Failed external lookups |
| `hit_ratio` | Percentage served from cache |
| `miss_ratio` | Percentage requiring external lookup |
| `new_book_ratio` | Percentage of requests for new books |

## Using the Data

**High hit ratio (>80%)**: Your integration benefits from caching. Most requests are fast and don't hit external services.

**Low hit ratio (<50%)**: You're frequently requesting new/unique books. Consider if caching TTL meets your needs.

**High new_book_ratio**: Many unique books being discovered. Good for identifying popular titles in your use case.

## Access

The endpoint requires authentication via Bearer token:

```bash
curl -H "Authorization: Bearer YOUR_ADMIN_API_KEY" \
  https://bookcover.longitood.com/debug/cache-stats
```

**Authentication Required**: Set `ADMIN_API_KEY` environment variable to enable access. Requests without the correct Bearer token will receive a 401 Unauthorized response.

This endpoint is for administrative monitoring purposes only.
