package auth

import (
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/config/cmd"
	"github.com/micro/go-micro/v2/util/ctx"
	"github.com/micro/micro/v2/plugin"
	"net/http"

	"dmicro/common/log"
	"dmicro/common/util"
	passport "dmicro/srv/passport/api"
)

var (
	passportClient passport.PassportService
)

type auth struct {
	opts Options
}

func newPlugin(opts ...Option) plugin.Plugin {
	options := newOptions(opts...)
	return &auth{
		opts: options,
	}
}

func (*auth) Flags() []cli.Flag {
	return nil
}

func (*auth) Commands() []*cli.Command {
	return nil
}

func (a *auth) Handler() plugin.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Debugf("auth plugins received: %s %s", r.Method, r.RequestURI)

			if a.opts.skipperFunc(r) {
				h.ServeHTTP(w, r)
				return
			}

			cx := ctx.FromRequest(r)

			_, err := passportClient.ValidateToken(cx, &passport.TokenRequest{})
			if err != nil {
				log.Error(err)
				util.WriteError(w, err)
				return
			}
			// 运行到此说明token认证通过
			h.ServeHTTP(w, r)
		})
	}
}

func (*auth) Init(*cli.Context) error {
	passportClient = passport.NewPassportService("go.micro.srv.passport", *cmd.DefaultCmd.Options().Client)
	return nil
}

func (*auth) String() string {
	return "auth"
}

func NewPlugin(opts ...Option) plugin.Plugin {
	return newPlugin(opts...)
}
