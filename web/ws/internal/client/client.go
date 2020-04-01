package client

import (
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/client/grpc"
	"github.com/micro/go-micro/v2/client/selector"
	"github.com/micro/go-micro/v2/web"

	ums "dmicro/srv/ums/api"
)

var (
	UmsClient ums.UmsService
)

// ums 服务客户端
func Init(service web.Service) {
	cli := grpc.NewClient(
		client.Registry(service.Options().Registry),
		client.Selector(selector.NewSelector(selector.Registry(service.Options().Registry))),
	)
	UmsClient = ums.NewUmsService("go.micro.srv.ums", cli)
}
