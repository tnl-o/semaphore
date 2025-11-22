# Dependency Update Guide

This document provides guidance on updating dependencies for the Semaphore UI project.

## Go Dependencies

### Checking for Updates

```bash
# List all dependencies
go list -m all

# Check for available updates
go list -u -m all

# Check for security vulnerabilities
go list -json -m all | nancy sleuth
```

### Updating Dependencies

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

### Best Practices

1. **Update incrementally**: Update one package at a time and test thoroughly
2. **Check changelogs**: Review breaking changes before updating major versions
3. **Run tests**: Always run the full test suite after updating dependencies
4. **Check security**: Use `go list -json -m all | nancy sleuth` to check for vulnerabilities

## npm Dependencies

### Checking for Updates

```bash
cd web

# Check outdated packages
npm outdated

# Check for security vulnerabilities
npm audit

# Check for security vulnerabilities with fix suggestions
npm audit fix --dry-run
```

### Updating Dependencies

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

### Best Practices

1. **Update incrementally**: Update one package at a time, especially for major versions
2. **Check changelogs**: Review breaking changes before updating major versions
3. **Run tests**: Always run the full test suite after updating dependencies
4. **Check compatibility**: Ensure Vue 2 compatibility when updating Vue ecosystem packages
5. **Review security**: Check npm audit reports regularly

## Automated Updates with Renovate

The project uses Renovate for automated dependency updates. Configuration is in `renovate.json`.

### Renovate Features

- Automatic PR creation for dependency updates
- Grouping related updates
- Security vulnerability fixes
- Custom update schedules

### Manual Renovate Trigger

Renovate runs automatically on a schedule, but you can also trigger it manually in the GitHub repository settings.

## Testing After Updates

After updating dependencies, always:

1. Run the test suite: `task test`
2. Build the project: `task build`
3. Run linting: `task lint`
4. Test manually in development environment
5. Check for deprecation warnings

## Rollback

If an update causes issues:

### Go

```bash
# Revert go.mod and go.sum
git checkout go.mod go.sum

# Or restore from previous commit
git checkout HEAD~1 -- go.mod go.sum
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

