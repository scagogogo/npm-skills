---
name: npm
description: NPM Registry client — query and manage packages, versions, dist-tags, downloads, access, stars, tokens, webhooks, orgs, and audit. Supports publish, unpublish, deprecate, custom registry, 8 mirrors, proxy, and auth token. Trigger phrases: "npm package", "npm publish", "npm registry", "search npm", "npm stats", "npm mirror", "npm 版本", "npm 包", "npm 镜像", "npm 发布", "npm 审计".
version: 1.0.0
---

# NPM Skills — NPM Registry Client

Complete NPM Registry client — query packages, manage dist-tags, publish/unpublish, control access, audit security, and more. All CLI commands output JSON to stdout for easy parsing by AI agents.

## Installation

### Option 1: Claude Code Plugin (Recommended for AI agents)

```bash
claude plugin marketplace add scagogogo/npm-skills
claude plugin install npm@npm-skills
```

### Option 2: Download from GitHub Release

Pre-built binaries are available for all major platforms from the [GitHub Releases page](https://github.com/scagogogo/npm-skills/releases).

```bash
# Linux (x86_64)
curl -sL https://github.com/scagogogo/npm-skills/releases/latest/download/npm-skills_0.2.0_linux_x86_64.tar.gz | tar -xz
sudo mv npm-skills npm-mcp-server /usr/local/bin/

# macOS (Apple Silicon)
curl -sL https://github.com/scagogogo/npm-skills/releases/latest/download/npm-skills_0.2.0_aarch64.tar.gz | tar -xz
sudo mv npm-skills npm-mcp-server /usr/local/bin/

# macOS (Intel)
curl -sL https://github.com/scagogogo/npm-skills/releases/latest/download/npm-skills_0.2.0_x86_64.tar.gz | tar -xz
sudo mv npm-skills npm-mcp-server /usr/local/bin/

# Windows (PowerShell)
Invoke-WebRequest -Uri "https://github.com/scagogogo/npm-skills/releases/latest/download/npm-skills_0.2.0_windows_x86_64.zip" -OutFile "npm-skills.zip"
Expand-Archive npm-skills.zip
```

**Supported platforms:** Linux (amd64, arm64, 386, arm, mips*, ppc64*, riscv64, s390x, loong64), macOS (amd64, arm64), Windows (amd64, 386), FreeBSD, OpenBSD, NetBSD, Illumos, Solaris.

### Option 3: Build from Source

On first use or for development, build the CLI tool:

```bash
cd {{SKILL_DIR}} && bash scripts/install.sh
```

This compiles `npm-skills` and `npm-mcp-server` to `~/.local/bin/`.

### Option 4: Go Install

```bash
go install github.com/scagogogo/npm-skills/cmd/npm-skills@latest
go install github.com/scagogogo/npm-skills/cmd/mcp-server@latest
```

### Verify Installation

```bash
npm-skills --version
npm-skills mirrors
```

## Quick Reference (Most Common Commands)

> These cover 90% of typical NPM queries. Scroll down for advanced operations.

```bash
npm-skills package-summary <name>            # Lightweight package info (recommended)
npm-skills search <query> -l 10              # Search packages
npm-skills versions <name>                   # List all versions
npm-skills dist-tags get <name>              # Get dist-tags (latest, next, beta)
npm-skills download-stats <name> -p last-month  # Download stats
npm-skills mirrors                           # List available mirrors
```

## When to Use This Skill

Use this skill when the user asks about:

| Category | What | Key Commands |
|----------|------|-------------|
| **Package Info** | Description, versions, dependencies, maintainers | `package`, `package-summary`, `pkg-version` |
| **Search** | Find packages by keyword | `search` |
| **Downloads** | Popularity, trends, bulk comparison | `download-stats`, `download-range`, `download-stats-bulk`, `download-stats-date` |
| **Versions** | List versions, get latest, specific version | `versions`, `pkg-version` |
| **Dist-Tags** | Get/set/delete version aliases (latest, next, beta) | `dist-tags get/set/delete` |
| **Auth** | Check login status | `whoami`, `user login` |
| **Publishing** | Publish or unpublish packages | `publish`, `unpublish`, `deprecate` |
| **Access** | Package permissions, collaborators | `access get/set/collaborators/grant/revoke` |
| **Stars** | Star/unstar packages | `star add/remove/list/stargazers` |
| **Tokens** | Manage API tokens | `token list/get/create/delete` |
| **Audit** | Security audit, advisories | `audit quick/bulk/advisory/advisories` |
| **Orgs/Teams** | Organization and team management | `org get/create/delete/members/teams/...` |
| **Webhooks** | Create and manage package webhooks | `hook list/get/create/update/delete` |
| **Registry** | Registry health, mirror sources | `registry-info`, `mirrors`, `config` |
| **Tarball** | Download .tgz files | `download` |

Trigger phrases: `npm package`, `npm publish`, `NPM registry`, `search npm`, `npm stats`, `npm mirror`, `npm 版本`, `npm 包`, `npm 镜像`, `npm 发布`, `npm 审计`

## Global Flags

All commands support these flags:

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--mirror` | `-m` | `official` | Mirror source name (official\|taobao\|npm-mirror\|huawei\|tencent\|cnpm\|yarn\|npmjscom) |
| `--registry` | | | Custom registry URL (overrides --mirror, for private registries) |
| `--token` | `-t` | | NPM auth token (env: NPM_TOKEN). Required for write operations. |
| `--proxy` | | | HTTP proxy URL |
| `--timeout` | | `120` | Request timeout in seconds |
| `--no-color` | | `false` | Disable colored output |

**Environment Variables** (used as defaults when flags not set):

| Variable | Equivalent Flag |
|----------|----------------|
| `NPM_MIRROR` | `--mirror` |
| `NPM_REGISTRY` | `--registry` |
| `NPM_TOKEN` | `--token` |
| `NPM_PROXY` | `--proxy` |

**Priority**: CLI flag > Environment variable > Default

---

## CLI Commands — Read Operations

### Get Package Information

```bash
npm-skills package-summary <name>           # Lightweight (recommended, KB response)
npm-skills package <name>                    # Full metadata (can be 10MB+)
npm-skills pkg <name> -m taobao             # Using mirror
npm-skills ps <name> --proxy http://127.0.0.1:7890
```

> **Tip**: Always prefer `package-summary` over `package` — much smaller and faster.

### Search Packages

```bash
npm-skills search <query>                   # Basic search
npm-skills s <query> -l 10                  # Limit results
npm-skills search <query> --from 20 -l 10   # Paginated
npm-skills search <query> --popularity 1.0   # Weight by popularity
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--limit` | `-l` | `20` | Max results |
| `--from` | | `0` | Pagination offset |
| `--quality` | | `0` | Quality weight (0.0-1.0) |
| `--popularity` | | `0` | Popularity weight (0.0-1.0) |
| `--maintenance` | | `0` | Maintenance weight (0.0-1.0) |

### Version Information

```bash
npm-skills versions <name>                  # All versions
npm-skills vs <name> --latest               # Latest version only
npm-skills pkg-version <name> <version>      # Specific version detail
```

### Dist-Tags (Read)

```bash
npm-skills dist-tags get <name>             # All dist-tags
npm-skills dt get <name> --abbreviated      # Lightweight endpoint
npm-skills tags get <name> -m npm-mirror
```

### Download Statistics

```bash
# Single package
npm-skills download-stats <name> -p last-month

# Daily trends
npm-skills download-range <name> -p last-week

# Custom date range
npm-skills download-stats-date <name> --start 2024-01-01 --end 2024-06-30

# Bulk comparison (up to 128 packages)
npm-skills download-stats-bulk react,vue,angular -p last-month
```

> **Note**: Download stats always query api.npmjs.org regardless of mirror/registry.

### Other Read Commands

```bash
npm-skills registry-info                    # Registry health/info
npm-skills mirrors                          # List mirror sources
npm-skills config                           # Show current configuration
npm-skills whoami --token npm_xxxxx         # Check auth status
npm-skills download <name> <ver> <dest>     # Download .tgz tarball
```

---

## CLI Commands — Write Operations (require --token)

All write operations require authentication. Use `--token` or set `NPM_TOKEN`.

### Publish / Unpublish / Deprecate

```bash
# Publish
npm-skills publish ./my-pkg-1.0.0.tgz --name my-pkg --version 1.0.0 -t npm_xxxxx

# Deprecate a version (safe, doesn't delete)
npm-skills deprecate my-pkg 1.0.0 -M "Use v2.0.0 instead" -t npm_xxxxx

# Unpublish a specific version (dangerous)
npm-skills unpublish my-pkg --version 1.0.0 -t npm_xxxxx

# Unpublish entire package (very dangerous)
npm-skills unpublish my-pkg --force -t npm_xxxxx
```

### Dist-Tags Management

```bash
npm-skills dist-tags set <name> <tag> --version <ver> -t npm_xxxxx   # Set a tag
npm-skills dist-tags delete <name> <tag> -t npm_xxxxx                 # Delete a tag
```

### Package Access & Collaborators

```bash
npm-skills access get <name> -t npm_xxxxx                             # Get access settings
npm-skills access set <name> --visibility public -t npm_xxxxx         # Set visibility
npm-skills access collaborators <name> -t npm_xxxxx                   # List collaborators
npm-skills access grant <name> <user> --permission read -t npm_xxxxx  # Grant access
npm-skills access revoke <name> <user> -t npm_xxxxx                   # Revoke access
```

### Stars

```bash
npm-skills star add <name> -t npm_xxxxx           # Star a package
npm-skills star remove <name> -t npm_xxxxx         # Unstar a package
npm-skills star list <username>                     # Packages starred by user
npm-skills star stargazers <name>                   # Users who starred a package
```

### Token Management

```bash
npm-skills token list -t npm_xxxxx                 # List tokens
npm-skills token get <id> -t npm_xxxxx              # Get token details
npm-skills token create --password mypass -t npm_xxxxx  # Create token
npm-skills token delete <id> -t npm_xxxxx            # Delete token
```

### Security Audit

```bash
npm-skills audit quick --deps "lodash=4.17.11,express=4.17.1"   # Quick audit
npm-skills audit bulk --advisories "lodash=<4.17.12"            # Bulk audit
npm-skills audit advisory 123                                    # Get advisory by ID
npm-skills audit advisories --package lodash                     # List advisories
```

### Organization & Team Management

```bash
# Organization
npm-skills org get <org> -t npm_xxxxx
npm-skills org create <org> -t npm_xxxxx
npm-skills org delete <org> -t npm_xxxxx
npm-skills org members <org> -t npm_xxxxx
npm-skills org member-add <org> <user> -t npm_xxxxx
npm-skills org member-remove <org> <user> -t npm_xxxxx
npm-skills org packages <org> -t npm_xxxxx

# Teams
npm-skills org team-list <org> -t npm_xxxxx
npm-skills org team-create <org> <team> -t npm_xxxxx
npm-skills org team-delete <org> <team> -t npm_xxxxx
npm-skills org team-members <org> <team> -t npm_xxxxx
npm-skills org team-member-add <org> <team> <user> -t npm_xxxxx
npm-skills org team-member-remove <org> <team> <user> -t npm_xxxxx
npm-skills org team-packages <org> <team> -t npm_xxxxx
```

### Webhooks

```bash
npm-skills hook list -t npm_xxxxx                               # List hooks
npm-skills hook get <id> -t npm_xxxxx                            # Get hook
npm-skills hook create --name my-hook --endpoint https://... -t npm_xxxxx  # Create
npm-skills hook update <id> --endpoint https://new... -t npm_xxxxx          # Update
npm-skills hook delete <id> -t npm_xxxxx                         # Delete
```

### User Operations

```bash
npm-skills user login -u myuser -p mypass                        # Login
npm-skills user signup -u myuser -p mypass --email me@x.com      # Register
npm-skills user get <username> -t npm_xxxxx                      # Get profile
```

---

## Mirror Sources

| Mirror | Name | Region | Best For |
|--------|------|--------|----------|
| https://registry.npmjs.org | `official` | Global | Default, most up-to-date |
| https://registry.npmmirror.com | `npm-mirror` | China | Fast access in China (recommended) |
| https://registry.npm.taobao.org | `taobao` | China | Legacy Taobao mirror |
| https://mirrors.huaweicloud.com/repository/npm | `huawei` | China | Huawei Cloud users |
| http://mirrors.cloud.tencent.com/npm | `tencent` | China | Tencent Cloud users |
| http://r.cnpmjs.org | `cnpm` | China | CNPM community mirror |
| https://registry.yarnpkg.com | `yarn` | Global | Yarn users |
| https://skimdb.npmjs.com | `npmjscom` | Global | CouchDB metadata mirror |

> Pass any URL directly: `--mirror https://your-registry.com` or `--registry https://your-registry.com`.

## Proxy & Network

```bash
# HTTP proxy (one-time)
npm-skills package react --proxy http://127.0.0.1:7890

# HTTP proxy (default)
export NPM_PROXY=http://127.0.0.1:7890

# China mirror (no proxy needed)
npm-skills package react -m npm-mirror
export NPM_MIRROR=npm-mirror

# Private registry
npm-skills package my-lib --registry https://npm.my-company.com
export NPM_REGISTRY=https://npm.my-company.com
```

## Common Workflows

### Check if a package exists and get its latest version
```bash
npm-skills package-summary axios
npm-skills versions react --latest
```

### Compare package popularity
```bash
npm-skills download-stats-bulk react,vue,angular -p last-month
```

### Find alternatives
```bash
npm-skills search "http client" -l 10 --popularity 1.0
```

### Publish a new version to a private registry
```bash
npm-skills publish ./my-pkg-1.0.0.tgz --name my-pkg --version 1.0.0 \
  --registry https://npm.my-company.com -t npm_xxxxx
```

### Set a dist-tag on a private registry
```bash
npm-skills dist-tags set my-pkg stable --version 1.0.0 \
  --registry https://npm.my-company.com -t npm_xxxxx
```

### Security audit dependencies
```bash
npm-skills audit quick --deps "lodash=4.17.11,express=4.17.1"
```

### Access NPM from China
```bash
export NPM_MIRROR=npm-mirror
npm-skills package antd
npm-skills search "vue component" -l 5
```

---

## SDK Usage (Go)

For programmatic access within Go code:

```go
import "github.com/scagogogo/npm-skills/pkg/registry"

// Default client (official registry)
client := registry.NewRegistry()

// China mirror
client = registry.NewNpmMirrorRegistry()

// Custom registry with auth, timeout, and user-agent
options := registry.NewOptions().
    SetRegistryURL("https://npm.my-company.com").
    SetToken("npm_xxxxx").
    SetProxy("http://proxy:8080").
    SetTimeout(30 * time.Second).
    SetUserAgent("my-app/1.0")
client = registry.NewRegistry(options)
```

### Core Read Methods

```go
pkg, _ := client.GetPackageInformation(ctx, "react")                    // Full metadata
pkg, _ := client.GetAbbreviatedPackageInformation(ctx, "react")         // Lightweight
results, _ := client.SearchPackages(ctx, "http client", 10)             // Basic search
results, _ := client.SearchPackagesWithOptions(ctx, "http client", ...) // Advanced search
ver, _ := client.GetPackageVersion(ctx, "react", "18.2.0")             // Specific version
versions, _ := client.GetPackageVersions(ctx, "react")                  // All versions
latest, _ := client.GetPackageLatestVersion(ctx, "react")               // Latest version
tags, _ := client.GetDistTags(ctx, "react")                             // Dist-tags
tag, _ := client.GetDistTag(ctx, "react", "latest")                     // Single tag
stats, _ := client.GetDownloadStats(ctx, "react", "last-week")          // Download stats
	rangeStats, _ := client.GetDownloadRangeStats(ctx, "react", "last-week")  // Daily trends
	rangeStats, _ := client.GetDownloadRangeStatsByDateRange(ctx, "react", "2024-01-01", "2024-06-30") // Custom date range
	bulkStats, _ := client.GetBulkDownloadStats(ctx, []string{"react", "vue"}, "last-week")  // Bulk (>128 auto-chunked)
info, _ := client.GetRegistryInformation(ctx)                           // Registry info
user, _ := client.WhoAmI(ctx)                                           // Auth check
```

### Write Methods (require token)

```go
options.SetToken("npm_xxxxx")

// Dist-tags
client.SetDistTag(ctx, "my-pkg", "next", "2.0.0-rc.1")
client.DeleteDistTag(ctx, "my-pkg", "beta")
client.SetDistTags(ctx, "my-pkg", map[string]string{"latest": "2.0.0"})

// Publish / Unpublish / Deprecate
client.PublishPackage(ctx, pkg)
client.PublishPackageFromTarball(ctx, "my-pkg", "1.0.0", tarballBytes, metadata)
client.DeprecateVersion(ctx, "my-pkg", "1.0.0", "Use v2.0.0")
client.UnpublishPackage(ctx, "my-pkg")
client.UnpublishPackageVersion(ctx, "my-pkg", "1.0.0")

// Access & Collaborators
access, _ := client.GetPackageAccess(ctx, "my-pkg")
client.SetPackageAccess(ctx, "my-pkg", &models.PackageAccessUpdate{Access: "public"})
collabs, _ := client.ListCollaborators(ctx, "my-pkg")
client.GrantAccess(ctx, "my-pkg", "username", models.PermissionWrite)
client.RevokeAccess(ctx, "my-pkg", "username")

// Stars
client.StarPackage(ctx, "react")
client.UnstarPackage(ctx, "react")
starred, _ := client.GetStarredByUser(ctx, "username")
stargazers, _ := client.GetStarredByPackage(ctx, "react")

// Tokens
tokens, _ := client.ListTokens(ctx)
token, _ := client.CreateToken(ctx, &models.TokenCreation{Password: "mypass"})
client.DeleteToken(ctx, "token-id")

// Audit
result, _ := client.QuickAudit(ctx, &models.QuickAuditRequest{Dependencies: deps})
result, _ := client.BulkAudit(ctx, advisories)
advisory, _ := client.GetAdvisory(ctx, 123)
advisories, _ := client.ListAdvisories(ctx, nil)

// Org & Team
org, _ := client.GetOrg(ctx, "my-org")
client.CreateOrg(ctx, "my-org")
client.DeleteOrg(ctx, "my-org")
members, _ := client.ListOrgMembers(ctx, "my-org")
client.AddOrgMember(ctx, "my-org", "username")
teams, _ := client.ListTeams(ctx, "my-org")
client.CreateTeam(ctx, "my-org", "core")
teamMembers, _ := client.ListTeamMembers(ctx, "my-org", "core")
client.AddTeamMember(ctx, "my-org", "core", "username")

// Hooks
hooks, _ := client.ListHooks(ctx, "my-pkg")
hook, _ := client.GetHook(ctx, "hook-id")
client.CreateHook(ctx, &models.HookCreation{Name: "my-hook", Endpoint: "https://..."})
client.UpdateHook(ctx, "hook-id", &models.HookUpdate{Endpoint: "https://new..."})
client.DeleteHook(ctx, "hook-id")

// User
loginResult, _ := client.Login(ctx, "myuser", "mypass")
client.CreateUser(ctx, &models.UserCreation{Name: "newuser", Password: "pass", Email: "me@x.com"})
profile, _ := client.GetUser(ctx, "username")
```

---

## Important Notes

1. **JSON output**: All CLI commands output JSON to stdout; status/info messages go to stderr. Use `jq` for filtering.
2. **Auth required**: Write operations need `--token` or `NPM_TOKEN`. Use `whoami` to verify.
3. **package-summary**: Prefer over `package` — much smaller and faster response.
4. **Private registry**: Use `--registry <url> --token <token>` for Verdaccio, Artifactory, GitHub Packages, etc.
5. **Download stats**: Always query api.npmjs.org regardless of mirror/registry settings. Use `Options.SetDownloadStatsURL()` to override in SDK.
6. **Scoped packages**: Scoped packages like `@nestjs/core` are fully supported in all operations.
7. **Mirror as URL**: Pass any URL as `--mirror https://my-registry.com` for direct use.
8. **Help**: `npm-skills help` or `npm-skills <command> --help`.
9. **Error handling**: SDK provides typed errors (`registry.ErrNotFound`, `registry.ErrUnauthorized`, `registry.ErrRateLimited`, etc.) for programmatic error checking with `errors.Is()`.
10. **Timeout**: SDK supports per-client timeout via `Options.SetTimeout()`. CLI uses `--timeout` flag (default 120s).
11. **Bulk auto-chunking**: `GetBulkDownloadStats` and `GetBulkDownloadRangeStats` automatically chunk requests >128 packages.

## Detailed API Reference

For complete SDK documentation including all types and methods, see [references/api.md](references/api.md).

## MCP Server

An MCP (Model Context Protocol) server exposes NPM registry operations as tools for AI agents — no shell invocation needed.

### Build

```bash
bash scripts/install.sh    # Builds both CLI and MCP server
```

### Claude Code Integration

```json
{
  "mcpServers": {
    "npm-registry": {
      "command": "npm-mcp-server",
      "args": ["--mirror", "npm-mirror"]
    }
  }
}
```

### MCP Tools (33 total)

**Read Tools:**

| Tool | Description |
|------|-------------|
| `npm_registry_info` | Registry status and statistics |
| `npm_mirrors` | List available mirror sources |
| `npm_package` | Full package metadata (large) |
| `npm_package_summary` | Lightweight package metadata (recommended) |
| `npm_search` | Search packages with pagination and weighting |
| `npm_version` | Specific version metadata |
| `npm_versions` | All published version numbers |
| `npm_latest_version` | Latest version number |
| `npm_dist_tags` | Distribution tags (latest, next, beta) |
| `npm_download_stats` | Download count for a period |
| `npm_download_range` | Daily download trend data |
| `npm_whoami` | Check auth status |

**Write Tools (destructive, require token):**

| Tool | Description |
|------|-------------|
| `npm_dist_tag_set` | Set a dist-tag to a version |
| `npm_dist_tag_delete` | Delete a dist-tag |
| `npm_dist_tags_set` | Batch set multiple dist-tags |
| `npm_star` | Star a package |
| `npm_unstar` | Unstar a package |
| `npm_stargazers` | Get users who starred a package |
| `npm_starred_by_user` | Get packages starred by user |
| `npm_access_get` | Get package access settings |
| `npm_collaborators` | List package collaborators |
| `npm_token_list` | List API tokens |
| `npm_audit_quick` | Quick security audit |
| `npm_audit_advisory` | Get security advisory by ID |
| `npm_hook_list` | List webhooks |
| `npm_hook_get` | Get webhook details |
| `npm_org_get` | Get organization info |
| `npm_org_members` | List org members |
| `npm_org_packages` | List org packages |
| `npm_team_list` | List teams in an org |
| `npm_team_members` | List team members |
| `npm_team_packages` | List team packages |
| `npm_changes` | Get registry changes feed |
