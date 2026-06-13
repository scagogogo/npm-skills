package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/scagogogo/npm-skills/pkg/models"
	"github.com/spf13/cobra"
)

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Manage NPM access tokens",
	Long: color.New(color.FgCyan).Sprintf("Manage NPM access tokens") + "\n\n" +
		"Subcommands: list, get, create, delete\n\n" +
		color.HiYellowString("Note: ") + "All token operations require authentication.",
	Example: `  npm-skills token list -t npm_xxxxx
  npm-skills token create --password mypass -t npm_xxxxx`,
}

var tokenListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all NPM access tokens",
	Long: color.New(color.FgCyan).Sprintf("List all NPM access tokens") + "\n\n" +
		"Requires authentication token.",
	Aliases: []string{"ls"},
	Example: `  npm-skills token list -t npm_xxxxx`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(); err != nil {
			return err
		}

		printInfo("Listing tokens on %s...", currentMirrorLabel())
		client := resolveClientWithToken()
		ctx, cancel := newContext()
		defer cancel()

		tokens, err := client.ListTokens(ctx)
		if err != nil {
			return fmt.Errorf("failed to list tokens: %w", err)
		}

		if err := outputJSON(tokens); err != nil {
			return err
		}
		printSuccess("✓ Found %d tokens", len(tokens))
		return nil
	},
}

var tokenGetCmd = &cobra.Command{
	Use:   "get <token-id>",
	Short: "Get details of a specific NPM access token",
	Long: color.New(color.FgCyan).Sprintf("Get details of a specific NPM access token") + "\n\n" +
		"Requires authentication token.",
	Example: `  npm-skills token get abc123 -t npm_xxxxx`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(); err != nil {
			return err
		}

		tokenID := args[0]
		printInfo("Getting token %s from %s...", tokenID, currentMirrorLabel())
		client := resolveClientWithToken()
		ctx, cancel := newContext()
		defer cancel()

		token, err := client.GetToken(ctx, tokenID)
		if err != nil {
			return fmt.Errorf("failed to get token: %w", err)
		}

		if err := outputJSON(token); err != nil {
			return err
		}
		printSuccess("✓ Token: %s (readonly=%v)", token.ID, token.Readonly)
		return nil
	},
}

var tokenCreateReadonly bool
var tokenCreatePassword string
var tokenCreateCIDR []string

var tokenCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new NPM access token",
	Long: color.New(color.FgCyan).Sprintf("Create a new NPM access token") + "\n\n" +
		"Requires authentication token and your current password.\n" +
		"New tokens can be set to read-only and restricted to specific IP ranges (CIDR).",
	Aliases: []string{"new"},
	Example: `  npm-skills token create --password mypass -t npm_xxxxx
  npm-skills token create --password mypass --readonly -t npm_xxxxx
  npm-skills token create --password mypass --cidr 192.168.1.0/24 -t npm_xxxxx`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(); err != nil {
			return err
		}

		if tokenCreatePassword == "" {
			return fmt.Errorf("--password is required")
		}

		printInfo("Creating token on %s...", currentMirrorLabel())
		client := resolveClientWithToken()
		ctx, cancel := newContext()
		defer cancel()

		opts := &models.TokenCreation{
			Password: tokenCreatePassword,
			Readonly: tokenCreateReadonly,
			CIDR:     tokenCreateCIDR,
		}

		token, err := client.CreateToken(ctx, opts)
		if err != nil {
			return fmt.Errorf("failed to create token: %w", err)
		}

		if err := outputJSON(token); err != nil {
			return err
		}
		printSuccess("✓ Token created: %s (readonly=%v)", token.ID, token.Readonly)
		return nil
	},
}

var tokenDeleteCmd = &cobra.Command{
	Use:   "delete <token-id>",
	Short: "Delete (revoke) an NPM access token",
	Long: color.New(color.FgCyan).Sprintf("Delete (revoke) an NPM access token") + "\n\n" +
		"Requires authentication token. The token will be immediately invalidated.\n\n" +
		color.HiYellowString("Note: ") + "You cannot delete the token you are currently using.",
	Aliases: []string{"rm", "revoke"},
	Example: `  npm-skills token delete abc123 -t npm_xxxxx`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(); err != nil {
			return err
		}

		tokenID := args[0]
		printInfo("Deleting token %s on %s...", tokenID, currentMirrorLabel())
		client := resolveClientWithToken()
		ctx, cancel := newContext()
		defer cancel()

		err := client.DeleteToken(ctx, tokenID)
		if err != nil {
			return fmt.Errorf("failed to delete token: %w", err)
		}

		if err := outputJSON(map[string]string{"id": tokenID, "status": "deleted"}); err != nil {
			return err
		}
		printSuccess("✓ Token %s deleted", tokenID)
		return nil
	},
}

func init() {
	tokenCreateCmd.Flags().StringVar(&tokenCreatePassword, "password", "", "Current user password (required)")
	tokenCreateCmd.Flags().BoolVar(&tokenCreateReadonly, "readonly", false, "Create a read-only token")
	tokenCreateCmd.Flags().StringArrayVar(&tokenCreateCIDR, "cidr", nil, "IP whitelist CIDR ranges (can be specified multiple times)")

	tokenCmd.AddCommand(tokenListCmd)
	tokenCmd.AddCommand(tokenGetCmd)
	tokenCmd.AddCommand(tokenCreateCmd)
	tokenCmd.AddCommand(tokenDeleteCmd)
	rootCmd.AddCommand(tokenCmd)
}