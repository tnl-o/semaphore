#!/bin/bash
#
# Script to check for outdated dependencies
#

set -e

echo "=== Checking Go Dependencies ==="
echo ""

# Check if go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed"
    exit 1
fi

echo "Current Go version:"
go version
echo ""

echo "Checking for outdated Go modules..."
echo "Note: This may take a while..."
echo ""

# Get list of direct dependencies
echo "Direct dependencies:"
go list -m -f '{{if not .Indirect}}{{.}}{{end}}' all | grep -v "^github.com/semaphoreui/semaphore" | head -20
echo ""

# Check for available updates (this is informational)
echo "To update dependencies, run:"
echo "  go get -u ./..."
echo "  go mod tidy"
echo ""

echo "=== Checking npm Dependencies ==="
echo ""

# Check if npm is installed
if ! command -v npm &> /dev/null; then
    echo "Warning: npm is not installed, skipping npm checks"
    exit 0
fi

if [ ! -f "web/package.json" ]; then
    echo "Warning: web/package.json not found"
    exit 0
fi

cd web

echo "Current npm version:"
npm --version
echo ""

echo "Checking for outdated npm packages..."
echo "Note: This may take a while..."
echo ""

# Check outdated packages
npm outdated || true

echo ""
echo "To update dependencies, run:"
echo "  cd web"
echo "  npm update"
echo "  npm audit fix"
echo ""

cd ..

