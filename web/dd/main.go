package main

import (
	"github.com/micro/go-micro/v2/web"
	"github.com/opentracing/opentracing-go"

	"dmicro/common/log"
	"dmicro/common/service"
	"dmicro/pkg/tracer"
	"dmicro/web/dd/internal/client"
	"dmicro/web/dd/internal/config"
	"dmicro/web/dd/internal/router"
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

	client.Init(svc)
	// Run svc
	if err := svc.Run(); err != nil {
		log.Fatal(err)
	}
}
