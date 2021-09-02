package httpserver

import (
	"log"
	"market4/internal/api/httpserver/md"
	v1 "market4/internal/api/v1"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(
	mux *chi.Mux,
	shopController *v1.Shop,
	categoryController *v1.Category,
	productController *v1.Product,
	priceController *v1.Price,
	usersController *v1.Users) chi.Mux {
	mux.Use(middleware.Logger)
	mux.Route("/api/v1", func(router chi.Router) {
		router.Get("/shops", shopController.ListAllShops)
		router.Post("/shops", shopController.AddShop)
		router.Put("/shops", shopController.EditShop)

		router.Get("/categories", categoryController.ListAllCategories)
		router.Post("/categories", categoryController.AddCategory)
		router.Put("/categories", categoryController.EditCategory)

		router.With(md.Auth("ADMIN")).Post("/products", productController.AddProduct)
		router.Put("/products", productController.EditProduct)
		router.Get("/products", productController.ListAllProducts)

		router.Post("/prices", priceController.AddPrice)
		router.Put("/prices", priceController.EditPrice)
		router.Get("/prices", priceController.ListAllPrices)

		router.Get("/categories/{categoryID:.+}/products", productController.SearchProductsByCategory)
		router.Get("/search/{product_name:.+}", productController.SearchProductByName)
		router.Get("/shops/{shopID:.+}/products", productController.SearchActiveProductsOfShop)

		router.Post("/users", usersController.AddUser)
		router.Put("/users", usersController.EditUser)
		router.Get("/users/token", usersController.Token)
	})

	log.Println("new router is activated")
	return *mux
}
