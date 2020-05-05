package main

import (
	"net/http"
	"strings"

	"github.com/micro/micro/v2/api"
	"github.com/micro/micro/v2/web"

	"dmicro/gate/micro/config"
	"dmicro/gate/micro/plugin/auth"
	"dmicro/gate/micro/plugin/breaker"
	"dmicro/gate/micro/plugin/metrics"
	"dmicro/gate/micro/plugin/trace"
	"dmicro/pkg/tracer"
)

func initPlugin() {
	initMetrics()
	initTrace()
	initBreaker()
	initAuth()
}

// 前缀检查
func checkPrefix(s string, prefixes ...string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}

// 注册认证函数：网关统一认证，API/URL 权限检查等。
// web: web 控制台
// api: api 网关
func initAuth() {
	_ = web.Register(auth.NewPlugin(
		auth.WithSkipperFunc(func(r *http.Request) bool {
			// micro控制台访问放行
			if r.URL.Path == "/" ||
				r.URL.Path == "/favicon.ico" ||
				checkPrefix(r.URL.Path, "/stats", "/client", "/registry", "/terminal") {
				return true
			}
			// 统一消息网关认证由自身负责，所以这里放行
			if strings.HasPrefix(r.URL.Path, "/ws") {
				return true
			}
			// 公共接口放行，即访客可访问
			if checkPrefix(r.URL.Path,
				"/dd/passport/SmsLogin",
				"/dd/passport/Sms",
				"/dd/passport/Login",
				"/dd/passport/OAuthLogin") {
				return true
			}

			return false
		}),
	))

	_ = api.Register(auth.NewPlugin(
		auth.WithSkipperFunc(func(r *http.Request) bool {
			// 公共接口放行，即访客可访问
			if checkPrefix(r.URL.Path,
				"/dd/passport/SmsLogin",
				"/dd/passport/Sms",
				"/dd/passport/Login",
				"/dd/passport/OAuthLogin") {
				return true
			}
			return false
		}),
	))

}

// breaker init
func initBreaker() {
	_ = web.Register(breaker.NewPlugin(
		breaker.WithSkipperFunc(func(r *http.Request) bool {
			// micro控制台访问放行
			if r.URL.Path == "/" ||
				r.URL.Path == "/favicon.ico" ||
				checkPrefix(r.URL.Path, "/stats", "/client", "/registry", "/terminal") {
				return true
			}

			return false
		}),
	))

	_ = api.Register(breaker.NewPlugin(
		breaker.WithSkipperFunc(func(r *http.Request) bool {
			return false
		}),
	))
}

// Init metrics
func initMetrics() {
	_ = web.Register(metrics.NewPlugin())

	_ = api.Register(metrics.NewPlugin())

}

// Init with WithSkipperFunc and WithTracer
func initTrace() {
	webTracer, _ := tracer.Init("go.micro.web", config.Tracer.Addr)
	_ = web.Register(trace.NewPlugin(
		trace.WithTracer(webTracer),
		trace.WithSkipperFunc(func(r *http.Request) bool {
			// micro控制台访问放行
			if r.URL.Path == "/" ||
				r.URL.Path == "/favicon.ico" ||
				checkPrefix(r.URL.Path, "/stats", "/client", "/registry", "/terminal") {
				return true
			}
			// 统一消息网关认证由自身负责，所以这里放行
			if strings.HasPrefix(r.URL.Path, "/ws") {
				return true
			}

			return false
		}),
	))

	apiTracer, _ := tracer.Init("go.micro.api", config.Tracer.Addr)
	_ = api.Register(trace.NewPlugin(
		trace.WithTracer(apiTracer),
	))
}
