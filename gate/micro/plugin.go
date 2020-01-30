package main

import (
	"dmicro/gate/micro/plugin/auth"
	"dmicro/gate/micro/plugin/metrics"
	"dmicro/gate/micro/plugin/trace"
	"github.com/micro/micro/web"
)

func init() {
	web.Register(metrics.NewPlugin())
	web.Register(trace.NewPlugin())
	//web.Register(breaker.NewPlugin())
	web.Register(auth.NewPlugin())
}
