# Improve Documentation - Complete API Coverage

> **For agentic workers:** REQUIRED SUB-SKILL: `superpowers:subagent-driven-development`
> Steps use checkbox (`- [ ]`) syntax.

**Goal:** 完善 npm-crawler 文档网站，补全所有 API 文档，重点补充新添加的 `DownloadTarball` 方法，并添加相关示例。

**Architecture:** 用户访问 VitePress 文档网站 → 浏览 API 文档 → 查看示例代码 → 复制使用。`DownloadTarball` 方法文档独立成节，同时在 API 概述和方法列表中补充。

**Tech Stack:** VitePress, Go, markdown

**Risks:**
- Task 1 纯文档修改，无风险
- Task 2 VitePress sidebar 配置需仔细检查格式 → 缓解：最小化改动，只添加一个链接

---

### Task 1: 补充 API 文档 — 添加 DownloadTarball

**Depends on:** None
**Files:**
- Modify: `docs/api/registry.md`（补充 DownloadTarball 章节）
- Modify: `docs/en/api/registry.md`（补充英文版 DownloadTarball 章节）
- Modify: `docs/api/index.md`（方法列表补充 DownloadTarball）
- Modify: `docs/en/api/index.md`（英文版方法列表补充 DownloadTarball）

- [ ] **Step 1: 在 `docs/api/registry.md` 的 `GetDownloadStats` 方法后添加 `DownloadTarball` 文档**

文件: `docs/api/registry.md`

在 `### GetDownloadStats` 章节之后（大约第 165 行附近），添加以下内容：

```markdown
### `DownloadTarball(ctx context.Context, packageName, version, destPath string) error`

下载指定 NPM 包的 tarball 文件到本地路径。

**参数:**
- `ctx` - 上下文，用于取消和超时控制
- `packageName` - 要下载的包名称，例如 "react"、"lodash" 等
- `version` - 要下载的版本号，例如 "18.0.0"、"latest" 等
- `destPath` - 目标文件保存路径，例如 "./downloads/react-18.0.0.tgz"

**返回值:**
- `error` - 如果下载失败则返回错误

**示例:**
```go
ctx := context.Background()

// 下载 react 18.0.0 版本到本地文件
err := client.DownloadTarball(ctx, "react", "18.0.0", "./react.tgz")
if err != nil {
    return fmt.Errorf("下载 tarball 失败: %w", err)
}

// 使用 latest 下载最新版本
err = client.DownloadTarball(ctx, "vue", "latest", "./vue.tgz")
if err != nil {
    return fmt.Errorf("下载 tarball 失败: %w", err)
}

// 验证下载的文件
info, err := os.Stat("./react.tgz")
if err != nil {
    return err
}
fmt.Printf("文件大小: %d bytes\n", info.Size())
```

**使用 CNPM 镜像下载示例:**
```go
// 使用国内镜像下载，速度更快
options := registry.NewOptions().SetRegistryURL(registry.RegistryUrlCnpm)
client := registry.NewRegistry(options)

err := client.DownloadTarball(ctx, "axios", "1.0.0", "/tmp/axios.tgz")
if err != nil {
    log.Fatalf("下载失败: %v", err)
}

fmt.Println("下载成功！")
```
```

- [ ] **Step 2: 在 `docs/api/index.md` 的 `GetDownloadStats` 描述后添加 `DownloadTarball` 方法条目**

文件: `docs/api/index.md`

在"#### GetDownloadStats"条目之后（大约第 136 行附近），添加：

```markdown
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
```

- [ ] **Step 3: 在 `docs/en/api/registry.md` 的 `GetDownloadStats` 方法后添加 `DownloadTarball` 英文文档**

文件: `docs/en/api/registry.md`

在 `### GetDownloadStats` 章节之后（大约第 166 行附近），添加：

```markdown
### DownloadTarball
```go
func (r *Registry) DownloadTarball(ctx context.Context, packageName, version, destPath string) error
```

Downloads an NPM package tarball to a local file path.

**Parameters:**
- `ctx` - Context for cancellation and timeout control
- `packageName` - Name of the package to download
- `version` - Version to download (e.g., "18.0.0" or "latest")
- `destPath` - Local file path where the tarball will be saved

**Returns:**
- `error` - Error if the download fails

**Example:**
```go
ctx := context.Background()

// Download specific version
err := client.DownloadTarball(ctx, "react", "18.0.0", "./react.tgz")
if err != nil {
    return fmt.Errorf("download failed: %w", err)
}

// Download latest version
err = client.DownloadTarball(ctx, "vue", "latest", "./vue.tgz")
if err != nil {
    return fmt.Errorf("download failed: %w", err)
}

// Verify the downloaded file
info, err := os.Stat("./react.tgz")
if err != nil {
    return err
}
fmt.Printf("File size: %d bytes\n", info.Size())
```

**Using CNPM mirror for faster downloads in China:**
```go
options := registry.NewOptions().SetRegistryURL(registry.RegistryUrlCnpm)
client := registry.NewRegistry(options)

err := client.DownloadTarball(ctx, "axios", "1.0.0", "/tmp/axios.tgz")
if err != nil {
    log.Fatalf("Download failed: %v", err)
}

fmt.Println("Download successful!")
```
```

- [ ] **Step 4: 在 `docs/en/api/index.md` 的 `GetDownloadStats` 描述后添加 `DownloadTarball` 英文方法条目**

文件: `docs/en/api/index.md`

在"#### GetDownloadStats"条目之后（大约第 135 行附近），添加：

```markdown
#### DownloadTarball
Download an NPM package tarball to a local file.

```go
func (r *Registry) DownloadTarball(ctx context.Context, packageName, version, destPath string) error
```

**Parameters:**
- `packageName` - Package name
- `version` - Version number or tag (e.g., "latest")
- `destPath` - Local file path to save the tarball

**Example:**
```go
err := client.DownloadTarball(ctx, "react", "18.2.0", "./react.tgz")
```
```

- [ ] **Step 5: 验证 VitePress 文档构建**

Run: `cd /Users/cc11001100/github/scagogogo/npm-crawler/docs && npm run build 2>&1 | head -30`
Expected:
  - Exit code: 0
  - No error output
  - Output contains: "build complete" or similar

---

### Task 2: 添加 DownloadTarball 示例 + 更新 VitePress Sidebar

**Depends on:** Task 1
**Files:**
- Create: `docs/examples/download.md`（中文 tarball 下载示例）
- Create: `docs/en/examples/download.md`（英文 tarball 下载示例）
- Modify: `docs/.vitepress/config.ts`（sidebar 添加示例链接）

- [ ] **Step 1: 创建 `docs/examples/download.md` — 中文版 tarball 下载示例**

文件: `docs/examples/download.md`

```markdown
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

    "github.com/scagogogo/npm-crawler/pkg/registry"
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

    "github.com/scagogogo/npm-crawler/pkg/registry"
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

    "github.com/scagogogo/npm-crawler/pkg/registry"
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

    "github.com/scagogogo/npm-crawler/pkg/registry"
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
go get github.com/scagogogo/npm-crawler
go run main.go
```

## 下一步

- 查看 [基本用法示例](/examples/basic) 学习更多基础功能
- 阅读 [API 文档](/api/registry) 了解完整的 API 参考
- 了解 [镜像源配置](/examples/mirrors) 选择最快的下载源
```

- [ ] **Step 2: 创建 `docs/en/examples/download.md` — 英文版 tarball 下载示例**

文件: `docs/en/examples/download.md`

```markdown
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

    "github.com/scagogogo/npm-crawler/pkg/registry"
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

    "github.com/scagogogo/npm-crawler/pkg/registry"
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

    "github.com/scagogogo/npm-crawler/pkg/registry"
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

    "github.com/scagogogo/npm-crawler/pkg/registry"
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
go get github.com/scagogogo/npm-crawler
go run main.go
```

## Next Steps

- Check [Basic Usage Examples](/en/examples/basic) for more fundamentals
- Read [API Documentation](/en/api/registry) for complete API reference
- Explore [Mirror Configuration](/en/examples/mirrors) to choose the fastest download source
```

- [ ] **Step 3: 更新 VitePress sidebar — 在示例部分添加 download 示例链接**

文件: `docs/.vitepress/config.ts`

在中文 sidebar 的 "示例" 部分，将：
```ts
{
  text: '示例',
  items: [
    { text: '基本用法', link: '/examples/basic' },
    { text: '高级用法', link: '/examples/advanced' },
    { text: '镜像源配置', link: '/examples/mirrors' }
  ]
}
```

替换为：
```ts
{
  text: '示例',
  items: [
    { text: '基本用法', link: '/examples/basic' },
    { text: '高级用法', link: '/examples/advanced' },
    { text: '镜像源配置', link: '/examples/mirrors' },
    { text: '下载 Tarball', link: '/examples/download' }
  ]
}
```

在英文 sidebar 的 "Examples" 部分，将：
```ts
{
  text: 'Examples',
  items: [
    { text: 'Basic Usage', link: '/en/examples/basic' },
    { text: 'Advanced Usage', link: '/en/examples/advanced' },
    { text: 'Mirror Configuration', link: '/en/examples/mirrors' }
  ]
}
```

替换为：
```ts
{
  text: 'Examples',
  items: [
    { text: 'Basic Usage', link: '/en/examples/basic' },
    { text: 'Advanced Usage', link: '/en/examples/advanced' },
    { text: 'Mirror Configuration', link: '/en/examples/mirrors' },
    { text: 'Download Tarball', link: '/en/examples/download' }
  ]
}
```

- [ ] **Step 4: 验证 VitePress 站点构建**

Run: `cd /Users/cc11001100/github/scagogogo/npm-crawler/docs && npm run build 2>&1 | tail -10`
Expected:
  - Exit code: 0
  - Output contains: "build complete"

- [ ] **Step 5: 提交所有文档变更**

Run: `cd /Users/cc11001100/github/scagogogo/npm-crawler && git add docs/api/registry.md docs/en/api/registry.md docs/api/index.md docs/en/api/index.md docs/examples/download.md docs/en/examples/download.md docs/.vitepress/config.ts && git commit -m "docs: add DownloadTarball API documentation and examples"`
