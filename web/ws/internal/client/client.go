package client

import (
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/client/selector"
	"github.com/micro/go-micro/web"

	ums "dmicro/srv/ums/api"
)

var (
	UmsClient ums.UmsService
)

func Init(service web.Service) {
	cli := client.NewClient(
		client.Registry(service.Options().Registry),
		client.Selector(selector.NewSelector(selector.Registry(service.Options().Registry))),
	)
	UmsClient = ums.NewUmsService("", cli)
}
