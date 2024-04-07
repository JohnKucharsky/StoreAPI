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

	var idsStr []string
	for _, id := range m {
		s := strconv.Itoa(id)
		idsStr = append(idsStr, s)
	}

	rows, err := store.db.Query(
		ctx, `select orders.id as order_id, product.id as product_id,product.name as product_name, 
	 count(product.id) as quantity, shelf.name as shelf_name from orders 
	 left join order_product on orders.id = order_product.order_id 
	 left join product on order_product.product_id = product.id
	 left join shelf_product on product.id = shelf_product.product_id
	 left join shelf on shelf_product.shelf_id = shelf.id
	 where orders.id in(@ids) group by shelf.name,orders.id,product.id;`, pgx.NamedArgs{"ids": strings.Join(idsStr, "::int, ")},
	)
	fmt.Println(strings.Join(idsStr, "::int, "))
	if err != nil {
		return nil, err
	}

	res, err := pgx.CollectRows(
		rows, pgx.RowToAddrOfStructByName[domain.AssemblyInfo],
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}
