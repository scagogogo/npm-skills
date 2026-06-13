package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// StarPackage 收藏一个NPM包
//
// 需要认证Token。将指定包添加到用户的收藏列表。
// 收藏信息会体现在包的 users 字段中。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - packageName: 要收藏的包名称
//
// 返回值:
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	options := NewOptions().SetToken("npm_xxxxx")
//	registry := NewRegistry(options)
//	ctx := context.Background()
//	err := registry.StarPackage(ctx, "react")
func (x *Registry) StarPackage(ctx context.Context, packageName string) error {
	if err := x.requireToken(); err != nil {
		return err
	}

	// NPM star 操作需要先获取当前包信息，然后在 users 字段中设置自己的用户名
	pkg, err := x.GetPackageInformation(ctx, packageName)
	if err != nil {
		return fmt.Errorf("failed to get package for starring: %w", err)
	}

	// 获取当前用户名
	username, err := x.WhoAmI(ctx)
	if err != nil {
		return fmt.Errorf("failed to identify user for starring: %w", err)
	}

	// 设置 users 字段
	if pkg.Users == nil {
		pkg.Users = make(map[string]bool)
	}
	pkg.Users[username] = true

	// 更新 _rev
	rev, err := x.getRev(ctx, packageName)
	if err != nil {
		return fmt.Errorf("failed to get package revision: %w", err)
	}
	pkg.Rev = rev

	targetUrl := fmt.Sprintf("%s/%s", x.options.RegistryURL, url.PathEscape(packageName))
	_, err = x.sendJSON(ctx, http.MethodPut, targetUrl, pkg, http.StatusOK, http.StatusCreated)
	if err != nil {
		return fmt.Errorf("failed to star package '%s': %w", packageName, err)
	}
	return nil
}

// UnstarPackage 取消收藏一个NPM包
//
// 需要认证Token。将指定包从用户的收藏列表中移除。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - packageName: 要取消收藏的包名称
//
// 返回值:
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	options := NewOptions().SetToken("npm_xxxxx")
//	registry := NewRegistry(options)
//	ctx := context.Background()
//	err := registry.UnstarPackage(ctx, "react")
func (x *Registry) UnstarPackage(ctx context.Context, packageName string) error {
	if err := x.requireToken(); err != nil {
		return err
	}

	pkg, err := x.GetPackageInformation(ctx, packageName)
	if err != nil {
		return fmt.Errorf("failed to get package for unstarring: %w", err)
	}

	username, err := x.WhoAmI(ctx)
	if err != nil {
		return fmt.Errorf("failed to identify user for unstarring: %w", err)
	}

	// 从 users 字段中移除
	if pkg.Users != nil {
		delete(pkg.Users, username)
	}

	rev, err := x.getRev(ctx, packageName)
	if err != nil {
		return fmt.Errorf("failed to get package revision: %w", err)
	}
	pkg.Rev = rev

	targetUrl := fmt.Sprintf("%s/%s", x.options.RegistryURL, url.PathEscape(packageName))
	_, err = x.sendJSON(ctx, http.MethodPut, targetUrl, pkg, http.StatusOK, http.StatusCreated)
	if err != nil {
		return fmt.Errorf("failed to unstar package '%s': %w", packageName, err)
	}
	return nil
}

// GetStarredByUser 获取用户收藏的所有包
//
// 返回指定用户收藏的包名称列表。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - username: 用户名
//
// 返回值:
//   - []string: 用户收藏的包名称列表
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	registry := NewRegistry()
//	ctx := context.Background()
//	packages, err := registry.GetStarredByUser(ctx, "username")
func (x *Registry) GetStarredByUser(ctx context.Context, username string) ([]string, error) {
	targetUrl := fmt.Sprintf("%s/-/_view/starredByUser?key=%s", x.options.RegistryURL, url.QueryEscape(fmt.Sprintf("%q", username)))
	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to get starred packages for user '%s': %w", username, err)
	}

	var result struct {
		Rows []struct {
			Value string `json:"value"`
		} `json:"rows"`
	}
	if err := json.Unmarshal(bytes, &result); err != nil {
		return nil, fmt.Errorf("failed to parse starred packages: %w", err)
	}

	packages := make([]string, 0, len(result.Rows))
	for _, row := range result.Rows {
		packages = append(packages, row.Value)
	}
	return packages, nil
}

// GetStarredByPackage 获取收藏了指定包的所有用户
//
// 返回收藏了指定包的用户名列表。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - packageName: 包名称
//
// 返回值:
//   - []string: 收藏了该包的用户名列表
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	registry := NewRegistry()
//	ctx := context.Background()
//	users, err := registry.GetStarredByPackage(ctx, "react")
func (x *Registry) GetStarredByPackage(ctx context.Context, packageName string) ([]string, error) {
	targetUrl := fmt.Sprintf("%s/-/_view/starredByPackage?key=%s", x.options.RegistryURL, url.QueryEscape(fmt.Sprintf("%q", packageName)))
	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to get stargazers for package '%s': %w", packageName, err)
	}

	var result struct {
		Rows []struct {
			Value string `json:"value"`
		} `json:"rows"`
	}
	if err := json.Unmarshal(bytes, &result); err != nil {
		return nil, fmt.Errorf("failed to parse stargazers: %w", err)
	}

	users := make([]string, 0, len(result.Rows))
	for _, row := range result.Rows {
		users = append(users, row.Value)
	}
	return users, nil
}
