package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/semaphoreui/semaphore/db"
	"github.com/semaphoreui/semaphore/util"
	"github.com/semaphoreui/semaphore/api/helpers"
	log "github.com/sirupsen/logrus"
)

// HealthStatus represents the health status of the application
type HealthStatus struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Checks    map[string]Check  `json:"checks"`
}

// Check represents a health check result
type Check struct {
	Status   string        `json:"status"`
	Message  string        `json:"message,omitempty"`
	Duration time.Duration `json:"duration,omitempty"`
}

// healthHandler provides health check information
func healthHandler(w http.ResponseWriter, r *http.Request) {
	health := HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Checks:    make(map[string]Check),
	}

	// Check database
	dbCheck := checkDatabase(r)
	health.Checks["database"] = dbCheck
	if dbCheck.Status != "healthy" {
		health.Status = "unhealthy"
	}

	// Check Redis (if enabled)
	if util.Config.HA != nil && util.Config.HA.Enabled && util.Config.HA.Redis != nil {
		redisCheck := checkRedis()
		health.Checks["redis"] = redisCheck
		if redisCheck.Status != "healthy" {
			health.Status = "unhealthy"
		}
	}

	statusCode := http.StatusOK
	if health.Status != "healthy" {
		statusCode = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(health); err != nil {
		log.WithError(err).Error("Failed to encode health check response")
	}
}

// checkDatabase checks database connectivity
func checkDatabase(r *http.Request) Check {
	start := time.Now()

	// Try to get store from context
	storeVal := helpers.GetFromContext(r, "store")
	if storeVal == nil {
		return Check{
			Status:   "unknown",
			Message:  "Store not available in context",
			Duration: time.Since(start),
		}
	}

	store, ok := storeVal.(db.Store)
	if !ok || store == nil {
		return Check{
			Status:   "unknown",
			Message:  "Invalid store in context",
			Duration: time.Since(start),
		}
	}

	// Try a simple query to check database connectivity
	// Using a lightweight query that should work on all databases
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Use db.StoreSession to ensure proper session handling
	var checkErr error
	db.StoreSession(store, "health_check", func() {
		_, checkErr = store.GetUser(1) // This will return ErrNotFound if DB is working, or error if DB is down
	})

	duration := time.Since(start)

	// Check if context was cancelled (timeout)
	if ctx.Err() != nil {
		return Check{
			Status:   "unhealthy",
			Message:  "Database check timeout",
			Duration: duration,
		}
	}

	// ErrNotFound means DB is working, just user doesn't exist
	if checkErr != nil && checkErr != db.ErrNotFound {
		return Check{
			Status:   "unhealthy",
			Message:  checkErr.Error(),
			Duration: duration,
		}
	}

	return Check{
		Status:   "healthy",
		Duration: duration,
	}
}

// checkRedis checks Redis connectivity (if enabled)
func checkRedis() Check {
	start := time.Now()

	// TODO: Implement Redis health check when Redis client is available
	// For now, return healthy if Redis is configured
	// This should be implemented when Redis caching is added

	return Check{
		Status:   "healthy",
		Message:  "Redis check not yet implemented",
		Duration: time.Since(start),
	}
}

// livenessHandler provides liveness probe endpoint
// Always returns 200 if server is running
func livenessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "alive",
	})
}

// readinessHandler provides readiness probe endpoint
// Checks if application is ready to serve traffic
func readinessHandler(w http.ResponseWriter, r *http.Request) {
	health := HealthStatus{
		Status:    "ready",
		Timestamp: time.Now(),
		Checks:    make(map[string]Check),
	}

	// Check database
	dbCheck := checkDatabase(r)
	health.Checks["database"] = dbCheck
	if dbCheck.Status != "healthy" {
		health.Status = "not_ready"
	}

	statusCode := http.StatusOK
	if health.Status != "ready" {
		statusCode = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(health)
}

