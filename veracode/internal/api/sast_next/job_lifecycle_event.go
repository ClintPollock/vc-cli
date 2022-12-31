package sast_next

type JobLifecycleEvent struct {
	ID                 string  `json:"id,omitempty"`
	Time               float64 `json:"time,omitempty"`
	Source             string  `json:"source,omitempty"`
	Specversion        string  `json:"specversion,omitempty"`
	Type               string  `json:"type,omitempty"`
	Message            string  `json:"message,omitempty"`
	ScanID             string  `json:"scanId,omitempty"`
	JobID              string  `json:"jobId,omitempty"`
	Context            string  `json:"context,omitempty"`
	Status             string  `json:"status,omitempty"`
	Phase              int     `json:"phase,omitempty"`
	CurrentMemoryUsage int     `json:"currentMemoryUsage,omitempty"`
	MaxMemoryUsage     int     `json:"maxMemoryUsage,omitempty"`
	ProcessingTime     int     `json:"processingTime,omitempty"`
}
