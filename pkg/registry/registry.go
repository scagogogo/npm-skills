package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/crawler-go-go-go/go-requests"
	"github.com/scagogogo/npm-crawler/pkg/models"
)

// Registry NPM 注册表访问客户端，提供与 NPM Registry 交互的方法
// 可以使用不同的镜像源配置来创建实例，支持代理设置
type Registry struct {
	options *Options
}

// NewRegistry 创建一个新的 Registry 客户端实例
//
// 参数:
//   - options: 可选的配置选项，如未提供则使用默认配置（https://registry.npmjs.org）
//
// 返回值:
//   - *Registry: 新创建的 Registry 客户端实例
//
// 使用示例:
//
//	// 使用默认配置创建客户端
//	registry := NewRegistry()
//
//	// 使用自定义配置创建客户端
//	options := NewOptions().SetRegistryURL("https://registry.npmjs.org").SetProxy("http://proxy.example.com:8080")
//	registry := NewRegistry(options)
func NewRegistry(options ...*Options) *Registry {
	if len(options) == 0 {
		options = append(options, NewOptions())
	}
	return &Registry{
		options: options[0],
	}
}

// GetRegistryInformation 获取 NPM Registry 的状态信息
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//
// 返回值:
//   - *models.RegistryInformation: Registry 状态信息，包含数据库名称、文档数量、存储大小等
//   - error: 如果请求失败则返回错误
//
// 数据样例:
//
//	{
//	  "db_name": "registry",
//	  "doc_count": 1000000,
//	  "doc_del_count": 10000,
//	  "update_seq": 5000000,
//	  "purge_seq": 0,
//	  "compact_running": false,
//	  "sizes": {
//	    "file": 100000000,
//	    "active": 90000000,
//	    "external": 50000000
//	  },
//	  "disk_size": 100000000,
//	  "data_size": 90000000,
//	  "instance_start_time": "1641028800000"
//	}
//
// 使用示例:
//
//	registry := NewRegistry()
//	ctx := context.Background()
//	info, err := registry.GetRegistryInformation(ctx)
//	if err != nil {
//	  // 处理错误
//	}
//	fmt.Println("Registry 文档数量:", info.DocCount)
func (x *Registry) GetRegistryInformation(ctx context.Context) (*models.RegistryInformation, error) {
	targetUrl := x.options.RegistryURL
	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, err
	}
	return unmarshalJson[*models.RegistryInformation](bytes)
}

// GetPackageInformation 获取指定 NPM 包的详细信息
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - packageName: 要查询的包名称，例如 "react"、"lodash" 等
//
// 返回值:
//   - *models.Package: 包的详细信息，包含版本信息、维护者、依赖关系等
//   - error: 如果请求失败则返回错误
//
// 数据样例:
//
//	{
//	  "_id": "axios",
//	  "name": "axios",
//	  "description": "Promise based HTTP client",
//	  "dist-tags": {
//	    "latest": "1.0.0"
//	  },
//	  "versions": {
//	    "1.0.0": {
//	      "name": "axios",
//	      "version": "1.0.0",
//	      "description": "Promise based HTTP client",
//	      "dependencies": {
//	        "follow-redirects": "^1.15.0"
//	      },
//	      "dist": {
//	        "shasum": "abc123",
//	        "tarball": "https://registry.npmjs.org/axios/-/axios-1.0.0.tgz"
//	      }
//	    }
//	  }
//	}
//
// 使用示例:
//
//	registry := NewRegistry()
//	ctx := context.Background()
//	pkg, err := registry.GetPackageInformation(ctx, "axios")
//	if err != nil {
//	  // 处理错误
//	}
//	fmt.Println("包名:", pkg.Name)
//	fmt.Println("最新版本:", pkg.DistTags.Latest)
func (x *Registry) GetPackageInformation(ctx context.Context, packageName string) (*models.Package, error) {
	targetUrl := fmt.Sprintf("%s/%s", x.options.RegistryURL, packageName)
	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, err
	}
	return unmarshalJson[*models.Package](bytes)
}

// SearchPackages 搜索 NPM 包
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - query: 搜索关键字
//   - limit: 返回结果数量限制，默认为 20
//
// 返回值:
//   - *models.SearchResult: 搜索结果，包含匹配的包列表
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	registry := NewRegistry()
//	ctx := context.Background()
//	result, err := registry.SearchPackages(ctx, "react", 10)
//	if err != nil {
//		// 处理错误
//	}
//	for _, pkg := range result.Objects {
//		fmt.Println("包名:", pkg.Package.Name)
//	}
func (x *Registry) SearchPackages(ctx context.Context, query string, limit int) (*models.SearchResult, error) {
	if limit <= 0 {
		limit = 20
	}
	targetUrl := fmt.Sprintf("%s/-/v1/search?text=%s&size=%d", x.options.RegistryURL, url.QueryEscape(query), limit)
	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, err
	}
	return unmarshalJson[*models.SearchResult](bytes)
}

// GetPackageVersion 获取指定 NPM 包的特定版本信息
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - packageName: 要查询的包名称，例如 "react"、"lodash" 等
//   - version: 要查询的版本号，例如 "1.0.0"、"latest" 等
//
// 返回值:
//   - *models.Version: 指定版本的详细信息
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	registry := NewRegistry()
//	ctx := context.Background()
//	version, err := registry.GetPackageVersion(ctx, "react", "18.0.0")
//	if err != nil {
//		// 处理错误
//	}
//	fmt.Println("版本:", version.Version)
//	fmt.Println("依赖:", version.Dependencies)
func (x *Registry) GetPackageVersion(ctx context.Context, packageName, version string) (*models.Version, error) {
	targetUrl := fmt.Sprintf("%s/%s/%s", x.options.RegistryURL, packageName, version)
	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, err
	}
	return unmarshalJson[*models.Version](bytes)
}

// GetDownloadStats 获取指定 NPM 包的下载统计信息
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - packageName: 要查询的包名称
//   - period: 统计周期，例如 "last-day", "last-week", "last-month"
//
// 返回值:
//   - *models.DownloadStats: 下载统计信息
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	registry := NewRegistry()
//	ctx := context.Background()
//	stats, err := registry.GetDownloadStats(ctx, "react", "last-week")
//	if err != nil {
//		// 处理错误
//	}
//	fmt.Println("下载次数:", stats.Downloads)
func (x *Registry) GetDownloadStats(ctx context.Context, packageName, period string) (*models.DownloadStats, error) {
	baseURL := "https://api.npmjs.org/downloads"
	targetUrl := fmt.Sprintf("%s/point/%s/%s", baseURL, period, packageName)
	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, err
	}
	return unmarshalJson[*models.DownloadStats](bytes)
}

// DownloadTarball 下载指定 NPM 包的 tarball 文件到本地路径
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - packageName: 要下载的包名称，例如 "react"、"lodash" 等
//   - version: 要下载的版本号，例如 "18.0.0"、"latest" 等
//   - destPath: 目标文件保存路径，例如 "./downloads/react-18.0.0.tgz"
//
// 返回值:
//   - error: 如果下载失败则返回错误
//
// 使用示例:
//
//	registry := NewRegistry()
//	ctx := context.Background()
//	err := registry.DownloadTarball(ctx, "react", "18.0.0", "./react.tgz")
//	if err != nil {
//	  // 处理错误
//	}
func (x *Registry) DownloadTarball(ctx context.Context, packageName, version, destPath string) error {
	// 先获取包的版本信息，以获取 tarball URL
	versionInfo, err := x.GetPackageVersion(ctx, packageName, version)
	if err != nil {
		return err
	}

	// 从版本信息中获取 tarball URL
	tarballURL := versionInfo.Dist.Tarball
	if tarballURL == "" {
		return fmt.Errorf("tarball URL not found for package %s@%s", packageName, version)
	}

	// 使用 getBytes 获取 tarball 内容
	bytes, err := x.getBytes(ctx, tarballURL)
	if err != nil {
		return err
	}

	// 写入本地文件
	return os.WriteFile(destPath, bytes, 0644)
}

// GetOptions 获取当前 Registry 客户端的配置选项
//
// 返回值:
//   - *Options: 当前配置选项，包含 Registry URL 和代理设置等
//
// 使用示例:
//
//	registry := NewRegistry()
//	options := registry.GetOptions()
//	fmt.Println("当前 Registry URL:", options.RegistryURL)
func (x *Registry) GetOptions() *Options {
	return x.options
}

// unmarshalJson 将 JSON 字节数组解析为指定类型的对象
//
// 参数:
//   - bytes: 包含 JSON 数据的字节数组
//
// 返回值:
//   - T: 解析后的指定类型对象
//   - error: 如果解析失败则返回错误
//
// 注意: 这是一个泛型函数，T 可以是任何可解析的 JSON 目标类型
func unmarshalJson[T any](bytes []byte) (T, error) {
	var r T
	err := json.Unmarshal(bytes, &r)
	if err != nil {
		var zero T
		return zero, err
	}
	return r, nil
}

// getBytes 从指定 URL 获取响应数据的字节数组
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - targetUrl: 请求的目标 URL
//
// 返回值:
//   - []byte: 响应数据的字节数组
//   - error: 如果请求失败则返回错误
//
// 注意: 这是一个内部方法，支持代理设置和 Token 认证
func (x *Registry) getBytes(ctx context.Context, targetUrl string) ([]byte, error) {
	options := requests.NewOptions[any, []byte](targetUrl, requests.BytesResponseHandler())
	if x.options.Proxy != "" {
		options.AppendRequestSetting(requests.RequestSettingProxy(x.options.Proxy))
	}
	if x.options.Token != "" {
		options.AppendRequestSetting(requestSettingBearerToken(x.options.Token))
	}
	return requests.SendRequest[any, []byte](ctx, options)
}

// requestSettingBearerToken 创建一个设置 Bearer Token 认证头的 RequestSetting
func requestSettingBearerToken(token string) requests.RequestSetting {
	return func(client *http.Client, httpRequest *http.Request) error {
		httpRequest.Header.Set("Authorization", "Bearer "+token)
		return nil
	}
}
