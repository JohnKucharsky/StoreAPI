package auth

import (
	"github.com/gofiber/fiber/v2"
)

type (
	Service interface {
		SignUp(ctx *fiber.Ctx) error
		SignIn(ctx *fiber.Ctx) error
		RefreshAccessToken(ctx *fiber.Ctx) error
		DeserializeUser(ctx *fiber.Ctx) error
		GetMe(ctx *fiber.Ctx) error
		LogoutUser(ctx *fiber.Ctx) error
	}

	authService struct {
		repository StoreI
	}
)

func New(store *Store) Service {
	return &authService{repository: store}
}
