package order

import (
	"github.com/JohnKucharsky/StoreAPI/internal/domain"
	"github.com/JohnKucharsky/StoreAPI/internal/shared"
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
	var input domain.OrderInput
	if err := shared.BindBody(c, &input); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared.ErrorRes(err.Error()))
	}

	user := c.Locals("user").(*domain.User)
	orderDB, err := h.repository.Create(c.Context(), input, user.ID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}

	address, err := h.repository.GetAddress(c.Context(), orderDB.AddressID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}

	products, err := h.repository.GetProductsForOrder(c.Context(), orderDB.ID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}

	order := domain.OrderDbToOrder(orderDB, address, products)

	return c.Status(http.StatusOK).JSON(shared.SuccessRes(order))
}

func (h *service) GetMany(c *fiber.Ctx) error {
	user := c.Locals("user").(*domain.User)

	many, err := h.repository.GetMany(c.Context(), user.ID)
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

	user := c.Locals("user").(*domain.User)
	orderDB, err := h.repository.GetOne(c.Context(), inputID, user.ID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}

	address, err := h.repository.GetAddress(c.Context(), orderDB.AddressID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}

	products, err := h.repository.GetProductsForOrder(c.Context(), orderDB.ID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}

	order := domain.OrderDbToOrder(orderDB, address, products)

	return c.Status(http.StatusOK).JSON(shared.SuccessRes(order))
}

func (h *service) Update(c *fiber.Ctx) error {
	inputID, err := shared.GetID(c)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared.ErrorRes(err.Error()))
	}

	var input domain.OrderInput
	if err := shared.BindBody(c, &input); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared.ErrorRes(err.Error()))
	}

	orderDB, err := h.repository.Update(c.Context(), input, inputID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}

	address, err := h.repository.GetAddress(c.Context(), orderDB.AddressID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}

	products, err := h.repository.GetProductsForOrder(c.Context(), orderDB.ID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}

	order := domain.OrderDbToOrder(orderDB, address, products)

	return c.Status(http.StatusOK).JSON(shared.SuccessRes(order))
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
