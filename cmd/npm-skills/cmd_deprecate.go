package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var deprecateMessage string

var deprecateCmd = &cobra.Command{
	Use:   "deprecate <package> <version>",
	Short: "Deprecate a specific version of an NPM package",
	Long: color.New(color.FgCyan).Sprintf("Deprecate a specific version of an NPM package") + "\n\n" +
		"Marks a version as deprecated. Users will see a warning when installing.\n" +
		"Requires authentication token.\n\n" +
		color.HiYellowString("Note: ") + "Prefer deprecation over unpublish. It's safer and doesn't remove the version.",
	Aliases: []string{"dep"},
	Example: `  npm-skills deprecate my-package 1.0.0 --message "Use v2.0.0 instead" --token npm_xxxxx
  npm-skills dep my-package 1.0.0 -m "Security vulnerability, upgrade to 1.0.1" -t npm_xxxxx`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(); err != nil {
			return err
		}

		packageName := args[0]
		version := args[1]

		if deprecateMessage == "" {
			return fmt.Errorf("--message is required")
		}

		printInfo("Deprecating %s@%s on %s...", packageName, version, currentMirrorLabel())

		client := resolveClientWithToken()
		ctx, cancel := newContext()
		defer cancel()

		err := client.DeprecateVersion(ctx, packageName, version, deprecateMessage)
		if err != nil {
			return fmt.Errorf("failed to deprecate: %w", err)
		}

		result := map[string]string{
			"package": packageName,
			"version": version,
			"message": deprecateMessage,
			"status":  "deprecated",
		}
		if err := outputJSON(result); err != nil {
			return err
		}
		printSuccess("✓ Deprecated %s@%s: %s",
			color.New(color.FgWhite, color.Bold).Sprint(packageName), version, deprecateMessage)
		return nil
	},
}

func init() {
	deprecateCmd.Flags().StringVarP(&deprecateMessage, "message", "M", "", "Deprecation message (required)")
	rootCmd.AddCommand(deprecateCmd)
}
