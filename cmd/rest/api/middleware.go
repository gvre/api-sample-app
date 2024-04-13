package api

import (
	"net/http"
)

// LimitBody limits the size of incoming request body to 1MB.
func (s *Server) LimitBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const MB = 1024 * 1024
		r.Body = http.MaxBytesReader(w, r.Body, 1*MB)
		next.ServeHTTP(w, r)
	})
}

func (s *Server) EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}
