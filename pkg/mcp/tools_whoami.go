package mcp

import (
	"context"

	mcpserver "github.com/mark3labs/mcp-go/server"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/scagogogo/npm-skills/pkg/registry"
)

// registerWhoamiTools registers the npm_whoami tool
func registerWhoamiTools(client *registry.Registry, cfg Config) []mcpserver.ServerTool {
	var tools []mcpserver.ServerTool

	// npm_whoami
	tools = append(tools, mcpserver.ServerTool{
		Tool: mcp.NewTool("npm_whoami",
			mcp.WithDescription("Check current NPM authentication status. Returns the authenticated username. Requires the server to be configured with an auth token (--token flag or NPM_TOKEN env var)."),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			ctx, cancel := withTimeout(ctx, cfg)
			defer cancel()

			username, err := client.WhoAmI(ctx)
			if err != nil {
				return toolError("authentication check failed: %s", err.Error()), nil
			}

			return toolResult(map[string]string{
				"username": username,
				"status":   "authenticated",
			}), nil
		},
	})

	return tools
}
