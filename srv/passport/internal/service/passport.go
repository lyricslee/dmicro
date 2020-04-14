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

// token 的有效期，刷新时间。
const (
	TokenExpiresIn        = 2 * 3600
	RefreshTokenExpiresIn = 30 * 24 * 3600
)

// passport 依赖于 user 和 gid 这两个 service
type PassportService struct {
	engine     *xorm.Engine
	userClient user.UserService
	gidClient  gid.GidService // generate ids
}

var (
	passportService     *PassportService
	oncePassportService sync.Once
)

// 初始化 passportService
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

// SmsLogin 短信登录
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
	// 根据手机号获取用户信息
	if user, err = dao.GetUserRepo().GetByMobile(mobile); err != nil {
		log.Error(err)
		return
	}

	// 用户不存在注册新用户，存在更新 token 踢掉其它登录的设备
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

// 登录，验证 password 成功后，更新双 token 并且返回给客户端。
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
	// 登录的时候校验密码，CompareHashAndPassword 内部做了 salt 处理。
	if err = bcrypt.CompareHashAndPassword([]byte(user.Passwd), []byte(passwd)); err != nil {
		return nil, errors.ErrPasswordError
	}
	// TODO：踢掉其它登录的设备
	if token, err = p.updateToken(appid, user.Id, plat); err != nil {
		return nil, err
	}
	return
}

// 建议 token 是否过期
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

	// jwt 是 rsa 非对称加密的验证
	userInfo, err := jwt.Decode(config.AuthPublicKey, md["Token"])
	if err != nil {
		log.Error(err)
		return
	}

	tokenKey := fmt.Sprintf(constant.RedisKeyToken, userInfo.Appid, userInfo.Uid, userInfo.Plat)
	log.Debug(tokenKey)

	// 从 redis 中获取 tokenKey 并且作对比
	client := dao.GetRedisClient()
	if token, err = client.Get(tokenKey).Result(); err != nil {
		log.Error(err)
		return
	}

	// 对比 redis 中的 token 和内存中 session token 数据
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

//  更新 password，之后还需要更新数据库中的 token
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

	// 更新 password 之后还需要更新数据库中的 token
	log.Debugf("appid=%d uid=%d plat=%d", md.Appid, md.Uid, md.Plat)
	tokenInfo, err = p.updateToken(md.Appid, md.Uid, md.Plat)

	if err = session.Commit(); err != nil {
		log.Error(err)
		return
	}

	return
}

// 讲 2 个 token 存在 redis 数据库中， (tokenKey:token, refreshTokenKey:refreshToken)
// redis MSet 同时为多个键设置值， MSET 是一个原子性(atomic)操作， 所有给定键都会在同一时间内被设置，
// 不会出现某些键被设置了但是另一些键没有被设置的情况。
func (p *PassportService) setToken(tokenKey, token, refreshTokenKey, refreshToken string) (err error) {
	client := dao.GetRedisClient()
	_, err = client.Pipelined(func(pipe redis.Pipeliner) error {
		if err := pipe.MSet(tokenKey, token, refreshTokenKey, refreshToken).Err(); err != nil {
			return err
		}
		// 设置好 k:v 之后再分别设置 2 个 token 的有效期
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

// 更新双 token
func (p *PassportService) updateToken(appid int, uid int64, plat int) (tokenInfo *passport.TokenInfo, err error) {
	var (
		token        string
		refreshToken string
	)
	// 1. 生成新的 token 用 private key
	token, err = jwt.Encode(
		config.AuthPrivateKey,
		&jwt.UserInfo{Appid: appid, Uid: uid, Plat: plat},
		TokenExpiresIn,
	)
	if err != nil {
		return
	}

	// 2. 这里还生成了 refreshToken，通过 token+refreshToken 的组合来控制客户端有效登录时长。
	refreshToken, err = jwt.Encode(
		config.AuthPrivateKey,
		&jwt.UserInfo{Appid: appid, Uid: uid, Plat: plat},
		RefreshTokenExpiresIn,
	)
	if err != nil {
		return
	}

	// 3. token+refreshToken 双 token 设置
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

// 用户注册
func (p *PassportService) register(ctx context.Context, appid int, plat int, mobile string) (*passport.TokenInfo, error) {
	// 从 gid 服务中获取 id，一个 id 当做 user 表 id 一个当做消息表 topic id
	rsp, err := p.gidClient.GetMulti(ctx, &gid.MultiRequest{Count: 2})
	if err != nil {
		return nil, err
	}

	if len(rsp.Ids) != 2 {
		return nil, fmt.Errorf("the number of ids dose not match expected %v got %v", 2, len(rsp.Ids))
	}

	// Mysql new session
	session := p.engine.NewSession()
	defer func() {
		session.Close()
	}()

	var (
		tokenInfo *passport.TokenInfo
		msg       *topic.UserCreated
	)
	// 执行本地事务
	// 事务包含：1. 创建用户 2. 向 redis 中插入双 token 3. 本地事务消息表插入该事务
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
		// 向 redis 中插入双 token
		if err := p.setToken(tokenKey, token, refreshTokenKey, refreshToken); err != nil {
			return err
		}

		tokenInfo = &passport.TokenInfo{
			Token:        token,
			RefreshToken: refreshToken,
			ExpiresAt:    time.Now().Unix() + TokenExpiresIn,
		}

		// TODO: 分布式事务处理演示：上面代码已经插入 user 了，这里实际上不需要重复插入。
		// 这里的 user 事务消息仅仅是为了演示分布式事务。
		// 发布 topic.user.created 事件
		msg = &topic.UserCreated{
			Id:    rsp.Ids[1],
			Topic: constant.TOPIC_USER_CREATED,
			Info:  &topic.UserInfo{Uid: u.Id, Mobile: mobile},
		}

		// 本地事务消息表插入该事务，如果插入失败从 redis 中删除双 token
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
	// 用户注册成功后返回 token 信息
	return tokenInfo, nil
}
