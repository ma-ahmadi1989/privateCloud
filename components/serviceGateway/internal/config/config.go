package config

import (
	"log"

	"github.com/joho/godotenv"
)

type ServiceGatewayConfig struct {
	Redis RedisConfig
	API   API
	Kafka KafkaConfig
}

var ServiceGatewayConfigs *ServiceGatewayConfig

func init() {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Printf("failed to load configs from config file, default values will be used. error: %+v", err.Error())
	}

	ServiceGatewayConfigs = &ServiceGatewayConfig{
		Redis: LoadRedisConfig(),
		API:   LoadAPIConfig(),
		Kafka: LoadKafkaConfig(),
	}
}

func LoadConfigs() *ServiceGatewayConfig {
	return ServiceGatewayConfigs
}
