package client

import (
	"github.com/micro/go-micro"

	gid "dmicro/srv/gid/api"
	user "dmicro/srv/user/api"
)

var (
	gidClient  gid.GidService
	userClient user.UserService
)

func Init(s micro.Service) {
	gidClient = gid.NewGidService("", s.Client())
	userClient = user.NewUserService("", s.Client())
}

func GetGidClient() gid.GidService {
	return gidClient
}

func GetUserClient() user.UserService {
	return userClient
}
