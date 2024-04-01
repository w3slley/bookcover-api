package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HttpException struct {
  statusCode int
  message string
}


func BuildErrorResponse(w http.ResponseWriter, ex HttpException) []byte {
  jsonObj := map[string] string { "error": ex.message }
  res, err := json.Marshal(jsonObj)
  if err != nil {
    fmt.Println(err)
  }

  w.WriteHeader(ex.statusCode)
  return res
}
