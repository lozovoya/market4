package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"market4/internal/model"
	"market4/internal/views"
	"net/http"
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
	var data *model.Shop
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

	err = m.repo.EditShop(request.Context(), data)
	if err != nil {
		log.Println(fmt.Errorf("editShop: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (m *marketController) ListAllShops(writer http.ResponseWriter, request *http.Request) {

	log.Println("list all shops")
	shops, err := m.repo.ListAllShops(request.Context())
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

	id, err := m.repo.AddShop(request.Context(), data)

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
