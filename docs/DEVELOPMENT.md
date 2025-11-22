# Development Guide

This guide helps developers set up a development environment and contribute to Semaphore UI.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Development Workflow](#development-workflow)
- [Testing](#testing)
- [Code Standards](#code-standards)
- [Common Tasks](#common-tasks)
- [Troubleshooting](#troubleshooting)

## Prerequisites

### Required Software

- **Go**: 1.21+ ([Installation Guide](https://golang.org/doc/install))
- **Node.js**: 16+ ([Installation Guide](https://nodejs.org/))
- **Task Runner**: [Installation Guide](https://taskfile.dev/installation/)
- **Git**: For version control

### Optional Software

- **MySQL/PostgreSQL**: For database testing
- **Docker**: For containerized development
- **golangci-lint**: For code linting

## Development Setup

### 1. Clone Repository

```bash
git clone --recursive https://github.com/semaphoreui/semaphore.git
cd semaphore
```

### 2. Install Dependencies

```bash
# Install Task runner
go install github.com/go-task/task/v3/cmd/task@latest

# Install all dependencies (Go and npm)
task deps
```

**Note**: First run takes ~3 minutes. Never cancel during installation.

### 3. Build Application

```bash
task build
```

**Note**: Build takes ~1.5 minutes. Never cancel during build.

### 4. Setup Database

```bash
# Interactive setup (choose BoltDB option 2 for development)
./bin/semaphore setup
```

Or use environment variables:

```bash
export SEMAPHORE_DB_DIALECT=bolt
export SEMAPHORE_ADMIN=admin
export SEMAPHORE_ADMIN_PASSWORD=changeme
export SEMAPHORE_ADMIN_NAME=Admin
export SEMAPHORE_ADMIN_EMAIL=admin@localhost
./bin/semaphore setup
```

### 5. Run Application

```bash
./bin/semaphore server --config ./config.json
```

Access at:
- Web UI: http://localhost:3000
- API: http://localhost:3000/api/ping

Default credentials: `admin` / `changeme`

## Project Structure

```
.
├── api/                    # API handlers and routes
│   ├── helpers/           # Helper functions
│   ├── middleware/        # HTTP middleware
│   ├── projects/          # Project-related endpoints
│   └── ...
├── cli/                   # CLI application entry point
│   ├── cmd/              # CLI commands
│   └── main.go           # Main entry point
├── db/                    # Database layer
│   ├── bolt/             # BoltDB implementation
│   ├── sql/              # SQL implementation
│   └── *.go              # Database models
├── services/              # Business logic services
│   ├── tasks/            # Task management
│   ├── project/          # Project services
│   └── ...
├── util/                  # Utilities and configuration
├── web/                   # Frontend Vue.js application
│   ├── src/              # Source code
│   │   ├── components/   # Vue components
│   │   ├── views/        # Page views
│   │   ├── router/       # Vue Router
│   │   └── ...
│   └── public/           # Static assets
├── docs/                  # Documentation
├── deployment/            # Deployment configurations
├── go.mod                 # Go dependencies
├── Taskfile.yml          # Task runner configuration
└── README.md             # Project README
```

### Key Directories

- **`api/`**: HTTP API handlers, routes, and middleware
- **`db/`**: Database models and implementations
- **`services/`**: Business logic and orchestration
- **`util/`**: Shared utilities and configuration
- **`web/`**: Frontend Vue.js application
- **`cli/`**: Command-line interface

## Development Workflow

### 1. Create Feature Branch

```bash
git checkout -b feature/your-feature-name
```

### 2. Make Changes

- Write code following [Code Standards](#code-standards)
- Add tests for new functionality
- Update documentation if needed

### 3. Test Changes

```bash
# Run all tests
task test

# Run specific test
go test ./api/...

# Run frontend tests
cd web && npm run test:unit
```

### 4. Lint Code

```bash
# Backend linting
golangci-lint run

# Frontend linting
cd web && npm run lint
```

### 5. Build and Validate

```bash
# Build application
task build

# Run application and test manually
./bin/semaphore server --config ./config.json
```

### 6. Commit Changes

```bash
git add .
git commit -m "feat: add new feature"
```

Follow [Conventional Commits](https://www.conventionalcommits.org/) format.

### 7. Push and Create PR

```bash
git push origin feature/your-feature-name
```

Create pull request on GitHub.

## Testing

### Backend Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package
go test ./api/...
```

### Frontend Tests

```bash
cd web
npm run test:unit
```

### Integration Tests

Using Dredd for API integration tests:

```bash
# Build Dredd hooks
task dredd:hooks

# Install Dredd
npm install -g dredd

# Create config for Dredd
# Edit .dredd/config.json

# Start Semaphore server
./bin/semaphore server

# Run Dredd tests
dredd --config ./.dredd/dredd.local.yml
```

See [CONTRIBUTING.md](../CONTRIBUTING.md) for more details.

## Code Standards

### Go Code

- Follow [Effective Go](https://go.dev/doc/effective_go) guidelines
- Use `gofmt` for formatting
- Handle errors explicitly (no `//nolint:errcheck`)
- Use meaningful variable and function names
- Add comments for exported functions and types

### Vue.js Code

- Follow [Vue.js Style Guide](https://vuejs.org/style-guide/)
- Use ESLint for linting
- Follow component naming conventions
- Use Vuex for state management

### API Development

- Use proper HTTP status codes
- Validate all input data
- Return consistent JSON responses
- Update API documentation (`api-docs.yml`)

### Error Handling

```go
// Good
if err != nil {
    log.WithError(err).Error("Failed to process request")
    helpers.WriteError(w, err)
    return
}

// Bad
if err != nil {
    log.Error(err) // Missing context
    return // No error response to client
}
```

## Common Tasks

### Adding New API Endpoint

1. Add route in `api/router.go`
2. Create handler function
3. Add middleware if needed
4. Update API documentation
5. Add tests

### Adding New Database Model

1. Create model in `db/`
2. Add database methods in `db/sql/` or `db/bolt/`
3. Add to Store interface if needed
4. Add migrations if schema changes

### Adding New Frontend Component

1. Create component in `web/src/components/`
2. Add to view if needed
3. Update router if new page
4. Add tests

### Debugging

**Backend:**
- Use `log.Debug()` for debug messages
- Set log level: `export SEMAPHORE_LOG_LEVEL=debug`
- Use debugger (Delve) for step debugging

**Frontend:**
- Use browser DevTools
- Vue DevTools extension
- Console logging

## Troubleshooting

### Build Issues

**Issue: Go module errors**

```bash
go clean -modcache
go mod download
go mod tidy
```

**Issue: Frontend build errors**

```bash
cd web
rm -rf node_modules package-lock.json
npm install
```

### Database Issues

**Issue: Cannot connect to database**

- Check database is running
- Verify credentials in config
- Check network connectivity

**Issue: Migration errors**

```bash
./bin/semaphore migrate --config ./config.json
```

### Port Conflicts

**Issue: Port 3000 already in use**

Change port in `config.json`:

```json
{
  "port": ":3001"
}
```

## Resources

- [Architecture Documentation](./ARCHITECTURE.md)
- [Deployment Guide](./DEPLOYMENT.md)
- [Contributing Guide](../CONTRIBUTING.md)
- [API Documentation](../api-docs.yml)
- [Vue.js Documentation](https://vuejs.org/)
- [Go Documentation](https://go.dev/doc/)

## Getting Help

- [GitHub Issues](https://github.com/semaphoreui/semaphore/issues)
- [Discord Community](https://discord.gg/5R6k7hNGcH)
- [Documentation](https://docs.semaphoreui.com)

