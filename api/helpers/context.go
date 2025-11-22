package helpers

import (
	"context"
	"net/http"

	"github.com/semaphoreui/semaphore/db"
)

const (
	// CorrelationIDContextKey is the context key for correlation ID
	CorrelationIDContextKey contextKey = "correlation_id"
)

type contextKey string

func GetFromContext(r *http.Request, key string) any {
	return r.Context().Value(key)
}

func GetOkFromContext(r *http.Request, key string) (res any, ok bool) {
	res = r.Context().Value(key)
	return res, res != nil
}

func SetContextValue(r *http.Request, key string, value any) *http.Request {
	ctx := r.Context()
	// Convert string key to contextKey type if it matches known keys
	var ctxKey any = key
	if key == string(CorrelationIDContextKey) {
		ctxKey = CorrelationIDContextKey
	}
	ctx = context.WithValue(ctx, ctxKey, value)
	return r.WithContext(ctx)
}

func UserFromContext(r *http.Request) *db.User {
	return GetFromContext(r, "user").(*db.User)
}

func GetGlobalRole(r *http.Request) db.Role {
	return GetFromContext(r, "role").(db.Role)
}

// GetCorrelationID retrieves the correlation ID from the request context
func GetCorrelationID(r *http.Request) string {
	if correlationID, ok := r.Context().Value(CorrelationIDContextKey).(string); ok {
		return correlationID
	}
	return ""
}
