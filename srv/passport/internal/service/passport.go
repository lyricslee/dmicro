package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"dmicro/common/capx"
	"dmicro/common/constant"
	"dmicro/common/errors"
	"dmicro/common/log"
	topic "dmicro/common/proto/topic"
	"dmicro/pkg/tx"
	gid "dmicro/srv/gid/api"
	passport "dmicro/srv/passport/api"
	"dmicro/srv/passport/internal/client"
	"dmicro/srv/passport/internal/dao"
	"dmicro/srv/passport/internal/model"
	user "dmicro/srv/user/api"
)

type PassportService struct {
	engine     *xorm.Engine
	userClient user.UserService
	gidClient  gid.GidService
}

var (
	passportService     *PassportService
	oncePassportService sync.Once
)

func GetPassportService() *PassportService {
	oncePassportService.Do(func() {
		passportService = &PassportService{
			engine:     dao.GetEngine(),
			userClient: client.GetUserClient(),
			gidClient:  client.GetGidClient(),
		}
	})
	return passportService
}

func (p *PassportService) SmsLogin(ctx context.Context, mobile, code string, appid int) (token *passport.TokenInfo, err error) {
	var (
		user *model.User
	)
	// TODO: 校验验证码
	if user, err = dao.GetUserRepo().GetByMobile(mobile); err != nil {
		log.Error(err)
		return
	}

	if user == nil {
		token, err = p.register(ctx, mobile)
	} else {
		if token, err = p.updateToken(ctx, user.Id, appid); err != nil {
			return nil, err
		}
	}

	return
}

func (p *PassportService) Login(ctx context.Context, mobile, passwd string, appid int) (token *passport.TokenInfo, err error) {
	var (
		user *model.User
	)
	if user, err = dao.GetUserRepo().GetByMobile(mobile); err != nil {
		log.Error(err)
		return nil, err
	}
	if user == nil {
		return nil, errors.ErrUserNotExists
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Passwd), []byte(passwd)); err != nil {
		return nil, errors.ErrPasswordError
	}
	if token, err = p.updateToken(ctx, user.Id, appid); err != nil {
		return nil, err
	}
	return
}

func (p *PassportService) ValidateToken(ctx context.Context, uid int64, token string) (err error) {
	log.Debugf("uid=%d token=%s", uid, token)
	var ut *model.UserToken

	if ut, err = dao.GetUserTokenRepo().GetByAccessToken(token); err != nil {
		return
	}
	if ut == nil || ut.Uid != uid {
		return errors.ErrInvalidToken
	}

	return
}

func (p *PassportService) SetPwd(ctx context.Context, uid int64, passwd string, appid int) (token *passport.TokenInfo, err error) {
	var (
		user       *model.User
		userToken  *model.UserToken
		passwdHash []byte
	)
	if user, err = dao.GetUserRepo().Get(uid); err != nil {
		return
	}

	if user == nil {
		return nil, errors.ErrUserNotExists
	}

	if passwdHash, err = bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost); err != nil {
		return
	}
	user.Passwd = string(passwdHash)

	if userToken, err = dao.GetUserTokenRepo().GetByUidAndAppId(uid, appid); err != nil {
		return
	}
	userToken.AccessToken = uuid.New().String()
	userToken.RefreshToken = uuid.New().String()

	session := p.engine.NewSession()
	defer func() {
		session.Close()
	}()

	if err = session.Begin(); err != nil {
		log.Error(err)
		return
	}

	if _, err = session.ID(user.Id).Update(user); err != nil {
		session.Rollback()
		log.Error(err)
		return
	}

	if _, err = session.ID(userToken.Id).Update(userToken); err != nil {
		session.Rollback()
		log.Error(err)
		return
	}

	if err = session.Commit(); err != nil {
		log.Error(err)
		return
	}

	token = &passport.TokenInfo{
		Uid:          userToken.Uid,
		Token:        userToken.AccessToken,
		RefreshToken: userToken.RefreshToken,
		ExpiredAt:    time.Now().Unix() + 8640000,
	}

	return
}

func (p *PassportService) updateToken(ctx context.Context, uid int64, appid int) (token *passport.TokenInfo, err error) {
	var (
		userToken *model.UserToken
	)
	userToken, err = dao.GetUserTokenRepo().GetByUidAndAppId(uid, appid)
	if err != nil {
		return
	}

	userToken.AccessToken = uuid.New().String()
	userToken.RefreshToken = uuid.New().String()
	userToken.ExpiresIn = 8640000
	if err = dao.GetUserTokenRepo().Update(userToken); err != nil {
		log.Error(err)
		return
	}

	token = &passport.TokenInfo{
		Uid:          userToken.Uid,
		Token:        userToken.AccessToken,
		RefreshToken: userToken.RefreshToken,
		ExpiredAt:    time.Now().Unix() + 8640000,
	}

	return
}

func (p *PassportService) register(ctx context.Context, mobile string) (*passport.TokenInfo, error) {
	rsp, err := p.gidClient.GetMulti(ctx, &gid.MultiRequest{Count: 3})
	if err != nil {
		return nil, err
	}

	if len(rsp.Ids) != 3 {
		return nil, fmt.Errorf("the number of ids dose not match expected %v got %v", 3, len(rsp.Ids))
	}

	session := p.engine.NewSession()
	defer func() {
		session.Close()
	}()

	var (
		tokenInfo *passport.TokenInfo
		msg       *topic.UserCreated
	)
	// 执行本地事务
	err = tx.ExecWithTransaction(session, func(session *xorm.Session) error {
		// user表
		u := &model.User{
			Id:     rsp.Ids[0],
			Mobile: mobile,
		}
		if _, err := session.InsertOne(u); err != nil {
			return err
		}

		u1 := uuid.New().String()
		u2 := uuid.New().String()
		ut := &model.UserToken{
			Id:           rsp.Ids[1],
			AppId:        1,
			Uid:          u.Id,
			ExpiresIn:    8640000,
			AccessToken:  u1,
			RefreshToken: u2,
		}

		if _, err := session.InsertOne(ut); err != nil {
			return err
		}

		// 分布式事务处理
		// 发布topic.user.created事件
		msg = &topic.UserCreated{
			Id:    rsp.Ids[2],
			Topic: constant.TOPIC_USER_CREATED,
			Info:  &topic.UserInfo{Uid: u.Id, Mobile: mobile},
		}

		if err = capx.TxStorePublished(session, rsp.Ids[2], constant.TOPIC_USER_CREATED, msg); err != nil {
			log.Error(err)
		}

		tokenInfo = &passport.TokenInfo{
			Uid:          u.Id,
			Token:        u1,
			RefreshToken: u2,
			ExpiredAt:    time.Now().Unix() + ut.ExpiresIn,
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// 发布消息，发布消息出错不需要处理
	capx.Publish(rsp.Ids[2], constant.TOPIC_USER_CREATED, msg)

	return tokenInfo, nil
}
