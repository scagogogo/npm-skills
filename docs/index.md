---
layout: home

hero:
  name: NPM Crawler
  text: 高性能 NPM Registry 客户端
  tagline: 支持多镜像源和代理配置的Go语言NPM客户端库
  image:
    src: https://cdn.worldvectorlogo.com/logos/npm-2.svg
    alt: NPM Logo
  actions:
    - theme: brand
      text: 快速开始
      link: /getting-started
    - theme: alt
      text: API文档
      link: /api/
    - theme: alt
      text: GitHub
      link: https://github.com/scagogogo/npm-skills

features:
  - icon: 🚀
    title: 高性能
    details: 基于Go的高并发特性，提供快速的NPM Registry访问
  - icon: 🌐
    title: 多镜像源支持
    details: 内置支持多种NPM镜像源，包括官方Registry、淘宝镜像、华为云镜像等
  - icon: 🔄
    title: 代理支持
    details: 可配置HTTP代理，适应各种网络环境
  - icon: 📦
    title: 完整类型
    details: 完整的Go类型定义，对应NPM包的各种元数据
  - icon: 🧪
    title: 全面测试
    details: 完整的单元测试覆盖，确保代码质量
  - icon: 📝
    title: 详细文档
    details: 中英双语注释和文档，易于使用和集成
---

## 安装

```bash
go get github.com/scagogogo/npm-skills
```

## 快速开始

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

func main() {
    // 创建默认Registry客户端
    registry := registry.NewRegistry()
    ctx := context.Background()
    
    // 获取包信息
    pkg, err := registry.GetPackageInformation(ctx, "react")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("包名: %s\n", pkg.Name)
    fmt.Printf("最新版本: %s\n", pkg.DistTags["latest"])
}
```

## 支持的镜像源

| 镜像源 | URL | 地区 | 创建方法 |
|--------|-----|------|----------|
| NPM 官方 | https://registry.npmjs.org | 全球 | `NewRegistry()` |
| 淘宝 NPM | https://registry.npm.taobao.org | 中国 | `NewTaoBaoRegistry()` |
| NPM Mirror | https://registry.npmmirror.com | 中国 | `NewNpmMirrorRegistry()` |
| 华为云 | https://repo.huaweicloud.com/repository/npm | 中国 | `NewHuaWeiCloudRegistry()` |
| 腾讯云 | https://mirrors.cloud.tencent.com/npm | 中国 | `NewTencentRegistry()` |

更多镜像源配置请参考 [镜像源配置指南](/examples/mirrors)。

## 为什么选择 NPM Crawler？

- **简单易用**: 提供简洁的API接口，快速集成到您的项目中
- **高性能**: 基于Go语言的高并发特性，处理大量请求时表现优异
- **灵活配置**: 支持多种镜像源和代理配置，适应不同的网络环境
- **类型安全**: 完整的Go类型定义，减少运行时错误
- **生产就绪**: 经过充分测试，可直接用于生产环境
