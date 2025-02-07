package vet

// Option is an option to Vet.
type Option interface {
	apply(*options)
}

// Level specifies the minimum level of validation errors to report on.
func Level(l Lvl) Option {
	return option(func(o *options) {
		o.lvl = l
	})
}

type options struct {
	lvl Lvl
}

type option func(*options)

func (o option) apply(opts *options) {
	o(opts)
}
