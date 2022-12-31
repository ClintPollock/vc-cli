package intrem

import (
    //"encoding/json"
    //"io/ioutil"
    //"log"
    "strconv"
    "errors"
    "fmt"
    greenlight_api "github.com/veracode/veracode-cli/internal/api/greenlight"
    intrem_api "github.com/veracode/veracode-cli/internal/api/intrem"
)

// -----------------------------------------------------------------------------
//  Convert issues to simplified flaw-to-fix format
// -----------------------------------------------------------------------------

// Convert Veracode Issue to a form suitable for Jaroona
//
// Input is a Finding structure read from results.json, which contains all the information
// about one findings. Output is a FlawToFix structure, which extracts just the information
// that the Jaroona back-end needs to generate a fix.
func ConvertIssueToFlawToFix (finding *greenlight_api.Finding) (*intrem_api.FlawToFix) {

  // -- Get the main information (file, function, and line)
  sourcefile := finding.Files.SourceFile

  // -- Convert the stack dump
  //    Only include the stack frames that are in the original source file
  //    Only include the useful parts: a varname expression and line number
  var steps []intrem_api.Step
  var frames []greenlight_api.Frame
  frames = finding.StackDumps.StackDumpArray[0].FrameArray
  for _, frame := range frames {
    if frame.SourceFile == sourcefile.File {
      line, err := strconv.Atoi(frame.SourceLineStr)
      if err != nil {
        // -- This means that the engine produced bad line numbers
        panic(err)
      }
      step := intrem_api.Step {
        Region: intrem_api.Region{
          StartLine: line,
          EndLine: line,
        },
        Expression: frame.VarNames,
      }
      steps = append(steps, step)
    }
  }

  // -- Create the new Flaw object
  newflaw := &intrem_api.FlawToFix{
    SourceFile: sourcefile.File,
    Function: sourcefile.FunctionName,
    Line: sourcefile.Line,
    CWEId: finding.CweID,
    Flow: steps,
  }

  return newflaw
}

// -----------------------------------------------------------------------------
//  Convert issues to SARIF format
//
//  This code mirrors the strategy used by the results-to-SARIF converter
//  on our public GitHub repo
// -----------------------------------------------------------------------------

func ConvertResultsToSarif(results greenlight_api.ScanFindings, flawid int) (*intrem_api.SARIFSimple, error){

    // -- Convert our results format to SARIF
    //    Produces a list of rules (corresponding to CWEIDs) and
    //    a list of results (corresponding to findings)
    var rules []*intrem_api.Rule
    var sarifresults []*intrem_api.Result
    if flawid == 0 {
      // -- When FlawID is 0, convert all of the findings to SARIF
      // -- Make a list of rules (unique CWE IDs)
      rulemap := make(map[string]*intrem_api.Rule)
      for _, finding := range results.Findings {
        _, prs := rulemap[finding.CweID]
        if ! prs {
          rulemap[finding.CweID] = IssueToRule(&finding)
        }
      }

      rules = make([]*intrem_api.Rule, len(rulemap))
      i := 0
      for _,v := range rulemap {
        rules[i] = v
        i = i + 1
      }

      // -- Convert each issue into a SARIF result
      if flawid == 0 {
        sarifresults = make([]*intrem_api.Result, len(results.Findings))
        for index, finding := range results.Findings {
          sarifresults[index] = IssueToResult(&finding)
        }
      }
    } else {
      // -- When given a specific flaw ID, just convert that one finding
      for _, finding := range results.Findings {
        if finding.IssueID == flawid {
          rules = []*intrem_api.Rule{ IssueToRule(&finding) }
          sarifresults = []*intrem_api.Result{ IssueToResult(&finding) }
        }
      }
    }

    if rules == nil || sarifresults == nil {
      err := errors.New(fmt.Sprintf("Could not file flaw with ID %d", flawid))
      return nil, err
    }

    var run intrem_api.Run
    run.Tool.Driver.Name = "Veracode Static Analysis Pipeline Scan"
    run.Tool.Driver.Rules = rules
    run.Results = sarifresults

    // -- Build SARIF structure
    sarif := intrem_api.SARIFSimple{
      Schema:  "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
      Version: "2.1.0",
      Runs: []intrem_api.Run { run },
    }

    return &sarif, nil
}

func IssueToRule(finding *greenlight_api.Finding) *intrem_api.Rule {
  newrule := &intrem_api.Rule{
    ID: finding.CweID,
    Name: finding.IssueType,
    HelpURI: "https://cwe.mitre.org/data/definitions/" + finding.CweID + ".html",
  }

  newrule.ShortDescription.Text = "CWE-" + finding.CweID + ": " + finding.IssueType
  newrule.Properties.Category = finding.IssueTypeID
  newrule.Properties.Tags = []string{finding.IssueTypeID}
  newrule.DefaultConfiguration.Level = "??"
  return newrule
}

func IssueToResult(finding *greenlight_api.Finding) *intrem_api.Result {
  newresult := &intrem_api.Result {
    Level: "??",
    Rank: finding.Severity,
  }

  sourcefile := finding.Files.SourceFile

  var loc intrem_api.Location
  loc.PhysicalLocation.ArtifactLocation.URI = ""
  loc.PhysicalLocation.Region.StartLine = sourcefile.Line
  loc.PhysicalLocation.Region.EndLine = sourcefile.Line
  loc.LogicalLocations = []intrem_api.LogicalLocation{
    intrem_api.LogicalLocation {
      Name: sourcefile.FunctionName,
      FullyQualifiedName: sourcefile.QualifiedFunctionName,
      Kind: "function",
    },
    intrem_api.LogicalLocation {
      FullyQualifiedName: finding.Title,
      Kind: "member",
    },
  }
  newresult.Locations = []intrem_api.Location{ loc }

  newresult.Message.Text = finding.DisplayText
  cweInt, _ := strconv.Atoi(finding.CweID)
  newresult.RuleID = cweInt

  return newresult
}
