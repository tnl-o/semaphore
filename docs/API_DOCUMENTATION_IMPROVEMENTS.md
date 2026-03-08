# API Documentation Improvements Plan

This document outlines planned improvements to the API documentation.

## Current State

- API documentation exists in `api-docs.yml` (Swagger/OpenAPI format)
- Basic endpoint descriptions
- Limited examples
- Missing error response documentation

## Planned Improvements

### 1. Add Request/Response Examples

**Priority: High**

Add comprehensive examples for all endpoints:

```yaml
paths:
  /api/projects:
    post:
      summary: Create a new project
      requestBody:
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/Project'
            examples:
              basic:
                summary: Basic project
                value:
                  name: "My Project"
                  alert: false
              with_alert:
                summary: Project with alerts
                value:
                  name: "Production"
                  alert: true
      responses:
        '201':
          description: Project created
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Project'
              examples:
                created:
                  summary: Created project
                  value:
                    id: 1
                    name: "My Project"
                    created: "2024-01-01T00:00:00Z"
```

### 2. Error Response Documentation

**Priority: High**

Document all possible error responses:

```yaml
components:
  responses:
    BadRequest:
      description: Bad request
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/Error'
          examples:
            validation_error:
              summary: Validation error
              value:
                error: "Validation failed: username is required"
            invalid_json:
              summary: Invalid JSON
              value:
                error: "Invalid JSON format"
    
    Unauthorized:
      description: Unauthorized
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/Error'
          examples:
            no_session:
              summary: No session
              value:
                error: "Not authenticated"
            expired_session:
              summary: Expired session
              value:
                error: "Session expired"
    
    Forbidden:
      description: Forbidden
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/Error'
          examples:
            insufficient_permissions:
              summary: Insufficient permissions
              value:
                error: "User does not have permission to perform this action"
    
    NotFound:
      description: Resource not found
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/Error'
          examples:
            project_not_found:
              summary: Project not found
              value:
                error: "Project not found"
    
    InternalServerError:
      description: Internal server error
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/Error'
          examples:
            server_error:
              summary: Server error
              value:
                error: "Internal server error"
```

### 3. Schema Improvements

**Priority: Medium**

Add detailed schema descriptions:

```yaml
components:
  schemas:
    Project:
      type: object
      required:
        - name
      properties:
        id:
          type: integer
          description: Unique project identifier
          example: 1
        name:
          type: string
          description: Project name
          minLength: 1
          maxLength: 255
          example: "My Project"
        created:
          type: string
          format: date-time
          description: Project creation timestamp
          example: "2024-01-01T00:00:00Z"
        alert:
          type: boolean
          description: Enable alerts for this project
          default: false
```

### 4. Authentication Documentation

**Priority: High**

Document authentication flow:

```yaml
components:
  securitySchemes:
    cookieAuth:
      type: apiKey
      in: cookie
      name: semaphore
      description: |
        Authentication is performed via secure cookies.
        
        To authenticate:
        1. POST to /api/auth/login with credentials
        2. Server sets secure cookie
        3. Include cookie in subsequent requests
        
        Example:
        ```bash
        curl -X POST http://localhost:3000/api/auth/login \
          -H "Content-Type: application/json" \
          -d '{"auth":"admin","password":"changeme"}' \
          -c cookies.txt
        
        curl http://localhost:3000/api/projects -b cookies.txt
        ```
```

### 5. Pagination Documentation

**Priority: Medium**

Document pagination parameters and responses:

```yaml
components:
  parameters:
    page:
      name: page
      in: query
      description: Page number (1-based)
      schema:
        type: integer
        minimum: 1
        default: 1
        example: 1
    
    per_page:
      name: per_page
      in: query
      description: Number of items per page
      schema:
        type: integer
        minimum: 1
        maximum: 200
        default: 50
        example: 50
  
  schemas:
    PaginatedResponse:
      type: object
      properties:
        data:
          type: array
          description: Array of items
        pagination:
          $ref: '#/components/schemas/Pagination'
    
    Pagination:
      type: object
      properties:
        page:
          type: integer
          description: Current page number
        per_page:
          type: integer
          description: Items per page
        total:
          type: integer
          description: Total number of items
        total_pages:
          type: integer
          description: Total number of pages
        has_next:
          type: boolean
          description: Whether there is a next page
        has_prev:
          type: boolean
          description: Whether there is a previous page
```

### 6. WebSocket Documentation

**Priority: Low**

Document WebSocket API:

```yaml
paths:
  /api/ws:
    get:
      summary: WebSocket connection
      description: |
        Establishes a WebSocket connection for real-time updates.
        
        Messages received:
        - Task output updates
        - Task status changes
        - System notifications
        
        Message format:
        ```json
        {
          "type": "log",
          "output": "Task output line",
          "time": "2024-01-01T00:00:00Z",
          "task_id": 1,
          "project_id": 1
        }
        ```
      security:
        - cookieAuth: []
```

### 7. Rate Limiting Documentation

**Priority: Medium**

Document rate limiting:

```yaml
components:
  headers:
    X-RateLimit-Limit:
      schema:
        type: integer
      description: Request limit per window
    X-RateLimit-Remaining:
      schema:
        type: integer
      description: Remaining requests in current window
    X-RateLimit-Reset:
      schema:
        type: integer
      description: Unix timestamp when rate limit resets
```

### 8. Code Examples

**Priority: High**

Add code examples for common languages:

```yaml
paths:
  /api/projects:
    get:
      summary: List projects
      x-code-samples:
        - lang: 'curl'
          source: |
            curl -X GET http://localhost:3000/api/projects \
              -H "Cookie: semaphore=..."
        - lang: 'javascript'
          source: |
            const response = await fetch('http://localhost:3000/api/projects', {
              credentials: 'include'
            });
            const projects = await response.json();
        - lang: 'python'
          source: |
            import requests
            
            session = requests.Session()
            session.post('http://localhost:3000/api/auth/login', json={
                'auth': 'admin',
                'password': 'changeme'
            })
            
            response = session.get('http://localhost:3000/api/projects')
            projects = response.json()
```

## Implementation Checklist

- [ ] Add request examples for all POST/PUT endpoints
- [ ] Add response examples for all endpoints
- [ ] Document all error responses
- [ ] Add schema descriptions
- [ ] Document authentication flow
- [ ] Document pagination
- [ ] Document WebSocket API
- [ ] Document rate limiting
- [ ] Add code examples (curl, JavaScript, Python)
- [ ] Validate documentation with Swagger UI
- [ ] Test all examples

## Tools

- **Swagger Editor**: For editing and validating OpenAPI spec
- **Swagger UI**: For interactive API documentation
- **Dredd**: For validating API documentation against implementation

## Resources

- [OpenAPI Specification](https://swagger.io/specification/)
- [Swagger Editor](https://editor.swagger.io/)
- [API Documentation Best Practices](https://swagger.io/resources/articles/adopting-an-api-first-approach/)

