# 基本用法示例

本页面包含了 NPM Crawler 的基本使用示例，帮助您快速上手。

## 示例 1: 获取包信息

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

func main() {
    // 创建客户端
    client := registry.NewRegistry()
    ctx := context.Background()
    
    // 获取包信息
    pkg, err := client.GetPackageInformation(ctx, "react")
    if err != nil {
        log.Fatal(err)
    }
    
    // 显示基本信息
    fmt.Printf("包名: %s\n", pkg.Name)
    fmt.Printf("描述: %s\n", pkg.Description)
    fmt.Printf("最新版本: %s\n", pkg.DistTags["latest"])
    fmt.Printf("许可证: %s\n", pkg.License)
    fmt.Printf("主页: %s\n", pkg.Homepage)
    fmt.Printf("作者: %s\n", pkg.Author.Name)
}
```

## 示例 2: 使用不同镜像源

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

func main() {
    ctx := context.Background()
    
    // 使用不同的镜像源
    mirrors := map[string]*registry.Registry{
        "官方源":   registry.NewRegistry(),
        "淘宝镜像":  registry.NewTaoBaoRegistry(),
        "NPM镜像": registry.NewNpmMirrorRegistry(),
        "华为云":   registry.NewHuaWeiCloudRegistry(),
    }
    
    packageName := "lodash"
    
    for name, client := range mirrors {
        pkg, err := client.GetPackageInformation(ctx, packageName)
        if err != nil {
            fmt.Printf("%s: 获取失败 - %v\n", name, err)
            continue
        }
        
        fmt.Printf("%s: %s (v%s)\n", name, pkg.Name, pkg.DistTags["latest"])
    }
}
```

## 示例 3: 代理配置

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

func main() {
    // 配置代理
    options := registry.NewOptions().
        SetRegistryURL("https://registry.npmjs.org").
        SetProxy("http://proxy.example.com:8080")
    
    client := registry.NewRegistry(options)
    ctx := context.Background()
    
    pkg, err := client.GetPackageInformation(ctx, "express")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("通过代理获取包信息: %s\n", pkg.Name)
}
```

## 示例 4: 搜索包

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
    
    // 搜索包
    result, err := client.SearchPackages(ctx, "react ui component", 5)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("找到 %d 个结果 (耗时: %s)\n", result.Total, result.Time)
    fmt.Println("---")
    
    for i, obj := range result.Objects {
        pkg := obj.Package
        fmt.Printf("%d. %s (v%s)\n", i+1, pkg.Name, pkg.Version)
        fmt.Printf("   描述: %s\n", pkg.Description)
        fmt.Printf("   评分: %.2f\n", obj.Score.Final)
        fmt.Printf("   链接: %s\n", pkg.Links.NPM)
        fmt.Println()
    }
}
```

## 示例 5: 获取特定版本信息

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
    
    // 获取特定版本
    version, err := client.GetPackageVersion(ctx, "react", "17.0.2")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("包: %s\n", version.Name)
    fmt.Printf("版本: %s\n", version.Version)
    fmt.Printf("描述: %s\n", version.Description)
    fmt.Printf("主文件: %s\n", version.Main)
    
    // 显示依赖
    if len(version.Dependencies) > 0 {
        fmt.Println("\n运行时依赖:")
        for dep, ver := range version.Dependencies {
            fmt.Printf("  %s: %s\n", dep, ver)
        }
    }
    
    // 显示开发依赖
    if len(version.DevDependencies) > 0 {
        fmt.Println("\n开发依赖:")
        for dep, ver := range version.DevDependencies {
            fmt.Printf("  %s: %s\n", dep, ver)
        }
    }
}
```

## 示例 6: 获取下载统计

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
    
    for _, pkg := range packages {
        fmt.Printf("\n=== %s ===\n", pkg)
        
        for _, period := range periods {
            stats, err := client.GetDownloadStats(ctx, pkg, period)
            if err != nil {
                fmt.Printf("%s: 获取失败 - %v\n", period, err)
                continue
            }
            
            fmt.Printf("%s: %d 次下载\n", period, stats.Downloads)
        }
    }
}
```

## 示例 7: 错误处理

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "log"
    "time"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

func main() {
    client := registry.NewRegistry()
    
    // 设置超时
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    pkg, err := client.GetPackageInformation(ctx, "nonexistent-package-12345")
    if err != nil {
        switch {
        case errors.Is(err, context.DeadlineExceeded):
            log.Println("请求超时")
        case errors.Is(err, context.Canceled):
            log.Println("请求被取消")
        default:
            log.Printf("获取包信息失败: %v\n", err)
        }
        return
    }
    
    fmt.Printf("包信息: %s\n", pkg.Name)
}
```

## 示例 8: 比较版本信息

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
    
    packageName := "lodash"
    versions := []string{"4.17.20", "4.17.21"}
    
    fmt.Printf("比较 %s 的不同版本:\n\n", packageName)
    
    for _, ver := range versions {
        version, err := client.GetPackageVersion(ctx, packageName, ver)
        if err != nil {
            log.Printf("获取版本 %s 失败: %v\n", ver, err)
            continue
        }
        
        fmt.Printf("版本 %s:\n", version.Version)
        fmt.Printf("  描述: %s\n", version.Description)
        fmt.Printf("  依赖数量: %d\n", len(version.Dependencies))
        fmt.Printf("  开发依赖数量: %d\n", len(version.DevDependencies))
        
        if version.Deprecated != "" {
            fmt.Printf("  ⚠️  已弃用: %s\n", version.Deprecated)
        }
        
        fmt.Println()
    }
}
```

## 示例 9: 获取注册表状态

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

func main() {
    // 检查多个注册表的状态
    registries := map[string]*registry.Registry{
        "NPM官方": registry.NewRegistry(),
        "淘宝镜像": registry.NewTaoBaoRegistry(),
        "华为云":  registry.NewHuaWeiCloudRegistry(),
    }
    
    ctx := context.Background()
    
    for name, client := range registries {
        fmt.Printf("=== %s ===\n", name)
        
        info, err := client.GetRegistryInformation(ctx)
        if err != nil {
            fmt.Printf("获取失败: %v\n\n", err)
            continue
        }
        
        fmt.Printf("数据库名: %s\n", info.DbName)
        fmt.Printf("包总数: %d\n", info.DocCount)
        fmt.Printf("已删除包: %d\n", info.DocDelCount)
        fmt.Printf("数据大小: %d MB\n", info.DataSize/(1024*1024))
        fmt.Printf("磁盘使用: %d MB\n", info.DiskSize/(1024*1024))
        fmt.Printf("实例启动时间: %s\n", info.InstanceStartTime)
        fmt.Println()
    }
}
```

## 示例 10: 批量处理

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
    
    packages := []string{
        "react", "vue", "angular", "lodash", "express",
        "axios", "moment", "underscore", "jquery", "bootstrap",
    }
    
    var wg sync.WaitGroup
    results := make(chan string, len(packages))
    
    // 并发获取包信息
    for _, pkg := range packages {
        wg.Add(1)
        go func(packageName string) {
            defer wg.Done()
            
            info, err := client.GetPackageInformation(ctx, packageName)
            if err != nil {
                results <- fmt.Sprintf("❌ %s: %v", packageName, err)
                return
            }
            
            results <- fmt.Sprintf("✅ %s: v%s (%s)", 
                info.Name, 
                info.DistTags["latest"], 
                info.Description[:50]+"...")
        }(pkg)
    }
    
    // 等待所有 goroutine 完成
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // 收集并显示结果
    fmt.Println("批量获取包信息结果:")
    for result := range results {
        fmt.Println(result)
    }
}
```

## 运行示例

保存任意示例代码到文件（如 `example.go`），然后运行：

```bash
go mod init npm-skills-example
go get github.com/scagogogo/npm-skills
go run example.go
```

## 下一步

- 查看 [高级用法示例](/examples/advanced) 学习更复杂的用法
- 阅读 [镜像源配置](/examples/mirrors) 了解如何优化网络访问
- 浏览 [API 文档](/api/) 了解所有可用功能
