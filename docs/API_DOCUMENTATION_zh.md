# NPM Crawler API 文档

<div align="center">

[![Go Reference](https://pkg.go.dev/badge/github.com/scagogogo/npm-skills.svg)](https://pkg.go.dev/github.com/scagogogo/npm-skills)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

_完整的 API 参考和使用指南_

[中文版本](API_DOCUMENTATION_zh.md) | [English Version](API_DOCUMENTATION.md)

</div>

## 目录

- [概述](#概述)
- [安装](#安装)
- [快速开始](#快速开始)
- [核心 API 参考](#核心-api-参考)
  - [Registry 客户端](#registry-客户端)
  - [配置选项](#配置选项)
  - [镜像源](#镜像源)
- [数据模型](#数据模型)
- [高级用法](#高级用法)
- [错误处理](#错误处理)
- [最佳实践](#最佳实践)
- [故障排除](#故障排除)
- [示例](#示例)

## 概述

NPM Crawler 提供了一个全面的 Go 客户端用于与 NPM 注册表交互。本文档涵盖了所有可用的 API、使用模式和将库集成到应用程序中的最佳实践。

### 主要特性

- **多注册表支持**: 官方 NPM、淘宝、华为云等多种镜像源
- **代理配置**: 支持企业环境的 HTTP 代理
- **上下文支持**: 取消和超时处理
- **类型安全**: 完整的 NPM 元数据 Go 类型定义
- **高性能**: 针对并发访问优化

## 安装

```bash
go get github.com/scagogogo/npm-skills
```

**要求:**
- Go 1.20 或更高版本
- 网络连接（或访问 NPM 注册表镜像）

## 快速开始

### 基本包信息获取

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
    // 使用默认设置创建客户端（官方 NPM 注册表）
    client := registry.NewRegistry()
    
    // 创建带超时的上下文
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // 获取包信息
    pkg, err := client.GetPackageInformation(ctx, "react")
    if err != nil {
        log.Fatalf("获取包信息失败: %v", err)
    }
    
    fmt.Printf("包名: %s\n", pkg.Name)
    fmt.Printf("描述: %s\n", pkg.Description)
    fmt.Printf("最新版本: %s\n", pkg.DistTags["latest"])
    fmt.Printf("许可证: %s\n", pkg.License)
}
```

### 使用备用镜像

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

func main() {
    // 使用淘宝镜像在中国地区更快访问
    client := registry.NewTaoBaoRegistry()
    
    ctx := context.Background()
    
    pkg, err := client.GetPackageInformation(ctx, "vue")
    if err != nil {
        log.Fatalf("错误: %v", err)
    }
    
    fmt.Printf("包名: %s（来自淘宝镜像）\n", pkg.Name)
}
```

## 核心 API 参考

### Registry 客户端

#### 创建 Registry 客户端

##### `NewRegistry(options ...*Options) *Registry`

创建带有可选配置的新 Registry 客户端。

**参数:**
- `options` (可选): 客户端的配置选项

**返回值:**
- `*Registry`: 配置好的注册表客户端

**示例:**
```go
// 默认配置（官方 NPM 注册表）
client := registry.NewRegistry()

// 自定义配置
options := registry.NewOptions().
    SetRegistryURL("https://registry.npmjs.org").
    SetProxy("http://proxy.company.com:8080")
client := registry.NewRegistry(options)
```

##### 镜像特定的构造函数

```go
// 淘宝 NPM 镜像（中国）
client := registry.NewTaoBaoRegistry()

// NPM Mirror（新淘宝域名）
client := registry.NewNpmMirrorRegistry()

// 华为云镜像（中国）
client := registry.NewHuaWeiCloudRegistry()

// 腾讯云镜像（中国）
client := registry.NewTencentRegistry()

// CNPM 镜像（中国）
client := registry.NewCnpmRegistry()

// Yarn 官方镜像
client := registry.NewYarnRegistry()

// NPM CouchDB 镜像
client := registry.NewNpmjsComRegistry()
```

#### 核心方法

##### `GetPackageInformation(ctx context.Context, packageName string) (*models.Package, error)`

获取 NPM 包的详细信息。

**参数:**
- `ctx`: 用于取消和超时控制的上下文
- `packageName`: 要查询的包名

**返回值:**
- `*models.Package`: 完整的包信息
- `error`: 请求失败时的错误

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

##### `GetRegistryInformation(ctx context.Context) (*models.RegistryInformation, error)`

获取 NPM 注册表的状态和元数据信息。

**参数:**
- `ctx`: 用于取消和超时控制的上下文

**返回值:**
- `*models.RegistryInformation`: 注册表状态信息
- `error`: 请求失败时的错误

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

##### `GetOptions() *Options`

返回注册表客户端的当前配置选项。

**返回值:**
- `*Options`: 当前配置选项

**示例:**
```go
options := client.GetOptions()
fmt.Printf("注册表 URL: %s\n", options.RegistryURL)
fmt.Printf("代理: %s\n", options.Proxy)
```

### 配置选项

#### `NewOptions() *Options`

创建具有默认值的新 Options 实例。

**默认值:**
- `RegistryURL`: "https://registry.npmjs.org"
- `Proxy`: "" (无代理)

**示例:**
```go
options := registry.NewOptions()
```

#### `SetRegistryURL(url string) *Options`

设置 NPM 注册表 URL。支持方法链式调用。

**参数:**
- `url`: 注册表 URL (例如, "https://registry.npmjs.org")

**返回值:**
- `*Options`: 更新后的选项对象，支持链式调用

**示例:**
```go
options := registry.NewOptions().
    SetRegistryURL("https://registry.npmmirror.com")
```

#### `SetProxy(proxyUrl string) *Options`

设置 HTTP 代理 URL。支持方法链式调用。

**参数:**
- `proxyUrl`: 代理 URL (例如, "http://proxy.corp.com:8080")

**返回值:**
- `*Options`: 更新后的选项对象，支持链式调用

**示例:**
```go
options := registry.NewOptions().
    SetProxy("http://user:pass@proxy.example.com:8080")

// 清除代理
options.SetProxy("")
```

#### `GetHttpClient() (*http.Client, error)`

使用当前配置创建 HTTP 客户端。

**返回值:**
- `*http.Client`: 配置好的 HTTP 客户端
- `error`: 代理 URL 无效时的错误

**示例:**
```go
options := registry.NewOptions().SetProxy("http://proxy.example.com:8080")
httpClient, err := options.GetHttpClient()
if err != nil {
    log.Fatalf("创建 HTTP 客户端失败: %v", err)
}

// 使用客户端进行自定义请求
resp, err := httpClient.Get("https://example.com")
```

### 镜像源

| 镜像源 | URL | 地区 | 推荐用途 |
|--------|-----|--------|-----------|
| **官方 NPM** | https://registry.npmjs.org | 全球 | 生产环境，全球应用 |
| **淘宝 NPM** | https://registry.npm.taobao.org | 中国 | 旧版应用（已废弃） |
| **NPM Mirror** | https://registry.npmmirror.com | 中国 | 中国地区应用 |
| **华为云** | https://mirrors.huaweicloud.com/repository/npm | 中国 | 华为云环境 |
| **腾讯云** | http://mirrors.cloud.tencent.com/npm | 中国 | 腾讯云环境 |
| **CNPM** | http://r.cnpmjs.org | 中国 | 社区维护 |
| **Yarn** | https://registry.yarnpkg.com | 全球 | Yarn 项目 |
| **NPM CouchDB** | https://skimdb.npmjs.com | 全球 | 元数据密集型应用 |

## 数据模型

### Package 模型

`Package` 结构包含完整的 NPM 包信息：

```go
type Package struct {
    ID             string                 `json:"_id"`            // 包 ID
    Rev            string                 `json:"_rev"`           // 修订号
    Name           string                 `json:"name"`           // 包名称
    Description    string                 `json:"description"`    // 包描述
    DistTags       map[string]string      `json:"dist-tags"`      // 分发标签
    Versions       map[string]Version     `json:"versions"`       // 所有版本
    Maintainers    []Maintainer           `json:"maintainers"`    // 维护者列表
    Time           map[string]string      `json:"time"`           // 发布时间
    Repository     Repository             `json:"repository"`     // 仓库信息
    ReadMe         string                 `json:"readme"`         // README 内容
    Homepage       string                 `json:"homepage"`       // 项目主页
    License        string                 `json:"license"`        // 许可证类型
    Keywords       []string               `json:"keywords"`       // 关键词
    Author         Author                 `json:"author"`         // 作者信息
    Contributors   []Contributor          `json:"contributors"`   // 贡献者
    // ... 其他字段
}
```

**常用使用模式:**

```go
// 获取最新版本
latestVersion := pkg.DistTags["latest"]

// 检查包是否已废弃
if pkg.Deprecated != "" {
    fmt.Printf("警告: 包已废弃: %s\n", pkg.Deprecated)
}

// 列出所有可用版本
for version := range pkg.Versions {
    fmt.Printf("可用版本: %s\n", version)
}

// 获取特定版本详情
if versionInfo, exists := pkg.Versions["1.0.0"]; exists {
    fmt.Printf("依赖: %+v\n", versionInfo.Dependencies)
    fmt.Printf("开发依赖: %+v\n", versionInfo.DevDependencies)
}
```

### Version 模型

```go
type Version struct {
    Name            string               `json:"name"`            // 包名称
    Version         string               `json:"version"`         // 版本号
    Description     string               `json:"description"`     // 描述
    Main            string               `json:"main"`            // 入口点
    Scripts         *Script              `json:"scripts"`         // NPM 脚本
    Dependencies    map[string]string    `json:"dependencies"`    // 运行时依赖
    DevDependencies map[string]string    `json:"devDependencies"` // 开发依赖
    Repository      *Repository          `json:"repository"`      // 仓库
    License         string               `json:"license"`         // 许可证
    Dist            *Dist                `json:"dist"`            // 分发信息
    // ... 其他字段
}
```

### RegistryInformation 模型

```go
type RegistryInformation struct {
    DbName            string `json:"db_name"`              // 数据库名称
    DocCount          int    `json:"doc_count"`            // 包总数
    DocDelCount       int    `json:"doc_del_count"`        // 已删除包数
    UpdateSeq         int    `json:"update_seq"`           // 更新序列
    CompactRunning    bool   `json:"compact_running"`      // 压缩状态
    DiskSize          int64  `json:"disk_size"`            // 磁盘使用
    DataSize          int64  `json:"data_size"`            // 数据大小
    InstanceStartTime string `json:"instance_start_time"`  // 启动时间
    // ... 其他字段
}
```

## 高级用法

### 并发包查询

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
    packages := []string{"react", "vue", "angular", "svelte", "preact"}
    
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    var wg sync.WaitGroup
    results := make(chan string, len(packages))
    
    for _, pkg := range packages {
        wg.Add(1)
        go func(packageName string) {
            defer wg.Done()
            
            info, err := client.GetPackageInformation(ctx, packageName)
            if err != nil {
                results <- fmt.Sprintf("获取 %s 错误: %v", packageName, err)
                return
            }
            
            results <- fmt.Sprintf("%s: %s (最新: %s)", 
                info.Name, info.Description, info.DistTags["latest"])
        }(pkg)
    }
    
    go func() {
        wg.Wait()
        close(results)
    }()
    
    for result := range results {
        fmt.Println(result)
    }
}
```

### 注册表健康监控

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

func monitorRegistry(client *registry.Registry, interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
            
            info, err := client.GetRegistryInformation(ctx)
            if err != nil {
                fmt.Printf("注册表健康检查失败: %v\n", err)
                cancel()
                continue
            }
            
            fmt.Printf("注册表状态: %s，包数: %d，压缩: %t\n",
                info.DbName, info.DocCount, info.CompactRunning)
            
            cancel()
        }
    }
}

func main() {
    client := registry.NewRegistry()
    monitorRegistry(client, 5*time.Minute)
}
```

### 包依赖分析

```go
func analyzeDependencies(client *registry.Registry, packageName, version string) {
    ctx := context.Background()
    
    pkg, err := client.GetPackageInformation(ctx, packageName)
    if err != nil {
        log.Fatalf("获取包信息失败: %v", err)
    }
    
    versionInfo, exists := pkg.Versions[version]
    if !exists {
        // 尝试最新版本
        if latest, ok := pkg.DistTags["latest"]; ok {
            versionInfo, exists = pkg.Versions[latest]
            version = latest
        }
    }
    
    if !exists {
        fmt.Printf("版本 %s 未找到\n", version)
        return
    }
    
    fmt.Printf("分析 %s@%s\n", packageName, version)
    fmt.Printf("运行时依赖 (%d):\n", len(versionInfo.Dependencies))
    
    for dep, ver := range versionInfo.Dependencies {
        fmt.Printf("  %s: %s\n", dep, ver)
    }
    
    fmt.Printf("\n开发依赖 (%d):\n", len(versionInfo.DevDependencies))
    
    for dep, ver := range versionInfo.DevDependencies {
        fmt.Printf("  %s: %s\n", dep, ver)
    }
}
```

### 自定义注册表轮询

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

type RegistryPool struct {
    registries []*registry.Registry
    current    int
}

func NewRegistryPool() *RegistryPool {
    return &RegistryPool{
        registries: []*registry.Registry{
            registry.NewRegistry(),              // 官方
            registry.NewNpmMirrorRegistry(),     // NPM Mirror
            registry.NewTaoBaoRegistry(),        // 淘宝
            registry.NewHuaWeiCloudRegistry(),   // 华为
        },
        current: 0,
    }
}

func (rp *RegistryPool) GetPackageWithFallback(ctx context.Context, packageName string) (*models.Package, error) {
    for i := 0; i < len(rp.registries); i++ {
        client := rp.registries[rp.current]
        
        pkg, err := client.GetPackageInformation(ctx, packageName)
        if err == nil {
            return pkg, nil
        }
        
        fmt.Printf("注册表 %d 失败，尝试下一个: %v\n", rp.current, err)
        rp.current = (rp.current + 1) % len(rp.registries)
    }
    
    return nil, fmt.Errorf("所有注册表对包 %s 都失败了", packageName)
}
```

## 错误处理

### 常见错误类型

1. **网络错误**: 连接超时、DNS 失败
2. **HTTP 错误**: 404 (包未找到)、503 (服务不可用)
3. **解析错误**: 无效 JSON 响应
4. **上下文错误**: 超时或取消

### 错误处理最佳实践

```go
func robustPackageQuery(client *registry.Registry, packageName string) error {
    // 创建带超时和取消的上下文
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    pkg, err := client.GetPackageInformation(ctx, packageName)
    if err != nil {
        // 处理不同错误类型
        switch {
        case ctx.Err() == context.DeadlineExceeded:
            return fmt.Errorf("包 %s 请求超时", packageName)
        case ctx.Err() == context.Canceled:
            return fmt.Errorf("包 %s 请求被取消", packageName)
        default:
            return fmt.Errorf("获取包 %s 失败: %w", packageName, err)
        }
    }
    
    // 验证接收到的数据
    if pkg.Name == "" {
        return fmt.Errorf("包 %s 接收到无效数据", packageName)
    }
    
    if len(pkg.DistTags) == 0 {
        return fmt.Errorf("包 %s 没有分发标签", packageName)
    }
    
    fmt.Printf("成功获取包: %s\n", pkg.Name)
    return nil
}
```

### 指数退避重试逻辑

```go
import (
    "math"
    "math/rand"
    "time"
)

func retryWithBackoff(fn func() error, maxRetries int) error {
    for attempt := 0; attempt < maxRetries; attempt++ {
        err := fn()
        if err == nil {
            return nil
        }
        
        if attempt == maxRetries-1 {
            return fmt.Errorf("最大重试次数已达: %w", err)
        }
        
        // 指数退避加抖动
        backoff := time.Duration(math.Pow(2, float64(attempt))) * time.Second
        jitter := time.Duration(rand.Intn(1000)) * time.Millisecond
        time.Sleep(backoff + jitter)
        
        fmt.Printf("重试第 %d 次，错误: %v\n", attempt+1, err)
    }
    
    return nil
}

// 用法
err := retryWithBackoff(func() error {
    _, err := client.GetPackageInformation(ctx, "some-package")
    return err
}, 3)
```

## 最佳实践

### 1. 上下文使用

始终使用适当超时的上下文：

```go
// 好: 操作特定超时
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

// 更好: 用户操作可取消上下文
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// 处理取消信号
go func() {
    <-interrupt
    cancel()
}()
```

### 2. 注册表选择

根据部署区域选择注册表：

```go
func selectOptimalRegistry() *registry.Registry {
    // 检测区域或使用配置
    region := os.Getenv("DEPLOYMENT_REGION")
    
    switch region {
    case "china", "cn":
        return registry.NewNpmMirrorRegistry()
    case "huawei-cloud":
        return registry.NewHuaWeiCloudRegistry()
    default:
        return registry.NewRegistry() // 官方 NPM
    }
}
```

### 3. 资源管理

```go
// 高频操作使用 sync.Pool
var clientPool = sync.Pool{
    New: func() interface{} {
        return registry.NewRegistry()
    },
}

func getClient() *registry.Registry {
    return clientPool.Get().(*registry.Registry)
}

func putClient(client *registry.Registry) {
    clientPool.Put(client)
}
```

### 4. 配置管理

```go
type Config struct {
    RegistryURL string `env:"NPM_REGISTRY_URL" default:"https://registry.npmjs.org"`
    ProxyURL    string `env:"HTTP_PROXY"`
    Timeout     time.Duration `env:"NPM_TIMEOUT" default:"30s"`
}

func createClientFromConfig(cfg Config) *registry.Registry {
    options := registry.NewOptions().SetRegistryURL(cfg.RegistryURL)
    
    if cfg.ProxyURL != "" {
        options.SetProxy(cfg.ProxyURL)
    }
    
    return registry.NewRegistry(options)
}
```

## 故障排除

### 常见问题和解决方案

#### 1. 连接超时

**问题**: 请求频繁超时

**解决方案**:
```go
// 增加超时时间
ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

// 使用更快的镜像
client := registry.NewNpmMirrorRegistry() // 中国用户

// 检查代理配置
options := registry.NewOptions().SetProxy("") // 禁用代理
client := registry.NewRegistry(options)
```

#### 2. 包未找到 (404)

**问题**: 包存在但返回 404

**解决方案**:
```go
// 检查包名拼写
pkg, err := client.GetPackageInformation(ctx, "react") // 不是 "React"

// 尝试不同注册表
client := registry.NewYarnRegistry()

// 首先在 npmjs.com 验证包存在
```

#### 3. 代理问题

**问题**: 企业代理阻塞请求

**解决方案**:
```go
// 使用认证代理
options := registry.NewOptions().
    SetProxy("http://username:password@proxy.corp.com:8080")

// 尝试不同代理格式
options.SetProxy("http://proxy.corp.com:8080")

// 为特定域名绕过代理（需要自定义 HTTP 客户端）
```

#### 4. 速率限制

**问题**: 请求过多导致速率限制

**解决方案**:
```go
// 请求间添加延迟
time.Sleep(100 * time.Millisecond)

// 使用不同镜像分散负载
registries := []*registry.Registry{
    registry.NewRegistry(),
    registry.NewNpmMirrorRegistry(),
    registry.NewTaoBaoRegistry(),
}

// 实现请求队列
type RequestQueue struct {
    ch chan func()
}

func (rq *RequestQueue) Execute(fn func()) {
    rq.ch <- fn
}
```

#### 5. 大包内存使用

**问题**: 查询大包时内存使用过高

**解决方案**:
```go
// 只访问需要的字段
pkg, err := client.GetPackageInformation(ctx, packageName)
if err != nil {
    return err
}

// 立即提取所需数据
result := struct {
    Name        string
    Description string
    Latest      string
}{
    Name:        pkg.Name,
    Description: pkg.Description,
    Latest:      pkg.DistTags["latest"],
}

// 清除原始包数据
pkg = nil
```

## 示例

### 示例 1: 包信息仪表板

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

type PackageInfo struct {
    Name         string    `json:"name"`
    Description  string    `json:"description"`
    Latest       string    `json:"latest"`
    License      string    `json:"license"`
    LastModified time.Time `json:"lastModified"`
    Homepage     string    `json:"homepage"`
}

func packageHandler(w http.ResponseWriter, r *http.Request) {
    packageName := r.URL.Query().Get("package")
    if packageName == "" {
        http.Error(w, "需要 package 参数", http.StatusBadRequest)
        return
    }
    
    client := registry.NewRegistry()
    ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
    defer cancel()
    
    pkg, err := client.GetPackageInformation(ctx, packageName)
    if err != nil {
        http.Error(w, fmt.Sprintf("包未找到: %v", err), http.StatusNotFound)
        return
    }
    
    // 解析最后修改时间
    lastModified, _ := time.Parse(time.RFC3339, pkg.Time["modified"])
    
    info := PackageInfo{
        Name:         pkg.Name,
        Description:  pkg.Description,
        Latest:       pkg.DistTags["latest"],
        License:      pkg.License,
        LastModified: lastModified,
        Homepage:     pkg.Homepage,
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(info)
}

func main() {
    http.HandleFunc("/package", packageHandler)
    fmt.Println("服务器启动在 :8080")
    http.ListenAndServe(":8080", nil)
}
```

### 示例 2: 依赖树构建器

```go
package main

import (
    "context"
    "fmt"
    "strings"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

type DependencyNode struct {
    Name         string            `json:"name"`
    Version      string            `json:"version"`
    Dependencies []*DependencyNode `json:"dependencies,omitempty"`
}

func buildDependencyTree(client *registry.Registry, ctx context.Context, packageName, version string, depth int, maxDepth int) (*DependencyNode, error) {
    if depth > maxDepth {
        return &DependencyNode{Name: packageName, Version: version}, nil
    }
    
    pkg, err := client.GetPackageInformation(ctx, packageName)
    if err != nil {
        return nil, fmt.Errorf("获取包 %s 失败: %w", packageName, err)
    }
    
    node := &DependencyNode{
        Name:    packageName,
        Version: version,
    }
    
    // 获取版本信息
    versionInfo, exists := pkg.Versions[version]
    if !exists {
        // 尝试最新版本
        if latest, ok := pkg.DistTags["latest"]; ok {
            versionInfo, exists = pkg.Versions[latest]
            node.Version = latest
        }
    }
    
    if !exists {
        return node, nil
    }
    
    // 构建依赖
    for depName, depVersion := range versionInfo.Dependencies {
        // 清理版本规范（移除 ^, ~, 等）
        cleanVersion := strings.TrimLeft(depVersion, "^~>=<")
        
        depNode, err := buildDependencyTree(client, ctx, depName, cleanVersion, depth+1, maxDepth)
        if err != nil {
            fmt.Printf("警告: 无法解析依赖 %s: %v\n", depName, err)
            continue
        }
        
        node.Dependencies = append(node.Dependencies, depNode)
    }
    
    return node, nil
}

func main() {
    client := registry.NewRegistry()
    ctx := context.Background()
    
    tree, err := buildDependencyTree(client, ctx, "express", "4.18.2", 0, 2)
    if err != nil {
        log.Fatalf("构建依赖树失败: %v", err)
    }
    
    printTree(tree, 0)
}

func printTree(node *DependencyNode, indent int) {
    prefix := strings.Repeat("  ", indent)
    fmt.Printf("%s%s@%s\n", prefix, node.Name, node.Version)
    
    for _, dep := range node.Dependencies {
        printTree(dep, indent+1)
    }
}
```

### 示例 3: 包安全审计

```go
package main

import (
    "context"
    "fmt"
    "regexp"
    "strings"
    "time"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

type SecurityAudit struct {
    PackageName        string    `json:"packageName"`
    CurrentVersion     string    `json:"currentVersion"`
    LatestVersion      string    `json:"latestVersion"`
    IsOutdated         bool      `json:"isOutdated"`
    IsDeprecated       bool      `json:"isDeprecated"`
    HasVulnerabilities bool      `json:"hasVulnerabilities"`
    LastUpdate         time.Time `json:"lastUpdate"`
    Maintainers        int       `json:"maintainersCount"`
    License            string    `json:"license"`
    Issues             []string  `json:"issues"`
}

func auditPackage(client *registry.Registry, ctx context.Context, packageName, currentVersion string) (*SecurityAudit, error) {
    pkg, err := client.GetPackageInformation(ctx, packageName)
    if err != nil {
        return nil, fmt.Errorf("获取包信息失败: %w", err)
    }
    
    audit := &SecurityAudit{
        PackageName:    packageName,
        CurrentVersion: currentVersion,
        LatestVersion:  pkg.DistTags["latest"],
        IsDeprecated:   pkg.Deprecated != "",
        Maintainers:    len(pkg.Maintainers),
        License:        pkg.License,
    }
    
    // 检查版本是否过时
    audit.IsOutdated = currentVersion != audit.LatestVersion
    
    // 解析最后更新时间
    if timeStr, exists := pkg.Time[audit.LatestVersion]; exists {
        audit.LastUpdate, _ = time.Parse(time.RFC3339, timeStr)
    }
    
    // 检查潜在安全问题
    audit.Issues = checkSecurityIssues(pkg)
    audit.HasVulnerabilities = len(audit.Issues) > 0
    
    return audit, nil
}

func checkSecurityIssues(pkg *models.Package) []string {
    var issues []string
    
    // 检查已废弃包
    if pkg.Deprecated != "" {
        issues = append(issues, fmt.Sprintf("包已废弃: %s", pkg.Deprecated))
    }
    
    // 检查描述中的可疑模式
    suspiciousPatterns := []string{
        "bitcoin", "cryptocurrency", "mining", "keylogger",
        "password", "steal", "backdoor",
    }
    
    description := strings.ToLower(pkg.Description)
    for _, pattern := range suspiciousPatterns {
        if strings.Contains(description, pattern) {
            issues = append(issues, fmt.Sprintf("描述中的可疑关键词: %s", pattern))
        }
    }
    
    // 检查可疑仓库 URL
    if pkg.Repository.Url != "" {
        if matched, _ := regexp.MatchString(`(bitbucket\.org|gitlab\.com|github\.com)`, pkg.Repository.Url); !matched {
            issues = append(issues, "仓库托管在非标准平台")
        }
    }
    
    // 检查维护者数量
    if len(pkg.Maintainers) == 0 {
        issues = append(issues, "包没有维护者")
    } else if len(pkg.Maintainers) > 10 {
        issues = append(issues, "包的维护者异常多")
    }
    
    // 检查缺失许可证
    if pkg.License == "" {
        issues = append(issues, "包没有指定许可证")
    }
    
    return issues
}

func main() {
    client := registry.NewRegistry()
    ctx := context.Background()
    
    packages := map[string]string{
        "express":    "4.17.1",
        "lodash":     "4.17.20",
        "react":      "17.0.0",
        "vue":        "2.6.12",
    }
    
    fmt.Println("安全审计报告")
    fmt.Println("============")
    
    for packageName, version := range packages {
        audit, err := auditPackage(client, ctx, packageName, version)
        if err != nil {
            fmt.Printf("审计 %s 失败: %v\n", packageName, err)
            continue
        }
        
        fmt.Printf("\n包: %s\n", audit.PackageName)
        fmt.Printf("当前版本: %s\n", audit.CurrentVersion)
        fmt.Printf("最新版本: %s\n", audit.LatestVersion)
        fmt.Printf("状态: ")
        
        if audit.IsDeprecated {
            fmt.Print("已废弃 ")
        }
        if audit.IsOutdated {
            fmt.Print("过时 ")
        }
        if audit.HasVulnerabilities {
            fmt.Print("有问题 ")
        }
        if !audit.IsDeprecated && !audit.IsOutdated && !audit.HasVulnerabilities {
            fmt.Print("正常")
        }
        fmt.Println()
        
        if len(audit.Issues) > 0 {
            fmt.Println("问题:")
            for _, issue := range audit.Issues {
                fmt.Printf("  - %s\n", issue)
            }
        }
    }
}
```

---

## 贡献

我们欢迎对本文档的改进贡献！请随时：

- 报告问题或不一致之处
- 建议新示例或用例
- 添加翻译
- 改进现有内容

更多信息，请参阅我们的[贡献指南](https://github.com/scagogogo/npm-skills/blob/main/README.md#贡献指南)。

## 许可证

本文档是 NPM Crawler 项目的一部分，采用 MIT 许可证。详情请参阅 [LICENSE](https://github.com/scagogogo/npm-skills/blob/main/LICENSE)。
