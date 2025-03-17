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
	Metadata cue.Value `json:"metadata,omitzero"`  // File's metadata.
	Start    int       `json:"start"`              // Line number after the opening delimiter.
	End      int       `json:"end"`                // Line number before the closing delimiter.
	Errors   []Error   `json:"errors,omitempty"`   // Any validation errors encountered.
	Warnings []Error   `json:"warnings,omitempty"` // Any validation warnings encountered.
}

type Error struct {
	Field   string `json:"field,omitempty"`   // Field containing a validaion error.
	Message string `json:"message,omitempty"` // Validation error.
}
