package auth

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/config/cmd"
	"github.com/micro/go-micro/v2/util/ctx"
	"github.com/micro/micro/v2/plugin"

	"dmicro/common/log"
	"dmicro/common/util"
	passport "dmicro/srv/passport/api"
)

type auth struct {
	opts           Options
	passportClient passport.PassportService
}

// newPlugin with auth
// Options 就是该插件的参数，目前只有 SkipperFunc 就是处理函数了。
func newPlugin(opts ...Option) plugin.Plugin {
	options := newOptions(opts...)
	return &auth{
		opts: options,
	}
}

// plugin.Plugin require Flags(), Commands(), Handler() and so on....
func (*auth) Flags() []cli.Flag {
	return nil
}

func (*auth) Commands() []*cli.Command {
	return nil
}

// auth handler
// Auth 认证模块，可能用的 Token 或者 Session 认证机制。
func (a *auth) Handler() plugin.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Debugf("auth plugins received: %s %s", r.Method, r.RequestURI)

			if a.opts.skipperFunc(r) {
				h.ServeHTTP(w, r)
				return
			}

			cx := ctx.FromRequest(r)
			var err error

			// 1. Token 认证机制
			if token := strings.Join(r.Header["Token"], ","); token != "" {
				// Token
				log.Debug("AuthToken...")
				rsp, e := a.passportClient.AuthToken(cx, &passport.AuthTokenRequest{})
				log.Debug(rsp)
				if e == nil {
					r.Header.Set("App-Id", strconv.Itoa(int(rsp.Appid)))
					r.Header.Set("Uid", strconv.FormatInt(rsp.Uid, 10))
					r.Header.Set("Plat", strconv.Itoa(int(rsp.Plat)))
				}
				err = e
			} else {
				// 2. SESSION 认证机制，从 Cookie 中读取 SESSION。
				// Cookie
				cookie, _ := r.Cookie("SESSION")
				val, _ := url.QueryUnescape(cookie.Value)
				rsp, e := a.passportClient.AuthCookie(cx, &passport.AuthCookieRequest{Cookie: val})
				log.Debug(rsp)
				if e == nil {
					r.Header.Set("App-Id", strconv.Itoa(int(rsp.Appid)))
					r.Header.Set("Uid", strconv.FormatInt(rsp.Uid, 10))
					r.Header.Set("Plat", strconv.Itoa(int(rsp.Plat)))
				}
				err = e

			}
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

func (a *auth) Init(*cli.Context) error {
	a.passportClient = passport.NewPassportService("go.micro.srv.passport", *cmd.DefaultCmd.Options().Client)
	return nil
}

func (*auth) String() string {
	return "auth"
}

func NewPlugin(opts ...Option) plugin.Plugin {
	return newPlugin(opts...)
}
