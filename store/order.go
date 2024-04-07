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

type OrderStore struct {
	db *pgxpool.Pool
}

func NewOrderStore(db *pgxpool.Pool) *OrderStore {
	return &OrderStore{
		db: db,
	}
}

func (store *OrderStore) Create(m domain.OrderInput) (
	*domain.OrderShort,
	error,
) {
	ctx := context.Background()

	rows, err := store.db.Query(
		ctx, `
        INSERT INTO orders (user_id, address_id, total, payment)
        VALUES (@user_id, @address_id, @total, @payment)
        RETURNING id, user_id, address_id, total, payment, created_at, updated_at`,
		pgx.NamedArgs{
			"user_id":    m.UserID,
			"address_id": m.AddressID,
			"total":      m.Total,
			"payment":    m.Payment,
		},
	)
	if err != nil {
		return nil, err
	}

	res, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToAddrOfStructByName[domain.OrderShort],
	)
	if err != nil {
		return nil, err
	}

	if len(m.Products) == 0 {
		return nil, errors.New("categories array is empty")
	}
	var productIds []int
	for _, product := range m.Products {
		for range product.Quantity {
			productIds = append(productIds, product.ProductID)
		}
	}
	if err := store.BulkInsertProducts(res.ID, productIds); err != nil {
		return nil, err
	}

	return res, nil
}

func (store *OrderStore) GetMany() ([]*domain.OrderShort, error) {
	ctx := context.Background()

	rows, err := store.db.Query(
		ctx, `
		select * from orders;
     `,
	)
	if err != nil {
		return nil, err
	}

	res, err := pgx.CollectRows(
		rows, pgx.RowToAddrOfStructByName[domain.OrderShort],
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (store *OrderStore) GetOne(id int) (*domain.OrderShort, error) {
	ctx := context.Background()

	rows, err := store.db.Query(
		ctx,
		`select * from orders where id = @id`,
		pgx.NamedArgs{"id": id},
	)
	if err != nil {
		return nil, err
	}

	res, err := pgx.CollectExactlyOneRow(
		rows, pgx.RowToAddrOfStructByName[domain.OrderShort],
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (store *OrderStore) Update(m domain.OrderInput, id int) (*domain.OrderShort, error) {
	ctx := context.Background()

	rows, err := store.db.Query(
		ctx,
		`UPDATE orders SET 
			user_id = @user_id,
			address_id = @address_id,
			total = @total,
			payment = @payment
             WHERE id = @id 
             returning  id, user_id,address_id,total,payment, created_at, updated_at`,
		pgx.NamedArgs{
			"id":         id,
			"user_id":    m.UserID,
			"address_id": m.AddressID,
			"total":      m.Total,
			"payment":    m.Payment,
		},
	)
	if err != nil {
		return nil, err
	}

	res, err := pgx.CollectExactlyOneRow(
		rows, pgx.RowToAddrOfStructByName[domain.OrderShort],
	)
	if err != nil {
		return nil, err
	}

	if len(m.Products) != 0 {
		var productIds []int
		for _, product := range m.Products {
			for range product.Quantity {
				productIds = append(productIds, product.ProductID)
			}
		}
		if err := store.BulkUpdateProducts(res.ID, productIds); err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (store *OrderStore) Delete(id int) (*int, error) {
	ctx := context.Background()

	rows, err := store.db.Query(
		ctx,
		`delete from orders where id = @id 
        returning id`,
		pgx.NamedArgs{
			"id": id,
		},
	)
	if err != nil {
		return nil, err
	}

	type idRes struct {
		ID int `db:"id"`
	}

	res, err := pgx.CollectExactlyOneRow(
		rows, pgx.RowToAddrOfStructByName[idRes],
	)
	if err != nil {
		return nil, err
	}

	return &res.ID, nil
}

func (store *OrderStore) BulkDeleteProducts(orderID int, products []int) error {
	ctx := context.Background()

	createParams := pgx.NamedArgs{
		"order_id": orderID,
	}
	var valuesStringArr []string

	for idx, product := range products {
		productString := strconv.Itoa(product)
		idxString := strconv.Itoa(idx + 1)

		valuesStringArr = append(valuesStringArr, fmt.Sprintf("@%s", fmt.Sprintf("ord%s", idxString)))
		createParams[fmt.Sprintf("ord%s", idxString)] = productString
	}

	sql := fmt.Sprintf(`
		delete from order_product where order_id = @order_id and
		order_product.product_id in (%s) `, strings.Join(valuesStringArr, ", "),
	)

	_, err := store.db.Exec(ctx, sql, createParams)
	if err != nil {
		return err
	}

	return nil
}

func (store *OrderStore) BulkInsertProducts(orderID int, products []int) error {
	ctx := context.Background()

	createParams := pgx.NamedArgs{
		"order_id": orderID,
	}
	var valuesStringArr []string

	for idx, product := range products {
		productString := strconv.Itoa(product)
		idxString := strconv.Itoa(idx + 1)

		valuesStringArr = append(valuesStringArr, fmt.Sprintf("(@order_id, @%s)", fmt.Sprintf("ord%s", idxString)))
		createParams[fmt.Sprintf("ord%s", idxString)] = productString
	}

	sql := fmt.Sprintf(`
		insert into order_product (order_id, product_id)
		values %s `, strings.Join(valuesStringArr, ", "),
	)

	_, err := store.db.Exec(ctx, sql, createParams)
	if err != nil {
		return err
	}

	return nil
}

func (store *OrderStore) BulkUpdateProducts(orderID int, products []int) error {
	ctx := context.Background()

	rows, err := store.db.Query(
		ctx, `select order_id,product_id from order_product where order_id = @id`, pgx.NamedArgs{"id": orderID},
	)
	if err != nil {
		return err
	}

	resProducts, err := pgx.CollectRows(
		rows, pgx.RowToAddrOfStructByName[domain.OrderProductDB],
	)
	if err != nil {
		return err
	}

	var productsDbIDs []int
	for _, product := range resProducts {
		productsDbIDs = append(productsDbIDs, product.ProductID)
	}

	productsToAdd, productsToDelete := lo.Difference(products, productsDbIDs)
	if len(productsToAdd) != 0 {
		if err := store.BulkInsertProducts(orderID, productsToAdd); err != nil {
			return nil
		}
	}
	if len(productsToDelete) != 0 {
		if err := store.BulkDeleteProducts(orderID, productsToDelete); err != nil {
			return nil
		}
	}

	return nil
}

func (store *OrderStore) GetProductsForOrder(id int) ([]*domain.ProductWithQty, error) {
	ctx := context.Background()

	rows, err := store.db.Query(
		ctx, `
		select product.*, count(product.id) quantity from order_product 
	left join product on product_id = product.id  where order_id = @id group by product.id;
     `, pgx.NamedArgs{"id": id},
	)
	if err != nil {
		return nil, err
	}

	res, err := pgx.CollectRows(
		rows, pgx.RowToAddrOfStructByName[domain.ProductWithQty],
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}
