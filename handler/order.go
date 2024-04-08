package handler

import (
	"github.com/JohnKucharsky/StoreAPI/domain"
	"github.com/JohnKucharsky/StoreAPI/utils"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func (h *Handler) CreateOrder(c *fiber.Ctx) error {
	var input domain.OrderInput
	if err := utils.BindBody(c, &input); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.ErrorRes(err.Error()))
	}

	one, err := h.orderStore.Create(input)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusCreated).JSON(utils.SuccessRes(one))
}

func (h *Handler) GetOrders(c *fiber.Ctx) error {
	many, err := h.orderStore.GetMany()
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(utils.SuccessRes(many))
}

func (h *Handler) GetOneOrder(c *fiber.Ctx) error {
	inputID, err := utils.GetID(c)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.ErrorRes(err.Error()))
	}

	orderDB, err := h.orderStore.GetOne(inputID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrorRes(err.Error()))
	}

	address, err := h.addressStore.GetOne(orderDB.AddressID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrorRes(err.Error()))
	}

	products, err := h.orderStore.GetProductsForOrder(orderDB.ID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrorRes(err.Error()))
	}

	response := domain.OrderDbToOrder(orderDB, address, products)

	return c.Status(http.StatusOK).JSON(utils.SuccessRes(response))
}

func (h *Handler) UpdateOrder(c *fiber.Ctx) error {
	inputID, err := utils.GetID(c)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.ErrorRes(err.Error()))
	}

	var input domain.OrderInput
	if err := utils.BindBody(c, &input); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.ErrorRes(err.Error()))
	}

	one, err := h.orderStore.Update(input, inputID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusCreated).JSON(utils.SuccessRes(one))
}

func (h *Handler) DeleteOrder(c *fiber.Ctx) error {
	inputID, err := utils.GetID(c)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.ErrorRes(err.Error()))
	}

	ID, err := h.orderStore.Delete(inputID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(utils.SuccessRes(fiber.Map{"id": ID}))
}
