/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package inspect

import (
	"context"
	"os"

	"github.com/docker/distribution/uuid"

	flag "github.com/spf13/pflag"

	"github.com/docker/docker/pkg/stdcopy"
	"github.com/veracode/veracode-cli/internal/api/identity/initialize"
	"github.com/veracode/veracode-cli/internal/verascanner"
)

// Call out to container to execute logic
func Run(ctx context.Context, args []string, flags *flag.FlagSet) error {
	var cmd = "inspect"
	logUUID := uuid.Generate()

	cntrArgs := []string{cmd}
	flagsStr := parseOutFlagsForContainer(flags)

	userData := initialize.Validate(cmd, flagsStr, logUUID, nil)

	initialize.InitializeDocker()
	initialize.InitializeCache()

	cntrArgs = append(append(cntrArgs, args...), flagsStr...)

	out, _, err := verascanner.ContainerRun(ctx, cntrArgs)
	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	if err != nil {
		panic(err)
	}

	initialize.ValidateWithCreds(cmd, flagsStr, logUUID, userData)

	return err
}

// Used to pass thru some CLI flags to the underlying container call
func parseOutFlagsForContainer(flags *flag.FlagSet) []string {
	cntrArgs := []string{}

	bval, err := flags.GetBool("verbose")
	if err == nil {
		if bval {
			cntrArgs = append([]string{"--verbose"}, cntrArgs...)
		}
	}

	bval, err = flags.GetBool("prettyprint")
	if err == nil {
		if bval {
			cntrArgs = append([]string{"--prettyprint"}, cntrArgs...)
		}
	}

	val, err := flags.GetString("out")
	if err == nil {
		cntrArgs = append([]string{"--out", val}, cntrArgs...)
	}

	return cntrArgs

}
