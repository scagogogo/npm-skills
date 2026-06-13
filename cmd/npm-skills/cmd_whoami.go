package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Check current NPM authentication status",
	Long: color.New(color.FgCyan).Sprintf("Check current NPM authentication status") + "\n\n" +
		"Verifies the authentication token and returns the logged-in username.\n" +
		"Requires a token set via --token flag or NPM_TOKEN environment variable.\n\n" +
		color.HiYellowString("Note: ") + "This command always queries the configured registry's /-/whoami endpoint.",
	Aliases: []string{"me"},
	Example: `  npm-skills whoami --token npm_xxxxx
  NPM_TOKEN=npm_xxxxx npm-skills whoami
  npm-skills whoami --token npm_xxxxx -m npm-mirror
  npm-skills me --token npm_xxxxx --proxy http://127.0.0.1:7890`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(); err != nil {
			return err
		}

		printInfo("Checking authentication on %s...", currentMirrorLabel())

		// Create client with token
		client := resolveClientWithToken()

		ctx, cancel := newContext()
		defer cancel()

		username, err := client.WhoAmI(ctx)
		if err != nil {
			return fmt.Errorf("authentication check failed: %w", err)
		}

		result := map[string]string{
			"username": username,
			"status":   "authenticated",
		}
		if err := outputJSON(result); err != nil {
			return err
		}
		printSuccess("✓ Authenticated as %s",
			color.New(color.FgWhite, color.Bold).Sprint(username))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}
