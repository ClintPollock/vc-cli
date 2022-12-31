/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/veracode/veracode-cli/cmd/research"
)

// researchCmd represents the research command
var researchCmd = &cobra.Command{
	Use:   "research",
	Short: "Runs a set of comparisons of tools against the target.",
	Long: `Research Mode.

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
		research.Run(cmd.Context(), args, cmd.Flags())
	},
}

func init() {
	rootCmd.AddCommand(researchCmd)

	researchCmd.PersistentFlags().String("format", "json", "SBOM format to write out to")
	researchCmd.PersistentFlags().String("out", "-", "Output file")

	researchCmd.Flags().BoolP("verbose", "v", false, "Verbose output")
	researchCmd.Flags().BoolP("prettyprint", "", false, "Pretty print JSON")

	researchCmd.Flags().BoolP("kitchen-sink", "", false, "Runs everything we have available against a target and create an uber report at the end.")

}
