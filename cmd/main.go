package main

import (
	"net/http"
	"strconv"
  "bookcover-api/handlers"
)

const PORT int = 8000

func main() {
  http.HandleFunc("/", handlers.Home)
  http.HandleFunc("/bookcover", handlers.Bookcover)

  http.ListenAndServe(":" + strconv.Itoa(PORT), nil)
}

