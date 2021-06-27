package repository

import (
	"context"
	"fmt"
	"log"
	"market4/internal/model"
)

func (s *Storage) ListAllCategories(ctx context.Context) ([]*model.Category, error) {
	categories := make([]*model.Category, 0)

	dbReq := "SELECT id, name, uri_name " +
		"FROM categories"
	rows, err := s.pool.Query(ctx, dbReq)
	if err != nil {
		return categories, fmt.Errorf("ListAllCategories: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var category model.Category
		rows.Scan(&category.Id, &category.Name, &category.Uri_name)
		if err != nil {
			log.Println(err)
			return categories, fmt.Errorf("ListAllCategories: %w", err)
		}
		categories = append(categories, &category)
	}

	return categories, nil
}

func (s *Storage) AddCategory(ctx context.Context, category *model.Category) (int, error) {

	dbReq := "INSERT INTO categories (name) " +
		"VALUES ($1) " +
		"RETURNING id"
	var id int
	err := s.pool.QueryRow(ctx, dbReq, category.Name).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("AddCategory: %w", err)
	}

	dbReq = "UPDATE categories " +
		"SET uri_name = $1 " +
		"WHERE id = $2"
	uri_name := fmt.Sprintf("%s-%d", category.Name, id)
	_, err = s.pool.Exec(ctx, dbReq, uri_name, id)
	if err != nil {
		return 0, fmt.Errorf("AddCategory: %w", err)
	}

	log.Printf("Category %d is added", id)
	return id, nil
}

func (s *Storage) EditCategory(ctx context.Context, category *model.Category) error {
	var dbReq = "UPDATE categories SET "

	if !IsEmpty(category.Name) {
		dbReq = fmt.Sprintf("%s name = '%s',", dbReq, category.Name)
	}

	dbReq = fmt.Sprintf("%s uri_name = '%s-%d', updated = CURRENT_TIMESTAMP WHERE id = %d", dbReq, category.Name, category.Id, category.Id)

	_, err := s.pool.Exec(ctx, dbReq)
	if err != nil {
		return fmt.Errorf("UpdateCategoryParameter: %w", err)
	}

	return nil
}
