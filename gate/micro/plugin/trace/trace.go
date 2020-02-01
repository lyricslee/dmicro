package trace

import (
	"github.com/micro/cli/v2"
	"github.com/micro/micro/v2/plugin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"net/http"

	"dmicro/common/log"
	"dmicro/common/util"
)

type trace struct {
	opts Options
}

func newPlugin(opts ...Option) plugin.Plugin {
	options := newOptions(opts...)
	return &trace{
		opts: options,
	}
}

func (*trace) Flags() []cli.Flag {
	return nil
}

func (*trace) Commands() []*cli.Command {
	return nil
}

func (t *trace) Handler() plugin.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Debugf("trace plugins received: %s %s", r.Method, r.RequestURI)

			if t.opts.skipperFunc(r) {
				h.ServeHTTP(w, r)
				return
			}

			spanCtx, _ := t.opts.tracer.Extract(opentracing.TextMap, opentracing.HTTPHeadersCarrier(r.Header))
			sp := t.opts.tracer.StartSpan(r.URL.Path, opentracing.ChildOf(spanCtx))
			defer sp.Finish()

			if err := sp.Tracer().Inject(
				sp.Context(),
				opentracing.TextMap,
				opentracing.HTTPHeadersCarrier(r.Header)); err != nil {
			}
			sp.SetTag("http.host", r.Host)
			ext.PeerAddress.Set(sp, r.RemoteAddr)
			ext.HTTPUrl.Set(sp, r.URL.Path)
			ext.HTTPMethod.Set(sp, r.Method)

			dw := &util.HttpWriter{ResponseWriter: w}
			h.ServeHTTP(dw, r)
			ext.HTTPStatusCode.Set(sp, uint16(dw.Status))
		})
	}
}

func (t *trace) Init(*cli.Context) error {
	return nil
}

func (*trace) String() string {
	return "trace"
}

func NewPlugin(opts ...Option) plugin.Plugin {
	return newPlugin(opts...)
}
