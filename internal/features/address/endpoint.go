package address

import (
	"github.com/JohnKucharsky/StoreAPI/internal/domain"
	"github.com/JohnKucharsky/StoreAPI/internal/shared"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func (h *addressService) Create(c *fiber.Ctx) error {
	var input domain.AddressInput
	if err := shared.BindBody(c, &input); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared.ErrorRes(err.Error()))
	}

	createdEntity, err := h.repository.Create(c.Context(), input)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusCreated).JSON(shared.SuccessRes(createdEntity))
}

func (h *addressService) Get(c *fiber.Ctx) error {
	address, err := h.repository.GetMany(c.Context())
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(shared.SuccessRes(address))
}

func (h *addressService) GetOne(c *fiber.Ctx) error {
	id, err := shared.GetID(c)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared.ErrorRes(err.Error()))
	}

	address, err := h.repository.GetOne(c.Context(), id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(shared.SuccessRes(address))
}

func (h *addressService) Update(c *fiber.Ctx) error {
	id, err := shared.GetID(c)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared.ErrorRes(err.Error()))
	}

	var req domain.AddressInput
	if err := shared.BindBody(c, &req); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared.ErrorRes(err.Error()))
	}

	updatedEntity, err := h.repository.Update(c.Context(), req, id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusCreated).JSON(shared.SuccessRes(updatedEntity))
}

func (h *addressService) Delete(c *fiber.Ctx) error {
	id, err := shared.GetID(c)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared.ErrorRes(err.Error()))
	}

	resId, err := h.repository.Delete(c.Context(), id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(shared.SuccessRes(fiber.Map{"id": resId}))
}
