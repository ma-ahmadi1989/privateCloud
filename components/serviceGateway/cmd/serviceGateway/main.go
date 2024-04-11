package main

import (
	"serviceGateway/internal/api"
	"serviceGateway/internal/config"
	"serviceGateway/internal/repository"
)

func main() {
	config.LoadConfigs()
	repository.InitRedisClient()

	api.StartAPIServer()

}
