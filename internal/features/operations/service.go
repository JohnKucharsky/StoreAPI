package operations

import (
	"github.com/gofiber/fiber/v2"
)

type (
	Service interface {
		GetAssemblyInfo(ctx *fiber.Ctx) error
	}

	operationsService struct {
		repository *OperationsStore
	}
)

func New(store *OperationsStore) Service {
	return &operationsService{repository: store}
}
