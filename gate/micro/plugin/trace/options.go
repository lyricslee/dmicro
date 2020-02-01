package trace

import (
	"dmicro/gate/micro/plugin"
	"github.com/opentracing/opentracing-go"
)

type Options struct {
	tracer      opentracing.Tracer
	skipperFunc plugin.SkipperFunc
}

type Option func(*Options)

func newOptions(opts ...Option) Options {
	opt := Options{
		tracer:      opentracing.GlobalTracer(),
		skipperFunc: plugin.DefaultSkipperFunc,
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

func WithTracer(tracer opentracing.Tracer) Option {
	return func(o *Options) {
		o.tracer = tracer
	}
}

func WithSkipperFunc(skipperFunc plugin.SkipperFunc) Option {
	return func(o *Options) {
		o.skipperFunc = skipperFunc
	}
}
