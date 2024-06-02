package operations

import (
	"errors"
	"github.com/JohnKucharsky/WarehouseAPI/internal/domain"
	"github.com/JohnKucharsky/WarehouseAPI/internal/shared"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valyala/fasthttp"
)

type (
	StoreI interface {
		GetAssemblyInfoByOrders(ctx *fasthttp.RequestCtx, orders []int) ([]domain.AssemblyInfo, error)
		PlaceProductsOnShelf(ctx *fasthttp.RequestCtx, input domain.PlaceRemoveProductInput) ([]*domain.ProductWithQty, error)
		RemoveProductsFromShelf(ctx *fasthttp.RequestCtx, input domain.PlaceRemoveProductInput) ([]*domain.ProductWithQty, error)
	}

	Store struct {
		db *pgxpool.Pool
	}
)

func (store *Store) PlaceProductsOnShelf(ctx *fasthttp.RequestCtx, input domain.PlaceRemoveProductInput) ([]*domain.ProductWithQty, error) {
	sqlShelfProduct := `select shelf_id, product_id, product_qty from shelf_product where shelf_id = @shelf_id`
	argsShelfProduct := pgx.NamedArgs{"shelf_id": input.ShelfID}
	shelfProduct, err := shared.GetManyRows[domain.ShelfProductDB](ctx, store.db, sqlShelfProduct, argsShelfProduct)
	if err != nil {
		return nil, err
	}

	var shelfProductIds []int
	shelfProductMap := make(map[int]domain.ShelfProductDB)
	for _, item := range shelfProduct {
		shelfProductMap[item.ProductID] = *item
		shelfProductIds = append(shelfProductIds, item.ProductID)
	}
	var productIds []int
	productsMap := make(map[int]domain.ProductIdQty)
	for _, item := range input.ProductsWithQty {
		productsMap[item.ProductID] = item
		productIds = append(productIds, item.ProductID)
	}
	productsToAdd := shared.DifferenceLeft(productIds, shelfProductIds)
	productsToAddQty := shared.IntersectUniq(shelfProductIds, productIds)

	for _, item := range productsToAdd {
		productToAdd := productsMap[item]
		sql := `insert into shelf_product (product_id, product_qty, shelf_id) values (@product_id, @product_qty, @shelf_id)`
		args := pgx.NamedArgs{"product_id": productToAdd.ProductID,
			"product_qty": productToAdd.Quantity,
			"shelf_id":    input.ShelfID}
		_, err := store.db.Exec(ctx, sql, args)
		if err != nil {
			return nil, err
		}
	}

	for _, item := range productsToAddQty {
		productToAddInput := productsMap[item]
		productToAddDb := shelfProductMap[item]

		sql := `update shelf_product SET 
			product_qty = @product_qty
             WHERE shelf_id = @shelf_id and product_id = @product_id`
		args := pgx.NamedArgs{
			"product_qty": productToAddInput.Quantity + productToAddDb.ProductQty,
			"shelf_id":    input.ShelfID,
			"product_id":  item}

		_, err := store.db.Exec(ctx, sql, args)
		if err != nil {
			return nil, err
		}
	}

	sqlProduct := `select product.*, shelf_product.product_qty from shelf_product left join product on
    product.id=shelf_product.product_id where shelf_product.shelf_id = @id`
	argsProduct := pgx.NamedArgs{"id": input.ShelfID}
	productWithQty, err := shared.GetManyRows[domain.ProductWithQty](ctx, store.db, sqlProduct, argsProduct)
	if err != nil {
		return nil, err
	}

	return productWithQty, nil
}

func (store *Store) RemoveProductsFromShelf(ctx *fasthttp.RequestCtx, input domain.PlaceRemoveProductInput) ([]*domain.ProductWithQty, error) {
	sqlShelfProduct := `select shelf_id, product_id, product_qty from shelf_product where shelf_id = @shelf_id`
	argsShelfProduct := pgx.NamedArgs{"shelf_id": input.ShelfID}
	shelfProduct, err := shared.GetManyRows[domain.ShelfProductDB](ctx, store.db, sqlShelfProduct, argsShelfProduct)
	if err != nil {
		return nil, err
	}

	var shelfProductIds []int
	shelfProductMap := make(map[int]domain.ShelfProductDB)
	for _, item := range shelfProduct {
		shelfProductMap[item.ProductID] = *item
		shelfProductIds = append(shelfProductIds, item.ProductID)
	}
	var productIds []int
	productsMap := make(map[int]domain.ProductIdQty)
	for _, item := range input.ProductsWithQty {
		productsMap[item.ProductID] = item
		productIds = append(productIds, item.ProductID)
	}

	intersect := shared.IntersectUniq(productIds, shelfProductIds)

	for _, item := range intersect {
		inputProduct := productsMap[item]
		shlfProduct := shelfProductMap[item]
		if shlfProduct.ProductQty-inputProduct.Quantity > 0 {
			sql := `update shelf_product SET 
			product_qty = @product_qty
             WHERE shelf_id = @shelf_id and product_id = @product_id`
			args := pgx.NamedArgs{"product_id": inputProduct.ProductID,
				"product_qty": shlfProduct.ProductQty - inputProduct.Quantity,
				"shelf_id":    input.ShelfID}

			_, err := store.db.Exec(ctx, sql, args)
			if err != nil {
				return nil, err
			}
		} else {
			sql := `delete from shelf_product where shelf_id = @shelf_id and product_id = @product_id`
			args := pgx.NamedArgs{"shelf_id": input.ShelfID, "product_id": item}

			_, err := store.db.Exec(ctx, sql, args)
			if err != nil {
				return nil, err
			}
		}
	}

	sqlProduct := `select product.*, shelf_product.product_qty from shelf_product left join product on
    product.id=shelf_product.product_id where shelf_product.shelf_id = @id`
	argsProduct := pgx.NamedArgs{"id": input.ShelfID}
	productWithQty, err := shared.GetManyRows[domain.ProductWithQty](ctx, store.db, sqlProduct, argsProduct)
	if err != nil {
		return nil, err
	}

	return productWithQty, nil
}

func NewOperationsStore(db *pgxpool.Pool) *Store {
	return &Store{
		db: db,
	}
}

func (store *Store) GetAssemblyInfoByOrders(ctx *fasthttp.RequestCtx, orders []int) (
	[]domain.AssemblyInfo,
	error,
) {
	sqlOrderProduct := `select id as product_id, name as product_name, serial, order_id, product_qty 
	from product left join order_product on product.id = order_product.product_id 
	where order_product.order_id = any ($1)`
	orderProduct, err := shared.GetManyRowsInIds[domain.OrderProductDB](ctx, store.db, sqlOrderProduct, orders)
	if err != nil {
		return nil, err
	}

	var productIdsToGetShelves []int
	productMap := make(map[int]*domain.OrderProductDB)
	for _, item := range orderProduct {
		productMap[item.ProductID] = item
		productIdsToGetShelves = append(productIdsToGetShelves, item.ProductID)
	}
	if len(productIdsToGetShelves) == 0 {
		return nil, errors.New("can't get shelves, put products on shelves")
	}
	sqlShelfProduct := `select shelf_id, product_id, product_qty from shelf_product where product_id = any ($1)`
	shelfProduct, err := shared.GetManyRowsInIds[domain.ShelfProductDB](ctx, store.db, sqlShelfProduct, productIdsToGetShelves)
	if err != nil {
		return nil, err
	}

	var shelvesIds []int
	productShelfMap := make(map[int]*domain.ShelfProductDB)
	for _, shlf := range shelfProduct {
		shelvesIds = append(shelvesIds, shlf.ShelfID)
		productShelfMap[shlf.ProductID] = shlf
	}
	sqlShelves := `select * from shelf where id = any ($1)`
	shelves, err := shared.GetManyRowsInIds[domain.Shelf](ctx, store.db, sqlShelves, productIdsToGetShelves)
	if err != nil {
		return nil, err
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
