package domain

import (
	"time"
)

type OrderStore interface {
	Create(m OrderInput) (*OrderShort, error)
	GetMany() ([]*OrderShort, error)
	GetOne(id int) (*OrderShort, error)
	Update(m OrderInput, id int) (*OrderShort, error)
	Delete(id int) (*int, error)
	GetProductsForOrder(id int) ([]*ProductWithQty, error)
}

type Order struct {
	ID        int              `json:"id"`
	Address   Address          `json:"address"`
	Payment   string           `json:"payment"`
	Products  []ProductWithQty `json:"products"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

type OrderShort struct {
	ID        int       `db:"id" json:"id"`
	AddressID int       `db:"address_id" json:"address_id"`
	Payment   string    `db:"payment" json:"payment"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type OrderProductDB struct {
	ProductID  int `db:"product_id"`
	ProductQty int `db:"product_qty"`
	OrderID    int `db:"order_id"`
}

type ProductIdQty struct {
	ProductID int `json:"product_id" validate:"required"`
	Quantity  int `json:"quantity" validate:"required"`
}

type OrderInput struct {
	AddressID int            `json:"address_id" validate:"required"`
	Payment   string         `json:"payment" validate:"required"`
	Products  []ProductIdQty `json:"products" validate:"required,dive"`
}

func OrderDbToOrder(orderDB *OrderShort, address *Address, products []*ProductWithQty) Order {
	var prs []ProductWithQty
	for _, product := range products {
		if product != nil {
			prs = append(prs, *product)
		}
	}

	var addr Address
	if address != nil {
		addr = *address
	}

	return Order{
		ID:        orderDB.ID,
		Address:   addr,
		Payment:   orderDB.Payment,
		Products:  prs,
		CreatedAt: orderDB.CreatedAt,
		UpdatedAt: orderDB.UpdatedAt,
	}
}
