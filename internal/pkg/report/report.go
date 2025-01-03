// Package report implements types to report the status of individual files.
package report

import (
	"cuelang.org/go/cue"
)

// Report represents a report of individual files.
type Report []*File

// File represents an individual file.
type File struct {
	Name     string    `json:"name,omitempty"`     // Name of the file.
	Metadata cue.Value `json:"metadata,omitempty"` // File's metadata.
	Errors   []string  `json:"errors,omitempty"`   // Any validation errors encountered.
}
