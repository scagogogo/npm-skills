# 快速开始

本指南将帮助您快速上手使用 NPM Crawler。

## 环境要求

- Go 1.20 或更高版本
- 网络连接（或访问 NPM 注册表镜像）

## 安装

使用 Go 模块安装 NPM Crawler：

```bash
go get github.com/scagogogo/npm-skills
```

## 基本使用

### 1. 创建Registry客户端

```go
package main

import (
    "github.com/scagogogo/npm-skills/pkg/registry"
)

func main() {
    // 使用默认配置（官方NPM Registry）
    client := registry.NewRegistry()
    
    // 或者使用自定义配置
    options := registry.NewOptions().
        SetRegistryURL("https://registry.npmjs.org").
        SetProxy("http://proxy.example.com:8080")
    client := registry.NewRegistry(options)
}
```

### 2. 获取包信息

```go
ctx := context.Background()

// 获取包的完整信息
pkg, err := client.GetPackageInformation(ctx, "react")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("包名: %s\n", pkg.Name)
fmt.Printf("描述: %s\n", pkg.Description)
fmt.Printf("最新版本: %s\n", pkg.DistTags["latest"])
fmt.Printf("主页: %s\n", pkg.Homepage)
```

### 3. 获取Registry状态

```go
info, err := client.GetRegistryInformation(ctx)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Registry: %s\n", info.DbName)
fmt.Printf("包总数: %d\n", info.DocCount)
fmt.Printf("数据大小: %d bytes\n", info.DataSize)
```

### 4. 搜索包

```go
// 搜索包，限制返回10个结果
result, err := client.SearchPackages(ctx, "react ui", 10)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("找到 %d 个结果\n", result.Total)
for _, obj := range result.Objects {
    fmt.Printf("- %s: %s\n", obj.Package.Name, obj.Package.Description)
}
```

### 5. 获取特定版本信息

```go
// 获取特定版本的详细信息
version, err := client.GetPackageVersion(ctx, "react", "18.2.0")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("版本: %s\n", version.Version)
fmt.Printf("依赖: %+v\n", version.Dependencies)
```

### 6. 获取下载统计

```go
// 获取最近一周的下载统计
stats, err := client.GetDownloadStats(ctx, "react", "last-week")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("包: %s\n", stats.Package)
fmt.Printf("下载次数: %d\n", stats.Downloads)
fmt.Printf("统计周期: %s 到 %s\n", stats.Start, stats.End)
```

## 使用不同的镜像源

NPM Crawler 内置支持多种镜像源，特别适合中国大陆用户：

```go
// 使用淘宝镜像
taobaoClient := registry.NewTaoBaoRegistry()

// 使用华为云镜像
huaweiClient := registry.NewHuaWeiCloudRegistry()

// 使用NPM Mirror
npmMirrorClient := registry.NewNpmMirrorRegistry()

// 使用腾讯云镜像
tencentClient := registry.NewTencentRegistry()

// 使用CNPM镜像
cnpmClient := registry.NewCnpmRegistry()

// 使用Yarn镜像
yarnClient := registry.NewYarnRegistry()
```

## 配置代理

如果您在企业网络环境中，可能需要配置HTTP代理：

```go
options := registry.NewOptions().
    SetRegistryURL("https://registry.npmjs.org").
    SetProxy("http://proxy.corp.com:8080")

client := registry.NewRegistry(options)

// 或者使用带认证的代理
options.SetProxy("http://username:password@proxy.corp.com:8080")
```

## 错误处理

建议在生产环境中添加适当的错误处理：

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
    client := registry.NewRegistry()
    
    // 设置超时
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    pkg, err := client.GetPackageInformation(ctx, "react")
    if err != nil {
        switch {
        case ctx.Err() == context.DeadlineExceeded:
            log.Fatal("请求超时")
        case ctx.Err() == context.Canceled:
            log.Fatal("请求被取消")
        default:
            log.Fatalf("获取包信息失败: %v", err)
        }
    }
    
    fmt.Printf("成功获取包信息: %s\n", pkg.Name)
}
```

## 并发使用

NPM Crawler 是线程安全的，可以在并发环境中使用：

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
    
    packages := []string{"react", "vue", "angular", "lodash", "express"}
    
    var wg sync.WaitGroup
    for _, pkgName := range packages {
        wg.Add(1)
        go func(name string) {
            defer wg.Done()
            
            pkg, err := client.GetPackageInformation(ctx, name)
            if err != nil {
                fmt.Printf("获取 %s 失败: %v\n", name, err)
                return
            }
            
            fmt.Printf("%s: %s\n", pkg.Name, pkg.DistTags["latest"])
        }(pkgName)
    }
    
    wg.Wait()
}
```

## 下一步

- 查看完整的 [API 文档](/api/) 了解所有可用功能
- 浏览 [示例代码](/examples/basic) 学习更多用法
- 访问 [GitHub 仓库](https://github.com/scagogogo/npm-skills) 查看源代码和贡献指南
