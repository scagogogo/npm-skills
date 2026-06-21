#!/usr/bin/env bash
# Install script for npm-skills CLI — builds and installs to ~/.local/bin/
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"

echo "Building npm-skills CLI from $PROJECT_ROOT ..."

cd "$PROJECT_ROOT"

# Build CLI
go build -o "${HOME}/.local/bin/npm-skills" ./cmd/npm-skills
echo "✓ npm-skills → ~/.local/bin/npm-skills"

# Build MCP server
go build -o "${HOME}/.local/bin/npm-mcp-server" ./cmd/mcp-server
echo "✓ npm-mcp-server → ~/.local/bin/npm-mcp-server"

echo "Done! Both binaries installed to ~/.local/bin/"
