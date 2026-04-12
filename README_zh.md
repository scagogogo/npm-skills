# NPM Crawler

<div align="center">

[切换到英文版](README.md)

<img src="https://cdn.worldvectorlogo.com/logos/npm-2.svg" width="180" alt="NPM Logo" style="filter: brightness(0.9);">

[![Go Tests](https://github.com/scagogogo/npm-crawler/actions/workflows/go-test.yml/badge.svg)](https://github.com/scagogogo/npm-crawler/actions/workflows/go-test.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/scagogogo/npm-crawler.svg)](https://pkg.go.dev/github.com/scagogogo/npm-crawler)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

_高性能的 NPM Registry 客户端，支持多镜像源和代理配置_

</div>

## 三种使用方式

NPM Crawler 以 **AI 原生 (AI-native) 为首要设计目标**，提供三种互补的交互方式：

### 1. 🤖 AI / Agent 模式（主要方式）

专为 AI 智能体和自动化工作流设计。AI 可直接理解和调用此工具的能力，无需人工干预。

```markdown
# 本仓库本身就是一个 Claude Code Skill
# 当你询问时，AI 会自动发现并使用它：
# - "查找 axios NPM 包的信息"
# - "下载 react 的 tarball"
# - "搜索 HTTP 客户端库"
# - "获取 vue 的下载统计"
```

**AI 触发词**: `npm package`, `NPM registry`, `search npm`, `download npm tarball`, `get npm stats`, `npm mirror`

Skill 清单 (`SKILL.md`) 提供渐进式披露：
- **即时上下文**: frontmatter 中的 name + description (~100 词)
- **核心指引**: 快速开始 + 能力说明 (~500 行)
- **深入参考**: 完整 API 文档在 `references/api.md` (按需加载)

### 2. 📦 Go SDK

用于程序化访问 NPM Registry 的嵌入式 Go 库：

```go
import "github.com/scagogogo/npm-crawler/pkg/registry"

client := registry.NewRegistry()
pkg, err := client.GetPackageInformation(ctx, "react")
```

### 3. 🖥️ CLI 工具

用于快速查询和脚本化的命令行界面：

```bash
cd examples/create_registry && go run main.go
cd examples/download_tarball && go run main.go
```

---

## 简介

NPM Crawler 是一个用 Go 语言编写的高性能 NPM Registry 客户端库，提供了简单易用的 API 来访问 NPM Registry 中的包信息。该库支持多种 NPM 镜像源，包括官方 Registry、淘宝镜像、华为云镜像等，同时支持代理配置，可以轻松应对各种网络环境。

## 功能特点

- 🤖 **AI 原生**: 作为 Skill 设计，支持 AI 智能体渐进式披露
- 🚀 **高性能**: 基于 Go 的高并发特性，提供快速的 NPM Registry 访问
- 🌐 **多镜像源支持**: 内置支持多种 NPM 镜像源
- 🔄 **代理支持**: 可配置 HTTP 代理，适应各种网络环境
- 📦 **完整类型**: 完整的 Go 类型定义，对应 NPM 包的各种元数据
- 🧪 **全面测试**: 完整的单元测试覆盖
- 📝 **详细文档**: 中英双语注释和文档

## 安装

```bash
go get github.com/scagogogo/npm-crawler
```

## 快速开始

### 基本使用

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/scagogogo/npm-crawler/pkg/registry"
)

func main() {
    // 创建默认 Registry 客户端 (使用官方 npmjs.org)
    client := registry.NewRegistry()
    
    // 或使用淘宝镜像
    // client := registry.NewTaoBaoRegistry()
    
    ctx := context.Background()
    
    // 获取包信息
    pkg, err := client.GetPackageInformation(ctx, "react")
    if err != nil {
        log.Fatalf("获取包信息失败: %v", err)
    }
    
    fmt.Printf("包名: %s\n", pkg.Name)
    // 输出: 包名: react
    
    fmt.Printf("描述: %s\n", pkg.Description)
    // 输出: 描述: React is a JavaScript library for building user interfaces.
    
    fmt.Printf("最新版本: %s\n", pkg.DistTags["latest"])
    // 输出: 最新版本: 18.2.0
    
    // 获取 Registry 信息
    info, err := client.GetRegistryInformation(ctx)
    if err != nil {
        log.Fatalf("获取 Registry 信息失败: %v", err)
    }
    
    fmt.Printf("Registry 名称: %s\n", info.DbName)
    // 输出: Registry 名称: registry
    
    fmt.Printf("包总数: %d\n", info.DocCount)
    // 输出: 包总数: 2400000
}
```

### 使用代理

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/scagogogo/npm-crawler/pkg/registry"
)

func main() {
    // 创建选项并配置代理
    options := registry.NewOptions().
        SetRegistryURL("https://registry.npmjs.org").
        SetProxy("http://your-proxy-server:8080")
    
    // 创建带代理的客户端
    client := registry.NewRegistry(options)
    
    ctx := context.Background()
    
    // 获取包信息
    pkg, err := client.GetPackageInformation(ctx, "react")
    if err != nil {
        log.Fatalf("获取包信息失败: %v", err)
    }
    
    fmt.Printf("包名: %s\n", pkg.Name)
    // 输出: 包名: react
    
    fmt.Printf("描述: %s\n", pkg.Description)
    // 输出: 描述: React is a JavaScript library for building user interfaces.
}
```

## API 文档

### Registry 相关

#### 创建 Registry 客户端

```go
// NewRegistry 创建一个新的 Registry 客户端实例
//
// 参数:
//   - options: 可选的配置选项，如未提供则使用默认配置
//
// 返回值:
//   - *Registry: 新创建的 Registry 客户端实例
func NewRegistry(options ...*Options) *Registry
```

#### 创建特定镜像源的客户端

```go
// 创建使用淘宝 NPM 镜像源的 Registry 客户端
func NewTaoBaoRegistry() *Registry

// 创建使用 NPM Mirror 镜像源的 Registry 客户端 (原淘宝镜像新域名)
func NewNpmMirrorRegistry() *Registry

// 创建使用华为云镜像源的 Registry 客户端
func NewHuaWeiCloudRegistry() *Registry

// 创建使用腾讯云镜像源的 Registry 客户端
func NewTencentRegistry() *Registry

// 创建使用 CNPM 镜像源的 Registry 客户端
func NewCnpmRegistry() *Registry

// 创建使用 Yarn 官方镜像源的 Registry 客户端
func NewYarnRegistry() *Registry

// 创建使用 npmjs.com 镜像源的 Registry 客户端
func NewNpmjsComRegistry() *Registry
```

#### 获取 Registry 信息

```go
// GetRegistryInformation 获取 NPM Registry 的状态信息
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//
// 返回值:
//   - *models.RegistryInformation: Registry 状态信息
//   - error: 如果请求失败则返回错误
func (x *Registry) GetRegistryInformation(ctx context.Context) (*models.RegistryInformation, error)
```

#### 获取包信息

```go
// GetPackageInformation 获取指定 NPM 包的详细信息
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - packageName: 要查询的包名称，例如 "react"、"lodash" 等
//
// 返回值:
//   - *models.Package: 包的详细信息
//   - error: 如果请求失败则返回错误
func (x *Registry) GetPackageInformation(ctx context.Context, packageName string) (*models.Package, error)
```

### 配置选项相关

#### 创建选项

```go
// NewOptions 创建并返回一个新的默认配置选项实例
//
// 默认配置:
//   - RegistryURL: "https://registry.npmjs.org"
//   - Proxy: 无代理设置
func NewOptions() *Options
```

#### 设置 Registry URL

```go
// SetRegistryURL 设置 NPM 仓库服务器的 URL 地址
//
// 参数:
//   - url: 一个有效的 NPM 仓库 URL 地址字符串
//
// 返回值:
//   - *Options: 更新后的选项对象 (支持链式调用)
func (o *Options) SetRegistryURL(url string) *Options
```

#### 设置代理

```go
// SetProxy 设置 HTTP 代理服务器的 URL 地址
//
// 参数:
//   - proxyUrl: HTTP 代理服务器的 URL 地址字符串
//
// 返回值:
//   - *Options: 更新后的选项对象 (支持链式调用)
func (o *Options) SetProxy(proxyUrl string) *Options
```

### 主要模型

#### Package

表示一个 NPM 包的完整信息结构：

```go
type Package struct {
    ID             string                 `json:"_id"`            // 包 ID
    Rev            string                 `json:"_rev"`           // 修订号
    Name           string                 `json:"name"`           // 包名称
    Description    string                 `json:"description"`    // 包描述
    DistTags       map[string]string      `json:"dist-tags"`      // 发布标签，如 latest
    Versions       map[string]Version     `json:"versions"`       // 版本信息映射
    Maintainers    []Maintainer           `json:"maintainers"`    // 维护者列表
    Time           map[string]string      `json:"time"`           // 时间信息
    Repository     Repository             `json:"repository"`     // 代码仓库信息
    ReadMe         string                 `json:"readme"`         // README 内容
    ReadMeFilename string                 `json:"readmeFilename"` // README 文件名
    Homepage       string                 `json:"homepage"`       // 项目主页
    Bugs           map[string]interface{} `json:"bugs"`           // 问题追踪信息
    License        string                 `json:"license"`        // 许可证
    Users          map[string]bool        `json:"users"`          // 用户信息
    Keywords       []string               `json:"keywords"`       // 关键词列表
    Author         Author                 `json:"author"`         // 作者信息
    Contributors   []Contributor          `json:"contributors"`   // 贡献者列表
    Deprecated     string                 `json:"deprecated"`     // 弃用说明
    Other          map[string]interface{} `json:"other"`          // 其他字段
}
```

#### Version

表示 NPM 包的特定版本信息：

```go
type Version struct {
    Name            string               `json:"name"`            // 包名称
    Version         string               `json:"version"`         // 版本号
    Description     string               `json:"description"`     // 版本描述
    Main            string               `json:"main"`            // 主入口文件
    Scripts         *Script              `json:"scripts"`         // 脚本命令
    Repository      *Repository          `json:"repository"`      // 代码仓库
    Keywords        []string             `json:"keywords"`        // 关键词列表
    Author          *User                `json:"author"`          // 作者信息
    License         string               `json:"license"`         // 许可证
    Bugs            *Bugs                `json:"bugs"`            // 问题追踪
    Homepage        string               `json:"homepage"`        // 项目主页
    Dependencies    map[string]string    `json:"dependencies"`    // 运行时依赖
    DevDependencies map[string]string    `json:"devDependencies"` // 开发依赖
    Dist            *Dist                `json:"dist"`            // 分发信息
    // 其他字段...
}
```

#### RegistryInformation

表示 NPM Registry 的状态信息：

```go
type RegistryInformation struct {
    DbName            string `json:"db_name"`              // 数据库名称
    DocCount          int    `json:"doc_count"`            // 文档(包)总数
    DocDelCount       int    `json:"doc_del_count"`        // 已删除的文档数
    UpdateSeq         int    `json:"update_seq"`           // 更新序列号
    PurgeSeq          int    `json:"purge_seq"`            // 清除序列号
    CompactRunning    bool   `json:"compact_running"`      // 是否正在压缩
    DiskSize          int64  `json:"disk_size"`            // 磁盘占用大小
    DataSize          int64  `json:"data_size"`            // 数据大小
    InstanceStartTime string `json:"instance_start_time"`  // 实例启动时间
    // 其他字段...
}
```

## 支持的镜像源

| 镜像源 | URL | 地域 | 创建方法 |
|-------|-----|------|---------|
| NPM 官方 | https://registry.npmjs.org | 全球 | `NewRegistry()` |
| 淘宝 NPM | https://registry.npm.taobao.org | 中国 | `NewTaoBaoRegistry()` |
| NPM Mirror | https://registry.npmmirror.com | 中国 | `NewNpmMirrorRegistry()` |
| 华为云 | https://mirrors.huaweicloud.com/repository/npm | 中国 | `NewHuaWeiCloudRegistry()` |
| 腾讯云 | http://mirrors.cloud.tencent.com/npm | 中国 | `NewTencentRegistry()` |
| CNPM | http://r.cnpmjs.org | 中国 | `NewCnpmRegistry()` |
| Yarn | https://registry.yarnpkg.com | 全球 | `NewYarnRegistry()` |
| NPM CouchDB | https://skimdb.npmjs.com | 全球 | `NewNpmjsComRegistry()` |

## 贡献指南

欢迎贡献代码！请遵循以下步骤：

1. Fork 本仓库
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 许可证

本项目采用 MIT 许可证 - 详情请参阅 [LICENSE](LICENSE) 文件。

## 致谢

- [NPM Registry](https://registry.npmjs.org) - 提供 API 和数据
- [Go Requests](https://github.com/crawler-go-go-go/go-requests) - HTTP 客户端库