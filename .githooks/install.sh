#!/bin/bash
#
# Install git hooks for Semaphore UI
#

set -e

HOOKS_DIR=".githooks"
GIT_HOOKS_DIR=".git/hooks"

if [ ! -d "$GIT_HOOKS_DIR" ]; then
    echo "Error: .git/hooks directory not found. Are you in the project root?"
    exit 1
fi

if [ ! -d "$HOOKS_DIR" ]; then
    echo "Error: .githooks directory not found. Are you in the project root?"
    exit 1
fi

echo "Installing git hooks..."

# Install pre-commit hook
if [ -f "$HOOKS_DIR/pre-commit" ]; then
    cp "$HOOKS_DIR/pre-commit" "$GIT_HOOKS_DIR/pre-commit"
    chmod +x "$GIT_HOOKS_DIR/pre-commit"
    echo "✓ Installed pre-commit hook"
else
    echo "✗ pre-commit hook not found"
fi

echo "Git hooks installed successfully!"

