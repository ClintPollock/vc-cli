package greenlight

/*
 *
 * Static Pipeline basic flow implementation
 *
 */

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	config_ctxt "github.com/veracode/veracode-cli/internal/api/configuration"
	greenlight_api "github.com/veracode/veracode-cli/internal/api/greenlight"
	"github.com/veracode/veracode-cli/internal/hmac"

	"github.com/veracode/veracode-cli/cmd/version"
)

func StartScan(ctx context.Context, app greenlight_api.AppContext, configCtx config_ctxt.Configuration, scan *greenlight_api.ScanInfo) (*greenlight_api.ScanInfo, error) {

	// ActionResults startResult;
	//
	// String jsonBody = "{\"scan_status\": \"STARTED\"}";
	//
	// AuthenticatedRequest startRequest =
	//     this.executeAPICommand(startUri, AuthenticatedRequestBuilder.PUT, jsonBody);
	// startResult = new ActionResults(startRequest.getResponseCode(), startRequest.getResponseString());
	//
	// return startResult;

	// Format the API URL.
	startPath := scan.Links.Start.Href
	if startPath == "" {
		startPath = "/scans/" + scan.Id
	}

	apiUrl := url.URL{
		Scheme: configCtx.ApiScheme,
		Host:   configCtx.BaseHostname,
		Path:   configCtx.PipelineUrlPath + startPath,
	}
	httpMethod := "PUT"

	client := http.Client{
		Timeout: time.Second * scanAPIClientTimeout,
	}
	// disable TLS checks if forced to
	if configCtx.SSLIgnore {
		unsafeHttpsTransport := http.DefaultTransport.(*http.Transport).Clone()
		unsafeHttpsTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		client = http.Client{
			Transport: unsafeHttpsTransport,
			Timeout:   time.Second * scanAPIClientTimeout,
		}
	}

	requestBody := "{\"scan_status\": \"STARTED\"}"
	requestBodyReader := strings.NewReader(requestBody)
	req, err := http.NewRequest(httpMethod, apiUrl.String(), requestBodyReader)
	if err != nil {
		return nil, err
	}

	creds := ctx.Value("credentials").(hmac.HmacCredentials)
	authHeader, err := hmac.CalculateAuthorizationHeader(&apiUrl, httpMethod, &creds)
	if err != nil {
		panic(err)
	}

	req.Header = http.Header{
		"Authorization":  {authHeader},
		"PLUGIN_VERSION": {version.App + "_" + version.Version},
		"Content-Type":   {"application/json"},
	}

	// Print our request
	if app.Debug {
		requestDump, err := httputil.DumpRequest(req, true)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("\n-----------------------\n startScan request:\n-----------------------\n%s\n", string(requestDump))
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if app.Debug {
		fmt.Printf("\n-----------------------\n startScan response:\n-----------------------\n%s\n", responseBody)
	}

	startScanResponse := greenlight_api.ScanInfo{}

	if resp.StatusCode == 200 || resp.StatusCode == 202 {

		json.Unmarshal(responseBody, &startScanResponse)
		if app.Debug {
			fmt.Printf("Start scan status: '%s'.\n", startScanResponse.Status)
		}
		return &startScanResponse, err
	} else {
		err = errors.New(fmt.Sprintf("Non-200 status code (%d) when starting a scan.", resp.StatusCode))
		return nil, err
	}

}
