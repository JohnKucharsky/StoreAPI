package domain

type OperationsStore interface {
	GetAssemblyInfoByOrders(m []int) ([]AssemblyInfo, error)
}

type AssemblyProduct struct {
	ProductID       int    `json:"product_id"`
	ProductName     string `json:"product_name"`
	OrderID         int    `json:"order_id"`
	Quantity        int    `json:"quantity"`
	AdditionalShelf *int   `json:"additional_shelf"`
}

type AssemblyInfo struct {
	Name     string            `json:"name"`
	Products []AssemblyProduct `json:"products"`
}

type ShelfProductDB struct {
	ShelfID    int `db:"shelf_id"`
	ProductID  int `db:"product_id"`
	ProductQty int `db:"product_qty"`
}

type OrdersListInput struct {
	OrdersIDs []int `json:"ordersIDs" validate:"required"`
}
