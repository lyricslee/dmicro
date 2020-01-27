package auth

import (
	"net/http"
	"strings"

	"github.com/micro/cli"
	"github.com/micro/go-micro/config/cmd"
	"github.com/micro/go-micro/util/ctx"
	"github.com/micro/micro/plugin"

	"dmicro/common/log"
	"dmicro/common/util"
	passport "dmicro/srv/passport/api"
)

var (
	passportClient passport.PassportService
)

type auth struct {
}

func (*auth) Flags() []cli.Flag {
	return nil
}

func (*auth) Commands() []cli.Command {
	return nil
}

func (*auth) Handler() plugin.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Debugf("auth plugins received: %s %s", r.Method, r.RequestURI)
			if r.URL.Path == "/" ||
				r.URL.Path == "/stats" ||
				r.URL.Path == "/client" ||
				r.URL.Path == "/registry" ||
				r.URL.Path == "/terminal" ||
				strings.HasPrefix(r.URL.Path, "/dd/passport/smslogin") ||
				strings.HasPrefix(r.URL.Path, "/dd/passport/sms") ||
				strings.HasPrefix(r.URL.Path, "/dd/passport/login") ||
				strings.HasPrefix(r.URL.Path, "/dd/passport/oauthlogin") {
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
	passportClient = passport.NewPassportService("", *cmd.DefaultCmd.Options().Client)
	return nil
}

func (*auth) String() string {
	return "auth"
}

func NewPlugin() plugin.Plugin {
	return new(auth)
}
