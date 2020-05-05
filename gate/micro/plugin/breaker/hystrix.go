package breaker

import (
	"dmicro/common/util"
	"errors"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/micro/cli/v2"
	"github.com/micro/micro/v2/plugin"
	"net/http"

	"dmicro/common/log"
	"dmicro/gate/micro/config"
)

type breaker struct {
	opts Options
}

// 按照 go-micro 要求来注册 Plugin
// Options 就是该插件的参数，目前只有 SkipperFunc 就是处理函数了。
func newPlugin(opts ...Option) plugin.Plugin {
	options := newOptions(opts...)
	return &breaker{
		opts: options,
	}
}

// 向 go-micro 注册 Plugin 的必要函数
func (*breaker) Flags() []cli.Flag {
	return nil
}

// 向 go-micro 注册 Plugin 的必要函数
func (*breaker) Commands() []*cli.Command {
	return nil
}

// handler 函数
func (b *breaker) Handler() plugin.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Debugf("breaker plugins received: %s %s", r.Method, r.RequestURI)

			if b.opts.skipperFunc(r) {
				h.ServeHTTP(w, r)
				return
			}

			name := r.Method + " " + r.RequestURI
			sct := &util.StatusCodeTracker{ResponseWriter: w}
			// hystrix.Do(name string, run runFunc, fallback fallbackFunc)
			// hystrix.Do 熔断发生
			err := hystrix.Do(name, func() error {
				defer func() {
					// recover with panic
					if r := recover(); r != nil {
						log.Errorf("panic recovered: %v", r)
					}
				}()

				h.ServeHTTP(sct.WrappedResponseWriter(), r)

				if sct.Status >= http.StatusBadRequest && sct.Status != 499 {
					errmsg := fmt.Sprintf("%d %s", sct.Status, http.StatusText(sct.Status))
					return errors.New(errmsg)
				}

				return nil
			}, func(err error) error {
				//log.Error(err)
				return err
			})
			if err != nil {
				log.Error(err)
				return
			}
		})
	}
}

// Init breaker with config
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
