/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/veracode/veracode-cli/cmd/configure"
)

// sbomCmd represents the sbom command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure configuration & credentials for the Veracode CLI",
	Long: `Configure configuration & credentials for the Veracode CLI.
`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initialize(cmd)
	},

	Run: func(cmd *cobra.Command, args []string) {
		configure.Run(cmd.Context(), args, cmd.Flags())
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Configure credentials for the Veracode CLI (alias for configure)",
	Long: `Configure credentials for the Veracode CLI.
`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("running initializeConfig()")
		return initializeConfig(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		configure.Run(cmd.Context(), args, cmd.Flags())
	},
}

// sbomCmd represents the sbom command
var configureSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set configuration using the Veracode CLI",
	Long: `Configure credentials for the Veracode CLI.
`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initialize(cmd)
	},

	Run: func(cmd *cobra.Command, args []string) {
		configure.Set(cmd.Context(), args, cmd.Flags())
	},
}

// sbomCmd represents the sbom command
var configureDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete configuration using the Veracode CLI",
	Long: `Delete credentials for the Veracode CLI.
`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initialize(cmd)
	},

	Run: func(cmd *cobra.Command, args []string) {
		configure.Delete(cmd.Context(), args, cmd.Flags())
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)
	rootCmd.AddCommand(initCmd)

	configureCmd.AddCommand(configureSetCmd)
	configureCmd.AddCommand(configureDeleteCmd)
}
