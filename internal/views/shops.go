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

		item.Id = shop.Id
		item.Name = shop.Name
		item.Address = shop.Address
		item.Lon = shop.Lon
		item.Lat = shop.Lat
		item.WorkingHours = shop.WorkingHours

		shopList.Items = append(shopList.Items, &item)
	}
	return &shopList, nil
}
