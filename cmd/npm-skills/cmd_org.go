package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var orgCmd = &cobra.Command{
	Use:   "org",
	Short: "Organization and team management",
	Long: color.New(color.FgCyan).Sprintf("Organization and team management") + "\n\n" +
		"Subcommands: get, create, delete, members, packages,\n" +
		"team-list, team-create, team-delete, team-members, team-packages",
	Example: `  npm-skills org get my-org -t npm_xxxxx
  npm-skills org members my-org -t npm_xxxxx`,
}

var orgGetCmd = &cobra.Command{
	Use:   "get <org>",
	Short: "Get organization details",
	Long: color.New(color.FgCyan).Sprintf("Get organization details") + "\n\n" +
		"Requires authentication token.",
	Example: `  npm-skills org get myorg -t npm_xxxxx`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(); err != nil {
			return err
		}

		orgName := args[0]
		printInfo("Getting org %s from %s...", orgName, currentMirrorLabel())
		client := resolveClientWithToken()
		ctx, cancel := newContext()
		defer cancel()

		org, err := client.GetOrg(ctx, orgName)
		if err != nil {
			return fmt.Errorf("failed to get org: %w", err)
		}

		if err := outputJSON(org); err != nil {
			return err
		}
		printSuccess("✓ Org: %s (scope: %s)", org.Name, org.Scope)
		return nil
	},
}

var orgCreateCmd = &cobra.Command{
	Use:   "create <org>",
	Short: "Create a new organization",
	Long: color.New(color.FgCyan).Sprintf("Create a new organization") + "\n\n" +
		"Requires authentication token.",
	Example: `  npm-skills org create myorg -t npm_xxxxx`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(); err != nil {
			return err
		}

		orgName := args[0]
		printInfo("Creating org %s on %s...", orgName, currentMirrorLabel())
		client := resolveClientWithToken()
		ctx, cancel := newContext()
		defer cancel()

		org, err := client.CreateOrg(ctx, orgName)
		if err != nil {
			return fmt.Errorf("failed to create org: %w", err)
		}

		if err := outputJSON(org); err != nil {
			return err
		}
		printSuccess("✓ Created org %s", color.New(color.FgWhite, color.Bold).Sprint(orgName))
		return nil
	},
}

var orgDeleteCmd = &cobra.Command{
	Use:   "delete <org>",
	Short: "Delete an organization",
	Long: color.New(color.FgCyan).Sprintf("Delete an organization") + "\n\n" +
		"Requires authentication token.\n\n" +
		color.HiRedString("WARNING: ") + "This is an irreversible operation.",
	Aliases: []string{"rm"},
	Example: `  npm-skills org delete myorg -t npm_xxxxx`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(); err != nil {
			return err
		}

		orgName := args[0]
		printInfo("Deleting org %s on %s...", orgName, currentMirrorLabel())
		client := resolveClientWithToken()
		ctx, cancel := newContext()
		defer cancel()

		err := client.DeleteOrg(ctx, orgName)
		if err != nil {
			return fmt.Errorf("failed to delete org: %w", err)
		}

		if err := outputJSON(map[string]string{"org": orgName, "status": "deleted"}); err != nil {
			return err
		}
		printSuccess("✓ Deleted org %s", color.New(color.FgWhite, color.Bold).Sprint(orgName))
		return nil
	},
}

var orgMembersCmd = &cobra.Command{
	Use:   "members <org>",
	Short: "List organization members",
	Long: color.New(color.FgCyan).Sprintf("List organization members") + "\n\n" +
		"Requires authentication token.",
	Aliases: []string{"member"},
	Example: `  npm-skills org members myorg -t npm_xxxxx`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(); err != nil {
			return err
		}

		orgName := args[0]
		printInfo("Listing members of org %s on %s...", orgName, currentMirrorLabel())
		client := resolveClientWithToken()
		ctx, cancel := newContext()
		defer cancel()

		members, err := client.ListOrgMembers(ctx, orgName)
		if err != nil {
			return fmt.Errorf("failed to list org members: %w", err)
		}

		if err := outputJSON(members); err != nil {
			return err
		}
		printSuccess("✓ %d members in org %s", len(members), orgName)
		return nil
	},
}

var orgMemberAddCmd = &cobra.Command{
	Use:   "member-add <org> <username>",
	Short: "Add a member to an organization",
	Long: color.New(color.FgCyan).Sprintf("Add a member to an organization") + "\n\n" +
		"Requires authentication token.",
	Example: `  npm-skills org member-add myorg username -t npm_xxxxx`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(); err != nil {
			return err
		}

		orgName := args[0]
		username := args[1]
		printInfo("Adding %s to org %s on %s...", username, orgName, currentMirrorLabel())
		client := resolveClientWithToken()
		ctx, cancel := newContext()
		defer cancel()

		err := client.AddOrgMember(ctx, orgName, username)
		if err != nil {
			return fmt.Errorf("failed to add org member: %w", err)
		}

		if err := outputJSON(map[string]string{"org": orgName, "user": username, "status": "added"}); err != nil {
			return err
		}
		printSuccess("✓ Added %s to org %s", username, orgName)
		return nil
	},
}

var orgMemberRemoveCmd = &cobra.Command{
	Use:   "member-remove <org> <username>",
	Short: "Remove a member from an organization",
	Long: color.New(color.FgCyan).Sprintf("Remove a member from an organization") + "\n\n" +
		"Requires authentication token.",
	Aliases: []string{"member-rm"},
	Example: `  npm-skills org member-remove myorg username -t npm_xxxxx`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(); err != nil {
			return err
		}

		orgName := args[0]
		username := args[1]
		printInfo("Removing %s from org %s on %s...", username, orgName, currentMirrorLabel())
		client := resolveClientWithToken()
		ctx, cancel := newContext()
		defer cancel()

		err := client.RemoveOrgMember(ctx, orgName, username)
		if err != nil {
			return fmt.Errorf("failed to remove org member: %w", err)
		}

		if err := outputJSON(map[string]string{"org": orgName, "user": username, "status": "removed"}); err != nil {
			return err
		}
		printSuccess("✓ Removed %s from org %s", username, orgName)
		return nil
	},
}

var orgPackagesCmd = &cobra.Command{
	Use:   "packages <org>",
	Short: "List packages owned by an organization",
	Long: color.New(color.FgCyan).Sprintf("List packages owned by an organization") + "\n\n" +
		"Requires authentication token.",
	Aliases: []string{"pkgs"},
	Example: `  npm-skills org packages myorg -t npm_xxxxx`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(); err != nil {
			return err
		}

		orgName := args[0]
		printInfo("Listing packages of org %s on %s...", orgName, currentMirrorLabel())
		client := resolveClientWithToken()
		ctx, cancel := newContext()
		defer cancel()

		packages, err := client.ListOrgPackages(ctx, orgName)
		if err != nil {
			return fmt.Errorf("failed to list org packages: %w", err)
		}

		if err := outputJSON(packages); err != nil {
			return err
		}
		printSuccess("✓ %d packages in org %s", len(packages), orgName)
		return nil
	},
}

// ===== 团队子命令 =====

var teamListCmd = &cobra.Command{
	Use:   "team-list <org>",
	Short: "List teams in an organization",
	Long: color.New(color.FgCyan).Sprintf("List teams in an organization") + "\n\n" +
		"Requires authentication token.",
	Aliases: []string{"teams"},
	Example: `  npm-skills org team-list myorg -t npm_xxxxx`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(); err != nil {
			return err
		}

		orgName := args[0]
		printInfo("Listing teams in org %s on %s...", orgName, currentMirrorLabel())
		client := resolveClientWithToken()
		ctx, cancel := newContext()
		defer cancel()

		teams, err := client.ListTeams(ctx, orgName)
		if err != nil {
			return fmt.Errorf("failed to list teams: %w", err)
		}

		if err := outputJSON(teams); err != nil {
			return err
		}
		printSuccess("✓ %d teams in org %s", len(teams), orgName)
		return nil
	},
}

var teamCreateCmd = &cobra.Command{
	Use:   "team-create <org> <team>",
	Short: "Create a team in an organization",
	Long: color.New(color.FgCyan).Sprintf("Create a team in an organization") + "\n\n" +
		"Requires authentication token.",
	Example: `  npm-skills org team-create myorg devteam -t npm_xxxxx`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(); err != nil {
			return err
		}

		orgName := args[0]
		teamName := args[1]
		printInfo("Creating team %s in org %s on %s...", teamName, orgName, currentMirrorLabel())
		client := resolveClientWithToken()
		ctx, cancel := newContext()
		defer cancel()

		team, err := client.CreateTeam(ctx, orgName, teamName)
		if err != nil {
			return fmt.Errorf("failed to create team: %w", err)
		}

		if err := outputJSON(team); err != nil {
			return err
		}
		printSuccess("✓ Created team %s in org %s", teamName, orgName)
		return nil
	},
}

var teamDeleteCmd = &cobra.Command{
	Use:   "team-delete <org> <team>",
	Short: "Delete a team from an organization",
	Long: color.New(color.FgCyan).Sprintf("Delete a team from an organization") + "\n\n" +
		"Requires authentication token.\n\n" +
		color.HiRedString("WARNING: ") + "This is an irreversible operation.",
	Aliases: []string{"team-rm"},
	Example: `  npm-skills org team-delete myorg devteam -t npm_xxxxx`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(); err != nil {
			return err
		}

		orgName := args[0]
		teamName := args[1]
		printInfo("Deleting team %s from org %s on %s...", teamName, orgName, currentMirrorLabel())
		client := resolveClientWithToken()
		ctx, cancel := newContext()
		defer cancel()

		err := client.DeleteTeam(ctx, orgName, teamName)
		if err != nil {
			return fmt.Errorf("failed to delete team: %w", err)
		}

		if err := outputJSON(map[string]string{"org": orgName, "team": teamName, "status": "deleted"}); err != nil {
			return err
		}
		printSuccess("✓ Deleted team %s from org %s", teamName, orgName)
		return nil
	},
}

var teamMembersCmd = &cobra.Command{
	Use:   "team-members <org> <team>",
	Short: "List members of a team",
	Long: color.New(color.FgCyan).Sprintf("List members of a team") + "\n\n" +
		"Requires authentication token.",
	Example: `  npm-skills org team-members myorg devteam -t npm_xxxxx`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(); err != nil {
			return err
		}

		orgName := args[0]
		teamName := args[1]
		printInfo("Listing members of team %s/%s on %s...", orgName, teamName, currentMirrorLabel())
		client := resolveClientWithToken()
		ctx, cancel := newContext()
		defer cancel()

		members, err := client.ListTeamMembers(ctx, orgName, teamName)
		if err != nil {
			return fmt.Errorf("failed to list team members: %w", err)
		}

		if err := outputJSON(members); err != nil {
			return err
		}
		printSuccess("✓ %d members in team %s/%s", len(members), orgName, teamName)
		return nil
	},
}

var teamMemberAddCmd = &cobra.Command{
	Use:   "team-member-add <org> <team> <username>",
	Short: "Add a member to a team",
	Long: color.New(color.FgCyan).Sprintf("Add a member to a team") + "\n\n" +
		"Requires authentication token. The user must already be a member of the organization.",
	Example: `  npm-skills org team-member-add myorg devteam username -t npm_xxxxx`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(); err != nil {
			return err
		}

		orgName := args[0]
		teamName := args[1]
		username := args[2]
		printInfo("Adding %s to team %s/%s on %s...", username, orgName, teamName, currentMirrorLabel())
		client := resolveClientWithToken()
		ctx, cancel := newContext()
		defer cancel()

		err := client.AddTeamMember(ctx, orgName, teamName, username)
		if err != nil {
			return fmt.Errorf("failed to add team member: %w", err)
		}

		if err := outputJSON(map[string]string{"org": orgName, "team": teamName, "user": username, "status": "added"}); err != nil {
			return err
		}
		printSuccess("✓ Added %s to team %s/%s", username, orgName, teamName)
		return nil
	},
}

var teamMemberRemoveCmd = &cobra.Command{
	Use:   "team-member-remove <org> <team> <username>",
	Short: "Remove a member from a team",
	Long: color.New(color.FgCyan).Sprintf("Remove a member from a team") + "\n\n" +
		"Requires authentication token.",
	Aliases: []string{"team-member-rm"},
	Example: `  npm-skills org team-member-remove myorg devteam username -t npm_xxxxx`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(); err != nil {
			return err
		}

		orgName := args[0]
		teamName := args[1]
		username := args[2]
		printInfo("Removing %s from team %s/%s on %s...", username, orgName, teamName, currentMirrorLabel())
		client := resolveClientWithToken()
		ctx, cancel := newContext()
		defer cancel()

		err := client.RemoveTeamMember(ctx, orgName, teamName, username)
		if err != nil {
			return fmt.Errorf("failed to remove team member: %w", err)
		}

		if err := outputJSON(map[string]string{"org": orgName, "team": teamName, "user": username, "status": "removed"}); err != nil {
			return err
		}
		printSuccess("✓ Removed %s from team %s/%s", username, orgName, teamName)
		return nil
	},
}

var teamPackagesCmd = &cobra.Command{
	Use:   "team-packages <org> <team>",
	Short: "List packages a team has access to",
	Long: color.New(color.FgCyan).Sprintf("List packages a team has access to") + "\n\n" +
		"Requires authentication token.",
	Example: `  npm-skills org team-packages myorg devteam -t npm_xxxxx`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(); err != nil {
			return err
		}

		orgName := args[0]
		teamName := args[1]
		printInfo("Listing packages for team %s/%s on %s...", orgName, teamName, currentMirrorLabel())
		client := resolveClientWithToken()
		ctx, cancel := newContext()
		defer cancel()

		packages, err := client.ListTeamPackages(ctx, orgName, teamName)
		if err != nil {
			return fmt.Errorf("failed to list team packages: %w", err)
		}

		if err := outputJSON(packages); err != nil {
			return err
		}
		printSuccess("✓ %d packages for team %s/%s", len(packages), orgName, teamName)
		return nil
	},
}

func init() {
	orgCmd.AddCommand(orgGetCmd)
	orgCmd.AddCommand(orgCreateCmd)
	orgCmd.AddCommand(orgDeleteCmd)
	orgCmd.AddCommand(orgMembersCmd)
	orgCmd.AddCommand(orgMemberAddCmd)
	orgCmd.AddCommand(orgMemberRemoveCmd)
	orgCmd.AddCommand(orgPackagesCmd)
	orgCmd.AddCommand(teamListCmd)
	orgCmd.AddCommand(teamCreateCmd)
	orgCmd.AddCommand(teamDeleteCmd)
	orgCmd.AddCommand(teamMembersCmd)
	orgCmd.AddCommand(teamMemberAddCmd)
	orgCmd.AddCommand(teamMemberRemoveCmd)
	orgCmd.AddCommand(teamPackagesCmd)
	rootCmd.AddCommand(orgCmd)
}