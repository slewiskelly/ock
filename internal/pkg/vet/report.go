package vet

// Report is contains all files which failed validation along with their
// corresponding errors.
type Report struct {
	// Files maps filenames to all validation errors encountered.
	Files map[string][]string `json:"files,omitempty"`
}
