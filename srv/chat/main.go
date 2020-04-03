package main

import (
	"dmicro/common/constant"
	"dmicro/common/log"
	"dmicro/common/service"
	"dmicro/pkg/micro/go-plugins/wrapper/trace/opentracing"
	"dmicro/pkg/tracer"
	"dmicro/srv/chat/internal/broker"
	"dmicro/srv/chat/internal/client"
	"dmicro/srv/chat/internal/config"
	"fmt"
	"github.com/micro/go-micro/v2"
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
	// 添加 trace 和 broker
	opts = append(opts, micro.WrapHandler(opentracing.NewHandlerWrapper(t)))
	opts = append(opts, micro.Broker(broker.GetBroker()))

	svc.Init(opts...)
	client.Init(svc)
	topic := fmt.Sprintf(constant.TOPIC_L2A_PREFIX, 1)
	// 注册 broker 订阅者，这里的消费为 L2A 就是 logic to app server
	micro.RegisterSubscriber(
		topic,
		svc.Server(),
		broker.HandleL2A,
	)

	// Run svc
	if err := svc.Run(); err != nil {
		log.Fatal(err)
	}
}
