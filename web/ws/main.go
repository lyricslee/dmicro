package main

import (
	"github.com/micro/go-micro/web"
	"github.com/opentracing/opentracing-go"

	"dmicro/common/log"
	"dmicro/common/service"
	"dmicro/pkg/tracer"
	"dmicro/web/ws/internal/broker"
	"dmicro/web/ws/internal/client"
	"dmicro/web/ws/internal/config"
	"dmicro/web/ws/internal/router"
)

func main() {
	// Config
	config.Init()
	// Logger
	log.Init(config.Logger)
	// New Service
	svc := service.NewWebService()

	var opts []web.Option
	opts = append(opts, web.Handler(router.Router()))
	if err := svc.Init(opts...); err != nil {
		log.Fatal(err)
	}
	// Tracer
	t, err := tracer.Init(config.Micro.ServerName, config.Tracer.Addr)
	if err != nil {
		log.Error(err)
	}
	opentracing.SetGlobalTracer(t)

	broker.Init()

	client.Init(svc)
	// Run svc
	if err := svc.Run(); err != nil {
		log.Fatal(err)
	}
}
