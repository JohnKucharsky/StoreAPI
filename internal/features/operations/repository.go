package operations

import (
	"errors"
	"github.com/JohnKucharsky/StoreAPI/internal/domain"
	"github.com/JohnKucharsky/StoreAPI/internal/shared"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valyala/fasthttp"
)

type (
	StoreI interface {
		GetAssemblyInfoByOrders(ctx *fasthttp.RequestCtx, m []int) ([]domain.AssemblyInfo, error)
	}

	Store struct {
		db *pgxpool.Pool
	}

	ProductShort struct {
		ID   int    `db:"id"`
		Name string `db:"name"`
	}

	ShelfProductDB struct {
		ShelfID    int `db:"shelf_id"`
		ProductID  int `db:"product_id"`
		ProductQty int `db:"product_qty"`
	}

	OrderProductDB struct {
		ProductID   int    `db:"product_id"`
		ProductName string `db:"product_name"`
		ProductQty  int    `db:"product_qty"`
		Serial      string `db:"serial"`
		OrderID     int    `db:"order_id"`
	}
)

func NewOperationsStore(db *pgxpool.Pool) *Store {
	return &Store{
		db: db,
	}
}

func (store *Store) GetAssemblyInfoByOrders(ctx *fasthttp.RequestCtx, m []int) (
	[]domain.AssemblyInfo,
	error,
) {
	sqlOrderProduct := `select id as product_id, name as product_name, serial, order_id, product_qty 
	from product left join order_product on product.id = order_product.product_id 
	where order_product.order_id = any ($1)`
	orderProduct, err := shared.GetManyRowsInIds[OrderProductDB](ctx, store.db, sqlOrderProduct, m)
	if err != nil {
		return nil, err
	}

	var productIdsToGetShelves []int
	for _, item := range orderProduct {
		productIdsToGetShelves = append(productIdsToGetShelves, item.ProductID)
	}
	if len(productIdsToGetShelves) == 0 {
		return nil, errors.New("can't get shelves, put products on shelves")
	}
	sqlShelfProduct := `select shelf_id, product_id, product_qty from shelf_product where product_id = any ($1)`
	shelfProduct, err := shared.GetManyRowsInIds[ShelfProductDB](ctx, store.db, sqlShelfProduct, productIdsToGetShelves)
	if err != nil {
		return nil, err
	}

	var shelvesIds []int
	for _, shlf := range shelfProduct {
		shelvesIds = append(shelvesIds, shlf.ShelfID)
	}
	sqlShelves := `select * from shelf where id = any ($1)`
	shelves, err := shared.GetManyRowsInIds[domain.Shelf](ctx, store.db, sqlShelves, productIdsToGetShelves)
	if err != nil {
		return nil, err
	}

	productShelfMap := make(map[int]*ShelfProductDB)
	for _, shelfP := range shelfProduct {
		productShelfMap[shelfP.ProductID] = shelfP
	}
	productMap := make(map[int]*OrderProductDB)
	for _, product := range orderProduct {
		productMap[product.ProductID] = product
	}

	var assemblyInfo []domain.AssemblyInfo
	for _, shlf := range shelves {
		var assemblyProducts []domain.AssemblyProduct
		for _, shlfProduct := range shelfProduct {
			if shlf.ID != shlfProduct.ShelfID {
				continue
			}
			assemblyProduct := domain.AssemblyProduct{
				ProductID:     shlfProduct.ProductID,
				ProductName:   productMap[shlfProduct.ProductID].ProductName,
				OrderID:       productMap[shlfProduct.ProductID].OrderID,
				QuantityShelf: productShelfMap[shlfProduct.ProductID].ProductQty,
				QuantityOrder: productMap[shlfProduct.ProductID].ProductQty,
				SerialNumber:  productMap[shlfProduct.ProductID].Serial,
			}
			assemblyProducts = append(assemblyProducts, assemblyProduct)
		}

		var inner = domain.AssemblyInfo{Name: shlf.Name, Products: assemblyProducts}
		assemblyInfo = append(assemblyInfo, inner)
	}

	return assemblyInfo, nil
}
