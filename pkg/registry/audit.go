package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/scagogogo/npm-skills/pkg/models"
)

// BulkAudit 批量审计依赖安全性
//
// 提交一组包名和版本范围，返回匹配的安全公告。
// 这是 npm audit 使用的底层端点。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - advisories: 包名到版本范围的映射，如 {"lodash": "<4.17.12", "express": "<4.17.3"}
//
// 返回值:
//   - map[string][]models.Advisory: 包名到安全公告列表的映射
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	registry := NewRegistry()
//	ctx := context.Background()
//	advisories, err := registry.BulkAudit(ctx, map[string][]string{
//	    "lodash": {"<4.17.12"},
//	})
func (x *Registry) BulkAudit(ctx context.Context, advisories map[string][]string) (map[string][]models.Advisory, error) {
	targetUrl := fmt.Sprintf("%s/-/npm/v1/security/advisories/bulk", x.options.RegistryURL)
	bytes, err := x.sendJSON(ctx, http.MethodPost, targetUrl, advisories, http.StatusOK, http.StatusCreated, http.StatusNoContent)
	if err != nil {
		return nil, fmt.Errorf("failed to bulk audit: %w", err)
	}
	return unmarshalJson[map[string][]models.Advisory](bytes)
}

// GetAdvisory 获取特定安全公告
//
// 通过公告 ID 获取安全公告的详细信息。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - advisoryID: 安全公告ID
//
// 返回值:
//   - *models.Advisory: 安全公告详情
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	registry := NewRegistry()
//	ctx := context.Background()
//	advisory, err := registry.GetAdvisory(ctx, 1234)
func (x *Registry) GetAdvisory(ctx context.Context, advisoryID int) (*models.Advisory, error) {
	targetUrl := fmt.Sprintf("%s/-/npm/v1/security/advisories/%d", x.options.RegistryURL, advisoryID)
	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to get advisory %d: %w", advisoryID, err)
	}
	return unmarshalJson[*models.Advisory](bytes)
}

// ListAdvisories 列出安全公告
//
// 返回安全公告列表，支持分页和按包名过滤。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - opts: 查询参数，包含分页和过滤条件
//
// 返回值:
//   - []models.Advisory: 安全公告列表
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	registry := NewRegistry()
//	ctx := context.Background()
//	advisories, err := registry.ListAdvisories(ctx, models.AdvisoryListOptions{
//	    PerPage: 20,
//	    AffectedPackage: "lodash",
//	})
func (x *Registry) ListAdvisories(ctx context.Context, opts models.AdvisoryListOptions) ([]models.Advisory, error) {
	params := url.Values{}
	if opts.Page > 0 {
		params.Set("page", fmt.Sprintf("%d", opts.Page))
	}
	if opts.PerPage > 0 {
		params.Set("per_page", fmt.Sprintf("%d", opts.PerPage))
	}
	if opts.AffectedPackage != "" {
		params.Set("affected_package", opts.AffectedPackage)
	}

	targetUrl := fmt.Sprintf("%s/-/npm/v1/security/advisories", x.options.RegistryURL)
	if len(params) > 0 {
		targetUrl = fmt.Sprintf("%s?%s", targetUrl, params.Encode())
	}

	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to list advisories: %w", err)
	}

	var result struct {
		Objects []models.Advisory `json:"objects"`
		Total   int               `json:"total"`
	}
	if err := json.Unmarshal(bytes, &result); err != nil {
		return nil, fmt.Errorf("failed to parse advisories: %w", err)
	}
	return result.Objects, nil
}

// QuickAudit 快速审计
//
// 提交依赖列表进行快速安全审计，返回漏洞统计。
// 比 BulkAudit 更快，但返回的信息较少。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - payload: 审计请求，包含依赖列表
//
// 返回值:
//   - *models.QuickAuditResult: 审计结果，包含漏洞统计
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	registry := NewRegistry()
//	ctx := context.Background()
//	result, err := registry.QuickAudit(ctx, &models.QuickAuditRequest{
//	    Dependencies: map[string]string{
//	        "lodash": "4.17.11",
//	        "express": "4.17.1",
//	    },
//	})
func (x *Registry) QuickAudit(ctx context.Context, payload *models.QuickAuditRequest) (*models.QuickAuditResult, error) {
	targetUrl := fmt.Sprintf("%s/-/npm/v1/security/audits/quick", x.options.RegistryURL)
	bytes, err := x.sendJSON(ctx, http.MethodPost, targetUrl, payload, http.StatusOK, http.StatusCreated)
	if err != nil {
		return nil, fmt.Errorf("failed to quick audit: %w", err)
	}
	return unmarshalJson[*models.QuickAuditResult](bytes)
}
