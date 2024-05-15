package operations

import (
	"errors"
	"github.com/JohnKucharsky/StoreAPI/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valyala/fasthttp"
)

type (
	OperationsStoreI interface {
		GetAssemblyInfoByOrders(ctx *fasthttp.RequestCtx, m []int) ([]domain.AssemblyInfo, error)
	}

	OperationsStore struct {
		db *pgxpool.Pool
	}
)

func NewOperationsStore(db *pgxpool.Pool) *OperationsStore {
	return &OperationsStore{
		db: db,
	}
}

func (store *OperationsStore) GetAssemblyInfoByOrders(ctx *fasthttp.RequestCtx, m []int) (
	[]domain.AssemblyInfo,
	error,
) {
	if len(m) == 0 {
		return nil, errors.New("you have to provide array of orders to get info")
	}

	// get order_product in ids
	orderProductRows, err := store.db.Query(
		ctx, `select product_id,order_id,product_qty from order_product where order_id = any ($1);`,
		m,
	)
	if err != nil {
		return nil, err
	}

	orderProduct, err := pgx.CollectRows(
		orderProductRows, pgx.RowToStructByName[domain.OrderProductDB],
	)
	if err != nil {
		return nil, err
	}
	// get order_product in ids end

	var idsToGetProducts []int
	for _, orderProductDB := range orderProduct {
		idsToGetProducts = append(idsToGetProducts, orderProductDB.ProductID)
	}
	if len(idsToGetProducts) == 0 {
		return nil, errors.New("you have to buy some products")
	}

	// get product in ids
	productsRows, err := store.db.Query(
		ctx, `select id, main_shelf_id, name from product where id = any ($1);`,
		idsToGetProducts,
	)
	if err != nil {
		return nil, err
	}

	type ProductShort struct {
		ID          int    `db:"id"`
		MainShelfID *int   `db:"main_shelf_id"`
		Name        string `db:"name"`
	}
	productRes, err := pgx.CollectRows(
		productsRows, pgx.RowToStructByName[ProductShort],
	)
	if err != nil {
		return nil, err
	}
	// get product in ids end

	var shelvesIds []int
	var productIdsToGetShelves []int

	for _, item := range productRes {
		productIdsToGetShelves = append(productIdsToGetShelves, item.ID)
		if item.MainShelfID != nil {
			shelvesIds = append(shelvesIds, *item.MainShelfID)
		}
	}

	if len(shelvesIds) == 0 && len(productIdsToGetShelves) == 0 {
		return nil, errors.New("can't get shelves, put products on shelves")
	}

	// get shelf_product in ids
	shelfProductRows, err := store.db.Query(
		ctx, `select shelf_id, product_id, product_qty from shelf_product where product_id = any ($1);`,
		productIdsToGetShelves,
	)
	if err != nil {
		return nil, err
	}

	shelfProduct, err := pgx.CollectRows(
		shelfProductRows, pgx.RowToStructByName[domain.ShelfProductDB],
	)
	if err != nil {
		return nil, err
	}
	// get shelf_product in ids end

	for _, shlf := range shelfProduct {
		shelvesIds = append(shelvesIds, shlf.ShelfID)
	}

	// get shelves in ids
	shelvesRows, err := store.db.Query(
		ctx, `select * from shelf where id = any ($1);`,
		shelvesIds,
	)
	if err != nil {
		return nil, err
	}

	shelves, err := pgx.CollectRows(
		shelvesRows, pgx.RowToStructByName[domain.Shelf],
	)
	if err != nil {
		return nil, err
	}
	// get shelves in ids end

	// gather all together
	productOrderMap := make(map[int]int)
	for _, orderProduct := range orderProduct {
		productOrderMap[orderProduct.ProductID] = orderProduct.OrderID
	}
	productShelfMap := make(map[int]*domain.ShelfProductDB)
	for _, shelfP := range shelfProduct {
		productShelfMap[shelfP.ProductID] = &shelfP
	}
	productMap := make(map[int]ProductShort)
	for _, product := range productRes {
		productMap[product.ID] = product
	}

	var assemblyInfo []domain.AssemblyInfo

	for _, shlf := range shelves {
		var assemblyProducts []domain.AssemblyProduct

		for _, product := range shelfProduct {
			if shlf.ID != product.ShelfID {
				continue
			}

			var orderID = productOrderMap[product.ProductID]
			var additionalShelf *int
			if productMap[product.ProductID].MainShelfID != nil {
				additionalShelf = &productShelfMap[product.ProductID].ShelfID
			}

			assemblyProduct := domain.AssemblyProduct{
				ProductID:       product.ProductID,
				ProductName:     productMap[product.ProductID].Name,
				OrderID:         orderID,
				Quantity:        productShelfMap[product.ProductID].ProductQty,
				AdditionalShelf: additionalShelf,
			}

			assemblyProducts = append(assemblyProducts, assemblyProduct)
		}

		var innerRes = domain.AssemblyInfo{Name: shlf.Name, Products: assemblyProducts}
		assemblyInfo = append(assemblyInfo, innerRes)
	}

	return assemblyInfo, nil
}
