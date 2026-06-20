package registry

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOptions(t *testing.T) {
	// 测试默认选项创建
	options := NewOptions()
	assert.NotNil(t, options)
	assert.Equal(t, DefaultRegistryURL, options.RegistryURL)
	assert.Empty(t, options.Proxy)
}

func TestSetRegistryURL(t *testing.T) {
	// 测试设置 Registry URL
	options := NewOptions()

	// 测试链式调用返回值
	result := options.SetRegistryURL("https://test-registry.example.com")
	assert.Equal(t, options, result, "应该返回自身以支持链式调用")

	// 测试实际设置的值
	assert.Equal(t, "https://test-registry.example.com", options.RegistryURL)

	// 测试设置为空字符串
	options.SetRegistryURL("")
	assert.Empty(t, options.RegistryURL)

	// 测试设置为非标准 URL
	options.SetRegistryURL("http://localhost:8080")
	assert.Equal(t, "http://localhost:8080", options.RegistryURL)
}

func TestSetProxy(t *testing.T) {
	// 测试设置代理
	options := NewOptions()

	// 测试链式调用返回值
	result := options.SetProxy("http://proxy.example.com:3128")
	assert.Equal(t, options, result, "应该返回自身以支持链式调用")

	// 测试实际设置的值
	assert.Equal(t, "http://proxy.example.com:3128", options.Proxy)

	// 测试设置为空字符串
	options.SetProxy("")
	assert.Empty(t, options.Proxy)

	// 测试设置为 socks5 代理
	options.SetProxy("socks5://127.0.0.1:1080")
	assert.Equal(t, "socks5://127.0.0.1:1080", options.Proxy)
}

func TestOptionsChaining(t *testing.T) {
	// 测试选项链式调用
	options := NewOptions().
		SetRegistryURL("https://custom-registry.org").
		SetProxy("http://proxy.example.org:8888")

	assert.NotNil(t, options)
	assert.Equal(t, "https://custom-registry.org", options.RegistryURL)
	assert.Equal(t, "http://proxy.example.org:8888", options.Proxy)

	// 测试链式调用中的顺序
	options = NewOptions().
		SetRegistryURL("https://first-registry.com").
		SetProxy("http://first-proxy.com").
		SetRegistryURL("https://second-registry.com").
		SetProxy("http://second-proxy.com")

	assert.Equal(t, "https://second-registry.com", options.RegistryURL)
	assert.Equal(t, "http://second-proxy.com", options.Proxy)
}

func TestOptionsUsage(t *testing.T) {
	// 测试在 Registry 中使用选项
	options := NewOptions().
		SetRegistryURL("https://test-usage.example.com").
		SetProxy("http://test-proxy.example.com")

	registry := NewRegistry(options)
	retrievedOptions := registry.GetOptions()

	assert.Equal(t, options, retrievedOptions)
	assert.Equal(t, "https://test-usage.example.com", retrievedOptions.RegistryURL)
	assert.Equal(t, "http://test-proxy.example.com", retrievedOptions.Proxy)
}

func TestGetHttpClient(t *testing.T) {
	// 测试无代理的情况
	options := NewOptions()
	client, err := options.GetHttpClient()
	assert.Nil(t, err)
	assert.NotNil(t, client)
	// 无代理时应该返回默认客户端
	assert.NotNil(t, client.Transport, "should create Transport for InsecureSkipVerify support")

	// 测试缓存：再次调用应该返回同一个客户端实例
	client2, err := options.GetHttpClient()
	assert.Nil(t, err)
	assert.Same(t, client, client2, "should return cached client instance")

	// 测试有效代理的情况 — 新 Options 实例，因为代理变更会 ResetHttpClient
	options = NewOptions()
	options.SetProxy("http://proxy.example.com:8080")
	client, err = options.GetHttpClient()
	assert.Nil(t, err)
	assert.NotNil(t, client)
	assert.NotEqual(t, http.DefaultClient, client, "使用代理时不应该返回默认客户端")

	// 测试无效代理URL的情况（使用包含无效字符的URL）
	options = NewOptions()
	options.SetProxy("http://proxy with spaces.com:8080")
	client, err = options.GetHttpClient()
	assert.NotNil(t, err, "包含空格的代理URL应该返回错误")
	assert.Nil(t, client, "无效代理URL时客户端应该为nil")

	// 错误后重试应该能重新初始化（sync.Once 被重置）
	options.SetProxy("http://proxy.example.com:8080")
	client, err = options.GetHttpClient()
	assert.Nil(t, err, "重置后应该能创建新客户端")
	assert.NotNil(t, client)

	// 测试代理URL格式错误的情况（使用无效的URL格式）
	options = NewOptions()
	options.SetProxy("://invalid-url")
	client, err = options.GetHttpClient()
	assert.NotNil(t, err, "格式错误的代理URL应该返回错误")
	assert.Nil(t, client, "格式错误的代理URL时客户端应该为nil")

	// 测试空字符串代理（应该等同于无代理）
	options = NewOptions()
	options.SetProxy("")
	client, err = options.GetHttpClient()
	assert.Nil(t, err)
	assert.NotNil(t, client)
	assert.NotNil(t, client, "空字符串代理应该返回可用客户端")
}

func TestResetHttpClient(t *testing.T) {
	options := NewOptions()

	// 获取初始客户端
	client1, err := options.GetHttpClient()
	assert.Nil(t, err)
	assert.NotNil(t, client1)

	// 未修改配置，应该返回同一实例
	client2, err := options.GetHttpClient()
	assert.Nil(t, err)
	assert.Same(t, client1, client2, "should return same cached instance")

	// 修改代理设置会自动触发 ResetHttpClient
	options.SetProxy("http://proxy.example.com:8080")
	client3, err := options.GetHttpClient()
	assert.Nil(t, err)
	assert.NotSame(t, client1, client3, "should return new instance after proxy change")

	// 手动重置
	options.ResetHttpClient()
	client4, err := options.GetHttpClient()
	assert.Nil(t, err)
	assert.NotSame(t, client3, client4, "should return new instance after manual reset")

	// 修改 InsecureSkipVerify 也会自动重置
	options2 := NewOptions()
	client5, _ := options2.GetHttpClient()
	options2.SetInsecureSkipVerify(true)
	client6, _ := options2.GetHttpClient()
	assert.NotSame(t, client5, client6, "should return new instance after InsecureSkipVerify change")
}

func TestSetToken(t *testing.T) {
	options := NewOptions()

	// 测试默认 token 为空
	assert.Empty(t, options.Token)

	// 测试设置 token
	result := options.SetToken("npm_test_token")
	assert.Equal(t, options, result, "应该返回自身以支持链式调用")
	assert.Equal(t, "npm_test_token", options.Token)

	// 测试清除 token
	options.SetToken("")
	assert.Empty(t, options.Token)

	// 测试链式调用
	options = NewOptions().
		SetRegistryURL("https://registry.npmjs.org").
		SetProxy("http://proxy:8080").
		SetToken("npm_chained_token")
	assert.Equal(t, "npm_chained_token", options.Token)
	assert.Equal(t, "http://proxy:8080", options.Proxy)
}

func TestRegistryWithToken(t *testing.T) {
	options := NewOptions().SetToken("npm_test_token")
	registry := NewRegistry(options)
	retrievedOpts := registry.GetOptions()
	assert.Equal(t, "npm_test_token", retrievedOpts.Token)
}

func TestOptionsStringMasking(t *testing.T) {
	// Test that String() masks sensitive fields
	options := NewOptions().
		SetRegistryURL("https://registry.npmjs.org").
		SetToken("npm_abcdef123456").
		SetBasicAuth("admin", "supersecret123")

	s := options.String()
	assert.Contains(t, s, "npm_****", "token should be masked showing first 4 chars")
	assert.NotContains(t, s, "npm_abcdef123456", "full token should not appear")
	assert.Contains(t, s, "supe****", "password should be masked showing first 4 chars")
	assert.NotContains(t, s, "supersecret123", "full password should not appear")
	assert.Contains(t, s, "admin", "username should not be masked")

	// Test short token (<=4 chars)
	options2 := NewOptions().SetToken("abc")
	s2 := options2.String()
	assert.Contains(t, s2, "****", "short token should be fully masked")
	assert.NotContains(t, s2, "abc", "even short token should not appear in full")

	// Test no auth
	options3 := NewOptions()
	s3 := options3.String()
	assert.Contains(t, s3, "Token:****", "empty token should show ****")
}
