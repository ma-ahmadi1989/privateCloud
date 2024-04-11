package config

import (
	"os"
	"strconv"
)

type RedisConfig struct {
	Host string
	Port string
	DB   int
}

func LoadRedisConfig() RedisConfig {
	var err error
	redisConfig := RedisConfig{}

	redisConfig.Host = os.Getenv("REDIS_BLACK_LIST_HOST")
	if redisConfig.Host == "" {
		redisConfig.Host = "127.0.0.1"
	}

	redisConfig.Port = os.Getenv("REDIS_BLACK_LIST_PORT")
	if redisConfig.Port == "" {
		redisConfig.Port = "6379"
	}

	redisConfig.DB, err = strconv.Atoi(os.Getenv("REDIS_BLACK_LIST_DB"))
	if err != nil {
		redisConfig.DB = 0
	}

	return redisConfig
}
