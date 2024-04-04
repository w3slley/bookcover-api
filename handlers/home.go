package handlers

import (
  "net/http"
)

func Home(w http.ResponseWriter, request *http.Request) {
  response := BuildErrorResponse(w, HttpException{  
    statusCode: http.StatusBadRequest, message: METHOD_NOT_SUPPORTED,
  })
  w.Write(response)
}

