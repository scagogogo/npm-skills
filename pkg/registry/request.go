package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/crawler-go-go-go/go-requests"
)

// requestSettingHeader 创建一个设置自定义 HTTP header 的 RequestSetting
func requestSettingHeader(key, value string) requests.RequestSetting {
	return func(client *http.Client, httpRequest *http.Request) error {
		httpRequest.Header.Set(key, value)
		return nil
	}
}

// requestSettingBasicAuth 创建一个设置 HTTP Basic Auth 的 RequestSetting
//
// 适用于使用用户名密码认证的私有仓库（如 Verdaccio、Artifactory）。
func requestSettingBasicAuth(username, password string) requests.RequestSetting {
	return func(client *http.Client, httpRequest *http.Request) error {
		httpRequest.SetBasicAuth(username, password)
		return nil
	}
}

// applyAuthSettings 向请求选项添加认证设置
//
// 认证优先级：Token > Basic Auth
// Token 通常从 npm token create 获取，Basic Auth 使用用户名密码。
func applyAuthSettings(x *Registry, options *requests.Options[any, []byte]) {
	if x.options.Token != "" {
		options.AppendRequestSetting(requestSettingHeader("Authorization", "Bearer "+x.options.Token))
	} else if x.options.Username != "" && x.options.Password != "" {
		options.AppendRequestSetting(requestSettingBasicAuth(x.options.Username, x.options.Password))
	}
}

// sendRequest 发送HTTP请求到指定URL，支持自定义方法和请求体
//
// 这是所有写操作（PUT/POST/DELETE）的底层传输方法。
// 对于只读GET请求，继续使用 getBytes() 方法。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - method: HTTP方法，如 "PUT", "POST", "DELETE"
//   - targetUrl: 请求的目标URL
//   - body: 请求体字节数组，可为nil
//   - acceptStatusCodes: 接受的HTTP状态码列表，默认为200和201
//
// 返回值:
//   - []byte: 响应体字节数组
//   - error: 如果请求失败则返回错误
func (x *Registry) sendRequest(ctx context.Context, method, targetUrl string, body []byte, acceptStatusCodes ...int) ([]byte, error) {
	// Apply timeout from options if set
	if x.options.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, x.options.Timeout)
		defer cancel()
	}

	// 默认接受 200 OK 和 201 Created
	if len(acceptStatusCodes) == 0 {
		acceptStatusCodes = []int{http.StatusOK, http.StatusCreated}
	}

	options := requests.NewOptions[any, []byte](targetUrl, requests.BytesResponseHandler(acceptStatusCodes...))
	options = options.WithMethod(method)

	// 设置请求体
	if body != nil {
		options = options.WithBody(body)
	}

	// 设置代理
	if x.options.Proxy != "" {
		options.AppendRequestSetting(requests.RequestSettingProxy(x.options.Proxy))
	}

	// 设置认证（Token 优先，Basic Auth 其次）
	applyAuthSettings(x, options)

	// 设置 User-Agent
	if x.options.UserAgent != "" {
		options.AppendRequestSetting(requestSettingHeader("User-Agent", x.options.UserAgent))
	}

	return requests.SendRequest[any, []byte](ctx, options)
}

// sendJSON 发送JSON格式的HTTP请求
//
// 将payload序列化为JSON后作为请求体发送，并设置Content-Type为application/json。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - method: HTTP方法，如 "PUT", "POST"
//   - targetUrl: 请求的目标URL
//   - payload: 要序列化为JSON的请求体对象
//   - acceptStatusCodes: 接受的HTTP状态码列表，默认为200和201
//
// 返回值:
//   - []byte: 响应体字节数组
//   - error: 如果请求失败则返回错误
func (x *Registry) sendJSON(ctx context.Context, method, targetUrl string, payload interface{}, acceptStatusCodes ...int) ([]byte, error) {
	// Apply timeout from options if set
	if x.options.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, x.options.Timeout)
		defer cancel()
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON payload: %w", err)
	}

	// 默认接受 200 OK 和 201 Created
	if len(acceptStatusCodes) == 0 {
		acceptStatusCodes = []int{http.StatusOK, http.StatusCreated}
	}

	options := requests.NewOptions[any, []byte](targetUrl, requests.BytesResponseHandler(acceptStatusCodes...))
	options = options.WithMethod(method)
	options = options.WithBody(body)

	// 设置Content-Type
	options.AppendRequestSetting(requestSettingHeader("Content-Type", "application/json"))

	// 设置代理
	if x.options.Proxy != "" {
		options.AppendRequestSetting(requests.RequestSettingProxy(x.options.Proxy))
	}

	// 设置认证（Token 优先，Basic Auth 其次）
	applyAuthSettings(x, options)

	// 设置 User-Agent
	if x.options.UserAgent != "" {
		options.AppendRequestSetting(requestSettingHeader("User-Agent", x.options.UserAgent))
	}

	return requests.SendRequest[any, []byte](ctx, options)
}

// requireAuth 验证是否已配置认证信息（Token 或 Basic Auth）
//
// 所有写操作都需要认证。Token 优先级高于 Basic Auth。
//
// 返回值:
//   - error: 如果未配置认证信息则返回错误
func (x *Registry) requireAuth() error {
	if !x.options.HasAuth() {
		return fmt.Errorf("authentication required: configure with options.SetToken() or options.SetBasicAuth() before calling write operations")
	}
	return nil
}

// requireToken 验证Token是否已设置
//
// 仅检查 Token，不检查 Basic Auth。适用于某些严格要求 Token 的操作。
//
// 返回值:
//   - error: 如果Token未设置则返回错误
func (x *Registry) requireToken() error {
	if x.options.Token == "" {
		return fmt.Errorf("authentication required: no token set. Configure with options.SetToken() before calling write operations")
	}
	return nil
}

// encodePackageName 对包名称进行URL编码，处理 scoped 包（如 @nestjs/core）
//
// NPM Registry 的 /-/package/ 端点要求对 scoped 包中的 "/" 进行编码。
// 例如 @nestjs/core 需要编码为 @nestjs%2Fcore。
func encodePackageName(name string) string {
	if strings.Contains(name, "/") {
		return strings.Replace(name, "/", "%2F", 1)
	}
	return name
}

// encodePackageNameForPath 对包名进行完整路径编码
func encodePackageNameForPath(name string) string {
	return url.PathEscape(name)
}

// requirePackageName 验证包名是否为空
func requirePackageName(packageName string) error {
	if packageName == "" {
		return fmt.Errorf("package name is required and must not be empty")
	}
	return nil
}
