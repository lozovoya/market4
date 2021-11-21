package v1

import (
	"encoding/json"
	"go.uber.org/zap"
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
	lg       *zap.Logger
}

func NewShop(shopRepo repository.Shop, lg *zap.Logger) *Shop {
	return &Shop{shopRepo: shopRepo, lg: lg}
}

func (s *Shop) EditShop(writer http.ResponseWriter, request *http.Request) {
	var data *model.Shop
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		s.lg.Error("editShop", zap.Error(err))
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if data.ID == 0 {
		s.lg.Error("editShop", zap.Error(err))
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = s.shopRepo.EditShop(request.Context(), data)
	if err != nil {
		s.lg.Error("editShop", zap.Error(err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
func (s *Shop) ListAllShops(writer http.ResponseWriter, request *http.Request) {
	shops, err := s.shopRepo.ListAllShops(request.Context())
	if err != nil {
		s.lg.Error("ListAllShops", zap.Error(err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	shopList, err := views.ShopList(&shops)
	if err != nil {
		s.lg.Error("ListAllShops", zap.Error(err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(shopList)
	if err != nil {
		s.lg.Error("ListAllShops", zap.Error(err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
func (s *Shop) AddShop(writer http.ResponseWriter, request *http.Request) {
	var data *model.Shop
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		s.lg.Error("AddShop", zap.Error(err))
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id, err := s.shopRepo.AddShop(request.Context(), data)
	if err != nil {
		s.lg.Error("AddShop", zap.Error(err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var reply struct {
		Id int `json:"id,string"`
	}
	reply.Id = id

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(reply)
	if err != nil {
		s.lg.Error("AddShop", zap.Error(err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
