package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/semaphoreui/semaphore/api/helpers"
	log "github.com/sirupsen/logrus"
)

// PanicRecoveryMiddleware recovers from panics and returns a proper error response
// This middleware should be applied early in the middleware chain to catch all panics
func PanicRecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic with context
				log.WithFields(log.Fields{
					"error":   err,
					"path":    r.URL.Path,
					"method":  r.Method,
					"remote":  r.RemoteAddr,
					"stack":   string(debug.Stack()),
				}).Error("Panic recovered in HTTP handler")

				// Return a generic error response to avoid exposing internal details
				helpers.WriteErrorStatus(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

