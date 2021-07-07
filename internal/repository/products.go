package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"market4/internal/model"
)

type productRepo struct {
	pool         *pgxpool.Pool
	categoryRepo Category
	shopRepo     Shop
}

func NewProductRepository(pool *pgxpool.Pool, categoryRepo Category, shopRepo Shop) Product {
	return &productRepo{pool: pool, categoryRepo: categoryRepo, shopRepo: shopRepo}
}

func (p *productRepo) AddProduct(ctx context.Context, product *model.Product, shopId int, categoryId int) (string, error) {

	if !p.categoryRepo.IfCategoryExists(ctx, categoryId) {
		err := errors.New("category doesn't exist")
		return "", fmt.Errorf("AddProduct: %w", err)
	}
	if !p.shopRepo.IfShopExists(ctx, shopId) {
		err := errors.New("shop doesn't exist")
		return "", fmt.Errorf("AddProduct: %w", err)
	}

	dbReq := "INSERT INTO products(sku, name, uri, description, is_active)" +
		"VALUES ($1,$2,$3,$4,$5)" +
		"RETURNING id"
	pType := fmt.Sprintf("/product/%s-%s", product.Type, product.SKU)
	var id = ""

	err := p.pool.QueryRow(ctx, dbReq, product.SKU, product.Name, pType, product.Description, true).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("AddProduct: %w", err)
	}

	err = p.setProductCategory(ctx, categoryId, id)
	if err != nil {
		return "", fmt.Errorf("AddProduct: %w", err)
	}

	return id, nil
}

func (p *productRepo) setProductCategory(ctx context.Context, categoryId int, productId string) error {

	dbReq := "INSERT INTO productcategory (category_id, product_id)" +
		" VALUES ($1, $2)"
	//dbReq = fmt.Sprintf("%s (%d,'%s')", dbReq, categoryId, productId)
	log.Println(dbReq)
	_, err := p.pool.Exec(ctx, dbReq, categoryId, productId)
	if err != nil {
		return fmt.Errorf("SetProductCategory: %w", err)
	}
	return nil
}
