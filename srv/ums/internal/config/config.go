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

	Redis *config.Redis
)

type natsBroker struct {
	Addrs []string `json:"addrs,omitempty"`
}

// 配置文件加载
func Init() {
	config.Init()

	Micro = config.GetMicro()
	Logger = config.GetLogger()

	Tracer = &config.Tracer{}
	Tracer.Addr = viper.GetString("tracer.addr")

	NatsBroker = &natsBroker{}
	NatsBroker.Addrs = viper.GetStringSlice("nats.addrs")

	Redis = &config.Redis{}
	Redis.Addr = viper.GetString("redis.addr")
	Redis.Password = viper.GetString("redis.password")
	Redis.DB = viper.GetInt("redis.db")
}
