package model

type Product struct {
	ID          string `json:"id,omitempty"`
	SKU         string `json:"sku"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	URI         string `json:"uri"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}
