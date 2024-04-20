package repository

import (
	"context"
	"log"
	"serviceGateway/internal/config"
	"time"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func InitRedisClient() {
	ctx := context.Background()
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.ServiceGatewayConfigs.Redis.Host + ":" + config.ServiceGatewayConfigs.Redis.Port,
		DB:       config.ServiceGatewayConfigs.Redis.DB,
		PoolSize: 10,
	})

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Panicf("failed to connect to redis, error: %+v", err.Error())
	}
}

func SetOnRedis(key string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()
	return RedisClient.Set(ctx, key, value, expiration).Err()
}

func GetFromRedis(key string) (string, error) {
	ctx := context.Background()
	result, err := RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return result, err
}
