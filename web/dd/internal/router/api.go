package router

import (
	"dmicro/web/dd/internal/controller"
	"dmicro/web/dd/internal/middleware/tracer"
	"github.com/gin-gonic/gin"
)

func RegisterAPI_v1(r *gin.Engine) {
	// 这里以中间件 middleware 的方式注册了 tracer 追踪
	r.Use(tracer.Tracer())
	// Group 参数
	v1 := r.Group("/")

	// 注册路由
	// TODO: 增加新的路由在此注册！！！！！
	RegisterPassportRouter(v1, &controller.PassportController{})
}
