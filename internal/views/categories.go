package views

import (
	"market4/internal/model"
)

func CategoriesList(categories []*model.Category) (*CategoriesListDTO, error) {

	if len(categories) == 0 {
		return nil, nil
	}

	var categoriesList CategoriesListDTO
	categoriesList.Total = len(categories)
	for _, category := range categories {
		var item model.Category

		item.Id = category.Id
		item.Name = category.Name
		item.Uri_name = category.Uri_name

		categoriesList.Items = append(categoriesList.Items, &item)
	}
	return &categoriesList, nil
}
