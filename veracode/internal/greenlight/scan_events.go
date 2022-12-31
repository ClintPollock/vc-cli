package greenlight

/*
 *
 * Static Pipeline basic flow implementation
 *
 */

// const scheme = "https"
// const host = "api.veracode.com"
//
// const pipelineUrlPath = "/pipeline_scan/v1"
// const policyUrlPath = "/appsec/v1/policies"
// const createScanPath = "/scans"
//
// const maxSleepTime = 3200
// const maxRetries = 3
// // const createScanClientTimeout = 30
//
// func watchFlawEvents(ctx context.Context, app greenlight_api.AppContext, configCtx config_ctxt.Configuration, scan *greenlight_api.ScanInfo) (*greenlight_api.ScanInfo, error) {
//
// 	var moveOn = false
// 	flawEvents, err := sast_next.GetScanFindingsEvents(ctx, app, configCtx, scan)
//
// 	go func() {
// 		for {
// 			msg := <-flawEvents
//
// 			finding := sast_next_api.ScanFindingEvent{}
// 			json.Unmarshal(msg.Data, &finding)
//
// 			color.Set(color.FgYellow)
// 			if app.Debug {
//
// 				var outputJson []byte
// 				outputJson, err = json.MarshalIndent(finding, "", "    ")
// 				fmt.Printf("Finding (%s)\n------------------\n%s\n\n", finding.ID, string(outputJson))
//
// 			} else {
// 				fmt.Printf("Finding (%s) of type \"%s\" with severity %d\n", finding.ID, finding.Issue.IssueTypeID, finding.Issue.Severity)
// 			}
// 			color.Unset()
//
// 		}
// 	}()
//
// 	jobLogEvents, err := sast_next.GetScanJobLogEvents(ctx, app, configCtx, scan)
// 	go func() {
// 		for {
// 			msg := <-jobLogEvents
//
// 			jobEvent := sast_next_api.JobEvent{}
// 			json.Unmarshal(msg.Data, &jobEvent)
//
// 			color.Set(color.FgGreen)
// 			if app.Debug {
//
// 				var outputJson []byte
// 				outputJson, err = json.MarshalIndent(jobEvent, "", "    ")
// 				fmt.Printf("Job event (%s)\n------------------\n%s\n\n", jobEvent.ID, string(outputJson))
//
// 			} else {
// 				fmt.Printf("Job event (%s).\n", jobEvent.ID)
// 			}
// 			color.Unset()
// 		}
// 	}()
//
// 	jobLifecycleEvents, err := sast_next.GetScanJobLifecycleEvents(ctx, app, configCtx, scan)
// 	go func() {
// 		for {
// 			msg := <-jobLifecycleEvents
//
// 			jobLifecycleEvent := sast_next_api.JobLifecycleEvent{}
// 			json.Unmarshal(msg.Data, &jobLifecycleEvent)
//
// 			color.Set(color.FgRed)
// 			if app.Debug {
//
// 				var outputJson []byte
// 				outputJson, err = json.MarshalIndent(jobLifecycleEvent, "", "    ")
// 				fmt.Printf("Job Lifecycle event (%s)\n------------------\n%s\n\n", jobLifecycleEvent.ID, string(outputJson))
//
// 			} else {
// 				fmt.Printf("Job Lifecycle event (%s) with message \"%s\".\n", jobLifecycleEvent.ID, jobLifecycleEvent.Message)
// 			}
// 			color.Unset()
// 		}
// 	}()
//
// 	lifecycleEvents, err := sast_next.GetScanLifecycleEvents(ctx, app, configCtx, scan)
// 	go func() {
// 		for {
// 			msg := <-lifecycleEvents
//
// 			lifecycleEvent := sast_next_api.ScanLifecycleEvent{}
// 			json.Unmarshal(msg.Data, &lifecycleEvent)
//
// 			color.Set(color.FgCyan)
// 			if app.Debug {
//
// 				var outputJson []byte
// 				outputJson, err = json.MarshalIndent(lifecycleEvent, "", "    ")
// 				fmt.Printf("Scan Lifecycle event (%s)\n------------------\n%s\n\n", lifecycleEvent.ID, string(outputJson))
//
// 			} else {
// 				fmt.Printf("Scan Lifecycle event (%s) with message \"%s\".\n", lifecycleEvent.ID, lifecycleEvent.Message)
// 			}
// 			color.Unset()
//
// 			/// enum PENDING, UPLOADING, STARTED, SUCCESS, FAILURE, CANCELLED, TIMEOUT, USER_TIMEOUT
//
// 			if lifecycleEvent.Status == "SUCCESS" || lifecycleEvent.Status == "FAILURE" || lifecycleEvent.Status == "CANCELLED" {
// 				moveOn = true
// 			}
// 		}
// 	}()
//
// 	for !moveOn {
// 		time.Sleep(1 * time.Second)
// 	}
//
// 	return scan, err
//
// }
