package registry

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/scagogogo/npm-skills/pkg/models"
	"github.com/stretchr/testify/assert"
)

// orgMockServer 专门处理 org/team/member 端点
func orgMockServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		w.Header().Set("Content-Type", "application/json")

		// /-/org/{org}/team/{team}/member (GET)
		if path == "/-/org/myorg/team/devs/member" && r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`["alice","bob"]`))
			return
		}
		// /-/org/{org}/team (GET 列出团队)
		if path == "/-/org/myorg/team" && r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"id":"1","name":"devs"}]`))
			return
		}
		// /-/org/{org}/team/{team} (PUT 创建团队)
		if path == "/-/org/myorg/team/devs" && r.Method == http.MethodPut {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id":"1","name":"devs"}`))
			return
		}
		// /-/org/{org}/team/{team} (DELETE 删除团队)
		if path == "/-/org/myorg/team/devs" && r.Method == http.MethodDelete {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"ok":true}`))
			return
		}
		// /-/org/{org}/team/{team}/member (PUT 添加成员, DELETE 移除成员)
		if path == "/-/org/myorg/team/devs/member" && (r.Method == http.MethodPut || r.Method == http.MethodDelete) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"ok":true}`))
			return
		}
		// /-/org/{org} (GET/PUT/DELETE 组织)
		if path == "/-/org/myorg" {
			switch r.Method {
			case http.MethodGet:
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"name":"myorg"}`))
				return
			case http.MethodPut:
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"name":"myorg"}`))
				return
			case http.MethodDelete:
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"ok":true}`))
				return
			}
		}
		// /-/org/{org}/member (GET 列出成员, PUT 添加, DELETE 移除)
		if path == "/-/org/myorg/member" {
			switch r.Method {
			case http.MethodGet:
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`["alice","bob"]`))
				return
			case http.MethodPut, http.MethodDelete:
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"ok":true}`))
				return
			}
		}
		// /-/org/{org}/package (GET 列出包)
		if path == "/-/org/myorg/package" && r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`["pkg1","pkg2"]`))
			return
		}
		// /-/org/{org}/team/{team}/package (GET 列出团队包)
		if path == "/-/org/myorg/team/devs/package" && r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`["pkg1"]`))
			return
		}

		// hook 端点
		if path == "/-/hooks/" && r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"id":"1","type":"hook","name":"n","endpoint":"https://e","events":["package"]}]`))
			return
		}
		if path == "/-/hooks/" && r.Method == http.MethodPost {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id":"1","type":"hook","name":"n","endpoint":"https://e","events":["package"]}`))
			return
		}
		if path == "/-/hooks/1" {
			switch r.Method {
			case http.MethodGet:
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"id":"1","type":"hook","name":"n","endpoint":"https://e","events":["package"]}`))
				return
			case http.MethodPut:
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"id":"1","type":"hook","name":"n","endpoint":"https://e","events":["package"]}`))
				return
			case http.MethodDelete:
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"ok":true}`))
				return
			}
		}

		// access 端点
		if path == "/-/package/my-pkg/access" && r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"package":"my-pkg","access":{"read":"public"}}`))
			return
		}
		if path == "/-/package/my-pkg/access" && r.Method == http.MethodPut {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"ok":true}`))
			return
		}
		if path == "/-/package/my-pkg/collaborators" && r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"name":"alice","permissions":"write"}]`))
			return
		}
		if path == "/-/package/my-pkg/collaborators/alice" && (r.Method == http.MethodPut || r.Method == http.MethodDelete) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"ok":true}`))
			return
		}

		// token 端点
		if path == "/-/npm/v1/tokens" && r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"objects":[{"id":"t1","token":"abcd12345678efgh","key":"k1","readonly":true}],"total":1}`))
			return
		}
		if path == "/-/npm/v1/tokens" && r.Method == http.MethodPost {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id":"t1","token":"abcd12345678efgh","key":"k1","readonly":true}`))
			return
		}
		if path == "/-/npm/v1/tokens/t1" {
			switch r.Method {
			case http.MethodGet:
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"id":"t1","token":"abcd12345678efgh","key":"k1","readonly":true}`))
				return
			case http.MethodDelete:
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"ok":true}`))
				return
			}
		}

		// 默认 200 + 空对象
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
}

// ====================================================================
// org.go - 补充 ListTeamMembers 和 error 分支
// ====================================================================

func TestListTeamMembersCov(t *testing.T) {
	server := orgMockServer(t)
	defer server.Close()
	reg := NewRegistry(NewOptions().SetRegistryURL(server.URL).SetToken("npm_xxx"))
	members, err := reg.ListTeamMembers(context.Background(), "myorg", "devs")
	assert.NoError(t, err)
	assert.Contains(t, members, "alice")
}

func TestListTeamMembersNoToken(t *testing.T) {
	server := orgMockServer(t)
	defer server.Close()
	reg := NewRegistry(NewOptions().SetRegistryURL(server.URL))
	_, err := reg.ListTeamMembers(context.Background(), "myorg", "devs")
	assert.Error(t, err)
}

func TestListTeamMembersError(t *testing.T) {
	reg := NewRegistry(NewOptions().SetToken("npm_xxx").SetRegistryURL("http://localhost:1"))
	_, err := reg.ListTeamMembers(context.Background(), "myorg", "devs")
	assert.Error(t, err)
}

func TestOrgErrorBranches(t *testing.T) {
	reg := NewRegistry(NewOptions().SetToken("npm_xxx").SetRegistryURL("http://localhost:1"))
	_, err := reg.GetOrg(context.Background(), "myorg")
	assert.Error(t, err)
	_, err = reg.CreateOrg(context.Background(), "myorg")
	assert.Error(t, err)
	assert.Error(t, reg.DeleteOrg(context.Background(), "myorg"))
	_, err = reg.ListOrgMembers(context.Background(), "myorg")
	assert.Error(t, err)
	assert.Error(t, reg.AddOrgMember(context.Background(), "myorg", "alice"))
	assert.Error(t, reg.RemoveOrgMember(context.Background(), "myorg", "alice"))
	_, err = reg.ListOrgPackages(context.Background(), "myorg")
	assert.Error(t, err)
	_, err = reg.ListTeams(context.Background(), "myorg")
	assert.Error(t, err)
	_, err = reg.CreateTeam(context.Background(), "myorg", "devs")
	assert.Error(t, err)
	assert.Error(t, reg.DeleteTeam(context.Background(), "myorg", "devs"))
	assert.Error(t, reg.AddTeamMember(context.Background(), "myorg", "devs", "alice"))
	assert.Error(t, reg.RemoveTeamMember(context.Background(), "myorg", "devs", "alice"))
	_, err = reg.ListTeamPackages(context.Background(), "myorg", "devs")
	assert.Error(t, err)
}

// ====================================================================
// hooks.go - error 分支
// ====================================================================

func TestHooksErrorBranches(t *testing.T) {
	reg := NewRegistry(NewOptions().SetToken("npm_xxx").SetRegistryURL("http://localhost:1"))
	_, err := reg.ListHooks(context.Background(), models.HookListOptions{})
	assert.Error(t, err)
	_, err = reg.GetHook(context.Background(), "1")
	assert.Error(t, err)
	_, err = reg.CreateHook(context.Background(), &models.HookCreation{Name: "n", Endpoint: "https://e"})
	assert.Error(t, err)
	_, err = reg.UpdateHook(context.Background(), "1", &models.HookUpdate{})
	assert.Error(t, err)
	assert.Error(t, reg.DeleteHook(context.Background(), "1"))
}

func TestHooksNoToken(t *testing.T) {
	reg := NewRegistry(NewOptions())
	_, err := reg.ListHooks(context.Background(), models.HookListOptions{})
	assert.Error(t, err)
}

// ====================================================================
// access.go - error 分支
// ====================================================================

func TestAccessErrorBranches(t *testing.T) {
	reg := NewRegistry(NewOptions().SetToken("npm_xxx").SetRegistryURL("http://localhost:1"))
	_, err := reg.GetPackageAccess(context.Background(), "my-pkg")
	assert.Error(t, err)
	assert.Error(t, reg.SetPackageAccess(context.Background(), "my-pkg", &models.PackageAccessUpdate{Access: "public"}))
	_, err = reg.ListCollaborators(context.Background(), "my-pkg")
	assert.Error(t, err)
	assert.Error(t, reg.GrantAccess(context.Background(), "my-pkg", "alice", models.PermissionWrite))
	assert.Error(t, reg.RevokeAccess(context.Background(), "my-pkg", "alice"))
}

// ====================================================================
// token.go - error 分支
// ====================================================================

func TestTokenErrorBranches(t *testing.T) {
	reg := NewRegistry(NewOptions().SetToken("npm_xxx").SetRegistryURL("http://localhost:1"))
	_, err := reg.ListTokens(context.Background())
	assert.Error(t, err)
	_, err = reg.GetToken(context.Background(), "t1")
	assert.Error(t, err)
	_, err = reg.CreateToken(context.Background(), &models.TokenCreation{Password: "p"})
	assert.Error(t, err)
	assert.Error(t, reg.DeleteToken(context.Background(), "t1"))
}

func TestTokenNoToken(t *testing.T) {
	reg := NewRegistry(NewOptions())
	_, err := reg.ListTokens(context.Background())
	assert.Error(t, err)
	_, err = reg.GetToken(context.Background(), "t1")
	assert.Error(t, err)
}

// ====================================================================
// download_stats.go - error 分支（requireDownloadStatsURL）
// ====================================================================

func TestDownloadStatsPrivateRegistryError(t *testing.T) {
	// 私有仓库且未设 DownloadStatsURL → ErrDownloadStatsNotAvailable
	reg := NewCustomRegistry("https://npm.private.com")
	_, err := reg.GetDownloadRangeStats(context.Background(), "my-pkg", "last-week")
	assert.Error(t, err)
	_, err = reg.GetDownloadStatsByDateRange(context.Background(), "my-pkg", "2024-01-01", "2024-01-07")
	assert.Error(t, err)
	_, err = reg.GetDownloadRangeStatsByDateRange(context.Background(), "my-pkg", "2024-01-01", "2024-01-07")
	assert.Error(t, err)
	_, err = reg.GetBulkDownloadStats(context.Background(), []string{"my-pkg"}, "last-week")
	assert.Error(t, err)
	_, err = reg.GetBulkDownloadRangeStats(context.Background(), []string{"my-pkg"}, "last-week")
	assert.Error(t, err)
	_, err = reg.GetBulkDownloadStatsByDateRange(context.Background(), []string{"my-pkg"}, "2024-01-01", "2024-01-07")
	assert.Error(t, err)
	_, err = reg.GetBulkDownloadRangeStatsByDateRange(context.Background(), []string{"my-pkg"}, "2024-01-01", "2024-01-07")
	assert.Error(t, err)
}

// ====================================================================
// unpublish.go - sendJSON 失败分支
// ====================================================================

func TestUnpublishPackageVersionSendError(t *testing.T) {
	// 用 mock 让 GetPackageInformation 成功但 getRev 失败
	// 实际上 getRev 失败已被 TestUnpublishPackageVersionGetPkgError 覆盖
	// 这里测 sendJSON 失败：用一个会 404 PUT 的 server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"_rev":"10-abc","name":"my-pkg","versions":{"1.0.0":{"name":"my-pkg","version":"1.0.0"}}}`))
			return
		}
		// PUT 返回 500 触发错误
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"server"}`))
	}))
	defer server.Close()
	reg := NewRegistry(NewOptions().SetRegistryURL(server.URL).SetToken("npm_xxx"))
	err := reg.UnpublishPackageVersion(context.Background(), "my-pkg", "1.0.0")
	assert.Error(t, err)
}

func TestDeprecateVersionSendError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"_rev":"10-abc","name":"my-pkg","versions":{"1.0.0":{"name":"my-pkg","version":"1.0.0"}}}`))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"server"}`))
	}))
	defer server.Close()
	reg := NewRegistry(NewOptions().SetRegistryURL(server.URL).SetToken("npm_xxx"))
	err := reg.DeprecateVersion(context.Background(), "my-pkg", "1.0.0", "use v2")
	assert.Error(t, err)
}

func TestStarPackageSendError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			// 返回带 _rev 的包（getRev 走 GET 也命中此分支）
			w.Write([]byte(`{"_rev":"10-abc","name":"my-pkg","versions":{"1.0.0":{"name":"my-pkg","version":"1.0.0"}},"users":{"alice":true}}`))
			return
		}
		// PUT 失败
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"server"}`))
	}))
	defer server.Close()
	reg := NewRegistry(NewOptions().SetRegistryURL(server.URL).SetToken("npm_xxx"))
	err := reg.StarPackage(context.Background(), "my-pkg")
	assert.Error(t, err)
	err = reg.UnstarPackage(context.Background(), "my-pkg")
	assert.Error(t, err)
}

func TestPublishPackageSendError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"_rev":"10-abc","name":"my-pkg","versions":{"1.0.0":{"name":"my-pkg","version":"1.0.0"}}}`))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"server"}`))
	}))
	defer server.Close()
	reg := NewRegistry(NewOptions().SetRegistryURL(server.URL).SetToken("npm_xxx"))
	err := reg.PublishPackage(context.Background(), &models.Package{Name: "my-pkg"})
	assert.Error(t, err)
}

func TestUnpublishPackageSendError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"_rev":"10-abc","name":"my-pkg","versions":{"1.0.0":{"name":"my-pkg","version":"1.0.0"}}}`))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"server"}`))
	}))
	defer server.Close()
	reg := NewRegistry(NewOptions().SetRegistryURL(server.URL).SetToken("npm_xxx"))
	err := reg.UnpublishPackage(context.Background(), "my-pkg")
	assert.Error(t, err)
}

// ====================================================================
// DownloadTarball - 补充 http.Get 失败分支
// ====================================================================

func TestDownloadTarballHTTPGetError(t *testing.T) {
	// tarball URL 指向不可达地址
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"name":"x","version":"1.0.0","dist":{"tarball":"http://localhost:1/x.tgz"}}`))
	}))
	defer server.Close()
	reg := NewRegistry(NewOptions().SetRegistryURL(server.URL))
	err := reg.DownloadTarball(context.Background(), "x", "1.0.0", "/tmp/test-download.tgz")
	assert.Error(t, err)
}
