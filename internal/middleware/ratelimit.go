package middleware

import (
	"fmt"
	"net"
	"net/http"
	"strconv"

	"bookcover-api/internal/cache"
	"bookcover-api/pkg/response"

	"github.com/bradfitz/gomemcache/memcache"
)

const (
	DefaultDailyLimit   = 100
	DefaultMonthlyLimit = 1000
	dailyTTL            = 86400   // 24 hours
	monthlyTTL          = 2592000 // 30 days
)

type RateLimitConfig struct {
	DailyLimit   int
	MonthlyLimit int
	Unlimited    bool
}

var (
	FreeTier = RateLimitConfig{
		DailyLimit:   DefaultDailyLimit,
		MonthlyLimit: DefaultMonthlyLimit,
	}
	ProTier = RateLimitConfig{
		Unlimited: true,
	}
)

func RateLimitMiddleware(cacheClient cache.CacheClient) Middleware {
	return RateLimitMiddlewareWithConfig(cacheClient, ProTier)
}

func RateLimitMiddlewareWithConfig(cacheClient cache.CacheClient, cfg RateLimitConfig) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			activeCfg := cfg

			if activeCfg.Unlimited {
				f(w, r)
				return
			}

			ip := getClientIP(r)

			dailyKey := fmt.Sprintf("ratelimit:%s:daily", ip)
			monthlyKey := fmt.Sprintf("ratelimit:%s:monthly", ip)

			dailyCount, err := incrementCounter(cacheClient, dailyKey, dailyTTL)
			if err != nil {
				f(w, r)
				return
			}

			monthlyCount, err := incrementCounter(cacheClient, monthlyKey, monthlyTTL)
			if err != nil {
				f(w, r)
				return
			}

			w.Header().Set("X-RateLimit-Limit-Daily", strconv.Itoa(activeCfg.DailyLimit))
			w.Header().Set("X-RateLimit-Remaining-Daily", strconv.Itoa(max(0, activeCfg.DailyLimit-int(dailyCount))))
			w.Header().Set("X-RateLimit-Limit-Monthly", strconv.Itoa(activeCfg.MonthlyLimit))
			w.Header().Set("X-RateLimit-Remaining-Monthly", strconv.Itoa(max(0, activeCfg.MonthlyLimit-int(monthlyCount))))

			if int(dailyCount) > activeCfg.DailyLimit || int(monthlyCount) > activeCfg.MonthlyLimit {
				w.Write(response.Error(w, http.StatusTooManyRequests, "Rate limit exceeded"))
				return
			}

			f(w, r)
		}
	}
}

func incrementCounter(c cache.CacheClient, key string, ttl int32) (uint64, error) {
	newVal, err := c.Increment(key, 1)
	if err == memcache.ErrCacheMiss {
		err = c.Add(&memcache.Item{
			Key:        key,
			Value:      []byte("1"),
			Expiration: ttl,
		})
		if err == memcache.ErrNotStored {
			// Another request created it; retry increment
			return c.Increment(key, 1)
		}
		if err != nil {
			return 0, err
		}
		return 1, nil
	}
	return newVal, err
}

func getClientIP(r *http.Request) string {
	if cfIP := r.Header.Get("CF-Connecting-IP"); cfIP != "" {
		return cfIP
	}
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// First IP in the chain is the real client
		for i := 0; i < len(xff); i++ {
			if xff[i] == ',' {
				return xff[:i]
			}
		}
		return xff
	}
	if xri := r.Header.Get("X-Real-Ip"); xri != "" {
		return xri
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
