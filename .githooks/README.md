# Git Hooks

This directory contains git hooks for the Semaphore UI project.

## Installation

To install the pre-commit hook, run:

```bash
# On Unix/Linux/macOS
cp .githooks/pre-commit .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit

# On Windows (PowerShell)
Copy-Item .githooks/pre-commit .git/hooks/pre-commit
```

Or use the install script:

```bash
# On Unix/Linux/macOS
chmod +x .githooks/install.sh
./.githooks/install.sh
```

## Available Hooks

### pre-commit

Runs before each commit:
- Runs `golangci-lint` to check code quality
- Checks code formatting with `gofmt`
- Validates that code follows project standards

## Requirements

- `golangci-lint` - Install with: `go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.57.2`
- `gofmt` - Usually comes with Go installation

## Skipping Hooks

If you need to skip hooks for a specific commit (not recommended), use:

```bash
git commit --no-verify
```

