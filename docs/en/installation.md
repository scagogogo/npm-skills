# Installation Guide

This guide provides detailed instructions on how to install and configure NPM Crawler in different environments.

## System Requirements

- **Go Version**: 1.20 or higher
- **Operating System**: Linux, macOS, Windows
- **Network**: Access to NPM registry or its mirror sources

## Go Environment Setup

### Install Go

If you haven't installed Go yet, please visit the [Go official website](https://golang.org/dl/) to download and install the latest version.

Verify installation:
```bash
go version
```

You should see output similar to `go version go1.21.0 linux/amd64`.

### Configure Go Modules

Ensure your project has Go modules enabled:

```bash
# In your project directory
go mod init your-project-name
```

## Install NPM Crawler

### Using go get (Recommended)

```bash
go get github.com/scagogogo/npm-skills
```

### Add to go.mod

You can also add the dependency directly in your `go.mod` file:

```go
module your-project

go 1.20

require (
    github.com/scagogogo/npm-skills latest
)
```

Then run:
```bash
go mod tidy
```

### Specify Version

If you need a specific version:

```bash
# Install specific version
go get github.com/scagogogo/npm-skills@v1.0.0

# Install latest pre-release version
go get github.com/scagogogo/npm-skills@latest
```

## Verify Installation

Create a simple test file to verify the installation:

```go
// test.go
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
    
    // Get simple package information
    pkg, err := client.GetPackageInformation(ctx, "lodash")
    if err != nil {
        log.Fatal("Installation verification failed:", err)
    }
    
    fmt.Printf("✅ NPM Crawler installed successfully!\n")
    fmt.Printf("Test package: %s\n", pkg.Name)
    fmt.Printf("Latest version: %s\n", pkg.DistTags["latest"])
}
```

Run the test:
```bash
go run test.go
```

If you see a success message, the installation is correct.

## Network Configuration

### Proxy Settings

If you're in a corporate network environment:

```go
options := registry.NewOptions().
    SetProxy("http://proxy.corp.com:8080")
client := registry.NewRegistry(options)
```

### Mirror Source Selection

Choose the appropriate mirror source based on your geographic location:

```go
// Recommended for users in mainland China
client := registry.NewNpmMirrorRegistry()  // NPM Mirror
// or
client := registry.NewTaoBaoRegistry()     // Taobao Mirror

// For global users
client := registry.NewRegistry()           // Official source
```

## Common Issues

### 1. Network Connection Issues

**Error**: `dial tcp: lookup registry.npmjs.org: no such host`

**Solution**:
- Check network connection
- Try using domestic mirror sources
- Configure proper proxy

```go
// Use domestic mirror
client := registry.NewNpmMirrorRegistry()
```

### 2. Proxy Authentication Issues

**Error**: `407 Proxy Authentication Required`

**Solution**:
```go
// Use proxy with authentication
options := registry.NewOptions().
    SetProxy("http://username:password@proxy.corp.com:8080")
```

### 3. Go Version Incompatibility

**Error**: `package github.com/scagogogo/npm-skills: module requires Go 1.20`

**Solution**:
- Upgrade Go to version 1.20 or higher
- Or use a compatible earlier version

### 4. Module Import Issues

**Error**: `cannot find module providing package github.com/scagogogo/npm-skills`

**Solution**:
```bash
# Clean module cache
go clean -modcache

# Re-download dependencies
go mod download

# Tidy dependencies
go mod tidy
```

## Docker Environment

If you're using Docker:

```dockerfile
FROM golang:1.21-alpine

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN go build -o main .

CMD ["./main"]
```

## CI/CD Integration

### GitHub Actions

```yaml
name: Test with NPM Crawler

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.21
    
    - name: Download dependencies
      run: go mod download
    
    - name: Run tests
      run: go test ./...
```

### GitLab CI

```yaml
stages:
  - test

test:
  image: golang:1.21
  stage: test
  script:
    - go mod download
    - go test ./...
```

## Performance Optimization

### Connection Reuse

NPM Crawler uses Go's standard HTTP client, which automatically supports connection reuse.

### Concurrency Control

```go
import "golang.org/x/sync/semaphore"

// Limit concurrency
sem := semaphore.NewWeighted(10)

for _, pkg := range packages {
    sem.Acquire(ctx, 1)
    go func(packageName string) {
        defer sem.Release(1)
        // Process package
    }(pkg)
}
```

## Development Environment Setup

### IDE Configuration

**VS Code**:
Install the Go extension and configure:

```json
{
    "go.useLanguageServer": true,
    "go.alternateTools": {
        "go": "/usr/local/go/bin/go"
    }
}
```

**GoLand**:
Ensure proper GOROOT and GOPATH configuration.

### Debug Configuration

```go
import "log"

// Enable verbose logging
log.SetFlags(log.LstdFlags | log.Lshortfile)

// Add logging at key points
pkg, err := client.GetPackageInformation(ctx, "react")
log.Printf("Get package info: %s, error: %v", "react", err)
```

## Next Steps

After installation, you can:

1. Check the [Getting Started Guide](/en/getting-started) to learn basic usage
2. Read the [API Documentation](/en/api/) to understand all features
3. Browse [example code](/en/examples/basic) to learn practical applications
4. Participate in the [GitHub project](https://github.com/scagogogo/npm-skills) to contribute code

If you encounter issues, please report them in [GitHub Issues](https://github.com/scagogogo/npm-skills/issues). 