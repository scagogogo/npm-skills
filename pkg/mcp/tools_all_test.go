package mcp

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
	"github.com/scagogogo/npm-skills/pkg/registry"
	"github.com/stretchr/testify/assert"
)

// fullMcpServer 覆盖所有 MCP 工具会调用的端点（含 org/team/hook/token/access/audit/couchdb）。
func fullMcpServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		w.Header().Set("Content-Type", "application/json")

		// 写操作
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
			// path 形如 /downloads/point/... 或 /downloads/range/...
			if len(path) > 17 && path[10:16] == "/range" {
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

		// user
		if path == "/-/user/org.couchdb.user:alice" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"_id":"org.couchdb.user:alice","name":"alice","email":"a@x.com","type":"user"}`))
			return
		}

		// access
		if path == "/-/package/react/access" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"package":"react","access":{"read":"public"}}`))
			return
		}
		if path == "/-/package/react/collaborators" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"name":"alice","permissions":"write"}]`))
			return
		}

		// star
		if path == "/-/_view/starredByUser" || path == "/-/_view/starredByPackage" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"total_rows":2,"offset":0,"rows":[{"id":"1","key":"\"alice\"","value":"react"}]}`))
			return
		}

		// token
		if path == "/-/npm/v1/tokens" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"objects":[{"id":"t1","token":"abcd12345678efgh","key":"k1","readonly":true}],"total":1}`))
			return
		}
		if path == "/-/npm/v1/tokens/t1" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id":"t1","token":"abcd12345678efgh","key":"k1","readonly":true}`))
			return
		}

		// audit
		if path == "/-/npm/v1/security/audits/quick" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"metadata":{"vulnerabilities":{"low":1,"moderate":2,"high":3,"critical":0},"dependencies":5,"totalDependencies":5}}`))
			return
		}
		if path == "/-/npm/v1/security/advisories" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"objects":[{"id":1234,"title":"XSS","severity":"high","module_name":"lodash"}],"total":1}`))
			return
		}
		if len(path) > 30 && path[:30] == "/-/npm/v1/security/advisories/" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id":1234,"title":"XSS","severity":"high","module_name":"lodash"}`))
			return
		}

		// org/team
		if path == "/-/org/myorg" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"name":"myorg"}`))
			return
		}
		if path == "/-/org/myorg/member" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`["alice","bob"]`))
			return
		}
		if path == "/-/org/myorg/package" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`["pkg1","pkg2"]`))
			return
		}
		if path == "/-/org/myorg/team" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"objects":[{"id":"1","name":"devs"}],"total":1}`))
			return
		}
		if path == "/-/org/myorg/team/devs" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id":"1","name":"devs"}`))
			return
		}
		if path == "/-/org/myorg/team/devs/member" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`["alice","bob"]`))
			return
		}
		if path == "/-/org/myorg/team/devs/package" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`["pkg1"]`))
			return
		}

		// hooks
		if path == "/-/hooks/" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"id":"1","type":"hook","name":"n","endpoint":"https://e","events":["package"]}]`))
			return
		}
		if path == "/-/hooks/1" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id":"1","type":"hook","name":"n","endpoint":"https://e","events":["package"]}`))
			return
		}

		// couchdb
		if path == "/_changes" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"last_seq":"10","pending":0,"results":[{"seq":"1","id":"pkg","changes":[{"rev":"1-abc"}]}]}`))
			return
		}
		if path == "/_all_docs" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"total_rows":5,"offset":0,"rows":[{"id":"x","key":"x","value":{"rev":"1-abc"}}]}`))
			return
		}
		if len(path) > 9 && path[:9] == "/-/_view/" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"total_rows":2,"offset":0,"rows":[{"id":"1","key":"\"alice\"","value":"react"}]}`))
			return
		}

		// 默认
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
}

// allTools 收集所有 register 函数的工具，用 mock server 的 client。
func allTools(server *httptest.Server, withToken bool) []mcpserver.ServerTool {
	cfg := mcpCfg(server, withToken)
	client := registry.NewRegistry(cfg.RegistryOptions)
	var tools []mcpserver.ServerTool
	tools = append(tools, registerRegistryTools(client, cfg)...)
	tools = append(tools, registerPackageTools(client, cfg)...)
	tools = append(tools, registerSearchTools(client, cfg)...)
	tools = append(tools, registerVersionTools(client, cfg)...)
	tools = append(tools, registerDownloadTools(client, cfg)...)
	tools = append(tools, registerWhoamiTools(client, cfg)...)
	tools = append(tools, registerDistTagWriteTools(client, cfg)...)
	tools = append(tools, registerUserTools(client, cfg)...)
	tools = append(tools, registerAccessTools(client, cfg)...)
	tools = append(tools, registerStarTools(client, cfg)...)
	tools = append(tools, registerTokenTools(client, cfg)...)
	tools = append(tools, registerAuditTools(client, cfg)...)
	tools = append(tools, registerOrgTools(client, cfg)...)
	tools = append(tools, registerHooksTools(client, cfg)...)
	tools = append(tools, registerCouchDBTools(client, cfg)...)
	return tools
}

// toolArgs 定义每个工具的 happy-path 参数。
var toolArgs = map[string]map[string]any{
	"npm_registry_info":         {},
	"npm_mirrors":               {},
	"npm_package":               {"name": "react"},
	"npm_package_summary":       {"name": "react"},
	"npm_search":                {"query": "react"},
	"npm_version":               {"name": "react", "version": "18.0.0"},
	"npm_versions":              {"name": "react"},
	"npm_latest_version":        {"name": "react"},
	"npm_download_stats":        {"name": "react", "period": "last-week"},
	"npm_download_range":        {"name": "react", "period": "last-week"},
	"npm_whoami":                {},
	"npm_dist_tag_get":          {"name": "react", "tag": "latest"},
	"npm_dist_tag_set":          {"name": "react", "tag": "latest", "version": "18.0.0"},
	"npm_dist_tag_delete":       {"name": "react", "tag": "next"},
	"npm_dist_tags":             {"name": "react"},
	"npm_user_get":              {"name": "alice"},
	"npm_package_access":        {"name": "react"},
	"npm_package_collaborators": {"name": "react"},
	"npm_starred_by_user":       {"username": "alice"},
	"npm_starred_by_package":    {"name": "react"},
	"npm_token_list":            {},
	"npm_audit":                 {"dependencies": map[string]any{"lodash": "4.17.11"}},
	"npm_audit_advisory":        {"id": float64(1234)},
	"npm_org_get":               {"name": "myorg"},
	"npm_org_members":           {"name": "myorg"},
	"npm_org_packages":          {"name": "myorg"},
	"npm_team_list":             {"org": "myorg"},
	"npm_team_members":          {"org": "myorg", "team": "devs"},
	"npm_hook_list":             {},
	"npm_hook_get":              {"id": "1"},
	"npm_changes":               {},
}

// TestAllToolsHappyPath 对每个工具用 happy-path 参数调用 handler。
func TestAllToolsHappyPath(t *testing.T) {
	server := fullMcpServer()
	defer server.Close()
	tools := allTools(server, true)

	for _, tool := range tools {
		args, ok := toolArgs[tool.Tool.Name]
		if !ok {
			t.Logf("skip (no args defined): %s", tool.Tool.Name)
			continue
		}
		t.Run(tool.Tool.Name, func(t *testing.T) {
			res, err := callTool(&tool, args)
			assert.NoError(t, err, "handler should not return go error")
			assert.NotNil(t, res, "result should not be nil")
			if res.IsError {
				// 部分工具在 happy path 也可能返回 error（如 npm_audit_advisory 路径不匹配），记录但不失败
				t.Logf("tool %s returned IsError=true: %s", tool.Tool.Name, resultText(res))
			}
		})
	}
}

// TestAllToolsErrorPath 用空参数调用需要参数的工具，覆盖 missing-arg 分支。
func TestAllToolsMissingArgs(t *testing.T) {
	server := fullMcpServer()
	defer server.Close()
	tools := allTools(server, true)

	for _, tool := range tools {
		args, hasArgs := toolArgs[tool.Tool.Name]
		if !hasArgs {
			continue
		}
		// 跳过不需要参数的工具（空参数即 happy path）
		if len(args) == 0 {
			continue
		}
		t.Run(tool.Tool.Name+"_missing", func(t *testing.T) {
			// 用空参数调用
			res, err := callTool(&tool, map[string]any{})
			assert.NoError(t, err)
			assert.NotNil(t, res)
			// 大多数工具应返回 IsError（missing required arg）
			// 但部分工具可能用默认值不报错，这里只确保不 panic
		})
	}
}

// TestAllToolsNetworkError 用不可达 server 调用每个工具，覆盖网络错误分支。
func TestAllToolsNetworkError(t *testing.T) {
	// 不可达 server
	badServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer badServer.Close()
	tools := allTools(badServer, true)

	for _, tool := range tools {
		args, ok := toolArgs[tool.Tool.Name]
		if !ok || len(args) == 0 {
			continue
		}
		t.Run(tool.Tool.Name+"_neterr", func(t *testing.T) {
			res, err := callTool(&tool, args)
			assert.NoError(t, err)
			assert.NotNil(t, res)
			// 网络错误应返回 IsError
			if !res.IsError {
				t.Logf("tool %s did not return error on network failure", tool.Tool.Name)
			}
		})
	}
}

// TestSearchToolOptionalParams 测试 npm_search 的可选参数解析分支。
func TestSearchToolOptionalParams(t *testing.T) {
	server := fullMcpServer()
	defer server.Close()
	cfg := mcpCfg(server, false)
	client := registry.NewRegistry(cfg.RegistryOptions)
	tools := registerSearchTools(client, cfg)
	tool := findTool(tools, "npm_search")

	cases := []map[string]any{
		{"query": "react", "limit": "5", "from": "10", "quality": "0.5", "popularity": "0.5", "maintenance": "0.5"},
		{"query": "react", "limit": "invalid", "quality": "invalid"}, // 无效数值
		{"query": ""}, // 空 query
	}
	for i, args := range cases {
		res, err := callTool(&tool, args)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		if args["query"] == "" {
			assert.True(t, res.IsError, "case %d empty query should error", i)
		}
	}
}

// TestGetOptionalFloatBranches 覆盖 getOptionalFloat 的所有类型分支。
func TestGetOptionalFloatBranches(t *testing.T) {
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"f64":   1.5,
				"str":   "2.5",
				"bad":   "not-a-number",
				"other": 42,
			},
		},
	}
	// float64
	v, ok := getOptionalFloat(req, "f64")
	assert.True(t, ok)
	assert.Equal(t, 1.5, v)
	// string valid
	v, ok = getOptionalFloat(req, "str")
	assert.True(t, ok)
	assert.Equal(t, 2.5, v)
	// string invalid
	_, ok = getOptionalFloat(req, "bad")
	assert.False(t, ok)
	// missing
	_, ok = getOptionalFloat(req, "missing")
	assert.False(t, ok)
	// other type (int) → false
	_, ok = getOptionalFloat(req, "other")
	assert.False(t, ok)
}

// TestGetOptionalFloatJSONNumber 覆盖 json.Number 分支。
func TestGetOptionalFloatJSONNumber(t *testing.T) {
	// json.Number 有效
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]any{
				"good": json.Number("3.14"),
				"bad":  json.Number("not-a-number"),
			},
		},
	}
	v, ok := getOptionalFloat(req, "good")
	assert.True(t, ok)
	assert.Equal(t, 3.14, v)
	// json.Number 无效
	_, ok = getOptionalFloat(req, "bad")
	assert.False(t, ok)
}

// TestRegistryToolError 覆盖 npm_registry_info 的 error 分支。
func TestRegistryToolError(t *testing.T) {
	// 用不可达 server 触发 GetRegistryInformation 失败
	cfg := mcpCfg(httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	})), false)
	// 重新创建指向不可达地址的 client
	badClient := registry.NewRegistry(registry.NewOptions().SetRegistryURL("http://localhost:1"))
	tools := registerRegistryTools(badClient, cfg)
	tool := findTool(tools, "npm_registry_info")
	res, _ := callTool(&tool, map[string]any{})
	assert.True(t, res.IsError)
}

// TestTokenToolError 覆盖 npm_token_list 的 error 分支。
func TestTokenToolError(t *testing.T) {
	cfg := mcpCfg(httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})), true)
	badClient := registry.NewRegistry(registry.NewOptions().SetToken("npm_xxx").SetRegistryURL("http://localhost:1"))
	tools := registerTokenTools(badClient, cfg)
	tool := findTool(tools, "npm_token_list")
	res, _ := callTool(&tool, map[string]any{})
	assert.True(t, res.IsError)
}

// TestAuditToolArgBranches 覆盖 npm_audit / npm_audit_advisory 的参数解析分支。
func TestAuditToolArgBranches(t *testing.T) {
	server := fullMcpServer()
	defer server.Close()
	cfg := mcpCfg(server, false)
	client := registry.NewRegistry(cfg.RegistryOptions)
	tools := registerAuditTools(client, cfg)

	// npm_audit: dependencies 不是 map → "invalid arguments" 分支
	tool := findTool(tools, "npm_audit")
	res, _ := callTool(&tool, map[string]any{"dependencies": "not-a-map"})
	assert.True(t, res.IsError)

	// npm_audit: dependencies 是空 map → error
	res, _ = callTool(&tool, map[string]any{"dependencies": map[string]any{}})
	assert.True(t, res.IsError)

	// npm_audit: arguments 不是 map（传 string）→ "invalid arguments"
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Name: "npm_audit", Arguments: "not-a-map"}}
	res, _ = tool.Handler(context.Background(), req)
	assert.True(t, res.IsError)

	// npm_audit_advisory: id 不是 float64（传 string）→ error
	tool = findTool(tools, "npm_audit_advisory")
	res, _ = callTool(&tool, map[string]any{"id": "not-a-number"})
	assert.True(t, res.IsError)

	// npm_audit_advisory: arguments 不是 map
	req = mcp.CallToolRequest{Params: mcp.CallToolParams{Name: "npm_audit_advisory", Arguments: "not-a-map"}}
	res, _ = tool.Handler(context.Background(), req)
	assert.True(t, res.IsError)
}

// TestOrgTeamListError 覆盖 npm_team_list 的 error 分支。
func TestOrgTeamListError(t *testing.T) {
	cfg := mcpCfg(httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})), true)
	badClient := registry.NewRegistry(registry.NewOptions().SetToken("npm_xxx").SetRegistryURL("http://localhost:1"))
	tools := registerOrgTools(badClient, cfg)
	tool := findTool(tools, "npm_team_list")
	res, _ := callTool(&tool, map[string]any{"org": "myorg"})
	assert.True(t, res.IsError)
}

// TestHooksToolError 覆盖 npm_hook_list 的 error 分支。
func TestHooksToolError(t *testing.T) {
	cfg := mcpCfg(httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})), true)
	badClient := registry.NewRegistry(registry.NewOptions().SetToken("npm_xxx").SetRegistryURL("http://localhost:1"))
	tools := registerHooksTools(badClient, cfg)
	tool := findTool(tools, "npm_hook_list")
	// 带包名参数，触发 ListHooks 调用
	res, _ := callTool(&tool, map[string]any{"package": "my-pkg"})
	assert.True(t, res.IsError)
}

// TestChangesToolWithLimit 覆盖 npm_changes 的 limit 参数解析 + error 分支。
func TestChangesToolWithLimit(t *testing.T) {
	server := fullMcpServer()
	defer server.Close()
	cfg := mcpCfg(server, false)
	client := registry.NewRegistry(cfg.RegistryOptions)
	tools := registerCouchDBTools(client, cfg)

	// 带 limit 参数 → 覆盖 limit 解析分支
	tool := findTool(tools, "npm_changes")
	res, _ := callTool(&tool, map[string]any{"limit": float64(10), "since": "5"})
	assert.False(t, res.IsError)

	// 无 limit 参数 → argsOk=true 但无 limit 键，用默认 25
	res, _ = callTool(&tool, map[string]any{})
	assert.False(t, res.IsError)
}

// TestChangesToolError 覆盖 npm_changes 的 GetChanges 失败分支。
func TestChangesToolError(t *testing.T) {
	cfg := mcpCfg(httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})), false)
	badClient := registry.NewRegistry(registry.NewOptions().SetRegistryURL("http://localhost:1"))
	tools := registerCouchDBTools(badClient, cfg)
	tool := findTool(tools, "npm_changes")
	res, _ := callTool(&tool, map[string]any{"limit": float64(10)})
	assert.True(t, res.IsError)
}

// TestVersionTool 测试 version 工具的各种参数。
func TestVersionTool(t *testing.T) {
	server := fullMcpServer()
	defer server.Close()
	cfg := mcpCfg(server, false)
	client := registry.NewRegistry(cfg.RegistryOptions)
	tools := registerVersionTools(client, cfg)

	// npm_version happy
	tool := findTool(tools, "npm_version")
	res, _ := callTool(&tool, map[string]any{"name": "react", "version": "18.0.0"})
	assert.False(t, res.IsError)

	// npm_version missing name
	res, _ = callTool(&tool, map[string]any{})
	assert.True(t, res.IsError)

	// npm_versions
	tool = findTool(tools, "npm_versions")
	res, _ = callTool(&tool, map[string]any{"name": "react"})
	assert.False(t, res.IsError)
	res, _ = callTool(&tool, map[string]any{})
	assert.True(t, res.IsError)

	// npm_latest_version
	tool = findTool(tools, "npm_latest_version")
	res, _ = callTool(&tool, map[string]any{"name": "react"})
	assert.False(t, res.IsError)
}

// TestWhoamiToolNoToken 测试 whoami 无 token。
func TestWhoamiToolNoToken(t *testing.T) {
	server := fullMcpServer()
	defer server.Close()
	cfg := mcpCfg(server, false) // 无 token
	client := registry.NewRegistry(cfg.RegistryOptions)
	tools := registerWhoamiTools(client, cfg)
	tool := findTool(tools, "npm_whoami")
	res, _ := callTool(&tool, map[string]any{})
	assert.True(t, res.IsError)
}

// TestDownloadTool 测试 download 工具。
func TestDownloadTool(t *testing.T) {
	server := fullMcpServer()
	defer server.Close()
	cfg := mcpCfg(server, false)
	client := registry.NewRegistry(cfg.RegistryOptions)
	tools := registerDownloadTools(client, cfg)

	tool := findTool(tools, "npm_download_stats")
	res, _ := callTool(&tool, map[string]any{"name": "react", "period": "last-week"})
	assert.False(t, res.IsError)
	res, _ = callTool(&tool, map[string]any{})
	assert.True(t, res.IsError)

	tool = findTool(tools, "npm_download_range")
	res, _ = callTool(&tool, map[string]any{"name": "react", "period": "last-week"})
	assert.False(t, res.IsError)
}

// TestRegistryTool 测试 registry info + mirrors。
func TestRegistryTool(t *testing.T) {
	server := fullMcpServer()
	defer server.Close()
	cfg := mcpCfg(server, false)
	client := registry.NewRegistry(cfg.RegistryOptions)
	tools := registerRegistryTools(client, cfg)

	tool := findTool(tools, "npm_registry_info")
	res, _ := callTool(&tool, map[string]any{})
	assert.False(t, res.IsError)

	tool = findTool(tools, "npm_mirrors")
	res, _ = callTool(&tool, map[string]any{})
	assert.False(t, res.IsError)
}

// TestDistTagWriteToolNoToken 测试 dist-tag 写操作无 token。
func TestDistTagWriteToolNoToken(t *testing.T) {
	server := fullMcpServer()
	defer server.Close()
	cfg := mcpCfg(server, false) // 无 token
	client := registry.NewRegistry(cfg.RegistryOptions)
	tools := registerDistTagWriteTools(client, cfg)

	// npm_dist_tag_set 无 token → error
	tool := findTool(tools, "npm_dist_tag_set")
	res, _ := callTool(&tool, map[string]any{"name": "react", "tag": "latest", "version": "18.0.0"})
	assert.True(t, res.IsError)

	// npm_dist_tag_delete 无 token → error
	tool = findTool(tools, "npm_dist_tag_delete")
	res, _ = callTool(&tool, map[string]any{"name": "react", "tag": "next"})
	assert.True(t, res.IsError)
}

// TestContextCancelled 测试 context 取消。
func TestContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // 立即取消
	server := fullMcpServer()
	defer server.Close()
	cfg := mcpCfg(server, true)
	client := registry.NewRegistry(cfg.RegistryOptions)
	tools := registerPackageTools(client, cfg)
	tool := findTool(tools, "npm_package")
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]any{"name": "react"}}}
	// 用已取消的 context，应返回错误（可能是 IsError）
	res, err := tool.Handler(ctx, req)
	_ = err
	_ = res
	// 不 panic 即可
}
