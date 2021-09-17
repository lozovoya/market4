package model

type Product struct {
	ID          string `json:"id,omitempty"`
	SKU         string `json:"sku,omitempty"`
	Name        string `json:"name,omitempty"`
	Type        string `json:"type,omitempty"`
	URI         string `json:"uri,omitempty"`
	Description string `json:"description,omitempty"`
	IsActive    bool   `json:"is_active"`
}
