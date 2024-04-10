package handler

import (
	"github.com/JohnKucharsky/StoreAPI/utils"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

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
