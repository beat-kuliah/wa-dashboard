package sqlc

import "github.com/jackc/pgx/v5/pgxpool"

func NewFromPool(pool *pgxpool.Pool) *Queries {
	return New(pool)
}
