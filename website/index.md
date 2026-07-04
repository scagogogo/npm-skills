---
layout: home

hero:
  name: NPM Skills
  text: 面向 AI 智能体的 NPM 客户端
  tagline: 查询、发布、审计、镜像、代理一体化 · 70+ SDK 方法 · 31 MCP 工具 · 26 CLI 命令
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
    details: 31 个 MCP 工具，供 Claude Code、Cursor、Windsurf 等任意 MCP 客户端调用。
  - icon: 🏗️
    title: 34 平台预编译
    details: GoReleaser 覆盖 Linux/macOS/Windows/FreeBSD/OpenBSD/NetBSD/Illumos/Solaris × 13 架构，开箱即用。
---

## 📋 一键复制给 AI Agent

把下面这段提示词复制粘贴给你的 AI Agent（Claude Code、Cursor、Windsurf、ChatGPT 等），Agent 会顺着指引自动安装并使用 NPM Skills，无需你逐条解释命令：

::: tip 复制方法
点击代码块右上角的 **复制** 按钮，粘贴到你的 AI Agent 对话框即可。Agent 会按提示词完成安装、验证并开始执行 NPM 相关任务。
:::

````md
---
title: 使用 NPM Skills 操作 NPM Registry
---

# 你的任务

你现在可以使用 **NPM Skills** —— 一个完整的 NPM Registry 客户端，支持包查询、版本管理、dist-tags、下载统计、发布、访问控制、安全审计、组织/团队、webhooks 等，内置 8 个镜像源（中国大陆免代理）与 HTTP/SOCKS5 代理。

## 第一步：安装（按你的环境任选一种，推荐方式 1）

### 方式 1：作为 Claude Code 插件（AI 原生，推荐）
```bash
claude plugin marketplace add scagogogo/npm-skills
claude plugin install npm@npm-skills
```
安装后你（AI）会自动发现 SKILL.md 并按需调用，无需手动 shell。

### 方式 2：下载预编译二进制（无 Claude Code 时）
```bash
# Linux x86_64
curl -sL https://github.com/scagogogo/npm-skills/releases/latest/download/npm-skills_0.2.0_linux_x86_64.tar.gz | tar -xz
sudo mv npm-skills npm-mcp-server /usr/local/bin/
# macOS Apple Silicon：把 linux_x86_64 换成 aarch64；Windows 用 zip 包
```

### 方式 3：从源码构建
```bash
bash scripts/install.sh   # 编译到 ~/.local/bin/
```

### 方式 4：go install
```bash
go install github.com/scagogogo/npm-skills/cmd/npm-skills@latest
go install github.com/scagogogo/npm-skills/cmd/mcp-server@latest
```

## 第二步：验证安装
```bash
npm-skills --version
npm-skills mirrors          # 列出 8 个镜像源
```

## 第三步：常用命令速查（90% 场景）

### 查包信息
```bash
npm-skills package-summary <name>        # 轻量（推荐，KB 级响应）
npm-skills package <name>                # 完整元数据（可能 10MB+）
npm-skills versions <name> --latest      # 最新版本
npm-skills pkg-version <name> <ver>      # 特定版本详情
```

### 搜索与统计
```bash
npm-skills search "<query>" -l 10                  # 搜索包
npm-skills download-stats <name> -p last-month    # 下载量
npm-skills download-range <name> -p last-week     # 每日趋势
npm-skills download-stats-bulk react,vue -p last-month  # 批量对比
```

### 镜像与代理（中国大陆用户重点）
```bash
npm-skills package react -m npm-mirror                              # 国内镜像，免代理
npm-skills package react --proxy http://127.0.0.1:7890              # HTTP 代理
npm-skills package my-lib --registry https://npm.my-company.com     # 私有仓库
```
环境变量：`NPM_MIRROR=npm-mirror`、`NPM_PROXY=...`、`NPM_TOKEN=...`

### 写操作（需要 token）
```bash
npm-skills publish ./pkg.tgz --name my-pkg --version 1.0.0 -t npm_xxxxx
npm-skills deprecate my-pkg 1.0.0 -M "Use v2.0.0" -t npm_xxxxx
npm-skills dist-tags set my-pkg next --version 2.0.0-rc.1 -t npm_xxxxx
npm-skills access get my-pkg -t npm_xxxxx
npm-skills audit quick --deps "lodash=4.17.11,express=4.17.1"
```

## 关键规则
1. **所有命令输出 JSON 到 stdout**，状态信息走 stderr —— 你可以用 `jq` 解析。
2. **读取操作免认证**；写入操作必须 `--token` 或 `NPM_TOKEN`。用 `whoami` 验证。
3. **优先 `package-summary`** 而非 `package`，响应小、快。
4. **下载统计始终查 api.npmjs.org**，与镜像设置无关。
5. 中国大陆用户默认 `-m npm-mirror`；受限网络叠 `--proxy`；私有仓库用 `--registry`。
6. 类型化错误：`ErrNotFound`/`ErrUnauthorized`/`ErrRateLimited` 等，可程序化处理。

## 当用户提出 NPM 相关需求时
直接选对应命令执行，例如：
- "查 axios 信息" → `npm-skills package-summary axios`
- "react 上个月下载量" → `npm-skills download-stats react -p last-month`
- "国内镜像查 vue" → `npm-skills package vue -m npm-mirror`
- "发布我的包" → `npm-skills publish ... -t $NPM_TOKEN`
- "审计依赖安全" → `npm-skills audit quick --deps "..."`

完整 26 命令见 `npm-skills --help`，SDK 文档见 https://scagogogo.github.io/npm-skills/api/
````

## 四种使用方式，一套核心

NPM Skills 把同一份 Registry 能力暴露为四种入口，按你的场景任选其一——底层共享同一个经过连接池优化的 Go 客户端：

```mermaid
flowchart TB
    subgraph Entry["接入方式"]
        AI["🤖 Claude Code 插件<br/>SKILL.md 渐进披露"]
        MCP["📡 MCP Server<br/>31 工具 · JSON-RPC"]
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
npm-skills package-summary react

# 搜索并限制结果数
npm-skills search "web framework" --limit 10

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
