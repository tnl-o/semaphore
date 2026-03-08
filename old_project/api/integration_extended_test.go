package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/semaphoreui/semaphore/api/helpers"
	"github.com/semaphoreui/semaphore/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBind_ValidJSON(t *testing.T) {
	type TestStruct struct {
		Name  string `json:"name" validate:"required,min=3,max=50"`
		Email string `json:"email" validate:"required,email"`
		Age   int    `json:"age" validate:"required,min=18,max=120"`
	}

	body := TestStruct{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   25,
	}

	jsonBody, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	var result TestStruct
	success := helpers.Bind(w, req, &result)

	assert.True(t, success)
	assert.Equal(t, body.Name, result.Name)
	assert.Equal(t, body.Email, result.Email)
	assert.Equal(t, body.Age, result.Age)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBind_InvalidJSON(t *testing.T) {
	type TestStruct struct {
		Name string `json:"name" validate:"required"`
	}

	invalidJSON := `{"name": "test" invalid}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	var result TestStruct
	success := helpers.Bind(w, req, &result)

	assert.False(t, success)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBind_ValidationErrors(t *testing.T) {
	type TestStruct struct {
		Name  string `json:"name" validate:"required,min=3,max=50"`
		Email string `json:"email" validate:"required,email"`
		Age   int    `json:"age" validate:"required,min=18"`
	}

	tests := []struct {
		name        string
		body        TestStruct
		expectError bool
		description string
	}{
		{
			name: "missing required field",
			body: TestStruct{
				Email: "test@example.com",
				Age:   25,
			},
			expectError: true,
			description: "should fail when required field is missing",
		},
		{
			name: "invalid email",
			body: TestStruct{
				Name:  "John Doe",
				Email: "invalid-email",
				Age:   25,
			},
			expectError: true,
			description: "should fail with invalid email format",
		},
		{
			name: "value too short",
			body: TestStruct{
				Name:  "Jo", // Less than min=3
				Email: "test@example.com",
				Age:   25,
			},
			expectError: true,
			description: "should fail when value is too short",
		},
		{
			name: "age too low",
			body: TestStruct{
				Name:  "John Doe",
				Email: "test@example.com",
				Age:   15, // Less than min=18
			},
			expectError: true,
			description: "should fail when age is below minimum",
		},
		{
			name: "valid data",
			body: TestStruct{
				Name:  "John Doe",
				Email: "test@example.com",
				Age:   25,
			},
			expectError: false,
			description: "should succeed with valid data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, err := json.Marshal(tt.body)
			require.NoError(t, err)

			req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			var result TestStruct
			success := helpers.Bind(w, req, &result)

			if tt.expectError {
				assert.False(t, success, tt.description)
				assert.Equal(t, http.StatusBadRequest, w.Code)
			} else {
				assert.True(t, success, tt.description)
			}
		})
	}
}

func TestBind_NoValidationTags(t *testing.T) {
	// Test that structs without validation tags still work
	type TestStruct struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	body := TestStruct{
		Name:  "John Doe",
		Email: "test@example.com",
	}

	jsonBody, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	var result TestStruct
	success := helpers.Bind(w, req, &result)

	assert.True(t, success)
	assert.Equal(t, body.Name, result.Name)
	assert.Equal(t, body.Email, result.Email)
}

func TestIntegrationMatch_HeaderMatch(t *testing.T) {
	body := []byte(`{}`)
	header := make(http.Header)
	header.Set("X-Custom-Header", "test-value")

	matched := Match(db.IntegrationMatcher{
		ID:            0,
		Name:          "Test",
		IntegrationID: 0,
		MatchType:     db.IntegrationMatchHeader,
		Method:        db.IntegrationMatchMethodEquals,
		Key:           "X-Custom-Header",
		Value:         "test-value",
	}, header, body)

	assert.True(t, matched)
}

func TestIntegrationMatch_HeaderMismatch(t *testing.T) {
	body := []byte(`{}`)
	header := make(http.Header)
	header.Set("X-Custom-Header", "different-value")

	matched := Match(db.IntegrationMatcher{
		ID:            0,
		Name:          "Test",
		IntegrationID: 0,
		MatchType:     db.IntegrationMatchHeader,
		Method:        db.IntegrationMatchMethodEquals,
		Key:           "X-Custom-Header",
		Value:         "test-value",
	}, header, body)

	assert.False(t, matched)
}

func TestIntegrationMatch_BodyContains(t *testing.T) {
	body := []byte(`{"hook_id": 4856239453, "action": "opened"}`)
	header := make(http.Header)

	matched := Match(db.IntegrationMatcher{
		ID:            0,
		Name:          "Test",
		IntegrationID: 0,
		MatchType:     db.IntegrationMatchBody,
		Method:        db.IntegrationMatchMethodContains,
		BodyDataType:  db.IntegrationBodyDataJSON,
		Key:           "action",
		Value:         "opened",
	}, header, body)

	assert.True(t, matched)
}

func TestIntegrationMatch_InvalidJSON(t *testing.T) {
	body := []byte(`{invalid json}`)
	header := make(http.Header)

	matched := Match(db.IntegrationMatcher{
		ID:            0,
		Name:          "Test",
		IntegrationID: 0,
		MatchType:     db.IntegrationMatchBody,
		Method:        db.IntegrationMatchMethodEquals,
		BodyDataType:  db.IntegrationBodyDataJSON,
		Key:           "hook_id",
		Value:         "4856239453",
	}, header, body)

	assert.False(t, matched)
}

func TestExtract_MissingKey(t *testing.T) {
	payload := []byte(`{"branch": "main"}`)
	header := make(http.Header)
	header.Set("X-Event", "push")

	extractValues := []db.IntegrationExtractValue{
		{
			Variable:     "MISSING_VAR",
			ValueSource:  db.IntegrationExtractBodyValue,
			BodyDataType: db.IntegrationBodyDataJSON,
			Key:          "nonexistent.key",
			VariableType: db.IntegrationVariableEnvironment,
		},
	}

	result := Extract(extractValues, header, payload)

	// Missing keys should not cause errors, just return empty result
	assert.Empty(t, result["MISSING_VAR"])
}

func TestExtract_EmptyExtractValues(t *testing.T) {
	payload := []byte(`{"branch": "main"}`)
	header := make(http.Header)

	extractValues := []db.IntegrationExtractValue{}

	result := Extract(extractValues, header, payload)

	assert.Empty(t, result)
}

func TestGetTaskDefinition_EmptyPayload(t *testing.T) {
	integration := db.Integration{
		ID:         1,
		ProjectID:  1,
		TemplateID: 1,
		TaskParams: &db.TaskParams{
			ProjectID:   1,
			Environment: `{"existing":"value"}`,
			Params:      db.MapStringAnyField{},
		},
	}

	header := make(http.Header)
	payload := []byte(`{}`)

	task, err := GetTaskDefinition(integration, payload, header, func(projectID, integrationID int) ([]db.IntegrationExtractValue, error) {
		return []db.IntegrationExtractValue{}, nil
	})

	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, integration.ProjectID, task.ProjectID)
	assert.Equal(t, integration.TemplateID, task.TemplateID)
}

func TestGetTaskDefinition_MissingTaskParams(t *testing.T) {
	integration := db.Integration{
		ID:         1,
		ProjectID:  1,
		TemplateID: 1,
		// TaskParams is nil
	}

	header := make(http.Header)
	payload := []byte(`{}`)

	task, err := GetTaskDefinition(integration, payload, header, func(projectID, integrationID int) ([]db.IntegrationExtractValue, error) {
		return []db.IntegrationExtractValue{}, nil
	})

	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, integration.ProjectID, task.ProjectID)
	assert.Equal(t, integration.TemplateID, task.TemplateID)
}

func TestExtract_NestedJSONPath(t *testing.T) {
	payload := []byte(`{"repository": {"owner": {"name": "octocat"}, "name": "Hello-World"}}`)
	header := make(http.Header)

	extractValues := []db.IntegrationExtractValue{
		{
			Variable:     "REPO_OWNER",
			ValueSource:  db.IntegrationExtractBodyValue,
			BodyDataType: db.IntegrationBodyDataJSON,
			Key:          "repository.owner.name",
			VariableType: db.IntegrationVariableEnvironment,
		},
		{
			Variable:     "REPO_NAME",
			ValueSource:  db.IntegrationExtractBodyValue,
			BodyDataType: db.IntegrationBodyDataJSON,
			Key:          "repository.name",
			VariableType: db.IntegrationVariableEnvironment,
		},
	}

	result := Extract(extractValues, header, payload)

	assert.Equal(t, "octocat", result["REPO_OWNER"])
	assert.Equal(t, "Hello-World", result["REPO_NAME"])
}

func TestExtract_ArrayIndex(t *testing.T) {
	payload := []byte(`{"commits": [{"id": "abc123", "message": "First commit"}, {"id": "def456", "message": "Second commit"}]}`)
	header := make(http.Header)

	extractValues := []db.IntegrationExtractValue{
		{
			Variable:     "FIRST_COMMIT_ID",
			ValueSource:  db.IntegrationExtractBodyValue,
			BodyDataType: db.IntegrationBodyDataJSON,
			Key:          "commits.0.id",
			VariableType: db.IntegrationVariableEnvironment,
		},
		{
			Variable:     "SECOND_COMMIT_MESSAGE",
			ValueSource:  db.IntegrationExtractBodyValue,
			BodyDataType: db.IntegrationBodyDataJSON,
			Key:          "commits.1.message",
			VariableType: db.IntegrationVariableEnvironment,
		},
	}

	result := Extract(extractValues, header, payload)

	assert.Equal(t, "abc123", result["FIRST_COMMIT_ID"])
	assert.Equal(t, "Second commit", result["SECOND_COMMIT_MESSAGE"])
}

func TestExtract_HeaderCaseInsensitive(t *testing.T) {
	payload := []byte(`{}`)
	header := make(http.Header)
	header.Set("X-Custom-Header", "test-value")
	header.Set("content-type", "application/json")

	extractValues := []db.IntegrationExtractValue{
		{
			Variable:     "CUSTOM_HEADER",
			ValueSource:  db.IntegrationExtractHeaderValue,
			Key:          "x-custom-header", // Lowercase
			VariableType: db.IntegrationVariableEnvironment,
		},
		{
			Variable:     "CONTENT_TYPE",
			ValueSource:  db.IntegrationExtractHeaderValue,
			Key:          "Content-Type", // Mixed case
			VariableType: db.IntegrationVariableEnvironment,
		},
	}

	result := Extract(extractValues, header, payload)

	// Header extraction should be case-insensitive
	assert.Equal(t, "test-value", result["CUSTOM_HEADER"])
	assert.Equal(t, "application/json", result["CONTENT_TYPE"])
}
