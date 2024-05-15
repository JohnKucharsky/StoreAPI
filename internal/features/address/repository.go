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
	AddressStoreI interface {
		Create(ctx *fasthttp.RequestCtx, m domain.AddressInput) (*domain.Address, error)
		GetMany(ctx *fasthttp.RequestCtx) ([]*domain.Address, error)
		GetOne(ctx *fasthttp.RequestCtx, id int) (*domain.Address, error)
		Update(ctx *fasthttp.RequestCtx, m domain.AddressInput, id int) (*domain.Address, error)
		Delete(ctx *fasthttp.RequestCtx, id int) (*int, error)
	}

	AddressStore struct {
		db *pgxpool.Pool
	}
)

func NewAddressStore(db *pgxpool.Pool) *AddressStore {
	return &AddressStore{
		db: db,
	}
}

func (store *AddressStore) Create(ctx *fasthttp.RequestCtx, m domain.AddressInput) (
	*domain.Address,
	error,
) {
	sql := `
        INSERT INTO address (city, street, house, floor, entrance, additional_info)
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

func (store *AddressStore) GetMany(ctx *fasthttp.RequestCtx) ([]*domain.Address, error) {
	rows, err := store.db.Query(
		ctx, `
		select * from address;
     `,
	)
	if err != nil {
		return nil, err
	}

	res, err := pgx.CollectRows(
		rows, pgx.RowToAddrOfStructByName[domain.Address],
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (store *AddressStore) GetOne(ctx *fasthttp.RequestCtx, id int) (*domain.Address, error) {
	rows, err := store.db.Query(
		ctx,
		`select * from address where id = @id`,
		pgx.NamedArgs{"id": id},
	)
	if err != nil {
		return nil, err
	}

	res, err := pgx.CollectExactlyOneRow(
		rows, pgx.RowToAddrOfStructByName[domain.Address],
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (store *AddressStore) Update(ctx *fasthttp.RequestCtx, m domain.AddressInput, id int) (*domain.Address, error) {
	rows, err := store.db.Query(
		ctx,
		`UPDATE address SET 
			city = @city,
			street = @street,
			house = @house,
			floor = @floor,
			entrance = @entrance,
			additional_info = @additional_info,
			updated_at = @updated_at
             WHERE id = @id 
             returning  id, city, street, house, floor, entrance, additional_info, created_at, updated_at`,
		pgx.NamedArgs{
			"id":              id,
			"city":            m.City,
			"street":          m.Street,
			"house":           m.House,
			"floor":           m.Floor,
			"entrance":        m.Entrance,
			"additional_info": m.AdditionalInfo,
			"updated_at":      time.Now(),
		},
	)
	if err != nil {
		return nil, err
	}

	res, err := pgx.CollectExactlyOneRow(
		rows, pgx.RowToAddrOfStructByName[domain.Address],
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (store *AddressStore) Delete(ctx *fasthttp.RequestCtx, id int) (*int, error) {
	rows, err := store.db.Query(
		ctx,
		`delete from address where id = @id 
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
