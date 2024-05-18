package domain

import (
	"errors"
	"github.com/google/uuid"
	"github.com/matthewhartstonge/argon2"
	"time"
)

type (
	User struct {
		ID         uuid.UUID `json:"id" db:"id"`
		Name       string    `json:"name" db:"name"`
		LastName   string    `json:"last_name" db:"last_name"`
		MiddleName *string   `json:"middle_name" db:"middle_name"`
		Email      string    `json:"email" db:"email"`
		Password   string    `json:"-" db:"password"`
		CreatedAt  time.Time `json:"created_at" db:"created_at"`
		UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	}

	SignUpInput struct {
		Name       string  `json:"name" validate:"required"`
		LastName   string  `json:"last_name" validate:"required"`
		MiddleName *string `json:"middle_name"`
		Email      string  `json:"email" validate:"required,email"`
		Password   string  `json:"password" validate:"required,min=8"`
	}

	SignInInput struct {
		Email    string `json:"email"  validate:"required"`
		Password string `json:"password"  validate:"required,min=8"`
	}
)

func (r *SignUpInput) HashPassword() error {
	if len(r.Password) == 0 {
		return errors.New("password should not be empty")
	}

	argon := argon2.DefaultConfig()
	encoded, err := argon.HashEncoded([]byte(r.Password))

	if err != nil {
		return err
	}
	r.Password = string(encoded)

	return nil
}

func (r *SignInInput) CheckPassword(plain string) (bool, error) {
	return argon2.VerifyEncoded([]byte(r.Password), []byte(plain))
}
