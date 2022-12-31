package configuration

type Configuration struct {
	ApiScheme    string //https
	SSLIgnore    bool   // false
	BaseHostname string // api.veracode.com

	PipelineUrlPath string // "/pipeline_scan/v1"
	EventsUrlPath   string
	PolicyUrlPath   string // "/appsec/v1/policies"

	Mock struct {
		NumberOfScannerInstances string //int
		NumberOfFindingsToThrow  string //int
		DelayInMsBetweenFindings string //int
		RunScannersInParallel    string //bool
		SkipUpload               bool
	}

	SAST struct {
		UseRealTimeFlawAPIs bool
	}
}

var DefaultConfiguration = Configuration{

	ApiScheme:    "https",
	SSLIgnore:    false,
	BaseHostname: "api.veracode.com",

	PipelineUrlPath: "/pipeline_scan/v1",
	EventsUrlPath:   "/v1/events",
	PolicyUrlPath:   "/appsec/v1/policies",

	//Mock: nil,
}
