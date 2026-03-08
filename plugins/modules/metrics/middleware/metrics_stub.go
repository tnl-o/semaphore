package middleware

import "net/http"

// MetricsMiddleware is a no-op middleware when Prometheus is not enabled
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
