package main

import (
	"fmt"
	"strconv"

	"github.com/fatih/color"
	"github.com/scagogogo/npm-skills/pkg/models"
	"github.com/spf13/cobra"
)

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Security audit for NPM packages",
	Long: color.New(color.FgCyan).Sprintf("Security audit for NPM packages") + "\n\n" +
		"Subcommands: quick, bulk, advisory, advisories\n\n" +
		"Check dependencies for known vulnerabilities using the NPM registry's\n" +
		"security advisory database.",
	Example: `  npm-skills audit quick --deps "lodash=4.17.11"
  npm-skills audit advisory 123`,
}

var auditQuickDeps string
var auditBulkAdvisories string

var auditQuickCmd = &cobra.Command{
	Use:   "quick",
	Short: "Quick security audit of dependencies",
	Long: color.New(color.FgCyan).Sprintf("Quick security audit of dependencies") + "\n\n" +
		"Submits a dependency list for a fast vulnerability check.\n" +
		"Returns vulnerability counts by severity level.",
	Example: `  npm-skills audit quick --deps "lodash=4.17.11,express=4.17.1"`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if auditQuickDeps == "" {
			return fmt.Errorf("--deps is required (format: name=version,name2=version2)")
		}

		// 解析依赖字符串
		deps := parseDepsString(auditQuickDeps)

		printInfo("Running quick audit on %s...", currentMirrorLabel())
		client := resolveClient()
		ctx, cancel := newContext()
		defer cancel()

		result, err := client.QuickAudit(ctx, &models.QuickAuditRequest{Dependencies: deps})
		if err != nil {
			return fmt.Errorf("quick audit failed: %w", err)
		}

		if err := outputJSON(result); err != nil {
			return err
		}
		v := result.Metadata.Vulnerabilities
		total := v.Low + v.Moderate + v.High + v.Critical
		if total == 0 {
			printSuccess("✓ No vulnerabilities found")
		} else {
			printInfo("Found %d vulnerabilities: %d low, %d moderate, %d high, %d critical",
				total, v.Low, v.Moderate, v.High, v.Critical)
		}
		return nil
	},
}

var auditBulkCmd = &cobra.Command{
	Use:   "bulk",
	Short: "Bulk security audit with detailed advisories",
	Long: color.New(color.FgCyan).Sprintf("Bulk security audit with detailed advisories") + "\n\n" +
		"Submits package names and version ranges, returns matching security advisories\n" +
		"with full details including CVE, severity, and patch information.",
	Example: `  npm-skills audit bulk --advisories "lodash=<4.17.12,express=<4.17.3"`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if auditBulkAdvisories == "" {
			return fmt.Errorf("--advisories is required (format: name=<version,name2=<version2)")
		}

		// 解析公告字符串为 map[string][]string
		advisories := parseAdvisoriesString(auditBulkAdvisories)

		printInfo("Running bulk audit on %s...", currentMirrorLabel())
		client := resolveClient()
		ctx, cancel := newContext()
		defer cancel()

		result, err := client.BulkAudit(ctx, advisories)
		if err != nil {
			return fmt.Errorf("bulk audit failed: %w", err)
		}

		if err := outputJSON(result); err != nil {
			return err
		}
		totalAdvisories := 0
		for _, a := range result {
			totalAdvisories += len(a)
		}
		printSuccess("✓ Found %d advisories across %d packages", totalAdvisories, len(result))
		return nil
	},
}

var auditAdvisoryCmd = &cobra.Command{
	Use:   "advisory <id>",
	Short: "Get details of a specific security advisory",
	Long: color.New(color.FgCyan).Sprintf("Get details of a specific security advisory") + "\n\n" +
		"Returns full details of a security advisory including overview,\n" +
		"recommendation, and affected versions.",
	Example: `  npm-skills audit advisory 1234`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		advisoryID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("advisory ID must be a number")
		}

		printInfo("Getting advisory %d from %s...", advisoryID, currentMirrorLabel())
		client := resolveClient()
		ctx, cancel := newContext()
		defer cancel()

		advisory, err := client.GetAdvisory(ctx, advisoryID)
		if err != nil {
			return fmt.Errorf("failed to get advisory: %w", err)
		}

		if err := outputJSON(advisory); err != nil {
			return err
		}
		printSuccess("✓ Advisory %d: %s [%s]", advisory.ID, advisory.Title, advisory.Severity)
		return nil
	},
}

var auditListPage int
var auditListPerPage int
var auditListPackage string

var auditListCmd = &cobra.Command{
	Use:   "advisories",
	Short: "List security advisories",
	Long: color.New(color.FgCyan).Sprintf("List security advisories") + "\n\n" +
		"Returns a paginated list of security advisories, optionally filtered by package name.",
	Example: `  npm-skills audit advisories
  npm-skills audit advisories --package lodash --per-page 10`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		printInfo("Listing advisories on %s...", currentMirrorLabel())
		client := resolveClient()
		ctx, cancel := newContext()
		defer cancel()

		opts := models.AdvisoryListOptions{
			Page:            auditListPage,
			PerPage:         auditListPerPage,
			AffectedPackage: auditListPackage,
		}

		advisories, err := client.ListAdvisories(ctx, opts)
		if err != nil {
			return fmt.Errorf("failed to list advisories: %w", err)
		}

		if err := outputJSON(advisories); err != nil {
			return err
		}
		printSuccess("✓ Found %d advisories", len(advisories))
		return nil
	},
}

func init() {
	auditQuickCmd.Flags().StringVar(&auditQuickDeps, "deps", "", "Dependencies to audit (format: name=version,name2=version2)")
	auditBulkCmd.Flags().StringVar(&auditBulkAdvisories, "advisories", "", "Advisories to check (format: name=<version,name2=<version2)")
	auditListCmd.Flags().IntVar(&auditListPage, "page", 0, "Page number")
	auditListCmd.Flags().IntVar(&auditListPerPage, "per-page", 20, "Results per page")
	auditListCmd.Flags().StringVar(&auditListPackage, "package", "", "Filter by affected package name")

	auditCmd.AddCommand(auditQuickCmd)
	auditCmd.AddCommand(auditBulkCmd)
	auditCmd.AddCommand(auditAdvisoryCmd)
	auditCmd.AddCommand(auditListCmd)
	rootCmd.AddCommand(auditCmd)
}