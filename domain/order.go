package domain

import (
	"github.com/google/uuid"
	"time"
)

type OrderStore interface {
	Create(m OrderInput) (*Order, error)
	GetMany() ([]*Order, error)
	GetOne(id int) (*Order, error)
	Update(m OrderInput, id int) (*Order, error)
	Delete(id int) (*int, error)
}

type Order struct {
	ID        int       `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	AddressID int       `json:"address_id" db:"address_id"`
	Total     string    `json:"total" db:"total"`
	Payment   string    `json:"payment" db:"payment"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type OrderInput struct {
	UserID    uuid.UUID `json:"user_id" validate:"required"`
	AddressID int       `json:"address_id" validate:"required"`
	Total     string    `json:"total" validate:"required"`
	Payment   string    `json:"payment" validate:"required"`
}
