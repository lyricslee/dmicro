package client

import (
	ums "dmicro/srv/ums/api"
	"github.com/micro/go-micro/v2"
)

var (
	UmsClient ums.UmsService
)

func Init(s micro.Service) {
	UmsClient = ums.NewUmsService("go.micro.srv.ums", s.Client())
}
