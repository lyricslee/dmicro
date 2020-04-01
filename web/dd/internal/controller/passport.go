package controller

import (
	"dmicro/common/log"
	"dmicro/pkg/context"
	passport "dmicro/srv/passport/api"
	"dmicro/web/dd/internal/client"
)

// 这里定义了 Passport Controller 模块，包含了各个处理函数。
// PassportController ...
type PassportController struct {
}

// context.DmContext 内部 gin.Context 对象，因为 dmicro 用了 gin http 框架。
// gin.Context 是 http 请求参数，我们使用 context.DmContext 对他进行了二次处理，
// Json 请求参数解析等。
func (this *PassportController) Login(mctx *context.DmContext) {
	log.Debug("Login...")
	request := &passport.LoginRequest{}

	// 请求参数 Json 转换成 passport.LoginRequest 对象
	if err := mctx.ParseJSON(request); err != nil {
		log.Error(err)
		mctx.ResponseError(err)
		return
	}
	// TODO: request相关参数校验
	// 这里可以引入专门的请求参数校验来做，比如数据类型，字段等。

	// RPC PassportClient.Login() 方法，做登录。
	// DmContext 转换成 golang.context,
	// tracer 追踪 context 调用, context 包是用来解决，资源多次调用释放的问题。
	// 注意这里 RPC 参数有 2 个： context.Context, LoginRequest
	response, err := client.PassportClient.Login(toContext(mctx), request)
	if err != nil {
		log.Error(err)
		mctx.ResponseError(err)
		return
	}

	// response 对象
	mctx.Response(response)
	return
}

// Sma 短信服务，流程与 Login 登录服务基本一样。
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

// 短信登录
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

// Oauth 登录
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

// Password
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
