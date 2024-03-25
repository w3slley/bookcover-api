package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const METHOD_NOT_SUPPORTED string = "Method not suported yet."

func Home(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  jsonObj :=map[string] string { "error": METHOD_NOT_SUPPORTED}
  res, err := json.Marshal(jsonObj)
  if err != nil {
    fmt.Println(err)
  } 

  w.WriteHeader(http.StatusBadRequest)
  w.Write(res)
}


