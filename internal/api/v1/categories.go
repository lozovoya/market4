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

type categoryDTO struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name"`
	URI_name string `json:"uri_name"`
}

type Category struct {
	categoryRepo repository.Category
}

func NewCategory(categoryRepo repository.Category) *Category {
	return &Category{categoryRepo: categoryRepo}
}

func (c *Category) ListAllCategories(writer http.ResponseWriter, request *http.Request) {

	categories, err := c.categoryRepo.ListAllCategories(request.Context())
	if err != nil {
		log.Println(fmt.Errorf("ListAllCategories: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	categoriesList, err := views.CategoriesList(categories)
	if err != nil {
		log.Println(fmt.Errorf("ListAllCategories: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(categoriesList)
	if err != nil {
		log.Println(fmt.Errorf("ListAllCategories: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (c *Category) AddCategory(writer http.ResponseWriter, request *http.Request) {
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

	id, err := c.categoryRepo.AddCategory(request.Context(), data)
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

func (c *Category) EditCategory(writer http.ResponseWriter, request *http.Request) {

	var data *model.Category
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(fmt.Errorf("editCategory: %w", err))
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	//todo проверить id на пустоту
	if IsEmpty(data.Name) {
		log.Println("field is empty")
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = c.categoryRepo.EditCategory(request.Context(), data)
	if err != nil {
		log.Println(fmt.Errorf("editCategory: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
