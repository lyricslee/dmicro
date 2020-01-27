package dao

import (
	"sync"

	"github.com/go-xorm/xorm"

	"dmicro/srv/user/internal/model"
)

type UserRepo struct {
	Engine *xorm.Engine
}

var (
	userRepo *UserRepo
	onceUser sync.Once
)

func GetUserRepo() *UserRepo {
	onceUser.Do(func() {
		userRepo = &UserRepo{}
	})
	return userRepo
}

func (r *UserRepo) Add(user *model.User) (err error) {
	if _, err = GetEngine().InsertOne(user); err != nil {
		return
	}
	return
}

func (r *UserRepo) Get(id int64) (*model.User, error) {
	return nil, nil
}

func (r *UserRepo) GetByMobile(mobile string) (*model.User, error) {
	return nil, nil
}

func (r *UserRepo) Update(user *model.User) error {
	return nil
}
