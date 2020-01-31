package client

import (
	"github.com/micro/go-micro/v2"

	gid "dmicro/srv/gid/api"
	user "dmicro/srv/user/api"
)

var (
	gidClient  gid.GidService
	userClient user.UserService
)

func Init(s micro.Service) {
	gidClient = gid.NewGidService("go.micro.srv.gid", s.Client())
	userClient = user.NewUserService("go.micro.srv.user", s.Client())
}

func GetGidClient() gid.GidService {
	return gidClient
}

func GetUserClient() user.UserService {
	return userClient
}
