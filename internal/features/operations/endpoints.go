package operations

import (
	"github.com/JohnKucharsky/StoreAPI/internal/domain"
	"github.com/JohnKucharsky/StoreAPI/internal/shared"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type (
	Service interface {
		GetAssemblyInfo(ctx *fiber.Ctx) error
		PlaceProductsOnShelf(ctx *fiber.Ctx) error
		RemoveProductsFromShelf(ctx *fiber.Ctx) error
	}

	service struct {
		repository StoreI
	}
)

func New(store *Store) Service {
	return &service{repository: store}
}

func (h *service) GetAssemblyInfo(c *fiber.Ctx) error {
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

	assemblyInfo, err := h.repository.GetAssemblyInfoByOrders(c.Context(), *ids)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(shared.SuccessRes(assemblyInfo))
}

func (h *service) PlaceProductsOnShelf(c *fiber.Ctx) error {
	var input domain.PlaceRemoveProductInput
	if err := shared.BindBody(c, &input); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared.ErrorRes(err.Error()))
	}

	productsWithQty, err := h.repository.PlaceProductsOnShelf(c.Context(), input)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusCreated).JSON(shared.SuccessRes(productsWithQty))
}

func (h *service) RemoveProductsFromShelf(c *fiber.Ctx) error {
	var input domain.PlaceRemoveProductInput
	if err := shared.BindBody(c, &input); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared.ErrorRes(err.Error()))
	}

	productsWithQty, err := h.repository.RemoveProductsFromShelf(c.Context(), input)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusCreated).JSON(shared.SuccessRes(productsWithQty))
}
