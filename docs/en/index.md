---
layout: home

hero:
  name: NPM Crawler
  text: High-performance NPM Registry Client
  tagline: Go NPM client library with multi-mirror source and proxy support
  image:
    src: https://cdn.worldvectorlogo.com/logos/npm-2.svg
    alt: NPM Logo
  actions:
    - theme: brand
      text: Getting Started
      link: /en/getting-started
    - theme: alt
      text: API Documentation
      link: /en/api/
    - theme: alt
      text: GitHub
      link: https://github.com/scagogogo/npm-skills

features:
  - icon: 🚀
    title: High Performance
    details: Based on Go's high concurrency features, providing fast NPM Registry access
  - icon: 🌐
    title: Multi-Mirror Support
    details: Built-in support for multiple NPM mirror sources, including official Registry, Taobao Mirror, Huawei Cloud Mirror, etc.
  - icon: 🔄
    title: Proxy Support
    details: Configurable HTTP proxy, adapting to various network environments
  - icon: 📦
    title: Complete Types
    details: Complete Go type definitions corresponding to various NPM package metadata
  - icon: 🧪
    title: Comprehensive Testing
    details: Complete unit test coverage ensuring code quality
  - icon: 📝
    title: Detailed Documentation
    details: Bilingual comments and documentation in Chinese and English, easy to use and integrate
---

## Installation

```bash
go get github.com/scagogogo/npm-skills
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

func main() {
    // Create default Registry client
    registry := registry.NewRegistry()
    ctx := context.Background()
    
    // Get package information
    pkg, err := registry.GetPackageInformation(ctx, "react")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Package Name: %s\n", pkg.Name)
    fmt.Printf("Latest Version: %s\n", pkg.DistTags["latest"])
}
```

## Supported Mirror Sources

| Mirror Source | URL | Region | Creation Method |
|---------------|-----|--------|-----------------|
| NPM Official | https://registry.npmjs.org | Global | `NewRegistry()` |
| Taobao NPM | https://registry.npm.taobao.org | China | `NewTaoBaoRegistry()` |
| NPM Mirror | https://registry.npmmirror.com | China | `NewNpmMirrorRegistry()` | 