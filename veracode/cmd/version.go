package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/veracode/veracode-cli/cmd/version"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Return the Veracode CLI version.",
	Long:  `Return the Veracode CLI version.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Veracode CLI v%s -- %s\n", version.Version, version.GitHash)
	},
}
