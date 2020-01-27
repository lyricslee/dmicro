package controller

import (
	"dmicro/common/log"
	"dmicro/pkg/context"
	passport "dmicro/srv/passport/api"
	"dmicro/web/dd/internal/client"
)

// PassportController ...
type PassportController struct {
}

func (this *PassportController) Login(mctx *context.DmContext) {
	log.Debug("Login...")
	request := &passport.LoginRequest{}
	if err := mctx.ParseJSON(request); err != nil {
		log.Error(err)
		mctx.ResponseError(err)
		return
	}
	// TODO: request相关参数校验
	response, err := client.PassportClient.Login(toContext(mctx), request)
	if err != nil {
		log.Error(err)
		mctx.ResponseError(err)
		return
	}

	mctx.Response(response)
	return
}

func (this *PassportController) Sms(mctx *context.DmContext) {
	log.Debug("Sms...")
	request := &passport.Request{}
	if err := mctx.ParseJSON(request); err != nil {
		log.Error(err)
		mctx.ResponseError(err)
		return
	}
	// TODO: request相关参数校验
	response, err := client.PassportClient.Sms(toContext(mctx), request)
	if err != nil {
		log.Error(err)
		mctx.ResponseError(err)
		return
	}

	mctx.Response(response)
	return
}

func (this *PassportController) SmsLogin(mctx *context.DmContext) {
	log.Debug("SmsLogin...")
	request := &passport.SmsLoginRequest{}
	if err := mctx.ParseJSON(request); err != nil {
		log.Error(err)
		mctx.ResponseError(err)
		return
	}

	// TODO: request相关参数校验

	response, err := client.PassportClient.SmsLogin(toContext(mctx), request)
	if err != nil {
		log.Error(err)
		mctx.ResponseError(err)
		return
	}

	mctx.Response(response)
	return
}

func (this *PassportController) OauthLogin(mctx *context.DmContext) {
	log.Debug("OauthLogin...")
	request := &passport.OAuthLoginRequest{}
	if err := mctx.ParseJSON(request); err != nil {
		log.Error(err)
		mctx.ResponseError(err)
		return
	}
	// TODO: request相关参数校验
	response, err := client.PassportClient.OAuthLogin(toContext(mctx), request)
	if err != nil {
		log.Error(err)
		mctx.ResponseError(err)
		return
	}

	mctx.Response(response)

	return
}

func (this *PassportController) SetPwd(mctx *context.DmContext) {
	log.Debug("SetPwd...")
	request := &passport.SetPwdRequest{}
	if err := mctx.ParseJSON(request); err != nil {
		log.Error(err)
		mctx.ResponseError(err)
		return
	}
	// TODO: request相关参数校验
	response, err := client.PassportClient.SetPwd(toContext(mctx), request)
	if err != nil {
		log.Error(err)
		mctx.ResponseError(err)
		return
	}

	mctx.Response(response)

	return
}
