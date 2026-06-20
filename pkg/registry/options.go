package registry

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// DefaultRegistryURL 是默认的 NPM 官方仓库地址
const DefaultRegistryURL = "https://registry.npmjs.org"

// Options 表示 NPM 仓库客户端的配置选项
//
// 支持自定义仓库地址、代理、认证、超时等配置。
// 所有字段均可通过链式 Setter 方法设置。
//
// HTTP 客户端通过 GetHttpClient() 获取时会自动缓存并复用，
// 以便利用连接池和 Keep-Alive 提升性能。
//
// 示例:
//
//	options := NewOptions().
//	    SetRegistryURL("https://npm.my-company.com").
//	    SetToken("npm_xxxxx").
//	    SetTimeout(30 * time.Second)
type Options struct {
	RegistryURL        string        // NPM 仓库服务器的 URL 地址
	Proxy              string        // HTTP 代理服务器的 URL，用于网络请求
	Token              string        // Bearer token for authenticated API requests
	Username           string        // 用户名，用于 Basic Auth 认证（私有仓库常用）
	Password           string        // 密码，用于 Basic Auth 认证（私有仓库常用）
	DownloadStatsURL   string        // 下载统计 API 的基础 URL，默认为 https://api.npmjs.org/downloads
	Timeout            time.Duration // 请求超时时间，默认为 0（不超时），由调用方通过 context 控制
	UserAgent          string        // HTTP User-Agent 头，默认为 "npm-skills-sdk"
	InsecureSkipVerify bool          // 是否跳过 TLS 证书验证（内网自签名证书场景）

	// 内部缓存字段，懒初始化
	httpClient     *http.Client
	httpClientOnce sync.Once
}

// NewOptions 创建默认的 Options 实例
//
// 默认配置：
//   - RegistryURL: https://registry.npmjs.org（NPM 官方仓库）
//   - UserAgent: npm-skills-sdk
//
// 返回值:
//   - *Options: 新创建的选项实例
func NewOptions() *Options {
	return &Options{
		RegistryURL: "https://registry.npmjs.org",
		UserAgent:   "npm-skills-sdk",
	}
}

// SetRegistryURL 设置 NPM 仓库服务器的 URL 地址
//
// 适用于自定义私有仓库（如 Verdaccio、Artifactory、GitHub Packages、Nexus 等）。
// 设置后将所有 API 请求发送到该地址。
//
// 参数:
//   - registryURL: 仓库 URL，如 "https://npm.my-company.com"
//
// 返回值:
//   - *Options: 更新后的选项对象 (支持链式调用)
//
// 使用示例:
//
//	// 连接到公司内部 Verdaccio 仓库
//	options := NewOptions().SetRegistryURL("http://verdaccio.internal:4873")
//	// 连接到 GitHub Packages
//	options := NewOptions().SetRegistryURL("https://npm.pkg.github.com")
func (o *Options) SetRegistryURL(registryURL string) *Options {
	o.RegistryURL = registryURL
	return o
}

// SetProxy 设置 HTTP 代理服务器
//
// 参数:
//   - proxy: 代理服务器 URL，如 "http://127.0.0.1:7890"
//
// 返回值:
//   - *Options: 更新后的选项对象 (支持链式调用)
func (o *Options) SetProxy(proxy string) *Options {
	o.Proxy = proxy
	o.ResetHttpClient()
	return o
}

// SetToken 设置 Bearer Token 认证
//
// 大多数私有仓库和所有写操作需要认证。Token 可以从
// npm token create 或仓库管理界面获取。
//
// 参数:
//   - token: Bearer token 字符串，如 "npm_xxxxx"
//
// 返回值:
//   - *Options: 更新后的选项对象 (支持链式调用)
//
// 使用示例:
//
//	options := NewOptions().
//	    SetRegistryURL("https://npm.pkg.github.com").
//	    SetToken("ghp_xxxxx")
func (o *Options) SetToken(token string) *Options {
	o.Token = token
	return o
}

// SetBasicAuth 设置 Basic Auth 认证
//
// 部分私有仓库（如 Verdaccio、Artifactory）使用用户名密码认证。
// 设置后，所有请求将自动携带 Authorization: Basic <encoded> 头。
// 如果同时设置了 Token，Token 优先级更高。
//
// 参数:
//   - username: 用户名
//   - password: 密码
//
// 返回值:
//   - *Options: 更新后的选项对象 (支持链式调用)
//
// 使用示例:
//
//	// 连接到公司内部 Verdaccio 仓库
//	options := NewOptions().
//	    SetRegistryURL("http://verdaccio.internal:4873").
//	    SetBasicAuth("admin", "secret123")
func (o *Options) SetBasicAuth(username, password string) *Options {
	o.Username = username
	o.Password = password
	return o
}

// SetDownloadStatsURL 设置下载统计 API 的基础 URL
//
// 大多数私有仓库不提供下载统计 API，可以将其设置为空字符串
// 以禁用下载统计功能，避免请求失败。
//
// 参数:
//   - downloadStatsURL: 下载统计 API URL，如 "https://api.npmjs.org/downloads"
//
// 返回值:
//   - *Options: 更新后的选项对象 (支持链式调用)
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

// SetInsecureSkipVerify 设置是否跳过 TLS 证书验证
//
// 内网私有仓库经常使用自签名证书，默认 Go HTTP 客户端会拒绝此类证书。
// 设置为 true 可以跳过证书验证，但请注意这会降低安全性。
//
// ⚠️ 仅建议在受控的内网环境中使用此选项，生产环境请配置正确的 TLS 证书。
//
// 参数:
//   - skip: 是否跳过 TLS 证书验证
//
// 返回值:
//   - *Options: 更新后的选项对象 (支持链式调用)
func (o *Options) SetInsecureSkipVerify(skip bool) *Options {
	o.InsecureSkipVerify = skip
	o.ResetHttpClient()
	return o
}

// HasAuth 返回是否配置了认证信息（Token 或 Basic Auth）
func (o *Options) HasAuth() bool {
	return o.Token != "" || (o.Username != "" && o.Password != "")
}

// GetHttpClient 获取配置了代理和 TLS 的 HTTP 客户端
//
// 返回的客户端会被缓存并复用，以利用连接池和 Keep-Alive 提升性能。
// 如果修改了 Proxy 或 InsecureSkipVerify 等影响传输层的配置，
// 需要调用 ResetHttpClient() 使缓存失效。
func (o *Options) GetHttpClient() (*http.Client, error) {
	var initErr error
	o.httpClientOnce.Do(func() {
		transport := &http.Transport{}

		// 配置代理
		if o.Proxy != "" {
			proxyURL, err := url.Parse(o.Proxy)
			if err != nil {
				initErr = fmt.Errorf("invalid proxy URL: %w", err)
				return
			}
			transport.Proxy = http.ProxyURL(proxyURL)
		}

		// 配置 TLS
		if o.InsecureSkipVerify {
			transport.TLSClientConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}

		o.httpClient = &http.Client{
			Transport: transport,
		}
	})
	if initErr != nil {
		// 重置 sync.Once 以便下次重试
		o.httpClientOnce = sync.Once{}
		return nil, initErr
	}
	return o.httpClient, nil
}

// ResetHttpClient 重置缓存的 HTTP 客户端
//
// 当修改了 Proxy、InsecureSkipVerify 等影响传输层的配置后，
// 需要调用此方法使缓存失效，下次 GetHttpClient() 将创建新的客户端。
func (o *Options) ResetHttpClient() {
	o.httpClientOnce = sync.Once{}
	o.httpClient = nil
}
