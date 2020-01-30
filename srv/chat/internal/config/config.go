package config

import (
	"github.com/spf13/viper"

	"dmicro/common/config"
)

var (
	Micro  *config.Micro
	Logger *config.Logger

	Tracer *config.Tracer

	NatsBroker *natsBroker
)

type natsBroker struct {
	Addrs []string `json:"addrs,omitempty"`
}

func Init() {
	config.Init()

	Micro = config.GetMicro()
	Logger = config.GetLogger()

	Tracer = &config.Tracer{}
	Tracer.Addr = viper.GetString("tracer.addr")

	NatsBroker = &natsBroker{}
	NatsBroker.Addrs = viper.GetStringSlice("nats.addrs")
}
