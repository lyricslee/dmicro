package handler

import (
	"context"
	"sync"

	"dmicro/common/log"
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
		log.Debug(err)
		return
	}
	log.Debugf("uid=%d mobile=%s", req.Uid, req.Mobile)
	return
}
