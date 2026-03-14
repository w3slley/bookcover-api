package metrics

import (
	"expvar"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"bookcover-api/internal/middleware"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bookcover_http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"path", "method", "status_code"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "bookcover_http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

// normalizePath replaces dynamic ISBN segments to avoid high-cardinality labels.
func normalizePath(path string) string {
	if strings.HasPrefix(path, "/bookcover/") {
		return "/bookcover/:isbn"
	}
	return path
}

// MetricsMiddleware returns a middleware that records request metrics.
func MetricsMiddleware() middleware.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}

			next(rec, r)

			path := normalizePath(r.URL.Path)
			duration := time.Since(start).Seconds()

			httpRequestsTotal.WithLabelValues(path, r.Method, strconv.Itoa(rec.statusCode)).Inc()
			httpRequestDuration.WithLabelValues(path, r.Method).Observe(duration)
		}
	}
}

type CacheMetrics struct {
	totalRequests  *expvar.Int
	cacheHits      *expvar.Int
	cacheMisses    *expvar.Int
	scrapingErrors *expvar.Int
	newBooksCached *expvar.Int
	mu             sync.RWMutex
}

var (
	instance *CacheMetrics
	once     sync.Once
)

func GetCacheMetrics() *CacheMetrics {
	once.Do(func() {
		instance = &CacheMetrics{
			totalRequests:  expvar.NewInt("cache_total_requests"),
			cacheHits:      expvar.NewInt("cache_hits"),
			cacheMisses:    expvar.NewInt("cache_misses"),
			scrapingErrors: expvar.NewInt("scraping_errors"),
			newBooksCached: expvar.NewInt("new_books_cached"),
		}
	})
	return instance
}

func (m *CacheMetrics) RecordRequest() {
	m.totalRequests.Add(1)
}

func (m *CacheMetrics) RecordCacheHit() {
	m.cacheHits.Add(1)
}

func (m *CacheMetrics) RecordCacheMiss() {
	m.cacheMisses.Add(1)
}

func (m *CacheMetrics) RecordScrapingError() {
	m.scrapingErrors.Add(1)
}

func (m *CacheMetrics) RecordNewBookCached() {
	m.newBooksCached.Add(1)
}

func (m *CacheMetrics) GetStats() Stats {
	return Stats{
		TotalRequests:  m.totalRequests.Value(),
		CacheHits:      m.cacheHits.Value(),
		CacheMisses:    m.cacheMisses.Value(),
		ScrapingErrors: m.scrapingErrors.Value(),
		NewBooksCached: m.newBooksCached.Value(),
	}
}

type Stats struct {
	TotalRequests  int64 `json:"total_requests"`
	CacheHits      int64 `json:"cache_hits"`
	CacheMisses    int64 `json:"cache_misses"`
	ScrapingErrors int64 `json:"scraping_errors"`
	NewBooksCached int64 `json:"new_books_cached"`
}

func (s Stats) HitRatio() float64 {
	if s.TotalRequests == 0 {
		return 0
	}
	return float64(s.CacheHits) / float64(s.TotalRequests)
}

func (s Stats) MissRatio() float64 {
	if s.TotalRequests == 0 {
		return 0
	}
	return float64(s.CacheMisses) / float64(s.TotalRequests)
}

func (s Stats) NewBookRatio() float64 {
	if s.TotalRequests == 0 {
		return 0
	}
	return float64(s.NewBooksCached) / float64(s.TotalRequests)
}
