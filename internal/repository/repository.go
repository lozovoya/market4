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
	SearchProductsByCategory(ctx context.Context, category int) ([]*model.Product, error)
	SearchProductsByName(ctx context.Context, productName string) (*model.Product, error)
	SearchProductsByShop(ctx context.Context, shopID int) ([]*model.Product, error)
}

type Price interface {
	AddPrice(ctx context.Context, p *model.Price) (*model.Price, error)
	EditPrice(ctx context.Context, p *model.Price) (*model.Price, error)
	ListAllPrices(ctx context.Context) ([]*model.Price, error)
	SearchPriceByProductID(ctx context.Context, productID string) (*model.Price, error)
	EditPriceByProductID(ctx context.Context, p *model.Price) (*model.Price, error)
}

type Users interface {
	AddUser(ctx context.Context, u *model.User) (*model.User, error)
	EditUser(ctx context.Context, u *model.User) (*model.User, error)
	GetUserRole(ctx context.Context, login string) (int, error)
	CheckCreds(ctx context.Context, u *model.User) bool
	GetUserID(ctx context.Context, login string) (int, error)
	GetRoleByID(ctx context.Context, roleID int) (string, error)
}
