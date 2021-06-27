package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"market4/internal/model"
	"market4/internal/views"
	"net/http"
	"strconv"
)

type ShopDTO struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Address      string `json:"address"`
	WorkingHours string `json:"working_hours"`
	Lon          string `json:"lon"`
	Lat          string `json:"lat"`
}

func (m *marketController) EditShop(writer http.ResponseWriter, request *http.Request) {
	var data ShopDTO
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(fmt.Errorf("addShop: %w", err))
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err != nil {
		log.Println(fmt.Errorf("editShop: %w", err))
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(data.Id)
	if err != nil {
		log.Println(fmt.Errorf("editShop: %w", err))
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	shop := model.Shop{
		Id:           id,
		Name:         data.Name,
		Address:      data.Address,
		Lon:          data.Lon,
		Lat:          data.Lat,
		WorkingHours: data.WorkingHours,
	}

	err = m.market.EditShop(request.Context(), &shop)
	if err != nil {
		log.Println(fmt.Errorf("editShop: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (m *marketController) ListAllShops(writer http.ResponseWriter, request *http.Request) {

	log.Println("list all shops")
	shops, err := m.market.ListAllShops(request.Context())
	if err != nil {
		return
	}

	shopList, err := views.ShopList(shops)
	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(shopList)
}

func (m *marketController) AddShop(writer http.ResponseWriter, request *http.Request) {

	var data *model.Shop
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(fmt.Errorf("addShop: %w", err))
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	shop := model.Shop{
		Name:         data.Name,
		Address:      data.Address,
		Lon:          data.Lon,
		Lat:          data.Lat,
		WorkingHours: data.WorkingHours,
	}

	id, err := m.market.AddShop(request.Context(), &shop)

	if err != nil {
		log.Println(fmt.Errorf("addShop: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	var reply struct {
		Id string `json:"id"`
	}
	reply.Id = strconv.Itoa(id)

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(reply)
	if err != nil {
		log.Println(fmt.Errorf("addShop: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
