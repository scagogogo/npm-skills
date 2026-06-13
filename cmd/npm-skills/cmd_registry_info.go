package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var registryInfoCmd = &cobra.Command{
	Use:   "registry-info",
	Short: "Get NPM registry status information",
	Long: color.New(color.FgCyan).Sprintf("Get NPM registry status information") + "\n\n" +
		"Retrieves database name, total package count, disk size, and other registry metadata.\n\n" +
		color.HiBlackString("Mirror: %s", mirrorNames()) + " (via --mirror or --registry flag)",
	Aliases: []string{"info"},
	Example: `  npm-skills registry-info
  npm-skills registry-info -m npm-mirror
  npm-skills registry-info --registry https://registry.npmmirror.com
  npm-skills info --proxy http://127.0.0.1:7890`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		printInfo("Fetching registry info from %s...", currentMirrorLabel())
		client := resolveClient()
		ctx, cancel := newContext()
		defer cancel()

		info, err := client.GetRegistryInformation(ctx)
		if err != nil {
			return fmt.Errorf("failed to get registry information: %w", err)
		}

		if err := outputJSON(info); err != nil {
			return err
		}
		printSuccess("✓ Registry information retrieved successfully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(registryInfoCmd)
}
