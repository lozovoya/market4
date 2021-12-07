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

type Category struct {
	categoryRepo repository.Category
	lg           *zap.Logger
	renderer     *render.Render
}

func NewCategory(categoryRepo repository.Category, lg *zap.Logger, renderer *render.Render) *Category {
	return &Category{categoryRepo: categoryRepo, lg: lg, renderer: renderer}
}
func (c *Category) ListAllCategories(writer http.ResponseWriter, request *http.Request) {
	categories, err := c.categoryRepo.ListAllCategories(request.Context())
	if err != nil {
		c.lg.Error("ListAllCategories", zap.Error(err))
		err = c.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			c.lg.Error("Auth", zap.Error(err))
		}
		return
	}

	categoriesList, err := views.MakeCategoriesList(categories)
	if err != nil {
		c.lg.Error("ListAllCategories", zap.Error(err))
		err = c.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			c.lg.Error("Auth", zap.Error(err))
		}
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(categoriesList)
	if err != nil {
		c.lg.Error("ListAllCategories", zap.Error(err))
		err = c.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			c.lg.Error("Auth", zap.Error(err))
		}
		return
	}
}

func (c *Category) AddCategory(writer http.ResponseWriter, request *http.Request) {
	var data *model.Category
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		c.lg.Error("AddCategory", zap.Error(err))
		err = c.renderer.JSON(writer, http.StatusBadRequest, map[string]string{"Error": "BadRequest"})
		if err != nil {
			c.lg.Error("Auth", zap.Error(err))
		}
		return
	}

	if IsEmpty(data.Name) {
		c.lg.Error("field is empty")
		err = c.renderer.JSON(writer, http.StatusBadRequest, map[string]string{"Error": "BadRequest"})
		if err != nil {
			c.lg.Error("Auth", zap.Error(err))
		}
		return
	}

	id, err := c.categoryRepo.AddCategory(request.Context(), data)
	if err != nil {
		c.lg.Error("AddCategory", zap.Error(err))
		err = c.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			c.lg.Error("Auth", zap.Error(err))
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
		c.lg.Error("AddCategory", zap.Error(err))
		err = c.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			c.lg.Error("Auth", zap.Error(err))
		}
		return
	}
}
func (c *Category) EditCategory(writer http.ResponseWriter, request *http.Request) {
	var data *model.Category
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		c.lg.Error("editCategory", zap.Error(err))
		err = c.renderer.JSON(writer, http.StatusBadRequest, map[string]string{"Error": "BadRequest"})
		if err != nil {
			c.lg.Error("Auth", zap.Error(err))
		}
		return
	}
	if data.ID == 0 {
		c.lg.Error("wrong ID")
		err = c.renderer.JSON(writer, http.StatusBadRequest, map[string]string{"Error": "BadRequest"})
		if err != nil {
			c.lg.Error("Auth", zap.Error(err))
		}
		return
	}

	if IsEmpty(data.Name) {
		c.lg.Error("field is empty")
		err = c.renderer.JSON(writer, http.StatusBadRequest, map[string]string{"Error": "BadRequest"})
		if err != nil {
			c.lg.Error("Auth", zap.Error(err))
		}
		return
	}

	err = c.categoryRepo.EditCategory(request.Context(), data)
	if err != nil {
		c.lg.Error("editCategory", zap.Error(err))
		err = c.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			c.lg.Error("Auth", zap.Error(err))
		}
		return
	}
}
