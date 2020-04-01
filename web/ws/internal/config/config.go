package config

import (
	"log"

	"github.com/spf13/viper"

	"dmicro/common/config"
)

var (
	Micro  *config.Micro
	Logger *config.Logger

	Tracer *config.Tracer

	NatsBroker *natsBroker

	Redis *config.Redis

	GateId int
)

type natsBroker struct {
	Addrs []string `json:"addrs,omitempty"`
}

// 读取配置项基本与 dd 服务一样，多了一个 nats broker 的配置。
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

	GateId = viper.GetInt("gate_id")
	if GateId == 0 {
		log.Fatal("gate id未配置")
	}
}
