# Registry 客户端

Registry 客户端是 NPM Crawler 的核心组件，提供与 NPM 注册表交互的所有功能。

## 创建客户端

### `NewRegistry(options ...*Options) *Registry`

创建一个新的 Registry 客户端实例。

**参数:**
- `options` - 可选的配置选项

**示例:**
```go
// 使用默认配置
client := registry.NewRegistry()

// 使用自定义配置
options := registry.NewOptions().
    SetRegistryURL("https://registry.npmjs.org").
    SetProxy("http://proxy.example.com:8080")
client := registry.NewRegistry(options)
```

## 核心方法

### `GetPackageInformation(ctx context.Context, packageName string) (*models.Package, error)`

获取指定 NPM 包的详细信息。

**参数:**
- `ctx` - 上下文，用于取消和超时控制
- `packageName` - 要查询的包名称

**返回值:**
- `*models.Package` - 完整的包信息
- `error` - 如果请求失败则返回错误

**示例:**
```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

pkg, err := client.GetPackageInformation(ctx, "lodash")
if err != nil {
    return fmt.Errorf("获取包信息失败: %w", err)
}

// 访问包数据
fmt.Printf("名称: %s\n", pkg.Name)
fmt.Printf("最新版本: %s\n", pkg.DistTags["latest"])
fmt.Printf("作者: %s\n", pkg.Author.Name)

// 访问特定版本
if version, exists := pkg.Versions["4.17.21"]; exists {
    fmt.Printf("版本 4.17.21 依赖: %+v\n", version.Dependencies)
}
```

### `GetRegistryInformation(ctx context.Context) (*models.RegistryInformation, error)`

获取 NPM 注册表的状态和元数据信息。

**参数:**
- `ctx` - 上下文，用于取消和超时控制

**返回值:**
- `*models.RegistryInformation` - 注册表状态信息
- `error` - 如果请求失败则返回错误

**示例:**
```go
info, err := client.GetRegistryInformation(ctx)
if err != nil {
    return fmt.Errorf("获取注册表信息失败: %w", err)
}

fmt.Printf("注册表: %s\n", info.DbName)
fmt.Printf("包总数: %d\n", info.DocCount)
fmt.Printf("数据库大小: %d 字节\n", info.DataSize)
fmt.Printf("磁盘使用: %d 字节\n", info.DiskSize)
```

### `SearchPackages(ctx context.Context, query string, limit int) (*models.SearchResult, error)`

搜索 NPM 包。

**参数:**
- `ctx` - 上下文
- `query` - 搜索关键字
- `limit` - 返回结果数量限制，默认为 20

**返回值:**
- `*models.SearchResult` - 搜索结果
- `error` - 如果请求失败则返回错误

**示例:**
```go
result, err := client.SearchPackages(ctx, "react", 10)
if err != nil {
    return fmt.Errorf("搜索失败: %w", err)
}

fmt.Printf("找到 %d 个结果\n", result.Total)
for _, obj := range result.Objects {
    pkg := obj.Package
    fmt.Printf("包名: %s\n", pkg.Name)
    fmt.Printf("版本: %s\n", pkg.Version)
    fmt.Printf("描述: %s\n", pkg.Description)
    fmt.Printf("评分: %.2f\n", obj.Score.Final)
    fmt.Println("---")
}
```

### `GetPackageVersion(ctx context.Context, packageName, version string) (*models.Version, error)`

获取指定包的特定版本信息。

**参数:**
- `ctx` - 上下文
- `packageName` - 包名称
- `version` - 版本号或标签（如 "1.0.0" 或 "latest"）

**返回值:**
- `*models.Version` - 版本详细信息
- `error` - 如果请求失败则返回错误

**示例:**
```go
version, err := client.GetPackageVersion(ctx, "react", "18.2.0")
if err != nil {
    return fmt.Errorf("获取版本信息失败: %w", err)
}

fmt.Printf("版本: %s\n", version.Version)
fmt.Printf("描述: %s\n", version.Description)
fmt.Printf("依赖: %+v\n", version.Dependencies)
fmt.Printf("开发依赖: %+v\n", version.DevDependencies)
```

### `GetDownloadStats(ctx context.Context, packageName, period string) (*models.DownloadStats, error)`

获取指定包的下载统计信息。

**参数:**
- `ctx` - 上下文
- `packageName` - 包名称
- `period` - 统计周期（"last-day", "last-week", "last-month"）

**返回值:**
- `*models.DownloadStats` - 下载统计信息
- `error` - 如果请求失败则返回错误

**示例:**
```go
stats, err := client.GetDownloadStats(ctx, "react", "last-week")
if err != nil {
    return fmt.Errorf("获取下载统计失败: %w", err)
}

fmt.Printf("包: %s\n", stats.Package)
fmt.Printf("下载次数: %d\n", stats.Downloads)
fmt.Printf("统计周期: %s 到 %s\n", stats.Start, stats.End)
```

### `GetOptions() *Options`

返回当前注册表客户端的配置选项。

**返回值:**
- `*Options` - 当前配置选项

## 镜像源

NPM Crawler 内置支持多种镜像源，特别适合中国大陆用户：

### 官方镜像

```go
// NPM 官方注册表 (全球)
client := registry.NewRegistry()

// Yarn 官方镜像 (全球)
client := registry.NewYarnRegistry()
```

### 中国镜像源

```go
// 淘宝 NPM 镜像 (中国)
client := registry.NewTaoBaoRegistry()

// NPM Mirror (中国)
client := registry.NewNpmMirrorRegistry()

// 华为云镜像 (中国)
client := registry.NewHuaWeiCloudRegistry()

// 腾讯云镜像 (中国)
client := registry.NewTencentRegistry()

// CNPM 镜像 (中国)
client := registry.NewCnpmRegistry()

// NPM CouchDB 镜像
client := registry.NewNpmjsComRegistry()
```

## 最佳实践

### 超时控制

```go
// 设置合理的超时时间
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

pkg, err := client.GetPackageInformation(ctx, "package-name")
```

### 错误处理

```go
pkg, err := client.GetPackageInformation(ctx, "package-name")
if err != nil {
    // 检查是否是网络错误
    if strings.Contains(err.Error(), "timeout") {
        log.Printf("请求超时，请检查网络连接")
    }
    // 检查是否是包不存在
    if strings.Contains(err.Error(), "404") {
        log.Printf("包不存在: %s", packageName)
    }
    return fmt.Errorf("获取包信息失败: %w", err)
}
```

### 并发访问

```go
// 使用 goroutine 并发获取多个包信息
packages := []string{"react", "vue", "angular"}
results := make(chan *models.Package, len(packages))

for _, pkg := range packages {
    go func(packageName string) {
        info, err := client.GetPackageInformation(ctx, packageName)
        if err != nil {
            log.Printf("获取 %s 失败: %v", packageName, err)
            results <- nil
            return
        }
        results <- info
    }(pkg)
}

// 收集结果
for i := 0; i < len(packages); i++ {
    result := <-results
    if result != nil {
        fmt.Printf("包: %s, 版本: %s\n", result.Name, result.DistTags["latest"])
    }
}
``` 