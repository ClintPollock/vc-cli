package greenlight

// @Expose
// private String app_id;
// @Expose
// private String binary_hash;
// @Expose
// private String binary_name;
// @Expose
// private int binary_size;
// @Expose
// private String dev_stage;
// @Expose
// private String plugin_version;
// @Expose
// private String project_name;
// @Expose
// private String project_ref;
// @Expose
// private String project_uri;
// @Expose
// private int scan_timeout;

// private String scan_id;
// private String scan_status;
// private String scan_message;

//
// Sample response
//

// {
//   "_links": {
//     "root": {"href": "/"},
//     "self": {"href": "/scans"},
//     "help": {"href": "https://help.veracode.com/reader/tS9CaFwL4_lbIEWWomsJoA/ovfZGgu96UINQxIuTqRDwg"},
//     "create": {"href": "/scans"},
//     "details": {"href": "/scans/8e7e90a2-37f0-4172-b4aa-16d91dde7a6f"},
//     "upload": {"href": "/scans/8e7e90a2-37f0-4172-b4aa-16d91dde7a6f/segments/0"},
//     "cancel": {"href": "/scans/8e7e90a2-37f0-4172-b4aa-16d91dde7a6f"}
//   },
//   "scan_id": "8e7e90a2-37f0-4172-b4aa-16d91dde7a6f",
//   "scan_status": "UPLOADING",
//   "api_version": 1,
//   "app_id": "12345678",
//   "project_name": "PROTOTYPE",
//   "project_uri": "https://bogus-uri.veracode.com/project",
//   "project_ref": "BOGUS PROJECT REF",
//   "commit_hash": null,
//   "dev_stage": "DEVELOPMENT",
//   "binary_name": "veracode/veracode",
//   "binary_size": 13237425,
//   "binary_hash": "9108415ec07770cf2c16eadc30083c3c190f72138e48f5c6ac56ebb78446c405",
//   "binary_segments_expected": 6,
//   "binary_segments_uploaded": 0,
//   "scan_timeout": 10,
//   "scan_duration": null,
//   "results_size": null,
//   "message": null,
//   "created": "2022-09-20T16:24:42.618466",
//   "changed": "2022-09-20T16:24:44.840257"
// }

type ScanInfo struct {
	Id       string `json:"scan_id,omitempty"`
	Status   string `json:"scan_status,omitempty"`
	Message  string `json:"message,omitempty"`
	Timeout  int    `json:"scan_timeout"`  // up to 60 (miniutes?)
	Duration int    `json:"scan_duration"` // up to 60 (miniutes?)

	Created string `json:"created,omitempty"`
	Changed string `json:"changed,omitempty"`

	ApiVersion    string `json:"api_version"`
	AppID         string `json:"app_id"`
	Stage         string `json:"dev_stage"` //['DEVELOPMENT', 'TESTING', 'RELEASE']"
	PluginVersion string `json:"plugin_version"`
	EmitStackDump string `json:"emit_stack_dump"`

	BinaryHash             string `json:"binary_hash"`
	BinaryName             string `json:"binary_name"`
	BinarySize             int64  `json:"binary_size"`
	BinarySegmentsExpected int    `json:"binary_segments_expected"`
	BinarySegmentsUploaded int    `json:"binary_segments_uploaded"`

	ProjectName string `json:"project_name"`
	ProjectRef  string `json:"project_ref"`
	ProjectURI  string `json:"project_uri"`

	CommitHash string `json:"commit_hash,omitempty"`

	ResultsSize int `json:"results_size,omitempty"`

	Links struct {
		Root struct {
			Href string `json:"href,omitempty"`
		} `json:"root,omitempty"`
		Self struct {
			Href string `json:"href,omitempty"`
		} `json:"self,omitempty"`
		Help struct {
			Href string `json:"href,omitempty"`
		} `json:"help,omitempty"`
		Create struct {
			Href string `json:"href,omitempty"`
		} `json:"create,omitempty"`
		Upload struct {
			Href string `json:"href,omitempty"`
		} `json:"upload,omitempty"`
		Start struct {
			Href string `json:"href,omitempty"`
		} `json:"start,omitempty"`
		Details struct {
			Href string `json:"href,omitempty"`
		} `json:"details,omitempty"`
		Findings struct {
			Href string `json:"href,omitempty"`
		} `json:"findings,omitempty"`
		Cancel struct {
			Href string `json:"href,omitempty"`
		} `json:"cancel,omitempty"`
	} `json:"_links,omitempty"`
}
