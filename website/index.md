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
