package product

import (
	"fmt"
	"github.com/JohnKucharsky/WarehouseAPI/internal/domain"
	"github.com/JohnKucharsky/WarehouseAPI/internal/shared"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valyala/fasthttp"
	"time"
)

type (
	StoreI interface {
		Create(ctx *fasthttp.RequestCtx, m domain.ProductInput) (*domain.Product, error)
		GetMany(ctx *fasthttp.RequestCtx, params *shared.ParsedPaginationParams, orderString string) ([]*domain.Product, int32, error)
		GetOne(ctx *fasthttp.RequestCtx, id int) (*domain.Product, error)
		Update(ctx *fasthttp.RequestCtx, m domain.ProductInput, id int) (*domain.Product, error)
		Delete(ctx *fasthttp.RequestCtx, id int) (*int, error)
	}

	Store struct {
		db *pgxpool.Pool
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

func (store *Store) GetMany(ctx *fasthttp.RequestCtx, pp *shared.ParsedPaginationParams, orderString string) ([]*domain.Product, int32, error) {
	sql := fmt.Sprintf(`select * from product %s limit @limit offset @offset`, orderString)
	args := pgx.NamedArgs{"limit": pp.Limit, "offset": pp.Offset}

	many, err := shared.GetManyRows[domain.Product](ctx, store.db, sql, args)
	if err != nil {
		return nil, 0, err
	}
	var total int32

	err = store.db.QueryRow(ctx, `select count(*) from product`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return many, total, nil
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
			picture_url = @picture_url,
			updated_at = @updated_at
             WHERE id = @id 
             returning  id, name, serial, price, model, picture_url, created_at, updated_at`
	args := pgx.NamedArgs{
		"id":          id,
		"name":        m.Name,
		"serial":      m.Serial,
		"price":       m.Price,
		"model":       m.Model,
		"picture_url": m.PictureURL,
		"updated_at":  time.Now(),
	}

	return shared.GetOneRow[domain.Product](ctx, store.db, sql, args)
}

func (store *Store) Delete(ctx *fasthttp.RequestCtx, id int) (*int, error) {
	sql := `delete from product where id = @id returning id`
	args := pgx.NamedArgs{
		"id": id,
	}

	one, err := shared.GetOneRow[domain.IdRes](ctx, store.db, sql, args)
	if err != nil {
		return nil, err
	}

	return &one.ID, nil
}
