package main

import (
	"log/slog"
	"os"
	"time"

	"bookcover-api/internal/metrics"
	"bookcover-api/internal/server"
)

func main() {
	go startMetricsLogger()

	if err := server.Start(); err != nil {
		slog.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}

func startMetricsLogger() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		stats := metrics.GetCacheMetrics().GetStats()
		slog.Info("cache metrics",
			"total_requests", stats.TotalRequests,
			"cache_hits", stats.CacheHits,
			"cache_misses", stats.CacheMisses,
			"new_books_cached", stats.NewBooksCached,
			"hit_ratio", stats.HitRatio(),
			"miss_ratio", stats.MissRatio(),
			"new_book_ratio", stats.NewBookRatio(),
		)
	}
}
