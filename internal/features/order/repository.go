package order

import (
	"errors"
	"fmt"
	"github.com/JohnKucharsky/StoreAPI/internal/domain"
	"github.com/JohnKucharsky/StoreAPI/internal/shared"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"
	"github.com/valyala/fasthttp"
	"strconv"
	"strings"
)

type (
	StoreI interface {
		Create(ctx *fasthttp.RequestCtx, m domain.OrderInput) (*domain.OrderShort, error)
		GetMany(ctx *fasthttp.RequestCtx) ([]*domain.OrderShort, error)
		GetOne(ctx *fasthttp.RequestCtx, id int) (*domain.OrderShort, error)
		GetAddress(ctx *fasthttp.RequestCtx, id int) (*domain.Address, error)
		Update(ctx *fasthttp.RequestCtx, m domain.OrderInput, id int) (*domain.OrderShort, error)
		Delete(ctx *fasthttp.RequestCtx, id int) (*int, error)
		GetProductsForOrder(ctx *fasthttp.RequestCtx, id int) ([]*domain.ProductWithQty, error)
	}

	Store struct {
		db *pgxpool.Pool
	}

	OrderProductDB struct {
		ProductID  int `db:"product_id"`
		ProductQty int `db:"product_qty"`
		OrderID    int `db:"order_id"`
	}

	idRes struct {
		ID int `db:"id"`
	}
)

func NewOrderStore(db *pgxpool.Pool) *Store {
	return &Store{
		db: db,
	}
}

func (store *Store) Create(ctx *fasthttp.RequestCtx, m domain.OrderInput) (
	*domain.OrderShort,
	error,
) {
	sql := `
        INSERT INTO orders (address_id, payment)
        VALUES (@address_id, @payment)
        RETURNING id, address_id, payment, created_at, updated_at`
	args := pgx.NamedArgs{
		"address_id": m.AddressID,
		"payment":    m.Payment,
	}

	order, err := shared.GetOneRow[domain.OrderShort](ctx, store.db, sql, args)
	if err != nil {
		return nil, err
	}

	if len(m.Products) == 0 {
		return nil, errors.New("products array is empty")
	}
	if err := store.BulkInsertProducts(ctx, order.ID, m.Products); err != nil {
		return nil, err
	}

	return order, nil
}

func (store *Store) GetMany(ctx *fasthttp.RequestCtx) ([]*domain.OrderShort, error) {
	sql := `select * from orders`

	return shared.GetManyRows[domain.OrderShort](ctx, store.db, sql, pgx.NamedArgs{})
}

func (store *Store) GetOne(ctx *fasthttp.RequestCtx, id int) (*domain.OrderShort, error) {
	sql := `select * from orders where id = @id`
	args := pgx.NamedArgs{"id": id}

	return shared.GetOneRow[domain.OrderShort](ctx, store.db, sql, args)
}

func (store *Store) GetAddress(ctx *fasthttp.RequestCtx, id int) (*domain.Address, error) {
	sql := `select * from address where id = @id`
	args := pgx.NamedArgs{"id": id}

	return shared.GetOneRow[domain.Address](ctx, store.db, sql, args)
}

func (store *Store) Update(ctx *fasthttp.RequestCtx, m domain.OrderInput, id int) (*domain.OrderShort, error) {
	sql := `UPDATE orders SET 
			address_id = @address_id,
			payment = @payment
             WHERE id = @id 
             returning  id,address_id,payment, created_at, updated_at`
	args := pgx.NamedArgs{
		"id":         id,
		"address_id": m.AddressID,
		"payment":    m.Payment,
	}

	res, err := shared.GetOneRow[domain.OrderShort](ctx, store.db, sql, args)
	if err != nil {
		return nil, err
	}

	if len(m.Products) != 0 {
		if err := store.BulkUpdateProducts(ctx, res.ID, m.Products); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("you should add at least one product")
	}

	return res, nil
}

func (store *Store) Delete(ctx *fasthttp.RequestCtx, id int) (*int, error) {
	sql := `delete from orders where id = @id returning id`
	args := pgx.NamedArgs{"id": id}

	one, err := shared.GetOneRow[idRes](ctx, store.db, sql, args)
	if err != nil {
		return nil, err
	}

	return &one.ID, nil
}

func (store *Store) BulkDeleteProducts(ctx *fasthttp.RequestCtx, orderID int, products []int) error {
	createParams := pgx.NamedArgs{
		"order_id": orderID,
	}
	var valuesStringArr []string

	for idx, prdct := range products {
		productString := strconv.Itoa(prdct)
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

func (store *Store) BulkInsertProducts(ctx *fasthttp.RequestCtx, orderID int, products []domain.ProductIdQty) error {
	createParams := pgx.NamedArgs{
		"order_id": orderID,
	}
	var valuesStringArr []string

	for idx, prdct := range products {
		pID := strconv.Itoa(prdct.ProductID)
		pQty := strconv.Itoa(prdct.Quantity)
		idxString := strconv.Itoa(idx + 1)

		valuesStringArr = append(valuesStringArr, fmt.Sprintf("(@order_id, @%s, @%s)",
			fmt.Sprintf("p%s", idxString),
			fmt.Sprintf("q%s", idxString)))
		createParams[fmt.Sprintf("p%s", idxString)] = pID
		createParams[fmt.Sprintf("q%s", idxString)] = pQty
	}

	sql := fmt.Sprintf(`
		insert into order_product (order_id, product_id,product_qty)
		values %s `, strings.Join(valuesStringArr, ", "),
	)

	_, err := store.db.Exec(ctx, sql, createParams)
	if err != nil {
		return err
	}

	return nil
}

func (store *Store) BulkUpdateProducts(ctx *fasthttp.RequestCtx, orderID int, products []domain.ProductIdQty) error {
	rows, err := store.db.Query(
		ctx, `select order_id,product_id,product_qty from order_product
    where order_id = @order_id`, pgx.NamedArgs{"order_id": orderID},
	)
	if err != nil {
		return err
	}

	orderProductDB, err := pgx.CollectRows(
		rows, pgx.RowToStructByName[OrderProductDB],
	)
	if err != nil {
		return err
	}

	// add or delete
	var productsDbIDs []int
	for _, prdct := range orderProductDB {
		productsDbIDs = append(productsDbIDs, prdct.ProductID)
	}
	var productsInputIDs []int
	for _, prdct := range products {
		productsInputIDs = append(productsInputIDs, prdct.ProductID)
	}
	productsIdsAdd, productsIdsToDelete := lo.Difference(productsInputIDs, productsDbIDs)
	var productsToAdd []domain.ProductIdQty
	for _, prdct := range products {
		for _, productID := range productsIdsAdd {
			if prdct.ProductID != productID {
				continue
			}
			productsToAdd = append(productsToAdd, prdct)

		}
	}
	if len(productsToAdd) != 0 {
		if err := store.BulkInsertProducts(ctx, orderID, productsToAdd); err != nil {
			return nil
		}
	}
	if len(productsIdsToDelete) != 0 {
		if err := store.BulkDeleteProducts(ctx, orderID, productsIdsToDelete); err != nil {
			return nil
		}
	}
	// add or delete end

	// change qty on products
	var filteredOrderProduct = lo.Filter(orderProductDB, func(item OrderProductDB, index int) bool {
		return !lo.Contains(productsIdsToDelete, item.ProductID)
	})

	var inputProductsMap = make(map[int]domain.ProductIdQty)
	for _, product := range products {
		inputProductsMap[product.ProductID] = product
	}

	for _, filtOrdProd := range filteredOrderProduct {
		if filtOrdProd.ProductQty != inputProductsMap[filtOrdProd.ProductID].Quantity {
			var product = inputProductsMap[filtOrdProd.ProductID]

			_, err := store.db.Exec(ctx, `
			UPDATE order_product SET 
			product_qty = @product_qty
             WHERE product_id = @product_id and order_id = @order_id`,
				pgx.NamedArgs{
					"product_qty": product.Quantity,
					"product_id":  product.ProductID,
					"order_id":    orderID,
				})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (store *Store) GetProductsForOrder(ctx *fasthttp.RequestCtx, id int) ([]*domain.ProductWithQty, error) {
	rows, err := store.db.Query(
		ctx, `
		select * from order_product  where order_id = @order_id;
     `, pgx.NamedArgs{"order_id": id},
	)
	if err != nil {
		return nil, err
	}

	orderProductDB, err := pgx.CollectRows(
		rows, pgx.RowToStructByName[OrderProductDB],
	)
	if err != nil {
		return nil, err
	}

	var idsToGetProducts []int
	for _, orderPDB := range orderProductDB {
		idsToGetProducts = append(idsToGetProducts, orderPDB.ProductID)
	}
	if len(idsToGetProducts) == 0 {
		return nil, errors.New("you have to buy some products")
	}

	// get product in ids
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
	// get product in ids end

	productMap := make(map[int]domain.Product)
	for _, product := range productRes {
		productMap[product.ID] = product
	}

	var response []*domain.ProductWithQty
	for _, orderProduct := range orderProductDB {
		var productWithQty = domain.ProductWithQty{
			Product:  productMap[orderProduct.ProductID],
			Quantity: orderProduct.ProductQty,
		}
		response = append(response, &productWithQty)
	}

	return response, nil
}
