package sast_next

type ScanLifecycleEvent struct {
	ID          string  `json:"id,omitempty"`
	Time        float64 `json:"time,omitempty"`
	Source      string  `json:"source,omitempty"`
	Specversion string  `json:"specversion,omitempty"`
	Type        string  `json:"type,omitempty"`
	Message     string  `json:"message,omitempty"`
	ScanID      string  `json:"scanId,omitempty"`
	Status      string  `json:"status,omitempty"`
	Context     string  `json:"context,omitempty"`
}
