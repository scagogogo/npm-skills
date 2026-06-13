package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/scagogogo/npm-skills/pkg/models"
	"github.com/spf13/cobra"
)

var couchdbCmd = &cobra.Command{
	Use:   "couchdb",
	Short: "CouchDB advanced query operations",
	Long: color.New(color.FgCyan).Sprintf("CouchDB advanced query operations") + "\n\n" +
		"Subcommands: changes, all-docs, view\n\n" +
		"Low-level CouchDB API access for mirroring, incremental sync,\n" +
		"and advanced data queries. Most users do not need these commands.",
	Example: `  npm-skills couchdb changes
  npm-skills couchdb all-docs my-registry`,
}

var changesSince string
var changesLimit int
var changesIncludeDocs bool

var changesCmd = &cobra.Command{
	Use:   "changes",
	Short: "Get registry changes feed",
	Long: color.New(color.FgCyan).Sprintf("Get registry changes feed") + "\n\n" +
		"Returns the CouchDB _changes feed, listing all document modifications.\n" +
		"Primarily used for mirror synchronization and incremental data fetching.\n\n" +
		"Use --since to resume from a previous last_seq value.",
	Example: `  npm-skills couchdb changes
  npm-skills couchdb changes --limit 10
  npm-skills couchdb changes --since 12345 --limit 50
  npm-skills couchdb changes --include-docs`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		printInfo("Fetching changes feed from %s...", currentMirrorLabel())
		client := resolveClient()
		ctx, cancel := newContext()
		defer cancel()

		opts := models.ChangesOptions{
			Since:       changesSince,
			Limit:       changesLimit,
			IncludeDocs: changesIncludeDocs,
		}

		result, err := client.GetChanges(ctx, opts)
		if err != nil {
			return fmt.Errorf("failed to get changes: %w", err)
		}

		if err := outputJSON(result); err != nil {
			return err
		}
		printSuccess("✓ %d changes (last_seq=%s, pending=%d)",
			len(result.Results), result.LastSeq, result.Pending)
		return nil
	},
}

var allDocsStartKey string
var allDocsEndKey string
var allDocsLimit int
var allDocsSkip int
var allDocsIncludeDocs bool
var allDocsDescending bool

var allDocsCmd = &cobra.Command{
	Use:   "all-docs",
	Short: "List all document IDs in the registry",
	Long: color.New(color.FgCyan).Sprintf("List all document IDs in the registry") + "\n\n" +
		"Returns the CouchDB _all_docs result with all document IDs and revision info.\n" +
		"Optionally include full document content with --include-docs (may be very large).\n\n" +
		"Use --start-key and --end-key for range queries.",
	Example: `  npm-skills couchdb all-docs --limit 20
  npm-skills couchdb all-docs --start-key "@nestjs" --end-key "@nestt" --limit 50
  npm-skills couchdb all-docs --include-docs --limit 5`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		printInfo("Fetching all-docs from %s...", currentMirrorLabel())
		client := resolveClient()
		ctx, cancel := newContext()
		defer cancel()

		opts := models.AllDocsOptions{
			StartKey:    allDocsStartKey,
			EndKey:      allDocsEndKey,
			Limit:       allDocsLimit,
			Skip:        allDocsSkip,
			IncludeDocs: allDocsIncludeDocs,
			Descending:  allDocsDescending,
		}

		result, err := client.GetAllDocs(ctx, opts)
		if err != nil {
			return fmt.Errorf("failed to get all-docs: %w", err)
		}

		if err := outputJSON(result); err != nil {
			return err
		}
		printSuccess("✓ %d rows (total=%d, offset=%d)",
			len(result.Rows), result.TotalRows, result.Offset)
		return nil
	},
}

var viewName string
var viewKey string
var viewStartKey string
var viewEndKey string
var viewLimit int
var viewSkip int
var viewGroup bool
var viewGroupLevel int
var viewDescending bool

var viewCmd = &cobra.Command{
	Use:   "view <view-name>",
	Short: "Query a CouchDB view",
	Long: color.New(color.FgCyan).Sprintf("Query a CouchDB view") + "\n\n" +
		"Queries a CouchDB view on the NPM registry.\n" +
		"Common views: starredByUser, starredByPackage, byKeyword, byUser.\n\n" +
		color.HiYellowString("Note: ") + "Not all registries support all views.",
	Example: `  npm-skills couchdb view starredByUser --key '"username"'
  npm-skills couchdb view starredByPackage --key '"react"'
  npm-skills couchdb view byKeyword --key '"http"' --limit 20`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		viewName := args[0]
		printInfo("Querying view %s on %s...", viewName, currentMirrorLabel())
		client := resolveClient()
		ctx, cancel := newContext()
		defer cancel()

		opts := models.ViewOptions{
			Key:        viewKey,
			StartKey:   viewStartKey,
			EndKey:     viewEndKey,
			Limit:      viewLimit,
			Skip:       viewSkip,
			Group:      viewGroup,
			GroupLevel: viewGroupLevel,
			Descending: viewDescending,
		}

		result, err := client.GetView(ctx, viewName, opts)
		if err != nil {
			return fmt.Errorf("failed to get view: %w", err)
		}

		if err := outputJSON(result); err != nil {
			return err
		}
		printSuccess("✓ %d rows (total=%d, offset=%d)",
			len(result.Rows), result.TotalRows, result.Offset)
		return nil
	},
}

func init() {
	changesCmd.Flags().StringVar(&changesSince, "since", "", "Start sequence (resume from last_seq)")
	changesCmd.Flags().IntVar(&changesLimit, "limit", 0, "Limit number of results")
	changesCmd.Flags().BoolVar(&changesIncludeDocs, "include-docs", false, "Include full document content")

	allDocsCmd.Flags().StringVar(&allDocsStartKey, "start-key", "", "Start key for range query")
	allDocsCmd.Flags().StringVar(&allDocsEndKey, "end-key", "", "End key for range query")
	allDocsCmd.Flags().IntVar(&allDocsLimit, "limit", 0, "Limit number of results")
	allDocsCmd.Flags().IntVar(&allDocsSkip, "skip", 0, "Skip number of results")
	allDocsCmd.Flags().BoolVar(&allDocsIncludeDocs, "include-docs", false, "Include full document content (may be large)")
	allDocsCmd.Flags().BoolVar(&allDocsDescending, "descending", false, "Reverse order")

	viewCmd.Flags().StringVar(&viewKey, "key", "", "Exact key to query")
	viewCmd.Flags().StringVar(&viewStartKey, "start-key", "", "Start key for range query")
	viewCmd.Flags().StringVar(&viewEndKey, "end-key", "", "End key for range query")
	viewCmd.Flags().IntVar(&viewLimit, "limit", 0, "Limit number of results")
	viewCmd.Flags().IntVar(&viewSkip, "skip", 0, "Skip number of results")
	viewCmd.Flags().BoolVar(&viewGroup, "group", false, "Group results")
	viewCmd.Flags().IntVar(&viewGroupLevel, "group-level", 0, "Group level for hierarchical grouping")
	viewCmd.Flags().BoolVar(&viewDescending, "descending", false, "Reverse order")

	couchdbCmd.AddCommand(changesCmd)
	couchdbCmd.AddCommand(allDocsCmd)
	couchdbCmd.AddCommand(viewCmd)
	rootCmd.AddCommand(couchdbCmd)
}
