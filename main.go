package main

import (
	"bookcover-api/internal/middlewares"
	"bookcover-api/internal/routes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/joho/godotenv"
)

const PORT int = 8000

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	http.HandleFunc("/", middlewares.JsonHeaderMiddleware(routes.Home))
	http.HandleFunc("/bookcover", middlewares.JsonHeaderMiddleware(routes.BookcoverSearch))
	http.HandleFunc("/bookcover/{isbn}", middlewares.JsonHeaderMiddleware(routes.BookcoverByIsbn))

	fmt.Printf("Server listening at port %d 🚀\n", PORT)
	http.ListenAndServe(":"+strconv.Itoa(PORT), nil)
}
