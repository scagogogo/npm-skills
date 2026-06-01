package registry

import (
	"context"
	"encoding/json"
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