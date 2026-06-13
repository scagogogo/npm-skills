package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/scagogogo/npm-skills/pkg/models"
)

// ListTokens 列出所有NPM访问令牌
//
// 需要认证Token。返回当前认证用户的所有访问令牌列表。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//
// 返回值:
//   - []models.Token: 令牌列表
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	options := NewOptions().SetToken("npm_xxxxx")
//	registry := NewRegistry(options)
//	ctx := context.Background()
//	tokens, err := registry.ListTokens(ctx)
func (x *Registry) ListTokens(ctx context.Context) ([]models.Token, error) {
	if err := x.requireToken(); err != nil {
		return nil, err
	}

	targetUrl := fmt.Sprintf("%s/-/npm/v1/tokens", x.options.RegistryURL)
	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to list tokens: %w", err)
	}

	var result struct {
		Objects []models.Token `json:"objects"`
		Total   int            `json:"total"`
	}
	if err := json.Unmarshal(bytes, &result); err != nil {
		return nil, fmt.Errorf("failed to parse tokens: %w", err)
	}
	return result.Objects, nil
}

// GetToken 获取指定NPM访问令牌
//
// 需要认证Token。返回指定 ID 的访问令牌详情。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - tokenID: 令牌ID
//
// 返回值:
//   - *models.Token: 令牌详情
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	options := NewOptions().SetToken("npm_xxxxx")
//	registry := NewRegistry(options)
//	ctx := context.Background()
//	token, err := registry.GetToken(ctx, "token-id-here")
func (x *Registry) GetToken(ctx context.Context, tokenID string) (*models.Token, error) {
	if err := x.requireToken(); err != nil {
		return nil, err
	}

	targetUrl := fmt.Sprintf("%s/-/npm/v1/tokens/%s", x.options.RegistryURL, tokenID)
	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to get token '%s': %w", tokenID, err)
	}
	return unmarshalJson[*models.Token](bytes)
}

// CreateToken 创建新的NPM访问令牌
//
// 需要认证Token。创建一个新的访问令牌，可设置为只读或读写权限。
// 还可配置 IP 白名单（CIDR）以限制令牌的使用范围。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - opts: 令牌创建参数，包含密码、是否只读、IP白名单等
//
// 返回值:
//   - *models.Token: 新创建的令牌
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	options := NewOptions().SetToken("npm_xxxxx")
//	registry := NewRegistry(options)
//	ctx := context.Background()
//	token, err := registry.CreateToken(ctx, &models.TokenCreation{
//	    Password: "my-password",
//	    Readonly: true,
//	})
func (x *Registry) CreateToken(ctx context.Context, opts *models.TokenCreation) (*models.Token, error) {
	if err := x.requireToken(); err != nil {
		return nil, err
	}

	targetUrl := fmt.Sprintf("%s/-/npm/v1/tokens", x.options.RegistryURL)
	bytes, err := x.sendJSON(ctx, http.MethodPost, targetUrl, opts, http.StatusOK, http.StatusCreated)
	if err != nil {
		return nil, fmt.Errorf("failed to create token: %w", err)
	}
	return unmarshalJson[*models.Token](bytes)
}

// DeleteToken 删除指定的NPM访问令牌
//
// 需要认证Token。撤销指定 ID 的访问令牌，令牌将立即失效。
// 注意：不能删除当前正在使用的令牌。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - tokenID: 要删除的令牌ID
//
// 返回值:
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	options := NewOptions().SetToken("npm_xxxxx")
//	registry := NewRegistry(options)
//	ctx := context.Background()
//	err := registry.DeleteToken(ctx, "token-id-to-delete")
func (x *Registry) DeleteToken(ctx context.Context, tokenID string) error {
	if err := x.requireToken(); err != nil {
		return err
	}

	targetUrl := fmt.Sprintf("%s/-/npm/v1/tokens/%s", x.options.RegistryURL, tokenID)
	_, err := x.sendRequest(ctx, http.MethodDelete, targetUrl, nil, http.StatusOK, http.StatusNoContent)
	if err != nil {
		return fmt.Errorf("failed to delete token '%s': %w", tokenID, err)
	}
	return nil
}
