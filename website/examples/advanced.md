# 高级用法示例

本页面展示 NPM Skills 的高级使用场景和技巧。

## 示例 1: 缓存实现

```go
package main

import (
    "context"
    "fmt"
    "sync"
    "time"

    "github.com/scagogogo/npm-skills/pkg/models"
    "github.com/scagogogo/npm-skills/pkg/registry"
)

// 包信息缓存
type PackageCache struct {
    cache map[string]*CacheEntry
    mutex sync.RWMutex
}

type CacheEntry struct {
    Package   *models.Package
    Timestamp time.Time
    TTL       time.Duration
}

func NewPackageCache() *PackageCache {
    return &PackageCache{
        cache: make(map[string]*CacheEntry),
    }
}

func (c *PackageCache) Get(packageName string) (*models.Package, bool) {
    c.mutex.RLock()
    defer c.mutex.RUnlock()
    
    entry, exists := c.cache[packageName]
    if !exists {
        return nil, false
    }
    
    // 检查是否过期
    if time.Since(entry.Timestamp) > entry.TTL {
        return nil, false
    }
    
    return entry.Package, true
}

func (c *PackageCache) Set(packageName string, pkg *models.Package, ttl time.Duration) {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    
    c.cache[packageName] = &CacheEntry{
        Package:   pkg,
        Timestamp: time.Now(),
        TTL:       ttl,
    }
}

// 带缓存的客户端
type CachedClient struct {
    client *registry.Registry
    cache  *PackageCache
}

func NewCachedClient() *CachedClient {
    return &CachedClient{
        client: registry.NewRegistry(),
        cache:  NewPackageCache(),
    }
}

func (cc *CachedClient) GetPackageInformation(ctx context.Context, packageName string) (*models.Package, error) {
    // 先检查缓存
    if pkg, hit := cc.cache.Get(packageName); hit {
        fmt.Printf("缓存命中: %s\n", packageName)
        return pkg, nil
    }
    
    // 缓存未命中，从网络获取
    fmt.Printf("从网络获取: %s\n", packageName)
    pkg, err := cc.client.GetPackageInformation(ctx, packageName)
    if err != nil {
        return nil, err
    }
    
    // 存入缓存，TTL 为 10 分钟
    cc.cache.Set(packageName, pkg, 10*time.Minute)
    
    return pkg, nil
}

func main() {
    client := NewCachedClient()
    ctx := context.Background()
    
    packages := []string{"react", "vue", "react", "angular", "react"}
    
    for _, pkg := range packages {
        info, err := client.GetPackageInformation(ctx, pkg)
        if err != nil {
            fmt.Printf("错误: %v\n", err)
            continue
        }
        
        fmt.Printf("包: %s, 版本: %s\n", info.Name, info.DistTags["latest"])
    }
}
```

## 示例 2: 重试机制

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "log"
    "math"
    "time"

    "github.com/scagogogo/npm-skills/pkg/models"
    "github.com/scagogogo/npm-skills/pkg/registry"
)

// 重试配置
type RetryConfig struct {
    MaxRetries int
    BaseDelay  time.Duration
    MaxDelay   time.Duration
    Multiplier float64
}

// 默认重试配置
var DefaultRetryConfig = RetryConfig{
    MaxRetries: 3,
    BaseDelay:  1 * time.Second,
    MaxDelay:   30 * time.Second,
    Multiplier: 2.0,
}

// 带重试的客户端
type RetryClient struct {
    client *registry.Registry
    config RetryConfig
}

func NewRetryClient(config RetryConfig) *RetryClient {
    return &RetryClient{
        client: registry.NewRegistry(),
        config: config,
    }
}

// 指数退避重试
func (rc *RetryClient) GetPackageInformation(ctx context.Context, packageName string) (*models.Package, error) {
    var lastErr error
    
    for attempt := 0; attempt <= rc.config.MaxRetries; attempt++ {
        if attempt > 0 {
            delay := rc.calculateDelay(attempt)
            fmt.Printf("第 %d 次重试 %s，等待 %v\n", attempt, packageName, delay)
            
            select {
            case <-time.After(delay):
            case <-ctx.Done():
                return nil, ctx.Err()
            }
        }
        
        pkg, err := rc.client.GetPackageInformation(ctx, packageName)
        if err == nil {
            if attempt > 0 {
                fmt.Printf("重试成功: %s\n", packageName)
            }
            return pkg, nil
        }
        
        lastErr = err
        
        // 某些错误不需要重试
        if isNonRetryableError(err) {
            break
        }
    }
    
    return nil, fmt.Errorf("重试 %d 次后仍然失败: %w", rc.config.MaxRetries, lastErr)
}

func (rc *RetryClient) calculateDelay(attempt int) time.Duration {
    delay := time.Duration(float64(rc.config.BaseDelay) * math.Pow(rc.config.Multiplier, float64(attempt-1)))
    
    if delay > rc.config.MaxDelay {
        delay = rc.config.MaxDelay
    }
    
    return delay
}

func isNonRetryableError(err error) bool {
    // 这里可以根据错误类型判断是否需要重试
    if errors.Is(err, context.Canceled) {
        return true
    }
    
    // 可以添加更多不需要重试的错误类型
    return false
}

func main() {
    client := NewRetryClient(DefaultRetryConfig)
    ctx := context.Background()
    
    // 测试正常包
    pkg, err := client.GetPackageInformation(ctx, "react")
    if err != nil {
        log.Printf("获取 react 失败: %v", err)
    } else {
        fmt.Printf("成功获取: %s v%s\n", pkg.Name, pkg.DistTags["latest"])
    }
    
    // 测试不存在的包（会触发重试）
    _, err = client.GetPackageInformation(ctx, "nonexistent-package-12345")
    if err != nil {
        log.Printf("获取不存在的包失败（预期）: %v", err)
    }
}
```

## 示例 3: 连接池和性能优化

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "sync"
    "time"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

// 连接池管理器
type ConnectionManager struct {
    clients map[string]*registry.Registry
    mutex   sync.RWMutex
}

func NewConnectionManager() *ConnectionManager {
    return &ConnectionManager{
        clients: make(map[string]*registry.Registry),
    }
}

func (cm *ConnectionManager) GetClient(registryURL string) *registry.Registry {
    cm.mutex.RLock()
    client, exists := cm.clients[registryURL]
    cm.mutex.RUnlock()
    
    if exists {
        return client
    }
    
    cm.mutex.Lock()
    defer cm.mutex.Unlock()
    
    // 双重检查
    if client, exists := cm.clients[registryURL]; exists {
        return client
    }
    
    // 创建新客户端
    options := registry.NewOptions().SetRegistryURL(registryURL)
    client = registry.NewRegistry(options)
    cm.clients[registryURL] = client
    
    return client
}

// 负载均衡客户端
type LoadBalancedClient struct {
    clients []string
    manager *ConnectionManager
    current int
    mutex   sync.Mutex
}

func NewLoadBalancedClient(registryURLs []string) *LoadBalancedClient {
    return &LoadBalancedClient{
        clients: registryURLs,
        manager: NewConnectionManager(),
    }
}

func (lbc *LoadBalancedClient) getNextClient() *registry.Registry {
    lbc.mutex.Lock()
    url := lbc.clients[lbc.current]
    lbc.current = (lbc.current + 1) % len(lbc.clients)
    lbc.mutex.Unlock()
    
    return lbc.manager.GetClient(url)
}

func (lbc *LoadBalancedClient) GetPackageInformation(ctx context.Context, packageName string) error {
    client := lbc.getNextClient()
    
    pkg, err := client.GetPackageInformation(ctx, packageName)
    if err != nil {
        return err
    }
    
    fmt.Printf("从 %s 获取: %s v%s\n", 
        client.GetOptions().RegistryURL, 
        pkg.Name, 
        pkg.DistTags["latest"])
    
    return nil
}

func main() {
    // 配置多个镜像源进行负载均衡
    registryURLs := []string{
        "https://registry.npmjs.org",
        "https://registry.npmmirror.com",
        "https://registry.npm.taobao.org",
    }
    
    client := NewLoadBalancedClient(registryURLs)
    ctx := context.Background()
    
    packages := []string{"react", "vue", "angular", "lodash", "express", "axios"}
    
    var wg sync.WaitGroup
    for _, pkg := range packages {
        wg.Add(1)
        go func(packageName string) {
            defer wg.Done()
            
            err := client.GetPackageInformation(ctx, packageName)
            if err != nil {
                fmt.Printf("获取 %s 失败: %v\n", packageName, err)
            }
        }(pkg)
    }
    
    wg.Wait()
}
```

## 示例 4: 限流器

```go
package main

import (
    "context"
    "fmt"
    "sync"
    "time"

    "github.com/scagogogo/npm-skills/pkg/registry"
    "golang.org/x/time/rate"
)

// 限流客户端
type RateLimitedClient struct {
    client  *registry.Registry
    limiter *rate.Limiter
}

func NewRateLimitedClient(requestsPerSecond int) *RateLimitedClient {
    return &RateLimitedClient{
        client:  registry.NewRegistry(),
        limiter: rate.NewLimiter(rate.Limit(requestsPerSecond), requestsPerSecond),
    }
}

func (rlc *RateLimitedClient) GetPackageInformation(ctx context.Context, packageName string) error {
    // 等待令牌
    if err := rlc.limiter.Wait(ctx); err != nil {
        return err
    }
    
    pkg, err := rlc.client.GetPackageInformation(ctx, packageName)
    if err != nil {
        return err
    }
    
    fmt.Printf("✅ %s: v%s\n", pkg.Name, pkg.DistTags["latest"])
    return nil
}

func main() {
    // 限制每秒 2 个请求
    client := NewRateLimitedClient(2)
    ctx := context.Background()
    
    packages := []string{
        "react", "vue", "angular", "lodash", "express",
        "axios", "moment", "underscore", "jquery", "bootstrap",
    }
    
    start := time.Now()
    
    var wg sync.WaitGroup
    for i, pkg := range packages {
        wg.Add(1)
        go func(index int, packageName string) {
            defer wg.Done()
            
            fmt.Printf("[%d] 开始请求: %s\n", index, packageName)
            err := client.GetPackageInformation(ctx, packageName)
            if err != nil {
                fmt.Printf("❌ [%d] %s: %v\n", index, packageName, err)
            }
        }(i, pkg)
    }
    
    wg.Wait()
    
    fmt.Printf("\n总耗时: %v\n", time.Since(start))
}
```

## 示例 5: 健康检查和故障转移

```go
package main

import (
    "context"
    "fmt"
    "sync"
    "time"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

// 镜像源状态
type MirrorStatus struct {
    URL       string
    Healthy   bool
    LastCheck time.Time
    Latency   time.Duration
}

// 健康检查管理器
type HealthChecker struct {
    mirrors map[string]*MirrorStatus
    mutex   sync.RWMutex
}

func NewHealthChecker(urls []string) *HealthChecker {
    hc := &HealthChecker{
        mirrors: make(map[string]*MirrorStatus),
    }
    
    for _, url := range urls {
        hc.mirrors[url] = &MirrorStatus{
            URL:     url,
            Healthy: true,
        }
    }
    
    return hc
}

func (hc *HealthChecker) CheckHealth(ctx context.Context) {
    var wg sync.WaitGroup
    
    for url := range hc.mirrors {
        wg.Add(1)
        go func(mirrorURL string) {
            defer wg.Done()
            hc.checkSingleMirror(ctx, mirrorURL)
        }(url)
    }
    
    wg.Wait()
}

func (hc *HealthChecker) checkSingleMirror(ctx context.Context, url string) {
    start := time.Now()
    
    options := registry.NewOptions().SetRegistryURL(url)
    client := registry.NewRegistry(options)
    
    // 使用简单的包查询来检查健康状态
    _, err := client.GetPackageInformation(ctx, "lodash")
    
    hc.mutex.Lock()
    defer hc.mutex.Unlock()
    
    status := hc.mirrors[url]
    status.LastCheck = time.Now()
    status.Latency = time.Since(start)
    status.Healthy = (err == nil)
    
    if status.Healthy {
        fmt.Printf("✅ %s: 健康 (延迟: %v)\n", url, status.Latency)
    } else {
        fmt.Printf("❌ %s: 不健康 - %v\n", url, err)
    }
}

func (hc *HealthChecker) GetHealthyMirrors() []string {
    hc.mutex.RLock()
    defer hc.mutex.RUnlock()
    
    var healthy []string
    for url, status := range hc.mirrors {
        if status.Healthy {
            healthy = append(healthy, url)
        }
    }
    
    return healthy
}

func (hc *HealthChecker) GetBestMirror() string {
    hc.mutex.RLock()
    defer hc.mutex.RUnlock()
    
    var best string
    var bestLatency time.Duration = time.Hour // 初始值设置得很大
    
    for url, status := range hc.mirrors {
        if status.Healthy && status.Latency < bestLatency {
            best = url
            bestLatency = status.Latency
        }
    }
    
    return best
}

// 自适应客户端
type AdaptiveClient struct {
    checker *HealthChecker
    clients map[string]*registry.Registry
    mutex   sync.RWMutex
}

func NewAdaptiveClient(urls []string) *AdaptiveClient {
    return &AdaptiveClient{
        checker: NewHealthChecker(urls),
        clients: make(map[string]*registry.Registry),
    }
}

func (ac *AdaptiveClient) getClient(url string) *registry.Registry {
    ac.mutex.RLock()
    client, exists := ac.clients[url]
    ac.mutex.RUnlock()
    
    if exists {
        return client
    }
    
    ac.mutex.Lock()
    defer ac.mutex.Unlock()
    
    if client, exists := ac.clients[url]; exists {
        return client
    }
    
    options := registry.NewOptions().SetRegistryURL(url)
    client = registry.NewRegistry(options)
    ac.clients[url] = client
    
    return client
}

func (ac *AdaptiveClient) GetPackageInformation(ctx context.Context, packageName string) error {
    // 定期健康检查
    go ac.checker.CheckHealth(ctx)
    
    // 获取最佳镜像
    bestMirror := ac.checker.GetBestMirror()
    if bestMirror == "" {
        return fmt.Errorf("没有可用的健康镜像")
    }
    
    client := ac.getClient(bestMirror)
    pkg, err := client.GetPackageInformation(ctx, packageName)
    if err != nil {
        return err
    }
    
    fmt.Printf("使用镜像 %s 获取: %s v%s\n", 
        bestMirror, 
        pkg.Name, 
        pkg.DistTags["latest"])
    
    return nil
}

func main() {
    mirrors := []string{
        "https://registry.npmjs.org",
        "https://registry.npmmirror.com",
        "https://registry.npm.taobao.org",
        "https://mirrors.huaweicloud.com/repository/npm",
    }
    
    client := NewAdaptiveClient(mirrors)
    ctx := context.Background()
    
    // 先进行健康检查
    fmt.Println("进行健康检查...")
    client.checker.CheckHealth(ctx)
    
    fmt.Println("\n开始获取包信息...")
    packages := []string{"react", "vue", "angular", "lodash", "express"}
    
    for _, pkg := range packages {
        err := client.GetPackageInformation(ctx, pkg)
        if err != nil {
            fmt.Printf("获取 %s 失败: %v\n", pkg, err)
        }
        
        time.Sleep(1 * time.Second) // 避免请求过快
    }
}
```

## 示例 6: 监控和指标收集

```go
package main

import (
    "context"
    "fmt"
    "sync"
    "time"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

// 指标收集器
type Metrics struct {
    TotalRequests     int64
    SuccessfulRequests int64
    FailedRequests    int64
    TotalLatency      time.Duration
    mutex             sync.RWMutex
}

func NewMetrics() *Metrics {
    return &Metrics{}
}

func (m *Metrics) RecordRequest(latency time.Duration, success bool) {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    
    m.TotalRequests++
    m.TotalLatency += latency
    
    if success {
        m.SuccessfulRequests++
    } else {
        m.FailedRequests++
    }
}

func (m *Metrics) GetStats() (total, success, failed int64, avgLatency time.Duration) {
    m.mutex.RLock()
    defer m.mutex.RUnlock()
    
    total = m.TotalRequests
    success = m.SuccessfulRequests
    failed = m.FailedRequests
    
    if total > 0 {
        avgLatency = m.TotalLatency / time.Duration(total)
    }
    
    return
}

func (m *Metrics) GetSuccessRate() float64 {
    m.mutex.RLock()
    defer m.mutex.RUnlock()
    
    if m.TotalRequests == 0 {
        return 0
    }
    
    return float64(m.SuccessfulRequests) / float64(m.TotalRequests) * 100
}

// 监控客户端
type MonitoredClient struct {
    client  *registry.Registry
    metrics *Metrics
}

func NewMonitoredClient() *MonitoredClient {
    return &MonitoredClient{
        client:  registry.NewRegistry(),
        metrics: NewMetrics(),
    }
}

func (mc *MonitoredClient) GetPackageInformation(ctx context.Context, packageName string) error {
    start := time.Now()
    
    pkg, err := mc.client.GetPackageInformation(ctx, packageName)
    latency := time.Since(start)
    
    success := err == nil
    mc.metrics.RecordRequest(latency, success)
    
    if success {
        fmt.Printf("✅ %s: v%s (耗时: %v)\n", 
            pkg.Name, 
            pkg.DistTags["latest"], 
            latency)
    } else {
        fmt.Printf("❌ %s: %v (耗时: %v)\n", 
            packageName, 
            err, 
            latency)
    }
    
    return err
}

func (mc *MonitoredClient) PrintStats() {
    total, success, failed, avgLatency := mc.metrics.GetStats()
    successRate := mc.metrics.GetSuccessRate()
    
    fmt.Printf("\n=== 统计信息 ===\n")
    fmt.Printf("总请求数: %d\n", total)
    fmt.Printf("成功请求: %d\n", success)
    fmt.Printf("失败请求: %d\n", failed)
    fmt.Printf("成功率: %.2f%%\n", successRate)
    fmt.Printf("平均延迟: %v\n", avgLatency)
}

// 定期打印统计信息
func (mc *MonitoredClient) StartMonitoring(interval time.Duration) {
    ticker := time.NewTicker(interval)
    go func() {
        for range ticker.C {
            mc.PrintStats()
        }
    }()
}

func main() {
    client := NewMonitoredClient()
    ctx := context.Background()
    
    // 启动监控
    client.StartMonitoring(10 * time.Second)
    
    packages := []string{
        "react", "vue", "angular", "lodash", "express",
        "axios", "moment", "underscore", "jquery", "bootstrap",
        "nonexistent-package-1", "nonexistent-package-2", // 故意添加不存在的包
    }
    
    var wg sync.WaitGroup
    
    for _, pkg := range packages {
        wg.Add(1)
        go func(packageName string) {
            defer wg.Done()
            client.GetPackageInformation(ctx, packageName)
        }(pkg)
        
        // 稍微延迟，避免过快请求
        time.Sleep(500 * time.Millisecond)
    }
    
    wg.Wait()
    
    // 最终统计
    client.PrintStats()
}
```

## 运行高级示例

这些高级示例需要额外的依赖：

```bash
go mod init npm-skills-advanced
go get github.com/scagogogo/npm-skills
go get golang.org/x/time/rate
go get golang.org/x/sync/semaphore
go run example.go
```

## 下一步

- 查看 [镜像源配置示例](/examples/mirrors) 了解如何优化不同地区的访问
- 阅读 [API 文档](/api/) 了解所有可用功能
- 查看 [GitHub 仓库](https://github.com/scagogogo/npm-skills) 获取更多示例
