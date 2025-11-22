package helpers

import (
	"net/http"

	"github.com/semaphoreui/semaphore/db"
	log "github.com/sirupsen/logrus"
)

// Logger returns a logrus entry with context from the request
// It includes correlation ID, user ID, and project ID if available
func Logger(r *http.Request) *log.Entry {
	entry := log.NewEntry(log.StandardLogger())

	// Add correlation ID
	if correlationID := GetCorrelationID(r); correlationID != "" {
		entry = entry.WithField("correlation_id", correlationID)
	}

	// Add user ID if available
	if user := GetFromContext(r, "user"); user != nil {
		if u, ok := user.(*db.User); ok {
			entry = entry.WithField("user_id", u.ID)
		}
	}

	// Add project ID if available
	if project := GetFromContext(r, "project"); project != nil {
		if p, ok := project.(*db.Project); ok {
			entry = entry.WithField("project_id", p.ID)
		}
	}

	// Add request path and method
	entry = entry.WithFields(log.Fields{
		"method": r.Method,
		"path":   r.URL.Path,
	})

	return entry
}

