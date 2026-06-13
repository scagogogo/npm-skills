package mcp

import (
	"context"
	"fmt"

	mcpserver "github.com/mark3labs/mcp-go/server"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/scagogogo/npm-skills/pkg/models"
	"github.com/scagogogo/npm-skills/pkg/registry"
)

// registerDistTagWriteTools 注册 dist-tags 写操作工具
func registerDistTagWriteTools(client *registry.Registry, cfg Config) []mcpserver.ServerTool {
	var tools []mcpserver.ServerTool

	// npm_dist_tag_get
	tools = append(tools, mcpserver.ServerTool{
		Tool: mcp.NewTool("npm_dist_tag_get",
			mcp.WithDescription("Get a specific dist-tag value for an NPM package. Returns the version string that the tag points to."),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("Package name, e.g. 'react', '@nestjs/core'"),
			),
			mcp.WithString("tag",
				mcp.Required(),
				mcp.Description("Tag name, e.g. 'latest', 'next', 'beta'"),
			),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			name := request.GetString("name", "")
			tag := request.GetString("tag", "")
			if name == "" || tag == "" {
				return toolError("package name and tag are required"), nil
			}

			ctx, cancel := withTimeout(ctx, cfg)
			defer cancel()

			version, err := client.GetDistTag(ctx, name, tag)
			if err != nil {
				return toolError("failed to get dist-tag '%s' for '%s': %s", tag, name, err.Error()), nil
			}
			return toolResult(map[string]string{
				"package": name,
				"tag":     tag,
				"version": version,
			}), nil
		},
	})

	// npm_dist_tag_set
	tools = append(tools, mcpserver.ServerTool{
		Tool: mcp.NewTool("npm_dist_tag_set",
			mcp.WithDescription("Set or update a dist-tag for an NPM package to point to a specific version. Requires authentication token configured at server level."),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("Package name"),
			),
			mcp.WithString("tag",
				mcp.Required(),
				mcp.Description("Tag name to set, e.g. 'next', 'beta'"),
			),
			mcp.WithString("version",
				mcp.Required(),
				mcp.Description("Version string the tag should point to"),
			),
			mcp.WithDestructiveHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
			mcp.WithOpenWorldHintAnnotation(true),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			name := request.GetString("name", "")
			tag := request.GetString("tag", "")
			version := request.GetString("version", "")
			if name == "" || tag == "" || version == "" {
				return toolError("package name, tag, and version are required"), nil
			}

			ctx, cancel := withTimeout(ctx, cfg)
			defer cancel()

			err := client.SetDistTag(ctx, name, tag, version)
			if err != nil {
				return toolError("failed to set dist-tag '%s' to '%s' for '%s': %s", tag, version, name, err.Error()), nil
			}
			return toolResult(map[string]string{
				"package": name,
				"tag":     tag,
				"version": version,
				"status":  "updated",
			}), nil
		},
	})

	// npm_dist_tag_delete
	tools = append(tools, mcpserver.ServerTool{
		Tool: mcp.NewTool("npm_dist_tag_delete",
			mcp.WithDescription("Delete a dist-tag from an NPM package. Requires authentication token configured at server level. WARNING: Deleting 'latest' tag may cause issues."),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("Package name"),
			),
			mcp.WithString("tag",
				mcp.Required(),
				mcp.Description("Tag name to delete"),
			),
			mcp.WithDestructiveHintAnnotation(true),
			mcp.WithOpenWorldHintAnnotation(true),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			name := request.GetString("name", "")
			tag := request.GetString("tag", "")
			if name == "" || tag == "" {
				return toolError("package name and tag are required"), nil
			}

			ctx, cancel := withTimeout(ctx, cfg)
			defer cancel()

			err := client.DeleteDistTag(ctx, name, tag)
			if err != nil {
				return toolError("failed to delete dist-tag '%s' for '%s': %s", tag, name, err.Error()), nil
			}
			return toolResult(map[string]string{
				"package": name,
				"tag":     tag,
				"status":  "deleted",
			}), nil
		},
	})

	return tools
}

// registerUserTools 注册用户与认证工具
func registerUserTools(client *registry.Registry, cfg Config) []mcpserver.ServerTool {
	var tools []mcpserver.ServerTool

	// npm_user_get
	tools = append(tools, mcpserver.ServerTool{
		Tool: mcp.NewTool("npm_user_get",
			mcp.WithDescription("Get user profile information from NPM registry. Requires authentication token."),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("Username to look up"),
			),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			name := request.GetString("name", "")
			if name == "" {
				return toolError("username is required"), nil
			}

			ctx, cancel := withTimeout(ctx, cfg)
			defer cancel()

			profile, err := client.GetUser(ctx, name)
			if err != nil {
				return toolError("failed to get user '%s': %s", name, err.Error()), nil
			}
			return toolResult(profile), nil
		},
	})

	return tools
}

// registerAccessTools 注册权限管理工具
func registerAccessTools(client *registry.Registry, cfg Config) []mcpserver.ServerTool {
	var tools []mcpserver.ServerTool

	// npm_package_access
	tools = append(tools, mcpserver.ServerTool{
		Tool: mcp.NewTool("npm_package_access",
			mcp.WithDescription("Get access/permissions settings for an NPM package. Requires authentication token."),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("Package name"),
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

			access, err := client.GetPackageAccess(ctx, name)
			if err != nil {
				return toolError("failed to get access for '%s': %s", name, err.Error()), nil
			}
			return toolResult(access), nil
		},
	})

	// npm_package_collaborators
	tools = append(tools, mcpserver.ServerTool{
		Tool: mcp.NewTool("npm_package_collaborators",
			mcp.WithDescription("List collaborators for an NPM package. Requires authentication token."),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("Package name"),
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

			collabs, err := client.ListCollaborators(ctx, name)
			if err != nil {
				return toolError("failed to list collaborators for '%s': %s", name, err.Error()), nil
			}
			return toolResult(collabs), nil
		},
	})

	return tools
}

// registerStarTools 注册收藏操作工具
func registerStarTools(client *registry.Registry, cfg Config) []mcpserver.ServerTool {
	var tools []mcpserver.ServerTool

	// npm_starred_by_user
	tools = append(tools, mcpserver.ServerTool{
		Tool: mcp.NewTool("npm_starred_by_user",
			mcp.WithDescription("Get list of packages starred by a specific NPM user."),
			mcp.WithString("username",
				mcp.Required(),
				mcp.Description("Username to look up starred packages"),
			),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			username := request.GetString("username", "")
			if username == "" {
				return toolError("username is required"), nil
			}

			ctx, cancel := withTimeout(ctx, cfg)
			defer cancel()

			packages, err := client.GetStarredByUser(ctx, username)
			if err != nil {
				return toolError("failed to get starred packages for '%s': %s", username, err.Error()), nil
			}
			return toolResult(map[string]interface{}{
				"username": username,
				"packages": packages,
				"count":    len(packages),
			}), nil
		},
	})

	// npm_starred_by_package
	tools = append(tools, mcpserver.ServerTool{
		Tool: mcp.NewTool("npm_starred_by_package",
			mcp.WithDescription("Get list of users who starred a specific NPM package."),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("Package name"),
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

			users, err := client.GetStarredByPackage(ctx, name)
			if err != nil {
				return toolError("failed to get stargazers for '%s': %s", name, err.Error()), nil
			}
			return toolResult(map[string]interface{}{
				"package": name,
				"users":   users,
				"count":   len(users),
			}), nil
		},
	})

	return tools
}

// registerTokenTools 注册 Token 管理工具
func registerTokenTools(client *registry.Registry, cfg Config) []mcpserver.ServerTool {
	var tools []mcpserver.ServerTool

	// npm_token_list
	tools = append(tools, mcpserver.ServerTool{
		Tool: mcp.NewTool("npm_token_list",
			mcp.WithDescription("List all NPM access tokens for the authenticated user. Requires authentication token."),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			ctx, cancel := withTimeout(ctx, cfg)
			defer cancel()

			tokens, err := client.ListTokens(ctx)
			if err != nil {
				return toolError("failed to list tokens: %s", err.Error()), nil
			}
			return toolResult(tokens), nil
		},
	})

	return tools
}

// registerAuditTools 注册安全审计工具
func registerAuditTools(client *registry.Registry, cfg Config) []mcpserver.ServerTool {
	var tools []mcpserver.ServerTool

	// npm_audit
	tools = append(tools, mcpserver.ServerTool{
		Tool: mcp.NewTool("npm_audit",
			mcp.WithDescription("Quick security audit of NPM dependencies. Submit a map of package names to versions and get vulnerability counts by severity."),
			mcp.WithObject("dependencies",
				mcp.Required(),
				mcp.Description("Map of package names to version strings, e.g. {\"lodash\": \"4.17.11\", \"express\": \"4.17.1\"}"),
			),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			argsMap, ok := request.Params.Arguments.(map[string]interface{})
			if !ok {
				return toolError("invalid arguments"), nil
			}
			deps, ok := argsMap["dependencies"].(map[string]interface{})
			if !ok || len(deps) == 0 {
				return toolError("dependencies map is required"), nil
			}

			// Convert to map[string]string
			depMap := make(map[string]string, len(deps))
			for k, v := range deps {
				depMap[k] = fmt.Sprintf("%v", v)
			}

			ctx, cancel := withTimeout(ctx, cfg)
			defer cancel()

			result, err := client.QuickAudit(ctx, &models.QuickAuditRequest{Dependencies: depMap})
			if err != nil {
				return toolError("audit failed: %s", err.Error()), nil
			}
			return toolResult(result), nil
		},
	})

	// npm_audit_advisory
	tools = append(tools, mcpserver.ServerTool{
		Tool: mcp.NewTool("npm_audit_advisory",
			mcp.WithDescription("Get a specific NPM security advisory by ID."),
			mcp.WithNumber("id",
				mcp.Required(),
				mcp.Description("Advisory ID"),
			),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			argsMap, ok := request.Params.Arguments.(map[string]interface{})
			if !ok {
				return toolError("invalid arguments"), nil
			}
			idFloat, ok := argsMap["id"].(float64)
			if !ok {
				return toolError("advisory ID is required and must be a number"), nil
			}
			id := int(idFloat)

			ctx, cancel := withTimeout(ctx, cfg)
			defer cancel()

			advisory, err := client.GetAdvisory(ctx, id)
			if err != nil {
				return toolError("failed to get advisory %d: %s", id, err.Error()), nil
			}
			return toolResult(advisory), nil
		},
	})

	return tools
}

// registerOrgTools 注册组织和团队工具
func registerOrgTools(client *registry.Registry, cfg Config) []mcpserver.ServerTool {
	var tools []mcpserver.ServerTool

	// npm_org_get
	tools = append(tools, mcpserver.ServerTool{
		Tool: mcp.NewTool("npm_org_get",
			mcp.WithDescription("Get NPM organization details. Requires authentication token."),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("Organization name"),
			),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			name := request.GetString("name", "")
			if name == "" {
				return toolError("organization name is required"), nil
			}

			ctx, cancel := withTimeout(ctx, cfg)
			defer cancel()

			org, err := client.GetOrg(ctx, name)
			if err != nil {
				return toolError("failed to get org '%s': %s", name, err.Error()), nil
			}
			return toolResult(org), nil
		},
	})

	// npm_org_members
	tools = append(tools, mcpserver.ServerTool{
		Tool: mcp.NewTool("npm_org_members",
			mcp.WithDescription("List members of an NPM organization. Requires authentication token."),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("Organization name"),
			),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			name := request.GetString("name", "")
			if name == "" {
				return toolError("organization name is required"), nil
			}

			ctx, cancel := withTimeout(ctx, cfg)
			defer cancel()

			members, err := client.ListOrgMembers(ctx, name)
			if err != nil {
				return toolError("failed to list members of org '%s': %s", name, err.Error()), nil
			}
			return toolResult(map[string]interface{}{
				"org":     name,
				"members": members,
				"count":   len(members),
			}), nil
		},
	})

	// npm_org_packages
	tools = append(tools, mcpserver.ServerTool{
		Tool: mcp.NewTool("npm_org_packages",
			mcp.WithDescription("List packages owned by an NPM organization. Requires authentication token."),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("Organization name"),
			),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			name := request.GetString("name", "")
			if name == "" {
				return toolError("organization name is required"), nil
			}

			ctx, cancel := withTimeout(ctx, cfg)
			defer cancel()

			packages, err := client.ListOrgPackages(ctx, name)
			if err != nil {
				return toolError("failed to list packages of org '%s': %s", name, err.Error()), nil
			}
			return toolResult(map[string]interface{}{
				"org":      name,
				"packages": packages,
				"count":    len(packages),
			}), nil
		},
	})

	// npm_team_list
	tools = append(tools, mcpserver.ServerTool{
		Tool: mcp.NewTool("npm_team_list",
			mcp.WithDescription("List teams in an NPM organization. Requires authentication token."),
			mcp.WithString("org",
				mcp.Required(),
				mcp.Description("Organization name"),
			),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			org := request.GetString("org", "")
			if org == "" {
				return toolError("organization name is required"), nil
			}

			ctx, cancel := withTimeout(ctx, cfg)
			defer cancel()

			teams, err := client.ListTeams(ctx, org)
			if err != nil {
				return toolError("failed to list teams of org '%s': %s", org, err.Error()), nil
			}
			return toolResult(map[string]interface{}{
				"org":   org,
				"teams": teams,
				"count": len(teams),
			}), nil
		},
	})

	// npm_team_members
	tools = append(tools, mcpserver.ServerTool{
		Tool: mcp.NewTool("npm_team_members",
			mcp.WithDescription("List members of a team in an NPM organization. Requires authentication token."),
			mcp.WithString("org",
				mcp.Required(),
				mcp.Description("Organization name"),
			),
			mcp.WithString("team",
				mcp.Required(),
				mcp.Description("Team name"),
			),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			org := request.GetString("org", "")
			team := request.GetString("team", "")
			if org == "" || team == "" {
				return toolError("organization name and team name are required"), nil
			}

			ctx, cancel := withTimeout(ctx, cfg)
			defer cancel()

			members, err := client.ListTeamMembers(ctx, org, team)
			if err != nil {
				return toolError("failed to list members of team '%s/%s': %s", org, team, err.Error()), nil
			}
			return toolResult(map[string]interface{}{
				"org":     org,
				"team":    team,
				"members": members,
				"count":   len(members),
			}), nil
		},
	})

	return tools
}

// registerHooksTools 注册 Webhook 工具
func registerHooksTools(client *registry.Registry, cfg Config) []mcpserver.ServerTool {
	var tools []mcpserver.ServerTool

	// npm_hook_list
	tools = append(tools, mcpserver.ServerTool{
		Tool: mcp.NewTool("npm_hook_list",
			mcp.WithDescription("List all NPM webhooks owned by the authenticated user. Requires authentication token."),
			mcp.WithString("package",
				mcp.Description("Optional: filter hooks by package name"),
			),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			pkg := request.GetString("package", "")

			ctx, cancel := withTimeout(ctx, cfg)
			defer cancel()

			opts := models.HookListOptions{Package: pkg}
			hooks, err := client.ListHooks(ctx, opts)
			if err != nil {
				return toolError("failed to list hooks: %s", err.Error()), nil
			}
			return toolResult(map[string]interface{}{
				"hooks": hooks,
				"count": len(hooks),
			}), nil
		},
	})

	// npm_hook_get
	tools = append(tools, mcpserver.ServerTool{
		Tool: mcp.NewTool("npm_hook_get",
			mcp.WithDescription("Get details of a specific NPM webhook by ID. Requires authentication token."),
			mcp.WithString("id",
				mcp.Required(),
				mcp.Description("Webhook ID"),
			),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(true),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			id := request.GetString("id", "")
			if id == "" {
				return toolError("hook ID is required"), nil
			}

			ctx, cancel := withTimeout(ctx, cfg)
			defer cancel()

			hook, err := client.GetHook(ctx, id)
			if err != nil {
				return toolError("failed to get hook '%s': %s", id, err.Error()), nil
			}
			return toolResult(hook), nil
		},
	})

	return tools
}

// registerCouchDBTools 注册 CouchDB 高级查询工具
func registerCouchDBTools(client *registry.Registry, cfg Config) []mcpserver.ServerTool {
	var tools []mcpserver.ServerTool

	// npm_changes
	tools = append(tools, mcpserver.ServerTool{
		Tool: mcp.NewTool("npm_changes",
			mcp.WithDescription("Get NPM registry changes feed. Returns recent package changes, useful for mirroring or incremental data sync."),
			mcp.WithString("since",
				mcp.Description("Start from this update sequence (from last_seq of previous call)"),
			),
			mcp.WithNumber("limit",
				mcp.Description("Maximum number of changes to return"),
			),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithIdempotentHintAnnotation(false),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			since := request.GetString("since", "")
			limit := 25 // default
			argsMap, argsOk := request.Params.Arguments.(map[string]interface{})
			if argsOk {
				if limitFloat, ok := argsMap["limit"].(float64); ok {
					limit = int(limitFloat)
				}
			}

			ctx, cancel := withTimeout(ctx, cfg)
			defer cancel()

			opts := models.ChangesOptions{
				Since: since,
				Limit: limit,
			}
			changes, err := client.GetChanges(ctx, opts)
			if err != nil {
				return toolError("failed to get changes: %s", err.Error()), nil
			}
			return toolResult(changes), nil
		},
	})

	return tools
}
