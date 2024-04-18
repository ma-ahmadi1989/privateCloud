package api

import (
	"fmt"
	"log"
	"os"
	"serviceGateway/internal/config"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func StartAPIServer() {

	apiService := fiber.New()
	apiService.Use(logger.New(logger.Config{
		Format:       "[${time}] status:${status} - Latency: ${latency} Method: ${method} Path: ${path}\n",
		TimeFormat:   time.RFC3339Nano,
		TimeZone:     "Local",
		TimeInterval: 500 * time.Millisecond,
		Output:       os.Stdout,
		Done: func(c *fiber.Ctx, logString []byte) {
			if c.Response().StatusCode() != fiber.StatusOK {
				log.Println(string(logString))
			}
		},
	}))
	apiService.Post("/", GetEvent)

	if err := apiService.Listen(config.ServiceGatewayConfigs.API.URI()); err != nil {
		log.Printf("api service failed, error: %+v", err.Error())
	}
}

func PrintResponseTime(c *fiber.Ctx) {
	// Start timer
	start := time.Now()

	// Proceed to the next middleware or handler
	c.Next()

	// Calculate response time
	responseTime := time.Since(start)

	// Log response time
	fmt.Printf("Response Time: %s\n", responseTime)

}
