package views

import "market4/internal/model"

type ShopListDTO struct {
	Total int           `json:"total"`
	Items []*model.Shop `json:"items"`
}

type CategoriesListDTO struct {
	Total int               `json:"total"`
	Items []*model.Category `json:"items"`
}

type Price struct {
	SalePrice     int `json:"sale_price"`
	FactoryPrice  int `json:"factory_price"`
	DiscountPrice int `json:"discount_price"`
}

type Product struct {
	ID          string   `json:"id,omitempty"`
	SKU         string   `json:"sku"`
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	URI         string   `json:"uri"`
	Description string   `json:"description"`
	IsActive    bool     `json:"is_active"`
	Prices      []*Price `json:"prices"`
}

type ProductsListDTO struct {
	Total int        `json:"total"`
	Items []*Product `json:"items"`
}

type PricesListDTO struct {
	Total int            `json:"total"`
	Items []*model.Price `json:"items"`
}
