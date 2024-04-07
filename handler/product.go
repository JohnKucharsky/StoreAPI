package handler

import (
	"github.com/JohnKucharsky/StoreAPI/domain"
	"github.com/JohnKucharsky/StoreAPI/utils"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func (h *Handler) CreateProduct(c *fiber.Ctx) error {
	var input domain.ProductInput
	if err := utils.BindBody(c, &input); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.ErrorRes(err.Error()))
	}

	product, err := h.productStore.Create(input)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrorRes(err.Error()))
	}

	if product.MainShelfID != nil {
		var assignShelfInput = domain.PlaceProductInput{
			ShelfID:    *product.MainShelfID,
			ProductIDs: []int{product.ID},
		}

		_, err := h.operationsStore.OneShelfToManyProducts(assignShelfInput)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(utils.ErrorRes(err.Error()))
		}
	} else {
		randomShelfID, err := h.shelfStore.GetRandomShelfID()
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(utils.ErrorRes(err.Error()))
		}

		if randomShelfID != nil {
			var assignShelfInput = domain.PlaceProductInput{
				ShelfID:    *randomShelfID,
				ProductIDs: []int{product.ID},
			}

			_, err = h.operationsStore.OneShelfToManyProducts(assignShelfInput)
			if err != nil {
				return c.Status(http.StatusBadRequest).JSON(utils.ErrorRes(err.Error()))
			}
		} else {
			return c.Status(http.StatusBadRequest).JSON(utils.ErrorRes("Unable to place product on shelf, please add shelves"))
		}
	}

	return c.Status(http.StatusCreated).JSON(utils.SuccessRes(product))
}

func (h *Handler) GetProducts(c *fiber.Ctx) error {
	many, err := h.productStore.GetMany()
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(utils.SuccessRes(many))
}

func (h *Handler) GetOneProduct(c *fiber.Ctx) error {
	inputID, err := utils.GetID(c)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.ErrorRes(err.Error()))
	}

	one, err := h.productStore.GetOne(inputID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(utils.SuccessRes(one))
}

func (h *Handler) UpdateProduct(c *fiber.Ctx) error {
	inputID, err := utils.GetID(c)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.ErrorRes(err.Error()))
	}

	var input domain.ProductInput
	if err := utils.BindBody(c, &input); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.ErrorRes(err.Error()))
	}

	product, err := h.productStore.Update(input, inputID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrorRes(err.Error()))
	}

	if product.MainShelfID != nil {
		var assignShelfInput = domain.PlaceProductInput{
			ShelfID:    *product.MainShelfID,
			ProductIDs: []int{product.ID},
		}

		_, err := h.operationsStore.OneShelfToManyProducts(assignShelfInput)

		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(utils.ErrorRes(err.Error()))
		}
	} else {
		randomShelfID, err := h.shelfStore.GetRandomShelfID()
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(utils.ErrorRes(err.Error()))
		}

		if randomShelfID != nil {
			var assignShelfInput = domain.PlaceProductInput{
				ShelfID:    *product.MainShelfID,
				ProductIDs: []int{*randomShelfID},
			}

			_, err = h.operationsStore.OneShelfToManyProducts(assignShelfInput)

			if err != nil {
				return c.Status(http.StatusBadRequest).JSON(utils.ErrorRes(err.Error()))
			}
		} else {
			return c.Status(http.StatusBadRequest).JSON(utils.ErrorRes("Unable to place product on shelf, please add shelves"))
		}
	}

	return c.Status(http.StatusCreated).JSON(utils.SuccessRes(product))
}

func (h *Handler) DeleteProduct(c *fiber.Ctx) error {
	inputID, err := utils.GetID(c)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.ErrorRes(err.Error()))
	}

	ID, err := h.productStore.Delete(inputID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(utils.SuccessRes(fiber.Map{"id": ID}))
}
