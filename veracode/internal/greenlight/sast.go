package greenlight

/*
 *
 * Static Pipeline basic flow implementation
 *
 */

const maxSleepTime = 3200
const maxRetries = 3

const scanAPIClientTimeout = 30

//
// Emulate current flow
//
//  https://gitlab.laputa.veracode.io/sms/pipeline-scan/pipeline-scan/-/blob/master/src/main/java/com/veracode/greenlight/tools/scanner/HttpDevopsEngineRunner.java#L97
// //
// func Scan(ctx context.Context, appCtx greenlight_api.AppContext, configCtx config_ctxt.Configuration, filename string) (*greenlight_api.ScanFindings, error) {
//
// 	artifact := Artifact{}
// 	artifact.Load(filename)
//
// 	scan, err := InitializeScanInfo(ctx, appCtx, configCtx, artifact)
//
// 	// If this is a new scan
// 	if scan.Id == "" {
// 		//
// 		// createScan
// 		//
// 		scan, err = CreateScan(ctx, appCtx, configCtx, scan)
// 		if err != nil {
// 			panic(err)
// 		}
// 		fmt.Printf("Created scan with ID: '%s'.\n", scan.Id)
//
// 		//
// 		// uploadScan
// 		//
// 		//fmt.Printf("configCtx.Mock.SkipUpload = %t.\n", configCtx.Mock.SkipUpload)
// 		if !configCtx.Mock.SkipUpload {
// 			fmt.Printf("Uploading file: '%s'.\n", artifact.Filename)
// 			scan, err = UploadArtifact(ctx, appCtx, configCtx, scan, artifact)
// 			if err != nil {
// 				panic(err)
// 			}
// 		}
// 		//
// 		// startScan
// 		//
// 		fmt.Printf("Starting scan '%s'.\n", scan.Id)
// 		scan, err = StartScan(ctx, appCtx, configCtx, scan)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}
//
// 	//
// 	// Wait for completion using getDetails
// 	//
// 	//fmt.Printf("Monitoring scan (%s) status ...\n", scan.Id)
// 	//scan, err = pollForScanCompletion(ctx, appCtx, configCtx, scan)
// 	_, err = watchFlawEvents(ctx, appCtx, configCtx, scan)
//
// 	scan, err = PollForScanCompletion(ctx, appCtx, configCtx, scan)
// 	if scan.Status == "FAILURE" {
// 		return nil, fmt.Errorf("Scan (%s) for file %s failed with message: %s", scan.Id, scan.BinaryName, scan.Message)
// 	}
// 	if scan.Status == "" {
// 		return nil, fmt.Errorf("Scan (%s) not found.", scan.Id)
// 	}
//
// 	if err != nil {
// 		return nil, err
// 	}
// 	//
// 	// Findings
// 	//
// 	findings, err := GetScanFindings(ctx, appCtx, configCtx, scan)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	return findings, err
// }
