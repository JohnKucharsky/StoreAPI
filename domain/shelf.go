package domain

import (
	"time"
)

type ShelfStore interface {
	Create(m ShelfInput) (*Shelf, error)
	GetMany() ([]*Shelf, error)
	GetOne(id int) (*Shelf, error)
	Update(m ShelfInput, id int) (*Shelf, error)
	Delete(id int) (*int, error)
}

type Shelf struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Destination string    `json:"destination" db:"destination"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type ShelfInput struct {
	Name        string `json:"name" validate:"required"`
	Destination string `json:"destination" validate:"required"`
}
