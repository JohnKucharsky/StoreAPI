package domain

import "time"

type (
	Product struct {
		ID         int       `json:"id" db:"id"`
		Name       string    `json:"name" db:"name"`
		Serial     string    `json:"serial" db:"serial"`
		Price      int       `json:"price" db:"price"`
		Model      *string   `json:"model" db:"model"`
		PictureURL string    `json:"picture_url" db:"picture_url"`
		CreatedAt  time.Time `json:"created_at" db:"created_at"`
		UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	}

	ProductInput struct {
		Name       string  `json:"name" validate:"required"`
		Serial     string  `json:"serial" validate:"required"`
		Price      int     `json:"price" validate:"required"`
		Model      *string `json:"model"`
		PictureURL string  `json:"picture_url" validate:"required"`
	}

	ProductWithQty struct {
		Product  Product `json:"product"`
		Quantity int     `json:"quantity" db:"product_qty"`
	}

	ProductIdQty struct {
		ProductID int `json:"product_id" validate:"required"`
		Quantity  int `json:"quantity" validate:"required"`
	}

	ProductShort struct {
		ID   int    `db:"id"`
		Name string `db:"name"`
	}

	ProductOrder struct {
		OrderBy   string `json:"order_by"`
		SortOrder string `json:"sort_order"`
	}
)
