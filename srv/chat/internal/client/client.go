package client

import (
	ums "dmicro/srv/ums/api"
	"github.com/micro/go-micro"
)

var (
	UmsClient ums.UmsService
)

func Init(s micro.Service) {
	UmsClient = ums.NewUmsService("", s.Client())
}
