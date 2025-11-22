# Dependency Update Guide

This document provides guidance on updating dependencies for the Semaphore UI project.

## Overview

The project uses:
- **Go modules** for backend dependencies (see `go.mod`)
- **npm** for frontend dependencies (see `web/package.json`)
- **Renovate** for automated dependency updates (see `renovate.json`)

## Automated Updates with Renovate

Renovate is configured to automatically create PRs for dependency updates. The configuration includes:

- **Grouped updates**: Related dependencies are grouped together
- **Scheduled updates**: Updates are checked weekly on Mondays
- **Auto-merge**: Patch updates are automatically merged (if tests pass)
- **Security priority**: Security updates are prioritized
- **Manual review**: Major version updates require manual review

### Renovate Configuration

See `renovate.json` for full configuration details.

## Manual Updates

### Go Dependencies

#### Checking for Updates

```bash
# List all dependencies
go list -m all

# Check for available updates
go list -u -m all

# Check for security vulnerabilities (requires nancy)
go list -json -m all | nancy sleuth
```

#### Updating Dependencies

```bash
# Update all dependencies to latest versions
go get -u ./...

# Update specific package
go get -u github.com/package/name@latest

# Update to specific version
go get -u github.com/package/name@v1.2.3

# Clean up unused dependencies
go mod tidy

# Verify dependencies
go mod verify
```

#### Best Practices

1. **Update incrementally**: Update one package at a time and test thoroughly
2. **Check changelogs**: Review breaking changes before updating major versions
3. **Run tests**: Always run the full test suite after updating dependencies
4. **Check security**: Use security scanning tools to check for vulnerabilities

### npm Dependencies

#### Checking for Updates

```bash
cd web

# Check outdated packages
npm outdated

# Check for security vulnerabilities
npm audit

# Check for security vulnerabilities with fix suggestions
npm audit fix --dry-run
```

#### Updating Dependencies

```bash
cd web

# Update all packages to latest versions (within semver range)
npm update

# Update specific package
npm install package-name@latest

# Update to specific version
npm install package-name@1.2.3

# Fix security vulnerabilities automatically
npm audit fix

# Fix security vulnerabilities manually (may require code changes)
npm audit fix --force
```

#### Best Practices

1. **Update incrementally**: Update one package at a time, especially for major versions
2. **Check changelogs**: Review breaking changes before updating major versions
3. **Run tests**: Always run the full test suite after updating dependencies
4. **Check compatibility**: Ensure Vue 2 compatibility when updating Vue ecosystem packages
5. **Review security**: Check npm audit reports regularly

## Testing After Updates

After updating dependencies, always:

1. **Run the test suite**: `task test`
2. **Build the project**: `task build`
3. **Run linting**: `task lint`
4. **Test manually**: Test in development environment
5. **Check for deprecation warnings**: Review build output for warnings

## Rollback

If an update causes issues:

### Go

```bash
# Revert go.mod and go.sum
git checkout go.mod go.sum

# Or restore from previous commit
git checkout HEAD~1 -- go.mod go.sum

# Re-download dependencies
go mod download
```

### npm

```bash
cd web

# Revert package.json and package-lock.json
git checkout package.json package-lock.json

# Or restore from previous commit
git checkout HEAD~1 -- package.json package-lock.json

# Reinstall
npm install
```

## Security Updates

Security updates are automatically prioritized by Renovate. For manual security updates:

### Go

```bash
# Check for known vulnerabilities
go list -json -m all | nancy sleuth

# Update vulnerable packages
go get -u vulnerable-package@patched-version
```

### npm

```bash
cd web

# Check for vulnerabilities
npm audit

# Fix automatically (if possible)
npm audit fix

# Fix manually (review changes)
npm audit fix --force
```

## Version Pinning

Some dependencies may need to be pinned to specific versions for compatibility:

### Go

```go
// In go.mod
require (
    github.com/package/name v1.2.3 // pinned for compatibility
)
```

### npm

```json
// In package.json
{
  "dependencies": {
    "package-name": "1.2.3"  // pinned for compatibility
  }
}
```

## Breaking Changes

When updating major versions:

1. **Read the changelog**: Review breaking changes
2. **Update code**: Make necessary code changes
3. **Update tests**: Update tests for new APIs
4. **Test thoroughly**: Run full test suite
5. **Document changes**: Update documentation if needed

## Monitoring

- **Renovate PRs**: Review and merge Renovate PRs regularly
- **Security alerts**: Monitor GitHub security alerts
- **Dependency status**: Check dependency status in CI/CD pipelines

## Resources

- [Go Modules Documentation](https://go.dev/ref/mod)
- [npm Documentation](https://docs.npmjs.com/)
- [Renovate Documentation](https://docs.renovatebot.com/)
- [Go Security Advisories](https://pkg.go.dev/vuln)
- [npm Security Advisories](https://www.npmjs.com/advisories)

