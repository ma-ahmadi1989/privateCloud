package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"serviceGateway/internal/config"
	"serviceGateway/internal/models"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

func StoreInKafka(event models.GitEvent) error {
	// Setup configuration
	kafkaProducerConfig := sarama.NewConfig()
	kafkaProducerConfig.Producer.Return.Successes = true

	// Create a new producer
	producer, err := sarama.NewSyncProducer(
		[]string{config.ServiceGatewayConfigs.Kafka.ConnectionString},
		kafkaProducerConfig)
	if err != nil {
		log.Printf("Failed to start kafka producer: %v", err)
	}

	// Close the producer after the messages are sent
	defer func() {
		if err := producer.Close(); err != nil {
			log.Printf("Failed to close the kafka producer: %v", err)
		}
	}()

	eventKey, err := GetRepoKey(event)
	if err != nil {
		log.Printf("failed to generate the event key, event: %+v error: %+v", event, err.Error())
	}

	eventJson, err := json.Marshal(event)
	if err != nil {
		log.Printf("failed to convert the event to json, key: %+v error: %+v", event, err.Error())
	}

	// Create a new message
	msg := &sarama.ProducerMessage{
		Topic: config.ServiceGatewayConfigs.Kafka.Topic,
		Key:   sarama.StringEncoder(eventKey),
		Value: sarama.ByteEncoder(eventJson),
	}

	// Send the message
	_, _, err = producer.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}

	return nil
}

func GetRepoKey(event models.GitEvent) (string, error) {
	keySource := fmt.Sprintf("%s%s",
		event.GitUserName,
		event.GitRepoName,
	)

	eventUUID, err := uuid.Parse(keySource)
	if err != nil {
		errorMessage := fmt.Sprintf("faild to generate the event key before storing in event queueu, event: %+v", event)
		return "", errors.New(errorMessage)
	}

	return eventUUID.String(), nil
}