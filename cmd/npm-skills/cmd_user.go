package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/scagogogo/npm-skills/pkg/models"
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:     "user",
	Short:   "User operations (login, signup, get)",
	Long:    color.New(color.FgCyan).Sprintf("User operations") + "\n\n" + "Subcommands: login, signup, get",
	Aliases: []string{"u"},
	Example: `  npm-skills user login -u myuser -p mypass
  npm-skills user get myuser`,
}

var loginUsername string
var loginPassword string
var signupUsername string
var signupPassword string
var signupEmail string

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to NPM registry",
	Long: color.New(color.FgCyan).Sprintf("Login to NPM registry") + "\n\n" +
		"Authenticates with the NPM registry and returns an authentication token.",
	Example: `  npm-skills user login --username myuser --password mypass
  npm-skills user login -u myuser -p mypass -m npm-mirror`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if loginUsername == "" || loginPassword == "" {
			return fmt.Errorf("--username and --password are required")
		}

		printInfo("Logging in to %s...", currentMirrorLabel())
		client := resolveClient()
		ctx, cancel := newContext()
		defer cancel()

		result, err := client.Login(ctx, loginUsername, loginPassword)
		if err != nil {
			return fmt.Errorf("login failed: %w", err)
		}

		if err := outputJSON(result); err != nil {
			return err
		}
		printSuccess("✓ Logged in as %s", color.New(color.FgWhite, color.Bold).Sprint(loginUsername))
		printInfo("Token: %s", result.Token)
		return nil
	},
}

var signupCmd = &cobra.Command{
	Use:     "signup",
	Short:   "Create a new NPM user account",
	Long:    color.New(color.FgCyan).Sprintf("Create a new NPM user account"),
	Example: `  npm-skills user signup --username myuser --password mypass --email me@example.com`,
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if signupUsername == "" || signupPassword == "" || signupEmail == "" {
			return fmt.Errorf("--username, --password, and --email are required")
		}

		printInfo("Creating user %s on %s...", signupUsername, currentMirrorLabel())
		client := resolveClient()
		ctx, cancel := newContext()
		defer cancel()

		result, err := client.CreateUser(ctx, &models.UserCreation{Name: signupUsername, Password: signupPassword, Email: signupEmail})
		if err != nil {
			return fmt.Errorf("signup failed: %w", err)
		}

		if err := outputJSON(result); err != nil {
			return err
		}
		printSuccess("✓ User %s created", color.New(color.FgWhite, color.Bold).Sprint(signupUsername))
		return nil
	},
}

var userGetCmd = &cobra.Command{
	Use:     "get <username>",
	Short:   "Get user profile information",
	Long:    color.New(color.FgCyan).Sprintf("Get user profile information"),
	Aliases: []string{"info"},
	Example: `  npm-skills user get myuser`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(); err != nil {
			return err
		}

		username := args[0]
		printInfo("Getting user %s from %s...", username, currentMirrorLabel())
		client := resolveClientWithToken()
		ctx, cancel := newContext()
		defer cancel()

		profile, err := client.GetUser(ctx, username)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}

		if err := outputJSON(profile); err != nil {
			return err
		}
		printSuccess("✓ User: %s (%s)", profile.Name, profile.Email)
		return nil
	},
}

func init() {
	loginCmd.Flags().StringVarP(&loginUsername, "username", "u", "", "Username")
	loginCmd.Flags().StringVarP(&loginPassword, "password", "p", "", "Password")
	signupCmd.Flags().StringVarP(&signupUsername, "username", "u", "", "Username")
	signupCmd.Flags().StringVarP(&signupPassword, "password", "p", "", "Password")
	signupCmd.Flags().StringVar(&signupEmail, "email", "", "Email address")

	userCmd.AddCommand(loginCmd)
	userCmd.AddCommand(signupCmd)
	userCmd.AddCommand(userGetCmd)
	rootCmd.AddCommand(userCmd)
}
