package views

import (
	"errors"
	"fmt"
	"market4/internal/model"
)

func PricesList(prices []*model.Price) (*PricesListDTO, error) {
	if len(prices) == 0 {
		err := errors.New("Price list is empty")
		return nil, fmt.Errorf("PricesList: %w", err)
	}

	var pricesList PricesListDTO
	pricesList.Total = len(prices)
	for _, price := range prices {
		var item model.Price

		item.ID = price.ID
		item.SalePrice = price.SalePrice
		item.FactoryPrice = price.FactoryPrice
		item.DiscountPrice = price.DiscountPrice
		item.IsActive = price.IsActive
		pricesList.Items = append(pricesList.Items, &item)
	}

	return &pricesList, nil
}
