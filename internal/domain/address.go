package domain

import "time"

type (
	Address struct {
		ID             int       `json:"id" db:"id"`
		City           string    `json:"city" db:"city"`
		Street         string    `json:"street" db:"street"`
		House          string    `json:"house" db:"house"`
		Floor          *int      `json:"floor" db:"floor"`
		Entrance       *int      `json:"entrance" db:"entrance"`
		AdditionalInfo *string   `json:"additional_info" db:"additional_info"`
		CreatedAt      time.Time `json:"created_at" db:"created_at"`
		UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
	}

	AddressInput struct {
		City           string  `json:"city" validate:"required"`
		Street         string  `json:"street" validate:"required"`
		House          string  `json:"house" validate:"required"`
		Floor          *int    `json:"floor"`
		Entrance       *int    `json:"entrance"`
		AdditionalInfo *string `json:"additional_info"`
	}
)
