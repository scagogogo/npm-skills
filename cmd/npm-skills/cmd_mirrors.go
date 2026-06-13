package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/scagogogo/npm-skills/pkg/registry"
	"github.com/spf13/cobra"
)

var mirrorsCmd = &cobra.Command{
	Use:   "mirrors",
	Short: "List available NPM mirror sources",
	Long: color.New(color.FgCyan).Sprintf("List available NPM mirror sources") + "\n\n" +
		"Shows all supported mirror sources with their URLs and descriptions.\n" +
		"Use the mirror name with --mirror flag in other commands.",
	Aliases: []string{"mirror", "ms"},
	Example: `  npm-skills mirrors`,
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		mirrors := registry.ListMirrors()

		// Print a nice table to stderr
		printHeader("Available NPM Mirror Sources")
		fmt.Fprintln(os.Stderr)
		for i, m := range mirrors {
			nameColor := color.New(color.FgGreen, color.Bold)
			urlColor := color.HiBlackString
			regionColor := color.New(color.FgYellow)
			descColor := color.HiWhiteString

			regionIcon := "🌍"
			if m.Region == "China" {
				regionIcon = "🇨🇳"
			}

			fmt.Fprintf(os.Stderr, "  %s %-12s  %s  %s %s  %s\n",
				color.HiBlackString("%d.", i+1),
				nameColor.Sprint(m.Name),
				urlColor(m.URL),
				regionIcon,
				regionColor.Sprint(m.Region),
				descColor(m.Description),
			)
		}
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, color.HiBlackString("  Usage: npm-skills <command> -m <name>"))
		fmt.Fprintln(os.Stderr, color.HiBlackString("  Or:    NPM_MIRROR=<name> npm-skills <command>"))
		fmt.Fprintln(os.Stderr)

		if err := outputJSON(mirrors); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(mirrorsCmd)
}
