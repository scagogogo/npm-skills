package registry

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// NewCustomRegistry 创建一个用于自定义/私有 NPM 仓库的客户端
//
// 这是连接内网私有仓库最便捷的方式，支持所有常见的私有仓库方案：
//   - Verdaccio: 轻量级私有 npm 代理仓库
//   - Artifactory: JFrog 企业级制品仓库
//   - GitHub Packages: GitHub 的包仓库服务
//   - Nexus: Sonatype 仓库管理器
//   - AWS CodeArtifact: AWS 托管制品仓库
//
// 参数:
//   - registryURL: 私有仓库的 URL 地址
//   - opts: 可选的配置函数，用于设置 Token、Basic Auth、代理等
//
// 返回值:
//   - *Registry: 配置好的客户端实例
//
// 使用示例:
//
//	// Verdaccio with Basic Auth
//	client := NewCustomRegistry("http://verdaccio.internal:4873",
//	    func(o *Options) { o.SetBasicAuth("admin", "secret") })
//
//	// GitHub Packages with Token
//	client := NewCustomRegistry("https://npm.pkg.github.com",
//	    func(o *Options) { o.SetToken("ghp_xxxxx") })
//
//	// Artifactory with self-signed cert
//	client := NewCustomRegistry("https://artifactory.corp.com/artifactory/api/npm/npm-local",
//	    func(o *Options) {
//	        o.SetBasicAuth("user", "pass")
//	        o.SetInsecureSkipVerify(true)
//	        o.SetTimeout(30 * time.Second)
//	    })
func NewCustomRegistry(registryURL string, opts ...func(*Options)) *Registry {
	o := NewOptions().SetRegistryURL(registryURL)
	for _, opt := range opts {
		opt(o)
	}
	return NewRegistry(o)
}

// RegistryHealthCheck 检查仓库是否可用
//
// 发送一个轻量级请求验证仓库的连通性和认证状态。
// 适用于在长时间运行的服务中验证仓库可达性。
//
// 返回值:
//   - bool: 仓库是否可达
//   - error: 如果不可达，返回具体错误
//
// 使用示例:
//
//	client := NewCustomRegistry("http://verdaccio.internal:4873",
//	    func(o *Options) { o.SetBasicAuth("admin", "secret") })
//	ok, err := client.RegistryHealthCheck(ctx)
//	if !ok {
//	    log.Fatalf("Registry unreachable: %v", err)
//	}
func (x *Registry) RegistryHealthCheck(ctx context.Context) (bool, error) {
	// Use a short timeout for health checks
	if x.options.Timeout == 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	httpClient, err := x.options.GetHttpClient()
	if err != nil {
		return false, fmt.Errorf("failed to create HTTP client: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, x.options.RegistryURL, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	// Apply auth
	if x.options.Token != "" {
		req.Header.Set("Authorization", "Bearer "+x.options.Token)
	} else if x.options.Username != "" && x.options.Password != "" {
		req.SetBasicAuth(x.options.Username, x.options.Password)
	}
	if x.options.UserAgent != "" {
		req.Header.Set("User-Agent", x.options.UserAgent)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("registry unreachable at %s: %w", x.options.RegistryURL, err)
	}
	defer resp.Body.Close()

	switch {
	case resp.StatusCode == http.StatusOK:
		return true, nil
	case resp.StatusCode == http.StatusUnauthorized:
		return false, fmt.Errorf("registry requires authentication (HTTP 401): configure SetToken() or SetBasicAuth()")
	case resp.StatusCode == http.StatusForbidden:
		return false, fmt.Errorf("access denied (HTTP 403): check your credentials and permissions")
	case resp.StatusCode == http.StatusNotFound:
		return false, fmt.Errorf("registry not found (HTTP 404): check the registry URL %s", x.options.RegistryURL)
	default:
		return false, fmt.Errorf("registry returned unexpected status %d", resp.StatusCode)
	}
}

// IsPrivateRegistry 判断当前配置是否为私有/自定义仓库
//
// 基于 RegistryURL 判断：如果不是已知的公共镜像源，则认为是私有仓库。
// 可用于决定是否跳过某些仅公共仓库支持的功能（如下载统计）。
//
// 返回值:
//   - bool: 是否为私有仓库
func (x *Registry) IsPrivateRegistry() bool {
	publicRegistries := map[string]bool{
		"https://registry.npmjs.org":                     true,
		"https://registry.npm.taobao.org":                true,
		"https://registry.npmmirror.com":                 true,
		"https://mirrors.huaweicloud.com/repository/npm": true,
		"http://mirrors.cloud.tencent.com/npm":           true,
		"http://r.cnpmjs.org":                            true,
		"https://registry.yarnpkg.com":                   true,
		"https://skimdb.npmjs.com":                       true,
	}
	return !publicRegistries[x.options.RegistryURL]
}
