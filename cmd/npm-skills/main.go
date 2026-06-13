package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// Build variables injected by goreleaser via -ldflags
var (
	version = "0.2.0"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

// Version returns the full version string
func Version() string {
	return version
}

// Global flags
var (
	globalMirror      string
	globalRegistry    string
	globalProxy       string
	globalToken       string
	globalTimeout     int
	globalNoColor     bool
)

var rootCmd = &cobra.Command{
	Use:   "npm-skills",
	Short: color.New(color.FgCyan, color.Bold).Sprint("NPM Registry CLI Tool"),
	Long: fmt.Sprintf(`
  %s
  %s

  %s  npm-skills <command> [flags]

  %s
    npm-skills package react
    npm-skills search "http client" -l 10
    npm-skills download-stats axios -p last-month
    npm-skills download lodash 4.17.21 ./lodash.tgz
    npm-skills mirrors

  %s
    npm-skills package react -m npm-mirror
    npm-skills package react --proxy http://127.0.0.1:7890
    npm-skills package react --registry https://registry.npmmirror.com
    NPM_REGISTRY=https://registry.npmmirror.com npm-skills package react`,
		color.New(color.FgCyan, color.Bold).Sprint("NPM Crawler — NPM Registry CLI Tool"),
		color.HiBlackString("Query package info, search, download stats, and tarballs with mirror & proxy support."),
		color.New(color.FgGreen).Sprint("Usage:"),
		color.New(color.FgGreen).Sprint("Examples:"),
		color.New(color.FgGreen).Sprint("Proxy / Mirror:"),
	),
	Version: version,
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Disable color if requested
		if globalNoColor {
			color.NoColor = true
		}

		// Apply environment variable defaults if flags not explicitly set
		if !cmd.Flags().Changed("proxy") && os.Getenv("NPM_PROXY") != "" {
			globalProxy = os.Getenv("NPM_PROXY")
		}
		if !cmd.Flags().Changed("registry") && os.Getenv("NPM_REGISTRY") != "" {
			globalRegistry = os.Getenv("NPM_REGISTRY")
		}
		if !cmd.Flags().Changed("mirror") && os.Getenv("NPM_MIRROR") != "" {
			globalMirror = os.Getenv("NPM_MIRROR")
		}
		if !cmd.Flags().Changed("token") && os.Getenv("NPM_TOKEN") != "" {
			globalToken = os.Getenv("NPM_TOKEN")
		}

		return nil
	},
}

func init() {
	rootCmd.SetVersionTemplate(fmt.Sprintf("%s %s\n",
		color.New(color.FgCyan, color.Bold).Sprint("npm-skills"),
		color.New(color.FgGreen).Sprintf("v%s", version),
	))

	// Global persistent flags
	rootCmd.PersistentFlags().StringVarP(&globalMirror, "mirror", "m", "official",
		"Mirror source: "+mirrorNames()+" (env: NPM_MIRROR)")
	rootCmd.PersistentFlags().StringVar(&globalRegistry, "registry", "",
		"Custom registry URL (overrides --mirror, env: NPM_REGISTRY)")
	rootCmd.PersistentFlags().StringVar(&globalProxy, "proxy", "",
		"HTTP proxy URL, e.g. http://127.0.0.1:7890 (env: NPM_PROXY)")
	rootCmd.PersistentFlags().IntVar(&globalTimeout, "timeout", 120,
		"Request timeout in seconds")
	rootCmd.PersistentFlags().BoolVar(&globalNoColor, "no-color", false,
		"Disable colored output")
	rootCmd.PersistentFlags().StringVarP(&globalToken, "token", "t", "",
		"NPM authentication token (env: NPM_TOKEN)")
}

// newContext creates a context with the global timeout
func newContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(globalTimeout)*time.Second)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, color.New(color.FgRed).Sprintf("✗ Error: %s", err))
		os.Exit(1)
	}
}
