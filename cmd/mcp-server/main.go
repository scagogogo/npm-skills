package main

import (
	"context"
	"fmt"
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

func main() {
	// Parse command-line flags
	registryURL := os.Getenv("NPM_REGISTRY")
	mirror := getEnvOrDefault("NPM_MIRROR", "official")
	proxy := os.Getenv("NPM_PROXY")
	token := os.Getenv("NPM_TOKEN")
	timeoutStr := getEnvOrDefault("NPM_TIMEOUT", "120")

	// Override env vars with command-line args if provided
	for i, arg := range os.Args[1:] {
		switch {
		case arg == "--registry" && i+1 < len(os.Args)-1:
			registryURL = os.Args[i+2]
		case arg == "--mirror" && i+1 < len(os.Args)-1:
			mirror = os.Args[i+2]
		case arg == "--proxy" && i+1 < len(os.Args)-1:
			proxy = os.Args[i+2]
		case arg == "--token" && i+1 < len(os.Args)-1:
			token = os.Args[i+2]
		case arg == "--timeout" && i+1 < len(os.Args)-1:
			timeoutStr = os.Args[i+2]
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
			printHelp()
			os.Exit(0)
		}
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

	if err := stdioServer.Listen(ctx, os.Stdin, os.Stdout); err != nil {
		log.Fatalf("Server error: %v\n", err)
	}
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

func printHelp() {
	fmt.Println("NPM Registry MCP Server")
	fmt.Println()
	fmt.Println("Exposes NPM registry operations as MCP tools for AI agents.")
	fmt.Println()
	fmt.Println("Usage: npm-mcp-server [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --registry URL     Custom registry URL (overrides --mirror)")
	fmt.Println("  --mirror NAME      Mirror source (default: official)")
	fmt.Println("                      Values: official|taobao|npm-mirror|huawei|tencent|cnpm|yarn|npmjscom")
	fmt.Println("  --proxy URL        HTTP proxy URL (e.g. http://127.0.0.1:7890)")
	fmt.Println("  --token TOKEN      NPM auth token (for whoami and private packages)")
	fmt.Println("  --timeout SECS     Request timeout in seconds (default: 120)")
	fmt.Println()
	fmt.Println("Environment variables (used as defaults):")
	fmt.Println("  NPM_REGISTRY       Custom registry URL")
	fmt.Println("  NPM_MIRROR         Mirror source name")
	fmt.Println("  NPM_PROXY          HTTP proxy URL")
	fmt.Println("  NPM_TOKEN          NPM auth token")
	fmt.Println("  NPM_TIMEOUT        Request timeout in seconds")
	fmt.Println()
	fmt.Println("Priority: CLI flag > Environment variable > Default")
	fmt.Println()
	fmt.Println("MCP Tools (12 total):")
	fmt.Println("  npm_registry_info     — Registry status and statistics")
	fmt.Println("  npm_mirrors           — List available mirror sources")
	fmt.Println("  npm_package           — Full package metadata (large response)")
	fmt.Println("  npm_package_summary   — Lightweight package metadata (recommended)")
	fmt.Println("  npm_search            — Search packages by keyword")
	fmt.Println("  npm_version           — Specific version metadata")
	fmt.Println("  npm_versions          — All published version numbers")
	fmt.Println("  npm_latest_version    — Latest version number")
	fmt.Println("  npm_dist_tags         — Distribution tags (latest, next, beta)")
	fmt.Println("  npm_download_stats    — Download count for a period")
	fmt.Println("  npm_download_range    — Daily download trend data")
	fmt.Println("  npm_whoami            — Check auth status (requires --token)")
	fmt.Println()
	fmt.Println("Claude Code integration:")
	fmt.Println("  Add to your settings:")
	fmt.Println(`  {
    "mcpServers": {
      "npm-registry": {
        "command": "npm-mcp-server",
        "args": ["--mirror", "npm-mirror"]
      }
    }
  }`)
}