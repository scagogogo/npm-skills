package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var unpublishForce bool
var unpublishVersion string

var unpublishCmd = &cobra.Command{
	Use:   "unpublish <package>",
	Short: "Unpublish an NPM package or specific version",
	Long: color.New(color.FgCyan).Sprintf("Unpublish an NPM package or specific version") + "\n\n" +
		"Removes a package or a specific version from the registry.\n" +
		"Requires authentication token.\n\n" +
		color.HiRedString("WARNING: ") + "This is a destructive and potentially irreversible operation.\n" +
		"Most registries only allow unpublish within 72 hours of publish.",
	Aliases: []string{"unpub"},
	Example: `  npm-skills unpublish my-package --version 1.0.0 --token npm_xxxxx
  npm-skills unpublish my-package --force --token npm_xxxxx`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(); err != nil {
			return err
		}

		packageName := args[0]
		client := resolveClientWithToken()
		ctx, cancel := newContext()
		defer cancel()

		if unpublishVersion != "" {
			// 取消发布指定版本
			printInfo("Unpublishing %s@%s from %s...", packageName, unpublishVersion, currentMirrorLabel())
			err := client.UnpublishPackageVersion(ctx, packageName, unpublishVersion)
			if err != nil {
				return fmt.Errorf("failed to unpublish version: %w", err)
			}
			result := map[string]string{
				"package": packageName,
				"version": unpublishVersion,
				"status":  "unpublished",
			}
			if err := outputJSON(result); err != nil {
				return err
			}
			printSuccess("✓ Unpublished %s@%s", color.New(color.FgWhite, color.Bold).Sprint(packageName), unpublishVersion)
		} else if unpublishForce {
			// 强制取消发布整个包
			printInfo("Unpublishing entire package %s from %s...", packageName, currentMirrorLabel())
			err := client.UnpublishPackage(ctx, packageName)
			if err != nil {
				return fmt.Errorf("failed to unpublish package: %w", err)
			}
			result := map[string]string{
				"package": packageName,
				"status":  "unpublished",
			}
			if err := outputJSON(result); err != nil {
				return err
			}
			printSuccess("✓ Unpublished %s", color.New(color.FgWhite, color.Bold).Sprint(packageName))
		} else {
			return fmt.Errorf("specify --version <version> to unpublish a specific version, or --force to unpublish entire package")
		}
		return nil
	},
}

func init() {
	unpublishCmd.Flags().StringVar(&unpublishVersion, "version", "", "Version to unpublish")
	unpublishCmd.Flags().BoolVar(&unpublishForce, "force", false, "Unpublish entire package (dangerous)")
	rootCmd.AddCommand(unpublishCmd)
}