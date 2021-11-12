package v1

import (
	"encoding/json"
	"go.uber.org/zap"
	"market4/internal/model"
	"market4/internal/repository"
	"market4/internal/views"
	"net/http"
)

type Category struct {
	categoryRepo repository.Category
	lg *zap.Logger
}

func NewCategory(categoryRepo repository.Category, lg *zap.Logger) *Category {
	return &Category{categoryRepo: categoryRepo, lg: lg}
}
func (c *Category) ListAllCategories(writer http.ResponseWriter, request *http.Request) {
	categories, err := c.categoryRepo.ListAllCategories(request.Context())
	if err != nil {
		c.lg.Error("ListAllCategories", zap.Error(err))

		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		http.
		return
	}

	categoriesList, err := views.CategoriesList(categories)
	if err != nil {
		c.lg.Error("ListAllCategories", zap.Error(err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(categoriesList)
	if err != nil {
		c.lg.Error("ListAllCategories", zap.Error(err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (c *Category) AddCategory(writer http.ResponseWriter, request *http.Request) {
	var data *model.Category
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		c.lg.Error("AddCategory", zap.Error(err))
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if IsEmpty(data.Name) {
		c.lg.Error("field is empty")
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id, err := c.categoryRepo.AddCategory(request.Context(), data)
	if err != nil {
		c.lg.Error("AddCategory", zap.Error(err))
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
		c.lg.Error("AddCategory", zap.Error(err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
func (c *Category) EditCategory(writer http.ResponseWriter, request *http.Request) {
	var data *model.Category
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		c.lg.Error("editCategory", zap.Error(err))
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if data.ID == 0 {
		c.lg.Error("wrong ID")
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if IsEmpty(data.Name) {
		c.lg.Error("field is empty")
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = c.categoryRepo.EditCategory(request.Context(), data)
	if err != nil {
		c.lg.Error("editCategory", zap.Error(err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
