# Mirror Configuration

This guide explains how to configure and use different NPM registry mirrors to optimize network access and improve download speeds.

## Available Mirrors

NPM Skills supports multiple registry mirrors out of the box:

### Official NPM Registry
```go
client := registry.NewRegistry()
// Uses: https://registry.npmjs.org
```

### Taobao Mirror (China)
```go
client := registry.NewTaoBaoRegistry()
// Uses: https://registry.npmmirror.com
```

### NPM Mirror
```go
client := registry.NewNpmMirrorRegistry()
// Uses: https://skimdb.npmjs.com/registry
```

### Huawei Cloud Mirror (China)
```go
client := registry.NewHuaWeiCloudRegistry()
// Uses: https://mirrors.huaweicloud.com/repository/npm
```

## Custom Registry Configuration

### Basic Custom Registry
```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

func main() {
    // Configure custom registry URL
    options := registry.NewOptions().
        SetRegistryURL("https://your-custom-registry.com")
    
    client := registry.NewRegistry(options)
    ctx := context.Background()
    
    pkg, err := client.GetPackageInformation(ctx, "lodash")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Package from custom registry: %s\n", pkg.Name)
}
```

### Registry with Authentication
```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

func main() {
    // Create HTTP client with authentication
    httpClient := &http.Client{
        Transport: &authTransport{
            token: "your-auth-token",
            base:  http.DefaultTransport,
        },
    }
    
    options := registry.NewOptions().
        SetRegistryURL("https://private-registry.com").
        SetHTTPClient(httpClient)
    
    client := registry.NewRegistry(options)
    ctx := context.Background()
    
    pkg, err := client.GetPackageInformation(ctx, "private-package")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Private package: %s\n", pkg.Name)
}

type authTransport struct {
    token string
    base  http.RoundTripper
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    req.Header.Set("Authorization", "Bearer "+t.token)
    return t.base.RoundTrip(req)
}
```

## Mirror Performance Testing

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

func testMirrorPerformance(name string, client *registry.Registry) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    start := time.Now()
    
    // Test package information retrieval
    pkg, err := client.GetPackageInformation(ctx, "react")
    if err != nil {
        fmt.Printf("❌ %s: Failed - %v\n", name, err)
        return
    }
    
    duration := time.Since(start)
    
    fmt.Printf("✅ %s: %v (Package: %s v%s)\n", 
        name, duration, pkg.Name, pkg.DistTags["latest"])
}

func main() {
    mirrors := map[string]*registry.Registry{
        "NPM Official":   registry.NewRegistry(),
        "Taobao Mirror":  registry.NewTaoBaoRegistry(),
        "NPM Mirror":     registry.NewNpmMirrorRegistry(),
        "Huawei Cloud":   registry.NewHuaWeiCloudRegistry(),
    }
    
    fmt.Println("🚀 Mirror Performance Test")
    fmt.Println("==========================")
    
    for name, client := range mirrors {
        testMirrorPerformance(name, client)
    }
}
```

## Automatic Mirror Selection

```go
package main

import (
    "context"
    "fmt"
    "sync"
    "time"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

type MirrorResult struct {
    Name     string
    Client   *registry.Registry
    Latency  time.Duration
    Success  bool
}

func testMirror(name string, client *registry.Registry, results chan<- MirrorResult) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    start := time.Now()
    _, err := client.GetPackageInformation(ctx, "lodash")
    latency := time.Since(start)
    
    results <- MirrorResult{
        Name:    name,
        Client:  client,
        Latency: latency,
        Success: err == nil,
    }
}

func selectBestMirror() *registry.Registry {
    mirrors := map[string]*registry.Registry{
        "NPM Official":   registry.NewRegistry(),
        "Taobao Mirror":  registry.NewTaoBaoRegistry(),
        "NPM Mirror":     registry.NewNpmMirrorRegistry(),
        "Huawei Cloud":   registry.NewHuaWeiCloudRegistry(),
    }
    
    results := make(chan MirrorResult, len(mirrors))
    var wg sync.WaitGroup
    
    // Test all mirrors concurrently
    for name, client := range mirrors {
        wg.Add(1)
        go func(n string, c *registry.Registry) {
            defer wg.Done()
            testMirror(n, c, results)
        }(name, client)
    }
    
    // Wait for all tests to complete
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // Find the fastest successful mirror
    var bestMirror *registry.Registry
    var bestLatency time.Duration = time.Hour // Start with a very high value
    
    fmt.Println("🔍 Testing mirror performance...")
    
    for result := range results {
        status := "❌"
        if result.Success {
            status = "✅"
            if result.Latency < bestLatency {
                bestLatency = result.Latency
                bestMirror = result.Client
            }
        }
        
        fmt.Printf("%s %s: %v\n", status, result.Name, result.Latency)
    }
    
    if bestMirror != nil {
        fmt.Printf("🏆 Selected fastest mirror with %v latency\n", bestLatency)
    } else {
        fmt.Println("⚠️  No mirrors available, using default")
        bestMirror = registry.NewRegistry()
    }
    
    return bestMirror
}

func main() {
    client := selectBestMirror()
    
    // Use the selected mirror
    ctx := context.Background()
    pkg, err := client.GetPackageInformation(ctx, "react")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("\n📦 Successfully retrieved: %s v%s\n", 
        pkg.Name, pkg.DistTags["latest"])
}
```

## Regional Mirror Recommendations

### For China Users
```go
// Recommended order for China users
mirrors := []*registry.Registry{
    registry.NewTaoBaoRegistry(),      // Primary: Taobao mirror
    registry.NewHuaWeiCloudRegistry(), // Fallback: Huawei Cloud
    registry.NewRegistry(),            // Last resort: Official
}
```

### For Global Users
```go
// Recommended order for global users
mirrors := []*registry.Registry{
    registry.NewRegistry(),         // Primary: Official NPM
    registry.NewNpmMirrorRegistry(), // Fallback: NPM mirror
}
```

## Mirror Health Monitoring

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

func monitorMirror(name string, client *registry.Registry) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    fmt.Printf("🔍 Monitoring %s...\n", name)
    
    // Test registry information
    start := time.Now()
    info, err := client.GetRegistryInformation(ctx)
    if err != nil {
        fmt.Printf("❌ %s: Registry info failed - %v\n", name, err)
        return
    }
    
    infoLatency := time.Since(start)
    
    // Test package query
    start = time.Now()
    _, err = client.GetPackageInformation(ctx, "lodash")
    if err != nil {
        fmt.Printf("❌ %s: Package query failed - %v\n", name, err)
        return
    }
    
    queryLatency := time.Since(start)
    
    // Test search functionality
    start = time.Now()
    _, err = client.SearchPackages(ctx, "react", 1)
    searchLatency := time.Since(start)
    
    fmt.Printf("✅ %s Health Report:\n", name)
    fmt.Printf("   📊 Total packages: %d\n", info.DocCount)
    fmt.Printf("   ⏱️  Registry info: %v\n", infoLatency)
    fmt.Printf("   ⏱️  Package query: %v\n", queryLatency)
    if err == nil {
        fmt.Printf("   ⏱️  Search query: %v\n", searchLatency)
    } else {
        fmt.Printf("   ❌ Search query: Failed\n")
    }
    fmt.Println()
}

func main() {
    mirrors := map[string]*registry.Registry{
        "NPM Official":   registry.NewRegistry(),
        "Taobao Mirror":  registry.NewTaoBaoRegistry(),
        "NPM Mirror":     registry.NewNpmMirrorRegistry(),
        "Huawei Cloud":   registry.NewHuaWeiCloudRegistry(),
    }
    
    fmt.Println("🏥 Mirror Health Monitoring")
    fmt.Println("===========================")
    
    for name, client := range mirrors {
        monitorMirror(name, client)
    }
}
```

## Best Practices

### 1. Mirror Selection Strategy
- Test multiple mirrors at startup
- Select based on latency and reliability
- Implement fallback mechanisms

### 2. Caching and Performance
- Use connection pooling for HTTP clients
- Implement response caching where appropriate
- Monitor and log mirror performance

### 3. Error Handling
- Implement retry logic with exponential backoff
- Have multiple mirrors as fallbacks
- Log mirror failures for monitoring

### 4. Security Considerations
- Use HTTPS for all registry connections
- Validate SSL certificates
- Implement proper authentication for private registries

## Next Steps

- Check [Basic Examples](/en/examples/basic) for fundamental usage
- Explore [Advanced Examples](/en/examples/advanced) for complex scenarios
- Review [API Documentation](/en/api/) for complete feature reference 