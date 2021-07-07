package v1

import (
	"encoding/json"
	"fmt"
	"log"
	"market4/internal/model"
	"market4/internal/repository"
	"net/http"
)

type ProductDTO struct {
	SKU         string `json:"sku"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Shop_ID     int    `json:"shop_id,string"`
	Category_ID int    `json:"category_id,string"`
}

type Product struct {
	productRepo repository.Product
}

func NewProduct(productRepo repository.Product) *Product {
	return &Product{productRepo: productRepo}
}

func (p *Product) AddProduct(writer http.ResponseWriter, request *http.Request) {

	var data *ProductDTO
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(fmt.Errorf("addProduct: %w", err))
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if IsEmpty(data.SKU) || IsEmpty(data.Name) || IsEmpty(data.Type) || IsEmpty(data.Description) {
		log.Println("field is empty")
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var product = model.Product{
		SKU:         data.SKU,
		Name:        data.Name,
		Type:        data.Type,
		Description: data.Description,
	}

	id, err := p.productRepo.AddProduct(request.Context(), &product, data.Shop_ID, data.Category_ID)
	if err != nil {
		log.Println(fmt.Errorf("addProduct: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	var reply struct {
		Id string `json:"id"`
	}
	reply.Id = id

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(reply)
	if err != nil {
		log.Println(fmt.Errorf("addCategory: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
