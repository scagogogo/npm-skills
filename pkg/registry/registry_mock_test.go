package registry

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockTestServer creates a mock NPM registry server for testing
func mockTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// axios 版本路由 - GET /axios/1.0.0
		if r.URL.Path == "/axios/1.0.0" {
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
					"tarball": "https://registry.npmjs.org/axios/-/axios-1.0.0.tgz"
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
