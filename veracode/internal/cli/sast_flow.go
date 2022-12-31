package cli

/*
 *
 * Static Pipeline basic flow implementation
 *
 */

import (
	"context"
	"fmt"
	"time"

	"github.com/wagoodman/go-partybus"

	config_ctxt "github.com/veracode/veracode-cli/internal/api/configuration"
	greenlight_api "github.com/veracode/veracode-cli/internal/api/greenlight"
	"github.com/veracode/veracode-cli/internal/greenlight"

	"github.com/veracode/veracode-cli/internal/sast_next"
)

const maxSleepTime = 3200
const maxRetries = 3
const createScanClientTimeout = 30

//
// CLI-based flow through scan process using multiple APIs
//
//

// Stupid semaphore for now
var scanComplete bool

func Scan(ctx context.Context, app greenlight_api.AppContext, config config_ctxt.Configuration, filename string) (*greenlight_api.ScanFindings, error) {

	var err error

	bus := partybus.NewBus()

	//
	// Load baseline
	//
	app.BaselineFindings, err = greenlight.LoadBaselineFindings(app)
	go handleEvents(app, bus)

	if err != nil {
		fmt.Println("Encountered an issue opening the baseline file (will be ignored): ", err)
	}

	//
	// Sequentially setup a scan, upload and start
	//
	scan, err := initiateScan(ctx, app, config, bus, filename)

	//
	// Wait for completion using getDetails or watch scan events
	//
	if config.SAST.UseRealTimeFlawAPIs {

		//_, err = watchFlawEvents(ctx, appCtx, configCtx, bus, scan)
		err := sast_next.AttachFlawEventsToBus(ctx, app, config, bus, scan)
		if err != nil {
			logError(bus, "Error while subscribing to realtime flaw events", err)
			return nil, err
		}

		// Cheap wait for completion step.
		for !scanComplete {
			time.Sleep(1 * time.Second)
		}
		//scan, err = greenlight.GetScanDetails(ctx, app, config, scan)

	} else {
		//fmt.Printf("Monitoring scan (%s) status ...\n", scan.Id)
		log(bus, fmt.Sprintf("Monitoring scan (%s) status ...", scan.Id))
		scan, err = greenlight.PollForScanCompletion(ctx, app, config, scan)
	}

	//
	// Once
	//
	scan, err = greenlight.GetScanDetails(ctx, app, config, scan)

	if scan.Status == "FAILURE" {
		return nil, fmt.Errorf("Scan (%s) for file %s failed with message: %s", scan.Id, scan.BinaryName, scan.Message)
	}
	if scan.Status == "" {
		return nil, fmt.Errorf("Scan (%s) not found.", scan.Id)
	}

	if err != nil {
		return nil, err
	}
	//
	// Findings
	//
	scanFindings, err := greenlight.GetScanFindings(ctx, app, config, scan)
	if err != nil {
		panic(err)
	}

	//
	// Filter findings based on baseline file
	//
	if app.BaselineFindings != nil {

		unfilteredScanFindings := scanFindings
		scanFindings.Findings = []greenlight_api.Finding{}

		for _, newf := range unfilteredScanFindings.Findings {
			for _, bf := range app.BaselineFindings.Findings {

				bId := bf.IssueID
				bHash := bf.FlawMatch.FlawHash
				nId := newf.IssueID
				nHash := newf.FlawMatch.FlawHash

				if nId != bId || nHash != bHash {
					scanFindings.Findings = append(scanFindings.Findings, newf)
				}

			}

		}
	}

	return scanFindings, err
}

// Synchronous routine to setup a scan
func initiateScan(ctx context.Context, appCtx greenlight_api.AppContext, configCtx config_ctxt.Configuration, bus *partybus.Bus, filename string) (*greenlight_api.ScanInfo, error) {

	artifact := greenlight.Artifact{}
	artifact.Load(filename)

	scan, err := greenlight.InitializeScanInfo(ctx, appCtx, configCtx, artifact)

	// If this is a new scan
	if scan.Id == "" {
		//
		// createScan
		//
		scan, err = greenlight.CreateScan(ctx, appCtx, configCtx, scan)
		bus.Publish(partybus.Event{Type: "ScanFlowEvent", Value: ScanFlowEvent{"CREATED", scan.Id}})

		if err != nil {
			logError(bus, "Error while creating scan", err)
			panic(err)
		}

		log(bus, fmt.Sprintf("Created scan with ID: '%s'.", scan.Id))

		//
		// uploadScan
		//
		//fmt.Printf("configCtx.Mock.SkipUpload = %t.\n", configCtx.Mock.SkipUpload)
		if !configCtx.Mock.SkipUpload {
			log(bus, fmt.Sprintf("Uploading file: '%s'.", artifact.Filename))

			scan, err = greenlight.UploadArtifact(ctx, appCtx, configCtx, scan, artifact)
			if err != nil {
				logError(bus, "Error while uploading binary", err)
				panic(err)
			}

			bus.Publish(partybus.Event{Type: "ScanFlowEvent", Value: ScanFlowEvent{"UPLOADED", scan.Id}})

		}
		//
		// startScan
		//
		//fmt.Printf("Starting scan '%s'.\n", scan.Id)

		scan, err = greenlight.StartScan(ctx, appCtx, configCtx, scan)
		if err != nil {
			panic(err)
		}
		log(bus, fmt.Sprintf("Starting scan '%s'.", scan.Id))
		bus.Publish(partybus.Event{Type: "ScanFlowEvent", Value: ScanFlowEvent{"STARTED", scan.Id}})

	}

	return scan, err
}
