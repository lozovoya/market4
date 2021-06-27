package repository

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type Storage struct {
	pool             *pgxpool.Pool
	marketRepository MarketRepository
}

func NewMarketRepository(pool *pgxpool.Pool) MarketRepository {
	return &Storage{
		pool: pool,
	}
}

func IsEmpty(field string) bool {
	return field == ""
}
