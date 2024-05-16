package product

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

	productService struct {
		repository StoreI
	}
)

func New(store *Store) Service {
	return &productService{repository: store}
}
