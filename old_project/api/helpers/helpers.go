package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/semaphoreui/semaphore/db"
	log "github.com/sirupsen/logrus"
)

func Store(r *http.Request) db.Store {
	return GetFromContext(r, "store").(db.Store)
}

func isXHR(w http.ResponseWriter, r *http.Request) bool {
	accept := r.Header.Get("Accept")
	return !strings.Contains(accept, "text/html")
}

// H just a string-to-anything map
type H map[string]any

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Bind decodes json into object and validates it
func Bind(w http.ResponseWriter, r *http.Request, out any) bool {
	err := json.NewDecoder(r.Body).Decode(out)
	if err != nil {
		log.WithError(err).Debug("Failed to decode JSON request body")
		WriteErrorStatus(w, "Invalid JSON format", http.StatusBadRequest)
		return false
	}

	// Validate the struct if it has validation tags
	// Only validate if the struct actually has validation tags to avoid false positives
	if err := validate.Struct(out); err != nil {
		// Check if this is a validation error or just a struct without validation tags
		if validationErrors, ok := err.(validator.ValidationErrors); ok && len(validationErrors) > 0 {
			validationErrorMsg := formatValidationErrors(err)
			log.WithFields(log.Fields{
				"errors": validationErrorMsg,
				"path":   r.URL.Path,
			}).Debug("Validation failed for request")
			WriteErrorStatus(w, fmt.Sprintf("Validation failed: %s", validationErrorMsg), http.StatusBadRequest)
			return false
		}
		// If it's not a validation error, ignore it (struct might not have validation tags)
	}

	return true
}

// formatValidationErrors formats validator errors into a readable string
func formatValidationErrors(err error) string {
	var errors []string
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errors = append(errors, fmt.Sprintf("%s: %s", e.Field(), getValidationErrorMessage(e)))
		}
	} else {
		errors = append(errors, err.Error())
	}
	return strings.Join(errors, ", ")
}

// getValidationErrorMessage returns a user-friendly error message for a validation error
func getValidationErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return fmt.Sprintf("must be at least %s characters", e.Param())
	case "max":
		return fmt.Sprintf("must be at most %s characters", e.Param())
	case "len":
		return fmt.Sprintf("must be exactly %s characters", e.Param())
	case "alphanum":
		return "must contain only alphanumeric characters"
	case "alpha":
		return "must contain only alphabetic characters"
	case "numeric":
		return "must be numeric"
	case "url":
		return "must be a valid URL"
	case "oneof":
		return fmt.Sprintf("must be one of: %s", e.Param())
	default:
		return fmt.Sprintf("failed validation for tag '%s'", e.Tag())
	}
}
