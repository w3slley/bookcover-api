package response

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// Success writes a successful JSON response with the given URL
func Success(w http.ResponseWriter, url string) []byte {
	var buffer bytes.Buffer
	enc := json.NewEncoder(&buffer)
	enc.SetEscapeHTML(false)
	enc.Encode(map[string]string{"url": url})
	w.WriteHeader(http.StatusOK)
	return buffer.Bytes()
}

// Error writes an error JSON response with the given status code and message
func Error(w http.ResponseWriter, statusCode int, message string) []byte {
	data, err := json.Marshal(map[string]string{"error": message})
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return nil
	}
	w.WriteHeader(statusCode)
	return data
}
