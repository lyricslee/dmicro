package breaker

import (
	"errors"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/micro/cli/v2"
	"github.com/micro/micro/v2/plugin"
	"net/http"

	"dmicro/common/log"
	"dmicro/common/util"
	"dmicro/gate/micro/config"
)

type breaker struct {
	opts Options
}

func newPlugin(opts ...Option) plugin.Plugin {
	options := newOptions(opts...)
	return &breaker{
		opts: options,
	}
}

func (*breaker) Flags() []cli.Flag {
	return nil
}

func (*breaker) Commands() []*cli.Command {
	return nil
}

func (b *breaker) Handler() plugin.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Debugf("breaker plugins received: %s %s", r.Method, r.RequestURI)

			if b.opts.skipperFunc(r) {
				h.ServeHTTP(w, r)
				return
			}

			name := r.Method + " " + r.RequestURI
			dw := &util.HttpWriter{ResponseWriter: w}
			err := hystrix.Do(name, func() error {
				defer func() {
					if r := recover(); r != nil {
						log.Errorf("panic recovered: %v", r)
					}
				}()
				h.ServeHTTP(dw, r)

				if dw.Status >= http.StatusBadRequest {
					errmsg := fmt.Sprintf("%d %s", dw.Status, http.StatusText(dw.Status))
					return errors.New(errmsg)
				}
				return nil
			}, func(err error) error {
				log.Error(err)
				util.WriteError(dw, err)
				return err
			})
			if err != nil {
				log.Error(err)
				return
			}
		})
	}
}

func (*breaker) Init(*cli.Context) error {
	if config.Hystrix.Timeout != 0 {
		hystrix.DefaultTimeout = config.Hystrix.Timeout
	}

	if config.Hystrix.MaxConcurrent != 0 {
		hystrix.DefaultMaxConcurrent = config.Hystrix.MaxConcurrent
	}

	if config.Hystrix.RequestVolumeThreshold != 0 {
		hystrix.DefaultVolumeThreshold = config.Hystrix.RequestVolumeThreshold
	}

	if config.Hystrix.SleepWindow != 0 {
		hystrix.DefaultSleepWindow = config.Hystrix.SleepWindow
	}

	if config.Hystrix.ErrorPercentThreshold != 0 {
		hystrix.DefaultErrorPercentThreshold = config.Hystrix.ErrorPercentThreshold
	}

	return nil
}

func (*breaker) String() string {
	return "breaker"
}

func NewPlugin(opts ...Option) plugin.Plugin {
	return newPlugin(opts...)
}
