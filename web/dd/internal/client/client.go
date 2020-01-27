package client

import (
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/client/selector"
	"github.com/micro/go-micro/web"

	passport "dmicro/srv/passport/api"
)

var (
	PassportClient passport.PassportService
)

func Init(service web.Service) {
	cli := client.NewClient(
		client.Registry(service.Options().Registry),
		client.Selector(selector.NewSelector(selector.Registry(service.Options().Registry))),
	)
	PassportClient = passport.NewPassportService("", cli)
}
