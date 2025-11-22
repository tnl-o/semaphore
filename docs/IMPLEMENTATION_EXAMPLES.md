# Implementation Examples

This document provides examples of how to use the newly implemented features.

## 1. Using Logger Helper

The `helpers.Logger` function provides structured logging with correlation IDs, user ID, and project ID automatically included.

### Example: Login Handler

```go
// api/login.go
func login(w http.ResponseWriter, r *http.Request) {
    logger := helpers.Logger(r)
    
    logger.Info("Login attempt")
    
    var loginData login
    if !helpers.Bind(w, r, &loginData) {
        logger.WithField("error", "validation_failed").Warn("Login validation failed")
        return
    }
    
    logger.WithFields(log.Fields{
        "username": loginData.Username,
        "method":   loginData.Method,
    }).Debug("Processing login")
    
    // ... login logic ...
    
    logger.WithField("user_id", user.ID).Info("Login successful")
}
```

### Example: Project Handler

```go
// api/projects/project.go
func (c *ProjectController) GetProject(w http.ResponseWriter, r *http.Request) {
    logger := helpers.Logger(r)
    project := helpers.GetFromContext(r, "project").(*db.Project)
    
    logger.WithFields(log.Fields{
        "project_id": project.ID,
        "project_name": project.Name,
    }).Info("Retrieving project")
    
    // ... handler logic ...
}
```

## 2. Using Toast Notifications

The `toast` helper provides convenient methods for showing notifications.

### Example: API Error Handling

```javascript
// web/src/components/SomeComponent.vue
import { toast } from '@/lib/toast';
import axios from 'axios';

export default {
  methods: {
    async saveData() {
      try {
        await axios.post('/api/some-endpoint', this.data);
        toast.success('Data saved successfully');
      } catch (err) {
        toast.apiError(err); // Automatically uses getErrorMessage
      }
    },
  },
};
```

### Example: Success/Error Notifications

```javascript
import { toast } from '@/lib/toast';

// Success notification
toast.success('Operation completed successfully');

// Error notification
toast.error('Something went wrong');

// Info notification
toast.info('Processing your request...');

// Warning notification
toast.warning('Please check your input');
```

### Example: Replacing EventBus Usage

**Before:**
```javascript
EventBus.$emit('i-snackbar', {
  color: 'error',
  text: getErrorMessage(err),
});
```

**After:**
```javascript
import { toast } from '@/lib/toast';

toast.apiError(err);
```

## 3. Using Loading Mixin

The `LoadingMixin` provides loading state management for components.

### Example: Component with Loading State

```vue
<template>
  <div>
    <v-progress-linear
      v-if="loading"
      indeterminate
      color="primary"
    ></v-progress-linear>
    
    <v-btn
      @click="loadData"
      :disabled="loading"
    >
      {{ loading ? loadingMessage || 'Loading...' : 'Load Data' }}
    </v-btn>
    
    <div v-if="!loading && data">
      <!-- Display data -->
    </div>
  </div>
</template>

<script>
import LoadingMixin from '@/mixins/LoadingMixin';
import axios from 'axios';

export default {
  mixins: [LoadingMixin],
  data() {
    return {
      data: null,
    };
  },
  methods: {
    async loadData() {
      await this.withLoading(async () => {
        const response = await axios.get('/api/data');
        this.data = response.data;
      }, 'Loading data...');
    },
  },
};
</script>
```

### Example: Manual Loading Control

```javascript
export default {
  mixins: [LoadingMixin],
  methods: {
    async complexOperation() {
      this.setLoading(true, 'Step 1: Initializing...');
      await this.step1();
      
      this.setLoading(true, 'Step 2: Processing...');
      await this.step2();
      
      this.setLoading(false);
    },
  },
};
```

## 4. Prometheus Metrics

Metrics are automatically collected by the `MetricsMiddleware`. To access metrics:

### Viewing Metrics

```bash
# Get Prometheus metrics
curl http://localhost:3000/api/metrics
```

### Metrics Available

- `http_request_duration_seconds` - Request duration histogram
- `http_requests_total` - Total request count
- `http_requests_in_flight` - Current in-flight requests

### Example Prometheus Query

```promql
# Request rate per second
rate(http_requests_total[5m])

# 95th percentile latency
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))

# Error rate
rate(http_requests_total{status_code=~"5.."}[5m])
```

## 5. Health Check Endpoints

### Basic Health Check

```bash
# Full health check
curl http://localhost:3000/api/health

# Liveness probe (always returns 200 if server is running)
curl http://localhost:3000/api/health/live

# Readiness probe (checks if ready to serve traffic)
curl http://localhost:3000/api/health/ready
```

### Example Response

```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "checks": {
    "database": {
      "status": "healthy",
      "duration": "2.5ms"
    }
  }
}
```

## 6. Correlation IDs

Correlation IDs are automatically added to all requests. They can be used for request tracing.

### Using Correlation ID in Logs

```go
func someHandler(w http.ResponseWriter, r *http.Request) {
    logger := helpers.Logger(r)
    // Correlation ID is automatically included in all log entries
    logger.Info("Processing request")
}
```

### Getting Correlation ID

```go
import "github.com/semaphoreui/semaphore/api/middleware"

correlationID := middleware.GetCorrelationID(r)
```

### Client-Side: Including Correlation ID

```javascript
// Include correlation ID in request headers
axios.interceptors.request.use((config) => {
  const correlationID = localStorage.getItem('correlation_id') || 
    generateUUID();
  localStorage.setItem('correlation_id', correlationID);
  config.headers['X-Correlation-ID'] = correlationID;
  return config;
});
```

## 7. Error Handling Best Practices

### Backend: Using Logger for Errors

```go
func someHandler(w http.ResponseWriter, r *http.Request) {
    logger := helpers.Logger(r)
    
    if err := doSomething(); err != nil {
        logger.WithError(err).Error("Failed to do something")
        helpers.WriteError(w, err)
        return
    }
    
    logger.Info("Operation successful")
}
```

### Frontend: Using Toast for Errors

```javascript
import { toast } from '@/lib/toast';
import { getErrorMessage } from '@/lib/error';

try {
  await axios.post('/api/endpoint', data);
  toast.success('Success!');
} catch (err) {
  // Option 1: Use toast helper
  toast.apiError(err);
  
  // Option 2: Manual error handling
  const message = getErrorMessage(err);
  toast.error(message);
}
```

## 8. Integration Example: Complete Handler

```go
func createProject(w http.ResponseWriter, r *http.Request) {
    logger := helpers.Logger(r)
    logger.Info("Creating new project")
    
    var projectData projectRequest
    if !helpers.Bind(w, r, &projectData) {
        logger.Warn("Invalid project data")
        return
    }
    
    store := helpers.Store(r)
    user := helpers.UserFromContext(r)
    
    logger.WithFields(log.Fields{
        "user_id": user.ID,
        "project_name": projectData.Name,
    }).Debug("Creating project")
    
    project, err := store.CreateProject(db.Project{
        Name: projectData.Name,
        // ... other fields
    })
    
    if err != nil {
        logger.WithError(err).Error("Failed to create project")
        helpers.WriteError(w, err)
        return
    }
    
    logger.WithField("project_id", project.ID).Info("Project created successfully")
    helpers.WriteJSON(w, http.StatusCreated, project)
}
```

## 9. Testing Examples

### Testing with Correlation ID

```go
func TestHandler(t *testing.T) {
    req := httptest.NewRequest("GET", "/api/test", nil)
    req.Header.Set("X-Correlation-ID", "test-correlation-id")
    
    // Handler will use the provided correlation ID
    // or generate a new one if not provided
}
```

### Testing Metrics

```go
func TestMetrics(t *testing.T) {
    // Metrics are automatically collected
    // Can be verified by checking Prometheus metrics endpoint
    resp := httptest.NewRecorder()
    req := httptest.NewRequest("GET", "/api/metrics", nil)
    
    MetricsHandler().ServeHTTP(resp, req)
    
    // Verify metrics are exposed
    assert.Contains(t, resp.Body.String(), "http_requests_total")
}
```

