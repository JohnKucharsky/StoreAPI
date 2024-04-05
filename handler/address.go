package handler

import (
	"github.com/JohnKucharsky/StoreAPI/domain"
	"github.com/JohnKucharsky/StoreAPI/utils"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func (h *Handler) CreateAddress(c *fiber.Ctx) error {
	var input domain.AddressInput
	if err := utils.BindBody(c, &input); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.ErrorRes(err.Error()))
	}

	createdEntity, err := h.addressStore.Create(input)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusCreated).JSON(utils.SuccessRes(createdEntity))
}

func (h *Handler) GetAddresses(c *fiber.Ctx) error {
	address, err := h.addressStore.GetMany()
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(utils.SuccessRes(address))
}

func (h *Handler) GetOneAddress(c *fiber.Ctx) error {
	id, err := utils.GetID(c)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.ErrorRes(err.Error()))
	}

	address, err := h.addressStore.GetOne(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(utils.SuccessRes(address))
}

func (h *Handler) UpdateAddress(c *fiber.Ctx) error {
	id, err := utils.GetID(c)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.ErrorRes(err.Error()))
	}

	var req domain.AddressInput
	if err := utils.BindBody(c, &req); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.ErrorRes(err.Error()))
	}

	updatedEntity, err := h.addressStore.Update(req, id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusCreated).JSON(utils.SuccessRes(updatedEntity))
}

func (h *Handler) DeleteAddress(c *fiber.Ctx) error {
	id, err := utils.GetID(c)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.ErrorRes(err.Error()))
	}

	resId, err := h.addressStore.Delete(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(utils.SuccessRes(fiber.Map{"id": resId}))
}
