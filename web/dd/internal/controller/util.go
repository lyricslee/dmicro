package controller

import (
	"context"
	"strings"

	"github.com/micro/go-micro/v2/metadata"

	mxcontext "dmicro/pkg/context"
	"dmicro/web/dd/internal/middleware/tracer"
)

// DmContext 转换成 golang.context
// tracer 追踪 context 调用
// context 包是用来解决，资源多次调用释放的问题。
func toContext(mctx *mxcontext.DmContext) context.Context {
	ctx, _ := tracer.ContextWithSpan(mctx.Context())

	mda, _ := metadata.FromContext(ctx)
	md := metadata.Copy(mda)

	// 删除掉 trace_id
	delete(mctx.Context().Request.Header, "Uber-Trace-Id")
	// set headers
	for k, v := range mctx.Context().Request.Header {
		if _, ok := md[k]; !ok {
			md[k] = strings.Join(v, ",")
		}
	}
	ctx = metadata.NewContext(ctx, md)
	return ctx
}
