package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"market4/internal/model"
	"market4/internal/views"
	"net/http"
)

type categoryDTO struct {
	Id       string `json:"id,omitempty"`
	Name     string `json:"name"`
	Uri_name string `json:"uri_name"`
}

func (m *marketController) ListAllCategories(writer http.ResponseWriter, request *http.Request) {

	categories, err := m.repo.ListAllCategories(request.Context())
	if err != nil {
		log.Println(fmt.Errorf("getAllCategories: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	categoriesList, err := views.CategoriesList(categories)

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(categoriesList)
}

func (m *marketController) AddCategory(writer http.ResponseWriter, request *http.Request) {
	var data *model.Category
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(fmt.Errorf("addCategory: %w", err))
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if IsEmpty(data.Name) {
		log.Println("field is empty")
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id, err := m.repo.AddCategory(request.Context(), data)
	if err != nil {
		log.Println(fmt.Errorf("addCategory: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	var reply struct {
		Id int `json:"id,string"`
	}
	reply.Id = id

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(reply)
	if err != nil {
		log.Println(fmt.Errorf("addCategory: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (m *marketController) EditCategory(writer http.ResponseWriter, request *http.Request) {

	var data *model.Category
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(fmt.Errorf("editCategory: %w", err))
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = m.repo.EditCategory(request.Context(), data)
	if err != nil {
		log.Println(fmt.Errorf("editCategory: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
