package greenlight

/*
 *
 * Static Pipeline basic flow implementation
 *
 */

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	config_ctxt "github.com/veracode/veracode-cli/internal/api/configuration"
	greenlight_api "github.com/veracode/veracode-cli/internal/api/greenlight"
	"github.com/veracode/veracode-cli/internal/hmac"

	"github.com/veracode/veracode-cli/cmd/version"
)

const createScanPath = "/scans"

// const maxSleepTime = 3200
// const maxRetries = 3
// const createScanClientTimeout = 30
//
// //
// Generate the data needed to kickoff a createScan request
//
func InitializeScanInfo(ctx context.Context, app greenlight_api.AppContext, configCtx config_ctxt.Configuration, artifact Artifact) (*greenlight_api.ScanInfo, error) {

	scanRequest := greenlight_api.ScanInfo{
		AppID:         app.AppId,
		BinaryHash:    artifact.sha256,
		BinaryName:    artifact.Filename,
		BinarySize:    artifact.size,
		Stage:         app.Stage,
		PluginVersion: version.App + "_" + version.Version,

		ProjectName: app.ProjectName,
		ProjectRef:  app.ProjectRef,
		ProjectURI:  app.ProjectUrl,

		Timeout: 10,
	}

	if app.EmitStackDump {
		scanRequest.EmitStackDump = "true"
	} else {
		scanRequest.EmitStackDump = "false"
	}

	if app.ResumeScanId != "" {
		scanRequest.Id = app.ResumeScanId
		if app.Debug {
			fmt.Printf("Resume scan with ID = %s.\n", app.ResumeScanId)
		}
	}

	return &scanRequest, nil
}

func CreateScan(ctx context.Context, app greenlight_api.AppContext, configCtx config_ctxt.Configuration, scanRequest *greenlight_api.ScanInfo) (*greenlight_api.ScanInfo, error) {

	// Format the API URL.
	apiUrl := url.URL{
		Scheme: configCtx.ApiScheme,
		Host:   configCtx.BaseHostname,
		Path:   configCtx.PipelineUrlPath + createScanPath,
	}
	httpMethod := "POST"

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

	requestBody, err := json.Marshal(scanRequest)
	requestBodyReader := bytes.NewReader(requestBody)
	req, err := http.NewRequest(httpMethod, apiUrl.String(), requestBodyReader)
	if err != nil {
		fmt.Println("Error in createScan creating a new request ", err)
		return nil, err
	}

	creds := ctx.Value("credentials").(hmac.HmacCredentials)
	authHeader, err := hmac.CalculateAuthorizationHeader(&apiUrl, httpMethod, &creds)
	if err != nil {
		fmt.Println("Error in createScan calculating HMAC headers ", err)
		panic(err)
	}

	req.Header = http.Header{
		"Authorization":  {authHeader},
		"PLUGIN_VERSION": {version.App + "_" + version.Version},
		"Content-Type":   {"application/json"},
	}

	req.Header.Set("x-mock-number-of-scanner-instances", configCtx.Mock.NumberOfScannerInstances)
	req.Header.Set("x-mock-number-of-findings-to-throw", configCtx.Mock.NumberOfFindingsToThrow)
	req.Header.Set("x-mock-delay-in-ms-between-findings", configCtx.Mock.DelayInMsBetweenFindings)
	req.Header.Set("x-mock-run-scanners-in-parallel", configCtx.Mock.RunScannersInParallel)

	// Print our request
	if app.Debug {
		requestDump, err := httputil.DumpRequest(req, true)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("\n-----------------------\n createScan request:\n-----------------------\n%s\n\n", requestDump)
	}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error in createScan doing REST call ", err)
		return nil, err
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if app.Debug {
		fmt.Printf("\n-----------------------\n createScan response:\n-----------------------\n%s\n\n", responseBody)
	}

	createScanResponse := greenlight_api.ScanInfo{}

	if resp.StatusCode == 200 {

		json.Unmarshal(responseBody, &createScanResponse)
		if app.Debug {
			fmt.Printf("Created scan expects: '%d' segments.\n", createScanResponse.BinarySegmentsExpected)
		}

	} else {
		err = errors.New(fmt.Sprintf("Non-200 status code (%d) when creating a scan.", resp.StatusCode))
	}
	return &createScanResponse, err

}
