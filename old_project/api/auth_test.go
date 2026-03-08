package api

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/semaphoreui/semaphore/api/helpers"
	"github.com/semaphoreui/semaphore/db"
	"github.com/semaphoreui/semaphore/pkg/tz"
	"github.com/semaphoreui/semaphore/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSession_NoCookie(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	session, ok := getSession(req)

	assert.False(t, ok)
	assert.Nil(t, session)
}

func TestGetSession_InvalidCookie(t *testing.T) {
	// Setup cookie encryption
	setupTestCookie()

	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  "semaphore",
		Value: "invalid_cookie_value",
	})

	session, ok := getSession(req)

	assert.False(t, ok)
	assert.Nil(t, session)
}

func TestGetSession_InvalidTypeAssertions(t *testing.T) {
	// Setup cookie encryption
	setupTestCookie()

	// Test with wrong type in cookie (string instead of int)
	value := map[string]any{
		"user":    "not_an_int", // Should be int
		"session": 123,
	}

	encoded, err := util.Cookie.Encode("semaphore", value)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  "semaphore",
		Value: encoded,
	})

	session, ok := getSession(req)

	assert.False(t, ok)
	assert.Nil(t, session)
}

func TestGetSession_InvalidSessionType(t *testing.T) {
	// Setup cookie encryption
	setupTestCookie()

	// Test with wrong type for session (string instead of int)
	value := map[string]any{
		"user":    123,
		"session": "not_an_int", // Should be int
	}

	encoded, err := util.Cookie.Encode("semaphore", value)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  "semaphore",
		Value: encoded,
	})

	session, ok := getSession(req)

	assert.False(t, ok)
	assert.Nil(t, session)
}

func TestGetSession_MissingFields(t *testing.T) {
	// Setup cookie encryption
	setupTestCookie()

	// Test with missing user field
	value := map[string]any{
		"session": 123,
	}

	encoded, err := util.Cookie.Encode("semaphore", value)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  "semaphore",
		Value: encoded,
	})

	session, ok := getSession(req)

	assert.False(t, ok)
	assert.Nil(t, session)
}

func TestGetSession_ValidCookie(t *testing.T) {
	// Setup cookie encryption and mock store
	setupTestCookie()

	// Create a mock store
	mockStore := &MockStore{
		sessions: map[string]db.Session{
			"1:123": {
				ID:         123,
				UserID:     1,
				LastActive: tz.Now(),
			},
		},
	}

	value := map[string]any{
		"user":    1,
		"session": 123,
	}

	encoded, err := util.Cookie.Encode("semaphore", value)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  "semaphore",
		Value: encoded,
	})

	req = helpers.SetContextValue(req, "store", mockStore)

	session, ok := getSession(req)

	assert.True(t, ok)
	assert.NotNil(t, session)
	assert.Equal(t, 123, session.ID)
	assert.Equal(t, 1, session.UserID)
}

func TestGetSession_ExpiredSession(t *testing.T) {
	// Setup cookie encryption and mock store
	setupTestCookie()

	// Create a mock store with expired session
	expiredTime := tz.Now().Add(-8 * 24 * time.Hour) // 8 days ago
	mockStore := &MockStore{
		sessions: map[string]db.Session{
			"1:123": {
				ID:         123,
				UserID:     1,
				LastActive: expiredTime,
			},
		},
	}

	value := map[string]any{
		"user":    1,
		"session": 123,
	}

	encoded, err := util.Cookie.Encode("semaphore", value)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  "semaphore",
		Value: encoded,
	})

	req = helpers.SetContextValue(req, "store", mockStore)

	session, ok := getSession(req)

	assert.False(t, ok)
	assert.Nil(t, session)
}

func TestTypeAssertions_Safe(t *testing.T) {
	// Test that our safe type assertions work correctly
	tests := []struct {
		name     string
		value    any
		expected bool
	}{
		{
			name:     "valid int",
			value:    123,
			expected: true,
		},
		{
			name:     "int32",
			value:    int32(123),
			expected: false, // int32 is not int
		},
		{
			name:     "string",
			value:    "123",
			expected: false,
		},
		{
			name:     "float64",
			value:    float64(123),
			expected: false,
		},
		{
			name:     "nil",
			value:    nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, ok := tt.value.(int)
			if tt.expected {
				assert.True(t, ok)
				assert.Equal(t, 123, userID)
			} else {
				assert.False(t, ok)
			}
		})
	}
}

// Helper functions for testing

func setupTestCookie() {
	// Generate test cookie encryption keys
	hash := make([]byte, 32)
	for i := range hash {
		hash[i] = byte(i)
	}
	encryption := make([]byte, 32)
	for i := range encryption {
		encryption[i] = byte(i + 32)
	}

	util.CookieHash = base64.StdEncoding.EncodeToString(hash)
	util.CookieEncryption = base64.StdEncoding.EncodeToString(encryption)
	util.Cookie = securecookie.New(hash, encryption)
}

// MockStore is a simple mock implementation of db.Store for testing
type MockStore struct {
	db.Store
	sessions map[string]db.Session
	users    map[int]db.User
}

func (m *MockStore) GetSession(userID, sessionID int) (db.Session, error) {
	key := fmt.Sprintf("%d:%d", userID, sessionID)
	if session, ok := m.sessions[key]; ok {
		return session, nil
	}
	return db.Session{}, db.ErrNotFound
}

func (m *MockStore) ExpireSession(userID, sessionID int) error {
	key := fmt.Sprintf("%d:%d", userID, sessionID)
	delete(m.sessions, key)
	return nil
}

func (m *MockStore) GetUser(userID int) (db.User, error) {
	if user, ok := m.users[userID]; ok {
		return user, nil
	}
	return db.User{}, db.ErrNotFound
}

// Add other required methods as needed (stubs)
func (m *MockStore) TouchSession(userID, sessionID int) error  { return nil }
func (m *MockStore) VerifySession(userID, sessionID int) error { return nil }
func (m *MockStore) AddTotpVerification(userID int, url string, recoveryHash string) (db.UserTotp, error) {
	return db.UserTotp{}, nil
}
func (m *MockStore) DeleteTotpVerification(userID int, totpID int) error { return nil }
func (m *MockStore) AddEmailOtpVerification(userID int, code string) (db.UserEmailOtp, error) {
	return db.UserEmailOtp{}, nil
}
func (m *MockStore) DeleteEmailOtpVerification(userID int, totpID int) error { return nil }
