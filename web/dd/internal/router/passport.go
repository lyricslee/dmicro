package router

import (
	//"github.com/afex/breaker-go/breaker"
	"dmicro/web/dd/internal/controller"
	"github.com/gin-gonic/gin"
)

func RegisterPassportRouter(g *gin.RouterGroup, c *controller.PassportController) {
	POST(g, "/passport/login", c.Login, "用户登入")
	POST(g, "/passport/sms", c.Sms, "获取验证码")
	POST(g, "/passport/smslogin", c.SmsLogin, "短信验证码登录")
	POST(g, "/passport/oauthlogin", c.OauthLogin, "第三方账号登录")
	POST(g, "/passport/setpwd", c.SetPwd, "设置密码")
}
