/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package research

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
	var cmd = "research"
	logUUID := uuid.Generate()

	cntrArgs := []string{cmd}

	flagsStr := parseOutFlagsForContainer(flags)

	userData := initialize.Validate(cmd, flagsStr, logUUID, nil)

	initialize.InitializeDocker()
	initialize.InitializeCache()

	cntrArgs = append(append(cntrArgs, args...), flagsStr...)

	out, retval, err := verascanner.ContainerRun(ctx, cntrArgs)
	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	if err != nil {
		panic(err)
	}

	initialize.ValidateWithCreds(cmd, flagsStr, logUUID, userData)

	if retval > 0 {
		os.Exit(int(retval))
	}

	return err
}

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

	val, err := flags.GetString("format")
	if err == nil {
		cntrArgs = append([]string{"--format", val}, cntrArgs...)
	}

	val, err = flags.GetString("out")
	if err == nil {
		cntrArgs = append([]string{"--out", val}, cntrArgs...)
	}

	return cntrArgs

}
