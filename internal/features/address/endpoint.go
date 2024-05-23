package address

import (
	"github.com/JohnKucharsky/WarehouseAPI/internal/domain"
	"github.com/JohnKucharsky/WarehouseAPI/internal/shared"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type (
	Service interface {
		Create(ctx *fiber.Ctx) error
		Get(ctx *fiber.Ctx) error
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
	var input domain.AddressInput
	if err := shared.BindBody(c, &input); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared.ErrorRes(err.Error()))
	}

	one, err := h.repository.Create(c.Context(), input)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusCreated).JSON(shared.SuccessRes(one))
}

func (h *service) Get(c *fiber.Ctx) error {
	address, err := h.repository.GetMany(c.Context())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(shared.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(shared.SuccessRes(address))
}

func (h *service) GetOne(c *fiber.Ctx) error {
	id, err := shared.GetID(c)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared.ErrorRes(err.Error()))
	}

	address, err := h.repository.GetOne(c.Context(), id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(shared.SuccessRes(address))
}

func (h *service) Update(c *fiber.Ctx) error {
	id, err := shared.GetID(c)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared.ErrorRes(err.Error()))
	}

	var req domain.AddressInput
	if err := shared.BindBody(c, &req); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared.ErrorRes(err.Error()))
	}

	updatedEntity, err := h.repository.Update(c.Context(), req, id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusCreated).JSON(shared.SuccessRes(updatedEntity))
}

func (h *service) Delete(c *fiber.Ctx) error {
	id, err := shared.GetID(c)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared.ErrorRes(err.Error()))
	}

	ID, err := h.repository.Delete(c.Context(), id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(shared.SuccessRes(fiber.Map{"id": ID}))
}
