package dao

import (
	"sync"

	"dmicro/srv/passport/internal/model"
)

type UserTokenRepo struct {
}

var (
	userTokenRepo *UserTokenRepo
	onceUserToken sync.Once
)

func GetUserTokenRepo() *UserTokenRepo {
	onceUserToken.Do(func() {
		userTokenRepo = &UserTokenRepo{}
	})
	return userTokenRepo
}

func (r *UserTokenRepo) Add(v *model.UserToken) (err error) {
	if _, err = GetEngine().InsertOne(v); err != nil {
		return
	}
	return
}

func (r *UserTokenRepo) GetByUserIdAndAppId(userId int64, appId int) (v *model.UserToken, err error) {
	var (
		has bool
	)
	v = &model.UserToken{UserId: userId, AppId: appId}
	if has, err = GetEngine().Get(v); err != nil {
		return nil, err
	} else if !has {
		return nil, nil
	}
	return v, nil
}

func (r *UserTokenRepo) Update(v *model.UserToken) (err error) {
	if _, err = GetEngine().ID(v.Id).Update(v); err != nil {
		return
	}
	return
}

func (r *UserTokenRepo) GetByAccessToken(token string) (*model.UserToken, error) {
	v := &model.UserToken{AccessToken: token}
	has, err := GetEngine().Get(v)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, nil
	}
	return v, nil
}
