package main

import (
	"bookcover-api/handlers"
	"bookcover-api/middlewares"
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

  http.HandleFunc("/", middlewares.JsonHeaderMiddleware(handlers.Home))
  http.HandleFunc("/bookcover", middlewares.JsonHeaderMiddleware(handlers.BookcoverSearch))
  http.HandleFunc("/bookcover/{isbn}", middlewares.JsonHeaderMiddleware(handlers.BookcoverByIsbn))

  fmt.Printf("Server listening at port %d 🚀\n", PORT)
  http.ListenAndServe(":" + strconv.Itoa(PORT), nil)
}

