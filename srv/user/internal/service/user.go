package service

import (
	"sync"

	"dmicro/srv/user/internal/dao"
	"dmicro/srv/user/internal/model"
)

type UserService struct {
	userRepo *dao.UserRepo
}

var (
	userService     *UserService
	onceUserService sync.Once
)

func GetUserService() *UserService {
	onceUserService.Do(func() {
		userService = &UserService{}
	})
	return userService
}

// user Create
func (s *UserService) Create(uid int64, mobile string) error {
	u := &model.User{Id: uid, Mobile: mobile}
	return s.userRepo.Add(u)
}
