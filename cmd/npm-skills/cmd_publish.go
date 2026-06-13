package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/scagogogo/npm-skills/pkg/models"
	"github.com/spf13/cobra"
)

var publishName string
var publishVersion string
var publishDescription string

var publishCmd = &cobra.Command{
	Use:   "publish <tarball-path>",
	Short: "Publish an NPM package to the registry",
	Long: color.New(color.FgCyan).Sprintf("Publish an NPM package to the registry") + "\n\n" +
		"Publishes a package tarball (.tgz) to the configured NPM registry.\n" +
		"Requires authentication token.\n\n" +
		color.HiYellowString("Note: ") + "This is a destructive operation. Make sure you have permission to publish.",
	Aliases: []string{"pub"},
	Example: `  npm-skills publish ./my-package-1.0.0.tgz --token npm_xxxxx
  npm-skills publish ./my-package-1.0.0.tgz -t npm_xxxxx -m npm-mirror
  npm-skills publish ./my-package-1.0.0.tgz --registry https://npm.mycompany.com -t npm_xxxxx`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(); err != nil {
			return err
		}

		tarballPath := args[0]
		printInfo("Publishing %s to %s...", tarballPath, currentMirrorLabel())

		// 读取 tarball 文件
		tarballBytes, err := os.ReadFile(tarballPath)
		if err != nil {
			return fmt.Errorf("failed to read tarball: %w", err)
		}

		if publishName == "" || publishVersion == "" {
			return fmt.Errorf("--name and --version are required")
		}

		client := resolveClientWithToken()
		ctx, cancel := newContext()
		defer cancel()

		metadata := &models.PublishMetadata{
			Name:        publishName,
			Version:     publishVersion,
			Description: publishDescription,
		}

		err = client.PublishPackageFromTarball(ctx, publishName, publishVersion, tarballBytes, metadata)
		if err != nil {
			return fmt.Errorf("failed to publish: %w", err)
		}

		result := map[string]string{
			"package": publishName,
			"version": publishVersion,
			"status":  "published",
		}
		if err := outputJSON(result); err != nil {
			return err
		}
		printSuccess("✓ Published %s@%s", color.New(color.FgWhite, color.Bold).Sprint(publishName), publishVersion)
		return nil
	},
}

func init() {
	publishCmd.Flags().StringVar(&publishName, "name", "", "Package name (required)")
	publishCmd.Flags().StringVar(&publishVersion, "version", "", "Package version (required)")
	publishCmd.Flags().StringVar(&publishDescription, "description", "", "Package description")
	rootCmd.AddCommand(publishCmd)
}