# API Overview

NPM Crawler provides a simple yet powerful API for interacting with NPM Registry. This documentation covers all available API methods, parameters, and return values.

## Core API

### Registry Client

Registry is the core client of NPM Crawler, providing all functionality to access NPM Registry.

```go
import "github.com/scagogogo/npm-crawler/pkg/registry"

// Create default client
client := registry.NewRegistry()

// Use custom configuration
options := registry.NewOptions().
    SetRegistryURL("https://registry.npmjs.org").
    SetProxy("http://proxy.example.com:8080")
client := registry.NewRegistry(options)
```

## Main Methods

### Package Information Query

#### GetPackageInformation
Get complete information for a specified NPM package.

```go
func (r *Registry) GetPackageInformation(ctx context.Context, packageName string) (*models.Package, error)
```

**Parameters:**
- `ctx` - Context for cancellation and timeout control
- `packageName` - Name of the package to query

**Returns:**
- `*models.Package` - Complete package information
- `error` - Error information

**Example:**
```go
pkg, err := client.GetPackageInformation(ctx, "react")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Package: %s, Latest Version: %s\n", pkg.Name, pkg.DistTags["latest"])
```

#### GetPackageVersion
Get information for a specific version of a package.

```go
func (r *Registry) GetPackageVersion(ctx context.Context, packageName, version string) (*models.Version, error)
```

**Parameters:**
- `ctx` - Context
- `packageName` - Package name
- `version` - Version number (e.g., "1.0.0" or "latest")

**Example:**
```go
version, err := client.GetPackageVersion(ctx, "react", "18.2.0")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Version: %s, Dependencies: %+v\n", version.Version, version.Dependencies)
```

### Package Search

#### SearchPackages
Search NPM packages.

```go
func (r *Registry) SearchPackages(ctx context.Context, query string, limit int) (*models.SearchResult, error)
```

**Parameters:**
- `ctx` - Context
- `query` - Search keywords
- `limit` - Limit on number of results (default 20)

**Example:**
```go
result, err := client.SearchPackages(ctx, "react ui", 10)
if err != nil {
    log.Fatal(err)
}

for _, obj := range result.Objects {
    fmt.Printf("- %s: %s\n", obj.Package.Name, obj.Package.Description)
}
```

### Statistics

#### GetRegistryInformation
Get NPM Registry status information.

```go
func (r *Registry) GetRegistryInformation(ctx context.Context) (*models.RegistryInformation, error)
```

**Example:**
```go
info, err := client.GetRegistryInformation(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Registry: %s, Total Packages: %d\n", info.DbName, info.DocCount)
```

#### GetDownloadStats
Get download statistics for a package.

```go
func (r *Registry) GetDownloadStats(ctx context.Context, packageName, period string) (*models.DownloadStats, error)
```

**Parameters:**
- `packageName` - Package name
- `period` - Statistics period ("last-day", "last-week", "last-month")

**Example:**
```go
stats, err := client.GetDownloadStats(ctx, "react", "last-week")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Downloads: %d\n", stats.Downloads)
```

#### DownloadTarball
Download an NPM package tarball to a local file.

```go
func (r *Registry) DownloadTarball(ctx context.Context, packageName, version, destPath string) error
```

**Parameters:**
- `packageName` - Package name
- `version` - Version number or tag (e.g., "latest")
- `destPath` - Local file path to save the tarball

**Example:**
```go
err := client.DownloadTarball(ctx, "react", "18.2.0", "./react.tgz")
```

## Configuration Options

### Options
Configuration options for Registry client.

```go
type Options struct {
    RegistryURL string  // NPM registry URL
    Proxy       string  // HTTP proxy URL
}
```

**Methods:**
- `NewOptions()` - Create default options
- `SetRegistryURL(url string)` - Set Registry URL
- `SetProxy(proxyUrl string)` - Set proxy
- `GetHttpClient()` - Get configured HTTP client

## Mirror Source Support

NPM Crawler has built-in support for multiple mirror sources:

| Function | Mirror Source | URL |
|----------|---------------|-----|
| `NewRegistry()` | Official NPM | https://registry.npmjs.org |
| `NewTaoBaoRegistry()` | Taobao Mirror | https://registry.npm.taobao.org |
| `NewNpmMirrorRegistry()` | NPM Mirror | https://registry.npmmirror.com |
| `NewHuaWeiCloudRegistry()` | Huawei Cloud | https://mirrors.huaweicloud.com/repository/npm |
| `NewTencentRegistry()` | Tencent Cloud | http://mirrors.cloud.tencent.com/npm |
| `NewCnpmRegistry()` | CNPM | http://r.cnpmjs.org |
| `NewYarnRegistry()` | Yarn | https://registry.yarnpkg.com |

## Data Models

### Package
Represents complete NPM package information.

```go
type Package struct {
    ID           string                 `json:"_id"`
    Name         string                 `json:"name"`
    Description  string                 `json:"description"`
    DistTags     map[string]string      `json:"dist-tags"`
    Versions     map[string]Version     `json:"versions"`
    Maintainers  []Maintainer           `json:"maintainers"`
    Time         map[string]string      `json:"time"`
    Repository   Repository             `json:"repository"`
    Homepage     string                 `json:"homepage"`
    License      string                 `json:"license"`
    Keywords     []string               `json:"keywords"`
    Author       Author                 `json:"author"`
    // ... other fields
}
```

### Version
Represents specific version information of a package.

```go
type Version struct {
    Name            string               `json:"name"`
    Version         string               `json:"version"`
    Description     string               `json:"description"`
    Dependencies    map[string]string    `json:"dependencies"`
    DevDependencies map[string]string    `json:"devDependencies"`
    Dist            *Dist                `json:"dist"`
    // ... other fields
}
```

### SearchResult
Represents search results.

```go
type SearchResult struct {
    Objects []SearchObject `json:"objects"`
    Total   int            `json:"total"`
    Time    string         `json:"time"`
}
```

### DownloadStats
Represents download statistics.

```go
type DownloadStats struct {
    Downloads int    `json:"downloads"`
    Start     string `json:"start"`
    End       string `json:"end"`
    Package   string `json:"package"`
}
```

## Error Handling

All API methods return error type, proper error handling is recommended for production environments:

```go
import (
    "context"
    "errors"
    "time"
)

ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

pkg, err := client.GetPackageInformation(ctx, "react")
if err != nil {
    switch {
    case errors.Is(err, context.DeadlineExceeded):
        // Handle timeout
    case errors.Is(err, context.Canceled):
        // Handle cancellation
    default:
        // Handle other errors
    }
}
```

## Best Practices

1. **Use Context**: Always pass appropriate context to support cancellation and timeout
2. **Choose Appropriate Mirror**: Select the fastest mirror source based on network environment
3. **Concurrent Safety**: Registry client is concurrent-safe and can be used in multiple goroutines
4. **Error Handling**: Implement proper error handling and retry mechanisms
5. **Proxy Configuration**: Configure proxy correctly in enterprise environments

## Next Steps

- Check [Chinese detailed API documentation](/api/registry) or [English API documentation](/en/api/registry)
- Browse [example code](/en/examples/basic) to learn specific usage
- Visit [GitHub](https://github.com/scagogogo/npm-crawler) to view complete source code 