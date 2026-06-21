<div align="center">

# NPM Skills

[Switch to English](README.md)

<img src="https://cdn.worldvectorlogo.com/logos/npm-2.svg" width="180" alt="NPM Logo" style="filter: brightness(0.9);">

[![Go Tests](https://github.com/scagogogo/npm-skills/actions/workflows/go-test.yml/badge.svg)](https://github.com/scagogogo/npm-skills/actions/workflows/go-test.yml)
[![Release](https://github.com/scagogogo/npm-skills/actions/workflows/release.yml/badge.svg)](https://github.com/scagogogo/npm-skills/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/scagogogo/npm-skills.svg)](https://pkg.go.dev/github.com/scagogogo/npm-skills)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**面向 AI 智能体和开发者的 NPM Registry 客户端** — 查询包信息、管理 dist-tags、发布、审计等。

[四种接入方式](#-四种接入方式) · [下载](https://github.com/scagogogo/npm-skills/releases/latest) · [文档](https://pkg.go.dev/github.com/scagogogo/npm-skills)

</div>

---

## ⚡ 一键安装 Plugin

安装为 Claude Code 插件 — AI 智能体会自动发现并使用它：

```bash
# 第 1 步：添加 marketplace
claude plugin marketplace add scagogogo/npm-skills

# 第 2 步：安装插件
claude plugin install npm@npm-skills
```

安装后，直接用自然语言向 Claude Code 提问即可：
- *"查找 axios NPM 包的信息"*
- *"下载 react 的 tarball"*
- *"搜索 HTTP 客户端库"*
- *"获取 vue 上个月的下载统计"*
- *"用国内镜像查看 NPM 注册表"*
- *"发布包到私有仓库"*
- *"审计我的依赖漏洞"*

---

## 🔌 四种接入方式

NPM Skills 以 **AI 原生优先**为设计理念，提供四种互补的方式与 NPM Registry 交互：

### 1. 🤖 Skill（AI 智能体）— 主要方式

本仓库是一个 **Claude Code 插件** — 安装后 AI 智能体会自动发现并使用它，无需手动调用 Shell。

**安装：**
```bash
claude plugin marketplace add scagogogo/npm-skills
claude plugin install npm@npm-skills
```

安装后，直接用自然语言向 Claude Code 提问即可：
- *"查找 axios NPM 包的信息"*
- *"下载 react 的 tarball"*
- *"搜索 HTTP 客户端库"*
- *"获取 vue 上个月的下载统计"*
- *"用国内镜像查看 NPM 注册表"*
- *"发布包到私有仓库"*
- *"审计我的依赖漏洞"*

**触发词**：`npm package`、`npm publish`、`NPM registry`、`search npm`、`npm stats`、`npm mirror`、`npm 版本`、`npm 包`、`npm 镜像`、`npm 发布`

Skill 清单 (`SKILL.md`) 采用渐进式披露设计：
- **即时上下文**：frontmatter 中的 name + description（约 100 词）
- **核心指引**：CLI 命令 + 使用模式
- **深入参考**：完整 API 文档在 `references/api.md`（按需加载）

### 2. 📦 Go SDK（开发者）

即插即用的 Go 库，提供完整类型安全：

```go
import "github.com/scagogogo/npm-skills/pkg/registry"

// 默认客户端（官方 Registry）
client := registry.NewRegistry()

// 自定义客户端
options := registry.NewOptions().
    SetRegistryURL("https://registry.npmjs.org").
    SetToken("npm_xxxxx").
    SetProxy("http://proxy:8080").
    SetTimeout(30 * time.Second)
client = registry.NewRegistry(options)

// 读取操作
pkg, _ := client.GetPackageInformation(ctx, "react")
versions, _ := client.GetPackageVersions(ctx, "react")
stats, _ := client.GetDownloadStats(ctx, "react", "last-week")
rangeStats, _ := client.GetDownloadRangeStatsByDateRange(ctx, "react", "2024-01-01", "2024-06-30")

// 写入操作（需要 token）
client.SetDistTag(ctx, "my-pkg", "next", "2.0.0-rc.1")
client.PublishPackage(ctx, pkg)
client.DeprecateVersion(ctx, "my-pkg", "1.0.0", "Use v2.0.0")

// 类型化错误，方便程序化处理
import "errors"
_, err := client.GetPackageInformation(ctx, "nonexistent")
if errors.Is(err, registry.ErrNotFound) {
    // 处理 404
}
```

### 3. 🖥️ CLI 工具

命令行界面，支持彩色输出、代理和镜像源。[所有主流平台](https://github.com/scagogogo/npm-skills/releases/latest)均有预编译二进制文件。

**安装：**
```bash
# 从 GitHub Release 下载（推荐）
# 参见：https://github.com/scagogogo/npm-skills/releases/latest

# 或从源码构建
bash scripts/install.sh

# 或 go install
go install github.com/scagogogo/npm-skills/cmd/npm-skills@latest
```

**使用：**
```bash
# 读取操作
npm-skills package-summary react            # 轻量包信息（推荐）
npm-skills package react                    # 完整包元数据
npm-skills search "http client" -l 10       # 搜索包
npm-skills versions react --latest          # 获取最新版本
npm-skills dist-tags get react              # 获取 dist-tags
npm-skills download-stats axios -p last-month  # 下载统计
npm-skills download lodash 4.17.21 ./lodash.tgz  # 下载 tarball
npm-skills mirrors                          # 列出镜像源
npm-skills whoami --token npm_xxxxx         # 检查认证状态

# 写入操作（需要 --token）
npm-skills publish ./pkg.tgz --name my-pkg --version 1.0.0 -t npm_xxxxx
npm-skills deprecate my-pkg 1.0.0 -M "Use v2" -t npm_xxxxx
npm-skills dist-tags set my-pkg stable --version 1.0.0 -t npm_xxxxx
npm-skills access get my-pkg -t npm_xxxxx
npm-skills star add react -t npm_xxxxx
npm-skills audit quick --deps "lodash=4.17.11"

# 镜像 & 代理 & 私有仓库
npm-skills package react -m npm-mirror                                    # 国内镜像
npm-skills package react --proxy http://127.0.0.1:7890                    # HTTP 代理
npm-skills package my-lib --registry https://npm.my-company.com -t npm_x  # 私有仓库

# 环境变量
export NPM_MIRROR=npm-mirror
export NPM_PROXY=http://127.0.0.1:7890
export NPM_REGISTRY=https://npm.company.com
npm-skills package react    # 自动使用环境变量

npm-skills --help           # 显示所有 26 个命令
```

### 4. 📡 MCP 服务器（AI 工具链）

MCP（Model Context Protocol）服务器，将 NPM Registry 操作暴露为工具，供任何 MCP 兼容的 AI 客户端使用 — Claude Code、Cursor、Windsurf 等。

**安装：**
```bash
bash scripts/install.sh   # 同时构建 CLI 和 MCP 服务器
```

**配置（Claude Code / Cursor / 任何 MCP 客户端）：**
```json
{
  "mcpServers": {
    "npm-registry": {
      "command": "npm-mcp-server",
      "args": ["--mirror", "npm-mirror"]
    }
  }
}
```

**33 个 MCP 工具**可用，包括：

| 读取工具 | 写入工具 |
|---|---|
| `npm_registry_info`、`npm_mirrors`、`npm_package`、`npm_package_summary`、`npm_search`、`npm_version`、`npm_versions`、`npm_latest_version`、`npm_dist_tags`、`npm_download_stats`、`npm_download_range`、`npm_whoami` | `npm_dist_tag_set`、`npm_dist_tag_delete`、`npm_dist_tags_set`、`npm_star`、`npm_unstar`、`npm_stargazers`、`npm_access_get`、`npm_collaborators`、`npm_token_list`、`npm_audit_quick`、`npm_audit_advisory`、`npm_hook_list`、`npm_hook_get`、`npm_org_get`、`npm_org_members`、`npm_org_packages`、`npm_team_list`、`npm_team_members`、`npm_team_packages`、`npm_changes` |

---

## ✨ 功能特点

- 🤖 **AI 原生优先**：以 Skill 设计，支持 AI 智能体渐进式披露
- 🚀 **高性能**：基于 Go 的并发请求和流式下载
- 🌐 **8 个镜像源**：内置官方、国内和全球镜像支持
- 🔄 **代理支持**：HTTP 代理配置，适应受限网络环境
- 📦 **完整 API 覆盖**：70+ SDK 方法，覆盖所有主要 NPM Registry 端点
- 🛡️ **类型化错误**：`ErrNotFound`、`ErrUnauthorized`、`ErrRateLimited` 等，支持 `errors.Is()`
- ⏱️ **超时控制**：通过 `Options.SetTimeout()` 设置客户端级超时
- 🔒 **认证支持**：Bearer Token 支持发布、取消发布和所有写操作
- 📊 **下载分析**：点统计、区间统计、批量统计（>128 个包自动分块）
- 🔍 **包搜索**：分页、质量/流行度/维护度评分
- 📡 **MCP 协议**：33 个工具供 AI 工具链使用
- 🏗️ **跨平台**：Linux、macOS、Windows、FreeBSD、OpenBSD、NetBSD、Illumos、Solaris 预编译二进制

## 📥 安装

### 下载二进制文件（推荐）

[最新 Release](https://github.com/scagogogo/npm-skills/releases/latest) 提供预编译二进制文件：

```bash
# Linux (x86_64)
curl -sL https://github.com/scagogogo/npm-skills/releases/latest/download/npm-skills_0.2.0_linux_x86_64.tar.gz | tar -xz
sudo mv npm-skills npm-mcp-server /usr/local/bin/

# macOS (Apple Silicon)
curl -sL https://github.com/scagogogo/npm-skills/releases/latest/download/npm-skills_0.2.0_aarch64.tar.gz | tar -xz
sudo mv npm-skills npm-mcp-server /usr/local/bin/

# Windows — 从 releases 页面下载 .zip 文件
```

### Go Install

```bash
go install github.com/scagogogo/npm-skills/cmd/npm-skills@latest
go install github.com/scagogogo/npm-skills/cmd/mcp-server@latest
```

### Go Module

```bash
go get github.com/scagogogo/npm-skills
```

## 🚀 快速开始

### Go SDK

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
    // 创建客户端（默认使用官方 Registry）
    client := registry.NewRegistry()
    ctx := context.Background()

    // 获取轻量包信息
    pkg, err := client.GetAbbreviatedPackageInformation(ctx, "react")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("包名: %s, 最新版本: %s\n", pkg.Name, pkg.DistTags["latest"])

    // 搜索包
    results, err := client.SearchPackages(ctx, "http client", 5)
    if err != nil {
        log.Fatal(err)
    }
    for _, obj := range results.Objects {
        fmt.Printf("  %s — %s\n", obj.Package.Name, obj.Package.Description)
    }

    // 下载统计
    stats, err := client.GetDownloadStats(ctx, "react", "last-week")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("React 上周下载量: %d\n", stats.Downloads)

    // 自定义仓库（带认证和超时）
    options := registry.NewOptions().
        SetRegistryURL("https://npm.my-company.com").
        SetToken("npm_xxxxx").
        SetTimeout(30 * time.Second)
    privateClient := registry.NewRegistry(options)
    _ = privateClient
}
```

## 🪞 支持的镜像源

| 镜像源 | URL | 地域 | SDK 方法 |
|-------|-----|------|---------|
| NPM 官方 | `https://registry.npmjs.org` | 全球 | `NewRegistry()` |
| NPM Mirror | `https://registry.npmmirror.com` | 中国 | `NewNpmMirrorRegistry()` |
| 淘宝 | `https://registry.npm.taobao.org` | 中国 | `NewTaoBaoRegistry()` |
| 华为云 | `https://mirrors.huaweicloud.com/repository/npm` | 中国 | `NewHuaWeiCloudRegistry()` |
| 腾讯云 | `http://mirrors.cloud.tencent.com/npm` | 中国 | `NewTencentRegistry()` |
| CNPM | `http://r.cnpmjs.org` | 中国 | `NewCnpmRegistry()` |
| Yarn | `https://registry.yarnpkg.com` | 全球 | `NewYarnRegistry()` |
| NPM CouchDB | `https://skimdb.npmjs.com` | 全球 | `NewNpmjsComRegistry()` |

## 🤝 贡献指南

欢迎贡献代码！请遵循以下步骤：

1. Fork 本仓库
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 — 详情请参阅 [LICENSE](LICENSE) 文件。

## 🙏 致谢

- [NPM Registry](https://registry.npmjs.org) — 提供 API 和数据
- [Go Requests](https://github.com/crawler-go-go-go/go-requests) — HTTP 客户端库
- [Cobra](https://github.com/spf13/cobra) — CLI 框架
- [MCP-Go](https://github.com/mark3labs/mcp-go) — MCP 服务器框架
