// Package vet provides functionalty to validate a file's metadata.
package vet

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/errors"
	"github.com/bmatcuk/doublestar/v4"

	"github.com/slewiskelly/ock/internal/pkg/get"
	"github.com/slewiskelly/ock/internal/pkg/report"
)

// Lvl represents a validation error level.
type Lvl int

const (
	LvlWarn  Lvl = 4
	LvlError Lvl = 8
)

// Vet validates all files rooted at the given path, against the given schema.
//
// The returned report contains all files which failed validation, along with
// their corresponding error(s).
func Vet(path string, schema cue.Value, opts ...Option) (report.Report, error) {
	o := &options{
		lvl: LvlWarn,
	}

	for _, opt := range opts {
		opt.apply(o)
	}

	if ok := doublestar.ValidatePathPattern(o.glob); !ok {
		return nil, errors.New("invalid globbing pattern")
	}

	if err := schema.Err(); err != nil {
		return nil, fmt.Errorf("invalid schema: %w", err)
	}

	schema = schema.LookupPath(cue.ParsePath("#Metadata"))

	var r report.Report

	err := filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if o.glob != "" && !doublestar.PathMatchUnvalidated(o.glob, p) {
			return nil
		}

		if !d.Type().IsRegular() || filepath.Ext(p) != ".md" {
			return nil
		}

		v, err := get.Get(p)
		if err != nil {
			r = append(r, &report.File{Name: p, Errors: []report.Error{{Message: err.Error()}}})
			return nil
		}

		// TODO(slewiskelly): Signals a lack of metadata, should be considered a
		// failure?
		if len(v) < 1 {
			return nil
		}

		if errs, wrns := validate(v[0].Metadata.Unify(schema), o.lvl); len(errs) > 0 || len(wrns) > 0 {
			r = append(r, &report.File{Name: p, Errors: errs, Warnings: wrns})
		}

		return nil
	})

	return r, err
}

func validate(v cue.Value, lvl Lvl) (errs, wrns []report.Error) {
	i, err := v.Fields()
	if err != nil {
		return []report.Error{{Message: fmt.Sprintf("Failed to validate: %v", err)}}, nil // TODO(slewiskelly): Reconsider.
	}

	for i.Next() {
		x := i.Value()

		// Recursively check fields if there is a nested structure.
		if _, err := x.Fields(); err == nil {
			e, w := validate(x, lvl)

			errs = append(errs, e...)
			wrns = append(wrns, w...)

			continue
		}

		if err := x.Validate(cue.Concrete(true)); err != nil {
			if a := x.Attribute("error"); a.NumArgs() > 0 {
				if lvl <= LvlError {
					errs = append(errs, report.Error{
						Field:   x.Path().String(),
						Message: a.Contents(),
					})
				}
				continue
			}

			if a := x.Attribute("warning"); a.NumArgs() > 0 {
				if lvl <= LvlWarn {
					wrns = append(wrns, report.Error{
						Field:   x.Path().String(),
						Message: a.Contents(),
					})
				}
				continue
			}

			errs = append(errs, report.Error{
				Field:   x.Path().String(),
				Message: errDetails(err).Error(),
			})
		}
	}

	return errs, wrns
}

func errDetails(e error) error {
	qe, ok := e.(errors.Error)
	if !ok {
		return e
	}

	var msgs []string

	qe = errors.Sanitize(qe)

	for _, err := range errors.Errors(qe) {
		f, a := err.Msg()
		msgs = append(msgs, fmt.Sprintf(f, a...))
	}

	return errors.New(strings.Join(msgs, "\n"))
}
