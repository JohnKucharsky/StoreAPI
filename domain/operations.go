package domain

type OperationsStore interface {
	OneShelfToManyProducts(m PlaceProductInput) (*PlaceProduct, error)
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
