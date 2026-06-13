package registry

import (
	"context"
	"fmt"
	"sort"

	"github.com/scagogogo/npm-skills/pkg/models"
)

// GetPackageVersions 获取指定 NPM 包的所有已发布版本号列表
//
// 使用精简 API (abbreviated) 获取包元数据，仅提取版本号列表，
// 比获取完整包信息更高效（完整包信息可达 10MB+）。
// 返回的版本号列表按 semver 排序。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - packageName: 要查询的包名称
//
// 返回值:
//   - []string: 所有已发布的版本号列表（按字典序排序）
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	registry := NewRegistry()
//	ctx := context.Background()
//	versions, err := registry.GetPackageVersions(ctx, "react")
//	if err != nil {
//	    // 处理错误
//	}
//	for _, v := range versions {
//	    fmt.Println(v)
//	}
func (x *Registry) GetPackageVersions(ctx context.Context, packageName string) ([]string, error) {
	// 使用精简 API 减少数据传输量
	pkg, err := x.GetAbbreviatedPackageInformation(ctx, packageName)
	if err != nil {
		return nil, fmt.Errorf("failed to get versions for '%s': %w", packageName, err)
	}

	versions := make([]string, 0, len(pkg.Versions))
	for v := range pkg.Versions {
		versions = append(versions, v)
	}
	sort.Strings(versions)

	return versions, nil
}

// GetPackageVersionCount 获取指定 NPM 包的已发布版本数量
//
// 比获取完整版本列表更轻量，适用于只需要数量信息的场景。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - packageName: 要查询的包名称
//
// 返回值:
//   - int: 已发布的版本数量
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	registry := NewRegistry()
//	ctx := context.Background()
//	count, err := registry.GetPackageVersionCount(ctx, "react")
//	if err != nil {
//	    // 处理错误
//	}
//	fmt.Printf("react has %d published versions\n", count)
func (x *Registry) GetPackageVersionCount(ctx context.Context, packageName string) (int, error) {
	pkg, err := x.GetAbbreviatedPackageInformation(ctx, packageName)
	if err != nil {
		return 0, fmt.Errorf("failed to get version count for '%s': %w", packageName, err)
	}
	return len(pkg.Versions), nil
}

// GetPackageLatestVersion 获取指定 NPM 包的最新版本号
//
// 通过 dist-tags 获取 latest 标签对应的版本号，不需要获取完整包信息。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - packageName: 要查询的包名称
//
// 返回值:
//   - string: 最新版本号
//   - error: 如果请求失败或包没有 latest 标签则返回错误
//
// 使用示例:
//
//	registry := NewRegistry()
//	ctx := context.Background()
//	latest, err := registry.GetPackageLatestVersion(ctx, "react")
//	if err != nil {
//	    // 处理错误
//	}
//	fmt.Println("Latest:", latest)
func (x *Registry) GetPackageLatestVersion(ctx context.Context, packageName string) (string, error) {
	tags, err := x.GetDistTagsAbbreviated(ctx, packageName)
	if err != nil {
		return "", fmt.Errorf("failed to get latest version for '%s': %w", packageName, err)
	}
	latest, ok := tags["latest"]
	if !ok {
		return "", fmt.Errorf("package '%s' has no 'latest' dist-tag", packageName)
	}
	return latest, nil
}

// GetPackageInformationSummary 获取包的摘要信息（名称、描述、最新版本、版本数量、dist-tags）
//
// 使用精简 API 获取数据，适用于列表展示或概览场景。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - packageName: 要查询的包名称
//
// 返回值:
//   - *models.Package: 精简的包信息
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	registry := NewRegistry()
//	ctx := context.Background()
//	pkg, err := registry.GetPackageInformationSummary(ctx, "react")
//	if err != nil {
//	    // 处理错误
//	}
//	fmt.Printf("%s@%s (%d versions)\n", pkg.Name, pkg.DistTags["latest"], len(pkg.Versions))
func (x *Registry) GetPackageInformationSummary(ctx context.Context, packageName string) (*models.Package, error) {
	return x.GetAbbreviatedPackageInformation(ctx, packageName)
}
