package auth

import "dmicro/gate/micro/plugin"

// Options 就是该插件的参数，目前只有 SkipperFunc 就是处理函数了。
// skipperFunc
type Options struct {
	skipperFunc plugin.SkipperFunc
}

type Option func(*Options)

// new opts
func newOptions(opts ...Option) Options {
	opt := Options{skipperFunc: plugin.DefaultSkipperFunc}
	for _, o := range opts {
		o(&opt) // set default parameters for OptionFunc
	}
	return opt
}

// register skipper func
func WithSkipperFunc(skipperFunc plugin.SkipperFunc) Option {
	return func(o *Options) {
		o.skipperFunc = skipperFunc
	}
}
