package address

import (
	"github.com/gofiber/fiber/v2"
)

type (
	Service interface {
		Create(ctx *fiber.Ctx) error
		Get(ctx *fiber.Ctx) error
		GetOne(ctx *fiber.Ctx) error
		Update(ctx *fiber.Ctx) error
		Delete(ctx *fiber.Ctx) error
	}

	addressService struct {
		repository *AddressStore
	}
)

func New(store *AddressStore) Service {
	return &addressService{repository: store}
}
