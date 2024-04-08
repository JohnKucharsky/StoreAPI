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
	[]*domain.AssemblyInfo,
	error,
) {
	ctx := context.Background()
	if len(m) == 0 {
		return nil, errors.New("you have to provide array of orders to get info")
	}

	// get orders in ids
	ordersRows, err := store.db.Query(
		ctx, `select product_id,order_id from order_product where order_id = any ($1);`,
		m,
	)
	if err != nil {
		return nil, err
	}

	orderProduct, err := pgx.CollectRows(
		ordersRows, pgx.RowToStructByName[domain.OrderProductDB],
	)
	if err != nil {
		return nil, err
	}
	fmt.Println(orderProduct)
	var idsToGetProducts []int
	for _, orderProductDB := range orderProduct {
		idsToGetProducts = append(idsToGetProducts, orderProductDB.ProductID)
	}

	if len(idsToGetProducts) == 0 {
		return nil, errors.New("you have to buy some products")
	}
	// get orders in ids

	// get products in ids
	productsRows, err := store.db.Query(
		ctx, `select * from product where id = any ($1);`,
		idsToGetProducts,
	)
	if err != nil {
		return nil, err
	}

	productRes, err := pgx.CollectRows(
		productsRows, pgx.RowToStructByName[domain.Product],
	)
	if err != nil {
		return nil, err
	}
	fmt.Println(productRes)

	return nil, nil
}
