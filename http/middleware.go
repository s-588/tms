package http

import (
	"fmt"
	"log/slog"
	"net/http"
)

func LogMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		slog.Info(fmt.Sprintf("%s %s %s", r.Method, r.URL.String(), r.RemoteAddr))
	}
}
