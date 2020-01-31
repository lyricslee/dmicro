package main

import (
	"time"

	"github.com/micro/go-micro/v2"
	ccmd "github.com/micro/go-micro/v2/config/cmd"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/micro/v2/cmd"

	"dmicro/common/log"
	"dmicro/gate/micro/config"
)

func main() {
	config.Init()
	log.Init(config.Logger)

	var opts []micro.Option
	if config.Micro.RegisterTTL > 0 {
		opts = append(opts, micro.RegisterTTL(time.Second*time.Duration(config.Micro.RegisterTTL)))
	}
	if config.Micro.RegisterInterval > 0 {
		opts = append(opts, micro.RegisterInterval(time.Second*time.Duration(config.Micro.RegisterInterval)))
	}
	if config.Micro.Registry != "" {
		r, ok := ccmd.DefaultRegistries[config.Micro.Registry]
		if !ok {
			log.Fatalf("Registry %s not found", config.Micro.Registry)
		}
		reg := r()
		if reg != nil {
			if err := reg.Init(registry.Addrs(config.Micro.RegistryAddress...)); err != nil {
				log.Fatalf("Error configuring registry: %v", err)
			}
			opts = append(opts, micro.Registry(reg))
			if err := ccmd.Init(ccmd.Registry(&reg)); err != nil {
				log.Fatalf("Error configuring registry: %v", err)
			}
		}
	}
	cmd.Init(opts...)
}
