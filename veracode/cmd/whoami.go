/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/veracode/veracode-cli/internal/user"
)

// clearCacheCmd represents the clearCache command
var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Identifies your user int hte Veracode platform",
	Long: `Identifies your user int hte Veracode platform

You will need to have working credentials for this to complete successfully.`,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initialize(cmd)
	},

	Run: func(cmd *cobra.Command, args []string) {
		creds := readCredentials(cmd.Context())

		var u user.User
		err := u.Login(&creds)
		if err != nil {
			panic(err)
		}
		err = u.Validate()

		fmt.Printf("%s (%s)\n", u.Data.UserName, u.Data.UserID)
	},
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}
