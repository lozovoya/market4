package repository

import (
	"context"
	"fmt"
	"market4/internal/model"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type priceRepo struct {
	pool *pgxpool.Pool
}

func NewPriceRepository(pool *pgxpool.Pool) Price {
	return &priceRepo{pool: pool}
}

func (price *priceRepo) AddPrice(ctx context.Context, p *model.Price) (model.Price, error) {
	dbReq := "INSERT INTO prices (sale_price, factory_price, discount_price, product_id, is_active)" +
		"VALUES ($1, $2, $3, $4, $5)" +
		"RETURNING sale_price, factory_price, discount_price"
	var newPrice model.Price
	err := price.pool.QueryRow(ctx,
		dbReq,
		p.SalePrice,
		p.FactoryPrice,
		p.DiscountPrice,
		p.ProductID,
		p.IsActive).Scan(&newPrice.SalePrice,
		&newPrice.FactoryPrice,
		&newPrice.DiscountPrice)
	if err != nil {
		return newPrice, fmt.Errorf("AddPrice: %w", err)
	}
	return newPrice, nil
}

func (price *priceRepo) EditPrice(ctx context.Context, p *model.Price) (model.Price, error) {
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
		p.ID).Scan(&result.ID,
		&result.SalePrice,
		&result.FactoryPrice,
		&result.DiscountPrice,
		&result.IsActive,
		&result.ProductID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return result, nil
		}
		return result, fmt.Errorf("EditPrice: %w", err)
	}
	return result, nil
}

func (price *priceRepo) EditPriceByProductID(ctx context.Context, p *model.Price) (model.Price, error) {
	var dbReq = "UPDATE prices " +
		"SET sale_price=$1, factory_price=$2, discount_price=$3, is_active=$4, updated=CURRENT_TIMESTAMP " +
		"WHERE product_id = $5" +
		"RETURNING id, sale_price, factory_price, discount_price, is_active, product_id"
	var result model.Price
	err := price.pool.QueryRow(
		ctx,
		dbReq,
		p.SalePrice,
		p.FactoryPrice,
		p.DiscountPrice,
		p.IsActive,
		p.ProductID).Scan(&result.ID, &result.SalePrice, &result.FactoryPrice, &result.DiscountPrice, &result.IsActive, &result.ProductID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return result, nil
		}
		return result, fmt.Errorf("EditPriceByProductID: %w", err)
	}
	return result, nil
}

func (price *priceRepo) ListAllPrices(ctx context.Context) ([]model.Price, error) {
	prices := make([]model.Price, 0)

	dbReq := "SELECT id, sale_price, factory_price, discount_price, is_active, product_id " +
		"FROM prices"
	rows, err := price.pool.Query(ctx, dbReq)
	if err != nil {
		if err == pgx.ErrNoRows {
			return prices, nil
		}
		return prices, fmt.Errorf("ListAllPrices: %w", err)
	}
	for rows.Next() {
		var result model.Price
		err = rows.Scan(&result.ID,
			&result.SalePrice,
			&result.FactoryPrice,
			&result.DiscountPrice,
			&result.IsActive,
			&result.ProductID)
		if err != nil {
			return prices, fmt.Errorf("ListAllPrices: %w", err)
		}

		prices = append(prices, result)
	}
	return prices, nil
}
func (price *priceRepo) SearchPriceByProductID(ctx context.Context, productID string) (model.Price, error) {
	dbReq := fmt.Sprintf("SELECT id, sale_price, factory_price, discount_price, product_id, is_active "+
		"FROM prices "+
		"WHERE product_id = '%s'",
		productID)
	var productPrice model.Price
	err := price.pool.QueryRow(ctx, dbReq).Scan(
		&productPrice.ID,
		&productPrice.SalePrice,
		&productPrice.FactoryPrice,
		&productPrice.DiscountPrice,
		&productPrice.ProductID,
		&productPrice.IsActive)
	if err != nil {
		if err == pgx.ErrNoRows {
			return productPrice, nil
		}
		return productPrice, fmt.Errorf("SearchPriceByProductID: %w", err)
	}
	return productPrice, nil
}
