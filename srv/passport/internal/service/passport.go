package service

import (
	"context"
	"dmicro/common/typ"
	"dmicro/common/util"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
	"github.com/micro/go-micro/v2/metadata"
	"golang.org/x/crypto/bcrypt"

	"dmicro/common/capx"
	"dmicro/common/constant"
	"dmicro/common/errors"
	"dmicro/common/log"
	topic "dmicro/common/proto/topic"
	"dmicro/common/util/jwt"
	"dmicro/pkg/tx"
	gid "dmicro/srv/gid/api"
	passport "dmicro/srv/passport/api"
	"dmicro/srv/passport/internal/client"
	"dmicro/srv/passport/internal/config"
	"dmicro/srv/passport/internal/dao"
	"dmicro/srv/passport/internal/model"
	user "dmicro/srv/user/api"
)

const (
	TokenExpiresIn        = 2 * 3600
	RefreshTokenExpiresIn = 30 * 24 * 3600
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

func (p *PassportService) SmsLogin(ctx context.Context, appid, plat int, mobile, code string) (token *passport.TokenInfo, err error) {
	var (
		user *model.User
	)
	if appid == 0 || plat == 0 {
		err = fmt.Errorf("appid or plat can't be zero")
		log.Error(err)
		return
	}
	// TODO: 校验验证码
	if user, err = dao.GetUserRepo().GetByMobile(mobile); err != nil {
		log.Error(err)
		return
	}

	if user == nil {
		token, err = p.register(ctx, appid, plat, mobile)
	} else {
		// TODO：踢掉其它登录的设备
		if token, err = p.updateToken(appid, user.Id, plat); err != nil {
			return nil, err
		}
	}

	return
}

func (p *PassportService) Login(ctx context.Context, appid, plat int, mobile, passwd string) (token *passport.TokenInfo, err error) {
	var (
		user *model.User
	)
	if appid == 0 || plat == 0 {
		err = fmt.Errorf("appid or plat can't be zero")
		log.Error(err)
		return
	}

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
	// TODO：踢掉其它登录的设备
	if token, err = p.updateToken(appid, user.Id, plat); err != nil {
		return nil, err
	}
	return
}

func (p *PassportService) AuthToken(ctx context.Context) (rsp *passport.AuthTokenResponse, err error) {
	var (
		token string
	)

	md, ok := metadata.FromContext(ctx)
	if !ok {
		log.Error("metadata.FromContext error")
		err = fmt.Errorf("AuthToken: metadata.FromContext error")
		return
	}

	userInfo, err := jwt.Decode(config.AuthPublicKey, md["Token"])
	if err != nil {
		log.Error(err)
		return
	}

	tokenKey := fmt.Sprintf(constant.RedisKeyToken, userInfo.Appid, userInfo.Uid, userInfo.Plat)
	log.Debug(tokenKey)

	client := dao.GetRedisClient()
	if token, err = client.Get(tokenKey).Result(); err != nil {
		log.Error(err)
		return
	}

	if token != md["Token"] {
		err = errors.ErrInvalidToken
		log.Debug(err)
		return
	}

	log.Debugf("appid=%d uid=%d plat=%d", userInfo.Appid, userInfo.Uid, userInfo.Plat)
	rsp = &passport.AuthTokenResponse{
		Appid: int32(userInfo.Appid),
		Uid:   userInfo.Uid,
		Plat:  int32(userInfo.Plat),
	}
	return
}

func (p *PassportService) SetPwd(ctx context.Context, passwd string) (tokenInfo *passport.TokenInfo, err error) {
	var (
		md         *typ.MetaData
		user       *model.User
		passwdHash []byte
	)

	md, err = util.GetMetaDataFromContext(ctx)
	if err != nil {
		return
	}

	if md.Appid == 0 || md.Uid == 0 || md.Plat == 0 {
		err = fmt.Errorf("appid or uid or plat can't be zero.")
		log.Error(err)
		return
	}

	if user, err = dao.GetUserRepo().Get(md.Uid); err != nil {
		log.Error(err)
		return
	}

	if user == nil {
		return nil, errors.ErrUserNotExists
	}

	if passwdHash, err = bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost); err != nil {
		log.Error(err)
		return
	}
	user.Passwd = string(passwdHash)

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

	log.Debugf("appid=%d uid=%d plat=%d", md.Appid, md.Uid, md.Plat)
	tokenInfo, err = p.updateToken(md.Appid, md.Uid, md.Plat)

	if err = session.Commit(); err != nil {
		log.Error(err)
		return
	}

	return
}

func (p *PassportService) setToken(tokenKey, token, refreshTokenKey, refreshToken string) (err error) {
	client := dao.GetRedisClient()
	_, err = client.Pipelined(func(pipe redis.Pipeliner) error {
		if err := pipe.MSet(tokenKey, token, refreshTokenKey, refreshToken).Err(); err != nil {
			return err
		}
		if err := pipe.Expire(tokenKey, time.Duration(TokenExpiresIn)*time.Second).Err(); err != nil {
			return err
		}
		if err := pipe.Expire(refreshTokenKey, time.Duration(RefreshTokenExpiresIn)*time.Second).Err(); err != nil {
			return err
		}
		return nil
	})
	return
}

func (p *PassportService) updateToken(appid int, uid int64, plat int) (tokenInfo *passport.TokenInfo, err error) {
	var (
		token        string
		refreshToken string
	)
	token, err = jwt.Encode(
		config.AuthPrivateKey,
		&jwt.UserInfo{Appid: appid, Uid: uid, Plat: plat},
		TokenExpiresIn,
	)
	if err != nil {
		return
	}

	refreshToken, err = jwt.Encode(
		config.AuthPrivateKey,
		&jwt.UserInfo{Appid: appid, Uid: uid, Plat: plat},
		RefreshTokenExpiresIn,
	)
	if err != nil {
		return
	}

	tokenInfo = &passport.TokenInfo{
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Unix() + TokenExpiresIn,
	}

	tokenKey := fmt.Sprintf(constant.RedisKeyToken, appid, uid, plat)
	refreshTokenKey := fmt.Sprintf(constant.RedisKeyRefreshToken, appid, uid, plat)
	err = p.setToken(tokenKey, token, refreshTokenKey, refreshToken)

	return
}

func (p *PassportService) register(ctx context.Context, appid int, plat int, mobile string) (*passport.TokenInfo, error) {
	rsp, err := p.gidClient.GetMulti(ctx, &gid.MultiRequest{Count: 2})
	if err != nil {
		return nil, err
	}

	if len(rsp.Ids) != 2 {
		return nil, fmt.Errorf("the number of ids dose not match expected %v got %v", 2, len(rsp.Ids))
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

		token, err := jwt.Encode(
			config.AuthPrivateKey,
			&jwt.UserInfo{Appid: appid, Uid: rsp.Ids[0], Plat: plat},
			TokenExpiresIn,
		)
		if err != nil {
			return err
		}

		refreshToken, err := jwt.Encode(
			config.AuthPrivateKey,
			&jwt.UserInfo{Appid: appid, Uid: rsp.Ids[0], Plat: plat},
			RefreshTokenExpiresIn,
		)
		if err != nil {
			return err
		}
		tokenKey := fmt.Sprintf(constant.RedisKeyToken, appid, rsp.Ids[0], plat)
		refreshTokenKey := fmt.Sprintf(constant.RedisKeyRefreshToken, appid, rsp.Ids[0], plat)
		if err := p.setToken(tokenKey, token, refreshTokenKey, refreshToken); err != nil {
			return err
		}

		tokenInfo = &passport.TokenInfo{
			Token:        token,
			RefreshToken: refreshToken,
			ExpiresAt:    time.Now().Unix() + TokenExpiresIn,
		}

		// 分布式事务处理
		// 发布topic.user.created事件
		msg = &topic.UserCreated{
			Id:    rsp.Ids[1],
			Topic: constant.TOPIC_USER_CREATED,
			Info:  &topic.UserInfo{Uid: u.Id, Mobile: mobile},
		}

		if err = capx.TxStorePublished(session, rsp.Ids[1], constant.TOPIC_USER_CREATED, msg); err != nil {
			client := dao.GetRedisClient()
			client.Del(tokenKey, refreshTokenKey)
			log.Error(err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// 发布消息，发布消息出错不需要处理
	_ = capx.Publish(rsp.Ids[1], constant.TOPIC_USER_CREATED, msg)

	return tokenInfo, nil
}
