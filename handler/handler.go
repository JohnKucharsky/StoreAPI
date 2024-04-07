package handler

import (
	"github.com/JohnKucharsky/StoreAPI/domain"
)

type Handler struct {
	userStore       domain.AuthStore
	addressStore    domain.AddressStore
	productStore    domain.ProductStore
	orderStore      domain.OrderStore
	shelfStore      domain.ShelfStore
	operationsStore domain.OperationsStore
}

func NewHandler(
	us domain.AuthStore,
	addrStore domain.AddressStore,
	productSt domain.ProductStore,
	orderSt domain.OrderStore,
	shelfSt domain.ShelfStore,
	opsStore domain.OperationsStore,
) *Handler {
	return &Handler{
		userStore:       us,
		addressStore:    addrStore,
		productStore:    productSt,
		orderStore:      orderSt,
		shelfStore:      shelfSt,
		operationsStore: opsStore,
	}
}
