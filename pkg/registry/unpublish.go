package registry

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// UnpublishPackage 完整取消发布一个NPM包（危险操作）
//
// 需要认证Token。此操作将删除整个包文档，包括所有版本。
// 大多数 Registry 只允许在发布后 72 小时内取消发布。
// 这是一个不可逆操作，请谨慎使用。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - packageName: 要取消发布的包名称
//
// 返回值:
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	options := NewOptions().SetToken("npm_xxxxx")
//	registry := NewRegistry(options)
//	ctx := context.Background()
//	err := registry.UnpublishPackage(ctx, "my-package")
func (x *Registry) UnpublishPackage(ctx context.Context, packageName string) error {
	if err := x.requireToken(); err != nil {
		return err
	}

	// 获取当前 _rev
	rev, err := x.getRev(ctx, packageName)
	if err != nil {
		return fmt.Errorf("failed to get package revision for unpublish: %w", err)
	}
	if rev == "" {
		return fmt.Errorf("package '%s' not found or already unpublished", packageName)
	}

	targetUrl := fmt.Sprintf("%s/%s/-rev/%s", x.options.RegistryURL, url.PathEscape(packageName), rev)
	_, err = x.sendRequest(ctx, http.MethodDelete, targetUrl, nil, http.StatusOK, http.StatusCreated)
	if err != nil {
		return fmt.Errorf("failed to unpublish package '%s': %w", packageName, err)
	}
	return nil
}

// UnpublishPackageVersion 取消发布指定包的特定版本
//
// 需要认证Token。此操作将删除包的特定版本，但保留其他版本。
// 需要先获取当前包文档的 _rev，然后通过修改文档来删除指定版本。
// 大多数 Registry 只允许在发布后 72 小时内取消发布。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - packageName: 包名称
//   - version: 要取消发布的版本号
//
// 返回值:
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	options := NewOptions().SetToken("npm_xxxxx")
//	registry := NewRegistry(options)
//	ctx := context.Background()
//	err := registry.UnpublishPackageVersion(ctx, "my-package", "1.0.0")
func (x *Registry) UnpublishPackageVersion(ctx context.Context, packageName, version string) error {
	if err := x.requireToken(); err != nil {
		return err
	}

	// 获取当前包的完整文档
	pkg, err := x.GetPackageInformation(ctx, packageName)
	if err != nil {
		return fmt.Errorf("failed to get package information for unpublish: %w", err)
	}

	// 删除指定版本
	if _, exists := pkg.Versions[version]; !exists {
		return fmt.Errorf("version '%s' not found in package '%s'", version, packageName)
	}
	delete(pkg.Versions, version)

	// 更新 _rev
	rev, err := x.getRev(ctx, packageName)
	if err != nil {
		return fmt.Errorf("failed to get package revision: %w", err)
	}
	pkg.Rev = rev

	// PUT 更新后的文档
	targetUrl := fmt.Sprintf("%s/%s", x.options.RegistryURL, url.PathEscape(packageName))
	_, err = x.sendJSON(ctx, http.MethodPut, targetUrl, pkg, http.StatusOK, http.StatusCreated)
	if err != nil {
		return fmt.Errorf("failed to unpublish version '%s' of package '%s': %w", version, packageName, err)
	}
	return nil
}
