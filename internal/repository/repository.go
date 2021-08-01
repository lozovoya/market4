package repository

import (
	"context"
	"market4/internal/model"
)

type Shop interface {
	ListAllShops(ctx context.Context) ([]*model.Shop, error)
	AddShop(ctx context.Context, s *model.Shop) (int, error)
	EditShop(ctx context.Context, s *model.Shop) error
	IfShopExists(ctx context.Context, shopID int) bool
}

type Category interface {
	ListAllCategories(ctx context.Context) ([]*model.Category, error)
	AddCategory(ctx context.Context, c *model.Category) (int, error)
	EditCategory(ctx context.Context, c *model.Category) error
	IfCategoryExists(ctx context.Context, categoryID int) bool
}

type Product interface {
	AddProduct(ctx context.Context, p *model.Product, shopId int, categoryId int) (*model.Product, error)
	EditProduct(ctx context.Context, p *model.Product, shopId int, categoryId int) (*model.Product, error)
	ListAllProducts(ctx context.Context) ([]*model.Product, error)
	IfProductExists(ctx context.Context, productID string) bool
}

type Price interface {
	AddPrice(ctx context.Context, p *model.Price, productID string) (int, error)
	EditPrice(ctx context.Context, p *model.Price, productID string) (*model.Price, error)
	ListAllPrices(ctx context.Context) ([]*model.Price, error)
}
