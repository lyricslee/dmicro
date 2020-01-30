package trace

import (
	"net/http"
	"strings"

	"github.com/micro/cli"
	"github.com/micro/micro/plugin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"dmicro/common/log"
	"dmicro/common/util"
	"dmicro/gate/micro/config"
	"dmicro/pkg/tracer"
)

type trace struct {
}

func (*trace) Flags() []cli.Flag {
	return nil
}

func (*trace) Commands() []cli.Command {
	return nil
}

func (*trace) Handler() plugin.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" ||
				r.URL.Path == "/stats" ||
				r.URL.Path == "/client" ||
				r.URL.Path == "/registry" ||
				r.URL.Path == "/terminal" ||
				strings.HasPrefix(r.URL.Path, "/ws") {
				h.ServeHTTP(w, r)
				return
			}

			log.Debugf("trace plugins received: %s %s", r.Method, r.RequestURI)
			spanCtx, _ := opentracing.GlobalTracer().Extract(opentracing.TextMap, opentracing.HTTPHeadersCarrier(r.Header))
			sp := opentracing.GlobalTracer().StartSpan(r.URL.Path, opentracing.ChildOf(spanCtx))
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

func (*trace) Init(*cli.Context) error {
	t, err := tracer.Init("go.micro.web", config.Tracer.Addr)
	if err != nil {
		log.Error(err)
		return nil
	}

	opentracing.SetGlobalTracer(t)
	return nil
}

func (*trace) String() string {
	return "trace"
}

func NewPlugin() plugin.Plugin {
	return new(trace)
}
