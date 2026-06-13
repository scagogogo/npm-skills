package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var pkgVersionCmd = &cobra.Command{
	Use:   "pkg-version <name> <version>",
	Short: "Get specific version information of an NPM package",
	Long: color.New(color.FgCyan).Sprintf("Get specific version information of an NPM package") + "\n\n" +
		"Retrieves version-specific details including dependencies, devDependencies,\n" +
		"scripts, dist info (tarball URL, shasum, integrity hash), and license.\n" +
		"Use \"latest\" as the version to get the latest version info.\n\n" +
		color.HiBlackString("Mirror: %s", mirrorNames()) + " (via --mirror or --registry flag)",
	Aliases: []string{"ver", "pv"},
	Example: `  npm-skills pkg-version react 18.2.0
  npm-skills pkg-version axios latest -m npm-mirror
  npm-skills ver lodash 4.17.21 --proxy http://127.0.0.1:7890`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		packageName := args[0]
		ver := args[1]

		printInfo("Fetching version info for %s from %s...",
			color.New(color.FgWhite, color.Bold).Sprintf("%s@%s", packageName, ver),
			currentMirrorLabel())
		client := resolveClient()
		ctx, cancel := newContext()
		defer cancel()

		v, err := client.GetPackageVersion(ctx, packageName, ver)
		if err != nil {
			return fmt.Errorf("failed to get version info for '%s@%s': %w", packageName, ver, err)
		}

		if err := outputJSON(v); err != nil {
			return err
		}
		printSuccess("✓ Version info for %s retrieved successfully",
			color.New(color.FgWhite, color.Bold).Sprintf("%s@%s", packageName, ver))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pkgVersionCmd)
}
