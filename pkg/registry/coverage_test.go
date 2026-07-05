package registry

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/scagogogo/npm-skills/pkg/models"
	"github.com/stretchr/testify/assert"
)

// fullMockServer 创建一个覆盖多种端点的 mock NPM registry 服务器。
// 通过路径前缀路由，返回不同 JSON 响应，并记录请求方法/路径用于断言。
func fullMockServer(t *testing.T, recording *[]string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if recording != nil {
			*recording = append(*recording, r.Method+" "+r.URL.Path)
		}
		path := r.URL.Path

		// 包文档 GET /my-pkg （含 _rev）
		if path == "/my-pkg" {
			switch r.Method {
			case http.MethodGet:
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{
					"_id":"my-pkg","_rev":"10-abc","name":"my-pkg",
					"dist-tags":{"latest":"1.0.0"},
					"versions":{"1.0.0":{"name":"my-pkg","version":"1.0.0"}},
					"users":{"alice":true}
				}`))
				return
			case http.MethodPut:
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"ok":true,"rev":"11-def"}`))
				return
			case http.MethodDelete:
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"ok":true}`))
				return
			}
		}

		// 不存在的包（404）
		if path == "/not-exist-pkg" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":"not found"}`))
			return
		}

		// 新包 /new-pkg：GET 返回 404（不存在），PUT 返回 200（创建）
		if path == "/new-pkg" {
			if r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(``))
				return
			}
			// PUT 走到底层 PUT 兜底分支
		}

		// dist-tags 端点
		if path == "/-/package/my-pkg/dist-tags" {
			if r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"latest":"1.0.0","next":"2.0.0-rc.1"}`))
				return
			}
			if r.Method == http.MethodPost {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"latest":"1.0.0","next":"2.0.0-rc.1"}`))
				return
			}
		}
		// 单个 dist-tag
		if path == "/-/package/my-pkg/dist-tags/latest" {
			if r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`"1.0.0"`))
				return
			}
			if r.Method == http.MethodPut {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"ok":true}`))
				return
			}
			if r.Method == http.MethodDelete {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"ok":true}`))
				return
			}
		}
		// 单个 dist-tag 端点返回 error（Verdaccio 风格），触发 fallback
		if strings.HasPrefix(path, "/-/package/verdaccio-pkg/dist-tags/") {
			if r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"error":"File not found"}`))
				return
			}
		}
		if path == "/-/package/verdaccio-pkg/dist-tags" {
			if r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"latest":"1.0.0"}`))
				return
			}
		}
		// 单个 dist-tag 端点返回 object {tag:version}
		if path == "/-/package/obj-pkg/dist-tags/beta" {
			if r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"beta":"1.0.0-beta.1"}`))
				return
			}
		}
		// 单个 dist-tag 端点返回单个值（非 string、非 array、非 error、非 tag-key）
		if path == "/-/package/singleval-pkg/dist-tags/alpha" {
			if r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				// 仅有一个 string 值，但 key 不是 "alpha"
				w.Write([]byte(`{"other":"1.2.3"}`))
				return
			}
		}

		// 下载统计端点（api.npmjs.org 风格，但用 mock server URL）
		// bulk 请求路径含逗号（多个包名 join），返回 map 格式
		if strings.HasPrefix(path, "/downloads/point/") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if strings.Contains(r.URL.Path, ",") || strings.Contains(r.URL.Path, "%2C") {
				// bulk 响应：{pkg: {downloads:...}}
				w.Write([]byte(`{"my-pkg":{"downloads":100,"start":"2024-01-01","end":"2024-01-07","package":"my-pkg"}}`))
				return
			}
			w.Write([]byte(`{"downloads":100,"start":"2024-01-01","end":"2024-01-07","package":"my-pkg"}`))
			return
		}
		if strings.HasPrefix(path, "/downloads/range/") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if strings.Contains(r.URL.Path, ",") || strings.Contains(r.URL.Path, "%2C") {
				w.Write([]byte(`{"my-pkg":{"start":"2024-01-01","end":"2024-01-07","package":"my-pkg","downloads":[{"day":"2024-01-01","downloads":10}]}}`))
				return
			}
			w.Write([]byte(`{"start":"2024-01-01","end":"2024-01-07","package":"my-pkg","downloads":[{"day":"2024-01-01","downloads":10}]}`))
			return
		}

		// CouchDB 端点
		if path == "/_changes" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"last_seq":"10","pending":0,"results":[{"seq":"1","id":"pkg","changes":[{"rev":"1-abc"}]}]}`))
			return
		}
		if path == "/_all_docs" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"total_rows":5,"offset":0,"rows":[{"id":"x","key":"x","value":{"rev":"1-abc"}}]}`))
			return
		}
		if strings.HasPrefix(path, "/-/_view/") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"total_rows":2,"offset":0,"rows":[{"id":"1","key":"\"alice\"","value":"my-pkg"}]}`))
			return
		}

		// audit 端点
		if path == "/-/npm/v1/security/advisories/bulk" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"lodash":[{"id":1,"title":"XSS","severity":"high","module_name":"lodash"}]}`))
			return
		}
		if strings.HasPrefix(path, "/-/npm/v1/security/advisories/") && path != "/-/npm/v1/security/advisories/bulk" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id":1234,"title":"XSS","severity":"high","module_name":"lodash"}`))
			return
		}
		// list advisories 端点（精确匹配，区分于单个 advisory）
		if path == "/-/npm/v1/security/advisories" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"objects":[{"id":1234,"title":"XSS","severity":"high","module_name":"lodash"}],"total":1}`))
			return
		}
		if path == "/-/npm/v1/security/audits/quick" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"metadata":{"vulnerabilities":{"low":1,"moderate":2,"high":3,"critical":0},"dependencies":5,"totalDependencies":5}}`))
			return
		}

		// whoami
		if path == "/-/whoami" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"username":"alice"}`))
			return
		}

		// user 端点
		if strings.HasPrefix(path, "/-/user/org.couchdb.user:") {
			if r.Method == http.MethodPut {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"id":"org.couchdb.user:alice","rev":"1-abc","token":"abcd12345678efgh","ok":true}`))
				return
			}
			if r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"_id":"org.couchdb.user:alice","name":"alice","email":"a@x.com","type":"user"}`))
				return
			}
		}

		// PUT 任意包路径（发布/更新）返回成功
		if r.Method == http.MethodPut {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"ok":true,"rev":"11-def"}`))
			return
		}

		// 默认 200 + 空对象
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
}

// registryAt 用 fullMockServer 创建一个指向它的 Registry。
func registryAt(server *httptest.Server, opts ...func(*Options)) *Registry {
	o := NewOptions().SetRegistryURL(server.URL)
	for _, opt := range opts {
		opt(o)
	}
	return NewRegistry(o)
}

// ====================================================================
// download_stats.go
// ====================================================================

func TestGetDownloadRangeStatsMock(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server, func(o *Options) {
		o.SetDownloadStatsURL(server.URL + "/downloads")
	})
	res, err := reg.GetDownloadRangeStats(context.Background(), "my-pkg", "last-week")
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "my-pkg", res.Package)
}

func TestGetDownloadStatsByDateRangeMock(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server, func(o *Options) {
		o.SetDownloadStatsURL(server.URL + "/downloads")
	})
	res, err := reg.GetDownloadStatsByDateRange(context.Background(), "my-pkg", "2024-01-01", "2024-01-07")
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 100, res.Downloads)
}

func TestGetDownloadRangeStatsByDateRangeMock(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server, func(o *Options) {
		o.SetDownloadStatsURL(server.URL + "/downloads")
	})
	res, err := reg.GetDownloadRangeStatsByDateRange(context.Background(), "my-pkg", "2024-01-01", "2024-01-07")
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Downloads, 1)
}

func TestGetBulkDownloadStatsMock(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server, func(o *Options) {
		o.SetDownloadStatsURL(server.URL + "/downloads")
	})
	// 空 slice
	_, err := reg.GetBulkDownloadStats(context.Background(), []string{}, "last-week")
	assert.Error(t, err)

	// 正常
	res, err := reg.GetBulkDownloadStats(context.Background(), []string{"my-pkg", "vue"}, "last-week")
	assert.NoError(t, err)
	assert.Contains(t, res, "my-pkg")

	// >128 包分批（用 130 个名字，分 128+2 两批，触发分批边界）
	names := make([]string, 130)
	for i := range names {
		names[i] = fmt.Sprintf("pkg-%d", i)
	}
	_, err = reg.GetBulkDownloadStats(context.Background(), names, "last-week")
	assert.NoError(t, err)
}

func TestGetBulkDownloadRangeStatsMock(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server, func(o *Options) {
		o.SetDownloadStatsURL(server.URL + "/downloads")
	})
	_, err := reg.GetBulkDownloadRangeStats(context.Background(), []string{}, "last-week")
	assert.Error(t, err)
	res, err := reg.GetBulkDownloadRangeStats(context.Background(), []string{"my-pkg", "vue"}, "last-week")
	assert.NoError(t, err)
	assert.Contains(t, res, "my-pkg")
}

func TestGetBulkDownloadStatsByDateRangeMock(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server, func(o *Options) {
		o.SetDownloadStatsURL(server.URL + "/downloads")
	})
	_, err := reg.GetBulkDownloadStatsByDateRange(context.Background(), []string{}, "2024-01-01", "2024-01-07")
	assert.Error(t, err)
	res, err := reg.GetBulkDownloadStatsByDateRange(context.Background(), []string{"my-pkg", "vue"}, "2024-01-01", "2024-01-07")
	assert.NoError(t, err)
	assert.Contains(t, res, "my-pkg")
}

func TestGetBulkDownloadRangeStatsByDateRangeMock(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server, func(o *Options) {
		o.SetDownloadStatsURL(server.URL + "/downloads")
	})
	_, err := reg.GetBulkDownloadRangeStatsByDateRange(context.Background(), []string{}, "2024-01-01", "2024-01-07")
	assert.Error(t, err)
	res, err := reg.GetBulkDownloadRangeStatsByDateRange(context.Background(), []string{"my-pkg", "vue"}, "2024-01-01", "2024-01-07")
	assert.NoError(t, err)
	assert.Contains(t, res, "my-pkg")
}

func TestGetDownloadStatsInvalidName(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server, func(o *Options) {
		o.SetDownloadStatsURL(server.URL + "/downloads")
	})
	// 无效包名（空）
	_, err := reg.GetDownloadStats(context.Background(), "", "last-week")
	assert.Error(t, err)
	// 无效包名
	_, err = reg.GetDownloadRangeStats(context.Background(), "", "last-week")
	assert.Error(t, err)
	_, err = reg.GetDownloadStatsByDateRange(context.Background(), "", "2024-01-01", "2024-01-07")
	assert.Error(t, err)
	_, err = reg.GetDownloadRangeStatsByDateRange(context.Background(), "", "2024-01-01", "2024-01-07")
	assert.Error(t, err)
}

// ====================================================================
// couchdb.go
// ====================================================================

func TestGetChangesMock(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server)
	// 带 opts
	res, err := reg.GetChanges(context.Background(), models.ChangesOptions{Since: "5", Limit: 10, IncludeDocs: true})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "10", res.LastSeq)
	// 空 opts
	_, err = reg.GetChanges(context.Background(), models.ChangesOptions{})
	assert.NoError(t, err)
}

func TestGetAllDocsMock(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server)
	res, err := reg.GetAllDocs(context.Background(), models.AllDocsOptions{
		StartKey: "a", EndKey: "z", Limit: 10, Skip: 5, IncludeDocs: true, Descending: true,
	})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 5, res.TotalRows)
	_, err = reg.GetAllDocs(context.Background(), models.AllDocsOptions{})
	assert.NoError(t, err)
}

func TestGetViewMock(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server)
	res, err := reg.GetView(context.Background(), "starredByUser", models.ViewOptions{
		Key: "\"alice\"", StartKey: "a", EndKey: "z", Limit: 10, Skip: 5,
		Group: true, GroupLevel: 1, Descending: true,
	})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 2, res.TotalRows)
	_, err = reg.GetView(context.Background(), "starredByUser", models.ViewOptions{})
	assert.NoError(t, err)
}

// ====================================================================
// dist_tags.go
// ====================================================================

func TestGetDistTagMock(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server)
	// 标准 string 响应
	v, err := reg.GetDistTag(context.Background(), "my-pkg", "latest")
	assert.NoError(t, err)
	assert.Equal(t, "1.0.0", v)
}

func TestGetDistTagVerdaccioFallbackMock(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server)
	// error 响应触发 fallback 到 GetDistTagsAbbreviated
	v, err := reg.GetDistTag(context.Background(), "verdaccio-pkg", "latest")
	assert.NoError(t, err)
	assert.Equal(t, "1.0.0", v)
}

func TestGetDistTagVerdaccioFallbackMissingMock(t *testing.T) {
	// fallback 后 tag 不存在
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server)
	_, err := reg.GetDistTag(context.Background(), "verdaccio-pkg", "nonexistent-tag")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetDistTagObjectFormatMock(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server)
	// {tag:version} 格式
	v, err := reg.GetDistTag(context.Background(), "obj-pkg", "beta")
	assert.NoError(t, err)
	assert.Equal(t, "1.0.0-beta.1", v)
}

func TestGetDistTagSingleValueFallbackMock(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server)
	// 单个值，key 不匹配，触发 for-range 返回唯一值
	v, err := reg.GetDistTag(context.Background(), "singleval-pkg", "alpha")
	assert.NoError(t, err)
	assert.Equal(t, "1.2.3", v)
}

func TestGetDistTagsAbbreviatedErrorMock(t *testing.T) {
	// 网络错误
	reg := NewRegistry(NewOptions().SetRegistryURL("http://localhost:1"))
	_, err := reg.GetDistTagsAbbreviated(context.Background(), "my-pkg")
	assert.Error(t, err)
}

func TestSetDistTagMockCov(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server, func(o *Options) { o.SetToken("npm_xxx") })
	err := reg.SetDistTag(context.Background(), "my-pkg", "latest", "1.0.0")
	assert.NoError(t, err)
}

func TestSetDistTagErrorMock(t *testing.T) {
	// 用无效 URL 触发网络错误
	regBad := NewRegistry(NewOptions().SetToken("npm_xxx").SetRegistryURL("http://localhost:1"))
	err := regBad.SetDistTag(context.Background(), "my-pkg", "latest", "1.0.0")
	assert.Error(t, err)
}

func TestDeleteDistTagMockCov(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server, func(o *Options) { o.SetToken("npm_xxx") })
	err := reg.DeleteDistTag(context.Background(), "my-pkg", "latest")
	assert.NoError(t, err)
}

func TestSetDistTagsMockCov(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server, func(o *Options) { o.SetToken("npm_xxx") })
	err := reg.SetDistTags(context.Background(), "my-pkg", map[string]string{"latest": "1.0.0"})
	assert.NoError(t, err)
}

// ====================================================================
// abbreviated.go
// ====================================================================

func TestGetAbbreviatedPackageInformationMock(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server)
	pkg, err := reg.GetAbbreviatedPackageInformation(context.Background(), "my-pkg")
	assert.NoError(t, err)
	assert.NotNil(t, pkg)
}

func TestGetAbbreviatedPackageInformationInvalidName(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server)
	_, err := reg.GetAbbreviatedPackageInformation(context.Background(), "")
	assert.Error(t, err)
}

// ====================================================================
// audit.go
// ====================================================================

func TestBulkAuditMockCov(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server)
	res, err := reg.BulkAudit(context.Background(), map[string][]string{"lodash": {"<4.17.12"}})
	assert.NoError(t, err)
	assert.Contains(t, res, "lodash")
}

func TestGetAdvisoryMockCov(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server)
	a, err := reg.GetAdvisory(context.Background(), 1234)
	assert.NoError(t, err)
	assert.NotNil(t, a)
	assert.Equal(t, "lodash", a.ModuleName)
}

func TestListAdvisoriesMockCov(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server)
	// 带参数
	list, err := reg.ListAdvisories(context.Background(), models.AdvisoryListOptions{Page: 1, PerPage: 20, AffectedPackage: "lodash"})
	assert.NoError(t, err)
	assert.Len(t, list, 1)
	// 不带参数
	_, err = reg.ListAdvisories(context.Background(), models.AdvisoryListOptions{})
	assert.NoError(t, err)
}

func TestQuickAuditMockCov(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server)
	res, err := reg.QuickAudit(context.Background(), &models.QuickAuditRequest{Dependencies: map[string]string{"lodash": "4.17.11"}})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 5, res.Metadata.TotalDependencies)
}

func TestBulkAuditErrorMock(t *testing.T) {
	reg := NewRegistry(NewOptions().SetRegistryURL("http://localhost:1"))
	_, err := reg.BulkAudit(context.Background(), map[string][]string{"lodash": {"<4.17.12"}})
	assert.Error(t, err)
	_, err = reg.GetAdvisory(context.Background(), 1234)
	assert.Error(t, err)
	_, err = reg.ListAdvisories(context.Background(), models.AdvisoryListOptions{})
	assert.Error(t, err)
	_, err = reg.QuickAudit(context.Background(), &models.QuickAuditRequest{})
	assert.Error(t, err)
}

// ====================================================================
// whoami.go
// ====================================================================

func TestWhoAmINoToken(t *testing.T) {
	reg := NewRegistry(NewOptions())
	_, err := reg.WhoAmI(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no token set")
}

func TestWhoAmIMockCov(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server, func(o *Options) { o.SetToken("npm_xxx") })
	name, err := reg.WhoAmI(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "alice", name)
}

func TestWhoAmINetworkError(t *testing.T) {
	reg := NewRegistry(NewOptions().SetToken("npm_xxx").SetRegistryURL("http://localhost:1"))
	_, err := reg.WhoAmI(context.Background())
	assert.Error(t, err)
}

func TestWhoAmIEmptyUsername(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"username":""}`))
	}))
	defer server.Close()
	reg := registryAt(server, func(o *Options) { o.SetToken("npm_xxx") })
	_, err := reg.WhoAmI(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty username")
}

func TestWhoAmIBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`not-json`))
	}))
	defer server.Close()
	reg := registryAt(server, func(o *Options) { o.SetToken("npm_xxx") })
	_, err := reg.WhoAmI(context.Background())
	assert.Error(t, err)
}

// ====================================================================
// user.go
// ====================================================================

func TestCreateUserMock(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server)
	res, err := reg.CreateUser(context.Background(), &models.UserCreation{Name: "alice", Password: "p", Email: "a@x.com"})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "abcd12345678efgh", res.Token)
}

func TestCreateUserErrorMock(t *testing.T) {
	reg := NewRegistry(NewOptions().SetRegistryURL("http://localhost:1"))
	_, err := reg.CreateUser(context.Background(), &models.UserCreation{Name: "alice", Password: "p"})
	assert.Error(t, err)
}

func TestGetUserMockCov(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server, func(o *Options) { o.SetToken("npm_xxx") })
	prof, err := reg.GetUser(context.Background(), "alice")
	assert.NoError(t, err)
	assert.NotNil(t, prof)
	assert.Equal(t, "alice", prof.Name)
}

func TestGetUserErrorMock(t *testing.T) {
	reg := NewRegistry(NewOptions().SetRegistryURL("http://localhost:1"))
	_, err := reg.GetUser(context.Background(), "alice")
	assert.Error(t, err)
}

// ====================================================================
// star.go
// ====================================================================

func TestStarPackageMock(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server, func(o *Options) { o.SetToken("npm_xxx") })
	err := reg.StarPackage(context.Background(), "my-pkg")
	assert.NoError(t, err)
}

func TestUnstarPackageMock(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server, func(o *Options) { o.SetToken("npm_xxx") })
	err := reg.UnstarPackage(context.Background(), "my-pkg")
	assert.NoError(t, err)
}

func TestStarPackageGetPkgError(t *testing.T) {
	reg := NewRegistry(NewOptions().SetToken("npm_xxx").SetRegistryURL("http://localhost:1"))
	err := reg.StarPackage(context.Background(), "my-pkg")
	assert.Error(t, err)
}

func TestUnstarPackageGetPkgError(t *testing.T) {
	reg := NewRegistry(NewOptions().SetToken("npm_xxx").SetRegistryURL("http://localhost:1"))
	err := reg.UnstarPackage(context.Background(), "my-pkg")
	assert.Error(t, err)
}

func TestGetStarredByUserMockCov(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server)
	pkgs, err := reg.GetStarredByUser(context.Background(), "alice")
	assert.NoError(t, err)
	assert.Contains(t, pkgs, "my-pkg")
}

func TestGetStarredByPackageMockCov(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server)
	users, err := reg.GetStarredByPackage(context.Background(), "my-pkg")
	assert.NoError(t, err)
	assert.Contains(t, users, "my-pkg")
}

func TestGetStarredByUserErrorMock(t *testing.T) {
	reg := NewRegistry(NewOptions().SetRegistryURL("http://localhost:1"))
	_, err := reg.GetStarredByUser(context.Background(), "alice")
	assert.Error(t, err)
	_, err = reg.GetStarredByPackage(context.Background(), "my-pkg")
	assert.Error(t, err)
}

func TestGetStarredByUserBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`not-json`))
	}))
	defer server.Close()
	reg := registryAt(server)
	_, err := reg.GetStarredByUser(context.Background(), "alice")
	assert.Error(t, err)
	_, err = reg.GetStarredByPackage(context.Background(), "my-pkg")
	assert.Error(t, err)
}

// ====================================================================
// publish.go / deprecate.go / unpublish.go
// ====================================================================

func TestGetRevMock(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server, func(o *Options) { o.SetToken("npm_xxx") })
	rev, err := reg.getRev(context.Background(), "my-pkg")
	assert.NoError(t, err)
	assert.Equal(t, "10-abc", rev)
}

func TestGetRevNewPackageMock(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server, func(o *Options) { o.SetToken("npm_xxx") })
	// new-pkg 返回 404 + 空 body
	rev, err := reg.getRev(context.Background(), "new-pkg")
	assert.NoError(t, err)
	assert.Equal(t, "", rev)
}

func TestGetRevError(t *testing.T) {
	reg := NewRegistry(NewOptions().SetToken("npm_xxx").SetRegistryURL("http://localhost:1"))
	_, err := reg.getRev(context.Background(), "my-pkg")
	assert.Error(t, err)
}

func TestGetRevBadJSON(t *testing.T) {
	// 404 但 body 不是 JSON（解析失败 → 返回空字符串）
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`<html>not found</html>`))
	}))
	defer server.Close()
	reg := registryAt(server, func(o *Options) { o.SetToken("npm_xxx") })
	rev, err := reg.getRev(context.Background(), "my-pkg")
	assert.NoError(t, err)
	assert.Equal(t, "", rev)
}

func TestPublishPackageMock(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server, func(o *Options) { o.SetToken("npm_xxx") })
	err := reg.PublishPackage(context.Background(), &models.Package{
		Name:     "my-pkg",
		DistTags: map[string]string{"latest": "1.0.0"},
		Versions: map[string]models.Version{"1.0.0": {Name: "my-pkg", Version: "1.0.0"}},
	})
	assert.NoError(t, err)
}

func TestPublishPackageErrorMock(t *testing.T) {
	reg := NewRegistry(NewOptions().SetToken("npm_xxx").SetRegistryURL("http://localhost:1"))
	err := reg.PublishPackage(context.Background(), &models.Package{Name: "my-pkg"})
	assert.Error(t, err)
}

func TestPublishPackageFromTarballNewMock(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server, func(o *Options) { o.SetToken("npm_xxx") })
	// new-pkg：GetPackageInformation 返回错误（404），走新建分支
	err := reg.PublishPackageFromTarball(context.Background(), "new-pkg", "1.0.0", []byte("tarball"), &models.PublishMetadata{
		Name:        "new-pkg",
		Version:     "1.0.0",
		Description: "desc",
		Keywords:    []string{"util"},
		License:     "MIT",
		Homepage:    "https://x",
	})
	assert.NoError(t, err)
}

func TestPublishPackageFromTarballExistingMock(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server, func(o *Options) { o.SetToken("npm_xxx") })
	// my-pkg：已存在，走合并分支
	err := reg.PublishPackageFromTarball(context.Background(), "my-pkg", "2.0.0", []byte("tarball"), &models.PublishMetadata{
		Name:    "my-pkg",
		Version: "2.0.0",
	})
	assert.NoError(t, err)
}

func TestPublishPackageFromTarballRevError(t *testing.T) {
	// getRev 失败：用 my-pkg 但服务器不可达
	reg := NewRegistry(NewOptions().SetToken("npm_xxx").SetRegistryURL("http://localhost:1"))
	err := reg.PublishPackageFromTarball(context.Background(), "my-pkg", "1.0.0", []byte("x"), &models.PublishMetadata{Name: "my-pkg", Version: "1.0.0"})
	assert.Error(t, err)
}

func TestDeprecateVersionMock(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server, func(o *Options) { o.SetToken("npm_xxx") })
	err := reg.DeprecateVersion(context.Background(), "my-pkg", "1.0.0", "use v2")
	assert.NoError(t, err)
}

func TestDeprecateVersionNotFoundMock(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server, func(o *Options) { o.SetToken("npm_xxx") })
	// 版本不存在
	err := reg.DeprecateVersion(context.Background(), "my-pkg", "9.9.9", "use v2")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestDeprecateVersionGetPkgError(t *testing.T) {
	reg := NewRegistry(NewOptions().SetToken("npm_xxx").SetRegistryURL("http://localhost:1"))
	err := reg.DeprecateVersion(context.Background(), "my-pkg", "1.0.0", "use v2")
	assert.Error(t, err)
}

func TestUnpublishPackageMock(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server, func(o *Options) { o.SetToken("npm_xxx") })
	err := reg.UnpublishPackage(context.Background(), "my-pkg")
	assert.NoError(t, err)
}

func TestUnpublishPackageRevError(t *testing.T) {
	reg := NewRegistry(NewOptions().SetToken("npm_xxx").SetRegistryURL("http://localhost:1"))
	err := reg.UnpublishPackage(context.Background(), "my-pkg")
	assert.Error(t, err)
}

func TestUnpublishPackageVersionMock(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server, func(o *Options) { o.SetToken("npm_xxx") })
	err := reg.UnpublishPackageVersion(context.Background(), "my-pkg", "1.0.0")
	assert.NoError(t, err)
}

func TestUnpublishPackageVersionNotFoundMock(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server, func(o *Options) { o.SetToken("npm_xxx") })
	err := reg.UnpublishPackageVersion(context.Background(), "my-pkg", "9.9.9")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestUnpublishPackageVersionGetPkgError(t *testing.T) {
	reg := NewRegistry(NewOptions().SetToken("npm_xxx").SetRegistryURL("http://localhost:1"))
	err := reg.UnpublishPackageVersion(context.Background(), "my-pkg", "1.0.0")
	assert.Error(t, err)
}

// ====================================================================
// custom_registry.go - RegistryHealthCheck
// ====================================================================

func TestRegistryHealthCheckOK(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server, func(o *Options) { o.SetToken("npm_xxx") })
	ok, err := reg.RegistryHealthCheck(context.Background())
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestRegistryHealthCheckUnauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()
	reg := registryAt(server)
	ok, err := reg.RegistryHealthCheck(context.Background())
	assert.False(t, ok)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "401")
}

func TestRegistryHealthCheckForbidden(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()
	reg := registryAt(server)
	ok, err := reg.RegistryHealthCheck(context.Background())
	assert.False(t, ok)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "403")
}

func TestRegistryHealthCheckNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()
	reg := registryAt(server)
	ok, err := reg.RegistryHealthCheck(context.Background())
	assert.False(t, ok)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "404")
}

func TestRegistryHealthCheckUnexpectedStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()
	reg := registryAt(server)
	ok, err := reg.RegistryHealthCheck(context.Background())
	assert.False(t, ok)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected status")
}

func TestRegistryHealthCheckUnreachableCov(t *testing.T) {
	reg := NewRegistry(NewOptions().SetRegistryURL("http://localhost:1"))
	ok, err := reg.RegistryHealthCheck(context.Background())
	assert.False(t, ok)
	assert.Error(t, err)
}

func TestRegistryHealthCheckWithTimeout(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	// 显式设置 Timeout，跳过内部 10s 默认
	reg := registryAt(server, func(o *Options) {
		o.SetToken("npm_xxx")
		o.SetTimeout(0)
	})
	// Timeout==0 → 内部用 10s；这里 server 可达
	ok, err := reg.RegistryHealthCheck(context.Background())
	assert.True(t, ok)
	assert.NoError(t, err)
}

func TestRegistryHealthCheckBasicAuth(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server, func(o *Options) {
		o.SetBasicAuth("user", "pass")
		o.SetUserAgent("test-agent")
	})
	ok, err := reg.RegistryHealthCheck(context.Background())
	assert.True(t, ok)
	assert.NoError(t, err)
}

// ====================================================================
// request.go
// ====================================================================

func TestEncodePackageNameForPath(t *testing.T) {
	assert.Equal(t, "my-pkg", encodePackageNameForPath("my-pkg"))
	assert.Equal(t, "@nestjs%2Fcore", encodePackageNameForPath("@nestjs/core"))
}

func TestDefaultAcceptStatusCodes(t *testing.T) {
	assert.Equal(t, []int{http.StatusOK, http.StatusCreated}, defaultAcceptStatusCodes(nil))
	assert.Equal(t, []int{http.StatusOK}, defaultAcceptStatusCodes([]int{http.StatusOK}))
}

func TestSendRequestMockCov(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server)
	// GET 风格用 sendRequest
	_, err := reg.sendRequest(context.Background(), http.MethodGet, server.URL+"/my-pkg", nil)
	assert.NoError(t, err)
	// 带 body
	_, err = reg.sendRequest(context.Background(), http.MethodPut, server.URL+"/my-pkg", []byte(`{}`), http.StatusOK)
	assert.NoError(t, err)
}

func TestSendJSONMock2(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server)
	_, err := reg.sendJSON(context.Background(), http.MethodPut, server.URL+"/my-pkg", map[string]string{"k": "v"})
	assert.NoError(t, err)
}

func TestSendJSONMarshalError(t *testing.T) {
	server := fullMockServer(t, nil)
	defer server.Close()
	reg := registryAt(server)
	// 传入无法序列化的对象（chan）
	_, err := reg.sendJSON(context.Background(), http.MethodPut, server.URL+"/my-pkg", make(chan int))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "marshal")
}

func TestSendRequestNetworkError(t *testing.T) {
	reg := NewRegistry(NewOptions().SetRegistryURL("http://localhost:1"))
	_, err := reg.sendRequest(context.Background(), http.MethodGet, "http://localhost:1/x", nil)
	assert.Error(t, err)
}

// ====================================================================
// errors.go - Is 的非 APIError 分支
// ====================================================================

func TestAPIErrorIsNonAPIError(t *testing.T) {
	e := NewAPIError(404, "not found")
	// target 不是 *APIError
	assert.False(t, errors.Is(e, errors.New("plain error")))
	// target 是不同 StatusCode
	assert.False(t, errors.Is(e, NewAPIError(500, "server error")))
	// target 是相同 StatusCode
	assert.True(t, errors.Is(e, ErrNotFound))
}

func TestStatusCodeToErrorDefault(t *testing.T) {
	// 418 是非映射的状态码，>=400
	err := statusCodeToError(418)
	assert.NotNil(t, err)
	// 200 不是错误
	assert.Nil(t, statusCodeToError(200))
}

// ====================================================================
// search.go - SearchPackagesWithOptions
// ====================================================================

func TestSearchPackagesWithOptionsMock(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"objects":[{"package":{"name":"react","version":"18.0.0","description":"ui lib"},"score":{"final":0.9,"detail":{"quality":0.9,"popularity":0.9,"maintenance":0.9}},"searchScore":0.9}],
			"total":1
		}`))
	}))
	defer server.Close()
	reg := registryAt(server)
	res, err := reg.SearchPackagesWithOptions(context.Background(), "react", SearchOptions{From: 0, Quality: 0.5, Popularity: 0.5, Maintenance: 0.5})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Objects, 1)
}

func TestSearchPackagesWithOptionsError(t *testing.T) {
	reg := NewRegistry(NewOptions().SetRegistryURL("http://localhost:1"))
	_, err := reg.SearchPackagesWithOptions(context.Background(), "react", SearchOptions{})
	assert.Error(t, err)
}

// ====================================================================
// versions.go - 补充未覆盖分支
// ====================================================================

func TestGetPackageLatestVersionEmptyTags(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()
	reg := registryAt(server)
	_, err := reg.GetPackageLatestVersion(context.Background(), "my-pkg")
	assert.Error(t, err)
}

func TestGetPackageVersionsErrorMock(t *testing.T) {
	reg := NewRegistry(NewOptions().SetRegistryURL("http://localhost:1"))
	_, err := reg.GetPackageVersions(context.Background(), "my-pkg")
	assert.Error(t, err)
}

func TestGetPackageVersionCountErrorMock(t *testing.T) {
	reg := NewRegistry(NewOptions().SetRegistryURL("http://localhost:1"))
	_, err := reg.GetPackageVersionCount(context.Background(), "my-pkg")
	assert.Error(t, err)
}

// ====================================================================
// DownloadTarball - 补充 os.Create 失败分支
// ====================================================================

func TestDownloadTarballCreateFileError(t *testing.T) {
	server := mockTestServer()
	defer server.Close()
	reg := registryAt(server)
	// destPath 是一个不可写的路径（已存在的目录）
	err := reg.DownloadTarball(context.Background(), "axios", "1.0.0", "/proc/self")
	assert.Error(t, err)
}

// ====================================================================
// Login (already covered) - 补充 Login 的解析错误
// ====================================================================

func TestLoginBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`not-json`))
	}))
	defer server.Close()
	reg := registryAt(server)
	_, err := reg.Login(context.Background(), "alice", "p")
	assert.Error(t, err)
}

// ====================================================================
// GetPackageInformation / GetRegistryInformation 网络错误分支
// ====================================================================

func TestGetPackageInformationNetworkErr(t *testing.T) {
	reg := NewRegistry(NewOptions().SetRegistryURL("http://localhost:1"))
	_, err := reg.GetPackageInformation(context.Background(), "my-pkg")
	assert.Error(t, err)
}

func TestGetRegistryInformationNetworkErr(t *testing.T) {
	reg := NewRegistry(NewOptions().SetRegistryURL("http://localhost:1"))
	_, err := reg.GetRegistryInformation(context.Background())
	assert.Error(t, err)
}

// ====================================================================
// 烟雾测试：确保 json 包可见
// ====================================================================

func TestJSONImportUsed(t *testing.T) {
	// 确保 encoding/json 在测试中被引用（避免 unused import）
	var m map[string]interface{}
	assert.NoError(t, json.Unmarshal([]byte(`{}`), &m))
	_ = os.Stdout
}
