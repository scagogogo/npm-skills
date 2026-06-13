#!/usr/bin/env bash
#
# install.sh - Build and install the npm-skills CLI tool
#
# This script is called when the Skill is installed via Claude Code.
# It compiles the Go CLI binary and makes it available on PATH.
#
# Usage:
#   ./scripts/install.sh [--dest DIR] [--proxy URL] [--mirror NAME]
#
# Options:
#   --dest DIR     Installation directory (default: ~/.local/bin)
#   --proxy URL    Set default HTTP proxy (e.g. http://127.0.0.1:7890)
#   --mirror NAME  Set default mirror (e.g. npm-mirror, taobao)
#   --force        Force rebuild even if binary exists
#

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

DEST_DIR="$HOME/.local/bin"
FORCE=false

# Parse arguments
while [[ $# -gt 0 ]]; do
  case "$1" in
    --dest)
      DEST_DIR="$2"
      shift 2
      ;;
    --force)
      FORCE=true
      shift
      ;;
    --help|-h)
      echo "Usage: $0 [--dest DIR] [--force]"
      echo ""
      echo "Build and install the npm-skills CLI tool."
      echo ""
      echo "Options:"
      echo "  --dest DIR   Installation directory (default: ~/.local/bin)"
      echo "  --force      Force rebuild even if binary exists"
      echo ""
      echo "Environment Variables (used at runtime, not install time):"
      echo "  NPM_REGISTRY   Custom registry URL (overrides --mirror)"
      echo "  NPM_MIRROR     Default mirror name (official|taobao|npm-mirror|...)"
      echo "  NPM_PROXY      HTTP proxy URL (e.g. http://127.0.0.1:7890)"
      exit 0
      ;;
    *)
      echo "Unknown option: $1" >&2
      exit 1
      ;;
  esac
done

CLI_BINARY="npm-skills"
MCP_BINARY="npm-mcp-server"
CLI_TARGET="$DEST_DIR/$CLI_BINARY"
MCP_TARGET="$DEST_DIR/$MCP_BINARY"

# Check if already installed
if [[ -f "$CLI_TARGET" ]] && [[ "$FORCE" != "true" ]]; then
  echo "✅ $CLI_BINARY is already installed at $CLI_TARGET"
  echo "   Use --force to rebuild, or run '$CLI_BINARY help' to verify."
  exit 0
fi

# Check Go installation
if ! command -v go &> /dev/null; then
  echo "❌ Go is not installed. Please install Go 1.24+ from https://go.dev/dl/"
  echo "   Alternatively, you can use the SDK directly: go get github.com/scagogogo/npm-skills"
  exit 1
fi

GO_VERSION=$(go version | awk '{print $3}')
echo "🔧 Building $CLI_BINARY with $GO_VERSION..."

# Create destination directory
mkdir -p "$DEST_DIR"

# Build the CLI
cd "$PROJECT_ROOT"
go build -o "$CLI_TARGET" ./cmd/npm-skills/

if [[ $? -eq 0 ]]; then
  chmod +x "$CLI_TARGET"
  echo "✅ $CLI_BINARY installed successfully to $CLI_TARGET"
  echo ""

  # Build the MCP server
  echo "🔧 Building $MCP_BINARY..."
  go build -o "$MCP_TARGET" ./cmd/mcp-server/

  if [[ $? -eq 0 ]]; then
    chmod +x "$MCP_TARGET"
    echo "✅ $MCP_BINARY installed successfully to $MCP_TARGET"
    echo ""
  else
    echo "⚠️  $MCP_BINARY build failed (MCP server is optional, CLI still works)"
    echo ""
  fi

  # Check if on PATH
  if echo ":$PATH:" | grep -q ":$DEST_DIR:"; then
    echo "✅ $DEST_DIR is on your PATH"
  else
    echo "⚠️  $DEST_DIR is not on your PATH. Add it with:"
    echo "   echo 'export PATH=\"$DEST_DIR:\$PATH\"' >> ~/.bashrc"
    echo "   source ~/.bashrc"
  fi

  echo ""
  echo "📋 Quick test:"
  echo "   $CLI_TARGET help"
  echo "   $CLI_TARGET mirrors"
  echo "   $CLI_TARGET package react"
  echo ""

  # Show MCP server integration tips
  echo "🔌 MCP Server integration (for AI agents):"
  echo "   Add to your Claude Code settings:"
  echo '   {'
  echo '     "mcpServers": {'
  echo '       "npm-registry": {'
  echo '         "command": "'$MCP_BINARY'"'
  echo '       }'
  echo '     }'
  echo '   }'
  echo ""

  # Show proxy/mirror tips for China users
  echo "🌐 Tips for China users (network acceleration):"
  echo "   # Use China mirror:"
  echo "   $CLI_TARGET package react -m npm-mirror"
  echo ""
  echo "   # Or set default mirror via environment variable:"
  echo "   export NPM_MIRROR=npm-mirror"
  echo "   $CLI_TARGET package react"
  echo ""
  echo "   # Use HTTP proxy (e.g. Clash, V2Ray):"
  echo "   $CLI_TARGET package react --proxy http://127.0.0.1:7890"
  echo ""
  echo "   # Or set default proxy via environment variable:"
  echo "   export NPM_PROXY=http://127.0.0.1:7890"
  echo "   $CLI_TARGET package react"

else
  echo "❌ Build failed. Please check the Go environment."
  exit 1
fi
