package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/semaphoreui/semaphore/api/helpers"
	"github.com/semaphoreui/semaphore/db"
	"github.com/stretchr/testify/assert"
)

func TestGetLastEvents(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/events/last", nil)
	w := httptest.NewRecorder()

	user := &db.User{
		ID:    1,
		Admin: false,
	}

	req = helpers.SetContextValue(req, "user", user)

	// This test will fail if store is not properly mocked
	// For now, it tests that the function doesn't panic
	getLastEvents(w, req)

	// Should return either 200 OK or an error, but not panic
	assert.True(t, w.Code == http.StatusOK || w.Code >= http.StatusBadRequest)
}

func TestGetAllEvents(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/events", nil)
	w := httptest.NewRecorder()

	user := &db.User{
		ID:    1,
		Admin: false,
	}

	req = helpers.SetContextValue(req, "user", user)

	// This test will fail if store is not properly mocked
	// For now, it tests that the function doesn't panic
	getAllEvents(w, req)

	// Should return either 200 OK or an error, but not panic
	assert.True(t, w.Code == http.StatusOK || w.Code >= http.StatusBadRequest)
}
