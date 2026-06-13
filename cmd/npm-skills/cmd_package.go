package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var packageCmd = &cobra.Command{
	Use:   "package <name>",
	Short: "Get NPM package information",
	Long: color.New(color.FgCyan).Sprintf("Get NPM package information") + "\n\n" +
		"Retrieves full package metadata including description, versions, dependencies,\n" +
		"maintainers, license, README content, and more.\n\n" +
		color.HiBlackString("Mirror: %s", mirrorNames()) + " (via --mirror or --registry flag)",
	Aliases: []string{"pkg"},
	Example: `  npm-skills package react
  npm-skills package axios -m taobao
  npm-skills pkg @nestjs/core --registry https://registry.npmmirror.com
  npm-skills package react --proxy http://127.0.0.1:7890`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		packageName := args[0]

		printInfo("Fetching package info for %s from %s...",
			color.New(color.FgWhite, color.Bold).Sprint(packageName),
			currentMirrorLabel())
		client := resolveClient()
		ctx, cancel := newContext()
		defer cancel()

		pkg, err := client.GetPackageInformation(ctx, packageName)
		if err != nil {
			return fmt.Errorf("failed to get package info for '%s': %w", packageName, err)
		}

		if err := outputJSON(pkg); err != nil {
			return err
		}
		printSuccess("✓ Package information for %s retrieved successfully",
			color.New(color.FgWhite, color.Bold).Sprint(packageName))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(packageCmd)
}