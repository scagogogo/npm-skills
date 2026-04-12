# NPM Crawler API Reference

## Registry Client

### Creating a Client

```go
import "github.com/scagogogo/npm-crawler/pkg/registry"

// Default (official npmjs.org)
client := registry.NewRegistry()

// Pre-configured mirrors
client := registry.NewTaoBaoRegistry()
client := registry.NewNpmMirrorRegistry()
client := registry.NewHuaWeiCloudRegistry()
client := registry.NewTencentRegistry()
client := registry.NewCnpmRegistry()
client := registry.NewYarnRegistry()
client := registry.NewNpmjsComRegistry()

// Custom configuration
options := registry.NewOptions().
    SetRegistryURL("https://registry.npmjs.org").
    SetProxy("http://proxy:8080")
client := registry.NewRegistry(options)
```

### GetRegistryInformation

Returns NPM registry status and statistics.

```go
info, err := client.GetRegistryInformation(ctx)
```

Returns `*models.RegistryInformation`:
- `DbName` - Database name
- `DocCount` - Total number of packages
- `DiskSize` - Disk usage in bytes
- `DataSize` - Data size in bytes
- `UpdateSeq` - Update sequence number
- `InstanceStartTime` - When the registry instance started

### GetPackageInformation

Get complete package metadata.

```go
pkg, err := client.GetPackageInformation(ctx, "react")
```

Returns `*models.Package`:
- `Name` - Package name
- `Description` - Package description
- `DistTags` - Map of dist-tags (e.g., `latest` -> "18.2.0")
- `Versions` - Map of version strings to Version objects
- `Maintainers` - List of maintainer objects
- `Time` - Time map (created, modified, per-version times)
- `Repository` - Repository information
- `ReadMe` - README content
- `License` - Package license
- `Deprecated` - Deprecation notice (string or bool)

### SearchPackages

Search NPM registry by query string.

```go
result, err := client.SearchPackages(ctx, "react framework", 20)
```

Parameters:
- `query` - Search keywords
- `limit` - Max results (default 20)

Returns `*models.SearchResult` with `Objects` array containing matched packages.

### GetPackageVersion

Get metadata for a specific version.

```go
version, err := client.GetPackageVersion(ctx, "react", "18.0.0")
```

Returns `*models.Version`:
- `Name` - Package name
- `Version` - Version string
- `Description` - Version description
- `Dependencies` - Runtime dependencies map
- `DevDependencies` - Development dependencies map
- `Dist` - Distribution info (tarball URL, shasum, size)
- `Scripts` - npm scripts
- `Repository` - Repository URL
- `License` - License string

### GetDownloadStats

Get download statistics for a package.

```go
stats, err := client.GetDownloadStats(ctx, "react", "last-week")
```

Period options: `last-day`, `last-week`, `last-month`

Returns `*models.DownloadStats`:
- `Downloads` - Number of downloads
- `Start` - Period start date
- `End` - Period end date
- `Package` - Package name

### DownloadTarball

Download package as .tgz file.

```go
err := client.DownloadTarball(ctx, "react", "18.2.0", "./react-18.2.0.tgz")
```

Parameters:
- `packageName` - Package name
- `version` - Version string (or "latest")
- `destPath` - Local file path to save tarball

## Models

### Package

```go
type Package struct {
    ID             string
    Rev            string
    Name           string
    Description    string
    DistTags       map[string]string
    Versions       map[string]*Version
    Maintainers    []Maintainer
    Time           map[string]string
    Repository     Repository
    ReadMe         string
    Homepage       string
    License        string
    Deprecated     interface{}  // string or bool
}
```

### Version

```go
type Version struct {
    Name            string
    Version         string
    Description     string
    Main            string
    Scripts         *Script
    Repository      *Repository
    Dependencies    map[string]string
    DevDependencies map[string]string
    Dist            *Dist
    License         string
}
```

### Dist

```go
type Dist struct {
    Tarball string  // Download URL
    Shasum  string
    Size    int
    Integrity string  // SHA512 integrity hash
}
```

## Complete Go SDK Example

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/scagogogo/npm-crawler/pkg/registry"
)

func main() {
    client := registry.NewRegistry()
    ctx := context.Background()

    // Get package info
    pkg, err := client.GetPackageInformation(ctx, "axios")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Package: %s\n", pkg.Name)
    fmt.Printf("Latest: %s\n", pkg.DistTags["latest"])

    // Search packages (use simple keywords, avoid spaces or special chars)
    results, err := client.SearchPackages(ctx, "axios", 5)
    if err != nil {
        log.Fatal(err)
    }
    for _, r := range results.Objects {
        fmt.Printf("Found: %s\n", r.Package.Name)
    }

    // Get download stats
    stats, err := client.GetDownloadStats(ctx, "axios", "last-month")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Downloads: %d\n", stats.Downloads)
}
```
