package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/scagogogo/npm-skills/pkg/models"
	"github.com/spf13/cobra"
)

var hookCmd = &cobra.Command{
	Use:   "hook",
	Short: "Manage NPM webhooks",
	Long: color.New(color.FgCyan).Sprintf("Manage NPM webhooks") + "\n\n" +
		"Subcommands: list, get, create, update, delete\n\n" +
		"Webhooks notify external services when packages are published or updated.\n" +
		"All hook operations require authentication token.",
	Example: `  npm-skills hook list -t npm_xxxxx
  npm-skills hook create --name my-hook --endpoint https://example.com/hook -t npm_xxxxx`,
}

var hookListPackage string
var hookListPage int
var hookListPerPage int

var hookListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all webhooks",
	Long: color.New(color.FgCyan).Sprintf("List all webhooks") + "\n\n" +
		"Requires authentication token. Optionally filter by package name.",
	Aliases: []string{"ls"},
	Example: `  npm-skills hook list -t npm_xxxxx
  npm-skills hook list --package react -t npm_xxxxx`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(); err != nil {
			return err
		}

		printInfo("Listing hooks on %s...", currentMirrorLabel())
		client := resolveClientWithToken()
		ctx, cancel := newContext()
		defer cancel()

		opts := models.HookListOptions{
			Package: hookListPackage,
			Page:    hookListPage,
			PerPage: hookListPerPage,
		}

		hooks, err := client.ListHooks(ctx, opts)
		if err != nil {
			return fmt.Errorf("failed to list hooks: %w", err)
		}

		if err := outputJSON(hooks); err != nil {
			return err
		}
		printSuccess("✓ Found %d hooks", len(hooks))
		return nil
	},
}

var hookGetCmd = &cobra.Command{
	Use:   "get <hook-id>",
	Short: "Get details of a specific webhook",
	Long: color.New(color.FgCyan).Sprintf("Get details of a specific webhook") + "\n\n" +
		"Requires authentication token.",
	Example: `  npm-skills hook get hook-id-here -t npm_xxxxx`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(); err != nil {
			return err
		}

		hookID := args[0]
		printInfo("Getting hook %s from %s...", hookID, currentMirrorLabel())
		client := resolveClientWithToken()
		ctx, cancel := newContext()
		defer cancel()

		hook, err := client.GetHook(ctx, hookID)
		if err != nil {
			return fmt.Errorf("failed to get hook: %w", err)
		}

		if err := outputJSON(hook); err != nil {
			return err
		}
		printSuccess("✓ Hook: %s -> %s", hook.Name, hook.Endpoint)
		return nil
	},
}

var hookCreateName string
var hookCreateEndpoint string
var hookCreateSecret string
var hookCreatePackage string
var hookCreateEvents []string

var hookCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new webhook",
	Long: color.New(color.FgCyan).Sprintf("Create a new webhook") + "\n\n" +
		"Requires authentication token. The webhook will send HTTP POST requests\n" +
		"to the specified endpoint when the target package is published or updated.",
	Example: `  npm-skills hook create --name my-hook --endpoint https://example.com/webhook --package react -t npm_xxxxx
  npm-skills hook create --name ci-hook --endpoint https://ci.example.com/npm-hook --secret mysecret --package my-pkg -t npm_xxxxx`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(); err != nil {
			return err
		}

		if hookCreateName == "" || hookCreateEndpoint == "" {
			return fmt.Errorf("--name and --endpoint are required")
		}

		printInfo("Creating hook on %s...", currentMirrorLabel())
		client := resolveClientWithToken()
		ctx, cancel := newContext()
		defer cancel()

		hook := &models.HookCreation{
			Name:     hookCreateName,
			Endpoint: hookCreateEndpoint,
			Secret:   hookCreateSecret,
			Package:  hookCreatePackage,
			Events:   hookCreateEvents,
			Active:   true,
		}

		result, err := client.CreateHook(ctx, hook)
		if err != nil {
			return fmt.Errorf("failed to create hook: %w", err)
		}

		if err := outputJSON(result); err != nil {
			return err
		}
		printSuccess("✓ Created hook %s -> %s", result.Name, result.Endpoint)
		return nil
	},
}

var hookUpdateEndpoint string
var hookUpdateSecret string
var hookUpdateEvents []string
var hookUpdateActive bool
var hookUpdateActiveSet bool

var hookUpdateCmd = &cobra.Command{
	Use:   "update <hook-id>",
	Short: "Update a webhook",
	Long: color.New(color.FgCyan).Sprintf("Update a webhook") + "\n\n" +
		"Requires authentication token. Only specified fields will be updated.",
	Aliases: []string{"edit"},
	Example: `  npm-skills hook update hook-id-here --endpoint https://new.example.com/webhook -t npm_xxxxx
  npm-skills hook update hook-id-here --secret newsecret -t npm_xxxxx`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(); err != nil {
			return err
		}

		hookID := args[0]
		printInfo("Updating hook %s on %s...", hookID, currentMirrorLabel())
		client := resolveClientWithToken()
		ctx, cancel := newContext()
		defer cancel()

		hook := &models.HookUpdate{
			Endpoint: hookUpdateEndpoint,
			Secret:   hookUpdateSecret,
			Events:   hookUpdateEvents,
		}
		if hookUpdateActiveSet {
			hook.Active = &hookUpdateActive
		}

		result, err := client.UpdateHook(ctx, hookID, hook)
		if err != nil {
			return fmt.Errorf("failed to update hook: %w", err)
		}

		if err := outputJSON(result); err != nil {
			return err
		}
		printSuccess("✓ Updated hook %s", hookID)
		return nil
	},
}

var hookDeleteCmd = &cobra.Command{
	Use:   "delete <hook-id>",
	Short: "Delete a webhook",
	Long: color.New(color.FgCyan).Sprintf("Delete a webhook") + "\n\n" +
		"Requires authentication token. The webhook will be permanently removed\n" +
		"and no further notifications will be sent.",
	Aliases: []string{"rm"},
	Example: `  npm-skills hook delete hook-id-here -t npm_xxxxx`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(); err != nil {
			return err
		}

		hookID := args[0]
		printInfo("Deleting hook %s on %s...", hookID, currentMirrorLabel())
		client := resolveClientWithToken()
		ctx, cancel := newContext()
		defer cancel()

		err := client.DeleteHook(ctx, hookID)
		if err != nil {
			return fmt.Errorf("failed to delete hook: %w", err)
		}

		if err := outputJSON(map[string]string{"id": hookID, "status": "deleted"}); err != nil {
			return err
		}
		printSuccess("✓ Deleted hook %s", hookID)
		return nil
	},
}

func init() {
	hookListCmd.Flags().StringVar(&hookListPackage, "package", "", "Filter by package name")
	hookListCmd.Flags().IntVar(&hookListPage, "page", 0, "Page number")
	hookListCmd.Flags().IntVar(&hookListPerPage, "per-page", 20, "Results per page")

	hookCreateCmd.Flags().StringVar(&hookCreateName, "name", "", "Hook name (required)")
	hookCreateCmd.Flags().StringVar(&hookCreateEndpoint, "endpoint", "", "Webhook endpoint URL (required)")
	hookCreateCmd.Flags().StringVar(&hookCreateSecret, "secret", "", "Webhook secret for signature verification")
	hookCreateCmd.Flags().StringVar(&hookCreatePackage, "package", "", "Package to monitor (empty = all packages)")
	hookCreateCmd.Flags().StringArrayVar(&hookCreateEvents, "events", nil, "Event types (default: all)")

	hookUpdateCmd.Flags().StringVar(&hookUpdateEndpoint, "endpoint", "", "New endpoint URL")
	hookUpdateCmd.Flags().StringVar(&hookUpdateSecret, "secret", "", "New secret")
	hookUpdateCmd.Flags().StringArrayVar(&hookUpdateEvents, "events", nil, "New event types")
	hookUpdateCmd.Flags().BoolVar(&hookUpdateActive, "active", false, "Set hook active/inactive")
	hookUpdateCmd.Flags().BoolVar(&hookUpdateActiveSet, "set-active", false, "Explicitly set the active flag")

	hookCmd.AddCommand(hookListCmd)
	hookCmd.AddCommand(hookGetCmd)
	hookCmd.AddCommand(hookCreateCmd)
	hookCmd.AddCommand(hookUpdateCmd)
	hookCmd.AddCommand(hookDeleteCmd)
	rootCmd.AddCommand(hookCmd)
}
