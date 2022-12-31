/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package fix

import (
	"context"
	"errors"
	"fmt"
	"log"

	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/veracode/veracode-cli/internal/hmac"

	"github.com/veracode/veracode-cli/internal/intrem"
	intrem_api "github.com/veracode/veracode-cli/internal/api/intrem"
)

// UNUSED allows unused variables to be included in Go programs
func UNUSED(x ...interface{}) {}
//
// Call out to container to execute logic
//
func Run(ctx context.Context, args []string, flags *flag.FlagSet) error {

	fmt.Printf("----------------------------------------------------------------\n")
	fmt.Printf("Veracode Intelligent Remediation service\n")
	fmt.Printf("Version 0.999\n\n")

	err := errors.New("")

	// -- Set up contextual information to pass around
	var fixCtx intrem_api.AppContext
	fixCtx.Debug, err = flags.GetBool("debug")

	// -- Choose which API deployment to use
	/*
	jaroona, err := flags.GetBool("jaroona")
	if jaroona {
		// -- Use the Jaroona test deployment (outside Veracode)
		fixCtx.HttpScheme = "http"
		fixCtx.APIHost    = "a7500143a934442178a717abeeca3abe-124806903.eu-central-1.elb.amazonaws.com:8080" // "18.158.216.3:8080"
		fixCtx.DoAuth     = false
	} else {
		// -- Use the production deployment
		fixCtx.HttpScheme = "https"
		fixCtx.APIHost    = "api.veracode.com"
		fixCtx.DoAuth     = true
	}
	*/

	apihost, err := flags.GetString("apihost")
	fixCtx.HttpScheme = "http"
	fixCtx.APIHost    = apihost
	fixCtx.DoAuth     = true

	// -- Source file to fix
	fixCtx.SourcePath = args[0]

	fixCtx.Reuse, err = flags.GetBool("reuse")
	 if err != nil {
		 return err
	 }

 	apply, err := flags.GetBool("apply")
	fixCtx.Choose = ! apply
	if err != nil {
		return err
	}

	// -- Get the ID of the Issue to fix
	fixCtx.IssueId, err = flags.GetInt("issueid")
	if err != nil {
		return err
	}

	// -- Get the name of the results file (results.json is the default)
	fixCtx.ResultsPath, err = flags.GetString("scanresults")
	if err != nil {
		return err
	}

	// -- Get the credentials...
	credentials := hmac.HmacCredentials{
		viper.GetString("credentials.veracode_api_key_id"),
		viper.GetString("credentials.veracode_api_key_secret"),
	}
	ctxWithCreds := context.WithValue(ctx, "credentials", credentials)

	// -- Do the main work...
	err = intrem.Fix(ctxWithCreds, &fixCtx)

	if err != nil {
		log.Fatal(err)
	}

	return err
}
