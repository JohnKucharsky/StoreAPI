package domain

type AssemblyProduct struct {
	ProductID     int    `json:"product_id"`
	ProductName   string `json:"product_name"`
	SerialNumber  string `json:"serial"`
	OrderID       int    `json:"order_id"`
	QuantityShelf int    `json:"quantity_shelf"`
	QuantityOrder int    `json:"quantity_order"`
}

type AssemblyInfo struct {
	Name     string            `json:"name"`
	Products []AssemblyProduct `json:"products"`
}

type OrdersListInput struct {
	OrdersIDs []int `json:"ordersIDs" validate:"required"`
}
