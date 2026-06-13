package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var downloadStatsDateCmd = &cobra.Command{
	Use:   "download-stats-date <name>",
	Short: "Get download stats for a custom date range",
	Long: color.New(color.FgCyan).Sprintf("Get download stats for a custom date range") + "\n\n" +
		"Retrieves download counts for a specific date range (YYYY-MM-DD format).\n" +
		"Use this instead of download-stats when you need stats for a specific period\n" +
		"that isn't covered by last-day/last-week/last-month.\n\n" +
		color.HiYellowString("Note: ") + "Download stats API (api.npmjs.org) is separate from the registry.\n" +
		"It always queries the official NPM API regardless of --mirror/--registry settings.",
	Aliases: []string{"dsd", "stats-date"},
	Example: `  npm-skills download-stats-date react --start 2024-01-01 --end 2024-01-31
  npm-skills dsd axios --start 2024-06-01 --end 2024-06-30
  npm-skills stats-date vue --start 2024-01-01 --end 2024-12-31 --proxy http://127.0.0.1:7890`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		packageName := args[0]

		if statsStartDate == "" || statsEndDate == "" {
			return fmt.Errorf("both --start and --end dates are required (format: YYYY-MM-DD)")
		}

		printInfo("Fetching download stats for %s (%s to %s)...",
			color.New(color.FgWhite, color.Bold).Sprint(packageName),
			statsStartDate, statsEndDate)

		// Download stats API uses api.npmjs.org, not the registry.
		client := resolveDownloadStatsClient()

		ctx, cancel := newContext()
		defer cancel()

		stats, err := client.GetDownloadStatsByDateRange(ctx, packageName, statsStartDate, statsEndDate)
		if err != nil {
			return fmt.Errorf("failed to get download stats for '%s': %w", packageName, err)
		}

		if err := outputJSON(stats); err != nil {
			return err
		}
		printSuccess("✓ Download stats for %s (%s to %s) retrieved successfully",
			color.New(color.FgWhite, color.Bold).Sprint(packageName),
			statsStartDate, statsEndDate)
		return nil
	},
}

var statsStartDate string
var statsEndDate string

func init() {
	downloadStatsDateCmd.Flags().StringVar(&statsStartDate, "start", "", "Start date (YYYY-MM-DD, required)")
	downloadStatsDateCmd.Flags().StringVar(&statsEndDate, "end", "", "End date (YYYY-MM-DD, required)")
	rootCmd.AddCommand(downloadStatsDateCmd)
}
