package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/scagogogo/npm-skills/pkg/models"
	"github.com/scagogogo/npm-skills/pkg/registry"
	"github.com/spf13/cobra"
)

var (
	searchLimit       int
	searchFrom        int
	searchQuality     float64
	searchPopularity  float64
	searchMaintenance float64
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search NPM packages by keyword",
	Long: color.New(color.FgCyan).Sprintf("Search NPM packages by keyword") + "\n\n" +
		"Searches the NPM registry for packages matching the given query.\n" +
		"Returns package name, version, description, and relevance score.\n\n" +
		"Advanced options:\n" +
		"  • --from: pagination offset (for browsing beyond first page)\n" +
		"  • --quality / --popularity / --maintenance: score weights (0.0-1.0)\n\n" +
		color.HiBlackString("Mirror: %s", mirrorNames()) + " (via --mirror or --registry flag)",
	Aliases: []string{"s"},
	Example: `  npm-skills search "http client"
  npm-skills search react -l 5
  npm-skills search "vue component" -m npm-mirror
  npm-skills s axios --from 20 -l 10
  npm-skills s "http client" --popularity 1.0 --quality 0.0 --maintenance 0.0
  npm-skills s react --proxy http://127.0.0.1:7890`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]

		client := resolveClient()
		ctx, cancel := newContext()
		defer cancel()

		// Check if advanced options are set
		hasAdvanced := searchFrom > 0 || searchQuality > 0 || searchPopularity > 0 || searchMaintenance > 0

		var result *models.SearchResult
		var err error

		if hasAdvanced {
			printInfo("Searching for %s (limit: %d, from: %d) on %s...",
				color.New(color.FgWhite, color.Bold).Sprintf("\"%s\"", query),
				searchLimit, searchFrom,
				currentMirrorLabel())

			result, err = client.SearchPackagesWithOptions(ctx, query, registry.SearchOptions{
				From:        searchFrom,
				Size:        searchLimit,
				Quality:     searchQuality,
				Popularity:  searchPopularity,
				Maintenance: searchMaintenance,
			})
		} else {
			printInfo("Searching for %s (limit: %d) on %s...",
				color.New(color.FgWhite, color.Bold).Sprintf("\"%s\"", query),
				searchLimit,
				currentMirrorLabel())

			result, err = client.SearchPackages(ctx, query, searchLimit)
		}

		if err != nil {
			return fmt.Errorf("failed to search for '%s': %w", query, err)
		}

		if err := outputJSON(result); err != nil {
			return err
		}
		printSuccess("✓ Found %d results for %s",
			result.Total,
			color.New(color.FgWhite, color.Bold).Sprintf("\"%s\"", query))
		return nil
	},
}

func init() {
	searchCmd.Flags().IntVarP(&searchLimit, "limit", "l", 20, "Maximum number of results")
	searchCmd.Flags().IntVar(&searchFrom, "from", 0, "Pagination offset (0-based)")
	searchCmd.Flags().Float64Var(&searchQuality, "quality", 0, "Quality weight (0.0-1.0)")
	searchCmd.Flags().Float64Var(&searchPopularity, "popularity", 0, "Popularity weight (0.0-1.0)")
	searchCmd.Flags().Float64Var(&searchMaintenance, "maintenance", 0, "Maintenance weight (0.0-1.0)")
	rootCmd.AddCommand(searchCmd)
}
