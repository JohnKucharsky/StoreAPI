package product

import (
	"github.com/JohnKucharsky/StoreAPI/internal/domain"
	"github.com/JohnKucharsky/StoreAPI/internal/shared"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"net/http"
)

type (
	Service interface {
		Create(ctx *fiber.Ctx) error
		GetMany(ctx *fiber.Ctx) error
		GetOne(ctx *fiber.Ctx) error
		Update(ctx *fiber.Ctx) error
		Delete(ctx *fiber.Ctx) error
	}

	service struct {
		repository StoreI
	}
)

func New(store *Store) Service {
	return &service{repository: store}
}

func (h *service) Create(c *fiber.Ctx) error {
	var input domain.ProductInput
	if err := shared.BindBody(c, &input); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared.ErrorRes(err.Error()))
	}

	one, err := h.repository.Create(c.Context(), input)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusCreated).JSON(shared.SuccessRes(one))
}

func (h *service) GetMany(c *fiber.Ctx) error {
	pp, err := shared.GetPaginationParams(c)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared.ErrorRes(err.Error()))
	}

	var orderBy = c.Query("order_by")
	if orderBy != "" && !lo.Contains([]string{"name", "serial", "updated_at"}, orderBy) {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared.ErrorRes("wrong term for order_by"))
	}
	var sortOrder = c.Query("sort_order")
	if sortOrder != "" && !lo.Contains([]string{"asc", "desc"}, sortOrder) {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared.ErrorRes("wrong term for sort_order"))
	}

	many, total, err := h.repository.GetMany(c.Context(), pp, shared.GetOrderString(orderBy, sortOrder))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}

	var pagination = shared.Pagination{
		Total:  total,
		Limit:  pp.Limit,
		Offset: pp.Offset,
	}
	return c.Status(http.StatusOK).JSON(shared.SuccessPaginatedRes(many, &pagination))
}

func (h *service) GetOne(c *fiber.Ctx) error {
	inputID, err := shared.GetID(c)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared.ErrorRes(err.Error()))
	}

	one, err := h.repository.GetOne(c.Context(), inputID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(shared.SuccessRes(one))
}

func (h *service) Update(c *fiber.Ctx) error {
	inputID, err := shared.GetID(c)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared.ErrorRes(err.Error()))
	}

	var input domain.ProductInput
	if err := shared.BindBody(c, &input); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared.ErrorRes(err.Error()))
	}

	one, err := h.repository.Update(c.Context(), input, inputID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusCreated).JSON(shared.SuccessRes(one))
}

func (h *service) Delete(c *fiber.Ctx) error {
	inputID, err := shared.GetID(c)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(shared.ErrorRes(err.Error()))
	}

	ID, err := h.repository.Delete(c.Context(), inputID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(shared.ErrorRes(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(shared.SuccessRes(fiber.Map{"id": ID}))
}
