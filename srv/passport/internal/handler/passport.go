package handler

import (
	"context"
	"sync"

	"dmicro/common/log"
	passport "dmicro/srv/passport/api"
	"dmicro/srv/passport/internal/service"
)

// handler 实际上是调用了对应的额 passport service API
// handler 调用 service 里面的具体方法来执行操作, 这样做的好处就是把接口和实现分离。类似于 http 中的 router 作用。
type PassportHandler struct {
	svc *service.PassportService
}

var (
	passportHandler     *PassportHandler
	oncePassportHandler sync.Once // 只初始化一次
)

// 只初始化一次
func GetPassportHandler() *PassportHandler {
	oncePassportHandler.Do(func() {
		passportHandler = &PassportHandler{svc: service.GetPassportService()}
	})
	return passportHandler
}

// api/passport.proto 对应的各个协议字段和接口
func (h *PassportHandler) Sms(ctx context.Context, req *passport.Request, rsp *passport.Response) error {
	log.Debugf("Sms: mobile=%s", req.Mobile)
	// TODO: 通过短信服务获取验证码
	rsp.Code = "8888"
	return nil
}

func (h *PassportHandler) SmsLogin(ctx context.Context, req *passport.SmsLoginRequest, rsp *passport.SmsLoginResponse) (err error) {
	log.Debugf("SmsLogin: mobile=%s code=%s", req.Mobile, req.Code)
	rsp.TokenInfo, err = h.svc.SmsLogin(ctx, int(req.Appid), int(req.Plat), req.Mobile, req.Code)
	return
}

func (h *PassportHandler) Login(ctx context.Context, req *passport.LoginRequest, rsp *passport.LoginResponse) (err error) {
	log.Debugf("Login: mobile=%s passwd=%s", req.Mobile, req.Passwd)
	rsp.TokenInfo, err = h.svc.Login(ctx, int(req.Appid), int(req.Plat), req.Mobile, req.Passwd)

	return
}

func (h *PassportHandler) OAuthLogin(ctx context.Context, req *passport.OAuthLoginRequest, rsp *passport.OAuthLoginResponse) error {
	return nil
}

func (h *PassportHandler) SetPwd(ctx context.Context, req *passport.SetPwdRequest, rsp *passport.SetPwdResponse) (err error) {
	log.Debug("SetPwd...")
	rsp.TokenInfo, err = h.svc.SetPwd(ctx, req.Passwd)
	return nil
}

func (h *PassportHandler) AuthToken(ctx context.Context, req *passport.AuthTokenRequest, rsp *passport.AuthTokenResponse) error {
	log.Debug("AuthToken...")
	r, err := h.svc.AuthToken(ctx)
	if r != nil {
		rsp.Appid = r.Appid
		rsp.Uid = r.Uid
		rsp.Plat = r.Plat
	}
	return err
}

func (h *PassportHandler) AuthCookie(ctx context.Context, req *passport.AuthCookieRequest, rsp *passport.AuthCookieResponse) (err error) {
	// TODO
	return
}
