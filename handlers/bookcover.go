package handlers

import (
	"io"
	"net/http"
)

func Bookcover(w http.ResponseWriter, r *http.Request) {
  io.WriteString(w, "Being implemented...")
}
