package routes

import (
	"net/http"

	"bookcover-api/internal/helpers"
)

func Home(w http.ResponseWriter, request *http.Request) {
	response := BuildErrorResponse(w, HttpException{
		statusCode: http.StatusBadRequest, message: helpers.ROUTE_NOT_SUPPORTED,
	})
	w.Write(response)
}
