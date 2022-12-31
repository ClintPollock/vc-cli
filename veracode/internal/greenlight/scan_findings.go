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
	"time"

	config_ctxt "github.com/veracode/veracode-cli/internal/api/configuration"
	greenlight_api "github.com/veracode/veracode-cli/internal/api/greenlight"
	"github.com/veracode/veracode-cli/internal/hmac"

	"github.com/veracode/veracode-cli/cmd/version"
)

func GetScanFindings(ctx context.Context, app greenlight_api.AppContext, configCtx config_ctxt.Configuration, scan *greenlight_api.ScanInfo) (*greenlight_api.ScanFindings, error) {

	// Format the API URL.
	findingsPath := scan.Links.Findings.Href
	if findingsPath == "" {
		findingsPath = "/scans/" + scan.Id + "/findings"
	}

	apiUrl := url.URL{
		Scheme: configCtx.ApiScheme,
		Host:   configCtx.BaseHostname,
		Path:   configCtx.PipelineUrlPath + findingsPath,
	}
	httpMethod := "GET"

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

	req, err := http.NewRequest(httpMethod, apiUrl.String(), nil)
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
	}

	// Print our request
	if app.Debug {
		requestDump, err := httputil.DumpRequest(req, false)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("\n-----------------------\n getScanFindings request:\n-----------------------\n%s\n", string(requestDump))
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
		fmt.Printf("\n-----------------------\n getScanFindings response:\n-----------------------\n%s\n", responseBody)
	}

	findings := greenlight_api.ScanFindings{}

	if resp.StatusCode == 200 || resp.StatusCode == 202 {

		json.Unmarshal(responseBody, &findings)
		if app.Debug {
			fmt.Printf("Scan Findings status: '%s'.\n", findings.ScanStatus)
		}
		return &findings, err
	} else {
		err = errors.New(fmt.Sprintf("Non-200 status code (%d) when starting a scan.", resp.StatusCode))
		return nil, err
	}

}
