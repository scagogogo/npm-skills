# Download Tarball Examples

This page demonstrates how to use the `DownloadTarball` method to download NPM package tarballs to local files.

## Example 1: Download Specific Version

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

func main() {
    client := registry.NewRegistry()
    ctx := context.Background()

    packageName := "axios"
    version := "1.0.0"
    destPath := "./axios-1.0.0.tgz"

    // Remove existing file if present
    os.Remove(destPath)

    fmt.Printf("Downloading %s@%s...\n", packageName, version)

    err := client.DownloadTarball(ctx, packageName, version, destPath)
    if err != nil {
        log.Fatalf("Download failed: %v", err)
    }

    // Verify the file
    info, err := os.Stat(destPath)
    if err != nil {
        log.Fatalf("File stat failed: %v", err)
    }

    fmt.Printf("Download successful!\n")
    fmt.Printf("File size: %d bytes\n", info.Size())
}
```

## Example 2: Using CNPM Mirror (Recommended for China Region)

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

func main() {
    // Use CNPM mirror for faster downloads
    options := registry.NewOptions().SetRegistryURL(registry.RegistryUrlCnpm)
    client := registry.NewRegistry(options)

    ctx := context.Background()
    packageName := "react"
    version := "18.2.0"
    destDir := "/tmp/npm-tarballs"

    // Ensure destination directory exists
    if err := os.MkdirAll(destDir, 0755); err != nil {
        log.Fatalf("Failed to create directory: %v", err)
    }

    destPath := fmt.Sprintf("%s/%s-%s.tgz", destDir, packageName, version)
    os.Remove(destPath)

    fmt.Printf("Downloading %s@%s using CNPM mirror...\n", packageName, version)

    err := client.DownloadTarball(ctx, packageName, version, destPath)
    if err != nil {
        log.Fatalf("Download failed: %v", err)
    }

    info, _ := os.Stat(destPath)
    fmt.Printf("Download successful! File size: %d bytes\n", info.Size())
}
```

## Example 3: Download Latest Version

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

func main() {
    client := registry.NewRegistry()
    ctx := context.Background()

    packageName := "vue"
    version := "latest"
    destPath := "./vue-latest.tgz"

    // Get version info first
    versionInfo, err := client.GetPackageVersion(ctx, packageName, version)
    if err != nil {
        log.Fatalf("Failed to get version info: %v", err)
    }

    fmt.Printf("Latest version: %s\n", versionInfo.Version)

    os.Remove(destPath)
    err = client.DownloadTarball(ctx, packageName, version, destPath)
    if err != nil {
        log.Fatalf("Download failed: %v", err)
    }

    info, _ := os.Stat(destPath)
    fmt.Printf("Download successful! File size: %d bytes\n", info.Size())
}
```

## Example 4: Batch Download Multiple Packages

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "sync"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

func main() {
    client := registry.NewRegistry()
    ctx := context.Background()

    packages := []struct {
        name    string
        version string
    }{
        {"react", "18.2.0"},
        {"vue", "3.4.0"},
        {"angular", "17.3.0"},
    }

    destDir := "/tmp/npm-packages"
    os.MkdirAll(destDir, 0755)

    var wg sync.WaitGroup

    for _, pkg := range packages {
        wg.Add(1)
        go func(name, version string) {
            defer wg.Done()

            destPath := filepath.Join(destDir, fmt.Sprintf("%s-%s.tgz", name, version))
            os.Remove(destPath)

            err := client.DownloadTarball(ctx, name, version, destPath)
            if err != nil {
                log.Printf("Failed to download %s@%s: %v\n", name, version, err)
                return
            }

            info, _ := os.Stat(destPath)
            fmt.Printf("✅ %s@%s: %d bytes\n", name, version, info.Size())
        }(pkg.name, pkg.version)
    }

    wg.Wait()
    fmt.Println("All downloads complete!")
}
```

## Run the Examples

```bash
go mod init tarball-example
go get github.com/scagogogo/npm-skills
go run main.go
```

## Next Steps

- Check [Basic Usage Examples](/en/examples/basic) for more fundamentals
- Read [API Documentation](/en/api/registry) for complete API reference
- Explore [Mirror Configuration](/en/examples/mirrors) to choose the fastest download source
