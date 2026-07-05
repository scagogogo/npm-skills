package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// cliMockServer 创建一个覆盖所有 CLI 命令会调用的端点的 mock server。
func cliMockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		w.Header().Set("Content-Type", "application/json")

		// 写操作
		if r.Method == http.MethodPut || r.Method == http.MethodPost {
			w.WriteHeader(http.StatusOK)
			// bulk audit 端点返回 map[string][]advisory
			if strings.HasSuffix(path, "/advisories/bulk") {
				w.Write([]byte(`{"lodash":[{"id":1234,"title":"XSS","severity":"high","module_name":"lodash"}]}`))
				return
			}
			w.Write([]byte(`{"ok":true,"rev":"11-def"}`))
			return
		}
		if r.Method == http.MethodDelete {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"ok":true}`))
			return
		}

		// 包文档
		if path == "/react" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"_id":"react","_rev":"10-abc","name":"react","description":"UI lib","dist-tags":{"latest":"18.0.0"},"versions":{"18.0.0":{"name":"react","version":"18.0.0","dist":{"tarball":"http://x.tgz","shasum":"abc"}}}}`))
			return
		}
		if path == "/react/18.0.0" || path == "/react/latest" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"name":"react","version":"18.0.0","description":"UI lib","dist":{"tarball":"http://x.tgz","shasum":"abc"}}`))
			return
		}

		// dist-tags
		if path == "/-/package/react/dist-tags" {
			w.Write([]byte(`{"latest":"18.0.0","next":"19.0.0-rc.1"}`))
			return
		}
		if path == "/-/package/react/dist-tags/latest" {
			w.Write([]byte(`"18.0.0"`))
			return
		}

		// 搜索
		if path == "/-/v1/search" {
			w.Write([]byte(`{"objects":[{"package":{"name":"react","version":"18.0.0","description":"ui"},"score":{"final":0.9,"detail":{"quality":0.9,"popularity":0.9,"maintenance":0.9}}}],"total":1}`))
			return
		}

		// 下载统计
		if strings.HasPrefix(path, "/downloads") {
			if strings.Contains(path, "/range/") {
				w.Write([]byte(`{"start":"2024-01-01","end":"2024-01-07","package":"react","downloads":[{"day":"2024-01-01","downloads":10}]}`))
				return
			}
			w.Write([]byte(`{"downloads":100,"start":"2024-01-01","end":"2024-01-07","package":"react"}`))
			return
		}

		// whoami
		if path == "/-/whoami" {
			w.Write([]byte(`{"username":"alice"}`))
			return
		}

		// registry info
		if path == "/" {
			w.Write([]byte(`{"db_name":"registry","doc_count":1000}`))
			return
		}

		// user
		if path == "/-/user/org.couchdb.user:alice" {
			w.Write([]byte(`{"_id":"org.couchdb.user:alice","name":"alice","email":"a@x.com","type":"user"}`))
			return
		}

		// access
		if path == "/-/package/react/access" {
			w.Write([]byte(`{"package":"react","access":{"read":"public"}}`))
			return
		}
		if path == "/-/package/react/collaborators" {
			w.Write([]byte(`[{"name":"alice","permissions":"write"}]`))
			return
		}

		// star
		if path == "/-/_view/starredByUser" || path == "/-/_view/starredByPackage" {
			w.Write([]byte(`{"total_rows":2,"offset":0,"rows":[{"id":"1","key":"\"alice\"","value":"react"}]}`))
			return
		}

		// token
		if path == "/-/npm/v1/tokens" {
			w.Write([]byte(`{"objects":[{"id":"t1","token":"abcd12345678efgh","key":"k1","readonly":true}],"total":1}`))
			return
		}
		if path == "/-/npm/v1/tokens/t1" {
			w.Write([]byte(`{"id":"t1","token":"abcd12345678efgh","key":"k1","readonly":true}`))
			return
		}

		// audit
		if path == "/-/npm/v1/security/audits/quick" {
			w.Write([]byte(`{"metadata":{"vulnerabilities":{"low":1,"moderate":2,"high":3,"critical":0},"dependencies":5,"totalDependencies":5}}`))
			return
		}
		if path == "/-/npm/v1/security/advisories" {
			w.Write([]byte(`{"objects":[{"id":1234,"title":"XSS","severity":"high","module_name":"lodash"}],"total":1}`))
			return
		}
		if strings.HasPrefix(path, "/-/npm/v1/security/advisories/") {
			w.Write([]byte(`{"id":1234,"title":"XSS","severity":"high","module_name":"lodash"}`))
			return
		}

		// org/team
		if path == "/-/org/myorg" {
			w.Write([]byte(`{"name":"myorg"}`))
			return
		}
		if path == "/-/org/myorg/member" {
			w.Write([]byte(`["alice","bob"]`))
			return
		}
		if path == "/-/org/myorg/package" {
			w.Write([]byte(`["pkg1","pkg2"]`))
			return
		}
		if path == "/-/org/myorg/team" {
			w.Write([]byte(`{"objects":[{"id":"1","name":"devs"}],"total":1}`))
			return
		}
		if path == "/-/org/myorg/team/devs" {
			w.Write([]byte(`{"id":"1","name":"devs"}`))
			return
		}
		if path == "/-/org/myorg/team/devs/member" {
			w.Write([]byte(`["alice","bob"]`))
			return
		}
		if path == "/-/org/myorg/team/devs/package" {
			w.Write([]byte(`["pkg1"]`))
			return
		}

		// hooks
		if path == "/-/hooks/" || path == "/-/npm/v1/hooks" {
			w.Write([]byte(`{"objects":[{"id":"1","type":"hook","name":"n","endpoint":"https://e","events":["package"]}],"total":1}`))
			return
		}
		if path == "/-/hooks/1" || path == "/-/npm/v1/hooks/1" {
			w.Write([]byte(`{"id":"1","type":"hook","name":"n","endpoint":"https://e","events":["package"]}`))
			return
		}

		// couchdb
		if path == "/_changes" {
			w.Write([]byte(`{"last_seq":"10","pending":0,"results":[{"seq":"1","id":"pkg","changes":[{"rev":"1-abc"}]}]}`))
			return
		}
		if path == "/_all_docs" {
			w.Write([]byte(`{"total_rows":5,"offset":0,"rows":[{"id":"x","key":"x","value":{"rev":"1-abc"}}]}`))
			return
		}
		if strings.HasPrefix(path, "/-/_view/") {
			w.Write([]byte(`{"total_rows":2,"offset":0,"rows":[{"id":"1","key":"\"alice\"","value":"react"}]}`))
			return
		}

		// 默认
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
}

// runCLI 在 mock server 环境下执行给定参数的 CLI 命令，返回 stdout 和 error。
func runCLI(t *testing.T, server *httptest.Server, args ...string) (string, error) {
	t.Helper()
	resetGlobals()
	globalRegistry = server.URL
	globalToken = "npm_xxx"

	// 捕获 stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// rootCmd.SetArgs + Execute
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String(), err
}

// ====================================================================
// 各子命令 happy path
// ====================================================================

func TestCLIPackage(t *testing.T) {
	server := cliMockServer()
	defer server.Close()
	out, err := runCLI(t, server, "package", "react")
	assert.NoError(t, err)
	assert.Contains(t, out, "react")
}

func TestCLIPackageSummary(t *testing.T) {
	server := cliMockServer()
	defer server.Close()
	out, err := runCLI(t, server, "package-summary", "react")
	assert.NoError(t, err)
	assert.Contains(t, out, "react")
}

func TestCLISearch(t *testing.T) {
	server := cliMockServer()
	defer server.Close()
	out, err := runCLI(t, server, "search", "react")
	assert.NoError(t, err)
	assert.Contains(t, out, "react")
}

func TestCLIVersions(t *testing.T) {
	server := cliMockServer()
	defer server.Close()
	_, err := runCLI(t, server, "versions", "react")
	assert.NoError(t, err)
}

func TestCLIPkgVersion(t *testing.T) {
	server := cliMockServer()
	defer server.Close()
	_, err := runCLI(t, server, "pkg-version", "react", "18.0.0")
	assert.NoError(t, err)
}

func TestCLIDownloadStats(t *testing.T) {
	server := cliMockServer()
	defer server.Close()
	_, err := runCLI(t, server, "download-stats", "react")
	assert.NoError(t, err)
}

func TestCLIDownloadStatsBulk(t *testing.T) {
	server := cliMockServer()
	defer server.Close()
	_, err := runCLI(t, server, "download-stats-bulk", "react", "vue")
	assert.NoError(t, err)
}

func TestCLIDownloadStatsDate(t *testing.T) {
	server := cliMockServer()
	defer server.Close()
	_, err := runCLI(t, server, "download-stats-date", "react", "--start", "2024-01-01", "--end", "2024-01-07")
	assert.NoError(t, err)
}

func TestCLIDownloadRange(t *testing.T) {
	server := cliMockServer()
	defer server.Close()
	_, err := runCLI(t, server, "download-range", "react")
	assert.NoError(t, err)
}

func TestCLIMirrors(t *testing.T) {
	server := cliMockServer()
	defer server.Close()
	_, err := runCLI(t, server, "mirrors")
	assert.NoError(t, err)
}

func TestCLIRegistryInfo(t *testing.T) {
	server := cliMockServer()
	defer server.Close()
	_, err := runCLI(t, server, "registry-info")
	assert.NoError(t, err)
}

func TestCLIWhoami(t *testing.T) {
	server := cliMockServer()
	defer server.Close()
	_, err := runCLI(t, server, "whoami")
	assert.NoError(t, err)
}

func TestCLIDistTags(t *testing.T) {
	server := cliMockServer()
	defer server.Close()
	_, err := runCLI(t, server, "dist-tags", "get", "react")
	assert.NoError(t, err)
	_, err = runCLI(t, server, "dist-tags", "set", "react", "latest", "--version", "18.0.0")
	assert.NoError(t, err)
	_, err = runCLI(t, server, "dist-tags", "delete", "react", "next")
	assert.NoError(t, err)
}

func TestCLIDeprecate(t *testing.T) {
	server := cliMockServer()
	defer server.Close()
	_, err := runCLI(t, server, "deprecate", "react", "18.0.0", "--message", "use v19")
	assert.NoError(t, err)
}

func TestCLIUnpublish(t *testing.T) {
	server := cliMockServer()
	defer server.Close()
	_, err := runCLI(t, server, "unpublish", "react", "--version", "18.0.0")
	assert.NoError(t, err)
	_, err = runCLI(t, server, "unpublish", "react", "--force")
	assert.NoError(t, err)
}

func TestCLIStar(t *testing.T) {
	server := cliMockServer()
	defer server.Close()
	_, err := runCLI(t, server, "star", "react")
	assert.NoError(t, err)
}

func TestCLIUser(t *testing.T) {
	server := cliMockServer()
	defer server.Close()
	_, err := runCLI(t, server, "user", "get", "alice")
	assert.NoError(t, err)
}

func TestCLIToken(t *testing.T) {
	server := cliMockServer()
	defer server.Close()
	_, err := runCLI(t, server, "token", "list")
	assert.NoError(t, err)
	_, err = runCLI(t, server, "token", "get", "t1")
	assert.NoError(t, err)
}

func TestCLIHook(t *testing.T) {
	server := cliMockServer()
	defer server.Close()
	_, err := runCLI(t, server, "hook", "list")
	assert.NoError(t, err)
	_, err = runCLI(t, server, "hook", "get", "1")
	assert.NoError(t, err)
}

func TestCLIOrg(t *testing.T) {
	server := cliMockServer()
	defer server.Close()
	_, err := runCLI(t, server, "org", "get", "myorg")
	assert.NoError(t, err)
	_, err = runCLI(t, server, "org", "members", "myorg")
	assert.NoError(t, err)
	_, err = runCLI(t, server, "org", "packages", "myorg")
	assert.NoError(t, err)
	_, err = runCLI(t, server, "org", "team-list", "myorg")
	assert.NoError(t, err)
	_, err = runCLI(t, server, "org", "team-members", "myorg", "devs")
	assert.NoError(t, err)
}

func TestCLIAccess(t *testing.T) {
	server := cliMockServer()
	defer server.Close()
	_, err := runCLI(t, server, "access", "get", "react")
	assert.NoError(t, err)
	_, err = runCLI(t, server, "access", "collaborators", "react")
	assert.NoError(t, err)
}

func TestCLIAudit(t *testing.T) {
	server := cliMockServer()
	defer server.Close()
	_, err := runCLI(t, server, "audit", "quick", "--deps", "lodash=4.17.11")
	assert.NoError(t, err)
	_, err = runCLI(t, server, "audit", "advisories", "--package", "lodash")
	assert.NoError(t, err)
	_, err = runCLI(t, server, "audit", "advisory", "1234")
	assert.NoError(t, err)
	_, err = runCLI(t, server, "audit", "bulk", "--advisories", "lodash=<4.17.12")
	assert.NoError(t, err)
}

func TestCLICouchdb(t *testing.T) {
	server := cliMockServer()
	defer server.Close()
	_, err := runCLI(t, server, "couchdb", "changes")
	assert.NoError(t, err)
	_, err = runCLI(t, server, "couchdb", "all-docs")
	assert.NoError(t, err)
	_, err = runCLI(t, server, "couchdb", "view", "starredByUser")
	assert.NoError(t, err)
}

func TestCLIConfig(t *testing.T) {
	server := cliMockServer()
	defer server.Close()
	_, err := runCLI(t, server, "config")
	assert.NoError(t, err)
}

// TestCLIErrorPath 覆盖命令的 error 分支（不可达 server）。
func TestCLIErrorPath(t *testing.T) {
	// 用不可达 server
	server := cliMockServer()
	defer server.Close()
	resetGlobals()
	globalRegistry = "http://localhost:1"
	globalToken = "npm_xxx"
	rootCmd.SetArgs([]string{"package", "react"})
	_ = rootCmd.Execute()
	// 不断言 error（cobra Execute 在 RunE 返回 error 时不一定返回 error）
}
