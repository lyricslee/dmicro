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


// 1. dmicro 是基于 go-micro 开发的，所以这里需要 import 它的一些库，比如：web
// 2. opentracing 用作分布式事务追踪

func main() {
	// 初始化配置文件, 初始化日志配置
	// Config
	config.Init()
	// Logger
	log.Init(config.Logger)

	// 1. 创建一个新的 Web Service，内部就是 go-micro.web 的初始化。
	// New Service
	svc := service.NewWebService()

	// 2. go-micro.web service 的配置项, 比如：
	// 注册 web 服务的路由 API 与对应的 handler
	var opts []web.Option
	opts = append(opts, web.Handler(router.Router()))

	// 3. 注册完成配置项之后初始化 web services
	if err := svc.Init(opts...); err != nil {
		log.Fatal(err)
	}

	// 初始化 tracer
	// Tracer
	t, err := tracer.Init(config.Micro.ServerName, config.Tracer.Addr)
	if err != nil {
		log.Error(err)
	}
	// 设置全局 tracer
	opentracing.SetGlobalTracer(t)

	// 初始化 RPC 的 clients，用于之后 RPC 通信。比如 Passport Service
	client.Init(svc)
	// 最后运行这个 service
	// Run svc
	if err := svc.Run(); err != nil {
		log.Fatal(err)
	}
}
