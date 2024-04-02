package main

import (
	"bookcover-api/handlers"
	"bookcover-api/middlewares"
	"fmt"
	"net/http"
	"strconv"
)

const PORT int = 8000

func main() {
  http.HandleFunc("/", middlewares.JsonHeaderMiddleware(handlers.Home))
  http.HandleFunc("/bookcover", middlewares.JsonHeaderMiddleware(handlers.Bookcover))

  fmt.Printf("Server listening at port %d 🚀\n", PORT)
  http.ListenAndServe(":" + strconv.Itoa(PORT), nil)
}

