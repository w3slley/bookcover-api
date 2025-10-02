package handler

import (
	"net/http"

	"bookcover-api/internal/config"
	"bookcover-api/pkg/response"
)

func Home(w http.ResponseWriter, r *http.Request) {
	w.Write(response.Error(w, http.StatusNotFound, config.RouteNotSupported))
}
