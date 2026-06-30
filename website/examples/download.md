# 下载 Tarball 示例

本页面展示如何使用 `DownloadTarball` 方法下载 NPM 包的 tarball 文件到本地。

## 示例 1: 下载指定版本

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

    // 删除已存在的文件（如果有）
    os.Remove(destPath)

    fmt.Printf("开始下载 %s@%s...\n", packageName, version)

    err := client.DownloadTarball(ctx, packageName, version, destPath)
    if err != nil {
        log.Fatalf("下载失败: %v", err)
    }

    // 验证文件
    info, err := os.Stat(destPath)
    if err != nil {
        log.Fatalf("文件状态检查失败: %v", err)
    }

    fmt.Printf("下载成功！\n")
    fmt.Printf("文件大小: %d bytes\n", info.Size())
}
```

## 示例 2: 使用 CNPM 镜像下载（中国大陆用户推荐）

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
    // 使用 CNPM 镜像，下载速度更快
    options := registry.NewOptions().SetRegistryURL(registry.RegistryUrlCnpm)
    client := registry.NewRegistry(options)

    ctx := context.Background()
    packageName := "react"
    version := "18.2.0"
    destDir := "/tmp/npm-tarballs"

    // 确保目标目录存在
    if err := os.MkdirAll(destDir, 0755); err != nil {
        log.Fatalf("创建目录失败: %v", err)
    }

    destPath := fmt.Sprintf("%s/%s-%s.tgz", destDir, packageName, version)
    os.Remove(destPath)

    fmt.Printf("使用 CNPM 镜像下载 %s@%s...\n", packageName, version)

    err := client.DownloadTarball(ctx, packageName, version, destPath)
    if err != nil {
        log.Fatalf("下载失败: %v", err)
    }

    info, _ := os.Stat(destPath)
    fmt.Printf("下载成功！文件大小: %d bytes\n", info.Size())
}
```

## 示例 3: 下载最新版本 (latest)

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

    // 获取版本信息
    versionInfo, err := client.GetPackageVersion(ctx, packageName, version)
    if err != nil {
        log.Fatalf("获取版本信息失败: %v", err)
    }

    fmt.Printf("最新版本: %s\n", versionInfo.Version)

    os.Remove(destPath)
    err = client.DownloadTarball(ctx, packageName, version, destPath)
    if err != nil {
        log.Fatalf("下载失败: %v", err)
    }

    info, _ := os.Stat(destPath)
    fmt.Printf("下载成功！文件大小: %d bytes\n", info.Size())
}
```

## 示例 4: 批量下载多个包

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
                log.Printf("下载 %s@%s 失败: %v\n", name, version, err)
                return
            }

            info, _ := os.Stat(destPath)
            fmt.Printf("✅ %s@%s: %d bytes\n", name, version, info.Size())
        }(pkg.name, pkg.version)
    }

    wg.Wait()
    fmt.Println("全部下载完成！")
}
```

## 运行示例

```bash
go mod init tarball-example
go get github.com/scagogogo/npm-skills
go run main.go
```

## 下一步

- 查看 [基本用法示例](/examples/basic) 学习更多基础功能
- 阅读 [API 文档](/api/registry) 了解完整的 API 参考
- 了解 [镜像源配置](/examples/mirrors) 选择最快的下载源
