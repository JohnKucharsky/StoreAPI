package address

import (
	"github.com/JohnKucharsky/StoreAPI/internal/domain"
	"github.com/JohnKucharsky/StoreAPI/internal/shared"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valyala/fasthttp"
	"time"
)

type (
	StoreI interface {
		Create(ctx *fasthttp.RequestCtx, m domain.AddressInput) (*domain.Address, error)
		GetMany(ctx *fasthttp.RequestCtx) ([]*domain.Address, error)
		GetOne(ctx *fasthttp.RequestCtx, id int) (*domain.Address, error)
		Update(ctx *fasthttp.RequestCtx, m domain.AddressInput, id int) (*domain.Address, error)
		Delete(ctx *fasthttp.RequestCtx, id int) (*int, error)
	}

	Store struct {
		db *pgxpool.Pool
	}
)

func NewAddressStore(db *pgxpool.Pool) *Store {
	return &Store{
		db: db,
	}
}

func (store *Store) Create(ctx *fasthttp.RequestCtx, m domain.AddressInput) (
	*domain.Address,
	error,
) {
	sql := `INSERT INTO address (city, street, house, floor, entrance, additional_info)
        VALUES (@city, @street, @house, @floor, @entrance, @additional_info)
        RETURNING id, city, street, house, floor, entrance, additional_info, created_at, updated_at`
	args := pgx.NamedArgs{
		"city":            m.City,
		"street":          m.Street,
		"house":           m.House,
		"floor":           m.Floor,
		"entrance":        m.Entrance,
		"additional_info": m.AdditionalInfo,
	}

	return shared.GetOneRow[domain.Address](ctx, store.db, sql, args)

}

func (store *Store) GetMany(ctx *fasthttp.RequestCtx) ([]*domain.Address, error) {
	sql := `select * from address`

	return shared.GetManyRows[domain.Address](ctx, store.db, sql, pgx.NamedArgs{})
}

func (store *Store) GetOne(ctx *fasthttp.RequestCtx, id int) (*domain.Address, error) {
	sql := `select * from address where id = @id`
	args := pgx.NamedArgs{"id": id}

	return shared.GetOneRow[domain.Address](ctx, store.db, sql, args)
}

func (store *Store) Update(ctx *fasthttp.RequestCtx, m domain.AddressInput, id int) (*domain.Address, error) {
	sql := `UPDATE address SET 
			city = @city,
			street = @street,
			house = @house,
			floor = @floor,
			entrance = @entrance,
			additional_info = @additional_info,
			updated_at = @updated_at
             WHERE id = @id 
             returning  id, city, street, house, floor, entrance, additional_info, created_at, updated_at`
	args := pgx.NamedArgs{
		"id":              id,
		"city":            m.City,
		"street":          m.Street,
		"house":           m.House,
		"floor":           m.Floor,
		"entrance":        m.Entrance,
		"additional_info": m.AdditionalInfo,
		"updated_at":      time.Now(),
	}

	return shared.GetOneRow[domain.Address](ctx, store.db, sql, args)
}

func (store *Store) Delete(ctx *fasthttp.RequestCtx, id int) (*int, error) {
	sql := `delete from address where id = @id 
        returning id`
	args := pgx.NamedArgs{"id": id}

	one, err := shared.GetOneRow[domain.IdRes](ctx, store.db, sql, args)
	if err != nil {
		return nil, err
	}

	return &one.ID, nil
}
