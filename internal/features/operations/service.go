package operations

import (
	"github.com/gofiber/fiber/v2"
)

type (
	Service interface {
		GetAssemblyInfo(ctx *fiber.Ctx) error
	}

	operationsService struct {
		repository StoreI
	}
)

func New(store *Store) Service {
	return &operationsService{repository: store}
}
