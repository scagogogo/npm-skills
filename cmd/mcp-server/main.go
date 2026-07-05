package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	mcpserver "github.com/mark3labs/mcp-go/server"

	npmMcp "github.com/scagogogo/npm-skills/pkg/mcp"
	"github.com/scagogogo/npm-skills/pkg/registry"
)

// parseArgs 解析命令行参数（覆盖 env 默认值），返回各项配置和是否请求 help。
// 抽离自 main 以便单元测试。
func parseArgs(args []string) (registryURL, mirror, proxy, token, timeoutStr string, helpRequested bool) {
	registryURL = os.Getenv("NPM_REGISTRY")
	mirror = getEnvOrDefault("NPM_MIRROR", "official")
	proxy = os.Getenv("NPM_PROXY")
	token = os.Getenv("NPM_TOKEN")
	timeoutStr = getEnvOrDefault("NPM_TIMEOUT", "120")

	// Override env vars with command-line args if provided
	for i, arg := range args {
		switch {
		case arg == "--registry" && i+1 < len(args):
			registryURL = args[i+1]
		case arg == "--mirror" && i+1 < len(args):
			mirror = args[i+1]
		case arg == "--proxy" && i+1 < len(args):
			proxy = args[i+1]
		case arg == "--token" && i+1 < len(args):
			token = args[i+1]
		case arg == "--timeout" && i+1 < len(args):
			timeoutStr = args[i+1]
		case strings.HasPrefix(arg, "--registry="):
			registryURL = strings.TrimPrefix(arg, "--registry=")
		case strings.HasPrefix(arg, "--mirror="):
			mirror = strings.TrimPrefix(arg, "--mirror=")
		case strings.HasPrefix(arg, "--proxy="):
			proxy = strings.TrimPrefix(arg, "--proxy=")
		case strings.HasPrefix(arg, "--token="):
			token = strings.TrimPrefix(arg, "--token=")
		case strings.HasPrefix(arg, "--timeout="):
			timeoutStr = strings.TrimPrefix(arg, "--timeout=")
		case arg == "--help" || arg == "-h":
			helpRequested = true
		}
	}
	return
}

// run 是主入口，抽离自 main 以便单元测试（main 本身因 os.Exit 无法直接测试）。
// 参数：args 为命令行参数（不含程序名），stdin/stdout 用于 stdio server。
// 返回进程退出码。
func run(args []string, stdin io.Reader, stdout io.Writer) int {
	registryURL, mirror, proxy, token, timeoutStr, helpRequested := parseArgs(args)
	if helpRequested {
		printHelp(stdout)
		return 0
	}

	timeout, err := time.ParseDuration(timeoutStr + "s")
	if err != nil {
		timeout = 120 * time.Second
	}

	// Build registry options
	opts := buildOptions(registryURL, mirror, proxy, token)

	cfg := npmMcp.Config{
		RegistryOptions: opts,
		Timeout:         timeout,
	}

	mcpSrv := npmMcp.NewServer(cfg)

	// Start stdio server
	stdioServer := mcpserver.NewStdioServer(mcpSrv)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-sigChan
		log.Printf("Received signal %v, shutting down...\n", sig)
		cancel()
	}()

	log.Printf("Starting NPM Registry MCP Server (mirror: %s, timeout: %s)\n", mirror, timeout)

	if err := stdioServer.Listen(ctx, stdin, stdout); err != nil {
		log.Printf("Server error: %v\n", err)
		return 1
	}
	return 0
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout))
}

func buildOptions(registryURL, mirror, proxy, token string) *registry.Options {
	opts := registry.NewOptions()

	if registryURL != "" {
		opts.SetRegistryURL(registryURL)
	} else {
		opts.SetRegistryURL(mirrorNameToURL(mirror))
	}

	if proxy != "" {
		opts.SetProxy(proxy)
	}

	if token != "" {
		opts.SetToken(token)
	}

	return opts
}

// mirrorNameToURL converts a mirror name to its registry URL
func mirrorNameToURL(name string) string {
	lowerName := strings.ToLower(name)
	for _, m := range registry.ListMirrors() {
		if strings.ToLower(m.Name) == lowerName {
			return m.URL
		}
	}

	// Special aliases
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

	return registry.DefaultRegistryURL
}

func getEnvOrDefault(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func printHelp(w io.Writer) {
	fmt.Fprintln(w, "NPM Registry MCP Server")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Exposes NPM registry operations as MCP tools for AI agents.")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Usage: npm-mcp-server [flags]")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Flags:")
	fmt.Fprintln(w, "  --registry URL     Custom registry URL (overrides --mirror)")
	fmt.Fprintln(w, "  --mirror NAME      Mirror source (default: official)")
	fmt.Fprintln(w, "                      Values: official|taobao|npm-mirror|huawei|tencent|cnpm|yarn|npmjscom")
	fmt.Fprintln(w, "  --proxy URL        HTTP proxy URL (e.g. http://127.0.0.1:7890)")
	fmt.Fprintln(w, "  --token TOKEN      NPM auth token (for whoami and private packages)")
	fmt.Fprintln(w, "  --timeout SECS     Request timeout in seconds (default: 120)")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Environment variables (used as defaults):")
	fmt.Fprintln(w, "  NPM_REGISTRY       Custom registry URL")
	fmt.Fprintln(w, "  NPM_MIRROR         Mirror source name")
	fmt.Fprintln(w, "  NPM_PROXY          HTTP proxy URL")
	fmt.Fprintln(w, "  NPM_TOKEN          NPM auth token")
	fmt.Fprintln(w, "  NPM_TIMEOUT        Request timeout in seconds")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Priority: CLI flag > Environment variable > Default")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "MCP Tools (12 total):")
	fmt.Fprintln(w, "  npm_registry_info     — Registry status and statistics")
	fmt.Fprintln(w, "  npm_mirrors           — List available mirror sources")
	fmt.Fprintln(w, "  npm_package           — Full package metadata (large response)")
	fmt.Fprintln(w, "  npm_package_summary   — Lightweight package metadata (recommended)")
	fmt.Fprintln(w, "  npm_search            — Search packages by keyword")
	fmt.Fprintln(w, "  npm_version           — Specific version metadata")
	fmt.Fprintln(w, "  npm_versions          — All published version numbers")
	fmt.Fprintln(w, "  npm_latest_version    — Latest version number")
	fmt.Fprintln(w, "  npm_dist_tags         — Distribution tags (latest, next, beta)")
	fmt.Fprintln(w, "  npm_download_stats    — Download count for a period")
	fmt.Fprintln(w, "  npm_download_range    — Daily download trend data")
	fmt.Fprintln(w, "  npm_whoami            — Check auth status (requires --token)")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Claude Code integration:")
	fmt.Fprintln(w, "  Add to your settings:")
	fmt.Fprintln(w, `  {
    "mcpServers": {
      "npm-registry": {
        "command": "npm-mcp-server",
        "args": ["--mirror", "npm-mirror"]
      }
    }
  }`)
}
