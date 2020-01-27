package router

import (
	"dmicro/web/dd/internal/controller"
	"dmicro/web/dd/internal/middleware/tracer"
	"github.com/gin-gonic/gin"
)

func RegisterAPI_v1(r *gin.Engine) {
	r.Use(tracer.Tracer())
	v1 := r.Group("/")

	// 注册路由
	// 增加新的路由在此注册！！！！！
	RegisterPassportRouter(v1, &controller.PassportController{})
}
