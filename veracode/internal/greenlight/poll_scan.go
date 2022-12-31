package greenlight

/*
 *
 * Static Pipeline basic flow implementation
 *
 */

import (
	"context"
	"fmt"
	"time"

	config_ctxt "github.com/veracode/veracode-cli/internal/api/configuration"
	greenlight_api "github.com/veracode/veracode-cli/internal/api/greenlight"
)

// const scheme = "https"
// const host = "api.veracode.com"
//
// const pipelineUrlPath = "/pipeline_scan/v1"
// const policyUrlPath = "/appsec/v1/policies"
// const createScanPath = "/scans"
//
// const maxSleepTime = 3200
// const maxRetries = 3
// const createScanClientTimeout = 30

func PollForScanCompletion(ctx context.Context, app greenlight_api.AppContext, configCtx config_ctxt.Configuration, scan *greenlight_api.ScanInfo) (*greenlight_api.ScanInfo, error) {

	// scanId := scan.Id

	var err error
	for true {

		scan, err := GetScanDetails(ctx, app, configCtx, scan)

		if err != nil {
			panic(err)
		}
		if app.Debug {
			fmt.Printf("Scan is in a %s status.\n", scan.Status)
		} else {
			if scan.Status == "PENDING" {
				fmt.Print("_")
			}
			if scan.Status == "STARTED" {
				fmt.Print("-")
			}
		}
		if scan.Status != "PENDING" && scan.Status != "STARTED" {
			fmt.Println("|")

			break
		}

		time.Sleep(3 * time.Second)

	}

	return scan, err

}
