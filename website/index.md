---
layout: home

hero:
  name: NPM Skills
  text: 面向 AI 智能体的 NPM 客户端
  tagline: 查询、发布、审计、镜像、代理一体化 · 70+ SDK 方法 · 33 MCP 工具 · 26 CLI 命令
  image:
    src: /architecture.svg
    alt: NPM Skills 架构图
  actions:
    - theme: brand
      text: 快速开始
      link: /getting-started
    - theme: alt
      text: CLI 命令手册
      link: /cli
    - theme: alt
      text: GitHub
      link: https://github.com/scagogogo/npm-skills

features:
  - icon: 🤖
    title: AI 原生优先
    details: 作为 Claude Code 插件安装后，AI 智能体自动发现并调用，无需手动 Shell。SKILL.md 采用渐进式披露。
  - icon: 🚀
    title: 高性能
    details: 基于 Go，HTTP 客户端 sync.Once 缓存复用连接池，并发请求与流式下载，CGO_ENABLED=0 纯静态二进制。
  - icon: 🌐
    title: 8 镜像源
    details: 内置官方、淘宝、华为云、腾讯云、CNPM、Yarn 等镜像，国内访问无需代理。
  - icon: 🔄
    title: 代理支持
    details: HTTP/HTTPS/SOCKS5 代理，适配受限网络环境；可跳过 TLS 验证用于内网自签名证书。
  - icon: 📦
    title: 全 API 覆盖
    details: 70+ SDK 方法覆盖包查询、版本、dist-tags、下载统计、访问控制、stars、tokens、webhooks、orgs、审计。
  - icon: 🛡️
    title: 类型化错误
    details: ErrNotFound / ErrUnauthorized / ErrRateLimited 等，支持 errors.Is() 程序化处理，敏感字段脱敏。
  - icon: 📡
    title: MCP 协议
    details: 33 个 MCP 工具，供 Claude Code、Cursor、Windsurf 等任意 MCP 客户端调用。
  - icon: 🏗️
    title: 34 平台预编译
    details: GoReleaser 覆盖 Linux/macOS/Windows/FreeBSD/OpenBSD/NetBSD/Illumos/Solaris × 13 架构，开箱即用。
---

## 四种使用方式，一套核心

NPM Skills 把同一份 Registry 能力暴露为四种入口，按你的场景任选其一——底层共享同一个经过连接池优化的 Go 客户端：

```mermaid
flowchart TB
    subgraph Entry["接入方式"]
        AI["🤖 Claude Code 插件<br/>SKILL.md 渐进披露"]
        MCP["📡 MCP Server<br/>33 工具 · JSON-RPC"]
        CLI["⌨️ CLI<br/>26 命令"]
        SDK["📦 Go SDK<br/>70+ 方法"]
    end

    Core["pkg/registry<br/>Registry 客户端"]
    HTTP["http.Client<br/>sync.Once 连接池复用"]
    NPM["NPM Registry / 8 镜像源"]

    AI --> MCP
    MCP --> Core
    CLI --> Core
    SDK --> Core
    Core --> HTTP
    HTTP --> NPM

    classDef core fill:#cb3837,stroke:#8b0000,color:#fff;
    class Core core;
```

## 30 秒上手

::: code-group

```bash [CLI]
# 查询包信息
npm-skills info react

# 搜索并按下载量排序
npm-skills search "web framework" --size 10

# 通过淘宝镜像下载 tarball（国内免代理）
npm-skills download lodash 4.17.21 ./ --registry https://registry.npmmirror.com
```

```go [Go SDK]
package main

import (
    "context"
    "fmt"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

func main() {
    client := registry.NewRegistry()
    pkg, err := client.GetPackageInformation(context.Background(), "react")
    if err != nil {
        panic(err)
    }
    fmt.Println("最新版本:", pkg.DistTags["latest"])
}
```

:::

准备好深入了？前往 [快速开始](/getting-started) 了解四种接入方式的完整配置，或直接查阅 [API 文档](/api/) 与 [CLI 命令手册](/cli)。
