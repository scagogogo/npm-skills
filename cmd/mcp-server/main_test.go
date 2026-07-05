package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"testing"
	"time"

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
	var buf bytes.Buffer
	assert.NotPanics(t, func() {
		printHelp(&buf)
	})
	assert.Contains(t, buf.String(), "NPM Registry MCP Server")
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

// ====================================================================
// parseArgs（抽离自 main 的参数解析）
// ====================================================================

func TestParseArgsDefaults(t *testing.T) {
	// 清空相关 env，验证默认值
	t.Setenv("NPM_REGISTRY", "")
	t.Setenv("NPM_MIRROR", "")
	t.Setenv("NPM_PROXY", "")
	t.Setenv("NPM_TOKEN", "")
	t.Setenv("NPM_TIMEOUT", "")

	reg, mirror, proxy, token, timeout, help := parseArgs(nil)
	assert.Empty(t, reg)
	assert.Equal(t, "official", mirror)
	assert.Empty(t, proxy)
	assert.Empty(t, token)
	assert.Equal(t, "120", timeout)
	assert.False(t, help)
}

func TestParseArgsEnvValues(t *testing.T) {
	t.Setenv("NPM_REGISTRY", "https://env-reg.com")
	t.Setenv("NPM_MIRROR", "taobao")
	t.Setenv("NPM_PROXY", "http://env-proxy:7890")
	t.Setenv("NPM_TOKEN", "env-token")
	t.Setenv("NPM_TIMEOUT", "60")

	reg, mirror, proxy, token, timeout, help := parseArgs(nil)
	assert.Equal(t, "https://env-reg.com", reg)
	assert.Equal(t, "taobao", mirror)
	assert.Equal(t, "http://env-proxy:7890", proxy)
	assert.Equal(t, "env-token", token)
	assert.Equal(t, "60", timeout)
	assert.False(t, help)
}

func TestParseArgsSpaceForm(t *testing.T) {
	// --flag value 形式（覆盖 i+1 < len(args) 分支）
	args := []string{"--registry", "https://r1.com", "--mirror", "huawei",
		"--proxy", "http://p:1", "--token", "t1", "--timeout", "30"}
	reg, mirror, proxy, token, timeout, help := parseArgs(args)
	assert.Equal(t, "https://r1.com", reg)
	assert.Equal(t, "huawei", mirror)
	assert.Equal(t, "http://p:1", proxy)
	assert.Equal(t, "t1", token)
	assert.Equal(t, "30", timeout)
	assert.False(t, help)
}

func TestParseArgsEqualForm(t *testing.T) {
	// --flag=value 形式
	args := []string{
		"--registry=https://r2.com", "--mirror=tencent",
		"--proxy=http://p:2", "--token=t2", "--timeout=45",
	}
	reg, mirror, proxy, token, timeout, help := parseArgs(args)
	assert.Equal(t, "https://r2.com", reg)
	assert.Equal(t, "tencent", mirror)
	assert.Equal(t, "http://p:2", proxy)
	assert.Equal(t, "t2", token)
	assert.Equal(t, "45", timeout)
	assert.False(t, help)
}

func TestParseArgsSpaceFormTrailingFlagNoValue(t *testing.T) {
	// --registry 在末尾，无后续值 → 不覆盖（i+1 < len(args) 为 false）
	args := []string{"--registry"}
	reg, _, _, _, _, _ := parseArgs(args)
	assert.Empty(t, reg)
}

func TestParseArgsHelpShortAndLong(t *testing.T) {
	_, _, _, _, _, help := parseArgs([]string{"-h"})
	assert.True(t, help)
	_, _, _, _, _, help = parseArgs([]string{"--help"})
	assert.True(t, help)
}

// ====================================================================
// run（main 抽离出的可测函数）
// ====================================================================

func TestRunHelp(t *testing.T) {
	var buf bytes.Buffer
	code := run([]string{"--help"}, strings.NewReader(""), &buf)
	assert.Equal(t, 0, code)
	assert.Contains(t, buf.String(), "NPM Registry MCP Server")
}

func TestRunStdioEOF(t *testing.T) {
	// 空 stdin → Listen 立即收到 io.EOF 返回 nil → run 返回 0
	// 用 io.Discard 作为 stdout，避免污染测试输出
	var stderr bytes.Buffer
	log.SetOutput(&stderr)
	defer log.SetOutput(os.Stderr)
	code := run(nil, strings.NewReader(""), io.Discard)
	assert.Equal(t, 0, code)
}

func TestRunInvalidTimeout(t *testing.T) {
	// 非法 timeout → 回退到 120s（time.Duration 显示为 "2m0s"），仍正常启动并因空 stdin 立即返回 0
	var stderr bytes.Buffer
	log.SetOutput(&stderr)
	defer log.SetOutput(os.Stderr)
	code := run([]string{"--timeout=not-a-number"}, strings.NewReader(""), io.Discard)
	assert.Equal(t, 0, code)
	assert.Contains(t, stderr.String(), "2m0s")
}

func TestRunSignalShutdown(t *testing.T) {
	// 用阻塞 stdin（pipe 读端不关闭）让 Listen 阻塞在 readNextLine，
	// 发 SIGINT 触发 signal goroutine → cancel context → Listen 返回 ctx.Err → run 返回 1。
	// 覆盖 signal goroutine 与 Listen err 分支（log.Printf + return 1）。
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)
	defer signal.Stop(sigChan)

	pr, pw := io.Pipe()
	defer pw.Close()
	defer pr.Close()

	var stderr bytes.Buffer
	log.SetOutput(&stderr)
	defer log.SetOutput(os.Stderr)

	codeCh := make(chan int, 1)
	go func() {
		codeCh <- run(nil, pr, io.Discard)
	}()

	// 等 Listen 进入阻塞读（goroutine 已注册 signal.Notify）
	time.Sleep(100 * time.Millisecond)

	// 触发 signal goroutine
	_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)

	select {
	case code := <-codeCh:
		assert.Equal(t, 1, code)
		assert.Contains(t, stderr.String(), "Server error")
	case <-time.After(3 * time.Second):
		t.Fatal("run did not return after SIGINT")
	}
}
