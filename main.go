package main

import (
	"fmt"
	"net/http"
	"strconv"

	"bookcover-api/internal/middlewares"
	"bookcover-api/internal/routes"
)

const PORT int = 8000

func main() {
	http.HandleFunc("/", middlewares.Chain(
		routes.Home,
		middlewares.JsonHeaderMiddleware(),
		middlewares.CorsHeaderMiddleware(),
	))
	http.HandleFunc("/bookcover", middlewares.Chain(
		routes.BookcoverSearch,
		middlewares.HttpMethod("GET"),
		middlewares.JsonHeaderMiddleware(),
		middlewares.CorsHeaderMiddleware(),
	))

	http.HandleFunc("/bookcover/{isbn}", middlewares.Chain(
		routes.BookcoverByIsbn,
		middlewares.HttpMethod("GET"),
		middlewares.JsonHeaderMiddleware(),
		middlewares.CorsHeaderMiddleware(),
	))

	fmt.Printf("Server listening at port %d 🚀\n", PORT)
	http.ListenAndServe(":"+strconv.Itoa(PORT), nil)
}
