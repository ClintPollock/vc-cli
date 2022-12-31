package sast_next

type ScanFindingEvent struct {
	ID          string  `json:"id,omitempty"`
	Time        float64 `json:"time,omitempty"`
	Source      string  `json:"source,omitempty"`
	Specversion string  `json:"specversion,omitempty"`
	Type        string  `json:"type,omitempty"`
	Message     string  `json:"message,omitempty"`
	ScanID      string  `json:"scanId,omitempty"`
	JobID       string  `json:"jobId,omitempty"`
	Context     string  `json:"context,omitempty"`
	Issue       struct {
		Title             string      `json:"title,omitempty"`
		IssueID           int         `json:"issueId,omitempty"`
		GreenlightFinding interface{} `json:"greenlightFinding,omitempty"`
		IssueTypeID       string      `json:"issueTypeId,omitempty"`
		IssueType         string      `json:"issueType,omitempty"`
		Severity          int         `json:"severity,omitempty"`
		MessageKey        interface{} `json:"messageKey,omitempty"`
		AnnotationID      interface{} `json:"annotationId,omitempty"`
		CweID             int         `json:"cweId,omitempty"`
		VcID              float64     `json:"vcId,omitempty"`
		ExploitLevel      interface{} `json:"exploitLevel,omitempty"`
		ConfidenceLevel   interface{} `json:"confidenceLevel,omitempty"`
		DisplayText       string      `json:"displayText,omitempty"`
		ModuleDisplayName interface{} `json:"moduleDisplayName,omitempty"`
		Files             struct {
			SourceFile struct {
				File                  string `json:"file,omitempty"`
				FileID                int    `json:"fileId,omitempty"`
				Line                  int    `json:"line,omitempty"`
				EndLine               int    `json:"endLine,omitempty"`
				FunctionName          string `json:"functionName,omitempty"`
				QualifiedFunctionName string `json:"qualifiedFunctionName,omitempty"`
				FunctionPrototype     string `json:"functionPrototype,omitempty"`
				EntryPoints           string `json:"entryPoints,omitempty"`
				Scope                 string `json:"scope,omitempty"`
				RelativeLocation      string `json:"relativeLocation,omitempty"`
				FunctionLine          int    `json:"functionLine,omitempty"`
			} `json:"sourceFile,omitempty"`
			RtFile struct {
				File         string `json:"file,omitempty"`
				FileID       int    `json:"fileId,omitempty"`
				Line         int    `json:"line,omitempty"`
				FunctionName string `json:"functionName,omitempty"`
				FunctionLine int    `json:"functionLine,omitempty"`
			} `json:"rtFile,omitempty"`
		} `json:"files,omitempty"`
		SummaryizeDepencency      interface{} `json:"summaryizeDepencency,omitempty"`
		IssueTemplate             interface{} `json:"issueTemplate,omitempty"`
		StackDumps                interface{} `json:"stackDumps,omitempty"`
		ExploitabilityAdjustments interface{} `json:"exploitabilityAdjustments,omitempty"`
		FlawMatch                 struct {
			ProcedureHash     int         `json:"procedureHash,omitempty"`
			PrototypeHash     int         `json:"prototypeHash,omitempty"`
			FlawHash          int         `json:"flawHash,omitempty"`
			FlawHashCount     int         `json:"flawHashCount,omitempty"`
			FlawHashOrdinal   int         `json:"flawHashOrdinal,omitempty"`
			CauseHash         int         `json:"causeHash,omitempty"`
			CauseHashCount    int         `json:"causeHashCount,omitempty"`
			CauseHashOrdinal  int         `json:"causeHashOrdinal,omitempty"`
			CauseHash2        int         `json:"causeHash2,omitempty"`
			CauseHash2Ordinal int         `json:"causeHash2Ordinal,omitempty"`
			CombinedExactHash interface{} `json:"combinedExactHash,omitempty"`
			CombinedFuzzyHash interface{} `json:"combinedFuzzyHash,omitempty"`
		} `json:"flawMatch,omitempty"`
		Mitigation interface{} `json:"mitigation,omitempty"`
		Internal   interface{} `json:"internal,omitempty"`
		ScopeID    interface{} `json:"scopeId,omitempty"`
		SomID      int         `json:"somId,omitempty"`
	} `json:"issue,omitempty"`
}
