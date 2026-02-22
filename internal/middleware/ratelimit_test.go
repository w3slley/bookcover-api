package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bradfitz/gomemcache/memcache"

	"bookcover-api/internal/cache"
)

// mockCache implements cache.CacheClient for rate limit tests
type mockCache struct {
	items map[string]*memcache.Item
}

func newMockCache() *mockCache {
	return &mockCache{items: make(map[string]*memcache.Item)}
}

func (m *mockCache) Get(key string) (*memcache.Item, error) {
	if item, ok := m.items[key]; ok {
		return item, nil
	}
	return nil, memcache.ErrCacheMiss
}

func (m *mockCache) Set(item *memcache.Item) error {
	m.items[item.Key] = item
	return nil
}

func (m *mockCache) Add(item *memcache.Item) error {
	if _, exists := m.items[item.Key]; exists {
		return memcache.ErrNotStored
	}
	m.items[item.Key] = item
	return nil
}

func (m *mockCache) Increment(key string, delta uint64) (uint64, error) {
	item, exists := m.items[key]
	if !exists {
		return 0, memcache.ErrCacheMiss
	}
	val := uint64(0)
	for _, b := range item.Value {
		val = val*10 + uint64(b-'0')
	}
	val += delta
	item.Value = []byte(fmt.Sprintf("%d", val))
	return val, nil
}

var _ cache.CacheClient = (*mockCache)(nil)

func okHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestRateLimit_DailyLimitExceeded(t *testing.T) {
	mc := newMockCache()
	cfg := RateLimitConfig{DailyLimit: 5, MonthlyLimit: 1000}
	mw := RateLimitMiddlewareWithConfig(mc, cfg)(okHandler)

	for i := 1; i <= 5; i++ {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/bookcover", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		mw(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("request %d: expected 200, got %d", i, rr.Code)
		}
		expectedRemaining := fmt.Sprintf("%d", 5-i)
		if got := rr.Header().Get("X-RateLimit-Remaining-Daily"); got != expectedRemaining {
			t.Errorf("request %d: expected remaining %s, got %s", i, expectedRemaining, got)
		}
	}

	// 6th request should be rate limited
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/bookcover", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	mw(rr, req)

	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("expected 429, got %d", rr.Code)
	}
	if got := rr.Header().Get("X-RateLimit-Remaining-Daily"); got != "0" {
		t.Errorf("expected remaining 0, got %s", got)
	}
}

func TestRateLimit_MonthlyLimitExceeded(t *testing.T) {
	mc := newMockCache()
	cfg := RateLimitConfig{DailyLimit: 1000, MonthlyLimit: 5}
	mw := RateLimitMiddlewareWithConfig(mc, cfg)(okHandler)

	for i := 1; i <= 5; i++ {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/bookcover", nil)
		req.RemoteAddr = "10.0.0.1:9999"
		mw(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("request %d: expected 200, got %d", i, rr.Code)
		}
		expectedRemaining := fmt.Sprintf("%d", 5-i)
		if got := rr.Header().Get("X-RateLimit-Remaining-Monthly"); got != expectedRemaining {
			t.Errorf("request %d: expected remaining %s, got %s", i, expectedRemaining, got)
		}
	}

	// 6th request should be rate limited
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/bookcover", nil)
	req.RemoteAddr = "10.0.0.1:9999"
	mw(rr, req)

	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("expected 429, got %d", rr.Code)
	}
	if got := rr.Header().Get("X-RateLimit-Remaining-Monthly"); got != "0" {
		t.Errorf("expected remaining 0, got %s", got)
	}
}

func TestRateLimit_DifferentIPsTrackedSeparately(t *testing.T) {
	mc := newMockCache()
	cfg := RateLimitConfig{DailyLimit: 2, MonthlyLimit: 1000}
	mw := RateLimitMiddlewareWithConfig(mc, cfg)(okHandler)

	// Exhaust limit for IP A
	for i := 0; i < 2; i++ {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/bookcover", nil)
		req.RemoteAddr = "1.1.1.1:1234"
		mw(rr, req)
	}

	// IP A should be blocked
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/bookcover", nil)
	req.RemoteAddr = "1.1.1.1:1234"
	mw(rr, req)
	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("IP A: expected 429, got %d", rr.Code)
	}

	// IP B should still work
	rr = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/bookcover", nil)
	req.RemoteAddr = "2.2.2.2:1234"
	mw(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("IP B: expected 200, got %d", rr.Code)
	}
}

func TestRateLimit_UsesXForwardedFor(t *testing.T) {
	mc := newMockCache()
	cfg := RateLimitConfig{DailyLimit: 1, MonthlyLimit: 1000}
	mw := RateLimitMiddlewareWithConfig(mc, cfg)(okHandler)

	// First request from forwarded IP
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/bookcover", nil)
	req.Header.Set("X-Forwarded-For", "203.0.113.50, 70.41.3.18")
	mw(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	// Second request from same forwarded IP should be blocked
	rr = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/bookcover", nil)
	req.Header.Set("X-Forwarded-For", "203.0.113.50, 70.41.3.18")
	mw(rr, req)
	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("expected 429, got %d", rr.Code)
	}
}
