package mcp

import (
	"context"

	mcpserver "github.com/mark3labs/mcp-go/server"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/scagogogo/npm-skills/pkg/registry"
)

// registerDownloadTools registers npm_download_stats and npm_download_range tools
func registerDownloadTools(client *registry.Registry, cfg Config) []mcpserver.ServerTool {
	var tools []mcpserver.ServerTool

	// npm_download_stats
	tools = append(tools, mcpserver.ServerTool{
		Tool: mcp.NewTool("npm_download_stats",
			mcp.WithDescription("Get download count statistics for an NPM package over a time period. Note: always queries api.npmjs.org regardless of mirror/registry settings."),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("Package name, e.g. 'react'"),
			),
			mcp.WithString("period",
				mcp.Description("Time period: last-day, last-week (default), last-month"),
				mcp.Enum("last-day", "last-week", "last-month"),
			),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			name := request.GetString("name", "")
			if name == "" {
				return toolError("package name is required"), nil
			}

			period := request.GetString("period", "last-week")

			ctx, cancel := withTimeout(ctx, cfg)
			defer cancel()

			stats, err := client.GetDownloadStats(ctx, name, period)
			if err != nil {
				return toolError("failed to get download stats for '%s': %s", name, err.Error()), nil
			}

			return toolResult(stats), nil
		},
	})

	// npm_download_range
	tools = append(tools, mcpserver.ServerTool{
		Tool: mcp.NewTool("npm_download_range",
			mcp.WithDescription("Get daily download counts for an NPM package over a time period. Returns an array of daily download counts, useful for trend visualization. Note: always queries api.npmjs.org regardless of mirror/registry settings."),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("Package name, e.g. 'react'"),
			),
			mcp.WithString("period",
				mcp.Description("Time period: last-day, last-week (default), last-month"),
				mcp.Enum("last-day", "last-week", "last-month"),
			),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			name := request.GetString("name", "")
			if name == "" {
				return toolError("package name is required"), nil
			}

			period := request.GetString("period", "last-week")

			ctx, cancel := withTimeout(ctx, cfg)
			defer cancel()

			stats, err := client.GetDownloadRangeStats(ctx, name, period)
			if err != nil {
				return toolError("failed to get download range for '%s': %s", name, err.Error()), nil
			}

			return toolResult(stats), nil
		},
	})

	return tools
}
