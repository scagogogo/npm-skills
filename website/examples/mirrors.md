# 镜像源配置示例

本页面展示如何在不同环境和地区中配置最佳的 NPM 镜像源。

## 镜像源对比

| 镜像源 | URL | 地区 | 特点 | 推荐场景 |
|--------|-----|------|------|----------|
| NPM 官方 | https://registry.npmjs.org | 全球 | 数据最权威，更新最及时 | 海外服务器，需要最新数据 |
| NPM Mirror | https://registry.npmmirror.com | 中国 | 淘宝镜像新域名，速度快 | 中国大陆用户首选 |
| 淘宝镜像 | https://registry.npm.taobao.org | 中国 | 老牌镜像，稳定可靠 | 兼容性要求高的场景 |
| 华为云 | https://mirrors.huaweicloud.com/repository/npm | 中国 | 企业级稳定性 | 华为云用户 |
| 腾讯云 | http://mirrors.cloud.tencent.com/npm | 中国 | 腾讯云优化 | 腾讯云用户 |
| CNPM | http://r.cnpmjs.org | 中国 | 社区维护 | 开发测试环境 |
| Yarn | https://registry.yarnpkg.com | 全球 | Yarn 生态集成 | Yarn 用户 |

## 示例 1: 自动镜像选择

```go
package main

import (
    "context"
    "fmt"
    "net"
    "time"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

// 镜像配置
type MirrorConfig struct {
    Name string
    URL  string
    Client *registry.Registry
}

// 智能镜像选择器
type SmartMirrorSelector struct {
    mirrors []MirrorConfig
}

func NewSmartMirrorSelector() *SmartMirrorSelector {
    mirrors := []MirrorConfig{
        {
            Name: "NPM官方",
            URL:  "https://registry.npmjs.org",
            Client: registry.NewRegistry(),
        },
        {
            Name: "NPM Mirror",
            URL:  "https://registry.npmmirror.com",
            Client: registry.NewNpmMirrorRegistry(),
        },
        {
            Name: "淘宝镜像",
            URL:  "https://registry.npm.taobao.org",
            Client: registry.NewTaoBaoRegistry(),
        },
        {
            Name: "华为云",
            URL:  "https://mirrors.huaweicloud.com/repository/npm",
            Client: registry.NewHuaWeiCloudRegistry(),
        },
    }
    
    return &SmartMirrorSelector{mirrors: mirrors}
}

// 测试镜像延迟
func (sms *SmartMirrorSelector) testMirrorLatency(ctx context.Context, mirror MirrorConfig) (time.Duration, error) {
    start := time.Now()
    
    // 使用简单的包查询来测试延迟
    _, err := mirror.Client.GetPackageInformation(ctx, "lodash")
    
    latency := time.Since(start)
    return latency, err
}

// 选择最快的镜像
func (sms *SmartMirrorSelector) SelectBestMirror(ctx context.Context) MirrorConfig {
    type result struct {
        mirror  MirrorConfig
        latency time.Duration
        err     error
    }
    
    results := make(chan result, len(sms.mirrors))
    
    // 并发测试所有镜像
    for _, mirror := range sms.mirrors {
        go func(m MirrorConfig) {
            latency, err := sms.testMirrorLatency(ctx, m)
            results <- result{m, latency, err}
        }(mirror)
    }
    
    var best MirrorConfig
    var bestLatency time.Duration = time.Hour // 初始值设置得很大
    
    // 收集结果
    for i := 0; i < len(sms.mirrors); i++ {
        res := <-results
        
        fmt.Printf("测试 %s: ", res.mirror.Name)
        if res.err != nil {
            fmt.Printf("失败 - %v\n", res.err)
            continue
        }
        
        fmt.Printf("延迟 %v\n", res.latency)
        
        if res.latency < bestLatency {
            best = res.mirror
            bestLatency = res.latency
        }
    }
    
    if best.Name == "" {
        // 如果所有镜像都失败，返回默认镜像
        fmt.Println("所有镜像测试失败，使用默认官方镜像")
        return sms.mirrors[0]
    }
    
    fmt.Printf("\n选择最佳镜像: %s (延迟: %v)\n", best.Name, bestLatency)
    return best
}

func main() {
    selector := NewSmartMirrorSelector()
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    fmt.Println("正在测试各镜像源延迟...")
    best := selector.SelectBestMirror(ctx)
    
    // 使用最佳镜像获取包信息
    fmt.Printf("\n使用 %s 获取包信息:\n", best.Name)
    packages := []string{"react", "vue", "angular"}
    
    for _, pkg := range packages {
        info, err := best.Client.GetPackageInformation(ctx, pkg)
        if err != nil {
            fmt.Printf("❌ %s: %v\n", pkg, err)
            continue
        }
        
        fmt.Printf("✅ %s: v%s\n", info.Name, info.DistTags["latest"])
    }
}
```

## 示例 2: 地理位置自适应

```go
package main

import (
    "context"
    "fmt"
    "net"
    "strings"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

// 地理位置检测器
type GeoDetector struct{}

func NewGeoDetector() *GeoDetector {
    return &GeoDetector{}
}

// 简单的地理位置检测（基于 DNS 解析）
func (gd *GeoDetector) DetectRegion() string {
    // 检测是否在中国大陆
    if gd.isInChina() {
        return "china"
    }
    
    // 可以添加更多地区检测逻辑
    return "global"
}

func (gd *GeoDetector) isInChina() bool {
    // 尝试解析中国特有的域名
    testDomains := []string{
        "baidu.com",
        "qq.com",
        "taobao.com",
    }
    
    successCount := 0
    for _, domain := range testDomains {
        if _, err := net.LookupHost(domain); err == nil {
            successCount++
        }
    }
    
    // 如果大部分中国域名都能解析，可能在中国大陆
    return successCount >= len(testDomains)/2
}

// 地理位置自适应客户端
type GeoAdaptiveClient struct {
    region   string
    client   *registry.Registry
    detector *GeoDetector
}

func NewGeoAdaptiveClient() *GeoAdaptiveClient {
    detector := NewGeoDetector()
    region := detector.DetectRegion()
    
    var client *registry.Registry
    
    switch region {
    case "china":
        fmt.Println("检测到中国大陆环境，使用 NPM Mirror")
        client = registry.NewNpmMirrorRegistry()
    default:
        fmt.Println("检测到海外环境，使用官方镜像")
        client = registry.NewRegistry()
    }
    
    return &GeoAdaptiveClient{
        region:   region,
        client:   client,
        detector: detector,
    }
}

func (gac *GeoAdaptiveClient) GetPackageInformation(ctx context.Context, packageName string) error {
    pkg, err := gac.client.GetPackageInformation(ctx, packageName)
    if err != nil {
        return err
    }
    
    fmt.Printf("从 %s 镜像获取: %s v%s\n", 
        gac.region, 
        pkg.Name, 
        pkg.DistTags["latest"])
    
    return nil
}

func main() {
    client := NewGeoAdaptiveClient()
    ctx := context.Background()
    
    packages := []string{"react", "vue", "lodash"}
    
    for _, pkg := range packages {
        err := client.GetPackageInformation(ctx, pkg)
        if err != nil {
            fmt.Printf("获取 %s 失败: %v\n", pkg, err)
        }
    }
}
```

## 示例 3: 环境变量配置

```go
package main

import (
    "context"
    "fmt"
    "os"
    "strings"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

// 环境配置管理器
type EnvConfigManager struct{}

func NewEnvConfigManager() *EnvConfigManager {
    return &EnvConfigManager{}
}

func (ecm *EnvConfigManager) GetRegistryURL() string {
    // 检查环境变量
    if url := os.Getenv("NPM_REGISTRY_URL"); url != "" {
        return url
    }
    
    // 检查 NPM 配置
    if url := os.Getenv("npm_config_registry"); url != "" {
        return url
    }
    
    // 检查部署环境
    if env := os.Getenv("DEPLOY_ENV"); env != "" {
        switch strings.ToLower(env) {
        case "china", "cn":
            return "https://registry.npmmirror.com"
        case "dev", "development":
            return "http://r.cnpmjs.org"
        }
    }
    
    // 默认官方镜像
    return "https://registry.npmjs.org"
}

func (ecm *EnvConfigManager) GetProxy() string {
    // 按优先级检查代理配置
    proxies := []string{
        os.Getenv("NPM_PROXY"),
        os.Getenv("HTTP_PROXY"),
        os.Getenv("http_proxy"),
        os.Getenv("HTTPS_PROXY"),
        os.Getenv("https_proxy"),
    }
    
    for _, proxy := range proxies {
        if proxy != "" {
            return proxy
        }
    }
    
    return ""
}

func (ecm *EnvConfigManager) CreateClient() *registry.Registry {
    registryURL := ecm.GetRegistryURL()
    proxy := ecm.GetProxy()
    
    fmt.Printf("使用注册表: %s\n", registryURL)
    if proxy != "" {
        fmt.Printf("使用代理: %s\n", proxy)
    }
    
    options := registry.NewOptions().SetRegistryURL(registryURL)
    
    if proxy != "" {
        options.SetProxy(proxy)
    }
    
    return registry.NewRegistry(options)
}

func main() {
    // 演示不同的环境变量配置
    fmt.Println("=== 环境变量配置示例 ===\n")
    
    // 示例 1: 默认配置
    fmt.Println("1. 默认配置:")
    manager := NewEnvConfigManager()
    client := manager.CreateClient()
    
    ctx := context.Background()
    pkg, err := client.GetPackageInformation(ctx, "lodash")
    if err != nil {
        fmt.Printf("获取失败: %v\n", err)
    } else {
        fmt.Printf("成功获取: %s v%s\n", pkg.Name, pkg.DistTags["latest"])
    }
    
    // 示例 2: 设置中国环境
    fmt.Println("\n2. 设置中国环境:")
    os.Setenv("DEPLOY_ENV", "china")
    manager2 := NewEnvConfigManager()
    client2 := manager2.CreateClient()
    
    pkg2, err := client2.GetPackageInformation(ctx, "react")
    if err != nil {
        fmt.Printf("获取失败: %v\n", err)
    } else {
        fmt.Printf("成功获取: %s v%s\n", pkg2.Name, pkg2.DistTags["latest"])
    }
    
    // 示例 3: 自定义注册表
    fmt.Println("\n3. 自定义注册表:")
    os.Setenv("NPM_REGISTRY_URL", "https://registry.npm.taobao.org")
    manager3 := NewEnvConfigManager()
    client3 := manager3.CreateClient()
    
    pkg3, err := client3.GetPackageInformation(ctx, "vue")
    if err != nil {
        fmt.Printf("获取失败: %v\n", err)
    } else {
        fmt.Printf("成功获取: %s v%s\n", pkg3.Name, pkg3.DistTags["latest"])
    }
}
```

## 示例 4: 多镜像故障转移

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

// 镜像配置
type Mirror struct {
    Name     string
    Client   *registry.Registry
    Priority int
    Timeout  time.Duration
}

// 故障转移客户端
type FailoverClient struct {
    mirrors []Mirror
}

func NewFailoverClient() *FailoverClient {
    mirrors := []Mirror{
        {
            Name:     "NPM Mirror (中国首选)",
            Client:   registry.NewNpmMirrorRegistry(),
            Priority: 1,
            Timeout:  10 * time.Second,
        },
        {
            Name:     "淘宝镜像 (中国备选)",
            Client:   registry.NewTaoBaoRegistry(),
            Priority: 2,
            Timeout:  10 * time.Second,
        },
        {
            Name:     "华为云 (中国备选)",
            Client:   registry.NewHuaWeiCloudRegistry(),
            Priority: 3,
            Timeout:  15 * time.Second,
        },
        {
            Name:     "官方镜像 (全球)",
            Client:   registry.NewRegistry(),
            Priority: 4,
            Timeout:  20 * time.Second,
        },
        {
            Name:     "Yarn 镜像 (全球备选)",
            Client:   registry.NewYarnRegistry(),
            Priority: 5,
            Timeout:  20 * time.Second,
        },
    }
    
    return &FailoverClient{mirrors: mirrors}
}

func (fc *FailoverClient) GetPackageInformation(ctx context.Context, packageName string) error {
    var lastErr error
    
    for _, mirror := range fc.mirrors {
        fmt.Printf("尝试使用 %s...\n", mirror.Name)
        
        // 为每个镜像设置超时
        mirrorCtx, cancel := context.WithTimeout(ctx, mirror.Timeout)
        
        pkg, err := mirror.Client.GetPackageInformation(mirrorCtx, packageName)
        cancel()
        
        if err != nil {
            fmt.Printf("❌ %s 失败: %v\n", mirror.Name, err)
            lastErr = err
            continue
        }
        
        fmt.Printf("✅ %s 成功: %s v%s\n", 
            mirror.Name, 
            pkg.Name, 
            pkg.DistTags["latest"])
        
        return nil
    }
    
    return fmt.Errorf("所有镜像都失败，最后错误: %w", lastErr)
}

func main() {
    client := NewFailoverClient()
    ctx := context.Background()
    
    packages := []string{"react", "nonexistent-package-12345", "vue"}
    
    for _, pkg := range packages {
        fmt.Printf("\n=== 获取包: %s ===\n", pkg)
        
        err := client.GetPackageInformation(ctx, pkg)
        if err != nil {
            fmt.Printf("最终失败: %v\n", err)
        }
        
        time.Sleep(2 * time.Second) // 避免请求过快
    }
}
```

## 示例 5: 企业内网配置

```go
package main

import (
    "context"
    "fmt"
    "net/url"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

// 企业配置管理器
type EnterpriseConfig struct {
    InternalRegistry string
    ExternalRegistry string
    ProxyURL         string
    UseInternal      bool
}

func NewEnterpriseConfig() *EnterpriseConfig {
    return &EnterpriseConfig{
        InternalRegistry: "http://npm.corp.internal:8080",
        ExternalRegistry: "https://registry.npmmirror.com",
        ProxyURL:         "http://proxy.corp.com:8080",
        UseInternal:      true,
    }
}

// 企业客户端
type EnterpriseClient struct {
    config         *EnterpriseConfig
    internalClient *registry.Registry
    externalClient *registry.Registry
}

func NewEnterpriseClient(config *EnterpriseConfig) *EnterpriseClient {
    // 内网客户端
    internalOptions := registry.NewOptions().
        SetRegistryURL(config.InternalRegistry)
    internalClient := registry.NewRegistry(internalOptions)
    
    // 外网客户端（通过代理）
    externalOptions := registry.NewOptions().
        SetRegistryURL(config.ExternalRegistry).
        SetProxy(config.ProxyURL)
    externalClient := registry.NewRegistry(externalOptions)
    
    return &EnterpriseClient{
        config:         config,
        internalClient: internalClient,
        externalClient: externalClient,
    }
}

func (ec *EnterpriseClient) GetPackageInformation(ctx context.Context, packageName string) error {
    // 首先尝试内网镜像
    if ec.config.UseInternal {
        fmt.Printf("尝试内网镜像: %s\n", packageName)
        
        pkg, err := ec.internalClient.GetPackageInformation(ctx, packageName)
        if err == nil {
            fmt.Printf("✅ 内网获取成功: %s v%s\n", 
                pkg.Name, 
                pkg.DistTags["latest"])
            return nil
        }
        
        fmt.Printf("❌ 内网获取失败: %v\n", err)
    }
    
    // fallback 到外网镜像
    fmt.Printf("尝试外网镜像: %s\n", packageName)
    
    pkg, err := ec.externalClient.GetPackageInformation(ctx, packageName)
    if err != nil {
        return fmt.Errorf("内外网都失败: %w", err)
    }
    
    fmt.Printf("✅ 外网获取成功: %s v%s\n", 
        pkg.Name, 
        pkg.DistTags["latest"])
    
    return nil
}

// 根据包名判断是否使用内网
func (ec *EnterpriseClient) isInternalPackage(packageName string) bool {
    // 企业内部包通常有特定前缀
    internalPrefixes := []string{
        "@company/",
        "@corp/",
        "@internal/",
    }
    
    for _, prefix := range internalPrefixes {
        if len(packageName) >= len(prefix) && packageName[:len(prefix)] == prefix {
            return true
        }
    }
    
    return false
}

func main() {
    config := NewEnterpriseConfig()
    client := NewEnterpriseClient(config)
    ctx := context.Background()
    
    packages := []string{
        "react",              // 公共包
        "@company/ui-kit",    // 企业内部包（假设）
        "lodash",             // 公共包
        "@internal/utils",    // 企业内部包（假设）
    }
    
    for _, pkg := range packages {
        fmt.Printf("\n=== 获取包: %s ===\n", pkg)
        
        // 根据包名调整策略
        if client.isInternalPackage(pkg) {
            fmt.Println("检测到内部包，优先使用内网")
            client.config.UseInternal = true
        } else {
            fmt.Println("检测到公共包，可以使用外网")
            client.config.UseInternal = false
        }
        
        err := client.GetPackageInformation(ctx, pkg)
        if err != nil {
            fmt.Printf("获取失败: %v\n", err)
        }
    }
}
```

## 配置文件示例

### Docker 环境配置

```dockerfile
# Dockerfile
FROM golang:1.21-alpine

# 设置镜像源环境变量
ENV NPM_REGISTRY_URL=https://registry.npmmirror.com
ENV DEPLOY_ENV=china

WORKDIR /app
COPY . .
RUN go build -o main .

CMD ["./main"]
```

### Kubernetes 配置

```yaml
# k8s-config.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: npm-config
data:
  NPM_REGISTRY_URL: "https://registry.npmmirror.com"
  DEPLOY_ENV: "china"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: npm-skills-app
spec:
  template:
    spec:
      containers:
      - name: app
        image: npm-skills-app:latest
        envFrom:
        - configMapRef:
            name: npm-config
```

### 开发环境配置

```bash
# .env 文件
NPM_REGISTRY_URL=https://registry.npmmirror.com
NPM_PROXY=http://proxy.corp.com:8080
DEPLOY_ENV=development
```

## 性能对比

运行以下脚本可以对比不同镜像源的性能：

```bash
go run mirror-benchmark.go
```

示例输出：
```
镜像源性能测试结果:
✅ NPM Mirror: 平均延迟 150ms, 成功率 100%
✅ 淘宝镜像: 平均延迟 200ms, 成功率 100%
✅ 华为云: 平均延迟 300ms, 成功率 98%
✅ 官方镜像: 平均延迟 800ms, 成功率 95%
```

## 最佳实践建议

1. **中国大陆用户**: 优先使用 NPM Mirror 或淘宝镜像
2. **海外用户**: 使用官方镜像或就近的 CDN
3. **企业环境**: 配置内网镜像和代理
4. **容器部署**: 通过环境变量配置镜像源
5. **开发测试**: 可以使用 CNPM 等测试镜像

## 下一步

- 返回 [基本用法示例](/examples/basic) 学习基础功能
- 查看 [高级用法示例](/examples/advanced) 了解更多技巧
- 阅读 [API 文档](/api/) 了解完整功能
