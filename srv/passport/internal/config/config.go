package config

import (
	"github.com/spf13/viper"

	"dmicro/common/config"
)

var (
	Micro  *config.Micro
	Logger *config.Logger

	Mysql  *config.Mysql
	Tracer *config.Tracer

	StanBroker *stanBroker
)

type stanBroker struct {
	Addrs       []string `json:"addrs,omitempty"`
	ClusterID   string   `json:"cluster_id,omitempty"`
	DurableName string   `json:"durable_name,omitempty"`
	Queue       string   `json:"queue,omitempty"`
}

func Init() {
	config.Init()

	Micro = config.GetMicro()
	Logger = config.GetLogger()

	Mysql = &config.Mysql{}
	Mysql.DataSource = viper.GetString("mysql.data_source")
	Mysql.MaxOpen = viper.GetInt("mysql.max_open")
	Mysql.MaxIdle = viper.GetInt("mysql.max_idle")

	StanBroker = &stanBroker{}
	StanBroker.Addrs = viper.GetStringSlice("stan.addrs")
	StanBroker.ClusterID = viper.GetString("stan.cluster_id")
	StanBroker.DurableName = viper.GetString("stan.durable_name")
	StanBroker.Queue = viper.GetString("stan.queue")

	Tracer = &config.Tracer{}
	Tracer.Addr = viper.GetString("tracer.addr")
}
