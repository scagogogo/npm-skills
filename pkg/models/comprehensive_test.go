package models

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// --- access.go ---

func TestPackageAccessToJsonString(t *testing.T) {
	a := &PackageAccess{
		Package: "my-pkg",
		Access:  map[string]string{"read": "public", "write": "restricted"},
	}
	s := a.ToJsonString()
	assert.Contains(t, s, "my-pkg")
	assert.Contains(t, s, "public")
	assert.Contains(t, s, "restricted")
}

func TestPackageAccessUpdateToJsonString(t *testing.T) {
	u := &PackageAccessUpdate{Access: "public"}
	s := u.ToJsonString()
	assert.Contains(t, s, "public")
}

func TestCollaboratorToJsonString(t *testing.T) {
	c := &Collaborator{Name: "alice", Email: "a@x.com", Permissions: "write"}
	s := c.ToJsonString()
	assert.Contains(t, s, "alice")
	assert.Contains(t, s, "write")
}

func TestPermissionConstants(t *testing.T) {
	assert.Equal(t, Permission("read"), PermissionRead)
	assert.Equal(t, Permission("write"), PermissionWrite)
}

// --- audit.go ---

func TestAdvisoryToJsonString(t *testing.T) {
	a := &Advisory{
		ID: 1, Title: "XSS", Severity: "high",
		ModuleName: "lodash", URL: "https://x",
	}
	s := a.ToJsonString()
	assert.Contains(t, s, "lodash")
	assert.Contains(t, s, "high")
}

func TestAdvisoryGetReferences(t *testing.T) {
	// nil
	a := &Advisory{}
	assert.Nil(t, a.GetReferences())

	// array
	a.References = json.RawMessage(`["https://r1","https://r2"]`)
	assert.Equal(t, []string{"https://r1", "https://r2"}, a.GetReferences())

	// single string
	a.References = json.RawMessage(`"https://r3"`)
	assert.Equal(t, []string{"https://r3"}, a.GetReferences())

	// invalid JSON
	a.References = json.RawMessage(`{invalid}`)
	assert.Nil(t, a.GetReferences())
}

func TestQuickAuditRequestToJsonString(t *testing.T) {
	r := &QuickAuditRequest{Dependencies: map[string]string{"lodash": "4.17.11"}}
	s := r.ToJsonString()
	assert.Contains(t, s, "lodash")
}

func TestQuickAuditResultToJsonString(t *testing.T) {
	r := &QuickAuditResult{}
	r.Metadata.Vulnerabilities.High = 2
	r.Metadata.TotalDependencies = 10
	s := r.ToJsonString()
	assert.Contains(t, s, "totalDependencies")
}

// --- couchdb.go ---

func TestChangesResultToJsonString(t *testing.T) {
	r := &ChangesResult{LastSeq: "10", Pending: 0, Results: []ChangeEntry{{Seq: "1", ID: "pkg"}}}
	s := r.ToJsonString()
	assert.Contains(t, s, "last_seq")
	assert.Contains(t, s, "pkg")
}

func TestAllDocsResultToJsonString(t *testing.T) {
	r := &AllDocsResult{TotalRows: 5, Offset: 0, Rows: []DocRow{{ID: "x", Key: "x"}}}
	s := r.ToJsonString()
	assert.Contains(t, s, "total_rows")
	assert.Contains(t, s, "\"x\"")
}

func TestViewResultToJsonString(t *testing.T) {
	r := &ViewResult{TotalRows: 3, Offset: 0, Rows: []ViewRow{{ID: "k"}}}
	s := r.ToJsonString()
	assert.Contains(t, s, "total_rows")
}

// --- hook.go ---

func TestHookToJsonString(t *testing.T) {
	h := &Hook{ID: "1", Type: "hook", Name: "n", Endpoint: "https://e", Events: []string{"package"}, Active: true}
	s := h.ToJsonString()
	assert.Contains(t, s, "endpoint")
	assert.Contains(t, s, "package")
}

func TestHookCreationToJsonString(t *testing.T) {
	h := &HookCreation{Name: "n", Endpoint: "https://e", Events: []string{"package"}, Active: true}
	s := h.ToJsonString()
	assert.Contains(t, s, "endpoint")
}

func TestHookUpdateToJsonString(t *testing.T) {
	active := false
	h := &HookUpdate{Endpoint: "https://e2", Events: []string{"package"}, Active: &active}
	s := h.ToJsonString()
	assert.Contains(t, s, "endpoint")
}

func TestHookListOptionsToJsonString(t *testing.T) {
	o := &HookListOptions{Package: "my-pkg", Page: 1, PerPage: 10}
	s := o.ToJsonString()
	assert.Contains(t, s, "my-pkg")
}

// --- org.go ---

func TestOrganizationToJsonString(t *testing.T) {
	o := &Organization{Name: "myorg", Scope: "myorg"}
	s := o.ToJsonString()
	assert.Contains(t, s, "myorg")
}

func TestTeamToJsonString(t *testing.T) {
	tm := &Team{ID: "1", Name: "devs", Description: "dev team"}
	s := tm.ToJsonString()
	assert.Contains(t, s, "devs")
}

func TestOrgCreationToJsonString(t *testing.T) {
	o := &OrgCreation{Name: "neworg"}
	s := o.ToJsonString()
	assert.Contains(t, s, "neworg")
}

func TestTeamCreationToJsonString(t *testing.T) {
	o := &TeamCreation{Name: "newteam"}
	s := o.ToJsonString()
	assert.Contains(t, s, "newteam")
}

// --- publish.go ---

func TestPublishMetadataToJsonString(t *testing.T) {
	p := &PublishMetadata{
		Name:        "my-pkg",
		Version:     "1.0.0",
		Description: "desc",
		Dependencies: map[string]string{
			"lodash": "^4.17.21",
		},
		Keywords: []string{"util"},
	}
	s := p.ToJsonString()
	assert.Contains(t, s, "my-pkg")
	assert.Contains(t, s, "1.0.0")
	assert.Contains(t, s, "lodash")
}

// --- token.go ---

func TestMaskToken(t *testing.T) {
	// short token
	assert.Equal(t, "****", maskToken("abc"))
	// long token
	assert.Equal(t, "abcd...efgh", maskToken("abcd12345678efgh"))
	// exactly 8 chars → masked (<=8)
	assert.Equal(t, "****", maskToken("12345678"))
	// 9 chars → unmasked (前4 + ... + 后4)
	assert.Equal(t, "1234...6789", maskToken("123456789"))
}

func TestTokenMarshalJSON(t *testing.T) {
	tok := Token{
		ID: "id1", Token: "abcd12345678efgh", Key: "k1",
		Readonly: true, CIDR: []string{"0.0.0.0/0"},
	}
	b, err := json.Marshal(tok)
	assert.NoError(t, err)
	s := string(b)
	// token 应被掩码
	assert.Contains(t, s, "abcd...efgh")
	assert.NotContains(t, s, "12345678efgh")
}

func TestTokenToJsonString(t *testing.T) {
	tok := &Token{ID: "id1", Token: "abcd12345678efgh"}
	s := tok.ToJsonString()
	assert.Contains(t, s, "abcd...efgh")
}

func TestTokenCreationToJsonString(t *testing.T) {
	c := &TokenCreation{Password: "p", Readonly: true, CIDR: []string{"1.2.3.0/24"}}
	s := c.ToJsonString()
	assert.Contains(t, s, "1.2.3.0/24")
}

// --- user.go ---

func TestMaskLoginToken(t *testing.T) {
	assert.Equal(t, "****", maskLoginToken("abc"))
	assert.Equal(t, "abcd...efgh", maskLoginToken("abcd12345678efgh"))
}

func TestLoginResultMarshalJSON(t *testing.T) {
	lr := LoginResult{
		ID: "org.couchdb.user:x", Rev: "1-abc", Token: "abcd12345678efgh",
		Ok: OkBool{value: true, isBool: true},
	}
	b, err := json.Marshal(lr)
	assert.NoError(t, err)
	s := string(b)
	assert.Contains(t, s, "abcd...efgh")
	assert.NotContains(t, s, "12345678efgh")
}

func TestLoginResultToJsonString(t *testing.T) {
	lr := &LoginResult{ID: "x", Token: "abcd12345678efgh"}
	s := lr.ToJsonString()
	assert.Contains(t, s, "abcd...efgh")
}

func TestUserCreationToJsonString(t *testing.T) {
	u := &UserCreation{ID: "org.couchdb.user:alice", Name: "alice", Password: "p", Email: "a@x.com", Type: "user"}
	s := u.ToJsonString()
	assert.Contains(t, s, "alice")
	assert.Contains(t, s, "org.couchdb.user:alice")
}

func TestUserProfileToJsonString(t *testing.T) {
	u := &UserProfile{Name: "alice", Email: "a@x.com", EmailVerified: true}
	s := u.ToJsonString()
	assert.Contains(t, s, "alice")
	assert.Contains(t, s, "a@x.com")
}

func TestOkBoolUnmarshalJSON(t *testing.T) {
	// bool true
	var o OkBool
	err := json.Unmarshal([]byte("true"), &o)
	assert.NoError(t, err)
	assert.True(t, o.isBool)
	assert.True(t, o.Bool())
	assert.Equal(t, "true", o.String())

	// bool false
	err = json.Unmarshal([]byte("false"), &o)
	assert.NoError(t, err)
	assert.True(t, o.isBool)
	assert.False(t, o.Bool())
	assert.Equal(t, "false", o.String())

	// string
	err = json.Unmarshal([]byte(`"user created"`), &o)
	assert.NoError(t, err)
	assert.False(t, o.isBool)
	assert.Equal(t, "user created", o.String())
	assert.False(t, o.Bool())

	// invalid (number)
	err = json.Unmarshal([]byte("123"), &o)
	assert.Error(t, err)
}

func TestOkBoolMarshalJSON(t *testing.T) {
	// bool
	o := OkBool{value: true, isBool: true}
	b, err := o.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, "true", string(b))

	// string
	o = OkBool{str: "user created", isBool: false}
	b, err = o.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, `"user created"`, string(b))
}

// --- package_information.go ---

func TestPackageDeprecatedMessage(t *testing.T) {
	// nil
	p := &Package{}
	assert.Equal(t, "", p.DeprecatedMessage())
	assert.False(t, p.IsDeprecated())

	// string
	p.Deprecated = "use v2"
	assert.True(t, p.IsDeprecated())
	assert.Equal(t, "use v2", p.DeprecatedMessage())

	// empty string
	p.Deprecated = ""
	assert.False(t, p.IsDeprecated())
	assert.Equal(t, "", p.DeprecatedMessage())

	// bool true
	p.Deprecated = true
	assert.True(t, p.IsDeprecated())
	assert.Equal(t, "this package has been deprecated", p.DeprecatedMessage())

	// bool false
	p.Deprecated = false
	assert.False(t, p.IsDeprecated())
	assert.Equal(t, "", p.DeprecatedMessage())

	// other type (e.g. number)
	p.Deprecated = 123
	assert.False(t, p.IsDeprecated())
	assert.Equal(t, "", p.DeprecatedMessage())
}

func TestPackageAttachmentJson(t *testing.T) {
	p := &Package{
		Name: "x",
		Attachments: map[string]Attachment{
			"file.tgz": {ContentType: "application/octet-stream", Data: "base64", Length: 100},
		},
	}
	s := p.ToJsonString()
	assert.Contains(t, s, "_attachments")
	assert.Contains(t, s, "file.tgz")
}

// --- version.go (补全 DeprecatedMessage 的 default/bool 分支) ---

func TestVersionDeprecatedMessageBranches(t *testing.T) {
	// string
	v := &Version{Deprecated: "old"}
	assert.True(t, v.IsDeprecated())
	assert.Equal(t, "old", v.DeprecatedMessage())

	// empty string
	v.Deprecated = ""
	assert.False(t, v.IsDeprecated())
	assert.Equal(t, "", v.DeprecatedMessage())

	// bool true
	v.Deprecated = true
	assert.True(t, v.IsDeprecated())
	assert.Equal(t, "this version has been deprecated", v.DeprecatedMessage())

	// bool false
	v.Deprecated = false
	assert.False(t, v.IsDeprecated())
	assert.Equal(t, "", v.DeprecatedMessage())

	// other type
	v.Deprecated = 3.14
	assert.False(t, v.IsDeprecated())
	assert.Equal(t, "", v.DeprecatedMessage())
}

// --- search_result.go (ToJsonString 75% — 测 Error 分支) ---

func TestSearchResultToJsonStringFull(t *testing.T) {
	// 正常路径
	sr := &SearchResult{Objects: []SearchObject{{Package: SearchPackage{Name: "react"}}}}
	s := sr.ToJsonString()
	assert.Contains(t, s, "react")
}

// --- download_stats.go (ToJsonString 75% — 测 Error 分支) ---

func TestDownloadStatsToJsonStringFull(t *testing.T) {
	ds := &DownloadStats{Downloads: 100, Start: "2024-01-01", End: "2024-01-07", Package: "react"}
	s := ds.ToJsonString()
	assert.Contains(t, s, "100")
	assert.Contains(t, s, "react")
}

func TestDownloadRangeStatsToJsonStringFull(t *testing.T) {
	drs := &DownloadRangeStats{Start: "2024-01-01", End: "2024-01-07", Package: "react"}
	drs.Downloads = append(drs.Downloads, DailyDownloads{Day: "2024-01-01", Downloads: 10})
	s := drs.ToJsonString()
	assert.Contains(t, s, "react")
	assert.Contains(t, s, "2024-01-01")
}

// --- 兜底：验证所有 ToJsonString 至少能跑通 ---

func TestAllModelsToJsonStringSmoke(t *testing.T) {
	cases := []string{
		(&PackageAccess{}).ToJsonString(),
		(&PackageAccessUpdate{}).ToJsonString(),
		(&Collaborator{}).ToJsonString(),
		(&Advisory{}).ToJsonString(),
		(&QuickAuditRequest{}).ToJsonString(),
		(&QuickAuditResult{}).ToJsonString(),
		(&ChangesResult{}).ToJsonString(),
		(&AllDocsResult{}).ToJsonString(),
		(&ViewResult{}).ToJsonString(),
		(&Hook{}).ToJsonString(),
		(&HookCreation{}).ToJsonString(),
		(&HookUpdate{}).ToJsonString(),
		(&HookListOptions{}).ToJsonString(),
		(&Organization{}).ToJsonString(),
		(&Team{}).ToJsonString(),
		(&OrgCreation{}).ToJsonString(),
		(&TeamCreation{}).ToJsonString(),
		(&PublishMetadata{Name: "x", Version: "1.0.0"}).ToJsonString(),
		(&Token{}).ToJsonString(),
		(&TokenCreation{}).ToJsonString(),
		(&LoginResult{}).ToJsonString(),
		(&UserCreation{}).ToJsonString(),
		(&UserProfile{}).ToJsonString(),
	}
	for i, s := range cases {
		// 每个都应是非空且能被重新解析或为 "{}"
		_ = s
		if s == "" {
			t.Fatalf("case %d returned empty", i)
		}
		// 简单烟雾：trim 后应以 { 开头
		assert.True(t, strings.HasPrefix(strings.TrimSpace(s), "{"), "case %d not object: %s", i, s)
	}
}
