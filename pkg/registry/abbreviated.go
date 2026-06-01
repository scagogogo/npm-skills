package registry

import (
	"context"
	"fmt"
	"net/http"

	"github.com/crawler-go-go-go/go-requests"
	"github.com/scagogogo/npm-crawler/pkg/models"
)

// requestSettingHeader 创建一个设置自定义 HTTP header 的 RequestSetting
func requestSettingHeader(key, value string) requests.RequestSetting {
	return func(client *http.Client, httpRequest *http.Request) error {
		httpRequest.Header.Set(key, value)
		return nil
	}
}

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
	targetUrl := fmt.Sprintf("%s/%s", x.options.RegistryURL, packageName)
	opts := requests.NewOptions[any, []byte](targetUrl, requests.BytesResponseHandler())
	opts.AppendRequestSetting(requestSettingHeader("Accept", "application/vnd.npm.install-v1+json"))
	if x.options.Proxy != "" {
		opts.AppendRequestSetting(requests.RequestSettingProxy(x.options.Proxy))
	}
	if x.options.Token != "" {
		opts.AppendRequestSetting(requestSettingHeader("Authorization", "Bearer "+x.options.Token))
	}
	bytes, err := requests.SendRequest[any, []byte](ctx, opts)
	if err != nil {
		return nil, err
	}
	return unmarshalJson[*models.Package](bytes)
}