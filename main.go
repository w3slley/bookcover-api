package main

import (
	"fmt"
	"log"
	"os"
	"net/http"
	"strconv"

	"bookcover-api/internal/middlewares"
	"bookcover-api/internal/routes"

	"github.com/joho/godotenv"
)

const PORT int = 8000

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file")
	}

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
	port := os.Getenv("PORT")
	if port == "" {
		port = strconv.Itoa(PORT)
	}

	fmt.Printf("Server listening at port %s 🚀\n", port)
	http.ListenAndServe(":"+port, nil)
}
