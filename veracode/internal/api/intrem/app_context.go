package intrem

type AppContext struct {
	ProjectId    string       // Returned from upload_code (maybe not needed here)
  HttpScheme   string       // Either http or https
  APIHost      string       // Allows us to switch between prod and test
  DoAuth       bool         // Jaroona deployment doesn't need auth
	SourcePath   string       // Path of the source file to fix
	IssueId      int          // ID of the issue to fix
	ResultsPath  string       // Path of the results file
	Reuse        bool         // If true, just choose from previously generated fixes
	Choose       bool         // If true, present the list of fixes; otherwise just take the first
	Debug        bool
}
