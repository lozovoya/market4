package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"market4/internal/model"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type productRepo struct {
	pool         *pgxpool.Pool
	categoryRepo Category
	shopRepo     Shop
	priceRepo    Price
}

func NewProductRepository(pool *pgxpool.Pool, categoryRepo Category, shopRepo Shop, priceRepo Price) Product {
	return &productRepo{pool: pool, categoryRepo: categoryRepo, shopRepo: shopRepo, priceRepo: priceRepo}
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
func (p *productRepo) AddProduct(ctx context.Context, product model.Product, shopId, categoryId int) (model.Product, error) {
	var result model.Product
	if !p.categoryRepo.IfCategoryExists(ctx, categoryId) {
		err := errors.New("category doesn't exist")
		return result, fmt.Errorf("AddProduct: %w", err)
	}
	if !p.shopRepo.IfShopExists(ctx, shopId) {
		err := errors.New("shop doesn't exist")
		return result, fmt.Errorf("AddProduct: %w", err)
	}

	dbReq := "INSERT INTO products(sku, name, uri, description, is_active)" +
		"VALUES ($1,$2,$3,$4,$5)" +
		"RETURNING id, sku, name, uri, description, is_active"
	uri := fmt.Sprintf("/product/%s-%s", product.Type, product.SKU)

	err := p.pool.QueryRow(ctx,
		dbReq,
		product.SKU,
		product.Name,
		uri,
		product.Description,
		true).Scan(&result.ID,
		&result.SKU,
		&result.Name,
		&result.URI,
		&result.Description,
		&result.IsActive)
	if err != nil {
		return result, fmt.Errorf("AddProduct: %w", err)
	}

	err = p.setProductCategory(ctx, categoryId, result.ID)
	if err != nil {
		return result, fmt.Errorf("AddProduct: %w", err)
	}

	err = p.setProductShop(ctx, shopId, result.ID)
	if err != nil {
		return result, fmt.Errorf("AddProduct: %w", err)
	}

	return result, nil
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
func (p *productRepo) EditProduct(ctx context.Context, product model.Product, shopID, categoryID int) (model.Product, error) {
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

	dbReq = fmt.Sprintf("%s is_active = %t, "+
		"updated = CURRENT_TIMESTAMP WHERE sku = '%s' "+
		"RETURNING id, sku, name, uri, description, is_active",
		dbReq, product.IsActive, product.SKU)

	var result model.Product
	err := p.pool.QueryRow(ctx, dbReq).Scan(&result.ID, &result.SKU, &result.Name, &result.URI, &result.Description, &result.IsActive)
	if err != nil {
		return result, fmt.Errorf("EditProduct: %w", err)
	}

	if shopID != 0 {
		dbReq = "UPDATE productshop SET shop_id = $1 WHERE product_id = $2"
		_, err = p.pool.Exec(ctx, dbReq, shopID, result.ID)
		if err != nil {
			return result, fmt.Errorf("EditProduct: %w", err)
		}
	}

	if categoryID != 0 {
		dbReq = "UPDATE productcategory SET category_id = $1 WHERE product_id = $2 "
		_, err = p.pool.Exec(ctx, dbReq, categoryID, result.ID)
		if err != nil {
			return result, fmt.Errorf("EditProduct: %w", err)
		}
	}
	return result, nil
}

func (p *productRepo) ListAllProducts(ctx context.Context) ([]model.Product, error) {
	products := make([]model.Product, 0)

	dbReq := "SELECT id, sku, name, uri, description, is_active " +
		"FROM products "
	rows, err := p.pool.Query(ctx, dbReq)
	if err != nil {
		if err == pgx.ErrNoRows {
			return products, nil
		}
		return products, fmt.Errorf("ListAllProducts: %w", err)
	}
	for rows.Next() {
		var product model.Product
		err = rows.Scan(&product.ID, &product.SKU, &product.Name, &product.URI, &product.Description, &product.IsActive)
		if err != nil {
			return products, fmt.Errorf("ListAllProducts: %w", err)
		}
		products = append(products, product)
	}
	return products, nil
}
func (p *productRepo) SearchProductsByCategory(ctx context.Context, category int) ([]model.Product, error) {
	products := make([]model.Product, 0)

	dbReq := "SELECT products.id, products.name, products.uri " +
		"FROM products " +
		"JOIN productcategory " +
		"ON products.id = productcategory.product_id " +
		"WHERE productcategory.category_id = $1 "
	rows, err := p.pool.Query(ctx, dbReq, category)
	if err != nil {
		if err == pgx.ErrNoRows {
			return products, nil
		}
		return products, fmt.Errorf("SearchProductsByCategory: %w", err)
	}
	for rows.Next() {
		var product model.Product
		err = rows.Scan(&product.ID, &product.Name, &product.URI)
		if err != nil {
			return products, fmt.Errorf("ListAllProducts: %w", err)
		}
		products = append(products, product)
	}
	return products, nil
}
func (p *productRepo) SearchProductsByName(ctx context.Context, productName string) (model.Product, error) {
	dbReq := "SELECT sku, name, uri, description, id " +
		"FROM products " +
		"WHERE name = $1"
	var product model.Product
	err := p.pool.QueryRow(ctx, dbReq, productName).Scan(&product.SKU, &product.Name, &product.URI, &product.Description, &product.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return product, nil
		}
		return product, fmt.Errorf("SearchProductsByName: %w", err)
	}

	return product, nil
}
func (p *productRepo) SearchProductsByShop(ctx context.Context, shopID int) ([]model.Product, error) {
	var products = make([]model.Product, 0)
	dbReq := "SELECT products.id, products.sku, products.name, products.uri, products.description, products.is_active " +
		"FROM products " +
		"JOIN productshop " +
		"ON products.id = productshop.product_id " +
		"WHERE productshop.shop_id = $1"

	rows, err := p.pool.Query(ctx, dbReq, shopID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return products, nil
		}
		return products, fmt.Errorf("SearchActiveProductsByShop: %w", err)
	}
	for rows.Next() {
		var product model.Product
		err = rows.Scan(
			&product.ID,
			&product.SKU,
			&product.Name,
			&product.URI,
			&product.Description,
			&product.IsActive,
		)
		if err != nil {
			return products, fmt.Errorf("SearchActiveProductsByShop: %w", err)
		}
		products = append(products, product)
	}
	return products, nil
}
