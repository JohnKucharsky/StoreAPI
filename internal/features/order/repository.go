package order

import (
	"errors"
	"fmt"
	"github.com/JohnKucharsky/WarehouseAPI/internal/domain"
	"github.com/JohnKucharsky/WarehouseAPI/internal/shared"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"
	"github.com/valyala/fasthttp"
	"strconv"
	"strings"
	"time"
)

type (
	StoreI interface {
		Create(ctx *fasthttp.RequestCtx, m domain.OrderInput, userID uuid.UUID) (*domain.OrderShort, error)
		GetMany(ctx *fasthttp.RequestCtx, userID uuid.UUID) ([]*domain.OrderShort, error)
		GetOne(ctx *fasthttp.RequestCtx, id int, userID uuid.UUID) (*domain.OrderShort, error)
		GetAddress(ctx *fasthttp.RequestCtx, id int) (*domain.Address, error)
		Update(ctx *fasthttp.RequestCtx, m domain.OrderInput, id int) (*domain.OrderShort, error)
		Delete(ctx *fasthttp.RequestCtx, id int) (*int, error)
		GetProductsForOrder(ctx *fasthttp.RequestCtx, id int) ([]*domain.ProductWithQty, error)
	}

	Store struct {
		db *pgxpool.Pool
	}
)

func NewOrderStore(db *pgxpool.Pool) *Store {
	return &Store{
		db: db,
	}
}

func (store *Store) Create(ctx *fasthttp.RequestCtx, m domain.OrderInput, userID uuid.UUID) (
	*domain.OrderShort,
	error,
) {
	sql := `
        INSERT INTO orders (address_id, payment,user_id)
        VALUES (@address_id, @payment, @user_id)
        RETURNING id, address_id, payment, created_at, updated_at`
	args := pgx.NamedArgs{
		"address_id": m.AddressID,
		"payment":    m.Payment,
		"user_id":    userID,
	}

	one, err := shared.GetOneRow[domain.OrderShort](ctx, store.db, sql, args)
	if err != nil {
		return nil, err
	}

	if len(m.Products) == 0 {
		return nil, errors.New("products array is empty")
	}
	if err := store.BulkInsertProducts(ctx, one.ID, m.Products); err != nil {
		return nil, err
	}

	return one, nil
}

func (store *Store) GetMany(ctx *fasthttp.RequestCtx, userID uuid.UUID) ([]*domain.OrderShort, error) {
	sql := `select id,address_id,payment,created_at,updated_at from orders where user_id = @user_id`

	return shared.GetManyRows[domain.OrderShort](ctx, store.db, sql, pgx.NamedArgs{"user_id": userID})
}

func (store *Store) GetOne(ctx *fasthttp.RequestCtx, id int, userID uuid.UUID) (*domain.OrderShort, error) {
	sql := `select id,address_id,payment,created_at,updated_at from orders where id = @id and user_id = @user_id`
	args := pgx.NamedArgs{"id": id, "user_id": userID}

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
			payment = @payment,
            updated_at = @updated_at
             WHERE id = @id 
             returning  id,address_id,payment, created_at, updated_at`
	args := pgx.NamedArgs{
		"id":         id,
		"address_id": m.AddressID,
		"payment":    m.Payment,
		"updated_at": time.Now(),
	}

	one, err := shared.GetOneRow[domain.OrderShort](ctx, store.db, sql, args)
	if err != nil {
		return nil, err
	}

	if len(m.Products) != 0 {
		if err := store.BulkUpdateProducts(ctx, one.ID, m.Products); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("you should add at least one product")
	}

	return one, nil
}

func (store *Store) Delete(ctx *fasthttp.RequestCtx, id int) (*int, error) {
	sql := `delete from orders where id = @id returning id`
	args := pgx.NamedArgs{"id": id}

	one, err := shared.GetOneRow[domain.IdRes](ctx, store.db, sql, args)
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
		values %s`, strings.Join(valuesStringArr, ", "),
	)

	_, err := store.db.Exec(ctx, sql, createParams)
	if err != nil {
		return err
	}

	return nil
}

func (store *Store) BulkUpdateProducts(ctx *fasthttp.RequestCtx, orderID int, products []domain.ProductIdQty) error {
	sql := `select order_id,product_id,product_qty from order_product where order_id = @order_id`
	args := pgx.NamedArgs{"order_id": orderID}

	orderProductDB, err := shared.GetManyRowsToStructByName[domain.OrderProductDBShort](ctx, store.db, sql, args)
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
	var filteredOrderProduct = lo.Filter(orderProductDB, func(item domain.OrderProductDBShort, index int) bool {
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
	sql := `select order_id,product_id,product_qty from order_product where order_id = @order_id`
	args := pgx.NamedArgs{"order_id": id}

	orderProductDB, err := shared.GetManyRowsToStructByName[domain.OrderProductDBShort](ctx, store.db, sql, args)
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
	sqlProducts := `select * from product where id = any ($1)`

	productRes, err := shared.GetManyRowsInIds[domain.Product](ctx, store.db, sqlProducts, idsToGetProducts)
	if err != nil {
		return nil, err
	}
	// get product in ids end

	productMap := make(map[int]domain.Product)
	for _, product := range productRes {
		productMap[product.ID] = *product
	}

	var productsWithQty []*domain.ProductWithQty
	for _, orderProduct := range orderProductDB {
		var productWithQty = domain.ProductWithQty{
			Product:  productMap[orderProduct.ProductID],
			Quantity: orderProduct.ProductQty,
		}
		productsWithQty = append(productsWithQty, &productWithQty)
	}

	return productsWithQty, nil
}
