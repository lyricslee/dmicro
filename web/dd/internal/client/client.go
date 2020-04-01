package client

import (
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/client/grpc"
	"github.com/micro/go-micro/v2/client/selector"
	"github.com/micro/go-micro/v2/web"

	passport "dmicro/srv/passport/api"
)

var (
	PassportClient passport.PassportService
)

// client 是 RPC 通信的客户端，比如这里的 PassportService 客户端。
// RPC 使用了 grpc 框架，这里需要初始化 Client 传入服务器的信息
// Selector 这里是负载均衡，因为线上同一个 service 可能运行多个实例。

func Init(service web.Service) {
	cli := grpc.NewClient(
		client.Registry(service.Options().Registry),
		client.Selector(selector.NewSelector(selector.Registry(service.Options().Registry))),
	)
	// 传入 PassportService 服务的名称，go-micro 会去 etcd 中查找对应的服务地址。（服务发现）
	PassportClient = passport.NewPassportService("go.micro.srv.passport", cli)
}
