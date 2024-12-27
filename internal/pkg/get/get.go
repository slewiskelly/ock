// Package get provides functionality for retrieving document metadata.
package get

import (
	"bufio"
	"errors"
	"os"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/encoding/yaml"
)

// Get retrieves metadata from the given file.
func Get(file string, opts ...Option) (*cue.Value, error) {
	o := &options{}

	for _, opt := range opts {
		opt.apply(o)
	}

	var b []byte

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	var closed bool

	for i := 0; s.Scan(); i++ {
		l := s.Bytes()

		if i == 0 {
			if string(l) != "---" {
				return nil, nil
			}
			continue
		}

		if i > 0 && string(l) == "---" {
			closed = true
			break
		}

		b = append(b, append(l, '\n')...)
	}

	if !closed {
		return nil, errors.New("frontmatter not closed")
	}

	y, err := yaml.Extract(file, b)
	if err != nil {
		return nil, err
	}

	v := cuecontext.New().BuildFile(y)
	if err := v.Err(); err != nil {
		return nil, err
	}

	if o.expr == "" {
		return &v, nil
	}

	// TODO(slewiskelly): Support filtering via expression.

	return &v, nil
}
