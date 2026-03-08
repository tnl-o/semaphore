# Monitoring and Logging Guide

This document outlines monitoring and logging strategies for the Semaphore UI project.

## Overview

Effective monitoring and logging are essential for:
- Debugging issues
- Performance analysis
- Security auditing
- System observability
- Proactive problem detection

## 1. Structured Logging

### Current State

- Uses `logrus` for logging
- Basic structured logging with fields
- No correlation IDs for request tracking
- Limited context propagation

### Implementation Plan

#### 1.1 Correlation IDs

Add correlation IDs to track requests across services and components.

**Middleware Implementation:**

```go
// api/middleware/correlation.go
package middleware

import (
    "github.com/google/uuid"
    "net/http"
)

const CorrelationIDHeader = "X-Correlation-ID"
const CorrelationIDContextKey = "correlation_id"

func CorrelationIDMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        correlationID := r.Header.Get(CorrelationIDHeader)
        if correlationID == "" {
            correlationID = uuid.New().String()
        }
        
        // Add to response header
        w.Header().Set(CorrelationIDHeader, correlationID)
        
        // Add to request context
        ctx := context.WithValue(r.Context(), CorrelationIDContextKey, correlationID)
        r = r.WithContext(ctx)
        
        next.ServeHTTP(w, r)
    })
}
```

**Logger Helper:**

```go
// api/helpers/logger.go
package helpers

import (
    "context"
    log "github.com/sirupsen/logrus"
)

func Logger(r *http.Request) *log.Entry {
    entry := log.NewEntry(log.StandardLogger())
    
    // Add correlation ID
    if correlationID := r.Context().Value(middleware.CorrelationIDContextKey); correlationID != nil {
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
    
    return entry
}
```

#### 1.2 Enhanced Context Logging

**Usage Example:**

```go
func someHandler(w http.ResponseWriter, r *http.Request) {
    logger := helpers.Logger(r)
    
    logger.WithFields(log.Fields{
        "action": "create_project",
        "name": projectName,
    }).Info("Creating project")
    
    // All logs from this request will include correlation_id
}
```

#### 1.3 Log Levels

- **DEBUG**: Detailed information for debugging
- **INFO**: General informational messages
- **WARN**: Warning messages for potential issues
- **ERROR**: Error messages for failures
- **PANIC**: Critical errors that cause panic

### Benefits

- Request tracing across services
- Easier debugging
- Better log aggregation
- Improved troubleshooting

## 2. Prometheus Metrics

### Current State

- No metrics collection
- No Prometheus integration
- Limited performance visibility

### Implementation Plan

#### 2.1 Metrics to Collect

**HTTP Metrics:**
- Request count by endpoint, method, status code
- Request duration (histogram)
- Request size
- Response size

**Application Metrics:**
- Active tasks count
- Task execution duration
- Task success/failure rate
- Database connection pool size
- Cache hit/miss rates

**System Metrics:**
- Memory usage
- CPU usage
- Goroutine count
- Database query duration

#### 2.2 Implementation

**Add Prometheus Client:**

```go
// go.mod
require (
    github.com/prometheus/client_golang v1.19.0
)
```

**Metrics Middleware:**

```go
// api/middleware/metrics.go
package middleware

import (
    "net/http"
    "strconv"
    "time"
    
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    httpRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
    
    httpRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "endpoint"},
    )
)

func MetricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Wrap response writer to capture status code
        rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
        
        next.ServeHTTP(rw, r)
        
        duration := time.Since(start).Seconds()
        endpoint := r.URL.Path
        method := r.Method
        status := strconv.Itoa(rw.statusCode)
        
        httpRequestsTotal.WithLabelValues(method, endpoint, status).Inc()
        httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration)
    })
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}
```

**Metrics Endpoint:**

```go
// api/metrics.go
package api

import (
    "net/http"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

func metricsHandler(w http.ResponseWriter, r *http.Request) {
    promhttp.Handler().ServeHTTP(w, r)
}

// In router.go:
metricsRouter := r.Path("/metrics").Subrouter()
metricsRouter.Use(adminMiddleware) // Protect metrics endpoint
metricsRouter.Methods("GET").HandlerFunc(metricsHandler)
```

#### 2.3 Custom Metrics

**Task Metrics:**

```go
var (
    tasksTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "tasks_total",
            Help: "Total number of tasks",
        },
        []string{"project_id", "template_id", "status"},
    )
    
    taskDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "task_duration_seconds",
            Help:    "Task execution duration in seconds",
            Buckets: []float64{1, 5, 10, 30, 60, 300, 600, 1800, 3600},
        },
        []string{"project_id", "template_id"},
    )
)
```

### Benefits

- Performance monitoring
- Alerting capabilities
- Capacity planning
- Trend analysis

## 3. Health Checks

### Current State

- Basic `/api/ping` endpoint returns "pong"
- No health checks for dependencies
- No readiness/liveness probes

### Implementation Plan

#### 3.1 Enhanced Health Check

**Health Check Response:**

```go
// api/health.go
package api

import (
    "encoding/json"
    "net/http"
    "time"
    
    "github.com/semaphoreui/semaphore/db"
    "github.com/semaphoreui/semaphore/util"
)

type HealthStatus struct {
    Status    string            `json:"status"`
    Timestamp time.Time         `json:"timestamp"`
    Checks    map[string]Check  `json:"checks"`
}

type Check struct {
    Status      string        `json:"status"`
    Message     string        `json:"message,omitempty"`
    Duration    time.Duration `json:"duration,omitempty"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    health := HealthStatus{
        Status:    "healthy",
        Timestamp: time.Now(),
        Checks:    make(map[string]Check),
    }
    
    // Check database
    dbCheck := checkDatabase()
    health.Checks["database"] = dbCheck
    if dbCheck.Status != "healthy" {
        health.Status = "unhealthy"
    }
    
    // Check Redis (if enabled)
    if util.Config.HA != nil && util.Config.HA.Enabled {
        redisCheck := checkRedis()
        health.Checks["redis"] = redisCheck
        if redisCheck.Status != "healthy" {
            health.Status = "unhealthy"
        }
    }
    
    // Check disk space
    diskCheck := checkDiskSpace()
    health.Checks["disk"] = diskCheck
    if diskCheck.Status != "healthy" {
        health.Status = "unhealthy"
    }
    
    statusCode := http.StatusOK
    if health.Status != "healthy" {
        statusCode = http.StatusServiceUnavailable
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(health)
}

func checkDatabase() Check {
    start := time.Now()
    
    store := db.StoreSession(nil, "health_check", func() db.Store {
        // Get store instance
        return nil // Replace with actual store
    })
    
    // Simple query to check database
    _, err := store.GetUser(1) // Or a simpler query
    
    duration := time.Since(start)
    
    if err != nil && err != db.ErrNotFound {
        return Check{
            Status:   "unhealthy",
            Message:  err.Error(),
            Duration: duration,
        }
    }
    
    return Check{
        Status:   "healthy",
        Duration: duration,
    }
}

func checkRedis() Check {
    start := time.Now()
    
    // Check Redis connection
    // Implementation depends on Redis client
    
    duration := time.Since(start)
    
    return Check{
        Status:   "healthy",
        Duration: duration,
    }
}

func checkDiskSpace() Check {
    // Check available disk space
    // Implementation depends on requirements
    
    return Check{
        Status: "healthy",
    }
}
```

#### 3.2 Separate Endpoints

**Liveness Probe (`/health/live`):**
- Checks if application is running
- Always returns 200 if server is up

**Readiness Probe (`/health/ready`):**
- Checks if application is ready to serve traffic
- Checks dependencies (DB, Redis, etc.)
- Returns 503 if not ready

**Startup Probe (`/health/startup`):**
- Checks if application has finished startup
- Useful for slow-starting applications

### Benefits

- Better orchestration support (Kubernetes)
- Dependency monitoring
- Faster failure detection
- Improved reliability

## 4. Distributed Tracing

### Current State

- No distributed tracing
- No request tracing across services
- Limited visibility into request flow

### Implementation Plan

#### 4.1 OpenTelemetry Integration

**Add OpenTelemetry:**

```go
// go.mod
require (
    go.opentelemetry.io/otel v1.24.0
    go.opentelemetry.io/otel/trace v1.24.0
    go.opentelemetry.io/otel/exporters/jaeger v1.17.0
    go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.49.0
)
```

**Tracing Setup:**

```go
// util/tracing.go
package util

import (
    "context"
    
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/jaeger"
    "go.opentelemetry.io/otel/propagation"
    "go.opentelemetry.io/otel/sdk/resource"
    "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func InitTracing(serviceName, jaegerEndpoint string) (*trace.TracerProvider, error) {
    exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerEndpoint)))
    if err != nil {
        return nil, err
    }
    
    tp := trace.NewTracerProvider(
        trace.WithBatcher(exp),
        trace.WithResource(resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceNameKey.String(serviceName),
        )),
    )
    
    otel.SetTracerProvider(tp)
    otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
        propagation.TraceContext{},
        propagation.Baggage{},
    ))
    
    return tp, nil
}
```

**Tracing Middleware:**

```go
// api/middleware/tracing.go
package middleware

import (
    "net/http"
    "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func TracingMiddleware(next http.Handler) http.Handler {
    return otelhttp.NewHandler(next, "semaphore-api")
}
```

### Benefits

- Request flow visibility
- Performance bottleneck identification
- Distributed system debugging
- Service dependency mapping

## Implementation Priority

1. **High Priority:**
   - Correlation IDs
   - Enhanced health checks
   - Basic Prometheus metrics

2. **Medium Priority:**
   - Advanced Prometheus metrics
   - Custom application metrics
   - OpenTelemetry tracing

3. **Low Priority:**
   - Advanced tracing features
   - Log aggregation setup
   - Metrics dashboards

## Resources

- [Logrus Documentation](https://github.com/sirupsen/logrus)
- [Prometheus Client Go](https://github.com/prometheus/client_golang)
- [OpenTelemetry Go](https://opentelemetry.io/docs/instrumentation/go/)
- [Jaeger Documentation](https://www.jaegertracing.io/docs/)

