package api

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsHandler returns Prometheus metrics endpoint handler
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}
