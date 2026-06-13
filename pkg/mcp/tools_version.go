package mcp

import (
	"context"

	mcpserver "github.com/mark3labs/mcp-go/server"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/scagogogo/npm-skills/pkg/registry"
)

// registerVersionTools registers npm_version, npm_versions, npm_latest_version, and npm_dist_tags tools
func registerVersionTools(client *registry.Registry, cfg Config) []mcpserver.ServerTool {
	var tools []mcpserver.ServerTool

	// npm_version
	tools = append(tools, mcpserver.ServerTool{
		Tool: mcp.NewTool("npm_version",
			mcp.WithDescription("Get metadata for a specific version of an NPM package, including dependencies, scripts, and distribution info."),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("Package name, e.g. 'react'"),
			),
			mcp.WithString("version",
				mcp.Required(),
				mcp.Description("Version string or dist-tag like 'latest'"),
			),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			name := request.GetString("name", "")
			version := request.GetString("version", "")
			if name == "" || version == "" {
				return toolError("package name and version are required"), nil
			}

			ctx, cancel := withTimeout(ctx, cfg)
			defer cancel()

			result, err := client.GetPackageVersion(ctx, name, version)
			if err != nil {
				return toolError("failed to get version '%s' of '%s': %s", version, name, err.Error()), nil
			}

			return toolResult(result), nil
		},
	})

	// npm_versions
	tools = append(tools, mcpserver.ServerTool{
		Tool: mcp.NewTool("npm_versions",
			mcp.WithDescription("List all published version numbers of an NPM package, sorted in ascending order."),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("Package name, e.g. 'react'"),
			),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			name := request.GetString("name", "")
			if name == "" {
				return toolError("package name is required"), nil
			}

			ctx, cancel := withTimeout(ctx, cfg)
			defer cancel()

			versions, err := client.GetPackageVersions(ctx, name)
			if err != nil {
				return toolError("failed to get versions for '%s': %s", name, err.Error()), nil
			}

			return toolResult(map[string]any{
				"package":       name,
				"version_count": len(versions),
				"versions":      versions,
			}), nil
		},
	})

	// npm_latest_version
	tools = append(tools, mcpserver.ServerTool{
		Tool: mcp.NewTool("npm_latest_version",
			mcp.WithDescription("Get the latest version number of an NPM package. Lightweight and fast — only queries the dist-tags endpoint."),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("Package name, e.g. 'react'"),
			),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			name := request.GetString("name", "")
			if name == "" {
				return toolError("package name is required"), nil
			}

			ctx, cancel := withTimeout(ctx, cfg)
			defer cancel()

			latest, err := client.GetPackageLatestVersion(ctx, name)
			if err != nil {
				return toolError("failed to get latest version for '%s': %s", name, err.Error()), nil
			}

			return toolResult(map[string]string{
				"package": name,
				"latest":  latest,
			}), nil
		},
	})

	// npm_dist_tags
	tools = append(tools, mcpserver.ServerTool{
		Tool: mcp.NewTool("npm_dist_tags",
			mcp.WithDescription("Get distribution tags (dist-tags) for an NPM package. Returns version aliases like 'latest', 'next', 'beta', etc."),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("Package name, e.g. 'react'"),
			),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			name := request.GetString("name", "")
			if name == "" {
				return toolError("package name is required"), nil
			}

			ctx, cancel := withTimeout(ctx, cfg)
			defer cancel()

			tags, err := client.GetDistTagsAbbreviated(ctx, name)
			if err != nil {
				return toolError("failed to get dist-tags for '%s': %s", name, err.Error()), nil
			}

			return toolResult(tags), nil
		},
	})

	return tools
}
