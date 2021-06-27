package repository

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type Storage struct {
	pool             *pgxpool.Pool
	marketRepository MarketRepository
}

type market struct {
	storage *Storage
}

func NewStorage(pool *pgxpool.Pool) *Storage {
	return &Storage{pool: pool}
}

func IsEmpty(field string) bool {
	return field == ""
}
