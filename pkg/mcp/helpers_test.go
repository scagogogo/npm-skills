package mcp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
	"github.com/scagogogo/npm-skills/pkg/models"
	"github.com/scagogogo/npm-skills/pkg/registry"
	"github.com/stretchr/testify/assert"
)

// mcpMockServer 创建一个 mock NPM registry，覆盖 MCP 工具会调用的所有端点。
func mcpMockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		w.Header().Set("Content-Type", "application/json")

		// PUT 任意包路径（发布/写操作）返回成功
		if r.Method == http.MethodPut {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"ok":true,"rev":"11-def"}`))
			return
		}
		if r.Method == http.MethodDelete {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"ok":true}`))
			return
		}
		if r.Method == http.MethodPost {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"ok":true}`))
			return
		}

		// 包文档
		if path == "/react" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"_id":"react","_rev":"10-abc","name":"react","description":"UI lib","dist-tags":{"latest":"18.0.0"},"versions":{"18.0.0":{"name":"react","version":"18.0.0"}},"users":{"alice":true}}`))
			return
		}
		if path == "/react/18.0.0" || path == "/react/latest" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"name":"react","version":"18.0.0","description":"UI lib","dist":{"tarball":"http://x.tgz","shasum":"abc"}}`))
			return
		}

		// dist-tags
		if path == "/-/package/react/dist-tags" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"latest":"18.0.0","next":"19.0.0-rc.1"}`))
			return
		}
		if path == "/-/package/react/dist-tags/latest" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`"18.0.0"`))
			return
		}

		// 搜索
		if path == "/-/v1/search" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"objects":[{"package":{"name":"react","version":"18.0.0","description":"ui"},"score":{"final":0.9,"detail":{"quality":0.9,"popularity":0.9,"maintenance":0.9}},"searchScore":0.9}],"total":1}`))
			return
		}

		// 下载统计
		if len(path) > 10 && path[:10] == "/downloads" {
			w.WriteHeader(http.StatusOK)
			if path[8] == 'r' {
				w.Write([]byte(`{"start":"2024-01-01","end":"2024-01-07","package":"react","downloads":[{"day":"2024-01-01","downloads":10}]}`))
				return
			}
			w.Write([]byte(`{"downloads":100,"start":"2024-01-01","end":"2024-01-07","package":"react"}`))
			return
		}

		// whoami
		if path == "/-/whoami" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"username":"alice"}`))
			return
		}

		// registry info
		if path == "/" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"db_name":"registry","doc_count":1000}`))
			return
		}

		// 默认
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
}

// mcpRegistry 创建一个指向 mock server 的 Registry，可选 token。
func mcpRegistry(server *httptest.Server, withToken bool) *registry.Registry {
	o := registry.NewOptions().SetRegistryURL(server.URL)
	if withToken {
		o.SetToken("npm_xxx")
	}
	o.SetDownloadStatsURL(server.URL + "/downloads")
	return registry.NewRegistry(o)
}

func mcpCfg(server *httptest.Server, withToken bool) Config {
	return Config{
		RegistryOptions: mcpRegistry(server, withToken).GetOptions(),
		Timeout:         10 * time.Second,
	}
}

// callTool 调用工具 handler 并返回结果。
func callTool(t *mcpserver.ServerTool, args map[string]any) (*mcp.CallToolResult, error) {
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      t.Tool.Name,
			Arguments: args,
		},
	}
	return t.Handler(context.Background(), req)
}

// findTool 按名字在工具列表中查找。
func findTool(tools []mcpserver.ServerTool, name string) mcpserver.ServerTool {
	for _, t := range tools {
		if t.Tool.Name == name {
			return t
		}
	}
	panic("tool not found: " + name)
}

// resultText 提取工具结果的文本内容。
func resultText(r *mcp.CallToolResult) string {
	if len(r.Content) == 0 {
		return ""
	}
	if tc, ok := r.Content[0].(mcp.TextContent); ok {
		return tc.Text
	}
	return ""
}

// ====================================================================
// 验证 NewServer 注册所有工具不 panic
// ====================================================================

func TestNewServerAllToolsRegistered(t *testing.T) {
	server := mcpMockServer()
	defer server.Close()
	cfg := mcpCfg(server, true)
	s := NewServer(cfg)
	assert.NotNil(t, s)
}

// ====================================================================
// tools_package.go
// ====================================================================

func TestPackageTool(t *testing.T) {
	server := mcpMockServer()
	defer server.Close()
	cfg := mcpCfg(server, false)
	client := registry.NewRegistry(cfg.RegistryOptions)
	tools := registerPackageTools(client, cfg)

	// happy path
	tool := findTool(tools, "npm_package")
	res, err := callTool(&tool, map[string]any{"name": "react"})
	assert.NoError(t, err)
	assert.False(t, res.IsError)
	assert.Contains(t, resultText(res), "react")

	// missing name
	res, err = callTool(&tool, map[string]any{})
	assert.NoError(t, err)
	assert.True(t, res.IsError)

	// network error
	tool2 := findTool(registerPackageTools(registry.NewRegistry(registry.NewOptions().SetRegistryURL("http://localhost:1")), cfg), "npm_package")
	res, err = callTool(&tool2, map[string]any{"name": "react"})
	assert.NoError(t, err)
	assert.True(t, res.IsError)
}

func TestPackageSummaryTool(t *testing.T) {
	server := mcpMockServer()
	defer server.Close()
	cfg := mcpCfg(server, false)
	client := registry.NewRegistry(cfg.RegistryOptions)
	tools := registerPackageTools(client, cfg)
	tool := findTool(tools, "npm_package_summary")

	// happy
	res, _ := callTool(&tool, map[string]any{"name": "react"})
	assert.False(t, res.IsError)

	// missing name
	res, _ = callTool(&tool, map[string]any{})
	assert.True(t, res.IsError)

	// error
	tool2 := findTool(registerPackageTools(registry.NewRegistry(registry.NewOptions().SetRegistryURL("http://localhost:1")), cfg), "npm_package_summary")
	res, _ = callTool(&tool2, map[string]any{"name": "react"})
	assert.True(t, res.IsError)
}

func TestTruncatePackage(t *testing.T) {
	// 小包：原样返回
	pkg := &models.Package{Name: "react", Description: "ui"}
	res := truncatePackage(pkg)
	assert.NotNil(t, res)

	// 大包：触发截断（README 超长 + 大 Versions）
	bigPkg := &models.Package{
		Name:   "big",
		ReadMe: string(make([]byte, 5000)),
	}
	// 填充大量版本使序列化超过 50KB
	bigPkg.Versions = make(map[string]models.Version)
	for i := 0; i < 2000; i++ {
		bigPkg.Versions["1.0."+itoa(i)] = models.Version{Name: "big", Version: "1.0." + itoa(i)}
	}
	res = truncatePackage(bigPkg)
	assert.NotNil(t, res)
	// 截断后应是 map[string]any 且带 _truncated
	m, ok := res.(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, true, m["_truncated"])
	// README 应被截断
	if readme, ok := m["readme"].(string); ok {
		assert.Contains(t, readme, "truncated")
	}
}

// TestTruncatePackageMarshalFail 覆盖 json.Marshal 失败分支。
// Deprecated 字段是 interface{}，设为 chan 会触发 marshal 错误（93 行）。
func TestTruncatePackageMarshalFail(t *testing.T) {
	pkg := &models.Package{Name: "x", Deprecated: make(chan int)}
	res := truncatePackage(pkg)
	// marshal 失败 → 返回原 pkg
	assert.Same(t, pkg, res)
}

func itoa(i int) string {
	return formatInt(i)
}

func formatInt(i int) string {
	// 简单整数转字符串
	if i == 0 {
		return "0"
	}
	neg := false
	if i < 0 {
		neg = true
		i = -i
	}
	var buf [20]byte
	pos := len(buf)
	for i > 0 {
		pos--
		buf[pos] = byte('0' + i%10)
		i /= 10
	}
	if neg {
		pos--
		buf[pos] = '-'
	}
	return string(buf[pos:])
}
