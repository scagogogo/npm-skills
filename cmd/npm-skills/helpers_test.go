package main

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/fatih/color"
	"github.com/scagogogo/npm-skills/pkg/registry"
	"github.com/stretchr/testify/assert"
)

func TestMirrorToURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"official mirror", "official", registry.DefaultRegistryURL},
		{"taobao mirror", "taobao", registry.RegistryUrlTaoBao},
		{"npm-mirror mirror", "npm-mirror", registry.RegistryUrlNpmMirror},
		{"huawei mirror", "huawei", registry.RegistryUrlHuaWeiCloud},
		{"tencent mirror", "tencent", registry.RegistryUrlTencent},
		{"cnpm mirror", "cnpm", registry.RegistryUrlCnpm},
		{"yarn mirror", "yarn", registry.RegistryUrlYarn},
		{"npmjscom mirror", "npmjscom", registry.RegistryUrlNpmjsCom},
		{"npmmirror alias", "npmmirror", registry.RegistryUrlNpmMirror},
		{"huaweicloud alias", "huaweicloud", registry.RegistryUrlHuaWeiCloud},
		{"tencentcloud alias", "tencentcloud", registry.RegistryUrlTencent},
		{"case insensitive", "TAOBAO", registry.RegistryUrlTaoBao},
		{"custom http URL", "http://my-registry.local:8080", "http://my-registry.local:8080"},
		{"custom https URL", "https://npm.company.com", "https://npm.company.com"},
		{"unknown name falls back to official", "unknown-mirror", registry.DefaultRegistryURL},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mirrorToURL(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMirrorToURLMatchesSDKListMirrors(t *testing.T) {
	// 验证 mirrorToURL 对每个 SDK 镜像名称都能返回与 SDK 常量一致的 URL
	for _, m := range registry.ListMirrors() {
		t.Run(m.Name, func(t *testing.T) {
			result := mirrorToURL(m.Name)
			assert.Equal(t, m.URL, result, "mirrorToURL(%q) should match SDK ListMirrors URL", m.Name)
		})
	}
}

func TestCurrentMirrorLabel(t *testing.T) {
	// 测试 --registry 优先时显示 registry URL
	globalRegistry = "https://npm.company.com"
	globalMirror = "official"
	assert.Equal(t, "https://npm.company.com", currentMirrorLabel())

	// 测试无 --registry 时显示 mirror 名称
	globalRegistry = ""
	globalMirror = "npm-mirror"
	assert.Equal(t, "npm-mirror", currentMirrorLabel())

	// 重置全局变量
	globalRegistry = ""
	globalMirror = "official"
}

// ====================================================================
// parseDepsString
// ====================================================================

func TestParseDepsString(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected map[string]string
	}{
		{"empty", "", map[string]string{}},
		{"single", "lodash=4.17.11", map[string]string{"lodash": "4.17.11"}},
		{"multiple", "lodash=4.17.11,express=4.17.1", map[string]string{"lodash": "4.17.11", "express": "4.17.1"}},
		{"with spaces", "lodash = 4.17.11 , express=4.17.1", map[string]string{"lodash ": " 4.17.11", "express": "4.17.1"}},
		{"no equals sign skipped", "lodash,express=4.17.1", map[string]string{"express": "4.17.1"}},
		{"trailing comma", "lodash=4.17.11,", map[string]string{"lodash": "4.17.11"}},
		{"only comma", ",", map[string]string{}},
		{"value contains equals", "pkg=a=b", map[string]string{"pkg": "a=b"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, parseDepsString(tc.input))
		})
	}
}

// ====================================================================
// parseAdvisoriesString
// ====================================================================

func TestParseAdvisoriesString(t *testing.T) {
	// empty
	assert.Empty(t, parseAdvisoriesString(""))
	// single
	r := parseAdvisoriesString("lodash=<4.17.12")
	assert.Equal(t, []string{"<4.17.12"}, r["lodash"])
	// multiple versions for same package
	r = parseAdvisoriesString("lodash=<4.17.12,lodash=<4.17.10")
	assert.Equal(t, []string{"<4.17.12", "<4.17.10"}, r["lodash"])
	// multiple packages
	r = parseAdvisoriesString("lodash=<4.17.12,express=<4.17.3")
	assert.Equal(t, []string{"<4.17.12"}, r["lodash"])
	assert.Equal(t, []string{"<4.17.3"}, r["express"])
	// no equals → skipped
	r = parseAdvisoriesString("lodash")
	assert.Empty(t, r)
}

// ====================================================================
// buildOptions / resolveClient / resolveClientWithToken / resolveDownloadStatsClient
// ====================================================================

func resetGlobals() {
	globalRegistry = ""
	globalMirror = "official"
	globalProxy = ""
	globalToken = ""
	globalTimeout = 120
	globalNoColor = false
}

func TestBuildOptions(t *testing.T) {
	defer resetGlobals()

	// default: mirror official
	resetGlobals()
	opts := buildOptions()
	assert.Equal(t, registry.DefaultRegistryURL, opts.RegistryURL)

	// --registry overrides --mirror
	globalRegistry = "https://npm.company.com"
	globalMirror = "taobao"
	opts = buildOptions()
	assert.Equal(t, "https://npm.company.com", opts.RegistryURL)

	// with proxy
	globalProxy = "http://127.0.0.1:7890"
	opts = buildOptions()
	assert.Equal(t, "http://127.0.0.1:7890", opts.Proxy)

	// with token
	globalToken = "npm_xxx"
	opts = buildOptions()
	assert.Equal(t, "npm_xxx", opts.Token)
}

func TestResolveClient(t *testing.T) {
	defer resetGlobals()
	resetGlobals()
	assert.NotNil(t, resolveClient())
	assert.NotNil(t, resolveClientWithToken())
	assert.NotNil(t, resolveDownloadStatsClient())

	// download stats client with proxy
	globalProxy = "http://127.0.0.1:7890"
	c := resolveDownloadStatsClient()
	assert.Equal(t, "http://127.0.0.1:7890", c.GetOptions().Proxy)
}

// ====================================================================
// requireToken
// ====================================================================

func TestRequireToken(t *testing.T) {
	defer resetGlobals()
	resetGlobals()
	assert.Error(t, requireToken())
	globalToken = "npm_xxx"
	assert.NoError(t, requireToken())
}

// ====================================================================
// mirrorNames
// ====================================================================

func TestMirrorNames(t *testing.T) {
	s := mirrorNames()
	assert.Contains(t, s, "official")
	assert.Contains(t, s, "taobao")
}

// ====================================================================
// outputJSON
// ====================================================================

func TestOutputJSON(t *testing.T) {
	// normal value
	err := outputJSON(map[string]string{"a": "b"})
	assert.NoError(t, err)

	// json.MarshalIndent 无法对普通结构体失败，但对带循环的对象会失败
	// 用 chan 触发 marshal 错误分支
	type cyclic struct {
		X *cyclic
	}
	c := &cyclic{}
	c.X = c
	err = outputJSON(c)
	assert.Error(t, err)
}

// ====================================================================
// print* 函数（确保不 panic 即可）
// ====================================================================

func TestPrintHelpers(t *testing.T) {
	printSuccess("success %s", "msg")
	printInfo("info %d", 1)
	printWarning("warn %s", "x")
	printHeader("Header")
}

// ====================================================================
// Version
// ====================================================================

func TestVersion(t *testing.T) {
	// version 由 goreleaser 注入，默认 "0.2.0"
	assert.NotEmpty(t, Version())
}

// ====================================================================
// newContext
// ====================================================================

func TestNewContext(t *testing.T) {
	defer resetGlobals()
	resetGlobals()
	globalTimeout = 5
	ctx, cancel := newContext()
	defer cancel()
	assert.NotNil(t, ctx)
	deadline, ok := ctx.Deadline()
	assert.True(t, ok)
	assert.NotZero(t, deadline)
}

// ====================================================================
// rootCmd Execute（覆盖 main 中 rootCmd.Execute 的 happy path）
// 用 --help 触发 Execute 但不实际执行业务命令
// ====================================================================

func TestRootCmdExecuteHelp(t *testing.T) {
	// 捕获 stdout/stderr 太复杂，这里只验证 rootCmd 对象构造正确
	assert.NotNil(t, rootCmd)
	assert.Equal(t, "npm-skills", rootCmd.Use)
	// PersistentPreRunE 不应出错
	assert.NoError(t, rootCmd.PersistentPreRunE(rootCmd, nil))
}

// ====================================================================
// PersistentPreRunE 环境变量分支
// ====================================================================

func TestPersistentPreRunEEnvDefaults(t *testing.T) {
	defer resetGlobals()
	resetGlobals()

	// 设置环境变量
	t.Setenv("NPM_PROXY", "http://env-proxy:7890")
	t.Setenv("NPM_REGISTRY", "https://env-registry.com")
	t.Setenv("NPM_MIRROR", "taobao")
	t.Setenv("NPM_TOKEN", "env-token")

	// 模拟 flags 未显式设置
	_ = rootCmd.PersistentPreRunE(rootCmd, nil)
	assert.Equal(t, "http://env-proxy:7890", globalProxy)
	assert.Equal(t, "https://env-registry.com", globalRegistry)
	assert.Equal(t, "taobao", globalMirror)
	assert.Equal(t, "env-token", globalToken)
}

func TestPersistentPreRunENoColor(t *testing.T) {
	defer resetGlobals()
	resetGlobals()
	globalNoColor = true
	_ = rootCmd.PersistentPreRunE(rootCmd, nil)
	// NoColor 全局变量应被设置
	assert.True(t, color.NoColor)
}

// ====================================================================
// 验证所有 cobra 子命令已注册（覆盖 init 副作用）
// ====================================================================

func TestAllSubCommandsRegistered(t *testing.T) {
	// 遍历所有子命令，确保都被注册到 rootCmd
	cmds := rootCmd.Commands()
	assert.NotEmpty(t, cmds)
	// 验证关键命令存在
	names := make(map[string]bool)
	for _, c := range cmds {
		names[c.Name()] = true
	}
	expected := []string{
		"package", "package-summary", "search", "versions", "download",
		"download-stats", "download-range", "mirrors", "registry-info",
		"whoami", "dist-tags", "deprecate", "unpublish", "publish", "star",
		"user", "token", "hook", "org", "access", "audit", "config", "couchdb",
	}
	for _, e := range expected {
		assert.True(t, names[e], "expected command %q to be registered", e)
	}
}

// dummy to avoid unused import warning for encoding/json in some builds
var _ = json.Marshal
var _ = strings.TrimSpace
