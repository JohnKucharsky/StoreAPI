package domain

import (
	"time"
)

type ProductStore interface {
	Create(m ProductInput) (*Product, error)
	GetMany() ([]*Product, error)
	GetOne(id int) (*Product, error)
	Update(m ProductInput, id int) (*Product, error)
	Delete(id int) (*int, error)
}

type Product struct {
	ID          int       `json:"id" db:"id"`
	MainShelfID *int      `json:"main_shelf_id" db:"main_shelf_id"`
	Name        string    `json:"name" db:"name"`
	Serial      string    `json:"serial" db:"serial"`
	Price       int       `json:"price" db:"price"`
	Model       *string   `json:"model" db:"model"`
	PictureURL  string    `json:"picture_url" db:"picture_url"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type ProductInput struct {
	MainShelfID    *int         `json:"main_shelf_id"`
	Name           string       `json:"name" validate:"required"`
	Serial         string       `json:"serial" validate:"required"`
	Price          int          `json:"price" validate:"required"`
	Model          *int         `json:"model"`
	PictureURL     string       `json:"picture_url" validate:"required"`
	AdditionalInfo *string      `json:"additional_info"`
	Shelves        []ShelfIdQty `json:"shelves" validate:"required,dive"`
}

type ShelfIdQty struct {
	ShelfID  int `json:"shelf_id" validate:"required"`
	Quantity int `json:"quantity" validate:"required"`
}

type ProductWithQty struct {
	Product
	Quantity int `json:"quantity" db:"quantity"`
}
