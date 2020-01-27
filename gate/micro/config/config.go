package config

import (
	"github.com/spf13/viper"

	"dmicro/common/config"
)

var (
	Micro  *config.Micro
	Logger *config.Logger

	Tracer *config.Tracer

	Hystrix *config.Hystrix
)

func Init() {
	config.Init()

	Micro = config.GetMicro()
	Logger = config.GetLogger()

	Tracer = &config.Tracer{}
	Tracer.Addr = viper.GetString("tracer.addr")

	Hystrix = &config.Hystrix{}
	Hystrix.Timeout = viper.GetInt("hystrix.timout")
	Hystrix.MaxConcurrent = viper.GetInt("hystrix.max_concurrent")
	Hystrix.RequestVolumeThreshold = viper.GetInt("hystrix.request_volume_threshold")
	Hystrix.SleepWindow = viper.GetInt("hystrix.sleep_window")
	Hystrix.ErrorPercentThreshold = viper.GetInt("hystrix.error_percent_threshold")
}
