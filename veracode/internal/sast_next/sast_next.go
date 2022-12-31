package sast_next

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/r3labs/sse"
	"github.com/wagoodman/go-partybus"

	sast_next_api "github.com/veracode/veracode-cli/internal/api/sast_next"

	config_ctxt "github.com/veracode/veracode-cli/internal/api/configuration"

	greenlight_api "github.com/veracode/veracode-cli/internal/api/greenlight"
	"github.com/veracode/veracode-cli/internal/hmac"

	"github.com/veracode/veracode-cli/cmd/version"
)

const eventsApiTimeout = 300

/*
 *
 * Static Pipeline basic flow implementation
 *
 */

const scan_events_hostname = "stage-static-scan-events-service.era1.dev.vnext.veracode.io"

const findings_events_path = "/v1/api/static-scan-finding-events/"

func AttachFlawEventsToBus(ctx context.Context, app greenlight_api.AppContext, configCtx config_ctxt.Configuration, bus *partybus.Bus, scan *greenlight_api.ScanInfo) error {

	jobLogEvents, err := GetScanJobLogEvents(ctx, app, configCtx, scan)
	if err != nil {
		panic(err)
	}
	jobLifecycleEvents, err := GetScanJobLifecycleEvents(ctx, app, configCtx, scan)
	if err != nil {
		panic(err)
	}
	scanLifecycleEvents, err := GetScanLifecycleEvents(ctx, app, configCtx, scan)
	if err != nil {
		panic(err)
	}
	flawEvents, err := GetScanFindingsEvents(ctx, app, configCtx, scan)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			msg := <-jobLogEvents

			jobEvent := sast_next_api.JobEvent{}
			json.Unmarshal(msg.Data, &jobEvent)

			bus.Publish(partybus.Event{Type: "JobEvent", Value: jobEvent})
		}
	}()

	go func() {
		for {
			msg := <-jobLifecycleEvents

			jobLifecycleEvent := sast_next_api.JobLifecycleEvent{}
			json.Unmarshal(msg.Data, &jobLifecycleEvent)

			bus.Publish(partybus.Event{Type: "JobLifecycleEvent", Value: jobLifecycleEvent})
		}
	}()

	go func() {
		for {
			msg := <-scanLifecycleEvents

			scanLifecycleEvent := sast_next_api.ScanLifecycleEvent{}
			json.Unmarshal(msg.Data, &scanLifecycleEvent)

			bus.Publish(partybus.Event{Type: "ScanLifecycleEvent", Value: scanLifecycleEvent})
			/// enum PENDING, UPLOADING, STARTED, SUCCESS, FAILURE, CANCELLED, TIMEOUT, USER_TIMEOUT

			// if lifecycleEvent.Status == "SUCCESS" || lifecycleEvent.Status == "FAILURE" || lifecycleEvent.Status == "CANCELLED" {
			// 	moveOn = true
			// }
		}
	}()

	if err != nil {
		panic(err)
	}
	go func() {
		for {
			msg := <-flawEvents

			finding := sast_next_api.ScanFindingEvent{}
			json.Unmarshal(msg.Data, &finding)

			bus.Publish(partybus.Event{Type: "ScanFindingEvent", Value: finding})
		}
	}()

	return err
}

// curl --location
//      --request GET "stage-static-scan-events-service.era1.dev.vnext.veracode.io/v1/api/static-scan-finding-events/?filterField=scanId&filterValue=$SCAN_ID" \
//      --header 'Accept: text/event-stream' \
//      --header 'Connection: keep-alive'

func GetScanFindingsEvents(ctx context.Context, app greenlight_api.AppContext, configCtx config_ctxt.Configuration, scan *greenlight_api.ScanInfo) (chan *sse.Event, error) {

	return subscribeToSSE(ctx, app, configCtx, "/static-scan-finding-events", "filterField=scanId&filterValue="+scan.Id)
}

// curl --location --request GET "stage-static-scan-events-service.era1.dev.vnext.veracode.io/v1/api/static-scan-lifecycle-events/?filterField=scanId&filterValue=$SCAN_ID" \
//   --header 'Accept: text/event-stream' \
//   --header 'Connection: keep-alive'

func GetScanLifecycleEvents(ctx context.Context, app greenlight_api.AppContext, configCtx config_ctxt.Configuration, scan *greenlight_api.ScanInfo) (chan *sse.Event, error) {

	return subscribeToSSE(ctx, app, configCtx, "/static-scan-lifecycle-events", "filterField=scanId&filterValue="+scan.Id)

}

// curl --location --request GET "stage-static-scan-events-service.era1.dev.vnext.veracode.io/v1/api/static-scan-job-log-events/?filterField=scanId&filterValue=$SCAN_ID" \
//   --header 'Accept: text/event-stream' \
//   --header 'Connection: keep-alive'

func GetScanJobLogEvents(ctx context.Context, app greenlight_api.AppContext, configCtx config_ctxt.Configuration, scan *greenlight_api.ScanInfo) (chan *sse.Event, error) {

	return subscribeToSSE(ctx, app, configCtx, "/static-scan-job-log-events", "filterField=scanId&filterValue="+scan.Id)

}

// curl --location --request GET "stage-static-scan-events-service.era1.dev.vnext.veracode.io/v1/api/static-scan-job-lifecycle-events/?filterField=scanId&filterValue=$SCAN_ID" \
//   --header 'Accept: text/event-stream' \
//   --header 'Connection: keep-alive'

func GetScanJobLifecycleEvents(ctx context.Context, app greenlight_api.AppContext, configCtx config_ctxt.Configuration, scan *greenlight_api.ScanInfo) (chan *sse.Event, error) {

	return subscribeToSSE(ctx, app, configCtx, "/static-scan-job-lifecycle-events", "filterField=scanId&filterValue="+scan.Id)

}

func subscribeToSSE(ctx context.Context, app greenlight_api.AppContext, configCtx config_ctxt.Configuration, path string, query string) (chan *sse.Event, error) {

	events := make(chan *sse.Event)

	// Format the API URL.
	apiUrl := url.URL{
		Scheme:   configCtx.ApiScheme,
		Host:     configCtx.BaseHostname,
		Path:     configCtx.EventsUrlPath + path,
		RawQuery: query,
	}
	httpMethod := http.MethodGet

	creds := ctx.Value("credentials").(hmac.HmacCredentials)
	authHeader, err := hmac.CalculateAuthorizationHeader(&apiUrl, httpMethod, &creds)
	if err != nil {
		panic(err)
	}

	client := http.Client{}

	// disable TLS checks if forced to
	if configCtx.SSLIgnore {
		unsafeHttpsTransport := http.DefaultTransport.(*http.Transport).Clone()
		unsafeHttpsTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		client = http.Client{
			Transport: unsafeHttpsTransport,
			Timeout:   time.Second * eventsApiTimeout,
		}
	}
	//
	// headers := map[string]string{}
	// headers["Authorization"] = authHeader
	// headers["PLUGIN_VERSION"] = version.App + "_" + version.Version
	// headers["Accept"] = "text/event-stream"
	// headers["Connection"] = "keep-alive"

	// func NewClient(url string) *Client {
	// 	return &Client{
	// 		URL:        url,
	// 		Connection: &http.Client{},
	// 		Headers:    make(map[string]string),
	// 		subscribed: make(map[chan *Event]chan bool),
	// 	}
	// }

	sseClient := sse.NewClient(apiUrl.String())
	sseClient.Connection = &client
	sseClient.Headers["Authorization"] = authHeader
	sseClient.Headers["PLUGIN_VERSION"] = version.App + "_" + version.Version
	sseClient.Headers["Accept"] = "text/event-stream"
	sseClient.Headers["Connection"] = "keep-alive"

	if app.Debug {
		fmt.Printf("\n-----------------------\n subscribeToSSE GET request:\n-----------------------\n%s\n", apiUrl.String())
	}

	sseClient.SubscribeChan("messages", events)
	if app.Debug {
		fmt.Printf("Successfully connected SSE stream to channel")
	}
	if err != nil {
		return nil, err
	}

	return events, nil

}
