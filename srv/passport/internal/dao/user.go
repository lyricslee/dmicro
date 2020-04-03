package dao

import (
	"sync"

	"dmicro/srv/passport/internal/model"
)

type UserRepo struct {
}

// User 这个表的 ORM 操作
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

func (r *UserRepo) Get(id int64) (user *model.User, err error) {
	var (
		has bool
	)
	user = &model.User{Id: id}
	if has, err = GetEngine().Get(user); err != nil {
		return nil, err
	} else if !has {
		return nil, nil
	}
	return user, nil
}

func (r *UserRepo) GetByMobile(mobile string) (user *model.User, err error) {
	var (
		has bool
	)
	user = &model.User{Mobile: mobile}
	if has, err = GetEngine().Get(user); err != nil {
		return nil, err
	} else if !has {
		return nil, nil
	}
	return user, nil
}

func (r *UserRepo) Update(user *model.User) (err error) {
	if _, err = GetEngine().Update(user); err != nil {
		return
	}
	return
}
