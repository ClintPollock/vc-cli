/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/veracode/veracode-cli/cmd/fix"
)

// Fix command runs the auto-remediation service
var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "Generate a fix for a given flaw.",
	Long: "Invoke the auto-remediation tool to generate a fix for the given flaw in the given source file.",
	Args: cobra.ExactArgs(1),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initialize(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		validateCredentials(cmd.Context())
		fix.Run(cmd.Context(), args, cmd.Flags())
	},
}

func init() {
	// -- Main flags
	fixCmd.Flags().String("apihost", "api.veracode.com", "Specify the host of the API endpoint")
	fixCmd.Flags().String("scanresults", "results.json", "results file from pipeline scan (JSON format)" )
	fixCmd.Flags().Int("issueid", -1, "ID of the issue to fix")
	fixCmd.Flags().BoolP("reuse", "r", false, "Reuse existing fix file instead of generating a new one")
	fixCmd.Flags().BoolP("apply", "a", false, "Automatically apply the top-ranked fix")

	// -- Debug flags
	fixCmd.Flags().BoolP("verbose", "v", false, "Verbose output")
	fixCmd.Flags().BoolP("prettyprint", "t", false, "Pretty print JSON")
	fixCmd.Flags().BoolP("debug", "d", false, "Be VERY noisy in running the command")

	rootCmd.AddCommand(fixCmd)
}
