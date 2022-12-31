/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package sast

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

	config_ctxt "github.com/veracode/veracode-cli/internal/api/configuration"
	app_context "github.com/veracode/veracode-cli/internal/api/greenlight"
	"github.com/veracode/veracode-cli/internal/cli"
	"github.com/veracode/veracode-cli/internal/hmac"
)

//
// Call out to container to execute logic
//
func Run(ctx context.Context, args []string, flags *flag.FlagSet) error {

	credentials := hmac.HmacCredentials{
		viper.GetString("credentials.veracode_api_key_id"),
		viper.GetString("credentials.veracode_api_key_secret"),
	}

	err := errors.New("")

	appContext := app_context.AppContext{}

	appContext.ResumeScanId, err = flags.GetString("scan-id")
	appContext.ProjectName, err = flags.GetString("project-name")
	appContext.ProjectUrl, err = flags.GetString("project-url")
	appContext.ProjectRef, err = flags.GetString("project-ref")
	appContext.AppId, err = flags.GetString("app-id")

	stage, err := flags.GetString("stage")
	stage = strings.ToUpper(stage) // Must be DEVELOPMENT', 'TESTING', 'RELEASE
	if !(stage == "DEVELOPMENT" || stage == "TESTING" || stage == "RELEASE") {
		err = errors.New("Scan stage parameter must be DEVELOPMENT, TESTING, or RELEASE")
		fmt.Println("ERROR parsing parameters: ", err)
		return err
	}
	appContext.Stage = stage

	appContext.Debug, err = flags.GetBool("debug")
	appContext.ShowDetails, err = flags.GetBool("show-details")

	appContext.FailOnSeverity, err = flags.GetFloat32("fail-on-severity")
	appContext.FailFast, err = flags.GetBool("fail-fast")
	appContext.BaselineFile, err = flags.GetString("baseline-file")

	appContext.EmitStackDump, err = flags.GetBool("emit-stack-dump")

	failOnCWE, err := flags.GetString("fail-on-cwe")
	appContext.FailOnCWEs = strings.Split(strings.ReplaceAll(failOnCWE, " ", ""), ",")

	config := config_ctxt.Configuration{}

	config.ApiScheme = viper.GetString("urls.scheme")
	config.BaseHostname = viper.GetString("urls.host")
	config.PipelineUrlPath = viper.GetString("urls.pipeline_api_path")
	config.EventsUrlPath = viper.GetString("urls.events_api_path")
	config.PolicyUrlPath = viper.GetString("urls.policy_api_path")
	config.SSLIgnore = viper.GetBool("urls.ssl-ignore")

	config.SAST.UseRealTimeFlawAPIs = viper.GetBool("sast.use_realtime_flaw_apis")

	config.Mock.NumberOfScannerInstances = viper.GetString("mock.number-of-scanner-instances")
	config.Mock.NumberOfFindingsToThrow = viper.GetString("mock.number-of-findings-to-throw")
	config.Mock.DelayInMsBetweenFindings = viper.GetString("mock.delay-in-ms-between-findings")
	config.Mock.RunScannersInParallel = viper.GetString("mock.run-scanners-in-parallel")
	config.Mock.SkipUpload = viper.GetBool("mock.skip-upload")

	sastCtx := context.WithValue(ctx, "credentials", credentials)
	findings, err := cli.Scan(sastCtx, appContext, config, args[0])
	if err != nil {
		panic(err)
	}

	fmt.Printf("Scan found %d findings.\n", len(findings.Findings))

	doPrettyPrint, err := flags.GetBool("prettyprint")

	var outputJson []byte
	if doPrettyPrint {
		outputJson, err = json.MarshalIndent(findings, "", "    ")
	} else {
		outputJson, err = json.Marshal(findings)
	}

	outfile, err := flags.GetString("out")
	if err != nil {
		panic(err)
	}

	if outfile == "-" {
		fmt.Println(string(outputJson))
	} else {
		f, err := os.Create(outfile)

		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		_, err = f.Write(outputJson)
		if err != nil {
			panic(err)
		}
	}
	return err
}
