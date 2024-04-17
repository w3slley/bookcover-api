package middlewares

import "net/http"

func JsonHeaderMiddleware(f http.HandlerFunc) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    f(w, r)
  }
}


