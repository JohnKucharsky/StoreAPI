package product

import (
	"github.com/JohnKucharsky/StoreAPI/internal/domain"
	"github.com/JohnKucharsky/StoreAPI/internal/shared"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valyala/fasthttp"
)

type (
	StoreI interface {
		Create(ctx *fasthttp.RequestCtx, m domain.ProductInput) (*domain.Product, error)
		GetMany(ctx *fasthttp.RequestCtx) ([]*domain.Product, error)
		GetOne(ctx *fasthttp.RequestCtx, id int) (*domain.Product, error)
		Update(ctx *fasthttp.RequestCtx, m domain.ProductInput, id int) (*domain.Product, error)
		Delete(ctx *fasthttp.RequestCtx, id int) (*int, error)
	}

	Store struct {
		db *pgxpool.Pool
	}

	idRes struct {
		ID int `db:"id"`
	}
)

func NewProductStore(db *pgxpool.Pool) *Store {
	return &Store{
		db: db,
	}
}

func (store *Store) Create(ctx *fasthttp.RequestCtx, m domain.ProductInput) (
	*domain.Product,
	error,
) {
	sql := `INSERT INTO product (name, serial, price, model, picture_url)
        VALUES (@name, @serial, @price, @model, @picture_url)
        RETURNING id,  name, serial, price, model, picture_url, created_at, updated_at`
	args := pgx.NamedArgs{
		"name":        m.Name,
		"serial":      m.Serial,
		"price":       m.Price,
		"model":       m.Model,
		"picture_url": m.PictureURL,
	}

	return shared.GetOneRow[domain.Product](ctx, store.db, sql, args)
}

func (store *Store) GetMany(ctx *fasthttp.RequestCtx) ([]*domain.Product, error) {
	sql := `select * from product`

	return shared.GetManyRows[domain.Product](ctx, store.db, sql, pgx.NamedArgs{})
}

func (store *Store) GetOne(ctx *fasthttp.RequestCtx, id int) (*domain.Product, error) {
	sql := `select * from product where id = @id`
	args := pgx.NamedArgs{"id": id}

	return shared.GetOneRow[domain.Product](ctx, store.db, sql, args)
}

func (store *Store) Update(ctx *fasthttp.RequestCtx, m domain.ProductInput, id int) (*domain.Product, error) {
	sql := `update product SET 
			name = @name,
			serial = @serial,
			price = @price,
			model = @model,
			picture_url = @picture_url
             WHERE id = @id 
             returning  id, name, serial, price, model, picture_url, created_at, updated_at`
	args := pgx.NamedArgs{
		"id":          id,
		"name":        m.Name,
		"serial":      m.Serial,
		"price":       m.Price,
		"model":       m.Model,
		"picture_url": m.PictureURL,
	}

	return shared.GetOneRow[domain.Product](ctx, store.db, sql, args)
}

func (store *Store) Delete(ctx *fasthttp.RequestCtx, id int) (*int, error) {
	sql := `delete from product where id = @id returning id`
	args := pgx.NamedArgs{
		"id": id,
	}

	res, err := shared.GetOneRow[idRes](ctx, store.db, sql, args)
	if err != nil {
		return nil, err
	}

	return &res.ID, nil
}
