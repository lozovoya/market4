package v1

import (
	"encoding/json"
	"github.com/unrolled/render"
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
	renderer *render.Render
}

func NewShop(shopRepo repository.Shop, lg *zap.Logger, renderer *render.Render) *Shop {
	return &Shop{shopRepo: shopRepo, lg: lg, renderer: renderer}
}

func (s *Shop) EditShop(writer http.ResponseWriter, request *http.Request) {
	var data *model.Shop
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		s.lg.Error("editShop", zap.Error(err))
		err = s.renderer.JSON(writer, http.StatusBadRequest, map[string]string{"Error": "BadRequest"})
		if err != nil {
			s.lg.Error("Auth", zap.Error(err))
		}
		return
	}
	if data.ID == 0 {
		s.lg.Error("editShop", zap.Error(err))
		err = s.renderer.JSON(writer, http.StatusBadRequest, map[string]string{"Error": "BadRequest"})
		if err != nil {
			s.lg.Error("Auth", zap.Error(err))
		}
		return
	}

	err = s.shopRepo.EditShop(request.Context(), data)
	if err != nil {
		s.lg.Error("editShop", zap.Error(err))
		err = s.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			s.lg.Error("Auth", zap.Error(err))
		}
		return
	}
}
func (s *Shop) ListAllShops(writer http.ResponseWriter, request *http.Request) {
	shops, err := s.shopRepo.ListAllShops(request.Context())
	if err != nil {
		s.lg.Error("ListAllShops", zap.Error(err))
		err = s.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			s.lg.Error("Auth", zap.Error(err))
		}
		return
	}

	shopList, err := views.MakeShopList(&shops)
	if err != nil {
		s.lg.Error("ListAllShops", zap.Error(err))
		err = s.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			s.lg.Error("Auth", zap.Error(err))
		}
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(shopList)
	if err != nil {
		s.lg.Error("ListAllShops", zap.Error(err))
		err = s.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			s.lg.Error("Auth", zap.Error(err))
		}
		return
	}
}
func (s *Shop) AddShop(writer http.ResponseWriter, request *http.Request) {
	var data *model.Shop
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		s.lg.Error("AddShop", zap.Error(err))
		err = s.renderer.JSON(writer, http.StatusBadRequest, map[string]string{"Error": "BadRequest"})
		if err != nil {
			s.lg.Error("Auth", zap.Error(err))
		}
		return
	}

	id, err := s.shopRepo.AddShop(request.Context(), data)
	if err != nil {
		s.lg.Error("AddShop", zap.Error(err))
		err = s.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			s.lg.Error("Auth", zap.Error(err))
		}
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
		err = s.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			s.lg.Error("Auth", zap.Error(err))
		}
		return
	}
}
