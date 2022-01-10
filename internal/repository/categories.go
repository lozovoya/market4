package repository

import (
	"context"
	"fmt"
	"log"
	"market4/internal/model"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type categoryRepo struct {
	pool *pgxpool.Pool
}

func NewCategoryRepository(pool *pgxpool.Pool) Category {
	return &categoryRepo{pool: pool}
}
func (c *categoryRepo) IfCategoryExists(ctx context.Context, category int) bool {
	dbReq := "SELECT id FROM categories WHERE id=$1"
	var id = 0
	err := c.pool.QueryRow(ctx, dbReq, category).Scan(&id)
	if err != nil {
		log.Println(fmt.Errorf("IfCategoryExists: %w", err))
		return false
	}
	if id != 0 {
		return true
	}
	return false
}

func (c *categoryRepo) ListAllCategories(ctx context.Context) ([]model.Category, error) {
	categories := make([]model.Category, 0)

	dbReq := "SELECT id, name, uri_name " +
		"FROM categories"
	rows, err := c.pool.Query(ctx, dbReq)
	if err != nil {
		if err == pgx.ErrNoRows {
			return categories, nil
		}
		return categories, fmt.Errorf("ListAllCategories: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var category model.Category
		err = rows.Scan(&category.ID, &category.Name, &category.URI_name)
		if err != nil {
			return categories, fmt.Errorf("ListAllCategories: %w", err)
		}
		categories = append(categories, category)
	}
	return categories, nil
}
func (c *categoryRepo) AddCategory(ctx context.Context, category *model.Category) (int, error) {
	dbReq := "INSERT INTO categories (name) " +
		"VALUES ($1) " +
		"RETURNING id"
	var id int
	err := c.pool.QueryRow(ctx, dbReq, category.Name).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("AddCategory: %w", err)
	}

	dbReq = "UPDATE categories " +
		"SET uri_name = $1 " +
		"WHERE id = $2"
	uri_name := fmt.Sprintf("%s-%d", category.Name, id)
	_, err = c.pool.Exec(ctx, dbReq, uri_name, id)
	if err != nil {
		return 0, fmt.Errorf("AddCategory: %w", err)
	}
	return id, nil
}
func (c *categoryRepo) EditCategory(ctx context.Context, category *model.Category) error {
	dbReq := fmt.Sprintf("UPDATE categories SET name = '%s', "+
		"uri_name = '%s-%d', updated = CURRENT_TIMESTAMP WHERE id = %d",
		category.Name, category.Name, category.ID, category.ID)
	_, err := c.pool.Exec(ctx, dbReq)
	if err != nil {
		return fmt.Errorf("UpdateCategoryParameter: %w", err)
	}
	return nil
}
