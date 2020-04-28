package main

import (
	"github.com/micro/go-micro/v2"

	"dmicro/common/log"
	"dmicro/common/service"
	"dmicro/pkg/micro/go-plugins/wrapper/trace/opentracing"
	"dmicro/pkg/tracer"
	ums "dmicro/srv/ums/api"
	"dmicro/srv/ums/internal/broker"
	"dmicro/srv/ums/internal/config"
	"dmicro/srv/ums/internal/handler"
)

/*
UMS 主要用来处理一些长连接业务，比如：聊天 推送 客服等。
加 MQ 是为了处理大量请求的时候，后台的服务处理不过来。
同事 MQ 也起到分发消息的作用。
*/
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
	// 注册 trace broker 等
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
