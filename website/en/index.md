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

## Four entry points, one core

NPM Skills exposes the same Registry capabilities through four entry points — pick whichever fits your scenario. They all share a single connection-pooled Go client underneath:

```mermaid
flowchart TB
    subgraph Entry["Entry points"]
        AI["🤖 Claude Code plugin<br/>SKILL.md progressive disclosure"]
        MCP["📡 MCP Server<br/>33 tools · JSON-RPC"]
        CLI["⌨️ CLI<br/>26 commands"]
        SDK["📦 Go SDK<br/>70+ methods"]
    end

    Core["pkg/registry<br/>Registry client"]
    HTTP["http.Client<br/>sync.Once pooled reuse"]
    NPM["NPM Registry / 8 mirrors"]

    AI --> MCP
    MCP --> Core
    CLI --> Core
    SDK --> Core
    Core --> HTTP
    HTTP --> NPM

    classDef core fill:#cb3837,stroke:#8b0000,color:#fff;
    class Core core;
```

## 30-second start

::: code-group

```bash [CLI]
# Query package info
npm-skills info react

# Search and sort by relevance
npm-skills search "web framework" --size 10

# Download a tarball via a mirror
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
    fmt.Println("Latest version:", pkg.DistTags["latest"])
}
```

:::

Ready to go deeper? Head to [Getting Started](/en/getting-started) for the full setup of all four entry points, or jump straight to the [API docs](/en/api/) and [CLI reference](/en/cli).
