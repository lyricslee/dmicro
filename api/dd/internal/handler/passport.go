package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	//mapi "github.com/micro/go-micro/v2/api"
	//hapi "github.com/micro/go-micro/v2/api/handler/api"
	api "github.com/micro/go-micro/v2/api/proto"

	pb "dmicro/api/dd/api"
	"dmicro/api/dd/internal/client"
	"dmicro/common/errors"
	"dmicro/common/log"
	passport "dmicro/srv/passport/api"
	"github.com/micro/go-micro/v2/server"
)

type Passport struct{}

func x(err error) (code int, body string) {
	ce := errors.Parse(err.Error())

	m := make(map[string]interface{})
	m["errno"] = ce.Errno
	m["errmsg"] = ce.Errmsg
	m["t"] = time.Now().UnixNano()

	if ce.Errno == -1 {
		code = http.StatusInternalServerError
	} else {
		code = 499
	}

	b, _ := json.Marshal(m)
	body = string(b)
	return
}

func y(obj interface{}) (code int, body string) {
	if obj == nil {
		obj = make(map[string]interface{})
	}

	m := make(map[string]interface{})
	m["errno"] = 0
	m["data"] = obj
	m["t"] = time.Now().UnixNano()

	code = http.StatusOK
	b, _ := json.Marshal(m)
	body = string(b)

	return
}

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
		code, body = x(err)
	} else {
		code, body = y(response)
	}

	rsp.StatusCode = int32(code)
	rsp.Body = body
	return nil
}

func (this *Passport) SmsLogin(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Debug("smslogin ...")
	return nil
}

func (this *Passport) SetPwd(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Debug("SetPwd ...")
	return nil
}

func (this *Passport) Login(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Debug("Login ...")
	return nil
}

func (this *Passport) OAuthLogin(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Debug("OAuthLogin ...")
	return nil
}

func registerPassportHandler(server server.Server) {

	_ = pb.RegisterPassportHandler(server, new(Passport))

	//_ = pb.RegisterPassportHandler(server, new(Passport),
	//	mapi.WithEndpoint(&mapi.Endpoint{
	//		// The RPC method
	//		Name: "Passport.Sms",
	//		// The HTTP paths. This can be a POSIX regex
	//		Path: []string{"/passport/sendsms"},
	//		// The HTTP Methods for this endpoint
	//		Method: []string{"POST"},
	//		// The API handler to use
	//		Handler: hapi.Handler,
	//	}),
	//	mapi.WithEndpoint(&mapi.Endpoint{
	//		// The RPC method
	//		Name: "Passport.SmsLogin",
	//		// The HTTP paths. This can be a POSIX regex
	//		Path: []string{"/passport/smslogin"},
	//		// The HTTP Methods for this endpoint
	//		Method: []string{"POST"},
	//		// The API handler to use
	//		Handler: hapi.Handler,
	//	}),
	//)
}
