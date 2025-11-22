package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/semaphoreui/semaphore/api/helpers"
	"github.com/semaphoreui/semaphore/db"
	"github.com/stretchr/testify/assert"
)

func TestClearCache_AdminUser(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/cache", nil)
	w := httptest.NewRecorder()

	adminUser := &db.User{
		ID:    1,
		Admin: true,
	}

	req = helpers.SetContextValue(req, "user", adminUser)

	clearCache(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestClearCache_NonAdminUser(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/cache", nil)
	w := httptest.NewRecorder()

	nonAdminUser := &db.User{
		ID:    1,
		Admin: false,
	}

	req = helpers.SetContextValue(req, "user", nonAdminUser)

	clearCache(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}
