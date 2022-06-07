package redis

import (
	config "github.com/titrxw/go-framework/src/Core/Config"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type RedisFactory struct {
	channelMap map[string]*redis.Client
}

func NewRedisFactory() *RedisFactory {
	return &RedisFactory{
		channelMap: make(map[string]*redis.Client),
	}
}

func (this *RedisFactory) Channel(channel string) *redis.Client {
	redis, exists := this.channelMap[channel]
	if !exists {
		panic("redis channel " + channel + " not exists")
	}

	return redis
}

func (this *RedisFactory) RegisterRedis(redisConfig config.Redis) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     redisConfig.Host + ":" + strconv.Itoa(redisConfig.Port),
		Username: redisConfig.Username,
		Password: redisConfig.Password,
		DB:       redisConfig.Db,
		PoolSize: redisConfig.PoolSize,
	})
}

func (this *RedisFactory) Register(maps map[string]config.Redis) {
	for key, value := range maps {
		this.channelMap[key] = this.RegisterRedis(value)
	}
}
