package mcp

import (
	"context"

	mcpserver "github.com/mark3labs/mcp-go/server"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/scagogogo/npm-skills/pkg/registry"
)

// registerRegistryTools registers npm_registry_info and npm_mirrors tools
func registerRegistryTools(client *registry.Registry, cfg Config) []mcpserver.ServerTool {
	var tools []mcpserver.ServerTool

	// npm_registry_info
	tools = append(tools, mcpserver.ServerTool{
		Tool: mcp.NewTool("npm_registry_info",
			mcp.WithDescription("Get NPM registry status and statistics including total package count, disk size, and update sequence."),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			ctx, cancel := withTimeout(ctx, cfg)
			defer cancel()

			info, err := client.GetRegistryInformation(ctx)
			if err != nil {
				return toolError("failed to get registry info: %s", err.Error()), nil
			}
			return toolResult(info), nil
		},
	})

	// npm_mirrors
	tools = append(tools, mcpserver.ServerTool{
		Tool: mcp.NewTool("npm_mirrors",
			mcp.WithDescription("List all available NPM mirror sources with their URLs, regions, and descriptions. Useful for users in China who need faster access."),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			mirrors := registry.ListMirrors()
			return toolResult(mirrors), nil
		},
	})

	return tools
}
