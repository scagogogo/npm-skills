<div align="center">

# NPM Skills

[Switch to 中文版](README_zh.md)

<img src="https://cdn.worldvectorlogo.com/logos/npm-2.svg" width="180" alt="NPM Logo" style="filter: brightness(0.9);">

[![Go Tests](https://github.com/scagogogo/npm-skills/actions/workflows/go-test.yml/badge.svg)](https://github.com/scagogogo/npm-skills/actions/workflows/go-test.yml)
[![Release](https://github.com/scagogogo/npm-skills/actions/workflows/release.yml/badge.svg)](https://github.com/scagogogo/npm-skills/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/scagogogo/npm-skills.svg)](https://pkg.go.dev/github.com/scagogogo/npm-skills)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**NPM Registry client for AI agents and developers** — query packages, manage dist-tags, publish, audit, and more.

[4 Ways to Integrate](#-four-ways-to-integrate) · [Download](https://github.com/scagogogo/npm-skills/releases/latest) · [Documentation](https://pkg.go.dev/github.com/scagogogo/npm-skills)

</div>

---

## ⚡ One-Click Plugin Install

Install as a Claude Code plugin — AI agents will automatically discover and use it:

```bash
# Step 1: Add the marketplace
claude plugin marketplace add scagogogo/npm-skills

# Step 2: Install the plugin
claude plugin install npm@npm-skills
```

After installation, just ask Claude Code naturally:
- *"Find info about the axios NPM package"*
- *"Download the react tarball"*
- *"Search for HTTP client libraries on NPM"*
- *"Get download stats for vue last month"*
- *"Check NPM registry using the China mirror"*
- *"Publish my package to a private registry"*
- *"Audit my dependencies for vulnerabilities"*

---

## 🔌 Four Ways to Integrate

NPM Skills is designed **AI-native first**, offering four complementary ways to interact with the NPM Registry:

### 1. 🤖 Skill (for AI Agents) — Primary

This repository is a **Claude Code Plugin** — install it and AI agents will automatically discover and use it. No shell invocation needed.

**Install:**
```bash
claude plugin marketplace add scagogogo/npm-skills
claude plugin install npm@npm-skills
```

After installation, just ask Claude Code naturally:
- *"Find info about the axios NPM package"*
- *"Download the react tarball"*
- *"Search for HTTP client libraries on NPM"*
- *"Get download stats for vue last month"*
- *"Check NPM registry using the China mirror"*
- *"Publish my package to a private registry"*
- *"Audit my dependencies for vulnerabilities"*

**Trigger phrases**: `npm package`, `npm publish`, `npm registry`, `search npm`, `npm stats`, `npm mirror`, `npm 版本`, `npm 包`, `npm 镜像`, `npm 发布`

The Skill manifest (`SKILL.md`) uses progressive disclosure:
- **Immediate context**: name + description in frontmatter (~100 words)
- **Core guidance**: CLI commands + usage patterns
- **Deep reference**: Full API docs in `references/api.md` (loaded on demand)

### 2. 📦 Go SDK (for Developers)

Drop-in Go library for programmatic access with full type safety:

```go
import "github.com/scagogogo/npm-skills/pkg/registry"

// Default client (official registry)
client := registry.NewRegistry()

// Custom client with options
options := registry.NewOptions().
    SetRegistryURL("https://registry.npmjs.org").
    SetToken("npm_xxxxx").
    SetProxy("http://proxy:8080").
    SetTimeout(30 * time.Second)
client = registry.NewRegistry(options)

// Read operations
pkg, _ := client.GetPackageInformation(ctx, "react")
versions, _ := client.GetPackageVersions(ctx, "react")
stats, _ := client.GetDownloadStats(ctx, "react", "last-week")
rangeStats, _ := client.GetDownloadRangeStatsByDateRange(ctx, "react", "2024-01-01", "2024-06-30")

// Write operations (require token)
client.SetDistTag(ctx, "my-pkg", "next", "2.0.0-rc.1")
client.PublishPackage(ctx, pkg)
client.DeprecateVersion(ctx, "my-pkg", "1.0.0", "Use v2.0.0")

// Typed errors for programmatic handling
import "errors"
_, err := client.GetPackageInformation(ctx, "nonexistent")
if errors.Is(err, registry.ErrNotFound) {
    // handle 404
}
```

### 3. 🖥️ CLI Tool

Command-line interface with colorful output, proxy & mirror support. Pre-built binaries available for [all major platforms](https://github.com/scagogogo/npm-skills/releases/latest).

**Install:**
```bash
# Download from GitHub Release (recommended)
# See: https://github.com/scagogogo/npm-skills/releases/latest

# Or build from source
bash scripts/install.sh

# Or go install
go install github.com/scagogogo/npm-skills/cmd/npm-skills@latest
```

**Usage:**
```bash
# Read operations
npm-skills package-summary react            # Lightweight package info (recommended)
npm-skills package react                    # Full package metadata
npm-skills search "http client" -l 10       # Search packages
npm-skills versions react --latest          # Get latest version
npm-skills dist-tags get react              # Get dist-tags
npm-skills download-stats axios -p last-month  # Download stats
npm-skills download lodash 4.17.21 ./lodash.tgz  # Download tarball
npm-skills mirrors                          # List mirror sources
npm-skills whoami --token npm_xxxxx         # Check auth status

# Write operations (require --token)
npm-skills publish ./pkg.tgz --name my-pkg --version 1.0.0 -t npm_xxxxx
npm-skills deprecate my-pkg 1.0.0 -M "Use v2" -t npm_xxxxx
npm-skills dist-tags set my-pkg stable --version 1.0.0 -t npm_xxxxx
npm-skills access get my-pkg -t npm_xxxxx
npm-skills star add react -t npm_xxxxx
npm-skills audit quick --deps "lodash=4.17.11"

# Mirror & Proxy & Private Registry
npm-skills package react -m npm-mirror                                    # China mirror
npm-skills package react --proxy http://127.0.0.1:7890                    # HTTP proxy
npm-skills package my-lib --registry https://npm.my-company.com -t npm_x  # Private registry

# Environment variables
export NPM_MIRROR=npm-mirror
export NPM_PROXY=http://127.0.0.1:7890
export NPM_REGISTRY=https://npm.company.com
npm-skills package react    # Auto-uses env vars

npm-skills --help           # Show all 26 commands
```

### 4. 📡 MCP Server (for AI Tool Chains)

An MCP (Model Context Protocol) server that exposes NPM registry operations as tools for any MCP-compatible AI client — Claude Code, Cursor, Windsurf, and more.

**Install:**
```bash
bash scripts/install.sh   # Builds both CLI and MCP server
```

**Configuration (Claude Code / Cursor / any MCP client):**
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

**33 MCP Tools** available, including:

| Read Tools | Write Tools |
|---|---|
| `npm_registry_info`, `npm_mirrors`, `npm_package`, `npm_package_summary`, `npm_search`, `npm_version`, `npm_versions`, `npm_latest_version`, `npm_dist_tags`, `npm_download_stats`, `npm_download_range`, `npm_whoami` | `npm_dist_tag_set`, `npm_dist_tag_delete`, `npm_dist_tags_set`, `npm_star`, `npm_unstar`, `npm_stargazers`, `npm_access_get`, `npm_collaborators`, `npm_token_list`, `npm_audit_quick`, `npm_audit_advisory`, `npm_hook_list`, `npm_hook_get`, `npm_org_get`, `npm_org_members`, `npm_org_packages`, `npm_team_list`, `npm_team_members`, `npm_team_packages`, `npm_changes` |

---

## ✨ Features

- 🤖 **AI-Native First**: Designed as a Skill with progressive disclosure for AI agents
- 🚀 **High Performance**: Go-based with concurrent requests and streaming downloads
- 🌐 **8 Mirror Sources**: Built-in support for official, China, and global mirrors
- 🔄 **Proxy Support**: HTTP proxy configuration for restricted networks
- 📦 **Full API Coverage**: 70+ SDK methods covering all major NPM Registry endpoints
- 🛡️ **Typed Errors**: `ErrNotFound`, `ErrUnauthorized`, `ErrRateLimited`, etc. with `errors.Is()` support
- ⏱️ **Timeout Control**: Per-client timeout via `Options.SetTimeout()`
- 🔒 **Auth Support**: Bearer token for publish, unpublish, and all write operations
- 📊 **Download Analytics**: Point stats, range stats, bulk stats with auto-chunking (>128 packages)
- 🔍 **Package Search**: Pagination, quality/popularity/maintenance scoring
- 📡 **MCP Protocol**: 33 tools for AI tool chains
- 🏗️ **Cross-Platform**: Pre-built binaries for Linux, macOS, Windows, FreeBSD, OpenBSD, NetBSD, Illumos, Solaris

## 📥 Installation

### Download Binary (Recommended)

Pre-built binaries are available from the [Latest Release](https://github.com/scagogogo/npm-skills/releases/latest):

```bash
# Linux (x86_64)
curl -sL https://github.com/scagogogo/npm-skills/releases/latest/download/npm-skills_0.2.0_linux_x86_64.tar.gz | tar -xz
sudo mv npm-skills npm-mcp-server /usr/local/bin/

# macOS (Apple Silicon)
curl -sL https://github.com/scagogogo/npm-skills/releases/latest/download/npm-skills_0.2.0_aarch64.tar.gz | tar -xz
sudo mv npm-skills npm-mcp-server /usr/local/bin/

# Windows — download the .zip from the releases page
```

### Go Install

```bash
go install github.com/scagogogo/npm-skills/cmd/npm-skills@latest
go install github.com/scagogogo/npm-skills/cmd/mcp-server@latest
```

### Go Module

```bash
go get github.com/scagogogo/npm-skills
```

## 🚀 Quick Start

### Go SDK

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

func main() {
    // Create client (official registry by default)
    client := registry.NewRegistry()
    ctx := context.Background()

    // Get lightweight package info
    pkg, err := client.GetAbbreviatedPackageInformation(ctx, "react")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Package: %s, Latest: %s\n", pkg.Name, pkg.DistTags["latest"])

    // Search packages
    results, err := client.SearchPackages(ctx, "http client", 5)
    if err != nil {
        log.Fatal(err)
    }
    for _, obj := range results.Objects {
        fmt.Printf("  %s — %s\n", obj.Package.Name, obj.Package.Description)
    }

    // Download stats
    stats, err := client.GetDownloadStats(ctx, "react", "last-week")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("React downloads (last week): %d\n", stats.Downloads)

    // Custom registry with auth & timeout
    options := registry.NewOptions().
        SetRegistryURL("https://npm.my-company.com").
        SetToken("npm_xxxxx").
        SetTimeout(30 * time.Second)
    privateClient := registry.NewRegistry(options)
    _ = privateClient
}
```

## 🪞 Supported Mirror Sources

| Mirror | URL | Region | SDK Method |
|--------|-----|--------|------------|
| NPM Official | `https://registry.npmjs.org` | Global | `NewRegistry()` |
| NPM Mirror | `https://registry.npmmirror.com` | China | `NewNpmMirrorRegistry()` |
| Taobao | `https://registry.npm.taobao.org` | China | `NewTaoBaoRegistry()` |
| Huawei Cloud | `https://mirrors.huaweicloud.com/repository/npm` | China | `NewHuaWeiCloudRegistry()` |
| Tencent Cloud | `http://mirrors.cloud.tencent.com/npm` | China | `NewTencentRegistry()` |
| CNPM | `http://r.cnpmjs.org` | China | `NewCnpmRegistry()` |
| Yarn | `https://registry.yarnpkg.com` | Global | `NewYarnRegistry()` |
| NPM CouchDB | `https://skimdb.npmjs.com` | Global | `NewNpmjsComRegistry()` |

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. Fork this repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Create a Pull Request

## 📄 License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgements

- [NPM Registry](https://registry.npmjs.org) — API and data
- [Go Requests](https://github.com/crawler-go-go-go/go-requests) — HTTP client library
- [Cobra](https://github.com/spf13/cobra) — CLI framework
- [MCP-Go](https://github.com/mark3labs/mcp-go) — MCP server framework
