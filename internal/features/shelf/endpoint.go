package shelf

import (
	"github.com/JohnKucharsky/StoreAPI/internal/domain"
	shared2 "github.com/JohnKucharsky/StoreAPI/internal/shared"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func (h *shelfService) Create(c *fiber.Ctx) error {
	var input domain.ShelfInput
	if err := shared2.BindBody(c, &input); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared2.ErrorRes(err.Error()))
	}

	one, err := h.repository.Create(c.Context(), input)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared2.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusCreated).JSON(shared2.SuccessRes(one))
}

func (h *shelfService) GetMany(c *fiber.Ctx) error {
	many, err := h.repository.GetMany(c.Context())
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared2.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(shared2.SuccessRes(many))
}

func (h *shelfService) GetOne(c *fiber.Ctx) error {
	inputID, err := shared2.GetID(c)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared2.ErrorRes(err.Error()))
	}

	one, err := h.repository.GetOne(c.Context(), inputID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared2.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(shared2.SuccessRes(one))
}

func (h *shelfService) Update(c *fiber.Ctx) error {
	inputID, err := shared2.GetID(c)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared2.ErrorRes(err.Error()))
	}

	var req domain.ShelfInput
	if err := shared2.BindBody(c, &req); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared2.ErrorRes(err.Error()))
	}

	one, err := h.repository.Update(c.Context(), req, inputID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared2.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusCreated).JSON(shared2.SuccessRes(one))
}

func (h *shelfService) Delete(c *fiber.Ctx) error {
	inputID, err := shared2.GetID(c)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared2.ErrorRes(err.Error()))
	}

	ID, err := h.repository.Delete(c.Context(), inputID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared2.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(shared2.SuccessRes(fiber.Map{"id": ID}))
}
