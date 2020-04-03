package main

import (
	"github.com/micro/go-micro/v2"

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

	// 初始化基本与 web/dd 这些类似，注册 RPC 对应的 handlers
	// Register Handler
	gid.RegisterGidHandler(svc.Server(), handler.GetGidHandler())

	// Run svc
	if err := svc.Run(); err != nil {
		log.Fatal(err)
	}
}
