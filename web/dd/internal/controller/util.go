package controller

import (
	"context"
	"strings"

	"github.com/micro/go-micro/metadata"

	mxcontext "dmicro/pkg/context"
	"dmicro/web/dd/internal/middleware/tracer"
)

func toContext(mctx *mxcontext.DmContext) context.Context {
	ctx, _ := tracer.ContextWithSpan(mctx.Context())

	mda, _ := metadata.FromContext(ctx)
	md := metadata.Copy(mda)

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
