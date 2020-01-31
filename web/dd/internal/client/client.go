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

func Init(service web.Service) {
	cli := grpc.NewClient(
		client.Registry(service.Options().Registry),
		client.Selector(selector.NewSelector(selector.Registry(service.Options().Registry))),
	)
	PassportClient = passport.NewPassportService("go.micro.srv.passport", cli)
}
