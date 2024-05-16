package shared

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valyala/fasthttp"
)

func GetManyRows[T any](ctx *fasthttp.RequestCtx, db *pgxpool.Pool, sql string, params pgx.NamedArgs) ([]*T, error) {
	rows, err := db.Query(ctx, sql, params)
	if err != nil {
		return nil, err
	}

	res, err := pgx.CollectRows(
		rows, pgx.RowToAddrOfStructByName[T],
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func GetManyRowsInIds[T any](ctx *fasthttp.RequestCtx, db *pgxpool.Pool, sql string, ids []int) ([]*T, error) {
	rows, err := db.Query(ctx, sql, ids)
	if err != nil {
		return nil, err
	}

	res, err := pgx.CollectRows(
		rows, pgx.RowToAddrOfStructByName[T],
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func GetOneRow[T any](ctx *fasthttp.RequestCtx, db *pgxpool.Pool, sql string, params pgx.NamedArgs) (*T, error) {
	rows, err := db.Query(ctx, sql, params)
	if err != nil {
		return nil, err
	}

	res, err := pgx.CollectExactlyOneRow(
		rows, pgx.RowToAddrOfStructByName[T],
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}
