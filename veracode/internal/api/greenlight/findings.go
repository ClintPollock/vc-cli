package greenlight

type ScanFindings struct {
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
	} `json:"_links,omitempty"`
	ScanID       string    `json:"scan_id"`
	ScanStatus   string    `json:"scan_status,omitempty"`
	Message      string    `json:"message,omitempty"`
	Modules      []string  `json:"modules"`
	ModulesCount int       `json:"modules_count"`
	Findings     []Finding `json:"findings,omitempty"`
}

type Finding struct {
	Title       string `json:"title,omitempty"`
	IssueID     int    `json:"issue_id,omitempty"`
	Gob         string `json:"gob,omitempty"`
	Severity    int    `json:"severity,omitempty"`
	IssueTypeID string `json:"issue_type_id,omitempty"`
	IssueType   string `json:"issue_type,omitempty"`
	CweID       string `json:"cwe_id,omitempty"`
	DisplayText string `json:"display_text,omitempty"`
	Files       struct {
		SourceFile struct {
			File                  string `json:"file,omitempty"`
			Line                  int    `json:"line,omitempty"`
			FunctionName          string `json:"function_name,omitempty"`
			QualifiedFunctionName string `json:"qualified_function_name,omitempty"`
			FunctionPrototype     string `json:"function_prototype,omitempty"`
			Scope                 string `json:"scope,omitempty"`
		} `json:"source_file,omitempty"`
	} `json:"files,omitempty"`
	FlawMatch struct {
		ProcedureHash     string `json:"procedure_hash,omitempty"`
		PrototypeHash     string `json:"prototype_hash,omitempty"`
		FlawHash          string `json:"flaw_hash,omitempty"`
		FlawHashCount     int    `json:"flaw_hash_count,omitempty"`
		FlawHashOrdinal   int    `json:"flaw_hash_ordinal,omitempty"`
		CauseHash         string `json:"cause_hash,omitempty"`
		CauseHashCount    int    `json:"cause_hash_count,omitempty"`
		CauseHashOrdinal  int    `json:"cause_hash_ordinal,omitempty"`
		CauseHash2        string `json:"cause_hash2,omitempty"`
		CauseHash2Ordinal string `json:"cause_hash2_ordinal,omitempty"`
	} `json:"flaw_match,omitempty"`
	StackDumps struct {
		StackDumpArray []StackDump `json:"stack_dump,omitempty"`
	} `json:"stack_dumps,omitempty"`
	FlawDetailsLink string `json:"flaw_details_link,omitempty"`
}

/*
<StackDumps>
	<StackDump>
		<Frame>
			<FrameId>0</FrameId>
			...
		</Frame>
		...
	</StackDump>
	...
</StackDumps>
*/

type StackDump struct {
	FrameArray []Frame `json:"Frame,omitempty"`
}

type Frame struct {
	FrameIDStr            string   `json:"FrameId,omitempty"`
	FunctionName          string   `json:"FunctionName,omitempty"`
	SourceFile            string   `json:"SourceFile,omitempty"`
	SourceLineStr         string   `json:"SourceLine,omitempty"`
	SourceFileIdStr       string   `json:"source_file_id,omitempty"`
	RTFile                string   `json:"rt_file,omitempty"`
	RTLine                string   `json:"rt_line,omitempty"`
	RTFileId              string   `json:"rt_file_id,omitempty"`
	StatementText         struct{} `json:"statement_text,omitempty"` // This needs some definition?
	VarNames              string   `json:"var_names,omitempty,omitempty"`
	QualifiedFunctionName string   `json:"qualified_function_name,omitempty"`
	FunctionPrototype     string   `json:"function_prototype,omitempty"`
	Scope                 string   `json:"scope,omitempty"`
	RelativeLocationStr   string   `json:"relative_location,omitempty"`
}
