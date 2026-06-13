# Getting Started

This guide will help you quickly get started with NPM Crawler.

## Requirements

- Go 1.20 or higher
- Network connection (or access to NPM registry mirrors)

## Installation

Install NPM Crawler using Go modules:

```bash
go get github.com/scagogogo/npm-skills
```

## Basic Usage

### 1. Create Registry Client

```go
package main

import (
    "github.com/scagogogo/npm-skills/pkg/registry"
)

func main() {
    // Use default configuration (official NPM Registry)
    client := registry.NewRegistry()
    
    // Or use custom configuration
    options := registry.NewOptions().
        SetRegistryURL("https://registry.npmjs.org").
        SetProxy("http://proxy.example.com:8080")
    client := registry.NewRegistry(options)
}
```

### 2. Get Package Information

```go
ctx := context.Background()

// Get complete package information
pkg, err := client.GetPackageInformation(ctx, "react")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Package Name: %s\n", pkg.Name)
fmt.Printf("Description: %s\n", pkg.Description)
fmt.Printf("Latest Version: %s\n", pkg.DistTags["latest"])
fmt.Printf("Homepage: %s\n", pkg.Homepage)
```

### 3. Get Registry Status

```go
info, err := client.GetRegistryInformation(ctx)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Registry: %s\n", info.DbName)
fmt.Printf("Total Packages: %d\n", info.DocCount)
fmt.Printf("Data Size: %d bytes\n", info.DataSize)
```

### 4. Search Packages

```go
// Search packages, limit to 10 results
result, err := client.SearchPackages(ctx, "react ui", 10)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Found %d results\n", result.Total)
for _, obj := range result.Objects {
    fmt.Printf("- %s: %s\n", obj.Package.Name, obj.Package.Description)
}
```

### 5. Get Specific Version Information

```go
// Get detailed information for a specific version
version, err := client.GetPackageVersion(ctx, "react", "18.2.0")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Version: %s\n", version.Version)
fmt.Printf("Dependencies: %+v\n", version.Dependencies)
```

### 6. Get Download Statistics

```go
// Get download statistics for the past week
stats, err := client.GetDownloadStats(ctx, "react", "last-week")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Package: %s\n", stats.Package)
fmt.Printf("Downloads: %d\n", stats.Downloads)
fmt.Printf("Period: %s to %s\n", stats.Start, stats.End)
```

## Using Different Mirror Sources

NPM Crawler has built-in support for multiple mirror sources, especially suitable for users in mainland China:

```go
// Use Taobao mirror
taobaoClient := registry.NewTaoBaoRegistry()

// Use Huawei Cloud mirror
huaweiClient := registry.NewHuaWeiCloudRegistry()

// Use NPM Mirror
npmMirrorClient := registry.NewNpmMirrorRegistry()

// Use Tencent Cloud mirror
tencentClient := registry.NewTencentRegistry()

// Use CNPM mirror
cnpmClient := registry.NewCnpmRegistry()

// Use Yarn mirror
yarnClient := registry.NewYarnRegistry()
```

## Configuring Proxy

If you're in a corporate network environment, you may need to configure an HTTP proxy:

```go
options := registry.NewOptions().
    SetRegistryURL("https://registry.npmjs.org").
    SetProxy("http://proxy.corp.com:8080")

client := registry.NewRegistry(options)

// Or use proxy with authentication
options.SetProxy("http://username:password@proxy.corp.com:8080")
```

## Error Handling

It's recommended to add proper error handling in production environments:

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
    client := registry.NewRegistry()
    
    // Set timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    pkg, err := client.GetPackageInformation(ctx, "react")
    if err != nil {
        switch {
        case ctx.Err() == context.DeadlineExceeded:
            log.Fatal("Request timeout")
        case ctx.Err() == context.Canceled:
            log.Fatal("Request canceled")
        default:
            log.Fatalf("Failed to get package information: %v", err)
        }
    }
    
    fmt.Printf("Successfully got package information: %s\n", pkg.Name)
}
```

## Concurrent Usage

NPM Crawler is thread-safe and can be used in concurrent environments:

```go
package main

import (
    "context"
    "fmt"
    "sync"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

func main() {
    client := registry.NewRegistry()
    ctx := context.Background()
    
    packages := []string{"react", "vue", "angular", "lodash", "express"}
    
    var wg sync.WaitGroup
    for _, pkgName := range packages {
        wg.Add(1)
        go func(name string) {
            defer wg.Done()
            
            pkg, err := client.GetPackageInformation(ctx, name)
            if err != nil {
                fmt.Printf("Failed to get %s: %v\n", name, err)
                return
            }
            
            fmt.Printf("%s: %s\n", pkg.Name, pkg.DistTags["latest"])
        }(pkgName)
    }
    
    wg.Wait()
}
```

## Next Steps

- Check the complete [API Documentation](/en/api/) to learn about all available features
- Browse [example code](/en/examples/basic) to learn more usage patterns
- Visit the [GitHub repository](https://github.com/scagogogo/npm-skills) to view source code and contribution guidelines 