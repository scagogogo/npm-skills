package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var distTagsAbbreviated bool

var distTagsCmd = &cobra.Command{
	Use:   "dist-tags",
	Short: "Manage distribution tags for an NPM package",
	Long: color.New(color.FgCyan).Sprintf("Manage distribution tags (dist-tags) for an NPM package") + "\n\n" +
		"Dist-tags are NPM's version aliases. The most common are:\n" +
		"  • latest — latest stable release\n" +
		"  • next   — upcoming version\n" +
		"  • beta   — beta release\n\n" +
		"Subcommands:\n" +
		"  get <package>          — get all dist-tags (default)\n" +
		"  set <package> <tag>    — set a dist-tag (requires --version and --token)\n" +
		"  delete <package> <tag> — delete a dist-tag (requires --token)\n\n" +
		color.HiBlackString("Mirror: %s", mirrorNames()) + " (via --mirror or --registry flag)",
	Aliases: []string{"tags", "dt"},
	Example: `  npm-skills dist-tags get react
  npm-skills dist-tags set react next --version 2.0.0-rc.1 --token npm_xxxxx
  npm-skills dist-tags delete react beta --token npm_xxxxx
  npm-skills dist-tags get @nestjs/core -m npm-mirror`,
}

var distTagsGetCmd = &cobra.Command{
	Use:   "get <package>",
	Short: "Get distribution tags for an NPM package",
	Long: color.New(color.FgCyan).Sprintf("Get distribution tags (dist-tags) for an NPM package") + "\n\n" +
		"By default uses the full package metadata endpoint.\n" +
		"Use --abbreviated flag for a faster, lightweight query.",
	Aliases: []string{"list", "ls"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		packageName := args[0]

		printInfo("Fetching dist-tags for %s...",
			color.New(color.FgWhite, color.Bold).Sprint(packageName))

		client := resolveClient()
		ctx, cancel := newContext()
		defer cancel()

		var tags map[string]string
		var err error

		if distTagsAbbreviated {
			tags, err = client.GetDistTagsAbbreviated(ctx, packageName)
		} else {
			tags, err = client.GetDistTags(ctx, packageName)
		}

		if err != nil {
			return fmt.Errorf("failed to get dist-tags for '%s': %w", packageName, err)
		}

		if err := outputJSON(tags); err != nil {
			return err
		}
		printSuccess("✓ Found %d dist-tags for %s",
			len(tags),
			color.New(color.FgWhite, color.Bold).Sprint(packageName))
		return nil
	},
}

var distTagSetVersion string

var distTagsSetCmd = &cobra.Command{
	Use:   "set <package> <tag>",
	Short: "Set a distribution tag (requires --token)",
	Long: color.New(color.FgCyan).Sprintf("Set a distribution tag for an NPM package") + "\n\n" +
		"Sets the specified tag to point to the given version.\n" +
		"Requires authentication token.\n\n" +
		color.HiYellowString("Note: ") + "If the tag already exists, it will be updated.",
	Aliases: []string{"add"},
	Example: `  npm-skills dist-tags set react next --version 2.0.0-rc.1 --token npm_xxxxx
  npm-skills dist-tags set my-pkg stable --version 1.0.0 -t npm_xxxxx
  npm-skills dist-tags add @myorg/mypkg beta --version 3.0.0-beta.1 -t npm_xxxxx`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(); err != nil {
			return err
		}
		if distTagSetVersion == "" {
			return fmt.Errorf("--version is required")
		}

		packageName := args[0]
		tag := args[1]

		printInfo("Setting dist-tag %s=%s on %s...",
			color.New(color.FgYellow).Sprint(tag),
			color.New(color.FgGreen).Sprint(distTagSetVersion),
			color.New(color.FgWhite, color.Bold).Sprint(packageName))

		client := resolveClientWithToken()
		ctx, cancel := newContext()
		defer cancel()

		err := client.SetDistTag(ctx, packageName, tag, distTagSetVersion)
		if err != nil {
			return fmt.Errorf("failed to set dist-tag: %w", err)
		}

		result := map[string]string{
			"package": packageName,
			"tag":     tag,
			"version": distTagSetVersion,
			"status":  "updated",
		}
		if err := outputJSON(result); err != nil {
			return err
		}
		printSuccess("✓ Set dist-tag %s=%s on %s",
			color.New(color.FgYellow).Sprint(tag),
			distTagSetVersion,
			color.New(color.FgWhite, color.Bold).Sprint(packageName))
		return nil
	},
}

var distTagsDeleteCmd = &cobra.Command{
	Use:   "delete <package> <tag>",
	Short: "Delete a distribution tag (requires --token)",
	Long: color.New(color.FgCyan).Sprintf("Delete a distribution tag from an NPM package") + "\n\n" +
		"Removes the specified tag from the package.\n" +
		"Requires authentication token.\n\n" +
		color.HiRedString("WARNING: ") + "Deleting the 'latest' tag may cause issues.",
	Aliases: []string{"rm", "remove"},
	Example: `  npm-skills dist-tags delete react beta --token npm_xxxxx
  npm-skills dist-tags rm my-pkg next -t npm_xxxxx`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(); err != nil {
			return err
		}

		packageName := args[0]
		tag := args[1]

		printInfo("Deleting dist-tag %s from %s...",
			color.New(color.FgRed).Sprint(tag),
			color.New(color.FgWhite, color.Bold).Sprint(packageName))

		client := resolveClientWithToken()
		ctx, cancel := newContext()
		defer cancel()

		err := client.DeleteDistTag(ctx, packageName, tag)
		if err != nil {
			return fmt.Errorf("failed to delete dist-tag: %w", err)
		}

		result := map[string]string{
			"package": packageName,
			"tag":     tag,
			"status":  "deleted",
		}
		if err := outputJSON(result); err != nil {
			return err
		}
		printSuccess("✓ Deleted dist-tag %s from %s",
			color.New(color.FgRed).Sprint(tag),
			color.New(color.FgWhite, color.Bold).Sprint(packageName))
		return nil
	},
}

func init() {
	distTagsGetCmd.Flags().BoolVarP(&distTagsAbbreviated, "abbreviated", "a", false,
		"Use lightweight dist-tags endpoint (faster, less data)")
	distTagsSetCmd.Flags().StringVar(&distTagSetVersion, "version", "",
		"Version to point the tag to (required)")

	distTagsCmd.AddCommand(distTagsGetCmd)
	distTagsCmd.AddCommand(distTagsSetCmd)
	distTagsCmd.AddCommand(distTagsDeleteCmd)
	rootCmd.AddCommand(distTagsCmd)
}