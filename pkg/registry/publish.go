package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/crawler-go-go-go/go-requests"
	"github.com/scagogogo/npm-skills/pkg/models"
)

// getRev 获取指定包的当前 CouchDB _rev 值
//
// CouchDB 使用 MVCC 机制，每次更新文档都必须提供当前的 _rev 值。
// 如果包不存在（新包发布），返回空字符串。
//
// 参数:
//   - ctx: 上下文
//   - packageName: 包名称
//
// 返回值:
//   - string: 当前 _rev 值，包不存在时为空字符串
//   - error: 如果请求失败（非 404 错误）则返回错误
func (x *Registry) getRev(ctx context.Context, packageName string) (string, error) {
	targetUrl := fmt.Sprintf("%s/%s", x.options.RegistryURL, url.PathEscape(packageName))

	// Custom request that accepts 404 (package doesn't exist) in addition to 200
	options := requests.NewOptions[any, []byte](targetUrl, requests.BytesResponseHandler(http.StatusOK, http.StatusNotFound))
	if x.options.Proxy != "" {
		options.AppendRequestSetting(requests.RequestSettingProxy(x.options.Proxy))
	}
	if x.options.Token != "" {
		options.AppendRequestSetting(requestSettingHeader("Authorization", "Bearer "+x.options.Token))
	}

	bytes, err := requests.SendRequest[any, []byte](ctx, options)
	if err != nil {
		return "", fmt.Errorf("failed to check package existence for '%s': %w", packageName, err)
	}

	// If no bytes or empty response, package doesn't exist yet
	if len(bytes) == 0 {
		return "", nil
	}

	var doc struct {
		Rev string `json:"_rev"`
	}
	if err := json.Unmarshal(bytes, &doc); err != nil {
		// Can't parse - likely a 404 error page, treat as package not existing
		return "", nil
	}
	return doc.Rev, nil
}

// PublishPackage 发布一个NPM包到Registry
//
// 需要认证Token。发布流程：
//  1. 获取当前包文档的 _rev（如果包已存在）
//  2. 将 _rev 合并到包文档中
//  3. PUT 完整包文档到 Registry
//
// 注意：这是底层发布方法，需要手动构造完整的 Package 文档。
// 大多数情况下，应使用 PublishPackageFromTarball 方法。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - pkg: 完整的包文档，必须包含至少一个版本
//
// 返回值:
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	options := NewOptions().SetToken("npm_xxxxx")
//	registry := NewRegistry(options)
//	ctx := context.Background()
//	pkg := &models.Package{
//	    Name: "my-package",
//	    DistTags: map[string]string{"latest": "1.0.0"},
//	    Versions: map[string]models.Version{
//	        "1.0.0": {Name: "my-package", Version: "1.0.0"},
//	    },
//	}
//	err := registry.PublishPackage(ctx, pkg)
func (x *Registry) PublishPackage(ctx context.Context, pkg *models.Package) error {
	if err := x.requireToken(); err != nil {
		return err
	}

	// 获取当前 _rev（如果包已存在）
	rev, err := x.getRev(ctx, pkg.Name)
	if err != nil {
		return fmt.Errorf("failed to get package revision: %w", err)
	}
	if rev != "" {
		pkg.Rev = rev
	}

	targetUrl := fmt.Sprintf("%s/%s", x.options.RegistryURL, url.PathEscape(pkg.Name))
	_, err = x.sendJSON(ctx, http.MethodPut, targetUrl, pkg, http.StatusOK, http.StatusCreated)
	if err != nil {
		return fmt.Errorf("failed to publish package '%s': %w", pkg.Name, err)
	}
	return nil
}

// PublishPackageFromTarball 通过tarball发布NPM包
//
// 需要认证Token。这是一个更高级的发布方法，自动处理文档构造和 _rev 管理。
// 发布流程：
//  1. 获取当前包文档的 _rev（如果包已存在）
//  2. 构造完整的包文档，包含元数据和 tarball 的 dist 信息
//  3. PUT 完整包文档到 Registry
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - packageName: 包名称
//   - version: 要发布的版本号
//   - tarballBytes: tarball 文件的字节数组
//   - metadata: 包的元数据信息
//
// 返回值:
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	options := NewOptions().SetToken("npm_xxxxx")
//	registry := NewRegistry(options)
//	ctx := context.Background()
//	tarballBytes, _ := os.ReadFile("my-package-1.0.0.tgz")
//	metadata := &models.PublishMetadata{
//	    Name:        "my-package",
//	    Version:     "1.0.0",
//	    Description: "A test package",
//	}
//	err := registry.PublishPackageFromTarball(ctx, "my-package", "1.0.0", tarballBytes, metadata)
func (x *Registry) PublishPackageFromTarball(ctx context.Context, packageName, version string, tarballBytes []byte, metadata *models.PublishMetadata) error {
	if err := x.requireToken(); err != nil {
		return err
	}

	// 尝试获取当前包文档（如果包已存在），以便合并新版本而非覆盖
	var pkg *models.Package
	existingPkg, err := x.GetPackageInformation(ctx, packageName)
	if err == nil && existingPkg != nil {
		// 包已存在，在现有文档基础上添加新版本
		pkg = existingPkg
		// 更新 dist-tags
		if pkg.DistTags == nil {
			pkg.DistTags = make(map[string]string)
		}
		pkg.DistTags["latest"] = version
	} else {
		// 包不存在，创建新文档
		pkg = &models.Package{
			Name:        packageName,
			Description: metadata.Description,
			DistTags:    map[string]string{"latest": version},
			Versions:    make(map[string]models.Version),
		}
	}

	// 添加新版本到文档中
	if pkg.Versions == nil {
		pkg.Versions = make(map[string]models.Version)
	}
	pkg.Versions[version] = models.Version{
		Name:            packageName,
		Version:         version,
		Description:     metadata.Description,
		Main:            metadata.Main,
		Scripts:         models.Script(metadata.Scripts),
		Keywords:        metadata.Keywords,
		License:         metadata.License,
		Homepage:        metadata.Homepage,
		Dependencies:    metadata.Dependencies,
		DevDependencies: metadata.DevDependencies,
	}

	// 更新包级别元数据（如果提供了）
	if metadata.Description != "" {
		pkg.Description = metadata.Description
	}
	if len(metadata.Keywords) > 0 {
		pkg.Keywords = metadata.Keywords
	}
	if metadata.License != "" {
		pkg.License = metadata.License
	}
	if metadata.Homepage != "" {
		pkg.Homepage = metadata.Homepage
	}
	if metadata.Repository != nil {
		pkg.Repository = *metadata.Repository
	}
	if metadata.Author != nil && metadata.Author.Name != "" {
		pkg.Author = *metadata.Author
	}

	// 获取当前 _rev（如果包已存在）
	rev, err := x.getRev(ctx, packageName)
	if err != nil {
		return fmt.Errorf("failed to get package revision: %w", err)
	}
	if rev != "" {
		pkg.Rev = rev
	}

	targetUrl := fmt.Sprintf("%s/%s", x.options.RegistryURL, url.PathEscape(packageName))
	_, err = x.sendJSON(ctx, http.MethodPut, targetUrl, pkg, http.StatusOK, http.StatusCreated)
	if err != nil {
		return fmt.Errorf("failed to publish package '%s@%s': %w", packageName, version, err)
	}
	return nil
}
