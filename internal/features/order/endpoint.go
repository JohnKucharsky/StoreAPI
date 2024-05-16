package order

import (
	"github.com/JohnKucharsky/StoreAPI/internal/domain"
	"github.com/JohnKucharsky/StoreAPI/internal/shared"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func (h *orderService) Create(c *fiber.Ctx) error {
	var input domain.OrderInput
	if err := shared.BindBody(c, &input); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared.ErrorRes(err.Error()))
	}

	orderDB, err := h.repository.Create(c.Context(), input)
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

	response := domain.OrderDbToOrder(orderDB, address, products)

	return c.Status(http.StatusOK).JSON(shared.SuccessRes(response))
}

func (h *orderService) GetMany(c *fiber.Ctx) error {
	many, err := h.repository.GetMany(c.Context())
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(shared.SuccessRes(many))
}

func (h *orderService) GetOne(c *fiber.Ctx) error {
	inputID, err := shared.GetID(c)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared.ErrorRes(err.Error()))
	}

	orderDB, err := h.repository.GetOne(c.Context(), inputID)
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

	response := domain.OrderDbToOrder(orderDB, address, products)

	return c.Status(http.StatusOK).JSON(shared.SuccessRes(response))
}

func (h *orderService) Update(c *fiber.Ctx) error {
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

	response := domain.OrderDbToOrder(orderDB, address, products)

	return c.Status(http.StatusOK).JSON(shared.SuccessRes(response))
}

func (h *orderService) Delete(c *fiber.Ctx) error {
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
