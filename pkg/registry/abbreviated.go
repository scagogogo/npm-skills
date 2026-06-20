package registry

import (
	"context"
	"fmt"
	"net/url"

	"github.com/scagogogo/npm-skills/pkg/models"
)

// GetAbbreviatedPackageInformation 获取指定包的精简元数据
//
// 使用 NPM 的 install-v1 Accept header，返回的元数据比完整包信息小得多
// （完整信息可达 10MB+，精简版通常几 KB），适合只需要版本列表和 dist-tags 的场景。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - packageName: 要查询的包名称
//
// 返回值:
//   - *models.Package: 精简的包信息（可能缺少 README、maintainers 等字段）
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	registry := NewRegistry()
//	ctx := context.Background()
//	pkg, err := registry.GetAbbreviatedPackageInformation(ctx, "react")
//	if err != nil {
//	    // 处理错误
//	}
//	fmt.Println("Latest:", pkg.DistTags["latest"])
func (x *Registry) GetAbbreviatedPackageInformation(ctx context.Context, packageName string) (*models.Package, error) {
	if err := requirePackageName(packageName); err != nil {
		return nil, err
	}
	targetUrl := fmt.Sprintf("%s/%s", x.options.RegistryURL, url.PathEscape(packageName))
	bytes, err := x.getBytesWithHeaders(ctx, targetUrl, map[string]string{
		"Accept": "application/vnd.npm.install-v1+json",
	})
	if err != nil {
		return nil, err
	}
	return unmarshalJson[*models.Package](bytes)
}
