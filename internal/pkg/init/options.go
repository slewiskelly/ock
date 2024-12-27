package init

// Option is an option to Init.
type Option interface {
	apply(*options)
}

// Force specifies that an existing schema file should be overwritten.
func Force(b bool) Option {
	return option(func(o *options) {
		o.force = b
	})
}

type options struct {
	force bool
}

type option func(*options)

func (o option) apply(opts *options) {
	o(opts)
}
