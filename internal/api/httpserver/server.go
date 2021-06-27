package httpserver

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"market4/internal/api/v1"
)

func NewRouter(mux chi.Mux, controller controllers.MarketController) chi.Mux {
	mux.Use(middleware.Logger)
	mux.Route("/api/v1", func(router chi.Router) {
		router.Get("/shops", controller.ListAllShops)
		router.Post("/shops", controller.AddShop)
		router.Put("/shops", controller.EditShop)

		//router.Get("/categories", controller.ListAllCategories)
		//router.Post("/categories", controller.AddCategory)
		//router.Put("/categories", controller.EditCategory)

	})

	log.Println("new router is activated")
	return mux
}
