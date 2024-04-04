package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type HttpException struct {
  statusCode int
  message string
}

func BuildSuccessResponse(w http.ResponseWriter, url string) []byte {
  var buffer bytes.Buffer
  enc := json.NewEncoder(&buffer)
  enc.SetEscapeHTML(false)
  enc.Encode(map[string] string { "url": url })
  w.WriteHeader(200)

  return buffer.Bytes()
}

func BuildErrorResponse(w http.ResponseWriter, ex HttpException) []byte {
  data, err := json.Marshal(map[string] string { "error": ex.message })
  if err != nil {
    fmt.Println(err)
  }

  w.WriteHeader(ex.statusCode)
  return data
}
