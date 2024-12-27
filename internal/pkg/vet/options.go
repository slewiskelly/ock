package vet

// Option is an option to Vet.
type Option interface {
	apply(*options)
}

type options struct{}

type option func(*options)

func (o option) apply(opts *options) {
	o(opts)
}
