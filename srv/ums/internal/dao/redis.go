package dao

import (
	"sync"

	"github.com/go-redis/redis"

	"dmicro/srv/ums/internal/config"
)

var (
	redisClient     *redis.Client
	onceRedisClient sync.Once
)

func GetClient() *redis.Client {
	onceRedisClient.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     config.Redis.Addr,
			Password: config.Redis.Password,
			DB:       config.Redis.DB,
		})
	})

	return redisClient
}
