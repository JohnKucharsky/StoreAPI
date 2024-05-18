package domain

import (
	"time"
)

type (
	Order struct {
		ID        int              `json:"id"`
		Address   Address          `json:"address"`
		Payment   string           `json:"payment"`
		Products  []ProductWithQty `json:"products"`
		CreatedAt time.Time        `json:"created_at"`
		UpdatedAt time.Time        `json:"updated_at"`
	}

	OrderInput struct {
		AddressID int            `json:"address_id" validate:"required"`
		Payment   string         `json:"payment" validate:"required"`
		Products  []ProductIdQty `json:"products" validate:"required,dive"`
	}

	OrderShort struct {
		ID        int       `db:"id" json:"id"`
		AddressID int       `db:"address_id" json:"address_id"`
		Payment   string    `db:"payment" json:"payment"`
		CreatedAt time.Time `db:"created_at" json:"created_at"`
		UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	}

	OrderProductDB struct {
		ProductID   int    `db:"product_id"`
		ProductName string `db:"product_name"`
		ProductQty  int    `db:"product_qty"`
		Serial      string `db:"serial"`
		OrderID     int    `db:"order_id"`
	}

	OrderProductDBShort struct {
		ProductID  int `db:"product_id"`
		ProductQty int `db:"product_qty"`
		OrderID    int `db:"order_id"`
	}
)

func OrderDbToOrder(orderDB *OrderShort, addrs *Address, products []*ProductWithQty) Order {
	var prs []ProductWithQty
	for _, product := range products {
		if product != nil {
			prs = append(prs, *product)
		}
	}

	var addr Address
	if addrs != nil {
		addr = *addrs
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
