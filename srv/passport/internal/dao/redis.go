package dao

import (
	"sync"

	"github.com/go-redis/redis"

	"dmicro/srv/passport/internal/config"
)

var (
	redisClient     *redis.Client
	onceRedisClient sync.Once
)

// redis 连接
func GetRedisClient() *redis.Client {
	onceRedisClient.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     config.Redis.Addr,
			Password: config.Redis.Password,
			DB:       config.Redis.DB,
		})
	})

	return redisClient
}
