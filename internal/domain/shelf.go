package domain

import (
	"time"
)

type (
	Shelf struct {
		ID        int       `json:"id" db:"id"`
		Name      string    `json:"name" db:"name"`
		CreatedAt time.Time `json:"created_at" db:"created_at"`
		UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	}

	ShelfInput struct {
		Name string `json:"name" validate:"required"`
	}

	ShelfProductDB struct {
		ShelfID    int `db:"shelf_id"`
		ProductID  int `db:"product_id"`
		ProductQty int `db:"product_qty"`
	}

	ShelfInfo struct {
		Shelf   Shelf
		Product []ProductWithQty
	}
)
