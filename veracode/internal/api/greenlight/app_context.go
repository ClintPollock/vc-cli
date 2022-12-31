package greenlight

type AppContext struct {
	ResumeScanId string
	ProjectName  string
	ProjectUrl   string
	ProjectRef   string
	AppId        string
	Stage        string

	FailOnSeverity float32
	FailOnCWEs     []string
	FailFast       bool
	BaselineFile   string

	BaselineFindings *ScanFindings
	ShowDetails      bool
	Debug            bool

	EmitStackDump    bool
}
