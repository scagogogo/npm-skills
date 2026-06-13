# 安装指南

本指南将详细介绍如何在不同环境中安装和配置 NPM Crawler。

## 系统要求

- **Go版本**: 1.20 或更高版本
- **操作系统**: Linux, macOS, Windows
- **网络**: 能够访问 NPM 注册表或其镜像源

## Go 环境准备

### 安装 Go

如果您还没有安装 Go，请访问 [Go 官网](https://golang.org/dl/) 下载并安装最新版本。

验证安装：
```bash
go version
```

应该输出类似 `go version go1.21.0 linux/amd64` 的信息。

### 配置 Go 模块

确保您的项目启用了 Go 模块：

```bash
# 在您的项目目录中
go mod init your-project-name
```

## 安装 NPM Crawler

### 使用 go get（推荐）

```bash
go get github.com/scagogogo/npm-skills
```

### 添加到 go.mod

您也可以直接在 `go.mod` 文件中添加依赖：

```go
module your-project

go 1.20

require (
    github.com/scagogogo/npm-skills latest
)
```

然后运行：
```bash
go mod tidy
```

### 指定版本

如果需要特定版本：

```bash
# 安装特定版本
go get github.com/scagogogo/npm-skills@v1.0.0

# 安装最新的预发布版本
go get github.com/scagogogo/npm-skills@latest
```

## 验证安装

创建一个简单的测试文件来验证安装：

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
    
    // 获取一个简单的包信息
    pkg, err := client.GetPackageInformation(ctx, "lodash")
    if err != nil {
        log.Fatal("安装验证失败:", err)
    }
    
    fmt.Printf("✅ NPM Crawler 安装成功！\n")
    fmt.Printf("测试包: %s\n", pkg.Name)
    fmt.Printf("最新版本: %s\n", pkg.DistTags["latest"])
}
```

运行测试：
```bash
go run test.go
```

如果看到成功消息，说明安装正确。

## 网络配置

### 代理设置

如果您在企业网络环境中：

```go
options := registry.NewOptions().
    SetProxy("http://proxy.corp.com:8080")
client := registry.NewRegistry(options)
```

### 镜像源选择

根据您的地理位置选择合适的镜像源：

```go
// 中国大陆用户推荐
client := registry.NewNpmMirrorRegistry()  // NPM Mirror
// 或
client := registry.NewTaoBaoRegistry()     // 淘宝镜像

// 全球用户
client := registry.NewRegistry()           // 官方源
```

## 常见问题

### 1. 网络连接问题

**错误**: `dial tcp: lookup registry.npmjs.org: no such host`

**解决方案**:
- 检查网络连接
- 尝试使用国内镜像源
- 配置正确的代理

```go
// 使用国内镜像
client := registry.NewNpmMirrorRegistry()
```

### 2. 代理认证问题

**错误**: `407 Proxy Authentication Required`

**解决方案**:
```go
// 使用带认证的代理
options := registry.NewOptions().
    SetProxy("http://username:password@proxy.corp.com:8080")
```

### 3. Go 版本不兼容

**错误**: `package github.com/scagogogo/npm-skills: module requires Go 1.20`

**解决方案**:
- 升级 Go 到 1.20 或更高版本
- 或者使用兼容的早期版本

### 4. 模块导入问题

**错误**: `cannot find module providing package github.com/scagogogo/npm-skills`

**解决方案**:
```bash
# 清理模块缓存
go clean -modcache

# 重新下载依赖
go mod download

# 整理依赖
go mod tidy
```

## Docker 环境

如果您在 Docker 中使用：

```dockerfile
FROM golang:1.21-alpine

WORKDIR /app

# 复制依赖文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN go build -o main .

CMD ["./main"]
```

## CI/CD 集成

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

## 性能优化

### 连接复用

NPM Crawler 使用 Go 的标准 HTTP 客户端，自动支持连接复用。

### 并发控制

```go
import "golang.org/x/sync/semaphore"

// 限制并发数
sem := semaphore.NewWeighted(10)

for _, pkg := range packages {
    sem.Acquire(ctx, 1)
    go func(packageName string) {
        defer sem.Release(1)
        // 处理包
    }(pkg)
}
```

## 开发环境设置

### IDE 配置

**VS Code**:
安装 Go 扩展并配置：

```json
{
    "go.useLanguageServer": true,
    "go.alternateTools": {
        "go": "/usr/local/go/bin/go"
    }
}
```

**GoLand**:
确保配置了正确的 GOROOT 和 GOPATH。

### 调试配置

```go
import "log"

// 启用详细日志
log.SetFlags(log.LstdFlags | log.Lshortfile)

// 在关键点添加日志
pkg, err := client.GetPackageInformation(ctx, "react")
log.Printf("获取包信息: %s, 错误: %v", "react", err)
```

## 下一步

安装完成后，您可以：

1. 查看 [快速开始指南](/getting-started) 学习基本用法
2. 阅读 [API 文档](/api/) 了解所有功能
3. 浏览 [示例代码](/examples/basic) 学习实际应用
4. 参与 [GitHub 项目](https://github.com/scagogogo/npm-skills) 贡献代码

如果遇到问题，请在 [GitHub Issues](https://github.com/scagogogo/npm-skills/issues) 中报告。
