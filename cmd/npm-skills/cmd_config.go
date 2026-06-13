package main

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show current effective configuration",
	Long: color.New(color.FgCyan).Sprintf("Show current effective configuration") + "\n\n" +
		"Displays the currently active registry URL, mirror, proxy, and timeout settings.\n" +
		"Useful for debugging which registry/mirror/proxy is actually being used\n" +
		"after applying CLI flags and environment variables.",
	Aliases: []string{"cfg", "conf"},
	Example: `  npm-skills config
  npm-skills config -m npm-mirror
  npm-skills config --registry https://registry.npmmirror.com --proxy http://127.0.0.1:7890`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := resolveClient()
		opts := client.GetOptions()

		// Determine the effective mirror name (for display only)
		effectiveMirror := currentMirrorLabel()

		config := map[string]interface{}{
			"registry_url": opts.RegistryURL,
			"mirror":       effectiveMirror,
			"proxy":        opts.Proxy,
			"timeout":      globalTimeout,
		}
		if opts.Proxy == "" {
			config["proxy"] = "(none)"
		}

		printInfo("Current effective configuration:")
		if err := outputJSON(config); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}