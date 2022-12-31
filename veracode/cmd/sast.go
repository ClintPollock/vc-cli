/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/veracode/veracode-cli/cmd/sast"
)

// sastCmd represents the sbom command
var sastCmd = &cobra.Command{
	Use:   "sast",
	Short: "Scan an artifact using Veracode SAST.",
	Long: `Scan an artifact using Veracode SAST.

`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initialize(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		validateCredentials(cmd.Context())
		sast.Run(cmd.Context(), args, cmd.Flags())
	},
}

func init() {

	rootCmd.AddCommand(sastCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	sastCmd.PersistentFlags().String("out", "-", "Output file")

	sastCmd.PersistentFlags().String("scan-id", "", "ID of scan to reconnect to if connection to SAST API is broken")

	sastCmd.PersistentFlags().String("project-name", "MyProject", "Name of the project")
	sastCmd.PersistentFlags().String("project-url", "https://github.com/my/project", "URL to project source code repository")
	sastCmd.PersistentFlags().String("project-ref", "my-project", "String that prpvides a unique reference to the project")
	sastCmd.PersistentFlags().String("app-id", "NONE", "String that provides a unique reference to a application for linking")
	sastCmd.PersistentFlags().String("stage", "DEVELOPMENT", "String that refers to lifecycle stage for this code. May be DEVELOPMENT, TESTING, or RELEASE")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	//sastCmd.Flags().BoolP("verbose", "v", false, "Verbose output")
	sastCmd.Flags().BoolP("prettyprint", "p", false, "Pretty print JSON")
	sastCmd.Flags().BoolP("debug", "d", false, "Be VERY noisy in runnin the command")
	sastCmd.Flags().BoolP("show-details", "", false, "Show issue details in output")

	sastCmd.Flags().Float32("fail-on-severity", 7.0, "Maximum severity beyond which the scan will be marked as failed.")
	sastCmd.Flags().String("fail-on-cwe", "80", "CWEs that automatically trigger a failed policy.")
	sastCmd.Flags().BoolP("fail-fast", "", false, "Whether to exit immediately once a policy failure is found.")
	sastCmd.Flags().String("baseline-file", "", "List of known flaws to ignore in the findings of this scan")

	sastCmd.Flags().BoolP("emit-stack-dump", "", false, "Request the SAST engine generate stack dumps")
}
