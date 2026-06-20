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

// applyRequestDefaults 向请求选项应用默认配置（代理、认证、User-Agent）
//
// 这是所有请求方法（getBytes、sendRequest、sendJSON）共享的配置逻辑。
func (x *Registry) applyRequestDefaults(options *requests.Options[any, []byte]) {
	if x.options.Proxy != "" {
		options.AppendRequestSetting(requests.RequestSettingProxy(x.options.Proxy))
	}
	applyAuthSettings(x, options)
	if x.options.UserAgent != "" {
		options.AppendRequestSetting(requestSettingHeader("User-Agent", x.options.UserAgent))
	}
}

// applyTimeout 向上下文应用超时设置
//
// 如果 Options.Timeout > 0，返回带超时的新上下文和取消函数；
// 否则返回原上下文。
func (x *Registry) applyTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if x.options.Timeout > 0 {
		return context.WithTimeout(ctx, x.options.Timeout)
	}
	return ctx, func() {}
}

// defaultAcceptStatusCodes 返回默认接受的 HTTP 状态码
func defaultAcceptStatusCodes(codes []int) []int {
	if len(codes) == 0 {
		return []int{http.StatusOK, http.StatusCreated}
	}
	return codes
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
	ctx, cancel := x.applyTimeout(ctx)
	defer cancel()

	options := requests.NewOptions[any, []byte](targetUrl, requests.BytesResponseHandler(defaultAcceptStatusCodes(acceptStatusCodes)...))
	options = options.WithMethod(method)

	if body != nil {
		options = options.WithBody(body)
	}

	x.applyRequestDefaults(options)
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
	ctx, cancel := x.applyTimeout(ctx)
	defer cancel()

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON payload: %w", err)
	}

	options := requests.NewOptions[any, []byte](targetUrl, requests.BytesResponseHandler(defaultAcceptStatusCodes(acceptStatusCodes)...))
	options = options.WithMethod(method)
	options = options.WithBody(body)
	options.AppendRequestSetting(requestSettingHeader("Content-Type", "application/json"))

	x.applyRequestDefaults(options)
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

// requirePackageName 验证包名是否合法
//
// NPM 包名规则:
//   - 不能为空
//   - 长度不超过 214 字符
//   - 只能包含小写字母、数字、连字符(-)、下划线(_)、点(.)
//   - 不能以点(.)或下划线(_)开头
//   - scoped 包格式为 @scope/name，scope 和 name 都遵循上述规则
func requirePackageName(packageName string) error {
	if packageName == "" {
		return fmt.Errorf("package name is required and must not be empty")
	}
	if len(packageName) > 214 {
		return fmt.Errorf("package name too long: %d characters (max 214)", len(packageName))
	}

	// Handle scoped packages: @scope/name
	if strings.HasPrefix(packageName, "@") {
		parts := strings.SplitN(packageName, "/", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid scoped package name: %q (expected @scope/name)", packageName)
		}
		scope := strings.TrimPrefix(parts[0], "@")
		name := parts[1]
		if err := validateNamePart(scope, "scope"); err != nil {
			return err
		}
		if err := validateNamePart(name, "name"); err != nil {
			return err
		}
		return nil
	}

	return validateNamePart(packageName, "name")
}

// validateNamePart 验证包名的一部分（scope 或 name）
func validateNamePart(part, label string) error {
	if part == "" {
		return fmt.Errorf("package %s must not be empty", label)
	}
	if part[0] == '.' || part[0] == '_' {
		return fmt.Errorf("package %s %q must not start with '.' or '_'", label, part)
	}
	for _, c := range part {
		if !isValidNameChar(c) {
			return fmt.Errorf("package %s %q contains invalid character %q (allowed: lowercase letters, digits, '-', '_', '.')", label, part, c)
		}
	}
	return nil
}

// isValidNameChar 检查字符是否为合法的包名字符
func isValidNameChar(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-' || c == '_' || c == '.'
}
