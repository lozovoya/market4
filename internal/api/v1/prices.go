package v1

import (
	"encoding/json"
	"market4/internal/model"
	"market4/internal/repository"
	"market4/internal/views"
	"net/http"

	"github.com/unrolled/render"
	"go.uber.org/zap"
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
	lg        *zap.Logger
	renderer  *render.Render
}

func NewPrice(priceRepo repository.Price, lg *zap.Logger, renderer *render.Render) *Price {
	return &Price{priceRepo: priceRepo, lg: lg, renderer: renderer}
}

func (price *Price) AddPrice(writer http.ResponseWriter, request *http.Request) {
	var data *PriceDTO
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		price.lg.Error("addPrice", zap.Error(err))
		err = price.renderer.JSON(writer, http.StatusBadRequest, map[string]string{"Error": "BadRequest"})
		if err != nil {
			price.lg.Error("Auth", zap.Error(err))
		}
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
		price.lg.Error("addPrice", zap.Error(err))
		err = price.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			price.lg.Error("Auth", zap.Error(err))
		}
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(addedPrice)
	if err != nil {
		price.lg.Error("addPrice", zap.Error(err))
		err = price.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			price.lg.Error("Auth", zap.Error(err))
		}
		return
	}
}

func (price *Price) EditPrice(writer http.ResponseWriter, request *http.Request) {
	var data *PriceDTO
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		price.lg.Error("EditPrice", zap.Error(err))
		err = price.renderer.JSON(writer, http.StatusBadRequest, map[string]string{"Error": "BadRequest"})
		if err != nil {
			price.lg.Error("Auth", zap.Error(err))
		}
		return
	}

	if data.ID == 0 {
		price.lg.Error("EditPrice: id is empty")
		err = price.renderer.JSON(writer, http.StatusBadRequest, map[string]string{"Error": "BadRequest"})
		if err != nil {
			price.lg.Error("Auth", zap.Error(err))
		}
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
		price.lg.Error("EditPrice", zap.Error(err))
		err = price.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			price.lg.Error("Auth", zap.Error(err))
		}
		return
	}

	if editedPrice.ID == 0 {
		return
	}
	var priceList = make([]model.Price, 0)
	priceList = append(priceList, editedPrice)
	result, err := views.MakePricesList(priceList)
	if err != nil {
		price.lg.Error("EditPrice", zap.Error(err))
		err = price.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			price.lg.Error("Auth", zap.Error(err))
		}
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(result)
	if err != nil {
		price.lg.Error("EditPrice", zap.Error(err))
		err = price.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			price.lg.Error("Auth", zap.Error(err))
		}
		return
	}
}

func (price *Price) ListAllPrices(writer http.ResponseWriter, request *http.Request) {
	prices, err := price.priceRepo.ListAllPrices(request.Context())
	if err != nil {
		price.lg.Error("ListAllPrices", zap.Error(err))
		err = price.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			price.lg.Error("Auth", zap.Error(err))
		}
		return
	}

	priceList, err := views.MakePricesList(prices)
	if err != nil {
		price.lg.Error("ListAllPrices", zap.Error(err))
		err = price.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			price.lg.Error("Auth", zap.Error(err))
		}
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(priceList)
	if err != nil {
		price.lg.Error("ListAllPrices", zap.Error(err))
		err = price.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			price.lg.Error("Auth", zap.Error(err))
		}
		return
	}
}
