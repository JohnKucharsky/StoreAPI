package shelf

import (
	"github.com/JohnKucharsky/StoreAPI/internal/domain"
	"github.com/JohnKucharsky/StoreAPI/internal/shared"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valyala/fasthttp"
)

type (
	ShelfStoreI interface {
		Create(ctx *fasthttp.RequestCtx, m domain.ShelfInput) (*domain.Shelf, error)
		GetMany(ctx *fasthttp.RequestCtx) ([]*domain.Shelf, error)
		GetOne(ctx *fasthttp.RequestCtx, id int) (*domain.Shelf, error)
		Update(ctx *fasthttp.RequestCtx, m domain.ShelfInput, id int) (*domain.Shelf, error)
		Delete(ctx *fasthttp.RequestCtx, id int) (*int, error)
	}

	ShelfStore struct {
		db *pgxpool.Pool
	}
)

func NewShelfStore(db *pgxpool.Pool) *ShelfStore {
	return &ShelfStore{
		db: db,
	}
}

func (store *ShelfStore) Create(ctx *fasthttp.RequestCtx, m domain.ShelfInput) (
	*domain.Shelf,
	error,
) {
	sql := `
        INSERT INTO shelf (name)
        VALUES (@name)
        RETURNING id, name, created_at, updated_at`
	args := pgx.NamedArgs{
		"name":        m.Name,
		"destination": m.Destination,
	}

	return shared.GetOneRow[domain.Shelf](ctx, store.db, sql, args)

}

func (store *ShelfStore) GetMany(ctx *fasthttp.RequestCtx) ([]*domain.Shelf, error) {
	sql := `select * from shelf`

	return shared.GetManyRows[domain.Shelf](ctx, store.db, sql, pgx.NamedArgs{})
}

func (store *ShelfStore) GetOne(ctx *fasthttp.RequestCtx, id int) (*domain.Shelf, error) {
	sql := `select * from shelf where id = @id`
	args := pgx.NamedArgs{"id": id}

	return shared.GetOneRow[domain.Shelf](ctx, store.db, sql, args)

}

func (store *ShelfStore) Update(ctx *fasthttp.RequestCtx, m domain.ShelfInput, id int) (*domain.Shelf, error) {
	sql := `UPDATE shelf SET 
			name = @name,
             WHERE id = @id 
             returning  id, name, created_at, updated_at`
	args := pgx.NamedArgs{
		"id":          id,
		"name":        m.Name,
		"destination": m.Destination,
	}

	return shared.GetOneRow[domain.Shelf](ctx, store.db, sql, args)
}

func (store *ShelfStore) Delete(ctx *fasthttp.RequestCtx, id int) (*int, error) {
	sql := `delete from shelf where id = @id 
        returning id`
	args := pgx.NamedArgs{
		"id": id,
	}

	type IdRes struct {
		ID int `db:"id"`
	}

	res, err := shared.GetOneRow[IdRes](ctx, store.db, sql, args)

	if err != nil {
		return nil, err
	}

	return &res.ID, nil
}
