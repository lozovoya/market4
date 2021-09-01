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

type ShopDTO struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Address      string `json:"address"`
	WorkingHours string `json:"working_hours"`
	LON          string `json:"lon"`
	LAT          string `json:"lat"`
}

type Shop struct {
	shopRepo repository.Shop
}

func NewShop(shopRepo repository.Shop) *Shop {
	return &Shop{shopRepo: shopRepo}
}

func (s *Shop) EditShop(writer http.ResponseWriter, request *http.Request) {
	var data *model.Shop
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(fmt.Errorf("editShop: %w", err))
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if data.ID == 0 {
		log.Println(fmt.Errorf("editShop: %w", err))
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = s.shopRepo.EditShop(request.Context(), data)
	if err != nil {
		log.Println(fmt.Errorf("editShop: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
func (s *Shop) ListAllShops(writer http.ResponseWriter, request *http.Request) {
	shops, err := s.shopRepo.ListAllShops(request.Context())
	if err != nil {
		log.Println(fmt.Errorf("ListAllShops: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	shopList, err := views.ShopList(shops)
	if err != nil {
		log.Println(fmt.Errorf("ListAllShops: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(shopList)
	if err != nil {
		log.Println(fmt.Errorf("ListAllShops: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
func (s *Shop) AddShop(writer http.ResponseWriter, request *http.Request) {
	var data *model.Shop
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(fmt.Errorf("addShop: %w", err))
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id, err := s.shopRepo.AddShop(request.Context(), data)
	if err != nil {
		log.Println(fmt.Errorf("addShop: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	var reply struct {
		Id int `json:"id,string"`
	}
	reply.Id = id

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(reply)
	if err != nil {
		log.Println(fmt.Errorf("addShop: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
