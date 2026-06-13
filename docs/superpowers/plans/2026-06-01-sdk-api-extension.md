# SDK API Extension Plan

> **For agentic workers:** REQUIRED SUB-SKILL: `superpowers:subagent-driven-development`
> Steps use checkbox (`- [ ]`) syntax.

**Goal:** 扩展 SDK API 封装，补全 NPM Registry 中与"查询工具"定位一致的只读端点：下载统计增强（区间/批量/日期范围）、搜索增强（分页+权重）、Dist-Tags 查询、精简包元数据、认证 Token 支持。

**Architecture:** 在 `Options` 中新增 `Token` 字段支持认证 → 将下载统计相关方法从 `registry.go` 拆分到独立文件 `download_stats.go` 并新增 3 个方法 → 将搜索方法拆到 `search.go` 并增强分页参数 → 新增 `dist_tags.go` 提供 dist-tags 只读查询 → 新增 `abbreviated.go` 支持精简包元数据 → 所有新方法都复用现有 `getBytes` + `unmarshalJson` 基础设施。

**Tech Stack:** Go 1.20, github.com/crawler-go-go-go/go-requests, stretchr/testify, net/http/httptest

**Risks:**
- Task 1 修改 `Options` 结构体，影响所有使用 `NewOptions()` 的代码 → 缓解：Token 字段零值为空字符串，向后兼容
- Task 2 从 `registry.go` 移动 `GetDownloadStats` 方法到新文件 → 缓解：Go 允许同一包内方法定义在不同文件，只需移动函数体，不改签名
- Task 3 从 `registry.go` 移动 `SearchPackages` 方法 → 缓解：同 Task 2

---

### Task 1: Add Token Support to Options

**Depends on:** None
**Files:**
- Modify: `pkg/registry/options.go:26-29`（Options 结构体新增 Token 字段）
- Modify: `pkg/registry/options.go:89-93`（新增 SetToken 方法）
- Modify: `pkg/registry/registry.go:326-333`（getBytes 支持 Token 认证头）
- Modify: `pkg/registry/options_test.go`（新增 Token 测试）

- [ ] **Step 1: 修改 Options 结构体 — 新增 Token 字段以支持认证 API**

文件: `pkg/registry/options.go:26-29`（替换 Options 结构体定义）

```go
type Options struct {
	RegistryURL string
	Proxy       string
	Token       string // Bearer token for authenticated API requests
}
```

- [ ] **Step 2: 新增 SetToken 方法到 options.go — 支持 token 的链式设置**

在 `pkg/registry/options.go` 第 93 行（`SetProxy` 方法之后）追加：

```go
// SetToken 设置用于认证 API 请求的 Bearer token
//
// 参数:
//   - token: Bearer token 字符串，通常从 npm token create 或 npm login 获取
//
// 返回值:
//   - *Options: 更新后的选项对象 (支持链式调用)
//
// 使用示例:
//
//	options := NewOptions().SetToken("npm_xxxxx")
//	registry := NewRegistry(options)
func (o *Options) SetToken(token string) *Options {
	o.Token = token
	return o
}
```

- [ ] **Step 3: 修改 getBytes 方法 — 在请求中附加 Authorization 头**

文件: `pkg/registry/registry.go:326-333`（替换 getBytes 方法）

```go
func (x *Registry) getBytes(ctx context.Context, targetUrl string) ([]byte, error) {
	options := requests.NewOptions[any, []byte](targetUrl, requests.BytesResponseHandler())
	if x.options.Proxy != "" {
		options.AppendRequestSetting(requests.RequestSettingProxy(x.options.Proxy))
	}
	if x.options.Token != "" {
		options.AppendRequestSetting(requests.RequestSettingHeader("Authorization", "Bearer "+x.options.Token))
	}
	return requests.SendRequest[any, []byte](ctx, options)
}
```

> 注意：`requests.RequestSettingHeader` 是假设 `go-requests` 库支持自定义 header 的 API。如实际库不支持此函数名，需查看 `go-requests` 包的 API 调整为正确的函数名。

- [ ] **Step 4: 新增 Token 相关测试**

在 `pkg/registry/options_test.go` 文件末尾追加：

```go
func TestSetToken(t *testing.T) {
	options := NewOptions()

	// 测试默认 token 为空
	assert.Empty(t, options.Token)

	// 测试设置 token
	result := options.SetToken("npm_test_token")
	assert.Equal(t, options, result, "应该返回自身以支持链式调用")
	assert.Equal(t, "npm_test_token", options.Token)

	// 测试清除 token
	options.SetToken("")
	assert.Empty(t, options.Token)

	// 测试链式调用
	options = NewOptions().
		SetRegistryURL("https://registry.npmjs.org").
		SetProxy("http://proxy:8080").
		SetToken("npm_chained_token")
	assert.Equal(t, "npm_chained_token", options.Token)
	assert.Equal(t, "http://proxy:8080", options.Proxy)
}

func TestRegistryWithToken(t *testing.T) {
	options := NewOptions().SetToken("npm_test_token")
	registry := NewRegistry(options)
	retrievedOpts := registry.GetOptions()
	assert.Equal(t, "npm_test_token", retrievedOpts.Token)
}
```

- [ ] **Step 5: 验证 Task 1**
Run: `cd /home/cc11001100/github/scagogogo/npm-skills && go build ./pkg/registry/ && go test ./pkg/registry/ -run "TestSetToken|TestRegistryWithToken" -v`
Expected:
  - Exit code: 0
  - Output contains: "PASS"

- [ ] **Step 6: 提交**
Run: `cd /home/cc11001100/github/scagogogo/npm-skills && git add pkg/registry/options.go pkg/registry/options_test.go pkg/registry/registry.go && git commit -m "feat(registry): add Token field to Options for authenticated API support"`

---

### Task 2: Download Stats API Extension

**Depends on:** Task 1
**Files:**
- Create: `pkg/registry/download_stats.go`（下载统计方法拆分+新增）
- Modify: `pkg/registry/registry.go:209-237`（删除原 GetDownloadStats 方法）
- Modify: `pkg/registry/registry_test.go`（新增下载统计测试）

- [ ] **Step 1: 创建 download_stats.go — 包含现有 GetDownloadStats + 3 个新方法**

```go
package registry

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/scagogogo/npm-skills/pkg/models"
)

// downloadStatsBaseURL 是 NPM 下载统计 API 的基础 URL
const downloadStatsBaseURL = "https://api.npmjs.org/downloads"

// GetDownloadStats 获取指定 NPM 包的下载统计信息（单个包、预定义周期）
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - packageName: 要查询的包名称
//   - period: 统计周期，例如 "last-day", "last-week", "last-month"
//
// 返回值:
//   - *models.DownloadStats: 下载统计信息
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	registry := NewRegistry()
//	ctx := context.Background()
//	stats, err := registry.GetDownloadStats(ctx, "react", "last-week")
//	if err != nil {
//	    // 处理错误
//	}
//	fmt.Println("下载次数:", stats.Downloads)
func (x *Registry) GetDownloadStats(ctx context.Context, packageName, period string) (*models.DownloadStats, error) {
	targetUrl := fmt.Sprintf("%s/point/%s/%s", downloadStatsBaseURL, period, packageName)
	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, err
	}
	return unmarshalJson[*models.DownloadStats](bytes)
}

// GetDownloadRangeStats 获取指定 NPM 包的每日下载统计（区间数据）
//
// 返回每日下载次数数组，适用于绘制下载趋势图。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - packageName: 要查询的包名称
//   - period: 统计周期，"last-day"、"last-week"、"last-month"
//
// 返回值:
//   - *models.DownloadRangeStats: 包含每日下载数据的区间统计信息
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	registry := NewRegistry()
//	ctx := context.Background()
//	stats, err := registry.GetDownloadRangeStats(ctx, "react", "last-week")
//	if err != nil {
//	    // 处理错误
//	}
//	for _, day := range stats.Downloads {
//	    fmt.Printf("%s: %d\n", day.Day, day.Downloads)
//	}
func (x *Registry) GetDownloadRangeStats(ctx context.Context, packageName, period string) (*models.DownloadRangeStats, error) {
	targetUrl := fmt.Sprintf("%s/range/%s/%s", downloadStatsBaseURL, period, packageName)
	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, err
	}
	return unmarshalJson[*models.DownloadRangeStats](bytes)
}

// GetDownloadStatsByDateRange 获取指定日期范围的下载统计
//
// 使用自定义日期范围而非预定义周期。日期格式为 YYYY-MM-DD。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - packageName: 要查询的包名称
//   - start: 开始日期，格式 YYYY-MM-DD
//   - end: 结束日期，格式 YYYY-MM-DD
//
// 返回值:
//   - *models.DownloadStats: 下载统计信息
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	registry := NewRegistry()
//	ctx := context.Background()
//	stats, err := registry.GetDownloadStatsByDateRange(ctx, "react", "2024-01-01", "2024-01-31")
func (x *Registry) GetDownloadStatsByDateRange(ctx context.Context, packageName, start, end string) (*models.DownloadStats, error) {
	targetUrl := fmt.Sprintf("%s/point/%s:%s/%s", downloadStatsBaseURL, start, end, packageName)
	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, err
	}
	return unmarshalJson[*models.DownloadStats](bytes)
}

// GetBulkDownloadStats 批量获取多个包的下载统计（最多 128 个包）
//
// 适用于比较多个包的下载量，一次请求替代多次单独查询。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - packageNames: 包名称切片，最多 128 个
//   - period: 统计周期，"last-day"、"last-week"、"last-month"
//
// 返回值:
//   - map[string]*models.DownloadStats: 包名到下载统计的映射
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	registry := NewRegistry()
//	ctx := context.Background()
//	stats, err := registry.GetBulkDownloadStats(ctx, []string{"react", "vue", "angular"}, "last-week")
func (x *Registry) GetBulkDownloadStats(ctx context.Context, packageNames []string, period string) (map[string]*models.DownloadStats, error) {
	if len(packageNames) == 0 {
		return nil, fmt.Errorf("packageNames must not be empty")
	}
	if len(packageNames) > 128 {
		return nil, fmt.Errorf("packageNames must not exceed 128, got %d", len(packageNames))
	}
	escaped := url.QueryEscape(strings.Join(packageNames, ","))
	targetUrl := fmt.Sprintf("%s/point/%s/%s", downloadStatsBaseURL, period, escaped)
	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, err
	}
	return unmarshalJson[map[string]*models.DownloadStats](bytes)
}

// GetBulkDownloadRangeStats 批量获取多个包的每日下载统计（最多 128 个包）
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - packageNames: 包名称切片，最多 128 个
//   - period: 统计周期
//
// 返回值:
//   - map[string]*models.DownloadRangeStats: 包名到区间统计的映射
//   - error: 如果请求失败则返回错误
func (x *Registry) GetBulkDownloadRangeStats(ctx context.Context, packageNames []string, period string) (map[string]*models.DownloadRangeStats, error) {
	if len(packageNames) == 0 {
		return nil, fmt.Errorf("packageNames must not be empty")
	}
	if len(packageNames) > 128 {
		return nil, fmt.Errorf("packageNames must not exceed 128, got %d", len(packageNames))
	}
	escaped := url.QueryEscape(strings.Join(packageNames, ","))
	targetUrl := fmt.Sprintf("%s/range/%s/%s", downloadStatsBaseURL, period, escaped)
	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, err
	}
	return unmarshalJson[map[string]*models.DownloadRangeStats](bytes)
}
```

- [ ] **Step 2: 从 registry.go 中删除原 GetDownloadStats 方法**

删除 `pkg/registry/registry.go` 第 209-237 行（原 `GetDownloadStats` 方法），该方法已移至 `download_stats.go`。

同时删除第 10 行中 `"net/url"` 的 import（如不再被其他方法使用），以及第 8 行中 `"os"` 的 import（如不再被其他方法使用）。

> 注意：删除后需验证编译通过，`import` 清理由编译器报错确定。

- [ ] **Step 3: 验证编译和已有测试**
Run: `cd /home/cc11001100/github/scagogogo/npm-skills && go build ./pkg/registry/ && go test ./pkg/registry/ -run "TestAllRegistryImplementations|TestMirrorRegistryCreation" -v`
Expected:
  - Exit code: 0
  - Output contains: "PASS"

- [ ] **Step 4: 提交**
Run: `cd /home/cc11001100/github/scagogogo/npm-skills && git add pkg/registry/download_stats.go pkg/registry/registry.go && git commit -m "feat(registry): add GetDownloadRangeStats, GetDownloadStatsByDateRange, GetBulkDownloadStats, GetBulkDownloadRangeStats"`

---

### Task 3: Search API Enhancement

**Depends on:** Task 1
**Files:**
- Create: `pkg/registry/search.go`（搜索方法拆分+增强）
- Modify: `pkg/registry/registry.go:145-177`（删除原 SearchPackages 方法）

- [ ] **Step 1: 创建 search.go — 包含增强版 SearchPackages 方法**

```go
package registry

import (
	"context"
	"fmt"
	"net/url"

	"github.com/scagogogo/npm-skills/pkg/models"
)

// SearchOptions 定义 NPM 搜索的可选参数
//
// 用于控制搜索行为，包括分页和结果权重调整。
type SearchOptions struct {
	From         int     // 分页偏移量，默认 0
	Size         int     // 返回结果数量，默认 20
	Quality      float64 // 质量权重，0.0-1.0
	Popularity   float64 // 流行度权重，0.0-1.0
	Maintenance  float64 // 维护性权重，0.0-1.0
}

// SearchPackages 搜索 NPM 包
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - query: 搜索关键字
//   - limit: 返回结果数量限制，默认为 20
//
// 返回值:
//   - *models.SearchResult: 搜索结果，包含匹配的包列表
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	registry := NewRegistry()
//	ctx := context.Background()
//	result, err := registry.SearchPackages(ctx, "react", 10)
func (x *Registry) SearchPackages(ctx context.Context, query string, limit int) (*models.SearchResult, error) {
	return x.SearchPackagesWithOptions(ctx, query, SearchOptions{Size: limit})
}

// SearchPackagesWithOptions 使用完整选项搜索 NPM 包
//
// 支持分页（from）和权重调整（quality/popularity/maintenance）。
//
// 参数:
//   - ctx: 上下文
//   - query: 搜索关键字
//   - opts: 搜索选项（分页、数量、权重）
//
// 返回值:
//   - *models.SearchResult: 搜索结果
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	registry := NewRegistry()
//	ctx := context.Background()
//
//	// 第 21-40 条结果
//	result, err := registry.SearchPackagesWithOptions(ctx, "http client", SearchOptions{
//	    From: 20,
//	    Size: 20,
//	})
//
//	// 侧重流行度的搜索
//	result, err = registry.SearchPackagesWithOptions(ctx, "http client", SearchOptions{
//	    Size:        10,
//	    Popularity:  1.0,
//	    Quality:     0.0,
//	    Maintenance: 0.0,
//	})
func (x *Registry) SearchPackagesWithOptions(ctx context.Context, query string, opts SearchOptions) (*models.SearchResult, error) {
	if opts.Size <= 0 {
		opts.Size = 20
	}

	params := url.Values{}
	params.Set("text", query)
	params.Set("size", fmt.Sprintf("%d", opts.Size))
	if opts.From > 0 {
		params.Set("from", fmt.Sprintf("%d", opts.From))
	}
	if opts.Quality > 0 {
		params.Set("quality", fmt.Sprintf("%.2f", opts.Quality))
	}
	if opts.Popularity > 0 {
		params.Set("popularity", fmt.Sprintf("%.2f", opts.Popularity))
	}
	if opts.Maintenance > 0 {
		params.Set("maintenance", fmt.Sprintf("%.2f", opts.Maintenance))
	}

	targetUrl := fmt.Sprintf("%s/-/v1/search?%s", x.options.RegistryURL, params.Encode())
	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, err
	}
	return unmarshalJson[*models.SearchResult](bytes)
}
```

- [ ] **Step 2: 从 registry.go 中删除原 SearchPackages 方法**

删除 `pkg/registry/registry.go` 第 145-177 行（原 `SearchPackages` 方法），该方法已移至 `search.go`。

同时删除 `registry.go` import 中 `"net/url"` （如仍被其他方法使用则保留）。

> 注意：`GetDownloadStats` 已在 Task 2 中移走，现在 `SearchPackages` 也移走，`registry.go` 的 import 列表会变短。编译后按报错清理。

- [ ] **Step 3: 验证编译**
Run: `cd /home/cc11001100/github/scagogogo/npm-skills && go build ./pkg/registry/ && go test ./pkg/registry/ -count=1 -v 2>&1 | tail -10`
Expected:
  - Exit code: 0
  - Output does NOT contain: "FAIL"

- [ ] **Step 4: 提交**
Run: `cd /home/cc11001100/github/scagogogo/npm-skills && git add pkg/registry/search.go pkg/registry/registry.go && git commit -m "feat(registry): add SearchPackagesWithOptions with pagination and score weighting"`

---

### Task 4: Dist-Tags Query API

**Depends on:** Task 1
**Files:**
- Create: `pkg/registry/dist_tags.go`

- [ ] **Step 1: 创建 dist_tags.go — 提供 dist-tags 只读查询方法**

```go
package registry

import (
	"context"
	"fmt"

	"github.com/scagogogo/npm-skills/pkg/models"
)

// GetDistTags 获取指定包的所有分发标签（dist-tags）
//
// Dist-tags 是 NPM 的版本别名机制，最常见的有 "latest"（最新稳定版）、
// "next"（下一个版本）、"beta" 等。此方法返回包的所有 dist-tags 映射。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - packageName: 要查询的包名称
//
// 返回值:
//   - map[string]string: 标签名到版本号的映射，如 {"latest": "18.2.0", "next": "19.0.0-rc.1"}
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	registry := NewRegistry()
//	ctx := context.Background()
//	tags, err := registry.GetDistTags(ctx, "react")
//	if err != nil {
//	    // 处理错误
//	}
//	fmt.Println("Latest:", tags["latest"])
//	fmt.Println("Next:", tags["next"])
func (x *Registry) GetDistTags(ctx context.Context, packageName string) (map[string]string, error) {
	// 使用已有的 GetPackageInformation 方法提取 dist-tags
	pkg, err := x.GetPackageInformation(ctx, packageName)
	if err != nil {
		return nil, fmt.Errorf("failed to get dist-tags for '%s': %w", packageName, err)
	}
	return pkg.DistTags, nil
}

// GetDistTagsAbbreviated 获取指定包的 dist-tags（使用精简 API）
//
// 使用 NPM 的 dist-tags 专用端点，只返回 dist-tags 数据，
// 不获取完整包信息，速度更快。
//
// 参数:
//   - ctx: 上下文
//   - packageName: 包名称（scoped 包需使用 URL 编码，如 @nestjs/core → @nestjs%2Fcore）
//
// 返回值:
//   - map[string]string: 标签名到版本号的映射
//   - error: 如果请求失败则返回错误
func (x *Registry) GetDistTagsAbbreviated(ctx context.Context, packageName string) (map[string]string, error) {
	targetUrl := fmt.Sprintf("%s/-/package/%s/dist-tags", x.options.RegistryURL, packageName)
	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to get dist-tags for '%s': %w", packageName, err)
	}
	return unmarshalJson[map[string]string](bytes)
}
```

- [ ] **Step 2: 验证编译**
Run: `cd /home/cc11001100/github/scagogogo/npm-skills && go build ./pkg/registry/`
Expected:
  - Exit code: 0

- [ ] **Step 3: 提交**
Run: `cd /home/cc11001100/github/scagogogo/npm-skills && git add pkg/registry/dist_tags.go && git commit -m "feat(registry): add GetDistTags and GetDistTagsAbbreviated methods"`

---

### Task 5: Abbreviated Package Metadata + WhoAmI

**Depends on:** Task 1
**Files:**
- Create: `pkg/registry/abbreviated.go`
- Create: `pkg/registry/whoami.go`

- [ ] **Step 1: 创建 abbreviated.go — 支持精简包元数据查询**

```go
package registry

import (
	"context"
	"fmt"

	"github.com/crawler-go-go-go/go-requests"
	"github.com/scagogogo/npm-skills/pkg/models"
)

// GetAbbreviatedPackageInformation 获取指定包的精简元数据
//
// 使用 NPM 的 install-v1 Accept header，返回的元数据比完整包信息小得多
// （完整信息可达 10MB+，精简版通常几 KB），适合只需要版本列表和 dist-tags 的场景。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - packageName: 要查询的包名称
//
// 返回值:
//   - *models.Package: 精简的包信息（可能缺少 README、maintainers 等字段）
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	registry := NewRegistry()
//	ctx := context.Background()
//	pkg, err := registry.GetAbbreviatedPackageInformation(ctx, "react")
//	if err != nil {
//	    // 处理错误
//	}
//	fmt.Println("Latest:", pkg.DistTags["latest"])
func (x *Registry) GetAbbreviatedPackageInformation(ctx context.Context, packageName string) (*models.Package, error) {
	targetUrl := fmt.Sprintf("%s/%s", x.options.RegistryURL, packageName)
	opts := requests.NewOptions[any, []byte](targetUrl, requests.BytesResponseHandler())
	opts.AppendRequestSetting(requests.RequestSettingHeader("Accept", "application/vnd.npm.install-v1+json"))
	if x.options.Proxy != "" {
		opts.AppendRequestSetting(requests.RequestSettingProxy(x.options.Proxy))
	}
	if x.options.Token != "" {
		opts.AppendRequestSetting(requests.RequestSettingHeader("Authorization", "Bearer "+x.options.Token))
	}
	bytes, err := requests.SendRequest[any, []byte](ctx, opts)
	if err != nil {
		return nil, err
	}
	return unmarshalJson[*models.Package](bytes)
}
```

- [ ] **Step 2: 创建 whoami.go — 检查当前认证状态**

```go
package registry

import (
	"context"
	"fmt"
)

// WhoAmI 检查当前认证身份
//
// 使用 NPM 的 /-/whoami 端点验证 Token 是否有效并返回用户名。
// 如果未设置 Token 或 Token 无效，将返回错误。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//
// 返回值:
//   - string: 当前认证用户名
//   - error: 如果未认证或请求失败则返回错误
//
// 使用示例:
//
//	options := NewOptions().SetToken("npm_xxxxx")
//	registry := NewRegistry(options)
//	ctx := context.Background()
//	username, err := registry.WhoAmI(ctx)
//	if err != nil {
//	    fmt.Println("未认证或 Token 无效:", err)
//	} else {
//	    fmt.Println("当前用户:", username)
//	}
func (x *Registry) WhoAmI(ctx context.Context) (string, error) {
	if x.options.Token == "" {
		return "", fmt.Errorf("no token set: configure with options.SetToken() before calling WhoAmI")
	}
	targetUrl := fmt.Sprintf("%s/-/whoami", x.options.RegistryURL)
	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return "", fmt.Errorf("whoami request failed: %w", err)
	}
	var result struct {
		Username string `json:"username"`
	}
	if err := json.Unmarshal(bytes, &result); err != nil {
		return "", fmt.Errorf("whoami response parse failed: %w", err)
	}
	if result.Username == "" {
		return "", fmt.Errorf("authentication failed: empty username in response")
	}
	return result.Username, nil
}
```

- [ ] **Step 3: 验证编译**
Run: `cd /home/cc11001100/github/scagogogo/npm-skills && go build ./pkg/registry/`
Expected:
  - Exit code: 0

- [ ] **Step 4: 全量回归测试**
Run: `cd /home/cc11001100/github/scagogogo/npm-skills && go test ./... -count=1 2>&1 | grep -E "^(ok|FAIL|\?)" `
Expected:
  - Output does NOT contain: "FAIL"

- [ ] **Step 5: 提交**
Run: `cd /home/cc11001100/github/scagogogo/npm-skills && git add pkg/registry/abbreviated.go pkg/registry/whoami.go && git commit -m "feat(registry): add GetAbbreviatedPackageInformation and WhoAmI methods"`
