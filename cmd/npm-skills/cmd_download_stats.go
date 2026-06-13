package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var statsPeriod string

var downloadStatsCmd = &cobra.Command{
	Use:   "download-stats <name>",
	Short: "Get download statistics for an NPM package",
	Long: color.New(color.FgCyan).Sprintf("Get download statistics for an NPM package") + "\n\n" +
		"Retrieves download counts for the specified time period.\n" +
		color.HiYellowString("Note: ") + "Download stats API (api.npmjs.org) is separate from the registry.\n" +
		"It always queries the official NPM API regardless of --mirror/--registry settings.\n\n" +
		color.HiBlackString("Periods: last-day | last-week (default) | last-month"),
	Aliases: []string{"stats"},
	Example: `  npm-skills download-stats react
  npm-skills download-stats axios -p last-month
  npm-skills stats vue -p last-day
  npm-skills stats react --proxy http://127.0.0.1:7890`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		packageName := args[0]

		printInfo("Fetching download stats for %s (%s)...",
			color.New(color.FgWhite, color.Bold).Sprint(packageName),
			statsPeriod)

		// Download stats API uses api.npmjs.org, not the registry.
		// But we still apply proxy settings from global flags.
		client := resolveDownloadStatsClient()

		ctx, cancel := newContext()
		defer cancel()

		stats, err := client.GetDownloadStats(ctx, packageName, statsPeriod)
		if err != nil {
			return fmt.Errorf("failed to get download stats for '%s': %w", packageName, err)
		}

		if err := outputJSON(stats); err != nil {
			return err
		}
		printSuccess("✓ Download stats for %s retrieved successfully",
			color.New(color.FgWhite, color.Bold).Sprint(packageName))
		return nil
	},
}

func init() {
	downloadStatsCmd.Flags().StringVarP(&statsPeriod, "period", "p", "last-week",
		"Time period: last-day, last-week, last-month")
	rootCmd.AddCommand(downloadStatsCmd)
}