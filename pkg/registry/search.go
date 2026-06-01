package registry

import (
	"context"
	"fmt"
	"net/url"

	"github.com/scagogogo/npm-crawler/pkg/models"
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