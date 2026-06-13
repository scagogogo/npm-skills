# NPM Crawler

A powerful and flexible Go library for interacting with NPM registries. NPM Crawler provides comprehensive functionality to query package information, search packages, retrieve download statistics, and work with multiple registry mirrors.

## ✨ Features

- 📦 **Complete Package Information**: Get detailed package metadata, versions, dependencies
- 🔍 **Package Search**: Search packages with filtering and scoring
- 📊 **Download Statistics**: Retrieve download counts and trends
- 🌐 **Multiple Registry Support**: Official NPM, Taobao, NPM Mirror, Huawei Cloud
- 🚀 **High Performance**: Concurrent processing with connection pooling
- 🛡️ **Robust Error Handling**: Comprehensive error handling and retry mechanisms
- 🔧 **Flexible Configuration**: Custom HTTP clients, proxies, authentication
- 📈 **Registry Health Monitoring**: Check registry status and performance

## 🚀 Quick Start

### Installation

```bash
go get github.com/scagogogo/npm-skills
```

### Basic Usage

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
    
    fmt.Printf("Package: %s\n", pkg.Name)
    fmt.Printf("Latest Version: %s\n", pkg.DistTags["latest"])
    fmt.Printf("Description: %s\n", pkg.Description)
}
```

## 📚 Documentation

### Examples
- [**Basic Usage**](examples/basic.md) - Get started with fundamental operations
- [**Advanced Usage**](examples/advanced.md) - Complex scenarios and patterns  
- [**Mirror Configuration**](examples/mirrors.md) - Optimize network access with mirrors

### API Reference
- [**Registry Client**](api/registry.md) - Main client interface and methods
- [**Data Models**](api/models.md) - Package and registry data structures
- [**Configuration**](api/configuration.md) - Client options and settings

## 🌍 Supported Registries

| Registry | URL | Region | Notes |
|----------|-----|--------|-------|
| **NPM Official** | `https://registry.npmjs.org` | Global | Official NPM registry |
| **Taobao Mirror** | `https://registry.npmmirror.com` | China | Fast mirror for China users |
| **NPM Mirror** | `https://skimdb.npmjs.com/registry` | Global | Alternative global mirror |
| **Huawei Cloud** | `https://mirrors.huaweicloud.com/repository/npm` | China | Enterprise-grade mirror |

### Quick Registry Usage

```go
// Official NPM
client := registry.NewRegistry()

// Taobao Mirror (recommended for China)
client := registry.NewTaoBaoRegistry()

// NPM Mirror
client := registry.NewNpmMirrorRegistry()

// Huawei Cloud
client := registry.NewHuaWeiCloudRegistry()
```

## 🛠️ Core Functionality

### Package Information
```go
pkg, err := client.GetPackageInformation(ctx, "lodash")
```

### Specific Version Details
```go
version, err := client.GetPackageVersion(ctx, "react", "18.2.0")
```

### Package Search
```go
results, err := client.SearchPackages(ctx, "react ui component", 10)
```

### Download Statistics
```go
stats, err := client.GetDownloadStats(ctx, "vue", "last-month")
```

### Registry Information
```go
info, err := client.GetRegistryInformation(ctx)
```

## ⚙️ Configuration Options

### Custom Registry URL
```go
options := registry.NewOptions().
    SetRegistryURL("https://custom-registry.com")
client := registry.NewRegistry(options)
```

### Proxy Configuration
```go
options := registry.NewOptions().
    SetProxy("http://proxy.example.com:8080")
client := registry.NewRegistry(options)
```

### Custom HTTP Client
```go
httpClient := &http.Client{Timeout: 30 * time.Second}
options := registry.NewOptions().SetHTTPClient(httpClient)
client := registry.NewRegistry(options)
```

## 🔧 Advanced Features

### Concurrent Package Processing
```go
var wg sync.WaitGroup
packages := []string{"react", "vue", "angular"}

for _, pkg := range packages {
    wg.Add(1)
    go func(packageName string) {
        defer wg.Done()
        info, err := client.GetPackageInformation(ctx, packageName)
        // Process package info...
    }(pkg)
}
wg.Wait()
```

### Mirror Performance Testing
```go
mirrors := map[string]*registry.Registry{
    "Official": registry.NewRegistry(),
    "Taobao":   registry.NewTaoBaoRegistry(),
}

for name, client := range mirrors {
    start := time.Now()
    _, err := client.GetPackageInformation(ctx, "test-package")
    duration := time.Since(start)
    fmt.Printf("%s: %v\n", name, duration)
}
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](../CONTRIBUTING.md) for details.

### Development Setup

1. Clone the repository:
```bash
git clone https://github.com/scagogogo/npm-skills.git
cd npm-skills
```

2. Install dependencies:
```bash
go mod download
```

3. Run tests:
```bash
go test ./...
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.

## 🔗 Links

- [GitHub Repository](https://github.com/scagogogo/npm-skills)
- [Go Package Documentation](https://pkg.go.dev/github.com/scagogogo/npm-skills)
- [Issue Tracker](https://github.com/scagogogo/npm-skills/issues)
- [中文文档](../README.md) 