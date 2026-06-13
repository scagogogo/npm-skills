package main

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var bulkStatsPeriod string
var bulkStatsRange bool

var downloadStatsBulkCmd = &cobra.Command{
	Use:   "download-stats-bulk <names...>",
	Short: "Get download stats for multiple packages at once",
	Long: color.New(color.FgCyan).Sprintf("Get download stats for multiple packages at once") + "\n\n" +
		"Batch query download statistics for up to 128 packages in a single request.\n" +
		"Pass package names separated by commas or as separate arguments.\n\n" +
		"Use --range flag to get daily download trends instead of totals.\n\n" +
		color.HiYellowString("Note: ") + "Download stats API (api.npmjs.org) is separate from the registry.\n" +
		"It always queries the official NPM API regardless of --mirror/--registry settings.\n\n" +
		color.HiBlackString("Periods: last-day | last-week (default) | last-month"),
	Aliases: []string{"dsb", "stats-bulk"},
	Example: `  npm-skills download-stats-bulk react,vue,angular
  npm-skills dsb react vue angular -p last-month
  npm-skills stats-bulk react,vue --range
  npm-skills dsb react lodash axios --proxy http://127.0.0.1:7890`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Collect package names from all arguments (support both comma-separated and space-separated)
		var packageNames []string
		for _, arg := range args {
			for _, name := range strings.Split(arg, ",") {
				name = strings.TrimSpace(name)
				if name != "" {
					packageNames = append(packageNames, name)
				}
			}
		}

		if len(packageNames) == 0 {
			return fmt.Errorf("at least one package name is required")
		}
		if len(packageNames) > 128 {
			return fmt.Errorf("maximum 128 packages allowed, got %d", len(packageNames))
		}

		printInfo("Fetching bulk download stats for %d packages (%s)...",
			len(packageNames), bulkStatsPeriod)

		// Download stats API uses api.npmjs.org, not the registry.
		client := resolveDownloadStatsClient()

		ctx, cancel := newContext()
		defer cancel()

		if bulkStatsRange {
			stats, err := client.GetBulkDownloadRangeStats(ctx, packageNames, bulkStatsPeriod)
			if err != nil {
				return fmt.Errorf("failed to get bulk download range stats: %w", err)
			}
			if err := outputJSON(stats); err != nil {
				return err
			}
		} else {
			stats, err := client.GetBulkDownloadStats(ctx, packageNames, bulkStatsPeriod)
			if err != nil {
				return fmt.Errorf("failed to get bulk download stats: %w", err)
			}
			if err := outputJSON(stats); err != nil {
				return err
			}
		}

		printSuccess("✓ Bulk download stats for %d packages retrieved successfully",
			len(packageNames))
		return nil
	},
}

func init() {
	downloadStatsBulkCmd.Flags().StringVarP(&bulkStatsPeriod, "period", "p", "last-week",
		"Time period: last-day, last-week, last-month")
	downloadStatsBulkCmd.Flags().BoolVar(&bulkStatsRange, "range", false,
		"Get daily download trends instead of totals")
	rootCmd.AddCommand(downloadStatsBulkCmd)
}
