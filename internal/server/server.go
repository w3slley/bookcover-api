package server

import (
	"fmt"
	"log"
	"net/http"

	"bookcover-api/internal/cache"
	"bookcover-api/internal/handler"
	"bookcover-api/internal/metrics"
	"bookcover-api/internal/middleware"
	"bookcover-api/internal/scraper"
	"bookcover-api/internal/service"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const port = 8000

func Start() error {
	if err := godotenv.Load(); err != nil {
		log.Print("Error loading .env file")
	}

	cacheClient := cache.GetCache()
	goodreadsScraper := scraper.NewGoodreads()
	bookcoverService := service.NewBookcoverService(goodreadsScraper, cacheClient)
	bookcoverHandler := handler.NewBookcoverHandler(bookcoverService)

	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/", middleware.Chain(
		handler.Home,
		metrics.MetricsMiddleware(),
		middleware.JsonHeaderMiddleware(),
		middleware.CorsHeaderMiddleware(),
	))

	http.HandleFunc("/bookcover", middleware.Chain(
		bookcoverHandler.Search,
		metrics.MetricsMiddleware(),
		middleware.HttpMethod("GET"),
		middleware.JsonHeaderMiddleware(),
		middleware.CorsHeaderMiddleware(),
	))

	http.HandleFunc("/bookcover/{isbn}", middleware.Chain(
		bookcoverHandler.ByISBN,
		metrics.MetricsMiddleware(),
		middleware.HttpMethod("GET"),
		middleware.JsonHeaderMiddleware(),
		middleware.CorsHeaderMiddleware(),
	))

	fmt.Printf("Server listening at port %d 🚀\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
