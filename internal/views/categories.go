package views

import (
	"market4/internal/model"
)

func CategoriesList(categories []model.Category) (*CategoriesListDTO, error) {
	if len(categories) == 0 {
		return nil, nil
	}

	var categoriesList CategoriesListDTO
	categoriesList.Total = len(categories)
	for _, category := range categories {
		var item model.Category

		item.ID = category.ID
		item.Name = category.Name
		item.URI_name = category.URI_name

		categoriesList.Items = append(categoriesList.Items, &item)
	}
	return &categoriesList, nil
}
