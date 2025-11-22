# Semaphore UI Architecture

This document describes the architecture of the Semaphore UI project, including its components, data flows, and design decisions.

## Overview

Semaphore UI is a modern web interface for managing DevOps tools like Ansible, Terraform, PowerShell, and Bash scripts. It consists of a Go backend API server and a Vue.js frontend application.

## High-Level Architecture

```
┌─────────────────┐
│   Web Browser   │
│  (Vue.js SPA)   │
└────────┬────────┘
         │ HTTP/WebSocket
         │
┌────────▼─────────────────────────────────────┐
│         Go API Server                        │
│  ┌──────────────────────────────────────┐   │
│  │  HTTP Router (Gorilla Mux)           │   │
│  │  ┌──────────┐  ┌──────────┐         │   │
│  │  │  Auth    │  │ Projects │  ...    │   │
│  │  └──────────┘  └──────────┘         │   │
│  └──────────────────────────────────────┘   │
│  ┌──────────────────────────────────────┐   │
│  │  Services Layer                      │   │
│  │  - Task Pool                         │   │
│  │  - Project Service                   │   │
│  │  - Integration Service               │   │
│  └──────────────────────────────────────┘   │
│  ┌──────────────────────────────────────┐   │
│  │  Database Layer                      │   │
│  │  - Store Interface                   │   │
│  │  - SQL/Bolt Implementations          │   │
│  └──────────────────────────────────────┘   │
└─────────────────────────────────────────────┘
         │
         │
┌────────▼────────┐  ┌──────────────┐
│   Database      │  │   Redis      │
│ (MySQL/Postgres │  │  (Optional)  │
│  /SQLite/Bolt)  │  │              │
└─────────────────┘  └──────────────┘
```

## Component Architecture

### Backend Components

#### 1. API Layer (`api/`)

The API layer handles HTTP requests and responses.

**Key Components:**
- **Router** (`api/router.go`): Defines all API routes and applies middleware
- **Handlers**: Request handlers for different resources
  - `api/auth.go`: Authentication and session management
  - `api/login.go`: Login and LDAP authentication
  - `api/projects/`: Project-related endpoints
  - `api/runners.go`: Runner management
  - `api/tasks/`: Task management
- **Middleware** (`api/middleware/`):
  - Rate limiting
  - Authentication
  - JSON response formatting
  - CORS handling
- **Helpers** (`api/helpers/`):
  - Context management
  - JSON binding and validation
  - Response writing
  - Query parameter parsing

**Request Flow:**
```
HTTP Request
    ↓
Middleware Chain (CORS, Auth, Rate Limit, etc.)
    ↓
Route Handler
    ↓
Service Layer
    ↓
Database Layer
    ↓
Response
```

#### 2. Services Layer (`services/`)

Business logic and orchestration.

**Key Services:**
- **Task Pool** (`services/tasks/`):
  - Manages task execution queue
  - Handles task lifecycle (queued → running → completed)
  - Coordinates with runners
- **Project Service** (`services/project/`):
  - Project management logic
  - Resource organization
- **Server Services** (`services/server/`):
  - Integration service
  - Access key encryption
  - Secret storage
  - Environment management
- **Schedules** (`services/schedules/`):
  - Cron-based task scheduling
  - Schedule execution

#### 3. Database Layer (`db/`)

Database abstraction and implementations.

**Key Components:**
- **Store Interface** (`db/Store.go`): Common interface for database operations
- **Implementations**:
  - `db/sql/`: SQL database implementation (MySQL, PostgreSQL, SQLite)
  - `db/bolt/`: BoltDB implementation (deprecated)
- **Models** (`db/*.go`):
  - User, Project, Task, Template
  - Inventory, Repository, AccessKey
  - Event, Session, Role

**Database Support:**
- MySQL 8.0+
- PostgreSQL 12+
- SQLite 3
- BoltDB (deprecated, for migration only)

#### 4. Utilities (`util/`)

Shared utilities and configuration.

**Key Components:**
- **Config** (`util/config.go`): Application configuration management
- **Encryption** (`util/encryption.go`): Encryption utilities
- **Mailer** (`util/mailer/`): Email sending
- **Error Logging** (`util/errorLogging.go`): Error logging helpers

### Frontend Components

#### 1. Vue.js Application (`web/`)

**Structure:**
- `web/src/`: Source code
  - `components/`: Reusable Vue components
  - `views/`: Page components
  - `router/`: Vue Router configuration
  - `store/`: State management (Vuex)
  - `api/`: API client
  - `utils/`: Utility functions
- `web/public/`: Static assets

**Key Technologies:**
- Vue.js 2.6.14
- Vuetify 2.6.10 (Material Design components)
- Vue Router 3.5.4
- Axios 1.12.0 (HTTP client)
- Vuex (state management)

## Data Flow

### Task Execution Flow

```
User creates/triggers task
    ↓
API receives request
    ↓
Task created in database
    ↓
Task added to Task Pool queue
    ↓
Task Pool assigns task to runner
    ↓
Runner executes task (Ansible/Terraform/etc.)
    ↓
Task output streamed via WebSocket
    ↓
Task status updated in database
    ↓
User receives notifications
```

### Authentication Flow

```
User submits credentials
    ↓
API validates credentials (local/LDAP/OIDC)
    ↓
Session created
    ↓
Secure cookie set
    ↓
User authenticated for subsequent requests
    ↓
Middleware validates session on each request
```

### Project Access Flow

```
User requests project resource
    ↓
ProjectMiddleware loads project
    ↓
Checks user membership in project
    ↓
Loads user role and permissions
    ↓
PermissionMiddleware validates access
    ↓
Request proceeds if authorized
```

## Key Design Decisions

### 1. Database Abstraction

**Decision:** Use interface-based database abstraction (`db.Store`)

**Rationale:**
- Support multiple database backends
- Easy testing with mock implementations
- Consistent API across databases

### 2. Service Layer Pattern

**Decision:** Separate business logic from API handlers

**Rationale:**
- Reusability across different interfaces
- Easier testing
- Clear separation of concerns

### 3. Task Pool Architecture

**Decision:** Centralized task pool for task management

**Rationale:**
- Efficient resource management
- Better control over concurrent tasks
- Easier monitoring and debugging

### 4. WebSocket for Real-time Updates

**Decision:** Use WebSocket for task output streaming

**Rationale:**
- Real-time feedback to users
- Efficient for long-running tasks
- Better user experience

### 5. Middleware Chain

**Decision:** Use middleware for cross-cutting concerns

**Rationale:**
- DRY principle
- Consistent behavior across endpoints
- Easy to add new functionality

## Security Architecture

### Authentication

- **Sessions**: Secure cookie-based sessions
- **LDAP**: Integration with LDAP/Active Directory
- **OIDC**: OpenID Connect support
- **TOTP**: Two-factor authentication support

### Authorization

- **Role-Based Access Control (RBAC)**: Project-level and global roles
- **Permission System**: Fine-grained permissions per resource
- **Project Isolation**: Users can only access projects they're members of

### Data Protection

- **Encryption**: Access keys encrypted at rest
- **Secure Cookies**: HTTP-only, secure cookies
- **Input Validation**: All inputs validated
- **Rate Limiting**: Protection against brute-force attacks

## Scalability Considerations

### Horizontal Scaling

- **Stateless API**: API server is stateless (except in-memory rate limiter)
- **Database**: Shared database for multiple instances
- **Redis**: Optional Redis for distributed state (HA mode)
- **Runners**: Can run on separate machines

### Performance Optimizations

- **Connection Pooling**: Database connection pooling
- **Caching**: Redis caching (planned)
- **Pagination**: List endpoints support pagination
- **Lazy Loading**: Frontend lazy loading (planned)

## Deployment Architecture

### Single Instance

```
┌─────────────────────┐
│  Semaphore Server   │
│  (API + Frontend)   │
└──────────┬──────────┘
           │
    ┌──────▼──────┐
    │  Database   │
    └─────────────┘
```

### High Availability

```
┌──────────────┐  ┌──────────────┐
│  Semaphore   │  │  Semaphore   │
│  Instance 1  │  │  Instance 2  │
└──────┬───────┘  └──────┬───────┘
       │                 │
       └────────┬────────┘
                │
    ┌───────────▼───────────┐
    │      Redis            │
    │  (Shared State)       │
    └───────────┬───────────┘
                │
    ┌───────────▼───────────┐
    │    Database           │
    │  (MySQL/Postgres)     │
    └───────────────────────┘
```

## Extension Points

### 1. Custom Apps

Support for custom application types beyond Ansible/Terraform:
- Implement `db_lib/App` interface
- Register in `db_lib/AppFactory`

### 2. Integrations

Webhook integrations for external systems:
- Implement integration handlers
- Register in integration service

### 3. Authentication Providers

Custom authentication providers:
- Implement OIDC provider configuration
- Add to authentication flow

## Future Architecture Considerations

### Planned Improvements

1. **Caching Layer**: Redis-based caching for frequently accessed data
2. **Metrics**: Prometheus metrics for monitoring
3. **Tracing**: OpenTelemetry for distributed tracing
4. **Microservices**: Potential split into separate services (if needed)
5. **Event Sourcing**: For audit and event replay

## Resources

- [Project Structure](./DEVELOPMENT.md#project-structure)
- [API Documentation](../api-docs.yml)
- [Deployment Guide](./DEPLOYMENT.md)
- [Development Guide](./DEVELOPMENT.md)

