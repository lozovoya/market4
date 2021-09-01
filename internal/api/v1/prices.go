package v1

import (
	"encoding/json"
	"fmt"
	"log"
	"market4/internal/model"
	"market4/internal/repository"
	"market4/internal/views"
	"net/http"
)

type PriceDTO struct {
	ID            int    `json:"id,omitempty,string"`
	SalePrice     int    `json:"sale_price,string"`
	FactoryPrice  int    `json:"factory_price,string"`
	DiscountPrice int    `json:"discount_price,string"`
	IsActive      bool   `json:"is_active,string"`
	ProductID     string `json:"product_id,omitempty"`
}

type Price struct {
	priceRepo repository.Price
}

func NewPrice(priceRepo repository.Price) *Price {
	return &Price{priceRepo: priceRepo}
}

func (price *Price) AddPrice(writer http.ResponseWriter, request *http.Request) {
	var data *PriceDTO
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(fmt.Errorf("addPrice: %w", err))
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var p = model.Price{
		SalePrice:     data.SalePrice,
		FactoryPrice:  data.FactoryPrice,
		DiscountPrice: data.DiscountPrice,
		IsActive:      data.IsActive,
		ProductID:     data.ProductID,
	}
	addedPrice, err := price.priceRepo.AddPrice(request.Context(), &p)
	if err != nil {
		log.Println(fmt.Errorf("AddPrice: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(addedPrice)
	if err != nil {
		log.Println(fmt.Errorf("editProduct: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (price *Price) EditPrice(writer http.ResponseWriter, request *http.Request) {
	var data *PriceDTO
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(fmt.Errorf("addPrice: %w", err))
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if data.ID == 0 {
		log.Println(fmt.Errorf("EditPrice: id is empty"))
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var p = model.Price{
		ID:            data.ID,
		SalePrice:     data.SalePrice,
		FactoryPrice:  data.FactoryPrice,
		DiscountPrice: data.DiscountPrice,
		IsActive:      data.IsActive,
		ProductID:     data.ProductID,
	}
	editedPrice, err := price.priceRepo.EditPrice(request.Context(), &p)
	if err != nil {
		log.Println(fmt.Errorf("EditPrice: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if editedPrice == nil {
		return
	}
	var priceList = make([]*model.Price, 0)
	priceList = append(priceList, editedPrice)
	result, err := views.PricesList(priceList)
	if err != nil {
		log.Println(fmt.Errorf("EditPrice: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(result)
	if err != nil {
		log.Println(fmt.Errorf("editPrice: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (price *Price) ListAllPrices(writer http.ResponseWriter, request *http.Request) {
	prices, err := price.priceRepo.ListAllPrices(request.Context())
	if err != nil {
		log.Println(fmt.Errorf("ListAllPrices: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	priceList, err := views.PricesList(prices)
	if err != nil {
		log.Println(fmt.Errorf("ListAllPrices: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(priceList)
	if err != nil {
		log.Println(fmt.Errorf("ListAllPrices: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
