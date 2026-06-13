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

	// 设置认证Token
	if x.options.Token != "" {
		options.AppendRequestSetting(requestSettingHeader("Authorization", "Bearer "+x.options.Token))
	}

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

	// 设置认证Token
	if x.options.Token != "" {
		options.AppendRequestSetting(requestSettingHeader("Authorization", "Bearer "+x.options.Token))
	}

	// 设置 User-Agent
	if x.options.UserAgent != "" {
		options.AppendRequestSetting(requestSettingHeader("User-Agent", x.options.UserAgent))
	}

	return requests.SendRequest[any, []byte](ctx, options)
}

// requireToken 验证Token是否已设置
//
// 所有写操作（发布、取消发布、设置dist-tag等）都需要认证Token。
// 此方法在写操作开头调用，确保Token已配置。
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
// 对于直接的包路径（如 /<package>），NPM Registry 可以处理原始的 scoped 包名，
// 但 /-/package/ 路径严格要求编码。
//
// 参数:
//   - name: 包名称，如 "react" 或 "@nestjs/core"
//
// 返回值:
//   - string: URL编码后的包名称
func encodePackageName(name string) string {
	// 对 scoped 包中的 "/" 进行 URL 编码
	// 保留 "@" 符号（NPM Registry 接受 @ 作为路径的一部分）
	if strings.Contains(name, "/") {
		return strings.Replace(name, "/", "%2F", 1)
	}
	return name
}

// encodePackageNameForPath 对包名进行完整路径编码
//
// 用于包名作为 URL 路径段的情况（如 /<package>/<version>）。
// NPM Registry 接受 @ 符号作为路径的一部分，但 "/" 需要编码。
//
// 参数:
//   - name: 包名称
//
// 返回值:
//   - string: URL编码后的包名称
func encodePackageNameForPath(name string) string {
	return url.PathEscape(name)
}

// requirePackageName 验证包名是否为空
//
// 所有接受包名参数的方法都应调用此函数进行验证。
//
// 参数:
//   - packageName: 包名称
//
// 返回值:
//   - error: 如果包名为空则返回错误
func requirePackageName(packageName string) error {
	if packageName == "" {
		return fmt.Errorf("package name is required and must not be empty")
	}
	return nil
}
