/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/veracode/veracode-cli/cmd/sbom"

	"github.com/spf13/cobra"
)

// sbomCmd represents the sbom command
var sbomCmd = &cobra.Command{
	Use:   "sbom",
	Short: "Generate an Software Bill Of Materials (SBOM) of an image, archive, repo or directory",
	Long: `Generate an Software Bill Of Materials (SBOM) of an image, archive, repo or directory.

`,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		err := initialize(cmd)
		if err != nil {
			return err
		}
		err = initializeDocker(cmd.Context())
		return err
	},

	Run: func(cmd *cobra.Command, args []string) {
		sbom.Run(cmd.Context(), args, cmd.Flags())
	},
}

func init() {

	rootCmd.AddCommand(sbomCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	sbomCmd.PersistentFlags().String("format", "json", "SBOM format to write out to (json, spdx, cyclonedx, github, directory, table)")
	sbomCmd.PersistentFlags().String("out", "-", "Output file")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	sbomCmd.Flags().BoolP("verbose", "v", false, "Verbose output")
	sbomCmd.Flags().BoolP("prettyprint", "t", false, "Pretty print JSON")

}
