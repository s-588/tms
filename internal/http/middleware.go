package http

import (
	"log/slog"
	"net/http"
)

func LogMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		slog.Info("", "method", r.Method, "url", r.URL.String(), "remote_addr", r.RemoteAddr)
	}
}
