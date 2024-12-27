// Package vet provides functionalty to validate a file's metadata.
package vet

import (
	"io/fs"
	"path/filepath"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/errors"

	"github.com/slewiskelly/ock/internal/pkg/get"
)

// Vet validates all files rooted at the given path, against the given schema.
//
// The returned report contains all files which failed validation, along with
// their corresponding error(s).
func Vet(path string, schema *cue.Value, opts ...Option) (*Report, error) {
	o := &options{}

	for _, opt := range opts {
		opt.apply(o)
	}

	r := &Report{Files: make(map[string][]string)}

	err := filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.Type().IsRegular() || filepath.Ext(p) != ".md" {
			return nil
		}

		v, err := get.Get(p)
		if err != nil {
			r.Files[p] = append(r.Files[p], err.Error())
			return nil
		}

		// TODO(slewiskelly): Signals a lack of metadata, should be considered a
		// failure?
		if v == nil {
			return nil
		}

		if err := (*v).Unify(*schema).Validate(); err != nil {
			r.Files[p] = append(r.Files[p], errs(err)...)
		}

		return nil
	})

	return r, err
}

func errs(err error) []string {
	var s []string

	for _, e := range errors.Errors(err) {
		s = append(s, e.Error())
	}

	return s
}
