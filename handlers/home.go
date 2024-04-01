package handlers

import (
  "net/http"
)

func Home(w http.ResponseWriter, request *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  response := BuildErrorResponse(w, HttpException{  
    statusCode: http.StatusBadRequest, message: METHOD_NOT_SUPPORTED,
  })
  w.Write(response)
}

