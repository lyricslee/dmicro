package breaker

import "dmicro/gate/micro/plugin"

type Options struct {
	skipperFunc plugin.SkipperFunc
}

// Option for Plugin
type Option func(*Options)

func newOptions(opts ...Option) Options {
	opt := Options{skipperFunc: plugin.DefaultSkipperFunc}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

// skipperFunc register
func WithSkipperFunc(skipperFunc plugin.SkipperFunc) Option {
	return func(o *Options) {
		o.skipperFunc = skipperFunc
	}
}