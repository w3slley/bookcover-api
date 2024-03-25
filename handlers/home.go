package handlers

import (
  "encoding/json"
  "fmt"
  "net/http"
)

const METHOD_NOT_SUPPORTED string = "Method not suported yet."

func BuildErrorResponse(w http.ResponseWriter) []byte {
  jsonObj :=map[string] string { "error": METHOD_NOT_SUPPORTED}
  res, err := json.Marshal(jsonObj)
  if err != nil {
    fmt.Println(err)
  }

  w.WriteHeader(http.StatusBadRequest)
  return res
}

func Home(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
   
  res := BuildErrorResponse(w)
  w.Write(res)
}

