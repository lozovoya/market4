package httpserver

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"market4/internal/api/v1"
)

func NewRouter(
	mux chi.Mux,
	shopController *v1.Shop,
	categoryController *v1.Category,
	productController *v1.Product) chi.Mux {

	mux.Use(middleware.Logger)
	mux.Route("/api/v1", func(router chi.Router) {
		router.Get("/shops", shopController.ListAllShops)
		router.Post("/shops", shopController.AddShop)
		router.Put("/shops", shopController.EditShop)

		router.Get("/categories", categoryController.ListAllCategories)
		router.Post("/categories", categoryController.AddCategory)
		router.Put("/categories", categoryController.EditCategory)

		router.Post("/products", productController.AddProduct)

	})

	log.Println("new router is activated")
	return mux
}
