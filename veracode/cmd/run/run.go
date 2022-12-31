/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package run

import (
	"context"
	"os"

	"github.com/docker/distribution/uuid"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/veracode/veracode-cli/internal/api/identity/initialize"
	"github.com/veracode/veracode-cli/internal/verascanner"
)

// Call out to container to execute logic
func Run(ctx context.Context, args []string) error {
	var cmd = "run"
	logUUID := uuid.Generate()

	userData := initialize.Validate(cmd, nil, logUUID, nil)

	initialize.InitializeDocker()
	initialize.InitializeCache()

	cntrArgs := []string{cmd}
	cntrArgs = append(cntrArgs, args...)

	out, _, err := verascanner.ContainerRun(ctx, cntrArgs)
	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	if err != nil {
		panic(err)
	}

	initialize.ValidateWithCreds(cmd, nil, logUUID, userData)

	return err
}
