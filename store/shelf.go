package store

import (
	"context"
	"github.com/JohnKucharsky/StoreAPI/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"
)

type ShelfStore struct {
	db *pgxpool.Pool
}

func NewShelfStore(db *pgxpool.Pool) *ShelfStore {
	return &ShelfStore{
		db: db,
	}
}

func (as *ShelfStore) Create(m domain.ShelfInput) (
	*domain.Shelf,
	error,
) {
	ctx := context.Background()

	rows, err := as.db.Query(
		ctx, `
        INSERT INTO shelf (name, destination)
        VALUES (@name, @destination)
        RETURNING id, name, destination, created_at, updated_at`,
		pgx.NamedArgs{
			"name":        m.Name,
			"destination": m.Destination,
		},
	)
	if err != nil {
		return nil, err
	}

	res, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToAddrOfStructByName[domain.Shelf],
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (as *ShelfStore) GetMany() ([]*domain.Shelf, error) {
	ctx := context.Background()

	rows, err := as.db.Query(
		ctx, `
		select * from shelf;
     `,
	)
	if err != nil {
		return nil, err
	}

	res, err := pgx.CollectRows(
		rows, pgx.RowToAddrOfStructByName[domain.Shelf],
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (as *ShelfStore) GetRandomShelfID() (*int, error) {
	ctx := context.Background()

	rows, err := as.db.Query(
		ctx, `
		select * from shelf limit 20;
     `,
	)
	if err != nil {
		return nil, err
	}

	res, err := pgx.CollectRows(
		rows, pgx.RowToAddrOfStructByName[domain.Shelf],
	)
	if err != nil {
		return nil, err
	}

	var allIDs []int
	for _, shelf := range res {
		allIDs = append(allIDs, shelf.ID)
	}

	random := lo.Sample(allIDs)

	return &random, nil
}

func (as *ShelfStore) GetOne(id int) (*domain.Shelf, error) {
	ctx := context.Background()

	rows, err := as.db.Query(
		ctx,
		`select * from shelf where id = @id`,
		pgx.NamedArgs{"id": id},
	)
	if err != nil {
		return nil, err
	}

	res, err := pgx.CollectExactlyOneRow(
		rows, pgx.RowToAddrOfStructByName[domain.Shelf],
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (as *ShelfStore) Update(m domain.ShelfInput, id int) (*domain.Shelf, error) {
	ctx := context.Background()

	rows, err := as.db.Query(
		ctx,
		`UPDATE shelf SET 
			name = @name,
			destination = @destination
             WHERE id = @id 
             returning  id, name, destination, created_at, updated_at`,
		pgx.NamedArgs{
			"id":          id,
			"name":        m.Name,
			"destination": m.Destination,
		},
	)
	if err != nil {
		return nil, err
	}

	res, err := pgx.CollectExactlyOneRow(
		rows, pgx.RowToAddrOfStructByName[domain.Shelf],
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (as *ShelfStore) Delete(id int) (*int, error) {
	ctx := context.Background()

	rows, err := as.db.Query(
		ctx,
		`delete from shelf where id = @id 
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
