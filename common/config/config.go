package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"

	"dmicro/common/config/env"
)

type Micro struct {
	RegisterTTL      int `json:"register_ttl,omitempty"`
	RegisterInterval int `json:"register_interval,omitempty"`

	Server          string `json:"server,omitempty"`
	ServerName      string `json:"server_name,omitempty"`
	ServerVersion   string `json:"server_version,omitempty"`
	ServerId        string `json:"server_id,omitempty"`
	ServerAddress   string `json:"server_address,omitempty"`
	ServerAdvertise string `json:"server_advertise,omitempty"`

	Registry        string   `json:"registry,omitempty"`
	RegistryAddress []string `json:"registry_address,omitempty"`

	Transport        string   `json:"transport,omitempty"`
	TransportAddress []string `json:"transport_address,omitempty"`

	Broker        string   `json:"broker,omitempty"`
	BrokerAddress []string `json:"broker_address,omitempty"`
}

type Logger struct {
	Level      string `json:"level,omitempty"`
	Filename   string `json:"filename,omitempty"`
	MaxSize    int    `json:"max_size,omitempty"`
	MaxBackups int    `json:"max_backups,omitempty"`
	MaxAge     int    `json:"max_age,omitempty"`
	Compress   bool   `json:"compress,omitempty"`
}

var (
	micro  *Micro
	logger *Logger
)

func GetMicro() *Micro {
	return micro
}

func GetLogger() *Logger {
	return logger
}

type Mysql struct {
	DataSource string `json:"data_source,omitempty"`
	MaxIdle    int    `json:"max_idle,omitempty"`
	MaxOpen    int    `json:"max_open,omitempty"`
}

type Tracer struct {
	Addr string `json:"tracer_addr,omitempty"`
}

type Redis struct {
	Addr     string `json:"addr,omitempty"`
	Password string `json:"password,omitempty"`
	MaxIdle  int    `json:"max_idle,omitempty"`
}

type Hystrix struct {
	Timeout                int `json:"timout,omitempty"`
	MaxConcurrent          int `json:"max_concurrent,omitempty"`
	RequestVolumeThreshold int `json:"request_volume_threshold,omitempty"`
	SleepWindow            int `json:"sleep_window,omitempty"`
	ErrorPercentThreshold  int `json:"error_percent_threshold,omitempty"`
}

func Init() {
	viper.SetConfigName("app")
	viper.AddConfigPath("./conf")
	viper.SetConfigType("yml")
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		log.Fatal(fmt.Sprintf("Read config file fail: %s", err.Error()))
	}
	for k, v := range viper.AllSettings() {
		viper.SetDefault(k, v)
	}
	env.DeployEnv = viper.GetString("deploy_env")
	if env.DeployEnv != "" {
		viper.SetConfigName(env.DeployEnv)
		viper.AddConfigPath("./conf")
		viper.SetConfigType("yaml")
		if err := viper.ReadInConfig(); err == nil {
			log.Println("Using config file:", viper.ConfigFileUsed())
		} else {
			log.Fatal(fmt.Sprintf("Read config file fail: %s", err.Error()))
		}
	}

	micro = &Micro{}
	micro.RegisterTTL = viper.GetInt("micro.register_ttl")
	micro.RegisterInterval = viper.GetInt("micro.register_interval")

	micro.Server = viper.GetString("micro.server")
	micro.ServerName = viper.GetString("micro.server_name")
	micro.ServerVersion = viper.GetString("micro.server_version")
	micro.ServerId = viper.GetString("micro.server_id")
	micro.ServerAddress = viper.GetString("micro.server_address")
	micro.ServerAdvertise = viper.GetString("micro.server_advertise")

	micro.Registry = viper.GetString("micro.registry")
	micro.RegistryAddress = viper.GetStringSlice("micro.registry_address")

	micro.Transport = viper.GetString("micro.transport")
	micro.TransportAddress = viper.GetStringSlice("micro.transport_address")

	micro.Broker = viper.GetString("micro.broker")
	micro.BrokerAddress = viper.GetStringSlice("micro.broker_address")

	logger = &Logger{}
	logger.Level = viper.GetString("logger.level")
	logger.Filename = viper.GetString("logger.filename")
	logger.MaxSize = viper.GetInt("logger.max_size")
	logger.MaxBackups = viper.GetInt("logger.max_backups")
	logger.MaxAge = viper.GetInt("logger.max_age")
	logger.Compress = viper.GetBool("logger.compress")
}
