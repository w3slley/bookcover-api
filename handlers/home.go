package handlers

import (
  "encoding/json"
  "fmt"
  "net/http"
)

const METHOD_NOT_SUPPORTED string = "Method not suported yet."

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

func Home(writer http.ResponseWriter, request *http.Request) {
  writer.Header().Set("Content-Type", "application/json")
  httpException := HttpException{  
    statusCode: http.StatusBadRequest, message: METHOD_NOT_SUPPORTED,
  } 
  response := BuildErrorResponse(writer, httpException)
  writer.Write(response)
}

