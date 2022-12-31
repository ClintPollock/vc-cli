/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/veracode/veracode-cli/cmd/inspect"
)

// inspectCmd represents the inspect command
var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Inspect a target to scan or generate an SBOM against.",
	Long: `Inspect a target to scan or generate an SBOM against.
.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		err := initialize(cmd)
		if err != nil {
			return err
		}
		err = initializeDocker(cmd.Context())
		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		inspect.Run(cmd.Context(), args, cmd.Flags())
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)

	inspectCmd.PersistentFlags().String("out", "o", "Output file")

	inspectCmd.Flags().BoolP("verbose", "v", false, "Verbose output")
	inspectCmd.Flags().BoolP("prettyprint", "t", false, "Pretty print JSON")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// inspectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// inspectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
