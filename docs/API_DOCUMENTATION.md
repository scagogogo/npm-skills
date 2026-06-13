# NPM Crawler API Documentation

<div align="center">

[![Go Reference](https://pkg.go.dev/badge/github.com/scagogogo/npm-skills.svg)](https://pkg.go.dev/github.com/scagogogo/npm-skills)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

_Complete API Reference and Usage Guide_

[中文版本](API_DOCUMENTATION_zh.md) | [English Version](API_DOCUMENTATION.md)

</div>

## Table of Contents

- [Overview](#overview)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Core API Reference](#core-api-reference)
  - [Registry Client](#registry-client)
  - [Configuration Options](#configuration-options)
  - [Mirror Sources](#mirror-sources)
- [Data Models](#data-models)
- [Advanced Usage](#advanced-usage)
- [Error Handling](#error-handling)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)
- [Examples](#examples)

## Overview

NPM Crawler provides a comprehensive Go client for interacting with NPM registries. This documentation covers all available APIs, usage patterns, and best practices for integrating the library into your applications.

### Key Features

- **Multiple Registry Support**: Official NPM, Taobao, Huawei Cloud, and more
- **Proxy Configuration**: HTTP proxy support for enterprise environments
- **Context Support**: Cancellation and timeout handling
- **Type Safety**: Complete Go type definitions for NPM metadata
- **High Performance**: Optimized for concurrent access

## Installation

```bash
go get github.com/scagogogo/npm-skills
```

**Requirements:**
- Go 1.20 or higher
- Internet connectivity (or access to NPM registry mirrors)

## Quick Start

### Basic Package Information Retrieval

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
    // Create client with default settings (official NPM registry)
    client := registry.NewRegistry()
    
    // Create context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Get package information
    pkg, err := client.GetPackageInformation(ctx, "react")
    if err != nil {
        log.Fatalf("Failed to get package info: %v", err)
    }
    
    fmt.Printf("Package: %s\n", pkg.Name)
    fmt.Printf("Description: %s\n", pkg.Description)
    fmt.Printf("Latest Version: %s\n", pkg.DistTags["latest"])
    fmt.Printf("License: %s\n", pkg.License)
}
```

### Using Alternative Mirrors

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

func main() {
    // Use Taobao mirror for faster access in China
    client := registry.NewTaoBaoRegistry()
    
    ctx := context.Background()
    
    pkg, err := client.GetPackageInformation(ctx, "vue")
    if err != nil {
        log.Fatalf("Error: %v", err)
    }
    
    fmt.Printf("Package: %s from Taobao mirror\n", pkg.Name)
}
```

## Core API Reference

### Registry Client

#### Creating Registry Clients

##### `NewRegistry(options ...*Options) *Registry`

Creates a new Registry client with optional configuration.

**Parameters:**
- `options` (optional): Configuration options for the client

**Returns:**
- `*Registry`: Configured registry client

**Example:**
```go
// Default configuration (official NPM registry)
client := registry.NewRegistry()

// Custom configuration
options := registry.NewOptions().
    SetRegistryURL("https://registry.npmjs.org").
    SetProxy("http://proxy.company.com:8080")
client := registry.NewRegistry(options)
```

##### Mirror-Specific Constructors

```go
// Taobao NPM Mirror (China)
client := registry.NewTaoBaoRegistry()

// NPM Mirror (New Taobao domain)
client := registry.NewNpmMirrorRegistry()

// Huawei Cloud Mirror (China)
client := registry.NewHuaWeiCloudRegistry()

// Tencent Cloud Mirror (China)
client := registry.NewTencentRegistry()

// CNPM Mirror (China)
client := registry.NewCnpmRegistry()

// Yarn Official Mirror
client := registry.NewYarnRegistry()

// NPM CouchDB Mirror
client := registry.NewNpmjsComRegistry()
```

#### Core Methods

##### `GetPackageInformation(ctx context.Context, packageName string) (*models.Package, error)`

Retrieves detailed information about an NPM package.

**Parameters:**
- `ctx`: Context for cancellation and timeout control
- `packageName`: Name of the package to query

**Returns:**
- `*models.Package`: Complete package information
- `error`: Error if the request fails

**Example:**
```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

pkg, err := client.GetPackageInformation(ctx, "lodash")
if err != nil {
    return fmt.Errorf("failed to get package info: %w", err)
}

// Access package data
fmt.Printf("Name: %s\n", pkg.Name)
fmt.Printf("Latest: %s\n", pkg.DistTags["latest"])
fmt.Printf("Author: %s\n", pkg.Author.Name)

// Access specific version
if version, exists := pkg.Versions["4.17.21"]; exists {
    fmt.Printf("Version 4.17.21 dependencies: %+v\n", version.Dependencies)
}
```

##### `GetRegistryInformation(ctx context.Context) (*models.RegistryInformation, error)`

Retrieves status and metadata information about the NPM registry.

**Parameters:**
- `ctx`: Context for cancellation and timeout control

**Returns:**
- `*models.RegistryInformation`: Registry status information
- `error`: Error if the request fails

**Example:**
```go
info, err := client.GetRegistryInformation(ctx)
if err != nil {
    return fmt.Errorf("failed to get registry info: %w", err)
}

fmt.Printf("Registry: %s\n", info.DbName)
fmt.Printf("Total Packages: %d\n", info.DocCount)
fmt.Printf("Database Size: %d bytes\n", info.DataSize)
fmt.Printf("Disk Usage: %d bytes\n", info.DiskSize)
```

##### `GetOptions() *Options`

Returns the current configuration options for the registry client.

**Returns:**
- `*Options`: Current configuration options

**Example:**
```go
options := client.GetOptions()
fmt.Printf("Registry URL: %s\n", options.RegistryURL)
fmt.Printf("Proxy: %s\n", options.Proxy)
```

### Configuration Options

#### `NewOptions() *Options`

Creates a new Options instance with default values.

**Default Values:**
- `RegistryURL`: "https://registry.npmjs.org"
- `Proxy`: "" (no proxy)

**Example:**
```go
options := registry.NewOptions()
```

#### `SetRegistryURL(url string) *Options`

Sets the NPM registry URL. Supports method chaining.

**Parameters:**
- `url`: Registry URL (e.g., "https://registry.npmjs.org")

**Returns:**
- `*Options`: Updated options object for chaining

**Example:**
```go
options := registry.NewOptions().
    SetRegistryURL("https://registry.npmmirror.com")
```

#### `SetProxy(proxyUrl string) *Options`

Sets the HTTP proxy URL. Supports method chaining.

**Parameters:**
- `proxyUrl`: Proxy URL (e.g., "http://proxy.corp.com:8080")

**Returns:**
- `*Options`: Updated options object for chaining

**Example:**
```go
options := registry.NewOptions().
    SetProxy("http://user:pass@proxy.example.com:8080")

// Clear proxy
options.SetProxy("")
```

#### `GetHttpClient() (*http.Client, error)`

Creates an HTTP client with the current configuration.

**Returns:**
- `*http.Client`: Configured HTTP client
- `error`: Error if proxy URL is invalid

**Example:**
```go
options := registry.NewOptions().SetProxy("http://proxy.example.com:8080")
httpClient, err := options.GetHttpClient()
if err != nil {
    log.Fatalf("Failed to create HTTP client: %v", err)
}

// Use the client for custom requests
resp, err := httpClient.Get("https://example.com")
```

### Mirror Sources

| Mirror | URL | Region | Recommended Use |
|--------|-----|--------|----------------|
| **Official NPM** | https://registry.npmjs.org | Global | Production, global applications |
| **Taobao NPM** | https://registry.npm.taobao.org | China | Legacy applications (deprecated) |
| **NPM Mirror** | https://registry.npmmirror.com | China | China-based applications |
| **Huawei Cloud** | https://mirrors.huaweicloud.com/repository/npm | China | Huawei Cloud environments |
| **Tencent Cloud** | http://mirrors.cloud.tencent.com/npm | China | Tencent Cloud environments |
| **CNPM** | http://r.cnpmjs.org | China | Community-maintained |
| **Yarn** | https://registry.yarnpkg.com | Global | Yarn-based projects |
| **NPM CouchDB** | https://skimdb.npmjs.com | Global | Metadata-heavy applications |

## Data Models

### Package Model

The `Package` struct contains complete NPM package information:

```go
type Package struct {
    ID             string                 `json:"_id"`            // Package ID
    Rev            string                 `json:"_rev"`           // Revision number
    Name           string                 `json:"name"`           // Package name
    Description    string                 `json:"description"`    // Package description
    DistTags       map[string]string      `json:"dist-tags"`      // Distribution tags
    Versions       map[string]Version     `json:"versions"`       // All versions
    Maintainers    []Maintainer           `json:"maintainers"`    // Maintainers list
    Time           map[string]string      `json:"time"`           // Publishing times
    Repository     Repository             `json:"repository"`     // Repository info
    ReadMe         string                 `json:"readme"`         // README content
    Homepage       string                 `json:"homepage"`       // Project homepage
    License        string                 `json:"license"`        // License type
    Keywords       []string               `json:"keywords"`       // Keywords
    Author         Author                 `json:"author"`         // Author info
    Contributors   []Contributor          `json:"contributors"`   // Contributors
    // ... other fields
}
```

**Common Usage Patterns:**

```go
// Get latest version
latestVersion := pkg.DistTags["latest"]

// Check if package is deprecated
if pkg.Deprecated != "" {
    fmt.Printf("Warning: Package is deprecated: %s\n", pkg.Deprecated)
}

// List all available versions
for version := range pkg.Versions {
    fmt.Printf("Available version: %s\n", version)
}

// Get specific version details
if versionInfo, exists := pkg.Versions["1.0.0"]; exists {
    fmt.Printf("Dependencies: %+v\n", versionInfo.Dependencies)
    fmt.Printf("Dev Dependencies: %+v\n", versionInfo.DevDependencies)
}
```

### Version Model

```go
type Version struct {
    Name            string               `json:"name"`            // Package name
    Version         string               `json:"version"`         // Version number
    Description     string               `json:"description"`     // Description
    Main            string               `json:"main"`            // Entry point
    Scripts         *Script              `json:"scripts"`         // NPM scripts
    Dependencies    map[string]string    `json:"dependencies"`    // Runtime deps
    DevDependencies map[string]string    `json:"devDependencies"` // Dev deps
    Repository      *Repository          `json:"repository"`      // Repository
    License         string               `json:"license"`         // License
    Dist            *Dist                `json:"dist"`            // Distribution info
    // ... other fields
}
```

### RegistryInformation Model

```go
type RegistryInformation struct {
    DbName            string `json:"db_name"`              // Database name
    DocCount          int    `json:"doc_count"`            // Total packages
    DocDelCount       int    `json:"doc_del_count"`        // Deleted packages
    UpdateSeq         int    `json:"update_seq"`           // Update sequence
    CompactRunning    bool   `json:"compact_running"`      // Compaction status
    DiskSize          int64  `json:"disk_size"`            // Disk usage
    DataSize          int64  `json:"data_size"`            // Data size
    InstanceStartTime string `json:"instance_start_time"`  // Start time
    // ... other fields
}
```

## Advanced Usage

### Concurrent Package Queries

```go
package main

import (
    "context"
    "fmt"
    "sync"
    "time"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

func main() {
    client := registry.NewRegistry()
    packages := []string{"react", "vue", "angular", "svelte", "preact"}
    
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    var wg sync.WaitGroup
    results := make(chan string, len(packages))
    
    for _, pkg := range packages {
        wg.Add(1)
        go func(packageName string) {
            defer wg.Done()
            
            info, err := client.GetPackageInformation(ctx, packageName)
            if err != nil {
                results <- fmt.Sprintf("Error getting %s: %v", packageName, err)
                return
            }
            
            results <- fmt.Sprintf("%s: %s (latest: %s)", 
                info.Name, info.Description, info.DistTags["latest"])
        }(pkg)
    }
    
    go func() {
        wg.Wait()
        close(results)
    }()
    
    for result := range results {
        fmt.Println(result)
    }
}
```

### Registry Health Monitoring

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

func monitorRegistry(client *registry.Registry, interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
            
            info, err := client.GetRegistryInformation(ctx)
            if err != nil {
                fmt.Printf("Registry health check failed: %v\n", err)
                cancel()
                continue
            }
            
            fmt.Printf("Registry Status: %s, Packages: %d, Compact: %t\n",
                info.DbName, info.DocCount, info.CompactRunning)
            
            cancel()
        }
    }
}

func main() {
    client := registry.NewRegistry()
    monitorRegistry(client, 5*time.Minute)
}
```

### Package Dependency Analysis

```go
func analyzeDependencies(client *registry.Registry, packageName, version string) {
    ctx := context.Background()
    
    pkg, err := client.GetPackageInformation(ctx, packageName)
    if err != nil {
        log.Fatalf("Failed to get package info: %v", err)
    }
    
    versionInfo, exists := pkg.Versions[version]
    if !exists {
        fmt.Printf("Version %s not found\n", version)
        return
    }
    
    fmt.Printf("Analyzing %s@%s\n", packageName, version)
    fmt.Printf("Runtime Dependencies (%d):\n", len(versionInfo.Dependencies))
    
    for dep, ver := range versionInfo.Dependencies {
        fmt.Printf("  %s: %s\n", dep, ver)
    }
    
    fmt.Printf("\nDev Dependencies (%d):\n", len(versionInfo.DevDependencies))
    
    for dep, ver := range versionInfo.DevDependencies {
        fmt.Printf("  %s: %s\n", dep, ver)
    }
}
```

### Custom Registry Rotation

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

type RegistryPool struct {
    registries []*registry.Registry
    current    int
}

func NewRegistryPool() *RegistryPool {
    return &RegistryPool{
        registries: []*registry.Registry{
            registry.NewRegistry(),              // Official
            registry.NewNpmMirrorRegistry(),     // NPM Mirror
            registry.NewTaoBaoRegistry(),        // Taobao
            registry.NewHuaWeiCloudRegistry(),   // Huawei
        },
        current: 0,
    }
}

func (rp *RegistryPool) GetPackageWithFallback(ctx context.Context, packageName string) (*models.Package, error) {
    for i := 0; i < len(rp.registries); i++ {
        client := rp.registries[rp.current]
        
        pkg, err := client.GetPackageInformation(ctx, packageName)
        if err == nil {
            return pkg, nil
        }
        
        fmt.Printf("Registry %d failed, trying next: %v\n", rp.current, err)
        rp.current = (rp.current + 1) % len(rp.registries)
    }
    
    return nil, fmt.Errorf("all registries failed for package %s", packageName)
}
```

## Error Handling

### Common Error Types

1. **Network Errors**: Connection timeouts, DNS failures
2. **HTTP Errors**: 404 (package not found), 503 (service unavailable)
3. **Parse Errors**: Invalid JSON response
4. **Context Errors**: Timeout or cancellation

### Error Handling Best Practices

```go
func robustPackageQuery(client *registry.Registry, packageName string) error {
    // Create context with timeout and cancellation
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    pkg, err := client.GetPackageInformation(ctx, packageName)
    if err != nil {
        // Handle different error types
        switch {
        case ctx.Err() == context.DeadlineExceeded:
            return fmt.Errorf("request timed out for package %s", packageName)
        case ctx.Err() == context.Canceled:
            return fmt.Errorf("request was canceled for package %s", packageName)
        default:
            return fmt.Errorf("failed to get package %s: %w", packageName, err)
        }
    }
    
    // Validate received data
    if pkg.Name == "" {
        return fmt.Errorf("received invalid package data for %s", packageName)
    }
    
    if len(pkg.DistTags) == 0 {
        return fmt.Errorf("package %s has no distribution tags", packageName)
    }
    
    fmt.Printf("Successfully retrieved package: %s\n", pkg.Name)
    return nil
}
```

### Retry Logic with Exponential Backoff

```go
import (
    "math"
    "math/rand"
    "time"
)

func retryWithBackoff(fn func() error, maxRetries int) error {
    for attempt := 0; attempt < maxRetries; attempt++ {
        err := fn()
        if err == nil {
            return nil
        }
        
        if attempt == maxRetries-1 {
            return fmt.Errorf("max retries exceeded: %w", err)
        }
        
        // Exponential backoff with jitter
        backoff := time.Duration(math.Pow(2, float64(attempt))) * time.Second
        jitter := time.Duration(rand.Intn(1000)) * time.Millisecond
        time.Sleep(backoff + jitter)
        
        fmt.Printf("Retry attempt %d after error: %v\n", attempt+1, err)
    }
    
    return nil
}

// Usage
err := retryWithBackoff(func() error {
    _, err := client.GetPackageInformation(ctx, "some-package")
    return err
}, 3)
```

## Best Practices

### 1. Context Usage

Always use contexts with appropriate timeouts:

```go
// Good: Specific timeout for the operation
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

// Better: Cancellable context for user operations
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// Handle cancellation signals
go func() {
    <-interrupt
    cancel()
}()
```

### 2. Registry Selection

Choose registries based on your deployment region:

```go
func selectOptimalRegistry() *registry.Registry {
    // Detect region or use configuration
    region := os.Getenv("DEPLOYMENT_REGION")
    
    switch region {
    case "china", "cn":
        return registry.NewNpmMirrorRegistry()
    case "huawei-cloud":
        return registry.NewHuaWeiCloudRegistry()
    default:
        return registry.NewRegistry() // Official NPM
    }
}
```

### 3. Resource Management

```go
// Use sync.Pool for high-frequency operations
var clientPool = sync.Pool{
    New: func() interface{} {
        return registry.NewRegistry()
    },
}

func getClient() *registry.Registry {
    return clientPool.Get().(*registry.Registry)
}

func putClient(client *registry.Registry) {
    clientPool.Put(client)
}
```

### 4. Configuration Management

```go
type Config struct {
    RegistryURL string `env:"NPM_REGISTRY_URL" default:"https://registry.npmjs.org"`
    ProxyURL    string `env:"HTTP_PROXY"`
    Timeout     time.Duration `env:"NPM_TIMEOUT" default:"30s"`
}

func createClientFromConfig(cfg Config) *registry.Registry {
    options := registry.NewOptions().SetRegistryURL(cfg.RegistryURL)
    
    if cfg.ProxyURL != "" {
        options.SetProxy(cfg.ProxyURL)
    }
    
    return registry.NewRegistry(options)
}
```

## Troubleshooting

### Common Issues and Solutions

#### 1. Connection Timeouts

**Problem**: Requests timeout frequently

**Solutions**:
```go
// Increase timeout
ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

// Use faster mirror
client := registry.NewNpmMirrorRegistry() // For China users

// Check proxy configuration
options := registry.NewOptions().SetProxy("") // Disable proxy
client := registry.NewRegistry(options)
```

#### 2. Package Not Found (404)

**Problem**: Package exists but returns 404

**Solutions**:
```go
// Check package name spelling
pkg, err := client.GetPackageInformation(ctx, "react") // not "React"

// Try different registry
client := registry.NewYarnRegistry()

// Verify package exists on npmjs.com first
```

#### 3. Proxy Issues

**Problem**: Corporate proxy blocking requests

**Solutions**:
```go
// Use authenticated proxy
options := registry.NewOptions().
    SetProxy("http://username:password@proxy.corp.com:8080")

// Try different proxy format
options.SetProxy("http://proxy.corp.com:8080")

// Bypass proxy for specific domains (requires custom HTTP client)
```

#### 4. Rate Limiting

**Problem**: Too many requests causing rate limits

**Solutions**:
```go
// Add delays between requests
time.Sleep(100 * time.Millisecond)

// Use different mirrors for load distribution
registries := []*registry.Registry{
    registry.NewRegistry(),
    registry.NewNpmMirrorRegistry(),
    registry.NewTaoBaoRegistry(),
}

// Implement request queuing
type RequestQueue struct {
    ch chan func()
}

func (rq *RequestQueue) Execute(fn func()) {
    rq.ch <- fn
}
```

#### 5. Memory Usage with Large Packages

**Problem**: High memory usage when querying large packages

**Solutions**:
```go
// Only access needed fields
pkg, err := client.GetPackageInformation(ctx, packageName)
if err != nil {
    return err
}

// Extract only required data immediately
result := struct {
    Name        string
    Description string
    Latest      string
}{
    Name:        pkg.Name,
    Description: pkg.Description,
    Latest:      pkg.DistTags["latest"],
}

// Clear original package data
pkg = nil
```

## Examples

### Example 1: Package Information Dashboard

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

type PackageInfo struct {
    Name         string    `json:"name"`
    Description  string    `json:"description"`
    Latest       string    `json:"latest"`
    License      string    `json:"license"`
    LastModified time.Time `json:"lastModified"`
    Homepage     string    `json:"homepage"`
}

func packageHandler(w http.ResponseWriter, r *http.Request) {
    packageName := r.URL.Query().Get("package")
    if packageName == "" {
        http.Error(w, "package parameter required", http.StatusBadRequest)
        return
    }
    
    client := registry.NewRegistry()
    ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
    defer cancel()
    
    pkg, err := client.GetPackageInformation(ctx, packageName)
    if err != nil {
        http.Error(w, fmt.Sprintf("Package not found: %v", err), http.StatusNotFound)
        return
    }
    
    // Parse last modified time
    lastModified, _ := time.Parse(time.RFC3339, pkg.Time["modified"])
    
    info := PackageInfo{
        Name:         pkg.Name,
        Description:  pkg.Description,
        Latest:       pkg.DistTags["latest"],
        License:      pkg.License,
        LastModified: lastModified,
        Homepage:     pkg.Homepage,
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(info)
}

func main() {
    http.HandleFunc("/package", packageHandler)
    fmt.Println("Server starting on :8080")
    http.ListenAndServe(":8080", nil)
}
```

### Example 2: Dependency Tree Builder

```go
package main

import (
    "context"
    "fmt"
    "strings"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

type DependencyNode struct {
    Name         string            `json:"name"`
    Version      string            `json:"version"`
    Dependencies []*DependencyNode `json:"dependencies,omitempty"`
}

func buildDependencyTree(client *registry.Registry, ctx context.Context, packageName, version string, depth int, maxDepth int) (*DependencyNode, error) {
    if depth > maxDepth {
        return &DependencyNode{Name: packageName, Version: version}, nil
    }
    
    pkg, err := client.GetPackageInformation(ctx, packageName)
    if err != nil {
        return nil, fmt.Errorf("failed to get package %s: %w", packageName, err)
    }
    
    node := &DependencyNode{
        Name:    packageName,
        Version: version,
    }
    
    // Get version info
    versionInfo, exists := pkg.Versions[version]
    if !exists {
        // Try latest version
        if latest, ok := pkg.DistTags["latest"]; ok {
            versionInfo, exists = pkg.Versions[latest]
            node.Version = latest
        }
    }
    
    if !exists {
        return node, nil
    }
    
    // Build dependencies
    for depName, depVersion := range versionInfo.Dependencies {
        // Clean version specification (remove ^, ~, etc.)
        cleanVersion := strings.TrimLeft(depVersion, "^~>=<")
        
        depNode, err := buildDependencyTree(client, ctx, depName, cleanVersion, depth+1, maxDepth)
        if err != nil {
            fmt.Printf("Warning: Could not resolve dependency %s: %v\n", depName, err)
            continue
        }
        
        node.Dependencies = append(node.Dependencies, depNode)
    }
    
    return node, nil
}

func main() {
    client := registry.NewRegistry()
    ctx := context.Background()
    
    tree, err := buildDependencyTree(client, ctx, "express", "4.18.2", 0, 2)
    if err != nil {
        log.Fatalf("Failed to build dependency tree: %v", err)
    }
    
    printTree(tree, 0)
}

func printTree(node *DependencyNode, indent int) {
    prefix := strings.Repeat("  ", indent)
    fmt.Printf("%s%s@%s\n", prefix, node.Name, node.Version)
    
    for _, dep := range node.Dependencies {
        printTree(dep, indent+1)
    }
}
```

### Example 3: Package Security Audit

```go
package main

import (
    "context"
    "fmt"
    "regexp"
    "strings"
    "time"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

type SecurityAudit struct {
    PackageName   string    `json:"packageName"`
    CurrentVersion string   `json:"currentVersion"`
    LatestVersion string    `json:"latestVersion"`
    IsOutdated    bool      `json:"isOutdated"`
    IsDeprecated  bool      `json:"isDeprecated"`
    HasVulnerabilities bool `json:"hasVulnerabilities"`
    LastUpdate    time.Time `json:"lastUpdate"`
    Maintainers   int       `json:"maintainersCount"`
    License       string    `json:"license"`
    Issues        []string  `json:"issues"`
}

func auditPackage(client *registry.Registry, ctx context.Context, packageName, currentVersion string) (*SecurityAudit, error) {
    pkg, err := client.GetPackageInformation(ctx, packageName)
    if err != nil {
        return nil, fmt.Errorf("failed to get package info: %w", err)
    }
    
    audit := &SecurityAudit{
        PackageName:   packageName,
        CurrentVersion: currentVersion,
        LatestVersion: pkg.DistTags["latest"],
        IsDeprecated:  pkg.Deprecated != "",
        Maintainers:   len(pkg.Maintainers),
        License:       pkg.License,
    }
    
    // Check if version is outdated
    audit.IsOutdated = currentVersion != audit.LatestVersion
    
    // Parse last update time
    if timeStr, exists := pkg.Time[audit.LatestVersion]; exists {
        audit.LastUpdate, _ = time.Parse(time.RFC3339, timeStr)
    }
    
    // Check for potential security issues
    audit.Issues = checkSecurityIssues(pkg)
    audit.HasVulnerabilities = len(audit.Issues) > 0
    
    return audit, nil
}

func checkSecurityIssues(pkg *models.Package) []string {
    var issues []string
    
    // Check for deprecated packages
    if pkg.Deprecated != "" {
        issues = append(issues, fmt.Sprintf("Package is deprecated: %s", pkg.Deprecated))
    }
    
    // Check for suspicious patterns in description
    suspiciousPatterns := []string{
        "bitcoin", "cryptocurrency", "mining", "keylogger",
        "password", "steal", "backdoor",
    }
    
    description := strings.ToLower(pkg.Description)
    for _, pattern := range suspiciousPatterns {
        if strings.Contains(description, pattern) {
            issues = append(issues, fmt.Sprintf("Suspicious keyword in description: %s", pattern))
        }
    }
    
    // Check for suspicious repository URLs
    if pkg.Repository.Url != "" {
        if matched, _ := regexp.MatchString(`(bitbucket\.org|gitlab\.com|github\.com)`, pkg.Repository.Url); !matched {
            issues = append(issues, "Repository hosted on non-standard platform")
        }
    }
    
    // Check maintainer count
    if len(pkg.Maintainers) == 0 {
        issues = append(issues, "Package has no maintainers")
    } else if len(pkg.Maintainers) > 10 {
        issues = append(issues, "Package has unusually many maintainers")
    }
    
    // Check for missing license
    if pkg.License == "" {
        issues = append(issues, "Package has no license specified")
    }
    
    return issues
}

func main() {
    client := registry.NewRegistry()
    ctx := context.Background()
    
    packages := map[string]string{
        "express":    "4.17.1",
        "lodash":     "4.17.20",
        "react":      "17.0.0",
        "vue":        "2.6.12",
    }
    
    fmt.Println("Security Audit Report")
    fmt.Println("====================")
    
    for packageName, version := range packages {
        audit, err := auditPackage(client, ctx, packageName, version)
        if err != nil {
            fmt.Printf("Failed to audit %s: %v\n", packageName, err)
            continue
        }
        
        fmt.Printf("\nPackage: %s\n", audit.PackageName)
        fmt.Printf("Current Version: %s\n", audit.CurrentVersion)
        fmt.Printf("Latest Version: %s\n", audit.LatestVersion)
        fmt.Printf("Status: ")
        
        if audit.IsDeprecated {
            fmt.Print("DEPRECATED ")
        }
        if audit.IsOutdated {
            fmt.Print("OUTDATED ")
        }
        if audit.HasVulnerabilities {
            fmt.Print("ISSUES ")
        }
        if !audit.IsDeprecated && !audit.IsOutdated && !audit.HasVulnerabilities {
            fmt.Print("OK")
        }
        fmt.Println()
        
        if len(audit.Issues) > 0 {
            fmt.Println("Issues:")
            for _, issue := range audit.Issues {
                fmt.Printf("  - %s\n", issue)
            }
        }
    }
}
```

---

## Contributing

We welcome contributions to improve this documentation! Please feel free to:

- Report issues or inconsistencies
- Suggest new examples or use cases
- Add translations
- Improve existing content

For more information, see our [Contributing Guide](https://github.com/scagogogo/npm-skills/blob/main/README.md#contribution-guide).

## License

This documentation is part of the NPM Crawler project and is licensed under the MIT License. See [LICENSE](https://github.com/scagogogo/npm-skills/blob/main/LICENSE) for details.
