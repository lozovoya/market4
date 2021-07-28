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

type ProductsListDTO struct {
	Total int              `json:"total"`
	Items []*model.Product `json:"items"`
}

type PricesListDTO struct {
	Total int            `json:"total"`
	Items []*model.Price `json:"items"`
}
