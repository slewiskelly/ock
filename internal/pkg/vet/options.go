package vet

// Option is an option to Vet.
type Option interface {
	apply(*options)
}

// Glob specifies a pattern to filter files that are attempted to be validated.
func Glob(pattern string) Option {
	return option(func(o *options) {
		o.glob = pattern
	})
}

// Level specifies the minimum level of validation errors to report on.
func Level(l Lvl) Option {
	return option(func(o *options) {
		o.lvl = l
	})
}

type options struct {
	glob string
	lvl  Lvl
}

type option func(*options)

func (o option) apply(opts *options) {
	o(opts)
}
