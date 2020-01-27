package tracer

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/metadata"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"dmicro/common/log"
)

const tracerContrextKey = "Tracer-context"

func Tracer() gin.HandlerFunc {
	return func(c *gin.Context) {
		var sp opentracing.Span

		spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.TextMap, opentracing.HTTPHeadersCarrier(c.Request.Header))
		if err != nil {
			sp = opentracing.GlobalTracer().StartSpan(c.Request.URL.Path)
		} else {
			sp = opentracing.GlobalTracer().StartSpan(c.Request.URL.Path, opentracing.ChildOf(spanCtx))
		}

		defer sp.Finish()
		md := make(map[string]string)
		if err := sp.Tracer().Inject(
			sp.Context(),
			opentracing.TextMap,
			opentracing.TextMapCarrier(md)); err != nil {
			log.Error(err)
		}

		ctx := context.TODO()
		ctx = opentracing.ContextWithSpan(ctx, sp)

		ctx = metadata.NewContext(ctx, md)
		c.Set(tracerContrextKey, ctx)

		// Jaeger在Golang中的使用
		// https://www.lizenghai.com/archives/6130.html
		sp.SetTag("http.host", c.Request.Host)
		ext.PeerAddress.Set(sp, c.Request.RemoteAddr)
		ext.HTTPUrl.Set(sp, c.Request.URL.Path)
		ext.HTTPMethod.Set(sp, c.Request.Method)
		c.Next()
		ext.HTTPStatusCode.Set(sp, uint16(c.Writer.Status()))
	}
}

func ContextWithSpan(c *gin.Context) (ctx context.Context, ok bool) {
	v, exist := c.Get(tracerContrextKey)
	if !exist {
		ok = false
		ctx = context.TODO()
		return
	}

	ctx, ok = v.(context.Context)
	return
}
