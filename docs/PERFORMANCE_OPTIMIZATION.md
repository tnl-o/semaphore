# Performance Optimization Guide

This document outlines performance optimization strategies and implementations for the Semaphore UI project.

## Overview

Performance optimization is crucial for providing a responsive user experience and efficient resource utilization. This guide covers:

1. Caching strategies
2. Database query optimization
3. Pagination
4. Frontend optimization

## 1. Caching

### Current State

- Redis is available for HA (High Availability) mode
- Basic cache clearing functionality exists (`/api/cache` endpoint)
- No application-level caching for frequently accessed data

### Recommended Implementation

#### 1.1 Redis-based Caching

Use Redis for distributed caching of:
- Project lists
- Template lists
- User permissions
- Frequently accessed configuration

**Implementation Plan:**

```go
// Create cache service
type CacheService interface {
    Get(key string) ([]byte, error)
    Set(key string, value []byte, ttl time.Duration) error
    Delete(key string) error
    Clear() error
}

// Redis implementation
type RedisCache struct {
    client *redis.Client
}
```

**Cache Keys:**
- `projects:list:{user_id}` - User's project list
- `templates:list:{project_id}` - Project templates
- `user:permissions:{user_id}:{project_id}` - User permissions
- `config:app` - Application configuration

**TTL Recommendations:**
- Project lists: 5 minutes
- Template lists: 10 minutes
- User permissions: 15 minutes
- Configuration: 30 minutes

#### 1.2 Cache Invalidation

Implement cache invalidation on:
- Project creation/update/deletion
- Template creation/update/deletion
- User permission changes
- Configuration updates

### Benefits

- Reduced database load
- Faster response times
- Better scalability
- Improved user experience

## 2. Database Query Optimization

### Current State

- Queries use Squirrel query builder
- Some pagination exists (Count, Offset parameters)
- Potential N+1 query patterns in some areas

### Optimization Strategies

#### 2.1 Index Analysis

**Required Indexes:**

```sql
-- Projects
CREATE INDEX idx_project_user_id ON project__user(user_id);
CREATE INDEX idx_project_user_project_id ON project__user(project_id);

-- Tasks
CREATE INDEX idx_task_template_id ON task(template_id);
CREATE INDEX idx_task_status ON task(status);
CREATE INDEX idx_task_created ON task(created);

-- Templates
CREATE INDEX idx_template_project_id ON project__template(project_id);

-- Users
CREATE INDEX idx_user_username ON `user`(username);
CREATE INDEX idx_user_email ON `user`(email);
```

**Check Existing Indexes:**

```sql
-- MySQL
SHOW INDEXES FROM table_name;

-- PostgreSQL
SELECT * FROM pg_indexes WHERE tablename = 'table_name';
```

#### 2.2 N+1 Query Prevention

**Problem Areas:**

1. **User Permissions** - Loading permissions for each project separately
2. **Task Details** - Loading related data in loops
3. **Template References** - Loading inventory/repository data separately

**Solutions:**

1. **Batch Loading:**
```go
// Instead of:
for _, project := range projects {
    permissions := getPermissions(userID, project.ID)
}

// Use:
permissions := getPermissionsBatch(userID, projectIDs)
```

2. **Eager Loading:**
```go
// Use JOINs to load related data in single query
SELECT t.*, u.name, tpl.name 
FROM task t
JOIN user u ON t.user_id = u.id
JOIN template tpl ON t.template_id = tpl.id
```

3. **Query Result Caching:**
```go
// Cache frequently accessed query results
cacheKey := fmt.Sprintf("query:%s:%v", query, args)
if cached, err := cache.Get(cacheKey); err == nil {
    return cached
}
```

#### 2.3 Query Performance Monitoring

**Add EXPLAIN Analysis:**

```go
func (d *SqlDb) ExplainQuery(query string, args ...interface{}) error {
    explainQuery := "EXPLAIN " + query
    rows, err := d.connection.Query(explainQuery, args...)
    if err != nil {
        return err
    }
    defer rows.Close()
    
    // Log or return EXPLAIN results
    // This helps identify slow queries
    return nil
}
```

**Slow Query Logging:**

Enable slow query logging in database configuration:
- MySQL: `slow_query_log = 1`, `long_query_time = 1`
- PostgreSQL: `log_min_duration_statement = 1000`

### Benefits

- Faster query execution
- Reduced database load
- Better scalability
- Improved response times

## 3. Pagination

### Current State

- Pagination parameters exist (`Count`, `Offset`)
- Not consistently applied across all endpoints
- Some endpoints may return large datasets

### Implementation Checklist

**Endpoints Requiring Pagination:**

- [ ] `/api/projects` - Project list
- [ ] `/api/users` - User list
- [ ] `/api/project/{id}/templates` - Template list
- [ ] `/api/project/{id}/tasks` - Task list
- [ ] `/api/project/{id}/users` - Project users
- [ ] `/api/events` - Event log
- [ ] `/api/project/{id}/events` - Project events

### Pagination Best Practices

1. **Default Limits:**
   - Default: 50 items per page
   - Maximum: 200 items per page
   - Minimum: 10 items per page

2. **Response Format:**
```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "per_page": 50,
    "total": 150,
    "total_pages": 3,
    "has_next": true,
    "has_prev": false
  }
}
```

3. **Cursor-based Pagination (for large datasets):**
```json
{
  "data": [...],
  "pagination": {
    "next_cursor": "eyJpZCI6MTIzfQ",
    "has_next": true
  }
}
```

### Benefits

- Reduced memory usage
- Faster response times
- Better user experience
- Improved scalability

## 4. Frontend Optimization

### Current State

- Vue 2 application
- Vuetify UI framework
- No code splitting
- No lazy loading

### Optimization Strategies

#### 4.1 Code Splitting

**Route-based Code Splitting:**

```javascript
// Instead of:
import Projects from './views/Projects.vue'

// Use:
const Projects = () => import('./views/Projects.vue')
```

**Component-based Code Splitting:**

```javascript
// For large components
const LargeComponent = () => import('./components/LargeComponent.vue')
```

#### 4.2 Lazy Loading

**Router Configuration:**

```javascript
const routes = [
  {
    path: '/projects',
    component: () => import('./views/Projects.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/templates',
    component: () => import('./views/Templates.vue'),
    meta: { requiresAuth: true }
  }
]
```

**Dynamic Imports:**

```javascript
// Load heavy libraries on demand
async function loadChartLibrary() {
  const { Chart } = await import('chart.js')
  return Chart
}
```

#### 4.3 Bundle Size Optimization

**Strategies:**

1. **Tree Shaking:**
   - Use ES6 imports
   - Remove unused dependencies
   - Use specific imports: `import { debounce } from 'lodash-es'`

2. **Dependency Analysis:**
```bash
npm run build -- --analyze
```

3. **Optimize Images:**
   - Use WebP format
   - Implement lazy loading for images
   - Use responsive images

4. **Minification:**
   - Enable production minification
   - Use compression (gzip/brotli)

#### 4.4 Performance Monitoring

**Web Vitals:**

- Largest Contentful Paint (LCP) < 2.5s
- First Input Delay (FID) < 100ms
- Cumulative Layout Shift (CLS) < 0.1

**Tools:**

- Lighthouse
- WebPageTest
- Chrome DevTools Performance

### Benefits

- Faster page loads
- Better user experience
- Reduced bandwidth usage
- Improved SEO

## Implementation Priority

1. **High Priority:**
   - Database query optimization (indexes, N+1 fixes)
   - Pagination for all list endpoints
   - Frontend code splitting

2. **Medium Priority:**
   - Redis caching implementation
   - Frontend lazy loading
   - Bundle size optimization

3. **Low Priority:**
   - Advanced caching strategies
   - Query result caching
   - Performance monitoring

## Monitoring and Metrics

### Key Metrics to Track

- API response times (p50, p95, p99)
- Database query execution times
- Cache hit rates
- Frontend bundle sizes
- Page load times

### Tools

- Application Performance Monitoring (APM)
- Database query analyzers
- Frontend performance tools
- Log aggregation

## Testing

### Performance Testing

1. **Load Testing:**
   - Use tools like k6, Apache JMeter
   - Test with realistic data volumes
   - Monitor resource usage

2. **Stress Testing:**
   - Identify breaking points
   - Test cache invalidation
   - Test database under load

3. **Benchmarking:**
   - Before/after comparisons
   - Track improvements over time
   - Set performance budgets

## Resources

- [Vue.js Performance Best Practices](https://vuejs.org/guide/best-practices/performance.html)
- [Database Indexing Best Practices](https://use-the-index-luke.com/)
- [Redis Caching Patterns](https://redis.io/docs/manual/patterns/)
- [Web Performance Best Practices](https://web.dev/performance/)

