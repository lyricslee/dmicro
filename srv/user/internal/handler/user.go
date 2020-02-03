package handler

import (
	"context"
	"fmt"
	"sync"

	user "dmicro/srv/user/api"
	"dmicro/srv/user/internal/service"
)

type UserHandler struct {
	svc *service.UserService
}

var (
	userHandler     *UserHandler
	onceUserHandler sync.Once
)

func GetUserHandler() *UserHandler {
	onceUserHandler.Do(func() {
		userHandler = &UserHandler{svc: service.GetUserService()}
	})
	return userHandler
}

func (h *UserHandler) Create(ctx context.Context, req *user.Request, rsp *user.Response) (err error) {
	if err = h.svc.Create(req.Uid, req.Mobile); err != nil {
		fmt.Println("xxxxxxxxxxxx")
		return
	}
	fmt.Println(req.Uid, req.Mobile)
	rsp = &user.Response{}
	return
}
