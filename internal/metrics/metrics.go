package metrics

import (
	"net/http"
	"strconv"
	"strings"
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
