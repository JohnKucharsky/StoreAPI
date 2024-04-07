package handler

import (
	"github.com/JohnKucharsky/StoreAPI/domain"
	"github.com/JohnKucharsky/StoreAPI/utils"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func (h *Handler) CreateShelf(c *fiber.Ctx) error {
	var input domain.ShelfInput
	if err := utils.BindBody(c, &input); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.ErrorRes(err.Error()))
	}

	one, err := h.shelfStore.Create(input)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusCreated).JSON(utils.SuccessRes(one))
}

func (h *Handler) GetShelves(c *fiber.Ctx) error {
	many, err := h.shelfStore.GetMany()
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(utils.SuccessRes(many))
}

func (h *Handler) GetOneShelf(c *fiber.Ctx) error {
	inputID, err := utils.GetID(c)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.ErrorRes(err.Error()))
	}

	one, err := h.shelfStore.GetOne(inputID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(utils.SuccessRes(one))
}

func (h *Handler) UpdateShelf(c *fiber.Ctx) error {
	inputID, err := utils.GetID(c)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.ErrorRes(err.Error()))
	}

	var req domain.ShelfInput
	if err := utils.BindBody(c, &req); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.ErrorRes(err.Error()))
	}

	one, err := h.shelfStore.Update(req, inputID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusCreated).JSON(utils.SuccessRes(one))
}

func (h *Handler) DeleteShelf(c *fiber.Ctx) error {
	inputID, err := utils.GetID(c)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.ErrorRes(err.Error()))
	}

	ID, err := h.shelfStore.Delete(inputID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(utils.SuccessRes(fiber.Map{"id": ID}))
}
