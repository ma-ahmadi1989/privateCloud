package api

import (
	"log"
	"serviceGateway/internal/config"

	"github.com/gofiber/fiber"
)

func StartAPIServer() {

	apiService := fiber.New()
	apiService.Post("/", GetEvent)

	if err := apiService.Listen(config.ServiceGatewayConfigs.API.URI()); err != nil {
		log.Printf("api service failed, error: %+v", err.Error())
	}
}
