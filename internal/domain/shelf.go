package domain

import (
	"time"
)

type Shelf struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type ShelfInput struct {
	Name        string `json:"name" validate:"required"`
	Destination string `json:"destination" validate:"required"`
}
