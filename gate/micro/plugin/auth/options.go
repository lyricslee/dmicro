package auth

import "dmicro/gate/micro/plugin"

type Options struct {
	skipperFunc plugin.SkipperFunc
}

type Option func(*Options)

func newOptions(opts ...Option) Options {
	opt := Options{skipperFunc: plugin.DefaultSkipperFunc}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

func WithSkipperFunc(skipperFunc plugin.SkipperFunc) Option {
	return func(o *Options) {
		o.skipperFunc = skipperFunc
	}
}
