package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/scagogogo/npm-skills/pkg/models"
	"github.com/spf13/cobra"
)

var accessCmd = &cobra.Command{
	Use:   "access",
	Short: "Package access and collaborator management",
	Long: color.New(color.FgCyan).Sprintf("Package access and collaborator management") + "\n\n" +
		"Subcommands: get, set, collaborators, grant, revoke",
	Example: `  npm-skills access get my-package -t npm_xxxxx
  npm-skills access set my-package --visibility public -t npm_xxxxx`,
}

var accessSetVisibility string

var accessGetCmd = &cobra.Command{
	Use:   "get <package>",
	Short: "Get package access settings",
	Long: color.New(color.FgCyan).Sprintf("Get package access settings") + "\n\n" +
		"Requires authentication token.",
	Example: `  npm-skills access get my-package -t npm_xxxxx`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(); err != nil {
			return err
		}
		packageName := args[0]
		printInfo("Getting access for %s...", packageName)
		client := resolveClientWithToken()
		ctx, cancel := newContext()
		defer cancel()
		access, err := client.GetPackageAccess(ctx, packageName)
		if err != nil {
			return fmt.Errorf("failed to get access: %w", err)
		}
		if err := outputJSON(access); err != nil {
			return err
		}
		printSuccess("✓ Access info for %s", packageName)
		return nil
	},
}

var accessSetCmd = &cobra.Command{
	Use:   "set <package>",
	Short: "Set package access (public/restricted)",
	Long: color.New(color.FgCyan).Sprintf("Set package access (public/restricted)") + "\n\n" +
		"Requires authentication token.\n\n" +
		color.HiYellowString("Note: ") + "Changing from public to restricted will make the package inaccessible to unauthorized users.",
	Example: `  npm-skills access set my-package --visibility public -t npm_xxxxx
  npm-skills access set my-package --visibility restricted -t npm_xxxxx`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(); err != nil {
			return err
		}
		packageName := args[0]
		if accessSetVisibility == "" {
			return fmt.Errorf("--visibility is required (public or restricted)")
		}
		if accessSetVisibility != "public" && accessSetVisibility != "restricted" {
			return fmt.Errorf("--visibility must be 'public' or 'restricted'")
		}
		printInfo("Setting %s to %s...", packageName, accessSetVisibility)
		client := resolveClientWithToken()
		ctx, cancel := newContext()
		defer cancel()
		err := client.SetPackageAccess(ctx, packageName, &models.PackageAccessUpdate{Access: accessSetVisibility})
		if err != nil {
			return fmt.Errorf("failed to set access: %w", err)
		}
		if err := outputJSON(map[string]string{"package": packageName, "access": accessSetVisibility, "status": "updated"}); err != nil {
			return err
		}
		printSuccess("✓ Access for %s set to %s", packageName, accessSetVisibility)
		return nil
	},
}

var collaboratorsCmd = &cobra.Command{
	Use:   "collaborators <package>",
	Short: "List package collaborators",
	Long: color.New(color.FgCyan).Sprintf("List package collaborators") + "\n\n" +
		"Requires authentication token.",
	Aliases: []string{"collabs"},
	Example: `  npm-skills access collaborators my-package -t npm_xxxxx
  npm-skills access collabs my-package -t npm_xxxxx`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(); err != nil {
			return err
		}
		packageName := args[0]
		printInfo("Listing collaborators for %s...", packageName)
		client := resolveClientWithToken()
		ctx, cancel := newContext()
		defer cancel()
		collabs, err := client.ListCollaborators(ctx, packageName)
		if err != nil {
			return fmt.Errorf("failed to list collaborators: %w", err)
		}
		if err := outputJSON(collabs); err != nil {
			return err
		}
		printSuccess("✓ %d collaborators for %s", len(collabs), packageName)
		return nil
	},
}

var grantPermission string

var grantCmd = &cobra.Command{
	Use:   "grant <package> <user>",
	Short: "Grant user access to a package",
	Long: color.New(color.FgCyan).Sprintf("Grant user access to a package") + "\n\n" +
		"Requires authentication token. For teams, use the format <org>:<team>.",
	Example: `  npm-skills access grant my-package username --permission write -t npm_xxxxx
  npm-skills access grant my-package myorg:devteam --permission read -t npm_xxxxx`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(); err != nil {
			return err
		}
		packageName := args[0]
		username := args[1]
		perm := models.Permission(grantPermission)
		if perm != models.PermissionRead && perm != models.PermissionWrite {
			return fmt.Errorf("--permission must be 'read' or 'write'")
		}
		printInfo("Granting %s access to %s on %s...", perm, username, packageName)
		client := resolveClientWithToken()
		ctx, cancel := newContext()
		defer cancel()
		err := client.GrantAccess(ctx, packageName, username, perm)
		if err != nil {
			return fmt.Errorf("failed to grant access: %w", err)
		}
		if err := outputJSON(map[string]string{"package": packageName, "user": username, "permission": string(perm), "status": "granted"}); err != nil {
			return err
		}
		printSuccess("✓ Granted %s access to %s on %s", perm, username, packageName)
		return nil
	},
}

var revokeCmd = &cobra.Command{
	Use:   "revoke <package> <user>",
	Short: "Revoke user access from a package",
	Long: color.New(color.FgCyan).Sprintf("Revoke user access from a package") + "\n\n" +
		"Requires authentication token.",
	Example: `  npm-skills access revoke my-package username -t npm_xxxxx`,
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(); err != nil {
			return err
		}
		packageName := args[0]
		username := args[1]
		printInfo("Revoking access for %s on %s...", username, packageName)
		client := resolveClientWithToken()
		ctx, cancel := newContext()
		defer cancel()
		err := client.RevokeAccess(ctx, packageName, username)
		if err != nil {
			return fmt.Errorf("failed to revoke access: %w", err)
		}
		if err := outputJSON(map[string]string{"package": packageName, "user": username, "status": "revoked"}); err != nil {
			return err
		}
		printSuccess("✓ Revoked access for %s on %s", username, packageName)
		return nil
	},
}

func init() {
	accessSetCmd.Flags().StringVar(&accessSetVisibility, "visibility", "", "Access level: public or restricted")
	grantCmd.Flags().StringVar(&grantPermission, "permission", "read", "Permission level: read or write")

	accessCmd.AddCommand(accessGetCmd)
	accessCmd.AddCommand(accessSetCmd)
	accessCmd.AddCommand(collaboratorsCmd)
	accessCmd.AddCommand(grantCmd)
	accessCmd.AddCommand(revokeCmd)
	rootCmd.AddCommand(accessCmd)
}
