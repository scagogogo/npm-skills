package registry

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockTestServer creates a mock NPM registry server for testing
func mockTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// axios 版本路由 - GET /axios/1.0.0
		if r.URL.Path == "/axios/1.0.0" {
			// tarball URL 指向 mock 服务器自身（通过 server.URL 动态设置）
			serverURL := "http://" + r.Host
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"name": "axios",
				"version": "1.0.0",
				"description": "Promise based HTTP client",
				"main": "index.js",
				"dependencies": {"follow-redirects": "^1.15.0"},
				"dist": {
					"shasum": "abc123",
					"tarball": "` + serverURL + `/axios/-/axios-1.0.0.tgz"
				}
			}`))
			return
		}

		// axios latest 版本路由 - GET /axios/latest
		if r.URL.Path == "/axios/latest" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"name": "axios",
				"version": "1.5.0",
				"description": "Promise based HTTP client",
				"main": "dist/axios.js",
				"dependencies": {"follow-redirects": "^1.15.0"},
				"dist": {
					"shasum": "def456",
					"tarball": "https://registry.npmjs.org/axios/-/axios-1.5.0.tgz"
				}
			}`))
			return
		}

		// axios tarball 路由 - GET /axios/-/axios-1.0.0.tgz
		if r.URL.Path == "/axios/-/axios-1.0.0.tgz" {
			w.Header().Set("Content-Type", "application/octet-stream")
			w.WriteHeader(http.StatusOK)
			// 写入简化的 gzip 头部（实际 npm tarball 是 gzip 压缩的）
			w.Write([]byte{0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00})
			return
		}

		// 默认返回空对象
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
}

func TestGetPackageVersionMock(t *testing.T) {
	server := mockTestServer()
	defer server.Close()

	registry := NewRegistry(NewOptions().SetRegistryURL(server.URL))

	// 测试获取有效版本
	version, err := registry.GetPackageVersion(context.Background(), "axios", "1.0.0")
	assert.Nil(t, err)
	assert.NotNil(t, version)
	assert.Equal(t, "axios", version.Name)
	assert.Equal(t, "1.0.0", version.Version)
	assert.Equal(t, "Promise based HTTP client", version.Description)
	assert.NotNil(t, version.Dependencies)
	assert.Contains(t, version.Dependencies, "follow-redirects")
}

func TestGetPackageVersionLatestMock(t *testing.T) {
	server := mockTestServer()
	defer server.Close()

	registry := NewRegistry(NewOptions().SetRegistryURL(server.URL))

	// 测试获取 latest 版本
	version, err := registry.GetPackageVersion(context.Background(), "axios", "latest")
	assert.Nil(t, err)
	assert.NotNil(t, version)
	assert.Equal(t, "axios", version.Name)
	assert.Equal(t, "1.5.0", version.Version)
}

func TestGetPackageVersionNetworkError(t *testing.T) {
	// 使用无效 URL 强制 getBytes 返回错误，覆盖 error 分支
	registry := NewRegistry(NewOptions().SetRegistryURL("http://localhost:1"))
	_, err := registry.GetPackageVersion(context.Background(), "axios", "1.0.0")
	assert.NotNil(t, err, "无效 URL 应该返回错误")
}

func TestDownloadTarballNetworkError(t *testing.T) {
	// 使用无效 URL 强制 getBytes 返回错误
	registry := NewRegistry(NewOptions().SetRegistryURL("http://localhost:1"))
	err := registry.DownloadTarball(context.Background(), "axios", "1.0.0", "/tmp/test.tgz")
	assert.NotNil(t, err, "无效 URL 应该返回错误")
}

func TestDownloadTarballEmptyTarballURL(t *testing.T) {
	// 测试 tarball URL 为空的情况
	// 创建一个只返回空 dist 的 mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"name": "axios",
			"version": "1.0.0",
			"description": "Promise based HTTP client",
			"dist": {}
		}`))
	}))
	defer server.Close()

	registry := NewRegistry(NewOptions().SetRegistryURL(server.URL))
	err := registry.DownloadTarball(context.Background(), "axios", "1.0.0", "/tmp/test.tgz")
	assert.NotNil(t, err, "空 tarball URL 应该返回错误")
	assert.Contains(t, err.Error(), "tarball URL not found")
}

func TestDownloadTarballSuccess(t *testing.T) {
	server := mockTestServer()
	defer server.Close()

	registry := NewRegistry(NewOptions().SetRegistryURL(server.URL))

	// 创建临时目录
	tmpDir := t.TempDir()
	destPath := tmpDir + "/axios-1.0.0.tgz"

	// 测试成功下载（tarball URL 指向 mock 服务器自身）
	err := registry.DownloadTarball(context.Background(), "axios", "1.0.0", destPath)
	assert.Nil(t, err, "下载应该成功")

	// 验证文件已创建
	data, readErr := os.ReadFile(destPath)
	assert.Nil(t, readErr, "文件应该可读")
	assert.NotEmpty(t, data, "文件应有内容")
	// 验证 gzip 头部
	assert.Equal(t, byte(0x1f), data[0])
	assert.Equal(t, byte(0x8b), data[1])
}
