# API 概述

NPM Skills 提供了一套简洁而强大的 API，用于与 NPM Registry 进行交互。本文档涵盖了所有可用的 API 方法、参数和返回值。

## 核心 API

### Registry 客户端

Registry 是 NPM Skills 的核心客户端，提供了访问 NPM Registry 的所有功能。

```go
import "github.com/scagogogo/npm-skills/pkg/registry"

// 创建默认客户端
client := registry.NewRegistry()

// 使用自定义配置
options := registry.NewOptions().
    SetRegistryURL("https://registry.npmjs.org").
    SetProxy("http://proxy.example.com:8080")
client := registry.NewRegistry(options)
```

## 主要方法

### 包信息查询

#### GetPackageInformation
获取指定 NPM 包的完整信息。

```go
func (r *Registry) GetPackageInformation(ctx context.Context, packageName string) (*models.Package, error)
```

**参数：**
- `ctx` - 上下文，用于取消和超时控制
- `packageName` - 要查询的包名称

**返回值：**
- `*models.Package` - 包的完整信息
- `error` - 错误信息

**示例：**
```go
pkg, err := client.GetPackageInformation(ctx, "react")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("包名: %s, 最新版本: %s\n", pkg.Name, pkg.DistTags["latest"])
```

#### GetPackageVersion
获取指定包的特定版本信息。

```go
func (r *Registry) GetPackageVersion(ctx context.Context, packageName, version string) (*models.Version, error)
```

**参数：**
- `ctx` - 上下文
- `packageName` - 包名称
- `version` - 版本号（如 "1.0.0" 或 "latest"）

**示例：**
```go
version, err := client.GetPackageVersion(ctx, "react", "18.2.0")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("版本: %s, 依赖: %+v\n", version.Version, version.Dependencies)
```

### 包搜索

#### SearchPackages
搜索 NPM 包。

```go
func (r *Registry) SearchPackages(ctx context.Context, query string, limit int) (*models.SearchResult, error)
```

**参数：**
- `ctx` - 上下文
- `query` - 搜索关键字
- `limit` - 返回结果数量限制（默认 20）

**示例：**
```go
result, err := client.SearchPackages(ctx, "react ui", 10)
if err != nil {
    log.Fatal(err)
}

for _, obj := range result.Objects {
    fmt.Printf("- %s: %s\n", obj.Package.Name, obj.Package.Description)
}
```

### 统计信息

#### GetRegistryInformation
获取 NPM Registry 的状态信息。

```go
func (r *Registry) GetRegistryInformation(ctx context.Context) (*models.RegistryInformation, error)
```

**示例：**
```go
info, err := client.GetRegistryInformation(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Registry: %s, 包总数: %d\n", info.DbName, info.DocCount)
```

#### GetDownloadStats
获取包的下载统计信息。

```go
func (r *Registry) GetDownloadStats(ctx context.Context, packageName, period string) (*models.DownloadStats, error)
```

**参数：**
- `packageName` - 包名称
- `period` - 统计周期（"last-day", "last-week", "last-month"）

**示例：**
```go
stats, err := client.GetDownloadStats(ctx, "react", "last-week")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("下载次数: %d\n", stats.Downloads)
```

#### DownloadTarball
下载 NPM 包的 tarball 文件到本地。

```go
func (r *Registry) DownloadTarball(ctx context.Context, packageName, version, destPath string) error
```

**参数：**
- `packageName` - 包名称
- `version` - 版本号或标签（如 "latest"）
- `destPath` - 本地保存路径

**示例：**
```go
err := client.DownloadTarball(ctx, "react", "18.2.0", "./react.tgz")
```

## 配置选项

### Options
配置 Registry 客户端的选项。

```go
type Options struct {
    RegistryURL string  // NPM 仓库 URL
    Proxy       string  // HTTP 代理 URL
}
```

**方法：**
- `NewOptions()` - 创建默认选项
- `SetRegistryURL(url string)` - 设置 Registry URL
- `SetProxy(proxyUrl string)` - 设置代理
- `GetHttpClient()` - 获取配置的 HTTP 客户端

## 镜像源支持

NPM Skills 内置支持多种镜像源：

| 函数 | 镜像源 | URL |
|------|--------|-----|
| `NewRegistry()` | 官方 NPM | https://registry.npmjs.org |
| `NewTaoBaoRegistry()` | 淘宝镜像 | https://registry.npm.taobao.org |
| `NewNpmMirrorRegistry()` | NPM Mirror | https://registry.npmmirror.com |
| `NewHuaWeiCloudRegistry()` | 华为云 | https://mirrors.huaweicloud.com/repository/npm |
| `NewTencentRegistry()` | 腾讯云 | http://mirrors.cloud.tencent.com/npm |
| `NewCnpmRegistry()` | CNPM | http://r.cnpmjs.org |
| `NewYarnRegistry()` | Yarn | https://registry.yarnpkg.com |

## 数据模型

### Package
表示 NPM 包的完整信息。

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
    // ... 其他字段
}
```

### Version
表示包的特定版本信息。

```go
type Version struct {
    Name            string               `json:"name"`
    Version         string               `json:"version"`
    Description     string               `json:"description"`
    Dependencies    map[string]string    `json:"dependencies"`
    DevDependencies map[string]string    `json:"devDependencies"`
    Dist            *Dist                `json:"dist"`
    // ... 其他字段
}
```

### SearchResult
表示搜索结果。

```go
type SearchResult struct {
    Objects []SearchObject `json:"objects"`
    Total   int            `json:"total"`
    Time    string         `json:"time"`
}
```

### DownloadStats
表示下载统计信息。

```go
type DownloadStats struct {
    Downloads int    `json:"downloads"`
    Start     string `json:"start"`
    End       string `json:"end"`
    Package   string `json:"package"`
}
```

## 错误处理

所有 API 方法都返回 error 类型，建议在生产环境中进行适当的错误处理：

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
        // 处理超时
    case errors.Is(err, context.Canceled):
        // 处理取消
    default:
        // 处理其他错误
    }
}
```

## 最佳实践

1. **使用上下文**: 始终传递适当的上下文以支持取消和超时
2. **选择合适的镜像**: 根据网络环境选择最快的镜像源
3. **并发安全**: Registry 客户端是并发安全的，可以在多个 goroutine 中使用
4. **错误处理**: 实现适当的错误处理和重试机制
5. **代理配置**: 在企业环境中正确配置代理

## 下一步

- 查看 [中文详细API文档](/api/registry) 或 [英文API文档](/en/api/registry)
- 浏览 [示例代码](/examples/basic) 学习具体用法
- 访问 [GitHub](https://github.com/scagogogo/npm-skills) 查看完整源代码
