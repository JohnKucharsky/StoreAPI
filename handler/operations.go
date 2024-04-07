package handler

import (
	"github.com/JohnKucharsky/StoreAPI/domain"
	"github.com/JohnKucharsky/StoreAPI/utils"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func (h *Handler) PutProductsOnShelf(c *fiber.Ctx) error {
	var input domain.PlaceProductInput
	if err := utils.BindBody(c, &input); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.ErrorRes(err.Error()))
	}

	one, err := h.operationsStore.OneShelfToManyProducts(input)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusCreated).JSON(utils.SuccessRes(one))
}

func (h *Handler) GetAssemblyInfo(c *fiber.Ctx) error {
	ids, err := utils.GetIntArrFromOriginalURL(c, "id")
	if err != nil || ids == nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.ErrorRes(err.Error()))
	}

	many, err := h.operationsStore.GetAssemblyInfoByOrders(*ids)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(utils.SuccessRes(many))
}
