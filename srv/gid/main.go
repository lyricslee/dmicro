package main

import (
	"github.com/micro/go-micro"

	"dmicro/common/log"
	"dmicro/common/service"
	"dmicro/pkg/micro/go-plugins/wrapper/trace/opentracing"
	"dmicro/pkg/tracer"
	gid "dmicro/srv/gid/api"
	"dmicro/srv/gid/internal/config"
	"dmicro/srv/gid/internal/handler"
)

func main() {
	// Config
	config.Init()
	// Logger
	log.Init(config.Logger)

	// New Service
	svc := service.NewService()

	var opts []micro.Option
	// Tracer
	t, err := tracer.Init(config.Micro.ServerName, config.Tracer.Addr)
	if err != nil {
		log.Error(err)
	}
	opts = append(opts, micro.WrapHandler(opentracing.NewHandlerWrapper(t)))

	svc.Init(opts...)

	// Register Handler
	gid.RegisterGidHandler(svc.Server(), handler.GetGidHandler())

	// Run svc
	if err := svc.Run(); err != nil {
		log.Fatal(err)
	}
}
