package product

import (
	"github.com/JohnKucharsky/StoreAPI/internal/domain"
	"github.com/JohnKucharsky/StoreAPI/internal/shared"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func (h *productService) Create(c *fiber.Ctx) error {
	var input domain.ProductInput
	if err := shared.BindBody(c, &input); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared.ErrorRes(err.Error()))
	}

	product, err := h.repository.Create(c.Context(), input)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusCreated).JSON(shared.SuccessRes(product))
}

func (h *productService) GetMany(c *fiber.Ctx) error {
	many, err := h.repository.GetMany(c.Context())
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(shared.SuccessRes(many))
}

func (h *productService) GetOne(c *fiber.Ctx) error {
	inputID, err := shared.GetID(c)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared.ErrorRes(err.Error()))
	}

	one, err := h.repository.GetOne(c.Context(), inputID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(shared.SuccessRes(one))
}

func (h *productService) Update(c *fiber.Ctx) error {
	inputID, err := shared.GetID(c)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared.ErrorRes(err.Error()))
	}

	var input domain.ProductInput
	if err := shared.BindBody(c, &input); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared.ErrorRes(err.Error()))
	}

	product, err := h.repository.Update(c.Context(), input, inputID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusCreated).JSON(shared.SuccessRes(product))
}

func (h *productService) Delete(c *fiber.Ctx) error {
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
