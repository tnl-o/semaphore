package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// RateLimiter implements a simple in-memory rate limiter
type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	limit    int
	window   time.Duration
	cleanup  *time.Ticker
}

// NewRateLimiter creates a new rate limiter
// limit: maximum number of requests allowed
// window: time window for rate limiting (e.g., 1 minute)
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
		cleanup:  time.NewTicker(5 * time.Minute), // Cleanup old entries every 5 minutes
	}

	// Start cleanup goroutine
	go rl.cleanupOldEntries()

	return rl
}

// cleanupOldEntries removes old entries from the map to prevent memory leaks
func (rl *RateLimiter) cleanupOldEntries() {
	for range rl.cleanup.C {
		rl.mu.Lock()
		now := time.Now()
		for key, times := range rl.requests {
			// Remove timestamps older than the window
			validTimes := []time.Time{}
			for _, t := range times {
				if now.Sub(t) < rl.window {
					validTimes = append(validTimes, t)
				}
			}
			if len(validTimes) == 0 {
				delete(rl.requests, key)
			} else {
				rl.requests[key] = validTimes
			}
		}
		rl.mu.Unlock()
	}
}

// Allow checks if a request from the given key should be allowed
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	// Get existing requests for this key
	times, exists := rl.requests[key]
	if !exists {
		rl.requests[key] = []time.Time{now}
		return true
	}

	// Remove old requests outside the window
	validTimes := []time.Time{}
	for _, t := range times {
		if t.After(cutoff) {
			validTimes = append(validTimes, t)
		}
	}

	// Check if we've exceeded the limit
	if len(validTimes) >= rl.limit {
		return false
	}

	// Add current request
	validTimes = append(validTimes, now)
	rl.requests[key] = validTimes
	return true
}

// GetRemaining returns the number of remaining requests for a key
func (rl *RateLimiter) GetRemaining(key string) int {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	times, exists := rl.requests[key]
	if !exists {
		return rl.limit
	}

	validCount := 0
	for _, t := range times {
		if t.After(cutoff) {
			validCount++
		}
	}

	remaining := rl.limit - validCount
	if remaining < 0 {
		return 0
	}
	return remaining
}

// GetResetTime returns when the rate limit will reset for a key
func (rl *RateLimiter) GetResetTime(key string) time.Time {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	times, exists := rl.requests[key]
	if !exists {
		return time.Now().Add(rl.window)
	}

	if len(times) == 0 {
		return time.Now().Add(rl.window)
	}

	// Find the oldest request
	oldest := times[0]
	for _, t := range times {
		if t.Before(oldest) {
			oldest = t
		}
	}

	return oldest.Add(rl.window)
}

// Stop stops the cleanup goroutine
func (rl *RateLimiter) Stop() {
	rl.cleanup.Stop()
}

// getClientIP extracts the client IP address from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (for proxies/load balancers)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		return forwarded
	}

	// Check X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	ip := r.RemoteAddr
	// Remove port if present
	if len(ip) > 0 && ip[len(ip)-1] == ']' {
		// IPv6 address
		return ip
	}
	for idx := len(ip) - 1; idx >= 0; idx-- {
		if ip[idx] == ':' {
			return ip[:idx]
		}
	}
	return ip
}

// RateLimitMiddleware creates a middleware that rate limits requests
func RateLimitMiddleware(limiter *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := getClientIP(r)
			key := clientIP

			// For authenticated endpoints, also use user ID if available
			// This allows per-user rate limiting in addition to per-IP
			if user := r.Context().Value("user"); user != nil {
				// We can't access user.ID here without importing db package
				// So we'll use IP-based limiting for now
				// In the future, this could be enhanced to use user ID
			}

			if !limiter.Allow(key) {
				remaining := limiter.GetRemaining(key)
				resetTime := limiter.GetResetTime(key)
				resetSeconds := int(time.Until(resetTime).Seconds())
				if resetSeconds < 0 {
					resetSeconds = 0
				}

				log.WithFields(log.Fields{
					"ip":       clientIP,
					"path":     r.URL.Path,
					"method":   r.Method,
					"remaining": remaining,
					"reset":    resetSeconds,
				}).Warn("Rate limit exceeded")

				w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limiter.limit))
				w.Header().Set("X-RateLimit-Remaining", "0")
				w.Header().Set("X-RateLimit-Reset", strconv.Itoa(resetSeconds))
				w.Header().Set("Retry-After", strconv.Itoa(resetSeconds))
				http.Error(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)
				return
			}

			remaining := limiter.GetRemaining(key)
			resetTime := limiter.GetResetTime(key)
			resetSeconds := int(time.Until(resetTime).Seconds())
			if resetSeconds < 0 {
				resetSeconds = 0
			}

			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limiter.limit))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
			w.Header().Set("X-RateLimit-Reset", strconv.Itoa(resetSeconds))

			next.ServeHTTP(w, r)
		})
	}
}

