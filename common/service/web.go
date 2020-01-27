package service

import (
	"log"
	"time"

	"github.com/micro/go-micro/config/cmd"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/web"

	"dmicro/common/config"
)

type webService struct {
	web.Service
}

func NewWebService() *webService {
	s := &webService{Service: web.NewService()}

	var opts []web.Option
	if config.GetMicro().ServerName != "" {
		opts = append(opts, web.Name(config.GetMicro().ServerName))
	}
	if config.GetMicro().ServerVersion != "" {
		opts = append(opts, web.Version(config.GetMicro().ServerVersion))
	}
	if config.GetMicro().RegisterTTL > 0 {
		opts = append(opts, web.RegisterTTL(time.Second*time.Duration(config.GetMicro().RegisterTTL)))
	}
	if config.GetMicro().RegisterInterval > 0 {
		opts = append(opts, web.RegisterInterval(time.Second*time.Duration(config.GetMicro().RegisterInterval)))
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
			opts = append(opts, web.Registry(reg))
		}
	}

	err := s.Init(opts...)
	if err != nil {
		log.Fatal(err)
	}
	return s
}
