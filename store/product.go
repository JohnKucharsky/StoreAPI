package store

import (
	"context"
	"github.com/JohnKucharsky/StoreAPI/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductStore struct {
	db *pgxpool.Pool
}

func NewProductStore(db *pgxpool.Pool) *ProductStore {
	return &ProductStore{
		db: db,
	}
}

func (as *ProductStore) Create(m domain.ProductInput) (
	*domain.Product,
	error,
) {
	ctx := context.Background()

	rows, err := as.db.Query(
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

	return res, nil
}

func (as *ProductStore) GetMany() ([]*domain.Product, error) {
	ctx := context.Background()

	rows, err := as.db.Query(
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

func (as *ProductStore) GetOne(id int) (*domain.Product, error) {
	ctx := context.Background()

	rows, err := as.db.Query(
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

func (as *ProductStore) Update(m domain.ProductInput, id int) (*domain.Product, error) {
	ctx := context.Background()

	rows, err := as.db.Query(
		ctx,
		`UPDATE product SET 
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
			"picture_url":   m.AdditionalInfo,
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

	return res, nil
}

func (as *ProductStore) Delete(id int) (*int, error) {
	ctx := context.Background()

	rows, err := as.db.Query(
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
