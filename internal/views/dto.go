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
	SKU         string   `json:"sku,omitempty"`
	Name        string   `json:"name,omitempty"`
	Type        string   `json:"type,omitempty"`
	URI         string   `json:"uri,omitempty"`
	Description string   `json:"description,omitempty"`
	IsActive    bool     `json:"is_active,omitempty"`
	Prices      []*Price `json:"prices,omitempty"`
}

type ProductsListDTO struct {
	Total int        `json:"total"`
	Items []*Product `json:"items"`
}

type PricesListDTO struct {
	Total int            `json:"total"`
	Items []*model.Price `json:"items"`
}
