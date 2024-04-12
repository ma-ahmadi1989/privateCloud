package config

import "os"

type KafkaConfig struct {
	Topic            string
	ConnectionString string
}

func LoadKafkaConfig() KafkaConfig {
	kafkaConfig := KafkaConfig{}

	kafkaConfig.Topic = os.Getenv("KAFKA_EVENT_QUEUE_TOPIC_NAME")
	if kafkaConfig.Topic == "" {
		kafkaConfig.Topic = "GIT_EVENTS"
	}

	kafkaConfig.ConnectionString = os.Getenv("KAFKA_EVENT_QUEUE_KAFKA_CONNECTION_STRING")
	if kafkaConfig.ConnectionString == "" {
		kafkaConfig.ConnectionString = "localhost:9092"
	}

	return kafkaConfig
}
