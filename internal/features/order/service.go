package order

import (
	"github.com/gofiber/fiber/v2"
)

type (
	Service interface {
		Create(ctx *fiber.Ctx) error
		GetMany(ctx *fiber.Ctx) error
		GetOne(ctx *fiber.Ctx) error
		Update(ctx *fiber.Ctx) error
		Delete(ctx *fiber.Ctx) error
	}

	orderService struct {
		repository *OrderStore
	}
)

func New(store *OrderStore) Service {
	return &orderService{repository: store}
}
