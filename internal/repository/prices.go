package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"market4/internal/model"
)

type priceRepo struct {
	pool *pgxpool.Pool
}

func NewPriceRepository(pool *pgxpool.Pool) Price {
	return &priceRepo{pool: pool}
}

func (price *priceRepo) AddPrice(ctx context.Context, p *model.Price) (int, error) {

	dbReq := "INSERT INTO prices (sale_price, factory_price, discount_price, product_id, is_active)" +
		"VALUES ($1, $2, $3, $4, $5)" +
		"RETURNING id"
	var id int
	err := price.pool.QueryRow(ctx, dbReq, p.SalePrice, p.FactoryPrice, p.DiscountPrice, p.ProductID, p.IsActive).Scan(&id)
	//TODO обработать ошибку отсутствия id продукта в БД
	if err != nil {
		return 0, fmt.Errorf("AddPrice: %w", err)
	}
	log.Printf("Price %d is added", id)
	return id, nil
}

func (price *priceRepo) EditPrice(ctx context.Context, p *model.Price) (*model.Price, error) {

	var dbReq = "UPDATE prices " +
		"SET sale_price=$1, factory_price=$2, discount_price=$3, is_active=$4, updated=CURRENT_TIMESTAMP " +
		"WHERE id = $5" +
		"RETURNING id, sale_price, factory_price, discount_price, is_active, product_id"
	var result model.Price
	err := price.pool.QueryRow(
		ctx,
		dbReq,
		p.SalePrice,
		p.FactoryPrice,
		p.DiscountPrice,
		p.IsActive,
		p.ID).Scan(&result.ID, &result.SalePrice, &result.FactoryPrice, &result.DiscountPrice, &result.IsActive, &result.ProductID)
	if err != nil {
		return nil, fmt.Errorf("EditPrice: %w", err)
	}

	log.Printf("Price for %s is updated", p.ProductID)
	return &result, nil
}

func (price *priceRepo) ListAllPrices(ctx context.Context) ([]*model.Price, error) {
	prices := make([]*model.Price, 0)

	dbReq := "SELECT id, sale_price, factory_price, discount_price, is_active, product_id " +
		"FROM prices"
	rows, err := price.pool.Query(ctx, dbReq)
	if err != nil {
		return prices, fmt.Errorf("ListAllPrices: %w", err)
	}
	for rows.Next() {
		var price model.Price
		err = rows.Scan(&price.ID,
			&price.SalePrice,
			&price.FactoryPrice,
			&price.DiscountPrice,
			&price.IsActive,
			&price.ProductID)
		if err != nil {
			return prices, fmt.Errorf("ListAllPrices: %w", err)
		}

		prices = append(prices, &price)
	}
	return prices, nil
}
