package websocket

import "time"

type serverOption struct {
	Authentication
	pattern     string
	MaxIdleTime time.Duration
}

type ServerOption func(option *serverOption)

func newServerOption(opts ...ServerOption) serverOption {
	o := serverOption{
		Authentication: new(authentication),
		pattern:        "/ws",
		MaxIdleTime:    defaultMaxIdleTime,
	}

	for _, opt := range opts {
		opt(&o)
	}

	return o
}

func WithServerAuthentication(authentication Authentication) ServerOption {
	return func(option *serverOption) {
		option.Authentication = authentication
	}
}

func WithServerPattern(pattern string) ServerOption {
	return func(option *serverOption) {
		option.pattern = pattern
	}
}

func WithServerMaxIdleTime(maxIdleTime time.Duration) ServerOption {
	return func(option *serverOption) {

		if maxIdleTime < 0 {
			return
		}
		option.MaxIdleTime = maxIdleTime
	}
}
