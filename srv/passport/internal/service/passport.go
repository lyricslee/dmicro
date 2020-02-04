package service

import (
	"context"
	"dmicro/common/typ"
	"dmicro/common/util"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
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

func (p *PassportService) SmsLogin(ctx context.Context, mobile, code string) (token *passport.TokenInfo, err error) {
	var (
		header *typ.Header
		user   *model.User
	)
	if header, err = util.GetHeaderFromContext(ctx); err != nil {
		return
	}
	// TODO: 校验验证码
	if user, err = dao.GetUserRepo().GetByMobile(mobile); err != nil {
		log.Error(err)
		return
	}

	if user == nil {
		token, err = p.register(ctx, header.Appid, header.Plat, mobile)
	} else {
		if token, err = p.updateToken(header.Appid, user.Id, header.Plat); err != nil {
			return nil, err
		}
	}

	return
}

func (p *PassportService) Login(ctx context.Context, mobile, passwd string) (token *passport.TokenInfo, err error) {
	var (
		header *typ.Header
		user   *model.User
	)
	if header, err = util.GetHeaderFromContext(ctx); err != nil {
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
	if token, err = p.updateToken(header.Appid, user.Id, header.Plat); err != nil {
		return nil, err
	}
	return
}

func (p *PassportService) AuthToken(ctx context.Context) (err error) {
	var (
		header *typ.Header
		token  string
	)
	if header, err = util.GetHeaderFromContext(ctx); err != nil {
		log.Error(err)
		return
	}

	log.Debugf("appid=%d uid=%d plat=%d token=%s", header.Appid, header.Uid, header.Plat, header.Token)
	tokenKey := fmt.Sprintf(constant.REDIS_KEY_TOKEN, header.Appid, header.Uid, header.Plat)
	log.Debug(tokenKey)
	client := dao.GetRedisClient()
	if token, err = client.Get(tokenKey).Result(); err != nil {
		log.Error(err)
		return
	}
	t := &passport.TokenInfo{}
	if err = json.Unmarshal([]byte(token), t); err != nil {
		log.Error(err)
		return err
	}
	if t.Uid != header.Uid || t.Token != header.Token {
		err = errors.ErrInvalidToken
		return
	}

	return
}

func (p *PassportService) SetPwd(ctx context.Context, passwd string) (tokenInfo *passport.TokenInfo, err error) {
	var (
		header     *typ.Header
		user       *model.User
		passwdHash []byte
	)
	header, err = util.GetHeaderFromContext(ctx)
	if err != nil {
		return
	}

	if user, err = dao.GetUserRepo().Get(header.Uid); err != nil {
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

	tokenInfo, err = p.updateToken(header.Appid, header.Uid, header.Plat)

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
		if err := pipe.Expire(tokenKey, time.Duration(2*3600)*time.Second).Err(); err != nil {
			return err
		}
		if err := pipe.Expire(refreshTokenKey, time.Duration(30*24)*time.Hour).Err(); err != nil {
			return err
		}
		return nil
	})
	return
}

func (p *PassportService) updateToken(appid int, uid int64, plat int) (tokenInfo *passport.TokenInfo, err error) {
	token := uuid.New().String()
	refreshToken := uuid.New().String()

	tokenInfo = &passport.TokenInfo{
		Uid:          uid,
		Token:        token,
		RefreshToken: refreshToken,
		ExpiredAt:    time.Now().Unix() + 2*3600,
	}

	if b, err := json.Marshal(tokenInfo); err != nil {
		return nil, err
	} else {
		tokenKey := fmt.Sprintf(constant.REDIS_KEY_TOKEN, appid, uid, plat)
		refreshTokenKey := fmt.Sprintf(constant.REDIS_KEY_REFRESH_TOKEN, appid, uid, plat)
		token = string(b)
		err = p.setToken(tokenKey, token, refreshTokenKey, refreshToken)

	}

	return
}

func (p *PassportService) register(ctx context.Context, appid int, plat int, mobile string) (*passport.TokenInfo, error) {
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

		token := uuid.New().String()
		refreshToken := uuid.New().String()

		tokenInfo = &passport.TokenInfo{
			Uid:          u.Id,
			Token:        token,
			RefreshToken: refreshToken,
			ExpiredAt:    time.Now().Unix() + 2*3600,
		}

		if b, err := json.Marshal(tokenInfo); err != nil {
			return err
		} else {
			tokenKey := fmt.Sprintf(constant.REDIS_KEY_TOKEN, appid, rsp.Ids[0], plat)
			refreshTokenKey := fmt.Sprintf(constant.REDIS_KEY_REFRESH_TOKEN, appid, rsp.Ids[0], plat)
			token = string(b)

			if err := p.setToken(tokenKey, token, refreshTokenKey, refreshToken); err != nil {
				return err
			}
		}

		// 分布式事务处理
		// 发布topic.user.created事件
		msg = &topic.UserCreated{
			Id:    rsp.Ids[2],
			Topic: constant.TOPIC_USER_CREATED,
			Info:  &topic.UserInfo{Uid: u.Id, Mobile: mobile},
		}

		if err = capx.TxStorePublished(session, rsp.Ids[2], constant.TOPIC_USER_CREATED, msg); err != nil {
			tokenKey := fmt.Sprintf(constant.REDIS_KEY_TOKEN, appid, rsp.Ids[0], plat)
			refreshTokenKey := fmt.Sprintf(constant.REDIS_KEY_REFRESH_TOKEN, appid, rsp.Ids[0], plat)
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
	_ = capx.Publish(rsp.Ids[2], constant.TOPIC_USER_CREATED, msg)

	return tokenInfo, nil
}
