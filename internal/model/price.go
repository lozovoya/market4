package model

type Price struct {
	ID            int    `json:"id,omitempty"`
	SalePrice     int    `json:"sale_price"`
	FactoryPrice  int    `json:"factory_price"`
	DiscountPrice int    `json:"discount_price"`
	IsActive      bool   `json:"is_active,omitempty"`
	ProductID     string `json:"product_id"`
}
