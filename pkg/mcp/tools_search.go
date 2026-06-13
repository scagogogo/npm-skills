package mcp

import (
	"context"
	"encoding/json"
	"strconv"

	mcpserver "github.com/mark3labs/mcp-go/server"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/scagogogo/npm-skills/pkg/registry"
)

// registerSearchTools registers the npm_search tool
func registerSearchTools(client *registry.Registry, cfg Config) []mcpserver.ServerTool {
	var tools []mcpserver.ServerTool

	// npm_search
	tools = append(tools, mcpserver.ServerTool{
		Tool: mcp.NewTool("npm_search",
			mcp.WithDescription("Search NPM packages by keyword with optional pagination and score weighting. Returns matching packages with relevance scores."),
			mcp.WithString("query",
				mcp.Required(),
				mcp.Description("Search keywords (quote multi-word queries, e.g. 'http client')"),
			),
			mcp.WithNumber("limit",
				mcp.Description("Maximum number of results (default 20)"),
			),
			mcp.WithNumber("from",
				mcp.Description("Pagination offset, 0-based (for browsing beyond first page)"),
			),
			mcp.WithNumber("quality",
				mcp.Description("Quality weight, 0.0-1.0 (emphasize code quality in results)"),
			),
			mcp.WithNumber("popularity",
				mcp.Description("Popularity weight, 0.0-1.0 (emphasize popular packages)"),
			),
			mcp.WithNumber("maintenance",
				mcp.Description("Maintenance weight, 0.0-1.0 (emphasize well-maintained packages)"),
			),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			query := request.GetString("query", "")
			if query == "" {
				return toolError("search query is required"), nil
			}

			// Parse optional parameters
			opts := registry.SearchOptions{Size: 20}

			if limitStr := request.GetString("limit", ""); limitStr != "" {
				if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
					opts.Size = limit
				}
			}
			if fromStr := request.GetString("from", ""); fromStr != "" {
				if from, err := strconv.Atoi(fromStr); err == nil && from > 0 {
					opts.From = from
				}
			}
			if qualityStr := request.GetString("quality", ""); qualityStr != "" {
				if quality, err := strconv.ParseFloat(qualityStr, 64); err == nil {
					opts.Quality = quality
				}
			}
			if popularityStr := request.GetString("popularity", ""); popularityStr != "" {
				if popularity, err := strconv.ParseFloat(popularityStr, 64); err == nil {
					opts.Popularity = popularity
				}
			}
			if maintenanceStr := request.GetString("maintenance", ""); maintenanceStr != "" {
				if maintenance, err := strconv.ParseFloat(maintenanceStr, 64); err == nil {
					opts.Maintenance = maintenance
				}
			}

			ctx, cancel := withTimeout(ctx, cfg)
			defer cancel()

			result, err := client.SearchPackagesWithOptions(ctx, query, opts)
			if err != nil {
				return toolError("failed to search for '%s': %s", query, err.Error()), nil
			}

			return toolResult(result), nil
		},
	})

	return tools
}

// getOptionalFloat extracts an optional float64 parameter from the tool request
func getOptionalFloat(request mcp.CallToolRequest, key string) (float64, bool) {
	args := request.GetArguments()
	val, ok := args[key]
	if !ok {
		return 0, false
	}
	switch v := val.(type) {
	case float64:
		return v, true
	case json.Number:
		f, err := v.Float64()
		if err != nil {
			return 0, false
		}
		return f, true
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, false
		}
		return f, true
	}
	return 0, false
}
