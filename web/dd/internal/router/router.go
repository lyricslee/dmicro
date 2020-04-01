package router

import (
	"net/http"
	"path"

	"github.com/gin-gonic/gin"

	"dmicro/pkg/context"
)

// HandlerFunc 处理函数
type HandlerFunc func(*context.DmContext)

// Handle registers a new request handle and middleware with the given path and method.
// 设置路由 title
// 注册 gin Group http 请求对应的 handler，Group 的概念就是 /api/v1 这种，
//方便我们批量的 /api/v1 /api/v2 而不用每个接口加 v1
func Handle(g *gin.RouterGroup, httpMethod string, relativePath string, handler HandlerFunc, title string) {
	context.SetRouterTitle(httpMethod, path.Join(g.BasePath(), relativePath), title)
	g.Handle(httpMethod, relativePath, func(c *gin.Context) {
		handler(context.New(c))
	})
}

// GET is a shortcut for router.Handle("GET", path, handle).
func GET(g *gin.RouterGroup, relativePath string, handler HandlerFunc, title string) {
	Handle(g, "GET", relativePath, handler, title)
}

// POST is a shortcut for router.Handle("POST", path, handle).
func POST(g *gin.RouterGroup, relativePath string, handler HandlerFunc, title string) {
	Handle(g, "POST", relativePath, handler, title)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle).
func DELETE(g *gin.RouterGroup, relativePath string, handler HandlerFunc, title string) {
	Handle(g, "DELETE", relativePath, handler, title)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle).
func PATCH(g *gin.RouterGroup, relativePath string, handler HandlerFunc, title string) {
	Handle(g, "PATCH", relativePath, handler, title)
}

// PUT is a shortcut for router.Handle("PUT", path, handle).
func PUT(g *gin.RouterGroup, relativePath string, handler HandlerFunc, title string) {
	Handle(g, "PUT", relativePath, handler, title)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle).
func OPTIONS(g *gin.RouterGroup, relativePath string, handler HandlerFunc, title string) {
	Handle(g, "OPTIONS", relativePath, handler, title)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle).
func HEAD(g *gin.RouterGroup, relativePath string, handler HandlerFunc, title string) {
	Handle(g, "HEAD", relativePath, handler, title)
}

// Any registers a route that matches all the HTTP methods.
// GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, CONNECT, TRACE.
func Any(g *gin.RouterGroup, relativePath string, handler HandlerFunc, title string) {
	Handle(g, "GET", relativePath, handler, title)
	Handle(g, "POST", relativePath, handler, title)
	Handle(g, "PUT", relativePath, handler, title)
	Handle(g, "PATCH", relativePath, handler, title)
	Handle(g, "HEAD", relativePath, handler, title)
	Handle(g, "OPTIONS", relativePath, handler, title)
	Handle(g, "DELETE", relativePath, handler, title)
	Handle(g, "CONNECT", relativePath, handler, title)
	Handle(g, "TRACE", relativePath, handler, title)
}

// Router 路由，返回 http.Handler
func Router() http.Handler {
	gin.SetMode("debug")
	// 初始化 gin Engine 对象
	r := gin.New()

	// http 405 errors, StatusMethodNotAllowed
	r.NoMethod(func(ctx *gin.Context) {
		ctx.String(http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		//context.New(ctx).ResponseError(fmt.Errorf(http.StatusText(http.StatusMethodNotAllowed)))
	})

	// http 404 error, StatusNotFound
	r.NoRoute(func(ctx *gin.Context) {
		ctx.String(http.StatusNotFound, http.StatusText(http.StatusNotFound))
		//context.New(ctx).ResponseError(fmt.Errorf(http.StatusText(http.StatusNotFound)))
	})

	// 我们也可以添加更多的 http 40* 错误处理

	// 注册/api/v1路由
	RegisterAPI_v1(r)

	return r
}
