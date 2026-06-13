# Advanced Usage Examples

This page demonstrates advanced usage patterns and features of NPM Crawler.

## Custom HTTP Client Configuration

```go
package main

import (
    "context"
    "crypto/tls"
    "fmt"
    "net/http"
    "time"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

func main() {
    // Create custom HTTP client with advanced configuration
    httpClient := &http.Client{
        Timeout: 30 * time.Second,
        Transport: &http.Transport{
            MaxIdleConns:        100,
            MaxIdleConnsPerHost: 10,
            IdleConnTimeout:     90 * time.Second,
            TLSClientConfig: &tls.Config{
                InsecureSkipVerify: false,
            },
        },
    }
    
    options := registry.NewOptions().
        SetHTTPClient(httpClient).
        SetRegistryURL("https://registry.npmjs.org")
    
    client := registry.NewRegistry(options)
    ctx := context.Background()
    
    pkg, err := client.GetPackageInformation(ctx, "react")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("Package: %s, Version: %s\n", pkg.Name, pkg.DistTags["latest"])
}
```

## Parallel Package Processing

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
    
    packages := []string{
        "react", "vue", "angular", "svelte", "solid-js",
        "lodash", "underscore", "ramda", "moment", "dayjs",
    }
    
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    type Result struct {
        Package string
        Version string
        Error   error
    }
    
    results := make(chan Result, len(packages))
    var wg sync.WaitGroup
    
    // Process packages concurrently
    for _, pkg := range packages {
        wg.Add(1)
        go func(packageName string) {
            defer wg.Done()
            
            info, err := client.GetPackageInformation(ctx, packageName)
            if err != nil {
                results <- Result{Package: packageName, Error: err}
                return
            }
            
            results <- Result{
                Package: packageName, 
                Version: info.DistTags["latest"],
            }
        }(pkg)
    }
    
    // Close results channel when all goroutines complete
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // Collect and display results
    fmt.Println("Package Processing Results:")
    fmt.Println("==========================")
    
    for result := range results {
        if result.Error != nil {
            fmt.Printf("❌ %s: %v\n", result.Package, result.Error)
        } else {
            fmt.Printf("✅ %s: v%s\n", result.Package, result.Version)
        }
    }
}
```

## Registry Health Check

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

func healthCheck(name string, client *registry.Registry) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    start := time.Now()
    
    // Test basic connectivity
    info, err := client.GetRegistryInformation(ctx)
    if err != nil {
        fmt.Printf("🔴 %s: UNHEALTHY - %v\n", name, err)
        return
    }
    
    latency := time.Since(start)
    
    // Additional health checks
    pkg, err := client.GetPackageInformation(ctx, "lodash")
    if err != nil {
        fmt.Printf("🟡 %s: PARTIAL - Registry info OK but package query failed: %v\n", name, err)
        return
    }
    
    fmt.Printf("🟢 %s: HEALTHY - %d packages, latency: %v\n", 
        name, info.DocCount, latency)
}

func main() {
    registries := map[string]*registry.Registry{
        "NPM Official":   registry.NewRegistry(),
        "Taobao Mirror":  registry.NewTaoBaoRegistry(),
        "NPM Mirror":     registry.NewNpmMirrorRegistry(),
        "Huawei Cloud":   registry.NewHuaWeiCloudRegistry(),
    }
    
    fmt.Println("Registry Health Check Results:")
    fmt.Println("===============================")
    
    for name, client := range registries {
        healthCheck(name, client)
    }
}
```

## Package Dependency Analysis

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

func analyzeDependencies(client *registry.Registry, packageName, version string) {
    ctx := context.Background()
    
    ver, err := client.GetPackageVersion(ctx, packageName, version)
    if err != nil {
        log.Printf("Failed to get version info: %v\n", err)
        return
    }
    
    fmt.Printf("\n📦 Analyzing %s@%s\n", packageName, version)
    fmt.Printf("=====================================\n")
    
    // Runtime dependencies
    if len(ver.Dependencies) > 0 {
        fmt.Printf("\n🔗 Runtime Dependencies (%d):\n", len(ver.Dependencies))
        for dep, depVer := range ver.Dependencies {
            fmt.Printf("  ├─ %s: %s\n", dep, depVer)
        }
    }
    
    // Dev dependencies
    if len(ver.DevDependencies) > 0 {
        fmt.Printf("\n🛠️  Dev Dependencies (%d):\n", len(ver.DevDependencies))
        for dep, depVer := range ver.DevDependencies {
            fmt.Printf("  ├─ %s: %s\n", dep, depVer)
        }
    }
    
    // Optional dependencies
    if len(ver.OptionalDependencies) > 0 {
        fmt.Printf("\n🔧 Optional Dependencies (%d):\n", len(ver.OptionalDependencies))
        for dep, depVer := range ver.OptionalDependencies {
            fmt.Printf("  ├─ %s: %s\n", dep, depVer)
        }
    }
    
    // Peer dependencies
    if len(ver.PeerDependencies) > 0 {
        fmt.Printf("\n👥 Peer Dependencies (%d):\n", len(ver.PeerDependencies))
        for dep, depVer := range ver.PeerDependencies {
            fmt.Printf("  ├─ %s: %s\n", dep, depVer)
        }
    }
    
    // Bundle dependencies
    if len(ver.BundleDependencies) > 0 {
        fmt.Printf("\n📦 Bundled Dependencies (%d):\n", len(ver.BundleDependencies))
        for _, dep := range ver.BundleDependencies {
            fmt.Printf("  ├─ %s\n", dep)
        }
    }
}

func main() {
    client := registry.NewRegistry()
    
    // Analyze popular packages
    packages := []struct {
        name    string
        version string
    }{
        {"react", "18.2.0"},
        {"vue", "3.3.0"},
        {"express", "4.18.2"},
    }
    
    for _, pkg := range packages {
        analyzeDependencies(client, pkg.name, pkg.version)
    }
}
```

## Package Download Trend Analysis

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
    
    fmt.Println("📈 Package Download Trend Analysis")
    fmt.Println("==================================")
    
    for _, pkg := range packages {
        fmt.Printf("\n📦 Package: %s\n", pkg)
        fmt.Println("─────────────────────")
        
        var dailyDownloads, weeklyDownloads, monthlyDownloads int
        
        for _, period := range periods {
            stats, err := client.GetDownloadStats(ctx, pkg, period)
            if err != nil {
                log.Printf("Failed to get %s stats for %s: %v\n", period, pkg, err)
                continue
            }
            
            switch period {
            case "last-day":
                dailyDownloads = stats.Downloads
                fmt.Printf("📅 Last Day:   %10d downloads\n", stats.Downloads)
            case "last-week":
                weeklyDownloads = stats.Downloads
                fmt.Printf("📅 Last Week:  %10d downloads\n", stats.Downloads)
            case "last-month":
                monthlyDownloads = stats.Downloads
                fmt.Printf("📅 Last Month: %10d downloads\n", stats.Downloads)
            }
        }
        
        // Calculate trends
        if weeklyDownloads > 0 && dailyDownloads > 0 {
            avgDailyFromWeek := weeklyDownloads / 7
            trend := float64(dailyDownloads-avgDailyFromWeek) / float64(avgDailyFromWeek) * 100
            
            trendIndicator := "📊"
            if trend > 10 {
                trendIndicator = "📈"
            } else if trend < -10 {
                trendIndicator = "📉"
            }
            
            fmt.Printf("%s Trend:      %+.1f%% vs weekly average\n", trendIndicator, trend)
        }
    }
}
```

## Next Steps

- Explore more [API Documentation](/en/api/) features
- Check [Mirror Configuration](/en/examples/mirrors) for network optimization
- View [Basic Examples](/en/examples/basic) for fundamental usage patterns 