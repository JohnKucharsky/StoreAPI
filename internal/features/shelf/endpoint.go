package shelf

import (
	"github.com/JohnKucharsky/WarehouseAPI/internal/domain"
	"github.com/JohnKucharsky/WarehouseAPI/internal/shared"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type (
	Service interface {
		Create(ctx *fiber.Ctx) error
		GetMany(ctx *fiber.Ctx) error
		GetOne(ctx *fiber.Ctx) error
		Update(ctx *fiber.Ctx) error
		Delete(ctx *fiber.Ctx) error
	}

	service struct {
		repository StoreI
	}
)

func New(store *Store) Service {
	return &service{repository: store}
}

func (h *service) Create(c *fiber.Ctx) error {
	var input domain.ShelfInput
	if err := shared.BindBody(c, &input); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared.ErrorRes(err.Error()))
	}

	one, err := h.repository.Create(c.Context(), input)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusCreated).JSON(shared.SuccessRes(one))
}

func (h *service) GetMany(c *fiber.Ctx) error {
	many, err := h.repository.GetMany(c.Context())
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(shared.SuccessRes(many))
}

func (h *service) GetOne(c *fiber.Ctx) error {
	inputID, err := shared.GetID(c)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared.ErrorRes(err.Error()))
	}

	shelfWithProductInfo, err := h.repository.GetOne(c.Context(), inputID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(shared.SuccessRes(shelfWithProductInfo))
}

func (h *service) Update(c *fiber.Ctx) error {
	inputID, err := shared.GetID(c)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared.ErrorRes(err.Error()))
	}

	var req domain.ShelfInput
	if err := shared.BindBody(c, &req); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared.ErrorRes(err.Error()))
	}

	one, err := h.repository.Update(c.Context(), req, inputID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusCreated).JSON(shared.SuccessRes(one))
}

func (h *service) Delete(c *fiber.Ctx) error {
	inputID, err := shared.GetID(c)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared.ErrorRes(err.Error()))
	}

	ID, err := h.repository.Delete(c.Context(), inputID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(shared.SuccessRes(fiber.Map{"id": ID}))
}
