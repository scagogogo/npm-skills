package registry

import (
	"net/http"
	"net/url"
)

// 默认 NPM 仓库地址
const DefaultRegistryURL = "https://registry.npmjs.org"

// Options 表示 Registry 客户端的配置选项
//
// 包含字段:
// - RegistryURL: NPM 仓库服务器的 URL 地址
// - Proxy: HTTP 代理服务器的 URL，用于网络请求
//
// 使用示例:
//
//	// 创建默认选项并自定义设置
//	options := NewOptions().
//		SetRegistryURL("https://registry.npmjs.org").
//		SetProxy("http://my-proxy.example.com:8080")
//
//	// 使用选项创建 Registry 客户端
//	registry := NewRegistry(options)
type Options struct {
	RegistryURL string
	Proxy       string
	Token       string // Bearer token for authenticated API requests
}

// NewOptions 创建并返回一个新的默认配置选项实例
//
// 默认配置:
// - RegistryURL: "https://registry.npmjs.org" (官方 NPM 仓库地址)
// - Proxy: 无代理设置
//
// 返回值:
//   - *Options: 配置有默认值的选项对象
//
// 使用示例:
//
//	options := NewOptions()
//	registry := NewRegistry(options)
func NewOptions() *Options {
	return &Options{
		RegistryURL: "https://registry.npmjs.org",
	}
}

// SetRegistryURL 设置 NPM 仓库服务器的 URL 地址
//
// 参数:
//   - url: 一个有效的 NPM 仓库 URL 地址字符串，例如:
//   - 官方仓库: "https://registry.npmjs.org"
//   - 淘宝镜像: "https://registry.npm.taobao.org"
//
// 返回值:
//   - *Options: 更新后的选项对象 (支持链式调用)
//
// 使用示例:
//
//	options := NewOptions().SetRegistryURL("https://registry.npm.taobao.org")
func (o *Options) SetRegistryURL(url string) *Options {
	o.RegistryURL = url
	return o
}

// SetProxy 设置 HTTP 代理服务器的 URL 地址
//
// 参数:
//   - proxyUrl: HTTP 代理服务器的 URL 地址字符串，例如:
//   - "http://proxy.example.com:8080"
//   - "http://username:password@proxy.example.com:8080"
//   - 传入空字符串可以清除之前设置的代理
//
// 返回值:
//   - *Options: 更新后的选项对象 (支持链式调用)
//
// 使用示例:
//
//	// 设置代理
//	options := NewOptions().SetProxy("http://proxy.corp.example.com:3128")
//
//	// 设置带认证的代理
//	options := NewOptions().SetProxy("http://user:pass@proxy.example.com:8080")
//
//	// 清除代理设置
//	options.SetProxy("")
func (o *Options) SetProxy(proxyUrl string) *Options {
	o.Proxy = proxyUrl
	return o
}

// SetToken 设置用于认证 API 请求的 Bearer token
//
// 参数:
//   - token: Bearer token 字符串，通常从 npm token create 或 npm login 获取
//
// 返回值:
//   - *Options: 更新后的选项对象 (支持链式调用)
//
// 使用示例:
//
//	options := NewOptions().SetToken("npm_xxxxx")
//	registry := NewRegistry(options)
func (o *Options) SetToken(token string) *Options {
	o.Token = token
	return o
}

// GetHttpClient 根据当前选项配置创建并返回一个 HTTP 客户端
//
// 如果设置了代理，返回的 HTTP 客户端将使用配置的代理服务器
// 如果没有设置代理，返回标准的 HTTP 客户端
//
// 返回值:
//   - *http.Client: 配置好的 HTTP 客户端
//   - error: 如果代理 URL 解析失败，返回错误
//
// 使用示例:
//
//	options := NewOptions().SetProxy("http://proxy.example.com:8080")
//	client, err := options.GetHttpClient()
//	if err != nil {
//		log.Fatalf("创建 HTTP 客户端失败: %v", err)
//	}
//	resp, err := client.Get("https://registry.npmjs.org/react")
func (o *Options) GetHttpClient() (*http.Client, error) {
	if o.Proxy == "" {
		return http.DefaultClient, nil
	}

	proxyUrl, err := url.Parse(o.Proxy)
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}

	return &http.Client{
		Transport: transport,
	}, nil
}
