package main

import (
	"github.com/micro/go-micro"

	"dmicro/common/log"
	"dmicro/common/service"
	"dmicro/pkg/micro/go-plugins/wrapper/trace/opentracing"
	"dmicro/pkg/tracer"
	ums "dmicro/srv/ums/api"
	"dmicro/srv/ums/internal/broker"
	"dmicro/srv/ums/internal/config"
	"dmicro/srv/ums/internal/handler"
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
	opts = append(opts, micro.Broker(broker.GetBroker()))

	svc.Init(opts...)

	// Register Handler
	ums.RegisterUmsHandler(svc.Server(), handler.GetUmsHandler())

	// Run svc
	if err := svc.Run(); err != nil {
		log.Fatal(err)
	}
}
