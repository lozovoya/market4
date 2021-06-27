package controllers

import (
	"market4/internal/repository"
	"net/http"
)

type marketController struct {
	market  repository.MarketRepository
	storage repository.Storage
}

type MarketController interface {
	ListAllShops(writer http.ResponseWriter, request *http.Request)
	AddShop(writer http.ResponseWriter, request *http.Request)
	EditShop(writer http.ResponseWriter, request *http.Request)

	//ListAllCategories(writer http.ResponseWriter, request *http.Request)
	//AddCategory(writer http.ResponseWriter, request *http.Request)
	//EditCategory(writer http.ResponseWriter, request *http.Request)
	//
	//ListAllProducts(writer http.ResponseWriter, request *http.Request)
	//AddProduct(writer http.ResponseWriter, request *http.Request)
	//EditProduct(writer http.ResponseWriter, request *http.Request)
}

func NewMarketController(storage repository.Storage) MarketController {
	return &marketController{storage: storage}
}

func (m *marketController) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	m.ServeHTTP(writer, request)
}
