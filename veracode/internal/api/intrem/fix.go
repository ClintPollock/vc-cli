package intrem

// This type represents a proposed fix returned from the auto-remediation
// system.

type FixResult struct {
  UnifiedDiffs []string    `json:"diffs,omitempty"`
}
