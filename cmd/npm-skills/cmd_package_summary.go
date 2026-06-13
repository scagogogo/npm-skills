package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var packageSummaryCmd = &cobra.Command{
	Use:   "package-summary <name>",
	Short: "Get lightweight NPM package metadata",
	Long: color.New(color.FgCyan).Sprintf("Get lightweight NPM package metadata") + "\n\n" +
		"Returns abbreviated package information using the install-v1 Accept header.\n" +
		"Much smaller response than 'package' (KB vs MB) — ideal for scripts and quick lookups.\n" +
		"May lack README, maintainers, and some other fields.\n\n" +
		color.HiBlackString("Mirror: %s", mirrorNames()) + " (via --mirror or --registry flag)",
	Aliases: []string{"ps", "pkgsum"},
	Example: `  npm-skills package-summary react
  npm-skills ps axios -m taobao
  npm-skills pkgsum @nestjs/core --proxy http://127.0.0.1:7890`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		packageName := args[0]

		printInfo("Fetching package summary for %s from %s...",
			color.New(color.FgWhite, color.Bold).Sprint(packageName),
			currentMirrorLabel())
		client := resolveClient()
		ctx, cancel := newContext()
		defer cancel()

		pkg, err := client.GetAbbreviatedPackageInformation(ctx, packageName)
		if err != nil {
			return fmt.Errorf("failed to get package summary for '%s': %w", packageName, err)
		}

		if err := outputJSON(pkg); err != nil {
			return err
		}
		printSuccess("✓ Package summary for %s retrieved successfully",
			color.New(color.FgWhite, color.Bold).Sprint(packageName))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(packageSummaryCmd)
}
