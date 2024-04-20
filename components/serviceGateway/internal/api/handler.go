package api

import (
	"errors"
	"fmt"
	"log"
	"serviceGateway/internal/models"
	"serviceGateway/internal/repository"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetEvent(context *fiber.Ctx) error {
	nowTime := time.Now()
	var event models.GitEvent
	if err := context.BodyParser(&event); err != nil {
		responseMessage := map[string]string{
			"error":   err.Error(),
			"message": "failed to parse the body",
		}
		return context.Status(fiber.StatusBadRequest).JSON(responseMessage)
	}
	log.Println("end of request parse", time.Since(nowTime))

	blackListed, err := IsInBlackList(event)
	if err != nil {
		responseMessage := map[string]string{
			"error":   "blacklist check failed",
			"message": "this event will not be proccessed due to internal error",
		}
		return context.Status(fiber.StatusInternalServerError).JSON(responseMessage)
	}
	log.Println("end of black list check", time.Since(nowTime))

	if blackListed {
		responseMessage := map[string]string{
			"error":   "blacklisted",
			"message": "this event will not be proccessed as is in black list",
		}
		return context.Status(fiber.StatusForbidden).JSON(responseMessage)
	}

	go repository.StoreInKafka(event)

	return context.Status(201).JSON("Created")

}

func IsInBlackList(event models.GitEvent) (bool, error) {
	repoKey, err := repository.GetRepoKey(event)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to detect the repo key, event: %+v, error: %+v", event, err.Error())
		log.Println(errorMessage)
		return false, errors.New(errorMessage)
	}

	redisCheckResult, err := repository.GetFromRedis(repoKey)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to check with backlist db, event: %+v, error: %+v", event, err.Error())
		log.Println(errorMessage)
		return false, errors.New(errorMessage)
	}

	if redisCheckResult != "" {
		return true, nil
	} else {
		return false, nil
	}
}
