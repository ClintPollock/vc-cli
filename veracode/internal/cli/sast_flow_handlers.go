package cli

/*
 *
 * Static Pipeline basic flow implementation
 *
 */

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/wagoodman/go-partybus"

	greenlight_api "github.com/veracode/veracode-cli/internal/api/greenlight"
	sast_next_api "github.com/veracode/veracode-cli/internal/api/sast_next"
)

//
// CLI-based flow through scan process using multiple APIs
//
//

type InfoEvent struct {
	message string
	err     error
}

type ScanFlowEvent struct {
	action  string
	message string
}

func log(bus *partybus.Bus, msg string) {
	bus.Publish(partybus.Event{Type: "InfoEvent", Value: InfoEvent{msg, nil}})
}
func logError(bus *partybus.Bus, msg string, err error) {
	bus.Publish(partybus.Event{Type: "InfoEvent", Value: InfoEvent{msg, err}})
}

func handleEvents(app greenlight_api.AppContext, bus *partybus.Bus) {

	subscription := bus.Subscribe()

	for event := range subscription.Events() {
		//fmt.Printf("Event: %+v\n", event)
		switch event.Value.(type) {

		case string:
			fmt.Println("STRING", event.Value.(string))

		case InfoEvent:
			fmt.Printf("LOG Message: %s\n", event.Value.(InfoEvent).message)

		case ScanFlowEvent:
			sfe := event.Value.(ScanFlowEvent)
			color.Set(color.FgBlack, color.BgWhite, color.Bold)
			fmt.Printf("Flow Event (%s) %s", sfe.action, sfe.message)
			color.Unset()
			fmt.Println("")

			if sfe.action == "STOP_WAIT" {
				scanComplete = true
			}

		case sast_next_api.ScanFindingEvent:

			sfe := event.Value.(sast_next_api.ScanFindingEvent)
			cweID := sfe.Issue.CweID
			severity := sfe.Issue.Severity

			passesPolicy := true
			duplicate := false

			if app.BaselineFindings != nil {
				for _, issue := range app.BaselineFindings.Findings {

					if fmt.Sprintf("%d", sfe.Issue.CweID) == issue.CweID &&
						fmt.Sprintf("%d", sfe.Issue.FlawMatch.FlawHash) == issue.FlawMatch.FlawHash {
						duplicate = true
					}
				}
			}
			color.Set(color.FgYellow)
			if app.Debug {

				var outputJson []byte
				//var err error
				outputJson, _ = json.MarshalIndent(sfe, "", "    ")
				fmt.Printf("Finding (%s)\n------------------\n%s", sfe.ID, string(outputJson))

			} else {
				if !duplicate {
					if app.ShowDetails {

						fileString := fmt.Sprintf("[%s]%s:%d", "", sfe.Issue.Files.SourceFile.File, sfe.Issue.Files.SourceFile.Line)
						fmt.Printf("Finding: (cweId: CWE-%d, type: %s, severity: %d, file: %s, id: %s) -- %s.", sfe.Issue.CweID, sfe.Issue.IssueTypeID, sfe.Issue.Severity, fileString, sfe.ID, sfe.Issue.IssueType)
					} else {
						fmt.Printf("Finding: (cweId: CWE-%d, type: %s, severity: %d).", sfe.Issue.CweID, sfe.Issue.IssueTypeID, sfe.Issue.Severity)
					}
				} else {
					color.Set(color.FgWhite, color.Faint)
					fmt.Printf("Duplicate Finding: (cweId: CWE-%d, type: %s, severity: %d).", sfe.Issue.CweID, sfe.Issue.IssueTypeID, sfe.Issue.Severity)
				}
			}
			color.Unset()
			fmt.Println()

			//
			// Evaluate simplisitic policy
			//
			if !duplicate {
				if float32(severity) >= app.FailOnSeverity {
					passesPolicy = false
				}
				for _, cwe := range app.FailOnCWEs {
					if fmt.Sprintf("%d", cweID) == cwe {
						passesPolicy = false
					}
				}

				if !passesPolicy {
					color.Set(color.FgBlack, color.BgYellow, color.Bold)
					fmt.Printf("Finding: (id: %s,  cweId: %d, severity: %d) FAILS policy.", sfe.ID, sfe.Issue.CweID, sfe.Issue.Severity)
					color.Unset()
					fmt.Println()

					if app.FailFast {
						color.Set(color.FgRed, color.BgWhite, color.Bold)
						fmt.Printf(" [EXITING] ")
						color.Unset()
						fmt.Println()

						os.Exit(-1)
					}
				}

			}

		case sast_next_api.JobEvent:

			je := event.Value.(sast_next_api.JobEvent)
			color.Set(color.FgGreen)
			if app.Debug {

				var outputJson []byte
				outputJson, _ = json.MarshalIndent(je, "", "    ")
				fmt.Printf("Job event (%s)\n------------------\n%s\n\n", je.ID, string(outputJson))

			} else {
				fmt.Printf("Job event (%s).\n", je.ID)
			}
			color.Unset()

		case sast_next_api.JobLifecycleEvent:

			jle := event.Value.(sast_next_api.JobLifecycleEvent)

			color.Set(color.FgRed)
			if app.Debug {
				outputJson, _ := json.MarshalIndent(jle, "", "    ")
				fmt.Printf("Job Lifecycle event (%s)\n------------------\n%s\n\n", jle.ID, string(outputJson))
			} else {
				fmt.Printf("Job Lifecycle event (id: %s, status: %s) with message \"%s\".\n", jle.ID, jle.Status, jle.Message)
			}
			color.Unset()

		case sast_next_api.ScanLifecycleEvent:

			sle := event.Value.(sast_next_api.ScanLifecycleEvent)

			color.Set(color.FgCyan)
			if app.Debug {
				outputJson, _ := json.MarshalIndent(sle, "", "    ")
				fmt.Printf("Scan Lifecycle event (%s)\n------------------\n%s\n\n", sle.ID, string(outputJson))
			} else {
				fmt.Printf("Scan Lifecycle event (id: %s, status: %s) with message \"%s\".\n", sle.ID, sle.Status, sle.Message)
			}
			color.Unset()

			if sle.Status == "SUCCESS" || sle.Status == "FAILURE" || sle.Status == "CANCELLED" {

				bus.Publish(partybus.Event{Type: "ScanFlowEvent", Value: ScanFlowEvent{"STOP_WAIT", sle.Status}})
			}

		default:
			fmt.Println("event occurred")
		}
	}

}
