package handler

import (
	"github.com/JohnKucharsky/StoreAPI/domain"
)

type Handler struct {
	userStore    domain.AuthStore
	addressStore domain.AddressStore
}

func NewHandler(
	us domain.AuthStore,
	addrStore domain.AddressStore,
) *Handler {
	return &Handler{
		userStore:    us,
		addressStore: addrStore,
	}
}
