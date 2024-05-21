package domain

type (
	AssemblyProduct struct {
		ProductID     int    `json:"product_id"`
		ProductName   string `json:"product_name"`
		SerialNumber  string `json:"serial"`
		OrderID       int    `json:"order_id"`
		QuantityShelf int    `json:"quantity_shelf"`
		QuantityOrder int    `json:"order_quantity"`
	}

	AssemblyInfo struct {
		Name     string            `json:"shelf_name"`
		Products []AssemblyProduct `json:"products"`
	}

	OrdersListInput struct {
		OrdersIDs []int `json:"orders_id" validate:"required"`
	}

	PlaceRemoveProductInput struct {
		ProductsWithQty []ProductIdQty `json:"products_with_qty" validate:"required,dive"`
		ShelfID         int            `json:"shelf_id" validate:"required"`
	}
)
