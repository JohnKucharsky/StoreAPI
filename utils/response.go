package utils

import (
	"github.com/JohnKucharsky/StoreAPI/domain"
	"github.com/gofiber/fiber/v2"
)

func ErrorRes(errorString string) fiber.Map {
	return fiber.Map{
		"message": errorString,
	}
}

func SuccessRes(data interface{}) fiber.Map {
	return fiber.Map{
		"data": data,
	}
}

func SuccessPaginatedRes(data interface{}, pagination *domain.Pagination) fiber.Map {
	return fiber.Map{
		"data":       data,
		"pagination": pagination,
	}
}
