package store

import (
	"context"
	"github.com/JohnKucharsky/StoreAPI/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AddressStore struct {
	db *pgxpool.Pool
}

func NewAddressStore(db *pgxpool.Pool) *AddressStore {
	return &AddressStore{
		db: db,
	}
}

func (as *AddressStore) Create(m domain.AddressInput) (
	*domain.Address,
	error,
) {
	ctx := context.Background()

	rows, err := as.db.Query(
		ctx, `
        INSERT INTO address (city, street, house, floor, entrance, additional_info)
        VALUES (@city, @street, @house, @floor, @entrance, @additional_info)
        RETURNING id, city, street, house, floor, entrance, additional_info, created_at, updated_at`,
		pgx.NamedArgs{
			"city":            m.City,
			"street":          m.Street,
			"house":           m.House,
			"floor":           m.Floor,
			"entrance":        m.Entrance,
			"additional_info": m.AdditionalInfo,
		},
	)
	if err != nil {
		return nil, err
	}

	res, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToAddrOfStructByName[domain.Address],
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (as *AddressStore) GetMany() ([]*domain.Address, error) {
	ctx := context.Background()

	rows, err := as.db.Query(
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

func (as *AddressStore) GetOne(id int) (*domain.Address, error) {
	ctx := context.Background()

	rows, err := as.db.Query(
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

func (as *AddressStore) Update(m domain.AddressInput, id int) (*domain.Address, error) {
	ctx := context.Background()

	rows, err := as.db.Query(
		ctx,
		`UPDATE address SET 
			city = @city,
			street = @street,
			house = @house,
			floor = @floor,
			entrance = @entrance,
			additional_info = @additional_info
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

func (as *AddressStore) Delete(id int) (*int, error) {
	ctx := context.Background()

	rows, err := as.db.Query(
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
