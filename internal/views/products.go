package views

import (
	"errors"
	"fmt"
	"market4/internal/model"
)

func ProductsList(products []*model.Product) (*ProductsListDTO, error) {

	if len(products) == 0 {
		err := errors.New("Product list is empty")
		return nil, fmt.Errorf("ProductList: %w", err)
	}

	var productsList ProductsListDTO
	productsList.Total = len(products)
	for _, product := range products {
		var item Product

		item.ID = product.ID
		item.SKU = product.SKU
		item.Name = product.Name
		item.URI = product.URI
		item.Description = product.Description
		item.IsActive = product.IsActive
		productsList.Items = append(productsList.Items, &item)
	}

	return &productsList, nil
}

func ProductsListWithPrices(products []*model.Product, prices []*model.Price) (*ProductsListDTO, error) {

	if len(products) == 0 {
		err := errors.New("Product list is empty")
		return nil, fmt.Errorf("ProductList: %w", err)
	}

	var productsList ProductsListDTO
	productsList.Total = len(products)
	for _, product := range products {
		var item Product

		item.ID = product.ID
		item.SKU = product.SKU
		item.Name = product.Name
		item.URI = product.URI
		item.Description = product.Description
		item.IsActive = product.IsActive

		for _, price := range prices {
			var itemPrice Price
			if item.ID == price.ProductID {
				itemPrice.SalePrice = price.SalePrice
				itemPrice.FactoryPrice = price.FactoryPrice
				itemPrice.DiscountPrice = price.DiscountPrice

				item.Prices = append(item.Prices, &itemPrice)
				break
			}
		}
		productsList.Items = append(productsList.Items, &item)
	}
	return &productsList, nil
}
