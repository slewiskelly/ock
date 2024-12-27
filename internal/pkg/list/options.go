package list

// Option is an option to List.
type Option interface {
	apply(*options)
}

// WithExpr specifies an expression used to filter files.
func WithExpr(e string) Option {
	return option(func(o *options) {
		o.expr = e
	})
}

type options struct {
	expr string
}

type option func(*options)

func (o option) apply(opts *options) {
	o(opts)
}
