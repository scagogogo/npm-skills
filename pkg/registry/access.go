package registry

import (
	"context"
	"fmt"
	"net/http"

	"github.com/scagogogo/npm-skills/pkg/models"
)

// GetPackageAccess 获取包的访问权限设置
//
// 需要认证Token。返回包的访问级别和权限信息。
// 对于私有包，此操作返回谁有读写权限。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - packageName: 包名称
//
// 返回值:
//   - *models.PackageAccess: 访问权限信息
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	options := NewOptions().SetToken("npm_xxxxx")
//	registry := NewRegistry(options)
//	ctx := context.Background()
//	access, err := registry.GetPackageAccess(ctx, "my-package")
func (x *Registry) GetPackageAccess(ctx context.Context, packageName string) (*models.PackageAccess, error) {
	if err := x.requireToken(); err != nil {
		return nil, err
	}

	targetUrl := fmt.Sprintf("%s/-/package/%s/access", x.options.RegistryURL, encodePackageName(packageName))
	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to get access for '%s': %w", packageName, err)
	}
	return unmarshalJson[*models.PackageAccess](bytes)
}

// SetPackageAccess 设置包的访问权限
//
// 需要认证Token。可将包设置为公开(public)或私有(restricted)。
// 注意：将包从 public 改为 restricted 后，未授权用户将无法访问。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - packageName: 包名称
//   - access: 访问级别更新请求，Access 字段为 "public" 或 "restricted"
//
// 返回值:
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	options := NewOptions().SetToken("npm_xxxxx")
//	registry := NewRegistry(options)
//	ctx := context.Background()
//	err := registry.SetPackageAccess(ctx, "my-package", &models.PackageAccessUpdate{Access: "restricted"})
func (x *Registry) SetPackageAccess(ctx context.Context, packageName string, access *models.PackageAccessUpdate) error {
	if err := x.requireToken(); err != nil {
		return err
	}

	targetUrl := fmt.Sprintf("%s/-/package/%s/access", x.options.RegistryURL, encodePackageName(packageName))
	_, err := x.sendJSON(ctx, http.MethodPost, targetUrl, access, http.StatusOK, http.StatusCreated)
	if err != nil {
		return fmt.Errorf("failed to set access for '%s': %w", packageName, err)
	}
	return nil
}

// ListCollaborators 列出包的协作者
//
// 需要认证Token。返回有权访问该包的所有用户和团队列表及其权限级别。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - packageName: 包名称
//
// 返回值:
//   - []models.Collaborator: 协作者列表
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	options := NewOptions().SetToken("npm_xxxxx")
//	registry := NewRegistry(options)
//	ctx := context.Background()
//	collabs, err := registry.ListCollaborators(ctx, "my-package")
func (x *Registry) ListCollaborators(ctx context.Context, packageName string) ([]models.Collaborator, error) {
	if err := x.requireToken(); err != nil {
		return nil, err
	}

	targetUrl := fmt.Sprintf("%s/-/package/%s/collaborators", x.options.RegistryURL, encodePackageName(packageName))
	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to list collaborators for '%s': %w", packageName, err)
	}
	return unmarshalJson[[]models.Collaborator](bytes)
}

// GrantAccess 授予用户或团队对包的访问权限
//
// 需要认证Token。可以给单个用户或组织团队授予读或写权限。
// 对于团队，格式为 "<org>:<team>"。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - packageName: 包名称
//   - user: 用户名或团队名（格式: "<org>:<team>"）
//   - permission: 权限级别，PermissionRead 或 PermissionWrite
//
// 返回值:
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	options := NewOptions().SetToken("npm_xxxxx")
//	registry := NewRegistry(options)
//	ctx := context.Background()
//	// 给用户授予写权限
//	err := registry.GrantAccess(ctx, "my-package", "username", models.PermissionWrite)
//	// 给团队授予读权限
//	err = registry.GrantAccess(ctx, "my-package", "myorg:devteam", models.PermissionRead)
func (x *Registry) GrantAccess(ctx context.Context, packageName, user string, permission models.Permission) error {
	if err := x.requireToken(); err != nil {
		return err
	}

	targetUrl := fmt.Sprintf("%s/-/package/%s/collaborators", x.options.RegistryURL, encodePackageName(packageName))
	payload := map[string]string{
		user: string(permission),
	}
	_, err := x.sendJSON(ctx, http.MethodPost, targetUrl, payload, http.StatusOK, http.StatusCreated)
	if err != nil {
		return fmt.Errorf("failed to grant access for '%s' on '%s': %w", user, packageName, err)
	}
	return nil
}

// RevokeAccess 撤销用户对包的访问权限
//
// 需要认证Token。从包的协作者列表中移除指定用户。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - packageName: 包名称
//   - user: 要撤销权限的用户名
//
// 返回值:
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	options := NewOptions().SetToken("npm_xxxxx")
//	registry := NewRegistry(options)
//	ctx := context.Background()
//	err := registry.RevokeAccess(ctx, "my-package", "username")
func (x *Registry) RevokeAccess(ctx context.Context, packageName, user string) error {
	if err := x.requireToken(); err != nil {
		return err
	}

	targetUrl := fmt.Sprintf("%s/-/package/%s/collaborators/%s", x.options.RegistryURL, encodePackageName(packageName), user)
	_, err := x.sendRequest(ctx, http.MethodDelete, targetUrl, nil, http.StatusOK, http.StatusNoContent)
	if err != nil {
		return fmt.Errorf("failed to revoke access for '%s' on '%s': %w", user, packageName, err)
	}
	return nil
}
