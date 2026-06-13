package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/scagogogo/npm-skills/pkg/models"
)

// ListHooks 列出所有 Webhook
//
// 需要认证Token。返回当前认证用户拥有的所有 Webhook 列表。
// 可选按包名过滤。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - opts: 查询参数，可选按包名过滤和分页
//
// 返回值:
//   - []models.Hook: Webhook 列表
//   - error: 如果请求失败则返回错误
func (x *Registry) ListHooks(ctx context.Context, opts models.HookListOptions) ([]models.Hook, error) {
	if err := x.requireToken(); err != nil {
		return nil, err
	}

	params := url.Values{}
	if opts.Package != "" {
		params.Set("package", opts.Package)
	}
	if opts.Page > 0 {
		params.Set("page", fmt.Sprintf("%d", opts.Page))
	}
	if opts.PerPage > 0 {
		params.Set("per_page", fmt.Sprintf("%d", opts.PerPage))
	}

	targetUrl := fmt.Sprintf("%s/-/npm/v1/hooks", x.options.RegistryURL)
	if len(params) > 0 {
		targetUrl = fmt.Sprintf("%s?%s", targetUrl, params.Encode())
	}

	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to list hooks: %w", err)
	}

	var result struct {
		Objects []models.Hook `json:"objects"`
		Total   int           `json:"total"`
	}
	if err := json.Unmarshal(bytes, &result); err != nil {
		return nil, fmt.Errorf("failed to parse hooks: %w", err)
	}
	return result.Objects, nil
}

// GetHook 获取指定 Webhook
//
// 需要认证Token。返回指定 ID 的 Webhook 详情。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - hookID: Webhook ID
//
// 返回值:
//   - *models.Hook: Webhook 详情
//   - error: 如果请求失败则返回错误
func (x *Registry) GetHook(ctx context.Context, hookID string) (*models.Hook, error) {
	if err := x.requireToken(); err != nil {
		return nil, err
	}

	targetUrl := fmt.Sprintf("%s/-/npm/v1/hooks/%s", x.options.RegistryURL, hookID)
	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to get hook '%s': %w", hookID, err)
	}
	return unmarshalJson[*models.Hook](bytes)
}

// CreateHook 创建 Webhook
//
// 需要认证Token。创建一个新的 Webhook，当指定包发布新版本时，
// 将向指定的 endpoint 发送 HTTP POST 请求。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - hook: Webhook 创建参数
//
// 返回值:
//   - *models.Hook: 创建的 Webhook
//   - error: 如果请求失败则返回错误
func (x *Registry) CreateHook(ctx context.Context, hook *models.HookCreation) (*models.Hook, error) {
	if err := x.requireToken(); err != nil {
		return nil, err
	}

	targetUrl := fmt.Sprintf("%s/-/npm/v1/hooks", x.options.RegistryURL)
	bytes, err := x.sendJSON(ctx, http.MethodPost, targetUrl, hook, http.StatusOK, http.StatusCreated)
	if err != nil {
		return nil, fmt.Errorf("failed to create hook: %w", err)
	}
	return unmarshalJson[*models.Hook](bytes)
}

// UpdateHook 更新 Webhook
//
// 需要认证Token。更新指定 Webhook 的配置，如 endpoint URL 或 secret。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - hookID: Webhook ID
//   - hook: Webhook 更新参数
//
// 返回值:
//   - *models.Hook: 更新后的 Webhook
//   - error: 如果请求失败则返回错误
func (x *Registry) UpdateHook(ctx context.Context, hookID string, hook *models.HookUpdate) (*models.Hook, error) {
	if err := x.requireToken(); err != nil {
		return nil, err
	}

	targetUrl := fmt.Sprintf("%s/-/npm/v1/hooks/%s", x.options.RegistryURL, hookID)
	bytes, err := x.sendJSON(ctx, http.MethodPut, targetUrl, hook, http.StatusOK)
	if err != nil {
		return nil, fmt.Errorf("failed to update hook '%s': %w", hookID, err)
	}
	return unmarshalJson[*models.Hook](bytes)
}

// DeleteHook 删除 Webhook
//
// 需要认证Token。删除指定的 Webhook。删除后不会再发送通知。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - hookID: Webhook ID
//
// 返回值:
//   - error: 如果请求失败则返回错误
func (x *Registry) DeleteHook(ctx context.Context, hookID string) error {
	if err := x.requireToken(); err != nil {
		return err
	}

	targetUrl := fmt.Sprintf("%s/-/npm/v1/hooks/%s", x.options.RegistryURL, hookID)
	_, err := x.sendRequest(ctx, http.MethodDelete, targetUrl, nil, http.StatusOK, http.StatusNoContent)
	if err != nil {
		return fmt.Errorf("failed to delete hook '%s': %w", hookID, err)
	}
	return nil
}
