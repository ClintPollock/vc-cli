/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"errors"
	"os"

	"github.com/docker/docker/pkg/stdcopy"
	"github.com/spf13/viper"
	"github.com/veracode/veracode-cli/internal/hmac"
	"github.com/veracode/veracode-cli/internal/user"
	"github.com/veracode/veracode-cli/internal/verascanner"
)

func validate(ctx context.Context) {
	validateCredentials(ctx)
	sanityCheckDocker(ctx)
}

func sanityCheckDocker(ctx context.Context) {
	// Run a simple echo command to pull and validate environment
	strArr := []string{"run", "bash", "echo", ""}
	out, _, err := verascanner.ContainerRun(ctx, strArr)
	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	if err != nil {
		panic(err)
	}
}

func validateCredentials(ctx context.Context) {

	var err error
	creds := readCredentials(ctx)

	if creds.Id == "" || creds.Secret == "" {
		err := errors.New("Empty credentials read from configuration")
		panic(err)
	}

	// Try to login with these credentials
	var u user.User
	err = u.Login(&creds)
	if err != nil {
		panic(err)
	}
	err = u.Validate()

	if err != nil {
		panic(err)
	}
}

func readCredentials(ctx context.Context) hmac.HmacCredentials {

	creds := hmac.HmacCredentials{}

	creds.Id = viper.GetString("credentials.veracode_api_key_id")
	creds.Secret = viper.GetString("credentials.veracode_api_key_secret")

	return creds
}

func validateFromCredentials(creds *hmac.HmacCredentials) (*user.User, error) {
	// Verify user via Identity API.
	var u user.User

	err := u.Login(creds)
	if err != nil {
		panic(err)
	}
	err = u.Validate()

	return &u, err
}
