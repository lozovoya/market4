package v1

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"market4/internal/model"
	"market4/internal/repository"
	"market4/internal/views"
	"net/http"
)

type ProductDTO struct {
	SKU         string      `json:"sku"`
	Name        string      `json:"name,omitempty"`
	Type        string      `json:"type,omitempty"`
	Description string      `json:"description,omitempty"`
	IsActive    bool        `json:"is_active,string,omitempty"`
	Shop_ID     int         `json:"shop_id,string,omitempty"`
	Category_ID int         `json:"category_id,string,omitempty"`
	Prices      []*PriceDTO `json:"prices,omitempty"`
}

type Product struct {
	productRepo repository.Product
	priceRepo   repository.Price
}

func NewProduct(productRepo repository.Product, priceRepo repository.Price) *Product {
	return &Product{productRepo: productRepo, priceRepo: priceRepo}
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

	addedProduct, err := p.productRepo.AddProduct(request.Context(), &product, data.Shop_ID, data.Category_ID)
	if err != nil {
		log.Println(fmt.Errorf("addProduct: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(addedProduct)
	if err != nil {
		log.Println(fmt.Errorf("addCategory: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (p *Product) EditProduct(writer http.ResponseWriter, request *http.Request) {
	var data *ProductDTO
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(fmt.Errorf("editProduct: %w", err))
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if IsEmpty(data.SKU) {
		log.Println(fmt.Errorf("editProduct: SKU field is empty"))
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var shopID, categoryID = 0, 0

	var product = model.Product{
		SKU:         data.SKU,
		Name:        data.Name,
		Type:        data.Type,
		Description: data.Description,
		IsActive:    data.IsActive,
	}

	shopID = data.Shop_ID
	categoryID = data.Category_ID

	editedProduct, err := p.productRepo.EditProduct(request.Context(), &product, shopID, categoryID)
	if err != nil {
		log.Println(fmt.Errorf("editProduct: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(editedProduct)
	if err != nil {
		log.Println(fmt.Errorf("editProduct: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

}

func (p *Product) ListAllProducts(writer http.ResponseWriter, request *http.Request) {
	products, err := p.productRepo.ListAllProducts(request.Context())
	if err != nil {
		log.Println(fmt.Errorf("ListAllProducts: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	prices, err := p.priceRepo.ListAllPrices(request.Context())
	if err != nil {
		log.Println(fmt.Errorf("ListAllPrices: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	productsList, err := views.ProductsList(products, prices)
	if err != nil {
		log.Println(fmt.Errorf("ListAllProducts: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(productsList)
	if err != nil {
		log.Println(fmt.Errorf("ListAllProducts: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (p *Product) SearchProductsByCategory(writer http.ResponseWriter, request *http.Request) {

	log.Println("search by categories")
	log.Printf("values %s", request.RequestURI)
	log.Printf("values %s", chi.URLParam(request, "categoryID"))
}
