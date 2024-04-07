package domain

type OperationsStore interface {
	OneShelfToManyProducts(m PlaceProductInput) (*PlaceProduct, error)
	GetAssemblyInfoByOrders(m []int) ([]*AssemblyInfo, error)
}

type AssemblyInfo struct {
	OrderID     int     `db:"order_id" json:"order_id"`
	ProductID   *int    `db:"product_id" json:"product_id"`
	ProductName *string `db:"product_name" json:"product_name"`
	Quantity    *int    `db:"quantity" json:"quantity"`
	ShelfName   *string `db:"shelf_name" json:"shelf_name"`
}

type ShelfProductDB struct {
	ShelfID   int `db:"shelf_id"`
	ProductID int `db:"product_id"`
}

type PlaceProduct struct {
	ShelfID  int       `json:"shelf_id"`
	Products []Product `json:"products"`
}

type PlaceProductInput struct {
	ShelfID    int   `json:"shelf_id" validate:"required"`
	ProductIDs []int `json:"products_ids" validate:"required,dive"`
}

type ShelfInfo struct {
	Shelf   Shelf            `json:"shelf"`
	Product []ProductWithQty `json:"product"`
}

type OrdersListInput struct {
	OrdersIDs []int `json:"ordersIDs" validate:"required"`
}
