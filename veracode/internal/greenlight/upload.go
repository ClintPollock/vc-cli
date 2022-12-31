package greenlight

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/veracode/veracode-cli/cmd/version"
	config_ctxt "github.com/veracode/veracode-cli/internal/api/configuration"
	greenlight_api "github.com/veracode/veracode-cli/internal/api/greenlight"
	"github.com/veracode/veracode-cli/internal/hmac"
)

const maxUploadPartRetries = 5
const uploadExpBackoffBase = 2.0
const uploadClientTimeout = 60

func UploadArtifact(ctx context.Context, app greenlight_api.AppContext, configCtx config_ctxt.Configuration, scan *greenlight_api.ScanInfo, artifact Artifact) (*greenlight_api.ScanInfo, error) {

	creds := ctx.Value("credentials").(hmac.HmacCredentials)

	segmentsExpected := scan.BinarySegmentsExpected
	chunkSize := int(math.Ceil(float64(artifact.size) / float64(segmentsExpected)))

	if segmentsExpected < 1 || chunkSize < 1 {
		err := errors.New(fmt.Sprintf("Upload sizing of the artifact received from upload service is inactionable (number of segments = %d)", segmentsExpected))
		return nil, err
	}
	// maxTimeout := scan.Timeout
	if app.Debug {
		fmt.Printf("Segments Expected: %d \n", segmentsExpected)
		fmt.Printf("Artifact chunk size: %d \n\n", chunkSize)
	}

	file, err := os.Open(artifact.Filename)
	if err != nil {
		fmt.Printf("Error opening file %s for reading.\n", artifact.Filename)
		return nil, err
	}
	defer file.Close()

	part := make([]byte, chunkSize)
	totalBytes := 0
	for i := 1; i <= segmentsExpected; i++ {

		nbytesread, err := file.Read(part)
		totalBytes = totalBytes + nbytesread

		if err != nil {
			if err != io.EOF {
				panic(err)
			}

			break
		}

		if app.Debug {
			fmt.Printf("bytes read %d (total: %d): ", nbytesread, totalBytes)
		}

		// Format the API URL.
		apiUrl := url.URL{
			Scheme: configCtx.ApiScheme,
			Host:   configCtx.BaseHostname,
			Path:   configCtx.PipelineUrlPath + scan.Links.Upload.Href,
		}
		scan, err = uploadArtifactPart(app, configCtx, creds, scan, &apiUrl, artifact.Filename, part)

	}
	if app.Debug {
		fmt.Println("\ntotal bytes read: ", totalBytes)
	}
	return scan, err
}

func uploadArtifactPart(app greenlight_api.AppContext, configCtx config_ctxt.Configuration, creds hmac.HmacCredentials, scan *greenlight_api.ScanInfo, apiUrl *url.URL, filename string, partBytes []byte) (*greenlight_api.ScanInfo, error) {

	// Format the API URL.
	httpMethod := "PUT"
	if app.Debug {
		fmt.Println("\nUploading part to URL: ", apiUrl.String())
	}

	client := http.Client{
		Timeout: time.Second * uploadClientTimeout,
	}
	// disable TLS checks if forced to
	if configCtx.SSLIgnore {
		unsafeHttpsTransport := http.DefaultTransport.(*http.Transport).Clone()
		unsafeHttpsTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		client = http.Client{
			Transport: unsafeHttpsTransport,
			Timeout:   time.Second * uploadClientTimeout,
		}
	}

	// -----------------------------------

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}

	part.Write(partBytes)

	//_ = writer.WriteField("FOO", "BAR")
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", apiUrl.String(), body)
	if err != nil {
		fmt.Println("PUT ERROR : ", err)

		panic(err)
	}

	req.Header = http.Header{
		"PLUGIN_VERSION": {version.App + "_" + version.Version},
		"Content-Type":   {writer.FormDataContentType()},
	}

	// FOR MAX_NUMBER_OF_RETRIES { ... }
	for tryN := 1; tryN <= maxUploadPartRetries; tryN++ {

		// Calculate auth header
		authHeader, err := hmac.CalculateAuthorizationHeader(apiUrl, httpMethod, &creds)
		if err != nil {
			panic(err)
		}
		req.Header.Add("Authorization", authHeader)

		// requestDump, err := httputil.DumpRequest(req, false)
		// if err != nil {
		// 	fmt.Println(err)
		// }
		//fmt.Printf("\n-----------------------\n uploadArtifactPart request:\n-----------------------\n%s\n", string(requestDump))

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Client %s error on upload: ", httpMethod)
			panic(err)
		}

		defer resp.Body.Close()

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		//fmt.Printf("\n-----------------------\n uploadArtifactPart response:\n-----------------------\n%s\n", respBody)

		createScanResponse := greenlight_api.ScanInfo{}

		if resp.StatusCode == 200 {

			json.Unmarshal(respBody, &createScanResponse)
			return &createScanResponse, err

		} else {

			if resp.StatusCode != 504 && resp.StatusCode != 500 && resp.StatusCode != 403 {
				err = errors.New(fmt.Sprintf("Non-200 status code (%d) when creating a scan.", resp.StatusCode))
				return nil, err
			}
		}

		// Exponential backoff
		sleepDuration := time.Duration(math.Pow(uploadExpBackoffBase, float64(tryN-1))) * time.Second
		time.Sleep(sleepDuration)

	}

	return nil, err
}
