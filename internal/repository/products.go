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

func (p *productRepo) IfProductExists(ctx context.Context, productID string) bool {
	dbReq := "SELECT id FROM products WHERE id=$1"
	var id = ""
	err := p.pool.QueryRow(ctx, dbReq, productID).Scan(&id)
	if err != nil {
		log.Println(fmt.Errorf("IfProductExists: %w", err))
		return false
	}
	if id != "" {
		return true
	}
	return false
}

func (p *productRepo) AddProduct(ctx context.Context, product *model.Product, shopId int, categoryId int) (*model.Product, error) {

	if !p.categoryRepo.IfCategoryExists(ctx, categoryId) {
		err := errors.New("category doesn't exist")
		return nil, fmt.Errorf("AddProduct: %w", err)
	}
	if !p.shopRepo.IfShopExists(ctx, shopId) {
		err := errors.New("shop doesn't exist")
		return nil, fmt.Errorf("AddProduct: %w", err)
	}

	dbReq := "INSERT INTO products(sku, name, uri, description, is_active)" +
		"VALUES ($1,$2,$3,$4,$5)" +
		"RETURNING id, sku, name, uri, description, is_active"
	uri := fmt.Sprintf("/product/%s-%s", product.Type, product.SKU)
	var result model.Product

	err := p.pool.QueryRow(ctx, dbReq, product.SKU, product.Name, uri, product.Description, true).Scan(&result.ID, &result.SKU, &result.Name, &result.URI, &result.Description, &result.IsActive)
	if err != nil {
		return nil, fmt.Errorf("AddProduct: %w", err)
	}

	err = p.setProductCategory(ctx, categoryId, result.ID)
	if err != nil {
		return nil, fmt.Errorf("AddProduct: %w", err)
	}

	err = p.setProductShop(ctx, shopId, result.ID)
	if err != nil {
		return nil, fmt.Errorf("AddProduct: %w", err)
	}

	return &result, nil
}

func (p *productRepo) setProductCategory(ctx context.Context, categoryId int, productId string) error {

	dbReq := "INSERT INTO productcategory (category_id, product_id)" +
		" VALUES ($1, $2)"
	_, err := p.pool.Exec(ctx, dbReq, categoryId, productId)
	if err != nil {
		return fmt.Errorf("SetProductCategory: %w", err)
	}
	return nil
}

func (p *productRepo) setProductShop(ctx context.Context, shopID int, productID string) error {
	dbReq := "INSERT INTO productshop (shop_id, product_id)" +
		"VALUES ($1, $2)"
	_, err := p.pool.Exec(ctx, dbReq, shopID, productID)
	if err != nil {
		return fmt.Errorf("SetProductShop: %w", err)
	}
	return nil
}

func (p *productRepo) EditProduct(ctx context.Context, product *model.Product, shopID int, categoryID int) (*model.Product, error) {

	var dbReq = "UPDATE products SET "

	if !IsEmpty(product.Name) {
		dbReq = fmt.Sprintf("%s name = '%s', ", dbReq, product.Name)
	}
	if !IsEmpty(product.Description) {
		dbReq = fmt.Sprintf("%s description = '%s', ", dbReq, product.Description)
	}
	if !IsEmpty(product.Type) {
		dbReq = fmt.Sprintf("%s uri = '/product/%s-%s', ", dbReq, product.Type, product.SKU)
	}

	dbReq = fmt.Sprintf("%s is_active = %t, updated = CURRENT_TIMESTAMP WHERE sku = '%s' RETURNING id, sku, name, uri, description, is_active", dbReq, product.IsActive, product.SKU)

	var result model.Product
	err := p.pool.QueryRow(ctx, dbReq).Scan(&result.ID, &result.SKU, &result.Name, &result.URI, &result.Description, &result.IsActive)
	if err != nil {
		return nil, fmt.Errorf("EditProduct: %w", err)
	}

	if shopID != 0 {
		dbReq = "UPDATE productshop SET shop_id = $1 WHERE product_id = $2"
		_, err = p.pool.Exec(ctx, dbReq, shopID, result.ID)
		if err != nil {
			return &result, fmt.Errorf("EditProduct: %w", err)
		}
	}

	if categoryID != 0 {
		dbReq = "UPDATE productcategory SET category_id = $1 WHERE product_id = $2 "
		_, err = p.pool.Exec(ctx, dbReq, categoryID, result.ID)
		if err != nil {
			return &result, fmt.Errorf("EditProduct: %w", err)
		}
	}

	log.Printf("Product %d updated", result.ID)
	return &result, nil
}

func (p *productRepo) ListAllProducts(ctx context.Context) ([]*model.Product, error) {
	products := make([]*model.Product, 0)

	dbReq := "SELECT products.id, products.sku, products.name, products.uri, products.description, products.is_active " +
		"prices.sale_price, prices.factory_price, prices.discount_price" +
		"FROM products, prices" +
		"WHERE products.id = prices.product_id "
	rows, err := p.pool.Query(ctx, dbReq)
	if err != nil {
		return products, fmt.Errorf("ListAllProducts: %w", err)
	}
	for rows.Next() {
		var product model.Product
		err = rows.Scan(&product.ID, &product.SKU, &product.Name, &product.URI, &product.Description, &product.IsActive)
		if err != nil {
			return products, fmt.Errorf("ListAllProducts: %w", err)
		}
		products = append(products, &product)
	}
	return products, nil
}
