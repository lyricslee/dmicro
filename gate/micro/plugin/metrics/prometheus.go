package metrics

import (
	"net/http"

	"github.com/micro/cli/v2"
	"github.com/micro/micro/v2/plugin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// https://www.cnblogs.com/chenqionghe/p/10494868.html
// metrics 监控指标
type metrics struct {
	opts Options
}

func newPlugin(opts ...Option) plugin.Plugin {
	options := newOptions(opts...)
	return &metrics{
		opts: options,
	}
}

func (*metrics) Flags() []cli.Flag {
	return nil
}

func (*metrics) Commands() []*cli.Command {
	return nil
}

func (*metrics) Handler() plugin.Handler {
	ph := promhttp.Handler()
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/metrics" {
				ph.ServeHTTP(w, r)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}

func (*metrics) Init(*cli.Context) error {
	return nil
}

func (*metrics) String() string {
	return "metrics"
}

func NewPlugin(opts ...Option) plugin.Plugin {
	return newPlugin(opts...)
}
