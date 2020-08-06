package api

import (
	"net/http"

	"github.com/gorilla/handlers"
)

// LimitBody limits the size of incoming request body to 1MB.
func (s *Server) LimitBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const MB = 1024 * 1024
		r.Body = http.MaxBytesReader(w, r.Body, 1*MB)
		next.ServeHTTP(w, r)
	})
}

// Compress gzip compresses HTTP responses for clients that support it
// via the 'Accept-Encoding' header.
func (s *Server) Compress(next http.Handler) http.Handler {
	return handlers.CompressHandler(next)
}
