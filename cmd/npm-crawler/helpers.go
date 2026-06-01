package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/scagogogo/npm-crawler/pkg/registry"
)

// resolveClient creates a registry client based on global flags (--mirror, --registry, --proxy)
// Priority: --registry > --mirror > default
func resolveClient() *registry.Registry {
	var opts *registry.Options

	// If --registry is set, use it directly (highest priority)
	if globalRegistry != "" {
		opts = registry.NewOptions().SetRegistryURL(globalRegistry)
	} else {
		// Otherwise, resolve mirror name to URL
		opts = registry.NewOptions().SetRegistryURL(mirrorToURL(globalMirror))
	}

	// Apply proxy
	if globalProxy != "" {
		opts.SetProxy(globalProxy)
	}

	return registry.NewRegistry(opts)
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
func outputJSON(v interface{}) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, color.RedString("✗ JSON marshal error: %v", err))
		os.Exit(1)
	}
	fmt.Println(string(data))
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
