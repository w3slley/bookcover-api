package main

import (
	"bookcover-api/internal/middlewares"
	"bookcover-api/internal/routes"
	"fmt"
	"net/http"
	"strconv"
)

const PORT int = 8000

func main() {
	http.HandleFunc("/", middlewares.JsonHeaderMiddleware(routes.Home))
	http.HandleFunc("/bookcover", middlewares.JsonHeaderMiddleware(routes.BookcoverSearch))
	http.HandleFunc("/bookcover/{isbn}", middlewares.JsonHeaderMiddleware(routes.BookcoverByIsbn))

	fmt.Printf("Server listening at port %d 🚀\n", PORT)
	http.ListenAndServe(":"+strconv.Itoa(PORT), nil)
}
