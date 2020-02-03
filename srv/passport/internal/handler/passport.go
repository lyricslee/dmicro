package handler

import (
	"context"
	"sync"

	"dmicro/common/log"
	"dmicro/common/typ"
	"dmicro/common/util"
	passport "dmicro/srv/passport/api"
	"dmicro/srv/passport/internal/service"
)

type PassportHandler struct {
	svc *service.PassportService
}

var (
	passportHandler     *PassportHandler
	oncePassportHandler sync.Once
)

func GetPassportHandler() *PassportHandler {
	oncePassportHandler.Do(func() {
		passportHandler = &PassportHandler{svc: service.GetPassportService()}
	})
	return passportHandler
}

func (h *PassportHandler) Sms(ctx context.Context, req *passport.Request, rsp *passport.Response) error {
	log.Debugf("Sms: mobile=%s", req.Mobile)
	// TODO: 通过短信服务获取验证码
	rsp.Code = "8888"
	return nil
}

func (h *PassportHandler) SmsLogin(ctx context.Context, req *passport.SmsLoginRequest, rsp *passport.SmsLoginResponse) (err error) {
	log.Debugf("SmsLogin: mobile=%s code=%s", req.Mobile, req.Code)
	var (
		header *typ.Header
	)
	if header, err = util.GetHeaderFromContext(ctx); err != nil {
		return
	}
	rsp.TokenInfo, err = h.svc.SmsLogin(ctx, req.Mobile, req.Code, header.AppId)
	return
}

func (h *PassportHandler) Login(ctx context.Context, req *passport.LoginRequest, rsp *passport.LoginResponse) (err error) {
	log.Debugf("Login: mobile=%s passwd=%s", req.Mobile, req.Passwd)
	var (
		header *typ.Header
	)
	header, err = util.GetHeaderFromContext(ctx)
	if err != nil {
		return err
	}
	if rsp.TokenInfo, err = h.svc.Login(ctx, req.Mobile, req.Passwd, header.AppId); err != nil {
		return err
	}
	return
}

func (h *PassportHandler) OAuthLogin(ctx context.Context, req *passport.OAuthLoginRequest, rsp *passport.OAuthLoginResponse) error {
	return nil
}

func (h *PassportHandler) SetPwd(ctx context.Context, req *passport.SetPwdRequest, rsp *passport.SetPwdResponse) error {
	log.Debug("SetPwd...")
	header, err := util.GetHeaderFromContext(ctx)
	if err != nil {
		return err
	}

	if rsp.TokenInfo, err = h.svc.SetPwd(ctx, header.Uid, req.Passwd, header.AppId); err != nil {
		return err
	}
	return nil
}

func (h *PassportHandler) AuthToken(ctx context.Context, req *passport.AuthTokenRequest, rsp *passport.AuthTokenResponse) error {
	log.Debug("AuthToken...")
	header, err := util.GetHeaderFromContext(ctx)
	if err != nil {
		return err
	}
	log.Debugf("AuthToken: uid=%d token=%s", header.Uid, header.Token)
	return h.svc.ValidateToken(ctx, header.Uid, header.Token)
}

func (h *PassportHandler) AuthCookie(ctx context.Context, req *passport.AuthCookieRequest, rsp *passport.AuthCookieResponse) error {
	// TODO
	return nil
}
