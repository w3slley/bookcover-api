package handler

import (
	"net/http"

	"bookcover-api/static"
)

func Home(w http.ResponseWriter, r *http.Request) {
	data, err := static.Files.ReadFile("index.html")
	if err != nil {
		http.Error(w, "page not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(data)
}
