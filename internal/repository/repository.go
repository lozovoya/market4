package repository

import (
	"context"
	"market4/internal/model"
)

type MarketRepository interface {
	ListAllShops(ctx context.Context) ([]*model.Shop, error)
	AddShop(ctx context.Context, s *model.Shop) (int, error)
	EditShop(ctx context.Context, s *model.Shop) error

	ListAllCategories(ctx context.Context) ([]*model.Category, error)
	AddCategory(ctx context.Context, c *model.Category) (int, error)
	EditCategory(ctx context.Context, c *model.Category) error
}
