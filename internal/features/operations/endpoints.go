package operations

import (
	"github.com/JohnKucharsky/StoreAPI/internal/shared"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func (h *operationsService) GetAssemblyInfo(c *fiber.Ctx) error {
	ids, err := shared.GetIntArrFromOriginalURL(c, "id")
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared.ErrorRes(err.Error()))

	}
	if ids == nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared.ErrorRes("can't get id's from originalURL"))
	}
	if len(*ids) == 0 {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes("you have to provide array of orders to get info"))
	}

	many, err := h.repository.GetAssemblyInfoByOrders(c.Context(), *ids)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(shared.SuccessRes(many))
}
