package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	mcpserver "github.com/mark3labs/mcp-go/server"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/scagogogo/npm-skills/pkg/models"
	"github.com/scagogogo/npm-skills/pkg/registry"
)

// registerPackageTools registers npm_package and npm_package_summary tools
func registerPackageTools(client *registry.Registry, cfg Config) []mcpserver.ServerTool {
	var tools []mcpserver.ServerTool

	// npm_package
	tools = append(tools, mcpserver.ServerTool{
		Tool: mcp.NewTool("npm_package",
			mcp.WithDescription("Get complete metadata for an NPM package including all versions, README, maintainers, and dependencies. WARNING: Response can be very large (10MB+). Use npm_package_summary for most queries."),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("Package name, e.g. 'react', '@nestjs/core'"),
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

			pkg, err := client.GetPackageInformation(ctx, name)
			if err != nil {
				return toolError("failed to get package '%s': %s", name, err.Error()), nil
			}

			return toolResult(truncatePackage(pkg)), nil
		},
	})

	// npm_package_summary
	tools = append(tools, mcpserver.ServerTool{
		Tool: mcp.NewTool("npm_package_summary",
			mcp.WithDescription("Get lightweight metadata for an NPM package (name, description, dist-tags, version list). Much smaller response than npm_package. Recommended for most queries."),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("Package name, e.g. 'react', '@nestjs/core'"),
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

			pkg, err := client.GetAbbreviatedPackageInformation(ctx, name)
			if err != nil {
				return toolError("failed to get package summary for '%s': %s", name, err.Error()), nil
			}

			// Abbreviated already omits README, but still may have large Versions map
			return toolResult(truncatePackage(pkg)), nil
		},
	})

	return tools
}

const (
	// truncateReadmeThreshold is the max README length before truncation
	truncateReadmeThreshold = 2000
	// truncateResponseThreshold is the max JSON response size before truncation
	truncateResponseThreshold = 50 * 1024 // 50KB
)

// truncatePackage truncates large Package responses to keep them manageable
// for MCP tool results. Truncates README and replaces Versions map with version keys.
func truncatePackage(pkg *models.Package) any {
	// First pass: serialize to check size
	bytes, err := json.Marshal(pkg)
	if err != nil {
		return pkg // fallback to original if serialization fails
	}

	// If under threshold, return as-is
	if len(bytes) <= truncateResponseThreshold {
		return pkg
	}

	// Need to truncate — create a map representation for manipulation
	result := make(map[string]any)
	raw, err := json.Marshal(pkg)
	if err != nil {
		return pkg
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return pkg
	}

	// Truncate README
	if readme, ok := result["readme"].(string); ok && len(readme) > truncateReadmeThreshold {
		omitted := len(readme) - truncateReadmeThreshold
		result["readme"] = readme[:truncateReadmeThreshold] + fmt.Sprintf("\n... [truncated, %d chars omitted]", omitted)
	}

	// Replace Versions map with just version keys
	if versions, ok := result["versions"].(map[string]any); ok {
		keys := make([]string, 0, len(versions))
		for k := range versions {
			keys = append(keys, k)
		}
		delete(result, "versions")
		result["version_keys"] = keys
	}

	// Add truncation metadata
	result["_truncated"] = true
	result["_original_size_bytes"] = len(bytes)

	return result
}
