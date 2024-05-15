package product

import (
	"errors"
	"fmt"
	"github.com/JohnKucharsky/StoreAPI/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"
	"github.com/valyala/fasthttp"
	"strconv"
	"strings"
)

type (
	ProductStoreI interface {
		Create(ctx *fasthttp.RequestCtx, m domain.ProductInput) (*domain.Product, error)
		GetMany(ctx *fasthttp.RequestCtx) ([]*domain.Product, error)
		GetOne(ctx *fasthttp.RequestCtx, id int) (*domain.Product, error)
		Update(ctx *fasthttp.RequestCtx, m domain.ProductInput, id int) (*domain.Product, error)
		Delete(ctx *fasthttp.RequestCtx, id int) (*int, error)
	}

	ProductStore struct {
		db *pgxpool.Pool
	}
)

func NewProductStore(db *pgxpool.Pool) *ProductStore {
	return &ProductStore{
		db: db,
	}
}

func (store *ProductStore) BulkDeleteShelves(ctx *fasthttp.RequestCtx, productID int, shelvesIDs []int) error {
	params := pgx.NamedArgs{
		"product_id": productID,
	}
	var valuesStringArr []string

	for idx, shelf := range shelvesIDs {
		shelfID := strconv.Itoa(shelf)
		idx := strconv.Itoa(idx + 1)

		valuesStringArr = append(valuesStringArr, fmt.Sprintf("@%s", fmt.Sprintf("p%s", idx)))
		params[fmt.Sprintf("p%s", idx)] = shelfID
	}

	sql := fmt.Sprintf(`
		delete from shelf_product where product_id = @product_id and
		shelf_product.shelf_id in (%s) `, strings.Join(valuesStringArr, ", "),
	)

	_, err := store.db.Exec(ctx, sql, params)
	if err != nil {
		return err
	}

	return nil
}

func (store *ProductStore) BulkInsertShelves(ctx *fasthttp.RequestCtx, productID int, shelves []domain.ShelfIdQty) error {
	params := pgx.NamedArgs{
		"product_id": productID,
	}
	var valuesStringArr []string

	for idx, shelf := range shelves {
		sID := strconv.Itoa(shelf.ShelfID)
		pQty := strconv.Itoa(shelf.Quantity)
		idxString := strconv.Itoa(idx + 1)

		valuesStringArr = append(valuesStringArr, fmt.Sprintf("(@product_id, @%s, @%s)",
			fmt.Sprintf("p%s", idxString),
			fmt.Sprintf("q%s", idxString)))
		params[fmt.Sprintf("p%s", idxString)] = sID
		params[fmt.Sprintf("q%s", idxString)] = pQty
	}

	sql := fmt.Sprintf(`
		insert into shelf_product (product_id, shelf_id, product_qty)
		values %s `, strings.Join(valuesStringArr, ", "),
	)

	_, err := store.db.Exec(ctx, sql, params)
	if err != nil {
		return err
	}

	return nil
}

func (store *ProductStore) BulkUpdateShelves(ctx *fasthttp.RequestCtx, productID int, shelves []domain.ShelfIdQty) error {
	rows, err := store.db.Query(
		ctx, `select shelf_id,product_id,product_qty from shelf_product
    where product_id = @product_id`, pgx.NamedArgs{"product_id": productID},
	)
	if err != nil {
		return err
	}

	shelfProductDB, err := pgx.CollectRows(
		rows, pgx.RowToStructByName[domain.ShelfProductDB],
	)
	if err != nil {
		return err
	}

	// add or delete
	var shelvesDbIDs []int
	for _, shelf := range shelfProductDB {
		shelvesDbIDs = append(shelvesDbIDs, shelf.ShelfID)
	}
	var shelvesInputIDs []int
	for _, shelf := range shelves {
		shelvesInputIDs = append(shelvesInputIDs, shelf.ShelfID)
	}
	shelvesIdsToAdd, shelvesIdsToDelete := lo.Difference(shelvesInputIDs, shelvesDbIDs)
	var shelvesToAdd []domain.ShelfIdQty
	for _, shelf := range shelves {
		for _, shelfID := range shelvesIdsToAdd {
			if shelf.ShelfID != shelfID {
				continue
			}
			shelvesToAdd = append(shelvesToAdd, shelf)

		}
	}
	if len(shelvesToAdd) != 0 {
		if err := store.BulkInsertShelves(ctx, productID, shelvesToAdd); err != nil {
			return nil
		}
	}
	if len(shelvesIdsToDelete) != 0 {
		if err := store.BulkDeleteShelves(ctx, productID, shelvesIdsToDelete); err != nil {
			return nil
		}
	}
	// add or delete and

	// change qty on products
	var filteredShelfProduct = lo.Filter(shelfProductDB, func(item domain.ShelfProductDB, index int) bool {
		return !lo.Contains(shelvesIdsToDelete, item.ShelfID)
	})

	var inputShelvesMap = make(map[int]domain.ShelfIdQty)
	for _, shelf := range shelves {
		inputShelvesMap[shelf.ShelfID] = shelf
	}

	for _, filtShelfProduct := range filteredShelfProduct {
		if filtShelfProduct.ProductQty != inputShelvesMap[filtShelfProduct.ShelfID].Quantity {
			var shelf = inputShelvesMap[filtShelfProduct.ShelfID]

			_, err := store.db.Exec(ctx, `
			UPDATE shelf_product SET 
			product_qty = @product_qty
             WHERE product_id = @product_id and shelf_id = @shelf_id`,
				pgx.NamedArgs{
					"product_qty": shelf.Quantity,
					"product_id":  productID,
					"shelf_id":    shelf.ShelfID,
				})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (store *ProductStore) Create(ctx *fasthttp.RequestCtx, m domain.ProductInput) (
	*domain.Product,
	error,
) {
	rows, err := store.db.Query(
		ctx, `
        INSERT INTO product (main_shelf_id, name, serial, price, model, picture_url)
        VALUES (@main_shelf_id, @name, @serial, @price, @model, @picture_url)
        RETURNING id, main_shelf_id, name, serial, price, model, picture_url, created_at, updated_at`,
		pgx.NamedArgs{
			"main_shelf_id": m.MainShelfID,
			"name":          m.Name,
			"serial":        m.Serial,
			"price":         m.Price,
			"model":         m.Model,
			"picture_url":   m.PictureURL,
		},
	)
	if err != nil {
		return nil, err
	}

	res, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToAddrOfStructByName[domain.Product],
	)
	if err != nil {
		return nil, err
	}

	if len(m.Shelves) != 0 {
		if err := store.BulkInsertShelves(ctx, res.ID, m.Shelves); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("you should add at least one shelf")
	}

	return res, nil
}

func (store *ProductStore) GetMany(ctx *fasthttp.RequestCtx) ([]*domain.Product, error) {
	rows, err := store.db.Query(
		ctx, `
		select * from product;
     `,
	)
	if err != nil {
		return nil, err
	}

	res, err := pgx.CollectRows(
		rows, pgx.RowToAddrOfStructByName[domain.Product],
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (store *ProductStore) GetOne(ctx *fasthttp.RequestCtx, id int) (*domain.Product, error) {
	rows, err := store.db.Query(
		ctx,
		`select * from product where id = @id`,
		pgx.NamedArgs{"id": id},
	)
	if err != nil {
		return nil, err
	}

	res, err := pgx.CollectExactlyOneRow(
		rows, pgx.RowToAddrOfStructByName[domain.Product],
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (store *ProductStore) Update(ctx *fasthttp.RequestCtx, m domain.ProductInput, id int) (*domain.Product, error) {
	rows, err := store.db.Query(
		ctx,
		`update product SET 
			main_shelf_id = @main_shelf_id,
			name = @name,
			serial = @serial,
			price = @price,
			model = @model,
			picture_url = @picture_url
             WHERE id = @id 
             returning  id, main_shelf_id, name, serial, price, model, picture_url, created_at, updated_at`,
		pgx.NamedArgs{
			"id":            id,
			"main_shelf_id": m.MainShelfID,
			"name":          m.Name,
			"serial":        m.Serial,
			"price":         m.Price,
			"model":         m.Model,
			"picture_url":   m.PictureURL,
		},
	)
	if err != nil {
		return nil, err
	}

	res, err := pgx.CollectExactlyOneRow(
		rows, pgx.RowToAddrOfStructByName[domain.Product],
	)
	if err != nil {
		return nil, err
	}

	if len(m.Shelves) != 0 {
		if err := store.BulkUpdateShelves(ctx, res.ID, m.Shelves); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("you should add at least one shelf")
	}

	return res, nil
}

func (store *ProductStore) Delete(ctx *fasthttp.RequestCtx, id int) (*int, error) {
	rows, err := store.db.Query(
		ctx,
		`delete from product where id = @id 
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
