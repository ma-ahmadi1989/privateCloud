package api

import (
	"serviceGateway/internal/models"

	"github.com/gofiber/fiber"
)

func GetEvent(context *fiber.Ctx) {
	event := models.GitEvent{}
	if err := context.BodyParser(event); err != nil {
		responseMessage := map[string]string{
			"error":   err.Error(),
			"message": "failed to parse the body",
		}
		context.Status(fiber.StatusBadRequest).JSON(responseMessage)
		return
	}

	if IsInBlackList(event) {
		responseMessage := map[string]string{
			"error":   "blacklisted",
			"message": "this event will not be proccessed as is in black list",
		}
		context.Status(fiber.StatusForbidden).JSON(responseMessage)
		return
	}

	if err := AddToQueue(event); err != nil {
		responseMessage := map[string]string{
			"error":   err.Error(),
			"message": "failed to accept the request",
		}
		context.Status(fiber.ErrInternalServerError.Code).JSON(responseMessage)
		return
	}

	context.Status(201)

}

func IsInBlackList(event models.GitEvent) bool {
	return true
}

func AddToQueue(even models.GitEvent) error {
	return nil
}
