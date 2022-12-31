/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/veracode/veracode-cli/cmd/scan"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Perform a scan of a container image and local build directory",
	Long: `Perform a scan of a container image and local build directory
Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		err := initialize(cmd)
		if err != nil {
			return err
		}
		err = initializeDocker(cmd.Context())
		return err
	},

	Run: func(cmd *cobra.Command, args []string) {
		scan.Run(cmd.Context(), args, cmd.Flags())
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)

	scanCmd.PersistentFlags().String("format", "json", "SBOM format to write out to")
	scanCmd.PersistentFlags().String("out", "-", "Output file")

	scanCmd.Flags().BoolP("verbose", "v", false, "Verbose output")
	scanCmd.Flags().BoolP("prettyprint", "", false, "Pretty print JSON")

	scanCmd.Flags().BoolP("kitchen-sink", "", false, "Runs everything we have available against a target and create an uber report at the end.")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
