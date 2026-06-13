package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download <name> <version> <dest>",
	Short: "Download an NPM package tarball (.tgz)",
	Long: color.New(color.FgCyan).Sprintf("Download an NPM package tarball (.tgz)") + "\n\n" +
		"Downloads the specified package version as a .tgz file to the given local path.\n" +
		"Use \"latest\" as the version string to download the latest version.\n\n" +
		color.HiBlackString("Mirror: %s", mirrorNames()) + " (via --mirror or --registry flag)",
	Aliases: []string{"dl"},
	Example: `  npm-skills download react 18.2.0 ./react-18.2.0.tgz
  npm-skills download lodash latest ./lodash.tgz -m npm-mirror
  npm-skills dl axios 1.0.0 /tmp/axios.tgz --proxy http://127.0.0.1:7890
  npm-skills download react 18.2.0 ./react.tgz --registry https://registry.npmmirror.com`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		packageName := args[0]
		ver := args[1]
		destPath := args[2]

		printInfo("Downloading %s to %s (source: %s)...",
			color.New(color.FgWhite, color.Bold).Sprintf("%s@%s", packageName, ver),
			color.New(color.FgWhite).Sprint(destPath),
			currentMirrorLabel())
		client := resolveClient()
		ctx, cancel := newContext()
		defer cancel()

		err := client.DownloadTarball(ctx, packageName, ver, destPath)
		if err != nil {
			return fmt.Errorf("failed to download '%s@%s': %w", packageName, ver, err)
		}

		result := map[string]interface{}{
			"package": packageName,
			"version": ver,
			"path":    destPath,
			"source":  currentMirrorLabel(),
			"status":  "downloaded",
		}
		if err := outputJSON(result); err != nil {
			return err
		}
		printSuccess("✓ %s downloaded successfully to %s",
			color.New(color.FgWhite, color.Bold).Sprintf("%s@%s", packageName, ver),
			color.New(color.FgWhite).Sprint(destPath))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}
