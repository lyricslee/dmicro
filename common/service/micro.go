package service

import (
	"log"
	"time"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/config/cmd"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/transport"

	"dmicro/common/config"
)

type service struct {
	micro.Service
}

func NewService() *service {
	s := &service{Service: micro.NewService()}
	var opts []micro.Option
	if config.GetMicro().ServerName != "" {
		opts = append(opts, micro.Name(config.GetMicro().ServerName))
	}
	if config.GetMicro().ServerVersion != "" {
		opts = append(opts, micro.Version(config.GetMicro().ServerVersion))
	}
	if config.GetMicro().RegisterTTL > 0 {
		opts = append(opts, micro.RegisterTTL(time.Second*time.Duration(config.GetMicro().RegisterTTL)))
	}
	if config.GetMicro().RegisterInterval > 0 {
		opts = append(opts, micro.RegisterInterval(time.Second*time.Duration(config.GetMicro().RegisterInterval)))
	}
	if config.GetMicro().Registry != "" {
		r, ok := cmd.DefaultRegistries[config.GetMicro().Registry]
		if !ok {
			log.Fatalf("Registry %s not found", config.GetMicro().Registry)
		}
		reg := r()
		if reg != nil {
			if err := reg.Init(registry.Addrs(config.GetMicro().RegistryAddress...)); err != nil {
				log.Fatalf("Error configuring registry: %v", err)
			}
			opts = append(opts, micro.Registry(reg))
		}
	}
	if config.GetMicro().Transport != "" {
		b, ok := cmd.DefaultTransports[config.GetMicro().Transport]
		if !ok {
			log.Fatalf("Transport %s not found", config.GetMicro().Transport)
		}
		br := b()
		if br != nil {
			if err := br.Init(transport.Addrs(config.GetMicro().TransportAddress...)); err != nil {
				log.Fatalf("Error configuring transport: %v", err)
			}
			opts = append(opts, micro.Transport(br))
		}
	}
	if config.GetMicro().Broker != "" {
		b, ok := cmd.DefaultBrokers[config.GetMicro().Broker]
		if !ok {
			log.Fatalf("Broker %s not found", config.GetMicro().Broker)
		}
		br := b()
		if br != nil {
			if err := br.Init(broker.Addrs(config.GetMicro().BrokerAddress...)); err != nil {
				log.Fatalf("Error configuring broker: %v", err)
			}
			opts = append(opts, micro.Broker(br))
		}
	}

	s.Init(opts...)

	return s
}
