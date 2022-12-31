package intrem

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
  "io"
  "log"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
  "mime/multipart"
  //"net/textproto"
	"path/filepath"
	"os"
	"os/exec"

	greenlight_api "github.com/veracode/veracode-cli/internal/api/greenlight"
  intrem_api "github.com/veracode/veracode-cli/internal/api/intrem"
	"github.com/veracode/veracode-cli/internal/hmac"
	"github.com/veracode/veracode-cli/cmd/version"

	// "github.com/sergi/go-diff/diffmatchpatch"
	// "github.com/bluekeyes/go-gitdiff/gitdiff"
)

// UNUSED allows unused variables to be included in Go programs
func UNUSED(x ...interface{}) {}

const intremUrlPath = "/intrem/v1"
const policyUrlPath = "/appsec/v1/policies"

const sleepTime = 6
const maxRetries = 100

const requestTimeout = 30

// -----------------------------------------------------------------------------
//
// Main FIX logic
//
// Generate fixes and apply them to the source files
//
func Fix(appCtx context.Context, fixCtx *intrem_api.AppContext) error {

	var err error
	var fix *intrem_api.FixResult

	// -- Two choices:
	//    (1) Generate new fixes (and store them) using the ML server
	//    (2) Load previously generated fixes from local file
	if ! fixCtx.Reuse {
		// -- Option 1: Do the patch generation process
		fix, err = GenerateFix(appCtx, fixCtx)
		if err != nil {
			return err
		}

		// -- Store the fixes for future use
		fixdata, err := json.MarshalIndent(fix, "", "  ")
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(fixFilename(fixCtx), fixdata, 0644)
		if err != nil {
			return err
		}
	} else {
		// -- Option 2: just try a different fix from the file
		if fixCtx.IssueId == -1 {
			return fmt.Errorf("An issue ID is required for non-generate mode")
		}

		fixdata, err := ioutil.ReadFile(fixFilename(fixCtx))
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				err = fmt.Errorf("No previously generated fixes for Issue %d", fixCtx.IssueId)
			}
			return err
		}

		err = json.Unmarshal(fixdata, &fix)
		if err != nil {
			return fmt.Errorf("No previous fixes: %s", err)
		}
	}

	if fix != nil {

		// -- Fix to apply
		var choice int

		if fixCtx.Choose {
			// -- Display the results, allow the user to choose
			for i, diff := range fix.UnifiedDiffs {
				fmt.Printf("\n--- FIX %d -------------------------------------------------------------\n\n", i+1)
				c := strings.Index(diff, "--- Reference patch")
				diffonly := diff[:c]
				// diffclean := strings.ReplaceAll(diffonly, "/* BAD */", "")

				fmt.Printf("%s\n", diffonly)
			}

			fmt.Printf("Select a fix to apply? [1-%d]: ", len(fix.UnifiedDiffs))
			fmt.Scanf("%d", &choice)
			choice = choice-1
		} else {
			// -- Automatically choose the top fix
			choice = 0
		}

		if choice >= 0 && choice < len(fix.UnifiedDiffs) {

			var diff = fix.UnifiedDiffs[choice]
			c := strings.Index(diff, "--- Reference patch")
			diffonly := diff[:c-8]

			fmt.Printf("Applying fix %d...\n", choice+1)

			// -- Parse the diff
			if fixCtx.Debug {
				log.Printf("--- Patch to apply ---\n%s\n", diffonly)
			}

			patchfilename := fmt.Sprintf("patch-%d-%d", choice, fixCtx.IssueId)
			patchdata := []byte(diffonly)
			err = ioutil.WriteFile(patchfilename, patchdata, 0644)
			if err != nil {
				return err
			}

			log.Printf("Backing up %s to %s\n", fixCtx.SourcePath, fmt.Sprintf("%s.backup", fixCtx.SourcePath))
			backup := exec.Command("cp", fixCtx.SourcePath, fmt.Sprintf("%s.backup", fixCtx.SourcePath))
			err = backup.Run()
			if err != nil {
				return err
			}

			log.Printf("Applying patch %s\n", patchfilename)
			cmd := exec.Command("git", "apply", "--recount", "--ignore-space-change", patchfilename)
    	err = cmd.Run()
			if err != nil {
				return err
			}

			/*
			files, preamble, err := gitdiff.Parse(strings.NewReader(diffonly))
			UNUSED(preamble)
			if err != nil {
					return err
			}

			code, err := os.Open(fixCtx.SourcePath)
			if err != nil {
					return err
			}

			// apply the changes in the patch to a source file
			if fixCtx.Debug {
				log.Printf("Apply fix: \n%s\n", files)
			}

			var output bytes.Buffer
			err = gitdiff.Apply(&output, code, files[0])
			if err != nil {
					return err
			}
			*/
		}
	}

	return err
}


func fixFilename(fixCtx *intrem_api.AppContext) string {
	// -- Filename for storing or loading fixes
	filename := fmt.Sprintf("veracode-fix-%d.json", fixCtx.IssueId)
	return filename
}

// -----------------------------------------------------------------------------
//
// Generate a patch file for the given flaw
// Input:
//    Fix context object with the following fields:
//        SourcePath  : the source file to fix
//        ResultsPath : the pipeline results file with the issues
//        IsseuId     : the ID of the issue to fix
// Output:
//    A JSON representation of a proposed patch (array of unified diffs)
//
// This function extracts the issue to fix from the results file, converts it to the JSON format that the Jaroona
// back end expects, and then waits for the results, which are provided as an array of strings. Each string is
// a potential patch in unified diff format, ordered from highest to lowest confidence.
func GenerateFix(appCtx context.Context, fixCtx *intrem_api.AppContext) (*intrem_api.FixResult, error) {

  // -- Read the results.json and convert to a ScanFindings structure
  if fixCtx.Debug {
    log.Printf("Read results file %s\n", fixCtx.ResultsPath)
  }

  content, err := ioutil.ReadFile(fixCtx.ResultsPath)
  if err != nil {
      return nil, err
  }

  if fixCtx.Debug {
    log.Printf("Parse JSON data\n")
  }
  var results greenlight_api.ScanFindings
  err = json.Unmarshal(content, &results)
  if err != nil {
    return nil, err
  }

  fmt.Printf("Read %d finding(s) from %s\n", len(results.Findings), fixCtx.ResultsPath)

  // -- Obsolete: Read in the results as JSON, convert to SARIF form
  //sarif, err := ConvertResultsToSarif(results, issueid)

	// -- Get the base filename of the source file
	sourcefile := filepath.Base(fixCtx.SourcePath)

	// -- Find a specific issue to fix
	var flawtofix *intrem_api.FlawToFix
	var finding greenlight_api.Finding

	// -- No issue specified: show a list
	if fixCtx.IssueId == -1 {
		fmt.Printf("Issues in source file %s:\n", sourcefile)
		var count = 0
		var finding greenlight_api.Finding
		for _, finding = range results.Findings {
			thisfile := filepath.Base(finding.Files.SourceFile.File)
			if thisfile == sourcefile {
				count = count + 1
				fmt.Printf("   IssueID %d: CWEId %s\n                 %s\n                 on line %d in function %s\n",
					finding.IssueID, finding.CweID, finding.IssueType, finding.Files.SourceFile.Line, finding.Files.SourceFile.QualifiedFunctionName)
			}
		}
		if count == 0 {
			err = fmt.Errorf("No issues in %s match source file %s\n", fixCtx.ResultsPath, sourcefile)
			return nil, err
		}
		fmt.Printf("Enter issue ID: ")
		fmt.Scanf("%d", &fixCtx.IssueId)
	}

  // -- Given a specific issue ID, just convert that one finding
  //    into a "FlawToFix" structure
  for _, finding = range results.Findings {
    if finding.IssueID == fixCtx.IssueId {
			thisfile := filepath.Base(finding.Files.SourceFile.File)
			if thisfile == sourcefile {
      	flawtofix = ConvertIssueToFlawToFix(&finding)
      	break
			} else {
				err = fmt.Errorf("Issue %d in source file %s does not match provided source file %s\n", fixCtx.IssueId, thisfile, sourcefile)
				return nil, err
			}
    }
  }

	if flawtofix == nil {
		err = fmt.Errorf("Could not find issue with ID %d", fixCtx.IssueId)
		return nil, err
	}

  if fixCtx.Debug {
    s, _ := json.MarshalIndent(flawtofix, "", "    ")
    log.Printf("Flaw to fix:\n%s\n", s)
  }

	fmt.Printf("Request auto-fix:\n")
	fmt.Printf("    CWEId:       %s   (%s)\n", flawtofix.CWEId, finding.IssueType)
	fmt.Printf("    Function:    %s\n", flawtofix.Function)
	fmt.Printf("    Source file: %s\n", flawtofix.SourceFile)
	fmt.Printf("    Line:        %d\n", flawtofix.Line)
	fmt.Printf("\n")

	// -- Submit a request for a fix (pass in the sarif and the source file)
	fmt.Printf("Connecting to auto-remediation service...\n")
	fixCtx.ProjectId, err = requestFix(appCtx, fixCtx, flawtofix /*results*/)
	if err != nil {
  	return nil, err
	}

	fmt.Printf("Request submitted. ID = %s\nWaiting for result...\n", fixCtx.ProjectId)

	done := false
	tries := 1
	var fixes *intrem_api.FixResult

	for ! done {
		if fixCtx.Debug {
			log.Printf("CHECK for fix: attempt %d\n", tries)
		} else {
			fmt.Printf(".")
		}
  	fixes, err = checkForResults(appCtx, fixCtx)
		if err != nil {
			return nil, err
		} else {
			if fixCtx.Debug {
				log.Printf("--> GOT %d fixes\n", len(fixes.UnifiedDiffs))
			}
			if len(fixes.UnifiedDiffs) > 0 {
				fmt.Printf("\n")
				done = true
			} else {
				tries = tries + 1
				if tries > maxRetries {
					fmt.Printf("\n")
					err = fmt.Errorf("\nNo fixes generated before timeout")
					done = true
				} else {
					time.Sleep(sleepTime * time.Second)
				}
			}
		}
	}

	return fixes, err
}

// -----------------------------------------------------------------------------
//
//  Request a fix
//
func requestFix(ctx context.Context, fixCtx *intrem_api.AppContext, flawtofix *intrem_api.FlawToFix /* *intrem_api.FlawToFix*/) (string, error) {

  // -- Create a multipart writer for the body.
  body := &bytes.Buffer{}
  writer := multipart.NewWriter(body)

  // -- Write out the flaw info part as JSON
  flawPart, _ := writer.CreateFormField("flawInfo")
  flawbytes, err := json.Marshal(flawtofix)
  io.Copy(flawPart, bytes.NewReader(flawbytes))

  // -- Attach the source file
  sourceBytes, err := ioutil.ReadFile(fixCtx.SourcePath)
  if err != nil {
      log.Fatal("Error when opening file: ", err)
  }

  sourcePart, _ := writer.CreateFormFile("sourceCode", fixCtx.SourcePath)
  io.Copy(sourcePart, bytes.NewReader(sourceBytes))

  // Close multipart writer.
  writer.Close()

  // Create a request for the upload_code (and generate fix) service
  apiUrl := url.URL{
		Scheme: fixCtx.HttpScheme,
		Host:   fixCtx.APIHost,
    Path:   intremUrlPath + "/project/upload_code",
  }

  // Set the timeout for the request
  client := http.Client{
    Timeout: time.Second * requestTimeout,
  }

  req, _ := http.NewRequest(http.MethodPost, apiUrl.String(), bytes.NewReader(body.Bytes()))
  req.Header = http.Header{
    "accept": {"application/json"},
    // "Authorization":  {authHeader},
    // "PLUGIN_VERSION": {version.App + "_" + version.Version},
    "Content-Type":   {writer.FormDataContentType()},
  }

	if fixCtx.DoAuth {
    // Pass in the credentials
    creds := ctx.Value("credentials").(hmac.HmacCredentials)
    authHeader, err := hmac.CalculateAuthorizationHeader(&apiUrl, "POST", &creds)
    if err != nil {
      panic(err)
    }
    req.Header.Set("Authorization", authHeader)
    req.Header.Set("PLUGIN_VERSION", version.App + "_" + version.Version)
  }

  // Print our request
  if fixCtx.Debug {
    requestDump, err := httputil.DumpRequest(req, true)
    if err != nil {
      return "", err
    }
		log.Printf("----------------------------------------------------------------\n")
    log.Printf("REQUEST: upload_code:\n%s\n", string(requestDump))
  }

  // --------> Do it!
  resp, err := client.Do(req)
  if err != nil {
    return "", err
  }

  responseBody, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return "", err
  }

	if fixCtx.Debug {
		log.Printf("----------------------------------------------------------------\n")
	  log.Printf("RESPONSE upload_code:\n%s\n", string(responseBody))
	}

  return string(responseBody), nil
}

// -----------------------------------------------------------------------------
//
// Check server for results
//
func checkForResults(ctx context.Context, fixCtx *intrem_api.AppContext) (*intrem_api.FixResult, error) {
  // access or construct the get details path
	resultsPath := fmt.Sprintf("/project/%s/results", fixCtx.ProjectId)

  // Create a request for the upload_code (and generate fix) service
  apiUrl := url.URL{
		Scheme: fixCtx.HttpScheme,
		Host:   fixCtx.APIHost,
    Path:   intremUrlPath + resultsPath,
  }

	httpMethod := "GET"

	client := http.Client{
		Timeout: time.Second * requestTimeout,
	}

	req, err := http.NewRequest(httpMethod, apiUrl.String(), nil)
	if err != nil {
		return nil, err
	}

  if fixCtx.DoAuth {
    // Pass in the credentials
    creds := ctx.Value("credentials").(hmac.HmacCredentials)
    authHeader, err := hmac.CalculateAuthorizationHeader(&apiUrl, "GET", &creds)
    if err != nil {
      panic(err)
    }
    req.Header.Set("Authorization", authHeader)
    req.Header.Set("PLUGIN_VERSION", version.App + "_" + version.Version)
  }

	// Print our request
	if fixCtx.Debug {
		requestDump, err := httputil.DumpRequest(req, false)
		if err != nil {
			return nil, err
		}
		log.Printf("----------------------------------------------------------------\n")
		log.Printf("REQUEST results:\n%s\n", string(requestDump))
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if fixCtx.Debug {
		log.Printf("----------------------------------------------------------------\n")
		log.Printf("RESPONSE results:\n%s\n", responseBody)
	}

	var results intrem_api.FixResult

	if resp.StatusCode == 200 || resp.StatusCode == 202 {
		if len(responseBody) > 0 {
			withstruct := "{ \"diffs\" : " + string(responseBody) + "}"
			err = json.Unmarshal([]byte(withstruct), &results)
		}
	}

  return &results, err
}
