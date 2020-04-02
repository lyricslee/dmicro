package tracer

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"dmicro/common/log"
)

const tracerContextKey = "Tracer-context"

// 使用 Jaeger 来做分布式链路追踪
// Jaeger在Golang中的使用 https://www.lizenghai.com/archives/6130.html
// span 就是 tracer 度量的最小单位，有开始时间和结束时间。
// Tracer.Inject() 发起 trace, 与 Tracer.Extract() 获取 trace, and Carrier. 用于 RPC
// Carrier 就是对应的提取对象 Binary 和 HTTPHeaders 类型。
func Tracer() gin.HandlerFunc {
	return func(c *gin.Context) {
		var sp opentracing.Span

		// 初始化 trace 用 c.Request.Header
		// 失败用 c.Request.URL.Path，成功继续记录。
		// GlobalTracer() 表示这个 tracer 是 RPC 跨进程的，我们通过 Extract() 方法来获取。
		// 这里是 TextMap 类型，实际上 Inject 还有
		spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.TextMap,
			opentracing.HTTPHeadersCarrier(c.Request.Header))
		if err != nil {
			sp = opentracing.GlobalTracer().StartSpan(c.Request.URL.Path)
		} else {
			sp = opentracing.GlobalTracer().StartSpan(c.Request.URL.Path, opentracing.ChildOf(spanCtx))
		}

		// defer 表示 sp.Finish() 在这个 Tracer() 函数执行完成后调用。
		defer sp.Finish()

		// Inject() 注册追踪，类型为 TextMap
		md := make(map[string]string)
		if err := sp.Tracer().Inject(
			sp.Context(),
			opentracing.TextMap,
			opentracing.TextMapCarrier(md)); err != nil {
			log.Error(err)
		}

		// 创建新的 context 并且增加 http header 的 trace 注入
		ctx := context.TODO() // TODO() returns a non-nil, empty Context
		ctx = opentracing.ContextWithSpan(ctx, sp)

		ctx = metadata.NewContext(ctx, md) // a new context
		c.Set(tracerContextKey, ctx) // "Tracer-context"

		// set_tag，在Span中记录请求的附加信息
		// ext 设置追踪的 RemoteAddr 地址，请求地址，请求方法信息等。
		sp.SetTag("http.host", c.Request.Host)
		ext.PeerAddress.Set(sp, c.Request.RemoteAddr)
		ext.HTTPUrl.Set(sp, c.Request.URL.Path)
		ext.HTTPMethod.Set(sp, c.Request.Method)
		c.Next()
		ext.HTTPStatusCode.Set(sp, uint16(c.Writer.Status()))
	}
}

func ContextWithSpan(c *gin.Context) (ctx context.Context, ok bool) {
	v, exist := c.Get(tracerContextKey)
	if !exist {
		ok = false
		ctx = context.TODO()
		return
	}

	ctx, ok = v.(context.Context)
	return
}
