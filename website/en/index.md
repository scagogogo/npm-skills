---
layout: home

hero:
  name: NPM Skills
  text: NPM Registry client for AI agents
  tagline: Query, publish, audit, mirrors, proxy in one · 70+ SDK methods · 33 MCP tools · 26 CLI commands
  image:
    src: /architecture.svg
    alt: NPM Skills architecture
  actions:
    - theme: brand
      text: Getting Started
      link: /en/getting-started
    - theme: alt
      text: CLI Reference
      link: /en/cli
    - theme: alt
      text: GitHub
      link: https://github.com/scagogogo/npm-skills

features:
  - icon: 🤖
    title: AI-Native First
    details: Installed as a Claude Code plugin, AI agents auto-discover and invoke it. SKILL.md uses progressive disclosure.
  - icon: 🚀
    title: High Performance
    details: Go-based, HTTP client cached via sync.Once for connection reuse, concurrent requests and streaming downloads, static CGO-free binaries.
  - icon: 🌐
    title: 8 Mirrors
    details: Built-in official, Taobao, Huawei Cloud, Tencent Cloud, CNPM, Yarn mirrors — no proxy needed in China.
  - icon: 🔄
    title: Proxy Support
    details: HTTP/HTTPS/SOCKS5 proxy for restricted networks; optional TLS skip for self-signed internal registries.
  - icon: 📦
    title: Full API Coverage
    details: 70+ SDK methods covering package query, versions, dist-tags, download stats, access, stars, tokens, webhooks, orgs, audit.
  - icon: 🛡️
    title: Typed Errors
    details: ErrNotFound / ErrUnauthorized / ErrRateLimited with errors.Is() support; sensitive fields masked in String().
  - icon: 📡
    title: MCP Protocol
    details: 33 MCP tools for any MCP-compatible client — Claude Code, Cursor, Windsurf.
  - icon: 🏗️
    title: 34-Platform Binaries
    details: GoReleaser covers Linux/macOS/Windows/FreeBSD/OpenBSD/NetBSD/Illumos/Solaris × 13 architectures, ready to run.
---
