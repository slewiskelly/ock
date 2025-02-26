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
	Errors   []Error   `json:"errors,omitempty"`   // Any validation errors encountered.
	Warnings []Error   `json:"warnings,omitempty"` // Any validation warnings encountered.
}

type Error struct {
	Field   string `json:"field,omitempty"`   // Field containing a validaion error.
	Message string `json:"message,omitempty"` // Validation error.
}
