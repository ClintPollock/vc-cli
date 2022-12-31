package sast_next

type JobEvent struct {
	ID          string  `json:"id,omitempty"`
	Time        float64 `json:"time,omitempty"`
	Source      string  `json:"source,omitempty"`
	Specversion string  `json:"specversion,omitempty"`
	Type        string  `json:"type,omitempty"`
	Message     string  `json:"message,omitempty"`
	ScanID      string  `json:"scanId,omitempty"`
	JobID       string  `json:"jobId,omitempty"`
	Context     string  `json:"context,omitempty"`
	Level       string  `json:"level,omitempty"`
}
