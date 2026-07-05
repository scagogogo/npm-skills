package main

import (
	"os"
	"os/exec"
	"testing"

	"github.com/scagogogo/npm-skills/pkg/registry"
	"github.com/stretchr/testify/assert"
)

// ====================================================================
// mirrorNameToURL
// ====================================================================

func TestMirrorNameToURL(t *testing.T) {
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
		{"empty falls back to official", "", registry.DefaultRegistryURL},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mirrorNameToURL(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMirrorNameToURLMatchesSDKListMirrors(t *testing.T) {
	for _, m := range registry.ListMirrors() {
		t.Run(m.Name, func(t *testing.T) {
			result := mirrorNameToURL(m.Name)
			assert.Equal(t, m.URL, result, "mirrorNameToURL(%q) should match SDK ListMirrors URL", m.Name)
		})
	}
}

// ====================================================================
// getEnvOrDefault
// ====================================================================

func TestGetEnvOrDefault(t *testing.T) {
	// 设置环境变量
	t.Setenv("MY_TEST_VAR", "value_from_env")
	assert.Equal(t, "value_from_env", getEnvOrDefault("MY_TEST_VAR", "default"))

	// 未设置时返回默认值
	assert.Equal(t, "default", getEnvOrDefault("UNSET_VAR_XYZ_123", "default"))

	// 空字符串环境变量视为未设置
	t.Setenv("EMPTY_VAR", "")
	assert.Equal(t, "default", getEnvOrDefault("EMPTY_VAR", "default"))
}

// ====================================================================
// buildOptions
// ====================================================================

func TestBuildOptions(t *testing.T) {
	// 仅 registryURL
	opts := buildOptions("https://custom.registry.com", "taobao", "", "")
	assert.Equal(t, "https://custom.registry.com", opts.RegistryURL)
	assert.Empty(t, opts.Proxy)
	assert.Empty(t, opts.Token)

	// registryURL 为空时用 mirror
	opts = buildOptions("", "taobao", "", "")
	assert.Equal(t, registry.RegistryUrlTaoBao, opts.RegistryURL)

	// 带 proxy 和 token
	opts = buildOptions("", "official", "http://127.0.0.1:7890", "npm_xxx")
	assert.Equal(t, registry.DefaultRegistryURL, opts.RegistryURL)
	assert.Equal(t, "http://127.0.0.1:7890", opts.Proxy)
	assert.Equal(t, "npm_xxx", opts.Token)

	// 全空 → official
	opts = buildOptions("", "", "", "")
	assert.Equal(t, registry.DefaultRegistryURL, opts.RegistryURL)
}

// ====================================================================
// printHelp
// ====================================================================

func TestPrintHelp(t *testing.T) {
	// 仅验证不 panic
	assert.NotPanics(t, func() {
		printHelp()
	})
}

// ====================================================================
// main（通过 exec 运行二进制，覆盖 --help 分支）
// ====================================================================

func TestMainHelpFlag(t *testing.T) {
	// 编译当前测试二进制并运行 --help
	bin := os.TempDir() + "/npm-mcp-server-test"
	build := exec.Command("go", "build", "-o", bin, ".")
	build.Dir = "."
	if out, err := build.CombinedOutput(); err != nil {
		t.Skipf("failed to build: %v\n%s", err, out)
	}
	defer os.Remove(bin)

	out, err := exec.Command(bin, "--help").CombinedOutput()
	assert.NoError(t, err)
	s := string(out)
	assert.Contains(t, s, "NPM Registry MCP Server")
	assert.Contains(t, s, "npm_mirrors")
}
