package config

import (
	"crypto/rsa"

	"github.com/dgrijalva/jwt-go/test"
	"github.com/spf13/viper"

	"dmicro/common/config"
)

var (
	Micro  *config.Micro
	Logger *config.Logger

	Mysql  *config.Mysql
	Tracer *config.Tracer

	StanBroker *stanBroker

	Redis *config.Redis

	AuthPrivateKey *rsa.PrivateKey
	AuthPublicKey  *rsa.PublicKey
)

type stanBroker struct {
	Addrs       []string `json:"addrs,omitempty"`
	ClusterID   string   `json:"cluster_id,omitempty"`
	DurableName string   `json:"durable_name,omitempty"`
	Queue       string   `json:"queue,omitempty"`
}

// 配置文件的加载
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

	Redis = &config.Redis{}
	Redis.Addr = viper.GetString("redis.addr")
	Redis.Password = viper.GetString("redis.password")
	Redis.DB = viper.GetInt("redis.db")

	AuthPrivateKey = test.LoadRSAPrivateKeyFromDisk("./conf/auth_key")
	AuthPublicKey = test.LoadRSAPublicKeyFromDisk("./conf/auth_key.pub")
}
