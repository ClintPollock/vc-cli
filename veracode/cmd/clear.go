/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/veracode/veracode-cli/internal/cache"
)

// clearCacheCmd represents the clearCache command
var clearCacheCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clears out the Veracode CLI cache",
	Long: `Clears out the Veracode CLI cache, typically at ~/.veracode/cache.

Run this command when you want to force a rerun of scans, and for the tool to
avoid using cached results.`,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initialize(cmd)
	},

	Run: func(cmd *cobra.Command, args []string) {
		cache.Clear()
		fmt.Printf("Cleared cache at %s\n", cache.Path())
	},
}

func init() {
	rootCmd.AddCommand(clearCacheCmd)

}
