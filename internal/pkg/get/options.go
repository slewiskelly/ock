package get

// Option is an option to Get.
type Option interface {
	apply(*options)
}

// Expr specifies an expression used to filter fields.
func Expr(e string) Option {
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
