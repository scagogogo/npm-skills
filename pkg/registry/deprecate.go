package registry

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// DeprecateVersion 标记指定包的特定版本为已弃用
//
// 需要认证Token。此操作不会删除版本，而是在版本元数据中设置 deprecated 字段，
// 用户在安装时会看到弃用警告。这是推荐的方式来标记不再维护的版本，
// 而不是使用 UnpublishPackageVersion。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - packageName: 包名称
//   - version: 要标记为弃用的版本号
//   - message: 弃用消息，通常会建议用户迁移到其他版本，如 "Use v2.0.0 instead"
//
// 返回值:
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	options := NewOptions().SetToken("npm_xxxxx")
//	registry := NewRegistry(options)
//	ctx := context.Background()
//	err := registry.DeprecateVersion(ctx, "my-package", "1.0.0", "Use v2.0.0 instead, see migration guide at https://example.com/migrate")
func (x *Registry) DeprecateVersion(ctx context.Context, packageName, version, message string) error {
	if err := x.requireToken(); err != nil {
		return err
	}

	// 获取当前包的完整文档
	pkg, err := x.GetPackageInformation(ctx, packageName)
	if err != nil {
		return fmt.Errorf("failed to get package information for deprecate: %w", err)
	}

	// 检查版本是否存在
	ver, exists := pkg.Versions[version]
	if !exists {
		return fmt.Errorf("version '%s' not found in package '%s'", version, packageName)
	}

	// 设置 deprecated 字段
	ver.Deprecated = message
	pkg.Versions[version] = ver

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
		return fmt.Errorf("failed to deprecate version '%s' of package '%s': %w", version, packageName, err)
	}
	return nil
}
