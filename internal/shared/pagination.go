package shared

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type Pagination struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type PaginationParams struct {
	Limit  string `json:"limit" validate:"omitempty,numeric"`
	Offset string `json:"offset" validate:"omitempty,numeric"`
}

type ParsedPaginationParams struct {
	Limit  int
	Offset *int
}

func GetPaginationParams(c *fiber.Ctx) (*ParsedPaginationParams, error) {
	var pp PaginationParams
	if err := BindQueries(c, &pp); err != nil {
		return nil, err
	}

	limit, _ := strconv.Atoi(pp.Limit)
	offset, _ := strconv.Atoi(pp.Offset)

	if limit > 0 {
		if offset > 0 {
			return &ParsedPaginationParams{
				Limit:  limit,
				Offset: &offset,
			}, nil
		}
		return &ParsedPaginationParams{
			Limit:  limit,
			Offset: nil,
		}, nil
	}

	return nil, nil
}
