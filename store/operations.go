package store

import (
	"context"
	"errors"
	"fmt"
	"github.com/JohnKucharsky/StoreAPI/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"
	"strconv"
	"strings"
)

type OperationsStore struct {
	db *pgxpool.Pool
}

func NewOperationsStore(db *pgxpool.Pool) *OperationsStore {
	return &OperationsStore{
		db: db,
	}
}

func (store *OperationsStore) BulkDeleteProducts(shelfID int, products []int) error {
	ctx := context.Background()

	createParams := pgx.NamedArgs{
		"shelf_id": shelfID,
	}
	var valuesStringArr []string

	for idx, product := range products {
		shelfString := strconv.Itoa(product)
		idxString := strconv.Itoa(idx + 1)

		valuesStringArr = append(valuesStringArr, fmt.Sprintf("@%s", fmt.Sprintf("p%s", idxString)))
		createParams[fmt.Sprintf("p%s", idxString)] = shelfString
	}

	sql := fmt.Sprintf(`
		delete from shelf_product where shelf_id = @shelf_id and
		shelf_product.product_id in (%s) `, strings.Join(valuesStringArr, ", "),
	)

	_, err := store.db.Exec(ctx, sql, createParams)
	if err != nil {
		return err
	}

	return nil
}

func (store *OperationsStore) BulkInsertProducts(shelfID int, products []int) error {
	ctx := context.Background()

	createParams := pgx.NamedArgs{
		"shelf_id": shelfID,
	}
	var valuesStringArr []string

	for idx, product := range products {
		shelfString := strconv.Itoa(product)
		idxString := strconv.Itoa(idx + 1)

		valuesStringArr = append(valuesStringArr, fmt.Sprintf("(@shelf_id, @%s)", fmt.Sprintf("s%s", idxString)))
		createParams[fmt.Sprintf("s%s", idxString)] = shelfString
	}

	sql := fmt.Sprintf(`
		insert into shelf_product (shelf_id, product_id)
		values %s `, strings.Join(valuesStringArr, ", "),
	)

	_, err := store.db.Exec(ctx, sql, createParams)
	if err != nil {
		return err
	}

	return nil
}

func (store *OperationsStore) BulkUpdateProducts(shelfID int, products []int) error {
	ctx := context.Background()

	rows, err := store.db.Query(
		ctx, `select shelf_id, product_id from shelf_product where shelf_id = @id`, pgx.NamedArgs{"id": shelfID},
	)
	if err != nil {
		return err
	}

	resProducts, err := pgx.CollectRows(
		rows, pgx.RowToAddrOfStructByName[domain.ShelfProductDB],
	)
	if err != nil {
		return err
	}

	var productDBids []int
	for _, product := range resProducts {
		productDBids = append(productDBids, product.ProductID)
	}

	productsToAdd, productsToDelete := lo.Difference(products, productDBids)
	if len(productsToAdd) != 0 {
		if err := store.BulkInsertProducts(shelfID, productsToAdd); err != nil {
			return nil
		}
	}
	if len(productsToDelete) != 0 {
		if err := store.BulkDeleteProducts(shelfID, productsToDelete); err != nil {
			return nil
		}
	}

	return nil
}

func (store *OperationsStore) OneShelfToManyProducts(m domain.PlaceProductInput) (
	*domain.PlaceProduct,
	error,
) {
	if err := store.BulkUpdateProducts(m.ShelfID, m.ProductIDs); err != nil {
		return nil, err
	}

	ctx := context.Background()

	rows, err := store.db.Query(
		ctx, `select product.* from shelf_product left join product on shelf_product.product_id = product.id
         where shelf_product.shelf_id = @id`, pgx.NamedArgs{"id": m.ShelfID},
	)
	if err != nil {
		return nil, err
	}

	res, err := pgx.CollectRows(
		rows, pgx.RowToStructByName[domain.Product],
	)
	if err != nil {
		return nil, err
	}

	return &domain.PlaceProduct{
		ShelfID:  m.ShelfID,
		Products: res,
	}, nil
}

func (store *OperationsStore) GetAssemblyInfoByOrders(m []int) (
	[]domain.AssemblyInfo,
	error,
) {
	ctx := context.Background()
	if len(m) == 0 {
		return nil, errors.New("you have to provide array of orders to get info")
	}

	// get order_product in ids
	orderProductRows, err := store.db.Query(
		ctx, `select product_id,order_id from order_product where order_id = any ($1);`,
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
		ctx, `select shelf_id, product_id from shelf_product where product_id = any ($1);`,
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

	for _, shelf := range shelfProduct {
		shelvesIds = append(shelvesIds, shelf.ShelfID)
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
	productShelfMap := make(map[int]*int)
	for _, shelfP := range shelfProduct {
		productShelfMap[shelfP.ProductID] = &shelfP.ShelfID
	}
	productMap := make(map[int]ProductShort)
	for _, product := range productRes {
		productMap[product.ID] = product
	}

	var assemblyInfo []domain.AssemblyInfo

	for _, shelf := range shelves {
		var assemblyProducts []domain.AssemblyProduct

		for _, product := range shelfProduct {
			if shelf.ID != product.ShelfID {
				continue
			}

			var orderID = productOrderMap[product.ProductID]
			var additionalShelf *int
			if productMap[product.ProductID].MainShelfID != nil {
				additionalShelf = productShelfMap[product.ProductID]
			}

			assemblyProduct := domain.AssemblyProduct{
				ProductID:       product.ProductID,
				ProductName:     productMap[product.ProductID].Name,
				OrderID:         orderID,
				Quantity:        0,
				AdditionalShelf: additionalShelf,
			}

			assemblyProducts = append(assemblyProducts, assemblyProduct)
		}

		var innerRes = domain.AssemblyInfo{Name: shelf.Name, Products: assemblyProducts}
		assemblyInfo = append(assemblyInfo, innerRes)
	}

	return assemblyInfo, nil
}
