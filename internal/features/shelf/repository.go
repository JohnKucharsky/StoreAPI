package shelf

import (
	"github.com/JohnKucharsky/WarehouseAPI/internal/domain"
	"github.com/JohnKucharsky/WarehouseAPI/internal/shared"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valyala/fasthttp"
	"time"
)

type (
	StoreI interface {
		Create(ctx *fasthttp.RequestCtx, m domain.ShelfInput) (*domain.Shelf, error)
		GetMany(ctx *fasthttp.RequestCtx) ([]*domain.Shelf, error)
		GetOne(ctx *fasthttp.RequestCtx, id int) (*domain.ShelfInfo, error)
		Update(ctx *fasthttp.RequestCtx, m domain.ShelfInput, id int) (*domain.Shelf, error)
		Delete(ctx *fasthttp.RequestCtx, id int) (*int, error)
	}

	Store struct {
		db *pgxpool.Pool
	}
)

func NewShelfStore(db *pgxpool.Pool) *Store {
	return &Store{
		db: db,
	}
}

func (store *Store) Create(ctx *fasthttp.RequestCtx, m domain.ShelfInput) (
	*domain.Shelf,
	error,
) {
	sql := `
        INSERT INTO shelf (name)
        VALUES (@name)
        RETURNING id, name, created_at, updated_at`
	args := pgx.NamedArgs{
		"name": m.Name,
	}

	return shared.GetOneRow[domain.Shelf](ctx, store.db, sql, args)
}

func (store *Store) GetMany(ctx *fasthttp.RequestCtx) ([]*domain.Shelf, error) {
	sql := `select * from shelf`

	return shared.GetManyRows[domain.Shelf](ctx, store.db, sql, pgx.NamedArgs{})
}

func (store *Store) GetOne(ctx *fasthttp.RequestCtx, id int) (*domain.ShelfInfo, error) {
	sql := `select * from shelf where id = @id`
	args := pgx.NamedArgs{"id": id}
	shelf, err := shared.GetOneRow[domain.Shelf](ctx, store.db, sql, args)
	if err != nil {
		return nil, err
	}
	sqlProduct := `select shelf_product.product_qty, product.*  
					from shelf_product left join product on
    				product.id=shelf_product.product_id where shelf_product.shelf_id = @id`
	argsProduct := pgx.NamedArgs{"id": shelf.ID}
	productWithQty, err := shared.GetManyRowsToStructByName[domain.ProductWithQty](ctx, store.db, sqlProduct, argsProduct)
	if err != nil {
		return nil, err
	}

	return &domain.ShelfInfo{
		Shelf:   *shelf,
		Product: productWithQty,
	}, nil
}

func (store *Store) Update(ctx *fasthttp.RequestCtx, m domain.ShelfInput, id int) (*domain.Shelf, error) {
	sql := `UPDATE shelf SET 
			name = @name,
			updated_at = @updated_at
             WHERE id = @id 
             returning  id, name, created_at, updated_at`
	args := pgx.NamedArgs{
		"id":         id,
		"name":       m.Name,
		"updated_at": time.Now(),
	}

	return shared.GetOneRow[domain.Shelf](ctx, store.db, sql, args)
}

func (store *Store) Delete(ctx *fasthttp.RequestCtx, id int) (*int, error) {
	sql := `delete from shelf where id = @id returning id`
	args := pgx.NamedArgs{
		"id": id,
	}

	one, err := shared.GetOneRow[domain.IdRes](ctx, store.db, sql, args)
	if err != nil {
		return nil, err
	}

	return &one.ID, nil
}
