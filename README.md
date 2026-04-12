# NPM Crawler

<div align="center">

[Switch to Chinese Version](README_zh.md)

<img src="https://cdn.worldvectorlogo.com/logos/npm-2.svg" width="180" alt="NPM Logo" style="filter: brightness(0.9);">

[![Go Tests](https://github.com/scagogogo/npm-crawler/actions/workflows/go-test.yml/badge.svg)](https://github.com/scagogogo/npm-crawler/actions/workflows/go-test.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/scagogogo/npm-crawler.svg)](https://pkg.go.dev/github.com/scagogogo/npm-crawler)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

_High-performance NPM Registry client with multi-mirror source and proxy support_

</div>

## Three Ways to Use

NPM Crawler is designed to be **AI-native first**, offering three complementary ways to interact:

### 1. 🤖 AI / Agent Mode (Primary)

Designed for AI agents and autonomous workflows. AI can directly understand and invoke this tool's capabilities without manual intervention.

```markdown
# This repository IS a Claude Code Skill
# AI agents automatically discover and use it when you ask:
# - "Find info about the axios NPM package"
# - "Download the react tarball"
# - "Search for HTTP client libraries"
# - "Get download stats for vue"
```

**AI Trigger Phrases**: `npm package`, `NPM registry`, `search npm`, `download npm tarball`, `get npm stats`, `npm mirror`

The Skill manifest (`SKILL.md`) provides AI-readable instructions with progressive disclosure:
- **Immediate context**: name + description in frontmatter (~100 words)
- **Core guidance**: Quick start + capabilities in body (~500 lines)
- **Deep reference**: Full API docs in `references/api.md` (loaded on demand)

### 2. 📦 Go SDK

Drop-in Go library for programmatic access to NPM Registry:

```go
import "github.com/scagogogo/npm-crawler/pkg/registry"

client := registry.NewRegistry()
pkg, err := client.GetPackageInformation(ctx, "react")
```

### 3. 🖥️ CLI Tool

Command-line interface for quick lookups and scripting:

```bash
cd examples/create_registry && go run main.go
cd examples/download_tarball && go run main.go
```

---

## Introduction

NPM Crawler is a high-performance NPM Registry client library written in Go, providing a simple and easy-to-use API to access package information in the NPM Registry. This library supports multiple NPM mirror sources, including the official Registry, Taobao mirror, Huawei Cloud mirror, etc., and also supports proxy configuration to easily handle various network environments.

## Features

- 🤖 **AI-Native**: Designed as a Skill for AI agents with progressive disclosure
- 🚀 **High Performance**: Based on Go's high concurrency features, providing fast NPM Registry access
- 🌐 **Multi-Mirror Source Support**: Built-in support for multiple NPM mirror sources
- 🔄 **Proxy Support**: Configurable HTTP proxy to adapt to various network environments
- 📦 **Complete Types**: Complete Go type definitions corresponding to various NPM package metadata
- 🧪 **Comprehensive Testing**: Complete unit test coverage
- 📝 **Detailed Documentation**: Bilingual annotations and documentation in both Chinese and English

## Installation

```bash
go get github.com/scagogogo/npm-crawler
```

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/scagogogo/npm-crawler/pkg/registry"
)

func main() {
    // Create default Registry client (using official npmjs.org)
    client := registry.NewRegistry()
    
    // Or use Taobao mirror
    // client := registry.NewTaoBaoRegistry()
    
    ctx := context.Background()
    
    // Get package information
    pkg, err := client.GetPackageInformation(ctx, "react")
    if err != nil {
        log.Fatalf("Failed to get package information: %v", err)
    }
    
    fmt.Printf("Package Name: %s\n", pkg.Name)
    // Output: Package Name: react
    
    fmt.Printf("Description: %s\n", pkg.Description)
    // Output: Description: React is a JavaScript library for building user interfaces.
    
    fmt.Printf("Latest Version: %s\n", pkg.DistTags["latest"])
    // Output: Latest Version: 18.2.0
    
    // Get Registry information
    info, err := client.GetRegistryInformation(ctx)
    if err != nil {
        log.Fatalf("Failed to get Registry information: %v", err)
    }
    
    fmt.Printf("Registry Name: %s\n", info.DbName)
    // Output: Registry Name: registry
    
    fmt.Printf("Total Packages: %d\n", info.DocCount)
    // Output: Total Packages: 2400000
}
```

### Using Proxy

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/scagogogo/npm-crawler/pkg/registry"
)

func main() {
    // Create options and configure proxy
    options := registry.NewOptions().
        SetRegistryURL("https://registry.npmjs.org").
        SetProxy("http://your-proxy-server:8080")
    
    // Create client with proxy
    client := registry.NewRegistry(options)
    
    ctx := context.Background()
    
    // Get package information
    pkg, err := client.GetPackageInformation(ctx, "react")
    if err != nil {
        log.Fatalf("Failed to get package information: %v", err)
    }
    
    fmt.Printf("Package Name: %s\n", pkg.Name)
    // Output: Package Name: react
    
    fmt.Printf("Description: %s\n", pkg.Description)
    // Output: Description: React is a JavaScript library for building user interfaces.
}
```

## API Documentation

### Registry Related

#### Creating Registry Client

```go
// NewRegistry creates a new Registry client instance
//
// Parameters:
//   - options: Optional configuration options, if not provided, default configuration will be used
//
// Return Value:
//   - *Registry: Newly created Registry client instance
func NewRegistry(options ...*Options) *Registry
```

#### Creating Clients for Specific Mirror Sources

```go
// Create Registry client using Taobao NPM mirror
func NewTaoBaoRegistry() *Registry

// Create Registry client using NPM Mirror (new domain for former Taobao mirror)
func NewNpmMirrorRegistry() *Registry

// Create Registry client using Huawei Cloud mirror
func NewHuaWeiCloudRegistry() *Registry

// Create Registry client using Tencent Cloud mirror
func NewTencentRegistry() *Registry

// Create Registry client using CNPM mirror
func NewCnpmRegistry() *Registry

// Create Registry client using Yarn official mirror
func NewYarnRegistry() *Registry

// Create Registry client using npmjs.com mirror
func NewNpmjsComRegistry() *Registry
```

#### Getting Registry Information

```go
// GetRegistryInformation gets the status information of NPM Registry
//
// Parameters:
//   - ctx: Context, can be used to cancel requests or set timeouts
//
// Return Value:
//   - *models.RegistryInformation: Registry status information
//   - error: Returns error if request fails
func (x *Registry) GetRegistryInformation(ctx context.Context) (*models.RegistryInformation, error)
```

#### Getting Package Information

```go
// GetPackageInformation gets detailed information of the specified NPM package
//
// Parameters:
//   - ctx: Context, can be used to cancel requests or set timeouts
//   - packageName: Name of the package to query, e.g. "react", "lodash", etc.
//
// Return Value:
//   - *models.Package: Detailed package information
//   - error: Returns error if request fails
func (x *Registry) GetPackageInformation(ctx context.Context, packageName string) (*models.Package, error)
```

### Configuration Options Related

#### Creating Options

```go
// NewOptions creates and returns a new default configuration options instance
//
// Default Configuration:
//   - RegistryURL: "https://registry.npmjs.org"
//   - Proxy: No proxy setting
func NewOptions() *Options
```

#### Setting Registry URL

```go
// SetRegistryURL sets the URL address of the NPM repository server
//
// Parameters:
//   - url: A valid NPM repository URL address string
//
// Return Value:
//   - *Options: Updated options object (supports method chaining)
func (o *Options) SetRegistryURL(url string) *Options
```

#### Setting Proxy

```go
// SetProxy sets the URL address of the HTTP proxy server
//
// Parameters:
//   - proxyUrl: HTTP proxy server URL address string
//
// Return Value:
//   - *Options: Updated options object (supports method chaining)
func (o *Options) SetProxy(proxyUrl string) *Options
```

### Main Models

#### Package

Represents the complete information structure of an NPM package:

```go
type Package struct {
    ID             string                 `json:"_id"`            // Package ID
    Rev            string                 `json:"_rev"`           // Revision number
    Name           string                 `json:"name"`           // Package name
    Description    string                 `json:"description"`    // Package description
    DistTags       map[string]string      `json:"dist-tags"`      // Distribution tags, such as latest
    Versions       map[string]Version     `json:"versions"`       // Version information mapping
    Maintainers    []Maintainer           `json:"maintainers"`    // Maintainer list
    Time           map[string]string      `json:"time"`           // Time information
    Repository     Repository             `json:"repository"`     // Repository information
    ReadMe         string                 `json:"readme"`         // README content
    ReadMeFilename string                 `json:"readmeFilename"` // README filename
    Homepage       string                 `json:"homepage"`       // Project homepage
    Bugs           map[string]interface{} `json:"bugs"`           // Bug tracking information
    License        string                 `json:"license"`        // License
    Users          map[string]bool        `json:"users"`          // User information
    Keywords       []string               `json:"keywords"`       // Keyword list
    Author         Author                 `json:"author"`         // Author information
    Contributors   []Contributor          `json:"contributors"`   // Contributor list
    Deprecated     string                 `json:"deprecated"`     // Deprecation notice
    Other          map[string]interface{} `json:"other"`          // Other fields
}
```

#### Version

Represents specific version information of an NPM package:

```go
type Version struct {
    Name            string               `json:"name"`            // Package name
    Version         string               `json:"version"`         // Version number
    Description     string               `json:"description"`     // Version description
    Main            string               `json:"main"`            // Main entry file
    Scripts         *Script              `json:"scripts"`         // Script commands
    Repository      *Repository          `json:"repository"`      // Repository
    Keywords        []string             `json:"keywords"`        // Keyword list
    Author          *User                `json:"author"`          // Author information
    License         string               `json:"license"`         // License
    Bugs            *Bugs                `json:"bugs"`            // Bug tracking
    Homepage        string               `json:"homepage"`        // Project homepage
    Dependencies    map[string]string    `json:"dependencies"`    // Runtime dependencies
    DevDependencies map[string]string    `json:"devDependencies"` // Development dependencies
    Dist            *Dist                `json:"dist"`            // Distribution information
    // Other fields...
}
```

#### RegistryInformation

Represents the status information of NPM Registry:

```go
type RegistryInformation struct {
    DbName            string `json:"db_name"`              // Database name
    DocCount          int    `json:"doc_count"`            // Total documents (packages)
    DocDelCount       int    `json:"doc_del_count"`        // Number of deleted documents
    UpdateSeq         int    `json:"update_seq"`           // Update sequence number
    PurgeSeq          int    `json:"purge_seq"`            // Purge sequence number
    CompactRunning    bool   `json:"compact_running"`      // Whether compaction is running
    DiskSize          int64  `json:"disk_size"`            // Disk usage size
    DataSize          int64  `json:"data_size"`            // Data size
    InstanceStartTime string `json:"instance_start_time"`  // Instance start time
    // Other fields...
}
```

## Supported Mirror Sources

| Mirror Source | URL | Region | Creation Method |
|---------------|-----|--------|-----------------|
| NPM Official | https://registry.npmjs.org | Global | `NewRegistry()` |
| Taobao NPM | https://registry.npm.taobao.org | China | `NewTaoBaoRegistry()` |
| NPM Mirror | https://registry.npmmirror.com | China | `NewNpmMirrorRegistry()` |
| Huawei Cloud | https://mirrors.huaweicloud.com/repository/npm | China | `NewHuaWeiCloudRegistry()` |
| Tencent Cloud | http://mirrors.cloud.tencent.com/npm | China | `NewTencentRegistry()` |
| CNPM | http://r.cnpmjs.org | China | `NewCnpmRegistry()` |
| Yarn | https://registry.yarnpkg.com | Global | `NewYarnRegistry()` |
| NPM CouchDB | https://skimdb.npmjs.com | Global | `NewNpmjsComRegistry()` |

## Contribution Guide

Contributions are welcome! Please follow these steps:

1. Fork this repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgements

- [NPM Registry](https://registry.npmjs.org) - Provides API and data
- [Go Requests](https://github.com/crawler-go-go-go/go-requests) - HTTP client library