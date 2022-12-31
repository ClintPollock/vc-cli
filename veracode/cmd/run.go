/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/veracode/veracode-cli/cmd/run"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Directly call a tool within the package docker container. For advanced usage only.",
	Long: `Directly call a tool within the package docker container. For advanced usage only.

Ignores all options past the "run" command and passes to the docker container's entrypoint, and
returns output from stdout and stderr. Requires knowledge of the composition of the container
for effeective usage.`,

	Args: cobra.MinimumNArgs(1),

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		err := initialize(cmd)
		if err != nil {
			return err
		}
		err = initializeDocker(cmd.Context())
		return err
	},

	Run: func(cmd *cobra.Command, args []string) {
		argsAfterRunCommand := os.Args[2:]
		run.Run(cmd.Context(), argsAfterRunCommand)
	},
	FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
