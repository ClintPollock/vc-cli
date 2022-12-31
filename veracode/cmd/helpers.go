/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*
* Various initialization and validation functions
*
*/
package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/veracode/veracode-cli/internal/cache"
	"github.com/veracode/veracode-cli/internal/verascanner"
)

const ENV_PREFIX = "VERACODE"
const CONFIG_FILENAME = "config"
const CONFIG_LOCAL_DIR = ".veracode"

//
// Generic initializer
//
func initialize(cmd *cobra.Command) error {
	initializeCache(cmd.Context())
	return initializeConfig(cmd)
}

// This references https://github.com/carolynvs/stingoftheviper .
// The approach helps us achieve the following precedence. flags > env variables >  config variables > default flag values
func initializeConfig(cmd *cobra.Command) error {

	viper.SetConfigType("yaml")
	viper.SetConfigName(CONFIG_FILENAME)
	viper.AddConfigPath("$HOME/" + CONFIG_LOCAL_DIR)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading in veracode config as yaml:", err)
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	viper.SetEnvPrefix(ENV_PREFIX) // This allows us to support VERACODE_<variable> environment variable pattern
	viper.AutomaticEnv()
	bindFlags(cmd)

	return nil
}

// Bind each cobra flag to its associated viper configuration (config file and environment variable)
func bindFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores, e.g. --favorite-color to STING_FAVORITE_COLOR
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			viper.BindEnv(f.Name, fmt.Sprintf("%s_%s", ENV_PREFIX, envVarSuffix))
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && viper.IsSet(f.Name) {
			val := viper.Get(f.Name)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}

// Pull veracode scanner image
func initializeDocker(ctx context.Context) error {
	return verascanner.ImagePull()
}

// creates cache path
func initializeCache(ctx context.Context) {
	_ = cache.Path()
}
