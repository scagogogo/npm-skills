package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/scagogogo/npm-skills/pkg/registry"
)

// parseDepsString 解析依赖字符串为 map[name]version
// 格式: "lodash=4.17.11,express=4.17.1"
func parseDepsString(raw string) map[string]string {
	deps := make(map[string]string)
	for _, pair := range strings.Split(raw, ",") {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			deps[parts[0]] = parts[1]
		}
	}
	return deps
}

// parseAdvisoriesString 解析公告字符串为 map[name][]version
// 格式: "lodash=<4.17.12,express=<4.17.3"
func parseAdvisoriesString(raw string) map[string][]string {
	advisories := make(map[string][]string)
	for _, pair := range strings.Split(raw, ",") {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			advisories[parts[0]] = append(advisories[parts[0]], parts[1])
		}
	}
	return advisories
}

// buildOptions 根据全局 flags 构建 Options（公共逻辑）
// Priority: --registry > --mirror > default
func buildOptions() *registry.Options {
	var opts *registry.Options

	if globalRegistry != "" {
		opts = registry.NewOptions().SetRegistryURL(globalRegistry)
	} else {
		opts = registry.NewOptions().SetRegistryURL(mirrorToURL(globalMirror))
	}

	if globalProxy != "" {
		opts.SetProxy(globalProxy)
	}

	if globalToken != "" {
		opts.SetToken(globalToken)
	}

	return opts
}

// resolveClient creates a registry client based on global flags (--mirror, --registry, --proxy)
func resolveClient() *registry.Registry {
	return registry.NewRegistry(buildOptions())
}

// resolveClientWithToken creates a registry client with authentication token
// Call requireToken() before this function to ensure token is set
func resolveClientWithToken() *registry.Registry {
	return registry.NewRegistry(buildOptions())
}

// resolveDownloadStatsClient creates a registry client for download stats API
// 下载统计 API 使用独立的 api.npmjs.org 端点，不走 Registry URL，
// 但仍然需要支持代理和可配置的下载统计 URL（用于私有仓库）
func resolveDownloadStatsClient() *registry.Registry {
	opts := registry.NewOptions()
	if globalProxy != "" {
		opts.SetProxy(globalProxy)
	}
	return registry.NewRegistry(opts)
}

// requireToken validates that a token is set, returns error if not
func requireToken() error {
	if globalToken == "" {
		return fmt.Errorf("authentication required: use --token flag or set NPM_TOKEN environment variable")
	}
	return nil
}

// mirrorToURL converts a mirror name to its registry URL
// 使用 SDK 的 ListMirrors() 作为唯一数据源
func mirrorToURL(name string) string {
	// 先在 SDK 镜像列表中查找
	lowerName := strings.ToLower(name)
	for _, m := range registry.ListMirrors() {
		if strings.ToLower(m.Name) == lowerName {
			return m.URL
		}
	}

	// 特殊别名支持（SDK 不包含，CLI 独有）
	switch lowerName {
	case "npmmirror":
		return registry.RegistryUrlNpmMirror
	case "huaweicloud":
		return registry.RegistryUrlHuaWeiCloud
	case "tencentcloud":
		return registry.RegistryUrlTencent
	}

	// If it looks like a URL, use it directly
	if strings.HasPrefix(strings.ToLower(name), "http://") ||
		strings.HasPrefix(strings.ToLower(name), "https://") {
		return name
	}

	// Fallback to official
	return registry.DefaultRegistryURL
}

// outputJSON prints a value as formatted JSON to stdout
// 返回 error 以便 Cobra 的 RunE 可以正确处理错误而不使用 os.Exit
func outputJSON(v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("JSON marshal error: %w", err)
	}
	fmt.Println(string(data))
	return nil
}

// printSuccess prints a green success message to stderr
func printSuccess(format string, args ...interface{}) {
	fmt.Fprintln(os.Stderr, color.New(color.FgGreen).Sprintf(format, args...))
}

// printInfo prints a cyan info message to stderr
func printInfo(format string, args ...interface{}) {
	fmt.Fprintln(os.Stderr, color.New(color.FgCyan).Sprintf(format, args...))
}

// printWarning prints a yellow warning message to stderr
func printWarning(format string, args ...interface{}) {
	fmt.Fprintln(os.Stderr, color.New(color.FgYellow).Sprintf(format, args...))
}

// printHeader prints a bold section header to stderr
func printHeader(title string) {
	fmt.Fprintln(os.Stderr, color.New(color.FgCyan, color.Bold).Sprintf("▸ %s", title))
}

// mirrorNames returns all known mirror names for help text
func mirrorNames() string {
	return "official|taobao|npm-mirror|huawei|tencent|cnpm|yarn|npmjscom"
}

// currentMirrorLabel returns a human-readable label for the current mirror/registry setting
func currentMirrorLabel() string {
	if globalRegistry != "" {
		return globalRegistry
	}
	return globalMirror
}