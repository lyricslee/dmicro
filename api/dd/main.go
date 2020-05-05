package main

import (
	"github.com/micro/go-micro/v2"

	"dmicro/api/dd/internal/client"
	"dmicro/api/dd/internal/config"
	"dmicro/api/dd/internal/handler"
	"dmicro/common/log"
	"dmicro/common/service"
	"dmicro/pkg/micro/go-plugins/wrapper/trace/opentracing"
	"dmicro/pkg/tracer"
	//passport "dmicro/srv/passport/api"
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
	// 注入 trace 信息
	opts = append(opts, micro.WrapHandler(opentracing.NewHandlerWrapper(t)))

	svc.Init(opts...)

	client.Init(svc)
	// Register Handler
	handler.RegisterHandler(svc.Server())

	// Run svc
	if err := svc.Run(); err != nil {
		log.Fatal(err)
	}
}
