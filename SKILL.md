---
name: npm-crawler
description: High-performance NPM Registry client for querying package information, searching packages, retrieving download statistics, and downloading tarballs. Use when you need to: (1) Look up NPM package details (description, version, dependencies, maintainers), (2) Search NPM packages by keyword, (3) Get download counts or statistics for packages, (4) Download NPM package tarballs (.tgz files), (5) Check NPM registry status or health, (6) Work with specific NPM mirror sources (Taobao, NPMMirror, Huawei Cloud, Tencent, CNPM, Yarn). Trigger phrases: "npm package", "NPM registry", "search npm", "download npm tarball", "get npm stats", "npm mirror".
---

# NPM Registry Crawler

A Go-based NPM Registry client with CLI and SDK support. Query package information, search packages, get download stats, and download tarballs.

## Quick Start

Build the example tools:
```bash
cd examples/create_registry && go run main.go
cd examples/download_tarball && go run main.go
```

## Capabilities

- **Get Package Information**: Retrieve full package details (description, versions, dependencies, maintainers)
- **Search Packages**: Search NPM registry by keyword
- **Get Download Stats**: Retrieve download counts (last-day, last-week, last-month)
- **Download Tarballs**: Download .tgz packages to local filesystem
- **Multi-Mirror Support**: Use official npmjs.org or regional mirrors (Taobao, NPMMirror, Huawei Cloud, Tencent, CNPM, Yarn)

## Usage Patterns

### Go SDK Usage

```go
import "github.com/scagogogo/npm-crawler/pkg/registry"

client := registry.NewRegistry()
pkg, err := client.GetPackageInformation(ctx, "react")
```

### Available Methods

| Method | Description |
|--------|-------------|
| `GetRegistryInformation` | Get registry status (doc count, disk size, etc.) |
| `GetPackageInformation` | Get full package metadata |
| `SearchPackages(query, limit)` | Search packages by keyword |
| `GetPackageVersion(name, version)` | Get specific version info |
| `GetDownloadStats(name, period)` | Get download stats (last-day/week/month) |
| `DownloadTarball(name, version, dest)` | Download .tgz to local path |

### Mirror Sources

```go
registry.NewRegistry()                    // Official npmjs.org
registry.NewTaoBaoRegistry()              // registry.npm.taobao.org
registry.NewNpmMirrorRegistry()           // registry.npmmirror.com
registry.NewHuaWeiCloudRegistry()         // mirrors.huaweicloud.com
registry.NewTencentRegistry()             // mirrors.cloud.tencent.com
```

### Custom Registry with Proxy

```go
options := registry.NewOptions().
    SetRegistryURL("https://registry.npmjs.org").
    SetProxy("http://proxy:8080")
client := registry.NewRegistry(options)
```

## CLI Examples

### Get Registry Info
```bash
cd examples/create_registry && go run main.go
```

### Download Package Tarball
```bash
cd examples/download_tarball && go run main.go
```

## Output Format

All methods return JSON-serializable structs. Use `.ToJsonString()` for debugging:
```go
fmt.Println(pkg.ToJsonString())
```

## Detailed API Reference

See [references/api.md](references/api.md) for complete SDK documentation.
