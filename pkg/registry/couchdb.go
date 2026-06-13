package registry

import (
	"context"
	"fmt"
	"net/url"

	"github.com/scagogogo/npm-skills/pkg/models"
)

// GetChanges 获取 Registry 的变更 Feed
//
// 返回 CouchDB 的 _changes Feed，包含所有文档的变更记录。
// 主要用于镜像同步和增量数据获取。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - opts: 查询参数，可选 since 和 limit
//
// 返回值:
//   - *models.ChangesResult: 变更结果
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	registry := NewRegistry()
//	ctx := context.Background()
//	// 获取最近10条变更
//	changes, err := registry.GetChanges(ctx, models.ChangesOptions{Limit: 10})
//	// 增量获取
//	changes, err = registry.GetChanges(ctx, models.ChangesOptions{Since: changes.LastSeq})
func (x *Registry) GetChanges(ctx context.Context, opts models.ChangesOptions) (*models.ChangesResult, error) {
	params := url.Values{}
	if opts.Since != "" {
		params.Set("since", opts.Since)
	}
	if opts.Limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", opts.Limit))
	}
	if opts.IncludeDocs {
		params.Set("include_docs", "true")
	}

	targetUrl := fmt.Sprintf("%s/_changes", x.options.RegistryURL)
	if len(params) > 0 {
		targetUrl = fmt.Sprintf("%s?%s", targetUrl, params.Encode())
	}

	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to get changes feed: %w", err)
	}
	return unmarshalJson[*models.ChangesResult](bytes)
}

// GetAllDocs 获取所有文档 ID 列表
//
// 返回 CouchDB 的 _all_docs 结果，包含所有文档的 ID 和修订版本。
// 可选返回完整文档内容（设置 IncludeDocs=true），但响应可能非常大。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - opts: 查询参数，可选范围查询和分页
//
// 返回值:
//   - *models.AllDocsResult: 文档列表结果
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	registry := NewRegistry()
//	ctx := context.Background()
//	// 获取以 @nestjs 开头的包
//	docs, err := registry.GetAllDocs(ctx, models.AllDocsOptions{
//	    StartKey: "@nestjs",
//	    EndKey:   "@nestjs￰",
//	    Limit:    20,
//	})
func (x *Registry) GetAllDocs(ctx context.Context, opts models.AllDocsOptions) (*models.AllDocsResult, error) {
	params := url.Values{}
	if opts.StartKey != "" {
		params.Set("startkey", opts.StartKey)
	}
	if opts.EndKey != "" {
		params.Set("endkey", opts.EndKey)
	}
	if opts.Limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", opts.Limit))
	}
	if opts.Skip > 0 {
		params.Set("skip", fmt.Sprintf("%d", opts.Skip))
	}
	if opts.IncludeDocs {
		params.Set("include_docs", "true")
	}
	if opts.Descending {
		params.Set("descending", "true")
	}

	targetUrl := fmt.Sprintf("%s/_all_docs", x.options.RegistryURL)
	if len(params) > 0 {
		targetUrl = fmt.Sprintf("%s?%s", targetUrl, params.Encode())
	}

	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to get all docs: %w", err)
	}
	return unmarshalJson[*models.AllDocsResult](bytes)
}

// GetView 查询 CouchDB 视图
//
// 查询 NPM Registry 上的 CouchDB 视图。
// 常用视图包括: starredByUser, starredByPackage, byKeyword, byUser 等。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - viewName: 视图名称，如 "starredByUser", "starredByPackage"
//   - opts: 查询参数
//
// 返回值:
//   - *models.ViewResult: 视图查询结果
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	registry := NewRegistry()
//	ctx := context.Background()
//	result, err := registry.GetView(ctx, "starredByUser", models.ViewOptions{
//	    Key: "\"username\"",
//	})
func (x *Registry) GetView(ctx context.Context, viewName string, opts models.ViewOptions) (*models.ViewResult, error) {
	params := url.Values{}
	if opts.Key != "" {
		params.Set("key", opts.Key)
	}
	if opts.StartKey != "" {
		params.Set("startkey", opts.StartKey)
	}
	if opts.EndKey != "" {
		params.Set("endkey", opts.EndKey)
	}
	if opts.Limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", opts.Limit))
	}
	if opts.Skip > 0 {
		params.Set("skip", fmt.Sprintf("%d", opts.Skip))
	}
	if opts.Group {
		params.Set("group", "true")
	}
	if opts.GroupLevel > 0 {
		params.Set("group_level", fmt.Sprintf("%d", opts.GroupLevel))
	}
	if opts.Descending {
		params.Set("descending", "true")
	}

	targetUrl := fmt.Sprintf("%s/-/_view/%s", x.options.RegistryURL, viewName)
	if len(params) > 0 {
		targetUrl = fmt.Sprintf("%s?%s", targetUrl, params.Encode())
	}

	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to get view '%s': %w", viewName, err)
	}
	return unmarshalJson[*models.ViewResult](bytes)
}
