package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var starCmd = &cobra.Command{
	Use:   "star",
	Short: "Star/unstar NPM packages",
	Long: color.New(color.FgCyan).Sprintf("Star/unstar NPM packages") + "\n\n" +
		"Subcommands: add, remove, list, stargazers",
	Example: `  npm-skills star add react -t npm_xxxxx
  npm-skills star list myuser`,
}

var starAddCmd = &cobra.Command{
	Use:   "add <package>",
	Short: "Star (favorite) an NPM package",
	Long: color.New(color.FgCyan).Sprintf("Star (favorite) an NPM package") + "\n\n" +
		"Adds the package to your starred packages list.\n" +
		"Requires authentication token.",
	Aliases: []string{"star"},
	Example: `  npm-skills star add react -t npm_xxxxx`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(); err != nil {
			return err
		}

		packageName := args[0]
		printInfo("Starring %s on %s...", packageName, currentMirrorLabel())
		client := resolveClientWithToken()
		ctx, cancel := newContext()
		defer cancel()

		err := client.StarPackage(ctx, packageName)
		if err != nil {
			return fmt.Errorf("failed to star package: %w", err)
		}

		if err := outputJSON(map[string]string{"package": packageName, "status": "starred"}); err != nil {
			return err
		}
		printSuccess("✓ Starred %s", color.New(color.FgWhite, color.Bold).Sprint(packageName))
		return nil
	},
}

var starRemoveCmd = &cobra.Command{
	Use:   "remove <package>",
	Short: "Unstar (unfavorite) an NPM package",
	Long: color.New(color.FgCyan).Sprintf("Unstar (unfavorite) an NPM package") + "\n\n" +
		"Removes the package from your starred packages list.\n" +
		"Requires authentication token.",
	Aliases: []string{"unstar"},
	Example: `  npm-skills star remove react -t npm_xxxxx`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(); err != nil {
			return err
		}

		packageName := args[0]
		printInfo("Unstarring %s on %s...", packageName, currentMirrorLabel())
		client := resolveClientWithToken()
		ctx, cancel := newContext()
		defer cancel()

		err := client.UnstarPackage(ctx, packageName)
		if err != nil {
			return fmt.Errorf("failed to unstar package: %w", err)
		}

		if err := outputJSON(map[string]string{"package": packageName, "status": "unstarred"}); err != nil {
			return err
		}
		printSuccess("✓ Unstarred %s", color.New(color.FgWhite, color.Bold).Sprint(packageName))
		return nil
	},
}

var starListCmd = &cobra.Command{
	Use:   "list <username>",
	Short: "List packages starred by a user",
	Long: color.New(color.FgCyan).Sprintf("List packages starred by a user") + "\n\n" +
		"Returns the package names that the specified user has starred.",
	Aliases: []string{"ls"},
	Example: `  npm-skills star list myuser`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		username := args[0]
		printInfo("Getting starred packages for %s from %s...", username, currentMirrorLabel())
		client := resolveClient()
		ctx, cancel := newContext()
		defer cancel()

		packages, err := client.GetStarredByUser(ctx, username)
		if err != nil {
			return fmt.Errorf("failed to get starred packages: %w", err)
		}

		if err := outputJSON(packages); err != nil {
			return err
		}
		printSuccess("✓ %d starred packages for %s", len(packages), username)
		return nil
	},
}

var stargazersCmd = &cobra.Command{
	Use:   "stargazers <package>",
	Short: "List users who starred a package",
	Long: color.New(color.FgCyan).Sprintf("List users who starred a package") + "\n\n" +
		"Returns usernames of all users who have starred the specified package.",
	Aliases: []string{"sg"},
	Example: `  npm-skills star stargazers react`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		packageName := args[0]
		printInfo("Getting stargazers for %s from %s...", packageName, currentMirrorLabel())
		client := resolveClient()
		ctx, cancel := newContext()
		defer cancel()

		users, err := client.GetStarredByPackage(ctx, packageName)
		if err != nil {
			return fmt.Errorf("failed to get stargazers: %w", err)
		}

		if err := outputJSON(users); err != nil {
			return err
		}
		printSuccess("✓ %d stargazers for %s", len(users), packageName)
		return nil
	},
}

func init() {
	starCmd.AddCommand(starAddCmd)
	starCmd.AddCommand(starRemoveCmd)
	starCmd.AddCommand(starListCmd)
	starCmd.AddCommand(stargazersCmd)
	rootCmd.AddCommand(starCmd)
}