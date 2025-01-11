package websocket

import "net/http"

type DialOptions func(option *DialOption)

type DialOption struct {
	Pattern string
	header  http.Header
}

func NewDialOption(opts ...DialOptions) DialOption {
	o := DialOption{
		Pattern: "/ws",
		header:  nil,
	}

	for _, opt := range opts {
		opt(&o)
	}

	return o
}
func WithPattern(pattern string) DialOptions {
	return func(option *DialOption) {
		option.Pattern = pattern
	}
}

func WithHeader(header http.Header) DialOptions {
	return func(option *DialOption) {
		option.header = header
	}
}
