package registry

import (
	"net/http"
	"net/url"
	"time"
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
	RegistryURL      string        // NPM 仓库服务器的 URL 地址
	Proxy            string        // HTTP 代理服务器的 URL，用于网络请求
	Token            string        // Bearer token for authenticated API requests
	DownloadStatsURL string        // 下载统计 API 的基础 URL，默认为 https://api.npmjs.org/downloads
	Timeout          time.Duration // 请求超时时间，默认为 0（不超时），由调用方通过 context 控制
	UserAgent        string        // HTTP User-Agent 头，默认为 "npm-skills-sdk"
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
		UserAgent:   "npm-skills-sdk",
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

// SetDownloadStatsURL 设置下载统计 API 的基础 URL
//
// 默认情况下，下载统计 API 使用 https://api.npmjs.org/downloads。
// 对于私有仓库或镜像，可以设置此字段以覆盖默认值。
//
// 参数:
//   - downloadStatsURL: 下载统计 API 的基础 URL
//
// 返回值:
//   - *Options: 返回自身以支持链式调用
func (o *Options) SetDownloadStatsURL(downloadStatsURL string) *Options {
	o.DownloadStatsURL = downloadStatsURL
	return o
}

// SetTimeout 设置请求超时时间
//
// 当设置了超时时间后，所有通过 Registry 客户端发出的请求都会自动应用此超时。
// 如果传入 0，则表示不设置超时（由调用方通过 context 自行控制）。
//
// 参数:
//   - timeout: 超时时间，例如 30*time.Second
//
// 返回值:
//   - *Options: 更新后的选项对象 (支持链式调用)
//
// 使用示例:
//
//	options := NewOptions().SetTimeout(30 * time.Second)
//	registry := NewRegistry(options)
func (o *Options) SetTimeout(timeout time.Duration) *Options {
	o.Timeout = timeout
	return o
}

// SetUserAgent 设置 HTTP User-Agent 头
//
// 默认为 "npm-skills-sdk"。某些 NPM 镜像或代理可能要求设置合理的 User-Agent。
//
// 参数:
//   - userAgent: User-Agent 字符串
//
// 返回值:
//   - *Options: 更新后的选项对象 (支持链式调用)
func (o *Options) SetUserAgent(userAgent string) *Options {
	o.UserAgent = userAgent
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
