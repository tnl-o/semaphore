package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/semaphoreui/semaphore/api/helpers"
)

const (
	// CorrelationIDHeader is the HTTP header name for correlation ID
	CorrelationIDHeader = "X-Correlation-ID"
)

// CorrelationIDMiddleware adds correlation ID to requests for request tracing
func CorrelationIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		correlationID := r.Header.Get(CorrelationIDHeader)
		if correlationID == "" {
			correlationID = uuid.New().String()
		}

		// Add to response header
		w.Header().Set(CorrelationIDHeader, correlationID)

		// Add to request context using helpers
		r = helpers.SetContextValue(r, string(helpers.CorrelationIDContextKey), correlationID)

		next.ServeHTTP(w, r)
	})
}
