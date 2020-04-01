package config

import (
	"github.com/spf13/viper"

	"dmicro/common/config"
)

var (
	Micro  *config.Micro
	Logger *config.Logger

	Tracer *config.Tracer
)

func Init() {
	// 从 yaml 或者配置中心中读取配置文件
	config.Init()

	// 读取配置文件后，初始化话 Micro Logger 等配置对象。
	Micro = config.GetMicro()
	Logger = config.GetLogger()

	Tracer = &config.Tracer{}
	// viper.go 内部定义了一个 var v *Viper 相当于单例模式
	Tracer.Addr = viper.GetString("tracer.addr")
}
