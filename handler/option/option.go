package option

import "github.com/tianhongw/tinyid/pkg/log"

type Options struct {
	Logger log.Logger
}

type Option func(*Options)

func WithLogger(logger log.Logger) Option {
	return func(opts *Options) {
		opts.Logger = logger
	}
}
