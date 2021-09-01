package views

import (
	"market4/internal/model"
)

func ShopList(shops []*model.Shop) (*ShopListDTO, error) {
	if len(shops) == 0 {
		return nil, nil
	}

	var shopList ShopListDTO
	shopList.Total = len(shops)

	for _, shop := range shops {
		var item model.Shop

		item.ID = shop.ID
		item.Name = shop.Name
		item.Address = shop.Address
		item.LON = shop.LON
		item.LAT = shop.LAT
		item.WorkingHours = shop.WorkingHours

		shopList.Items = append(shopList.Items, &item)
	}
	return &shopList, nil
}
