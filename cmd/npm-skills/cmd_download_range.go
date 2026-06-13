package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rangePeriod string

var downloadRangeCmd = &cobra.Command{
	Use:   "download-range <name>",
	Short: "Get daily download trends for an NPM package",
	Long: color.New(color.FgCyan).Sprintf("Get daily download trends for an NPM package") + "\n\n" +
		"Returns daily download counts for the specified period, useful for\n" +
		"visualizing download trends over time.\n\n" +
		color.HiYellowString("Note: ") + "Download stats API (api.npmjs.org) is separate from the registry.\n" +
		"It always queries the official NPM API regardless of --mirror/--registry settings.\n\n" +
		color.HiBlackString("Periods: last-day | last-week (default) | last-month"),
	Aliases: []string{"dr", "range"},
	Example: `  npm-skills download-range react
  npm-skills download-range axios -p last-month
  npm-skills dr vue -p last-day
  npm-skills range react --proxy http://127.0.0.1:7890`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		packageName := args[0]

		printInfo("Fetching daily download range for %s (%s)...",
			color.New(color.FgWhite, color.Bold).Sprint(packageName),
			rangePeriod)

		// Download stats API uses api.npmjs.org, not the registry.
		// But we still apply proxy settings from global flags.
		client := resolveDownloadStatsClient()

		ctx, cancel := newContext()
		defer cancel()

		stats, err := client.GetDownloadRangeStats(ctx, packageName, rangePeriod)
		if err != nil {
			return fmt.Errorf("failed to get download range for '%s': %w", packageName, err)
		}

		if err := outputJSON(stats); err != nil {
			return err
		}
		printSuccess("✓ Daily download trends for %s retrieved (%d days)",
			color.New(color.FgWhite, color.Bold).Sprint(packageName),
			len(stats.Downloads))
		return nil
	},
}

func init() {
	downloadRangeCmd.Flags().StringVarP(&rangePeriod, "period", "p", "last-week",
		"Time period: last-day, last-week, last-month")
	rootCmd.AddCommand(downloadRangeCmd)
}