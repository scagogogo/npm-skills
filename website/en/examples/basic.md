# Basic Usage Examples

This page contains basic usage examples of NPM Skills to help you get started quickly.

## Example 1: Get Package Information

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

func main() {
    // Create client
    client := registry.NewRegistry()
    ctx := context.Background()
    
    // Get package information
    pkg, err := client.GetPackageInformation(ctx, "react")
    if err != nil {
        log.Fatal(err)
    }
    
    // Display basic information
    fmt.Printf("Package: %s\n", pkg.Name)
    fmt.Printf("Description: %s\n", pkg.Description)
    fmt.Printf("Latest Version: %s\n", pkg.DistTags["latest"])
    fmt.Printf("License: %s\n", pkg.License)
    fmt.Printf("Homepage: %s\n", pkg.Homepage)
    fmt.Printf("Author: %s\n", pkg.Author.Name)
}
```

## Example 2: Using Different Mirror Sources

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

func main() {
    ctx := context.Background()
    
    // Use different mirror sources
    mirrors := map[string]*registry.Registry{
        "Official":      registry.NewRegistry(),
        "Taobao Mirror": registry.NewTaoBaoRegistry(),
        "NPM Mirror":    registry.NewNpmMirrorRegistry(),
        "Huawei Cloud":  registry.NewHuaWeiCloudRegistry(),
    }
    
    packageName := "lodash"
    
    for name, client := range mirrors {
        pkg, err := client.GetPackageInformation(ctx, packageName)
        if err != nil {
            fmt.Printf("%s: Failed to get - %v\n", name, err)
            continue
        }
        
        fmt.Printf("%s: %s (v%s)\n", name, pkg.Name, pkg.DistTags["latest"])
    }
}
```

## Example 3: Proxy Configuration

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

func main() {
    // Configure proxy
    options := registry.NewOptions().
        SetRegistryURL("https://registry.npmjs.org").
        SetProxy("http://proxy.example.com:8080")
    
    client := registry.NewRegistry(options)
    ctx := context.Background()
    
    pkg, err := client.GetPackageInformation(ctx, "express")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Got package info via proxy: %s\n", pkg.Name)
}
```

## Example 4: Search Packages

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

func main() {
    client := registry.NewRegistry()
    ctx := context.Background()
    
    // Search packages
    result, err := client.SearchPackages(ctx, "react ui component", 5)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Found %d results (time: %s)\n", result.Total, result.Time)
    fmt.Println("---")
    
    for i, obj := range result.Objects {
        pkg := obj.Package
        fmt.Printf("%d. %s (v%s)\n", i+1, pkg.Name, pkg.Version)
        fmt.Printf("   Description: %s\n", pkg.Description)
        fmt.Printf("   Score: %.2f\n", obj.Score.Final)
        fmt.Printf("   Link: %s\n", pkg.Links.NPM)
        fmt.Println()
    }
}
```

## Example 5: Get Specific Version Information

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

func main() {
    client := registry.NewRegistry()
    ctx := context.Background()
    
    // Get specific version
    version, err := client.GetPackageVersion(ctx, "react", "17.0.2")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Package: %s\n", version.Name)
    fmt.Printf("Version: %s\n", version.Version)
    fmt.Printf("Description: %s\n", version.Description)
    fmt.Printf("Main File: %s\n", version.Main)
    
    // Display dependencies
    if len(version.Dependencies) > 0 {
        fmt.Println("\nRuntime Dependencies:")
        for dep, ver := range version.Dependencies {
            fmt.Printf("  %s: %s\n", dep, ver)
        }
    }
    
    // Display dev dependencies
    if len(version.DevDependencies) > 0 {
        fmt.Println("\nDev Dependencies:")
        for dep, ver := range version.DevDependencies {
            fmt.Printf("  %s: %s\n", dep, ver)
        }
    }
}
```

## Example 6: Get Download Statistics

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

func main() {
    client := registry.NewRegistry()
    ctx := context.Background()
    
    packages := []string{"react", "vue", "angular"}
    periods := []string{"last-day", "last-week", "last-month"}
    
    for _, pkg := range packages {
        fmt.Printf("\n=== %s ===\n", pkg)
        
        for _, period := range periods {
            stats, err := client.GetDownloadStats(ctx, pkg, period)
            if err != nil {
                fmt.Printf("%s: Failed to get - %v\n", period, err)
                continue
            }
            
            fmt.Printf("%s: %d downloads\n", period, stats.Downloads)
        }
    }
}
```

## Example 7: Error Handling

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "log"
    "time"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

func main() {
    client := registry.NewRegistry()
    
    // Set timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    pkg, err := client.GetPackageInformation(ctx, "nonexistent-package-12345")
    if err != nil {
        switch {
        case errors.Is(err, context.DeadlineExceeded):
            log.Println("Request timeout")
        case errors.Is(err, context.Canceled):
            log.Println("Request canceled")
        default:
            log.Printf("Failed to get package info: %v\n", err)
        }
        return
    }
    
    fmt.Printf("Package info: %s\n", pkg.Name)
}
```

## Example 8: Compare Version Information

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

func main() {
    client := registry.NewRegistry()
    ctx := context.Background()
    
    packageName := "lodash"
    versions := []string{"4.17.20", "4.17.21"}
    
    fmt.Printf("Comparing different versions of %s:\n\n", packageName)
    
    for _, ver := range versions {
        version, err := client.GetPackageVersion(ctx, packageName, ver)
        if err != nil {
            log.Printf("Failed to get version %s: %v\n", ver, err)
            continue
        }
        
        fmt.Printf("Version %s:\n", version.Version)
        fmt.Printf("  Description: %s\n", version.Description)
        fmt.Printf("  Dependencies count: %d\n", len(version.Dependencies))
        fmt.Printf("  Dev dependencies count: %d\n", len(version.DevDependencies))
        
        if version.Deprecated != "" {
            fmt.Printf("  ⚠️  Deprecated: %s\n", version.Deprecated)
        }
        
        fmt.Println()
    }
}
```

## Example 9: Get Registry Status

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

func main() {
    // Check status of multiple registries
    registries := map[string]*registry.Registry{
        "NPM Official": registry.NewRegistry(),
        "Taobao Mirror": registry.NewTaoBaoRegistry(),
        "Huawei Cloud":  registry.NewHuaWeiCloudRegistry(),
    }
    
    ctx := context.Background()
    
    for name, client := range registries {
        fmt.Printf("=== %s ===\n", name)
        
        info, err := client.GetRegistryInformation(ctx)
        if err != nil {
            fmt.Printf("Failed to get: %v\n\n", err)
            continue
        }
        
        fmt.Printf("Database name: %s\n", info.DbName)
        fmt.Printf("Total packages: %d\n", info.DocCount)
        fmt.Printf("Deleted packages: %d\n", info.DocDelCount)
        fmt.Printf("Data size: %d MB\n", info.DataSize/(1024*1024))
        fmt.Printf("Disk usage: %d MB\n", info.DiskSize/(1024*1024))
        fmt.Printf("Instance start time: %s\n", info.InstanceStartTime)
        fmt.Println()
    }
}
```

## Example 10: Batch Processing

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
    
    packages := []string{
        "react", "vue", "angular", "lodash", "express",
        "axios", "moment", "underscore", "jquery", "bootstrap",
    }
    
    var wg sync.WaitGroup
    results := make(chan string, len(packages))
    
    // Concurrently get package information
    for _, pkg := range packages {
        wg.Add(1)
        go func(packageName string) {
            defer wg.Done()
            
            info, err := client.GetPackageInformation(ctx, packageName)
            if err != nil {
                results <- fmt.Sprintf("❌ %s: %v", packageName, err)
                return
            }
            
            results <- fmt.Sprintf("✅ %s: v%s (%s)", 
                info.Name, 
                info.DistTags["latest"], 
                info.Description[:50]+"...")
        }(pkg)
    }
    
    // Wait for all goroutines to complete
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // Collect and display results
    fmt.Println("Batch package information results:")
    for result := range results {
        fmt.Println(result)
    }
}
```

## Running Examples

Save any example code to a file (e.g., `example.go`), then run:

```bash
go mod init npm-skills-example
go get github.com/scagogogo/npm-skills
go run example.go
```

## Next Steps

- Check [Advanced Usage Examples](/en/examples/advanced) to learn more complex usage
- Read [Mirror Configuration](/en/examples/mirrors) to learn how to optimize network access
- Browse [API Documentation](/en/api/) to understand all available features 