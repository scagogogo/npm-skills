package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	mcpserver "github.com/mark3labs/mcp-go/server"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/scagogogo/npm-skills/pkg/registry"
)

// Config holds the configuration for the NPM MCP server
type Config struct {
	// RegistryOptions is the SDK options used to create the registry client
	RegistryOptions *registry.Options
	// Timeout is the default request timeout for all operations
	Timeout time.Duration
}

// NewServer creates a new MCP server with all NPM registry tools registered.
//
// The server exposes NPM registry operations as MCP tools that AI agents
// can call through the Model Context Protocol.
//
// Parameters:
//   - cfg: Server configuration including registry options and timeout
//
// Returns:
//   - *mcpserver.MCPServer: Configured MCP server with all tools registered
func NewServer(cfg Config) *mcpserver.MCPServer {
	s := mcpserver.NewMCPServer(
		"npm-registry",
		"0.1.0",
		mcpserver.WithToolCapabilities(true),
		mcpserver.WithInstructions(
			"NPM Registry MCP Server — query package info, search packages, get download stats, list versions, check dist-tags, "+
				"audit security, manage access, view orgs/teams, and more. "+
				"Use npm_package_summary instead of npm_package for most queries (much smaller response). "+
				"Download stats always query api.npmjs.org regardless of mirror/registry settings. "+
				"Write operations (dist-tag set/delete, star, etc.) require authentication token configured at server level.",
		),
	)

	// Create the registry client
	client := registry.NewRegistry(cfg.RegistryOptions)

	// Registry info tools
	s.AddTools(registerRegistryTools(client, cfg)...)

	// Package tools
	s.AddTools(registerPackageTools(client, cfg)...)

	// Search tools
	s.AddTools(registerSearchTools(client, cfg)...)

	// Version tools
	s.AddTools(registerVersionTools(client, cfg)...)

	// Download stats tools
	s.AddTools(registerDownloadTools(client, cfg)...)

	// WhoAmI tool
	s.AddTools(registerWhoamiTools(client, cfg)...)

	// Dist-tag write tools (requires token)
	s.AddTools(registerDistTagWriteTools(client, cfg)...)

	// User & auth tools
	s.AddTools(registerUserTools(client, cfg)...)

	// Access & collaborator tools
	s.AddTools(registerAccessTools(client, cfg)...)

	// Star tools
	s.AddTools(registerStarTools(client, cfg)...)

	// Token management tools
	s.AddTools(registerTokenTools(client, cfg)...)

	// Audit tools
	s.AddTools(registerAuditTools(client, cfg)...)

	// Org & team tools
	s.AddTools(registerOrgTools(client, cfg)...)

	// Hook tools
	s.AddTools(registerHooksTools(client, cfg)...)

	// CouchDB advanced tools
	s.AddTools(registerCouchDBTools(client, cfg)...)

	return s
}

// withTimeout returns a context with the configured timeout, deriving from the parent context
// so that client-side cancellation is properly propagated
func withTimeout(parent context.Context, cfg Config) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, cfg.Timeout)
}

// toolError creates an error result for a tool call
func toolError(format string, args ...any) *mcp.CallToolResult {
	return mcp.NewToolResultError(fmt.Sprintf(format, args...))
}

// toolResult creates a successful result with JSON-formatted content
func toolResult(data any) *mcp.CallToolResult {
	text := formatJSON(data)
	return mcp.NewToolResultText(text)
}

// formatJSON serializes data to indented JSON with truncation for large responses
func formatJSON(data any) string {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error": "failed to serialize response: %s"}`, err.Error())
	}

	result := string(bytes)

	// Hard limit: truncate at 100KB
	const maxResponseSize = 100 * 1024
	if len(result) > maxResponseSize {
		truncated := result[:maxResponseSize]
		return truncated + "\n... [RESPONSE TRUNCATED at 100KB]"
	}

	return result
}
