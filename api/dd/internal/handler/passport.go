package handler

import (
	"context"
	"encoding/json"

	api "github.com/micro/go-micro/v2/api/proto"
	"github.com/micro/go-micro/v2/server"

	pb "dmicro/api/dd/api"
	"dmicro/api/dd/internal/client"
	"dmicro/api/dd/internal/util"
	"dmicro/common/log"
	passport "dmicro/srv/passport/api"
)

type Passport struct{}

// 请求 handler 调用后台 Passport 对应的 RPC handler
func (this *Passport) Sms(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Debug("sms...")

	r := &passport.Request{}
	if err := json.Unmarshal([]byte(req.GetBody()), r); err != nil {
		return err
	}

	var (
		code int
		body string
	)
	response, err := client.PassportClient.Sms(ctx, r)
	if err != nil {
		code, body = util.MakeErrBody(err)
	} else {
		code, body = util.MakeBody(response)
	}

	rsp.StatusCode = int32(code)
	rsp.Body = body
	return nil
}

func (this *Passport) SmsLogin(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Debug("smslogin ...")
	r := &passport.SmsLoginRequest{}
	if err := json.Unmarshal([]byte(req.GetBody()), r); err != nil {
		return err
	}

	var (
		code int
		body string
	)
	response, err := client.PassportClient.SmsLogin(ctx, r)
	if err != nil {
		code, body = util.MakeErrBody(err)
	} else {
		code, body = util.MakeBody(response)
	}

	rsp.StatusCode = int32(code)
	rsp.Body = body

	return nil
}

func (this *Passport) SetPwd(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Debug("SetPwd ...")
	r := &passport.SetPwdRequest{}
	if err := json.Unmarshal([]byte(req.GetBody()), r); err != nil {
		return err
	}

	var (
		code int
		body string
	)
	response, err := client.PassportClient.SetPwd(ctx, r)
	if err != nil {
		code, body = util.MakeErrBody(err)
	} else {
		code, body = util.MakeBody(response)
	}

	rsp.StatusCode = int32(code)
	rsp.Body = body
	return nil
}

func (this *Passport) Login(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Debug("Login ...")
	r := &passport.LoginRequest{}
	if err := json.Unmarshal([]byte(req.GetBody()), r); err != nil {
		return err
	}

	var (
		code int
		body string
	)
	response, err := client.PassportClient.Login(ctx, r)
	if err != nil {
		code, body = util.MakeErrBody(err)
	} else {
		code, body = util.MakeBody(response)
	}

	rsp.StatusCode = int32(code)
	rsp.Body = body
	return nil
}

func (this *Passport) OAuthLogin(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Debug("OAuthLogin ...")
	r := &passport.OAuthLoginRequest{}
	if err := json.Unmarshal([]byte(req.GetBody()), r); err != nil {
		return err
	}

	var (
		code int
		body string
	)
	response, err := client.PassportClient.OAuthLogin(ctx, r)
	if err != nil {
		code, body = util.MakeErrBody(err)
	} else {
		code, body = util.MakeBody(response)
	}

	rsp.StatusCode = int32(code)
	rsp.Body = body
	return nil
}

func registerPassportHandler(server server.Server) {
	_ = pb.RegisterPassportHandler(server, new(Passport))
}
