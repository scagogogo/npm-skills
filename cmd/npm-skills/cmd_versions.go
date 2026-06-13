package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var versionsLatestOnly bool

var versionsCmd = &cobra.Command{
	Use:   "versions <name>",
	Short: "List all published versions of an NPM package",
	Long: color.New(color.FgCyan).Sprintf("List all published versions of an NPM package") + "\n\n" +
		"Returns a sorted list of all published version numbers.\n" +
		"Use --latest to show only the latest version.\n\n" +
		color.HiBlackString("Mirror: %s", mirrorNames()) + " (via --mirror or --registry flag)",
	Aliases: []string{"vers", "vs"},
	Example: `  npm-skills versions react
  npm-skills vs react --latest
  npm-skills versions @nestjs/core -m npm-mirror
  npm-skills vers lodash --proxy http://127.0.0.1:7890`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		packageName := args[0]

		client := resolveClient()
		ctx, cancel := newContext()
		defer cancel()

		if versionsLatestOnly {
			printInfo("Fetching latest version for %s...",
				color.New(color.FgWhite, color.Bold).Sprint(packageName))

			latest, err := client.GetPackageLatestVersion(ctx, packageName)
			if err != nil {
				return fmt.Errorf("failed to get latest version for '%s': %w", packageName, err)
			}

			result := map[string]string{
				"package": packageName,
				"latest":  latest,
			}
			if err := outputJSON(result); err != nil {
				return err
			}
			printSuccess("✓ Latest version of %s is %s",
				color.New(color.FgWhite, color.Bold).Sprint(packageName),
				color.New(color.FgGreen, color.Bold).Sprint(latest))
			return nil
		}

		printInfo("Fetching versions for %s...",
			color.New(color.FgWhite, color.Bold).Sprint(packageName))

		versions, err := client.GetPackageVersions(ctx, packageName)
		if err != nil {
			return fmt.Errorf("failed to get versions for '%s': %w", packageName, err)
		}

		result := map[string]interface{}{
			"package":       packageName,
			"version_count": len(versions),
			"versions":      versions,
		}
		if err := outputJSON(result); err != nil {
			return err
		}
		printSuccess("✓ Found %d versions for %s",
			len(versions),
			color.New(color.FgWhite, color.Bold).Sprint(packageName))
		return nil
	},
}

func init() {
	versionsCmd.Flags().BoolVarP(&versionsLatestOnly, "latest", "L", false,
		"Show only the latest version")
	rootCmd.AddCommand(versionsCmd)
}
