package domain

import (
	"time"
)

type Product struct {
	ID         int       `json:"id" db:"id"`
	Name       string    `json:"name" db:"name"`
	Serial     string    `json:"serial" db:"serial"`
	Price      int       `json:"price" db:"price"`
	Model      *string   `json:"model" db:"model"`
	PictureURL string    `json:"picture_url" db:"picture_url"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type ProductInput struct {
	Name       string `json:"name" validate:"required"`
	Serial     string `json:"serial" validate:"required"`
	Price      int    `json:"price" validate:"required"`
	Model      *int   `json:"model"`
	PictureURL string `json:"picture_url" validate:"required"`
}

type ProductWithQty struct {
	Product
	Quantity int `json:"quantity" db:"quantity"`
}
