package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"market4/internal/model"
)

type priceRepo struct {
	pool        *pgxpool.Pool
	productRepo Product
}

func NewPriceRepository(pool *pgxpool.Pool, productRepo Product) Price {
	return &priceRepo{pool: pool, productRepo: productRepo}
}

func (price *priceRepo) AddPrice(ctx context.Context, p *model.Price, productID string) (int, error) {

	if !price.productRepo.IfProductExists(ctx, productID) {
		err := errors.New("product doesn't exist")
		return 0, fmt.Errorf("AddPrice: %w", err)
	}

	dbReq := "INSERT INTO prices (sale_price, factory_price, discount_price, product_id, is_active)" +
		"VALUES ($1, $2, $3, $4, $5)" +
		"RETURNING id"
	var id int
	err := price.pool.QueryRow(ctx, dbReq, p.SalePrice, p.FactoryPrice, p.DiscountPrice, productID, p.IsActive).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("AddPrice: %w", err)
	}
	log.Printf("Price %d is added", id)
	return id, nil
}

func (price *priceRepo) EditPrice(ctx context.Context, p *model.Price, productID string) (*model.Price, error) {

	var dbReq = "UPDATE prices " +
		"SET sale_price=$1, factory_price=$2, discount_price=$3, is_active=$4, updated=CURRENT_TIMESTAMP " +
		"WHERE product_id = $5" +
		"RETURNING id, sale_price, factory_price, discount_price, is_active"
	var result model.Price
	err := price.pool.QueryRow(
		ctx,
		dbReq,
		p.SalePrice,
		p.FactoryPrice,
		p.DiscountPrice,
		p.IsActive,
		productID).Scan(&result.ID, &result.SalePrice, &result.FactoryPrice, &result.DiscountPrice, &result.IsActive)
	if err != nil {
		return nil, fmt.Errorf("EditPrice: %w", err)
	}

	log.Printf("Price for %s is updated", productID)
	return &result, nil
}
