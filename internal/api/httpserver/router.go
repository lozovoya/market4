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
	usersController *v1.Users,
	authController *v1.Auth) chi.Mux {
	mux.Use(middleware.Logger)
	mux.Route("/api/v1", func(router chi.Router) {
		RouterShop(router, shopController)
		RouterCategories(router, categoryController)
		RouterProduct(router, productController)
		RouterPrice(router, priceController)
		RouterUser(router, usersController)
		RouterAuth(router, authController)
	})

	log.Println("new router is activated")
	return *mux
}

func RouterShop(router chi.Router, shopController *v1.Shop) chi.Router {
	router.With(md.Auth("USER")).Get("/shops", shopController.ListAllShops)
	router.With(md.Auth("ADMIN")).Post("/shops", shopController.AddShop)
	router.With(md.Auth("ADMIN")).Put("/shops", shopController.EditShop)
	return router
}

func RouterCategories(router chi.Router, categoryController *v1.Category) chi.Router {
	router.With(md.Auth("USER")).Get("/categories", categoryController.ListAllCategories)
	router.With(md.Auth("ADMIN")).Post("/categories", categoryController.AddCategory)
	router.With(md.Auth("ADMIN")).Put("/categories", categoryController.EditCategory)
	return router
}

func RouterProduct(router chi.Router, productController *v1.Product) chi.Router {
	router.With(md.Auth("ADMIN")).Post("/products", productController.AddProduct)
	router.With(md.Auth("ADMIN")).Put("/products", productController.EditProduct)
	router.With(md.Auth("USER")).Get("/products", productController.ListAllProducts)
	router.With(md.Auth("USER")).Get("/categories/{categoryID:.+}/products", productController.SearchProductsByCategory)
	router.With(md.Auth("USER")).Get("/search/{product_name:.+}", productController.SearchProductByName)
	router.With(md.Auth("USER")).Get("/shops/{shopID:.+}/products", productController.SearchActiveProductsOfShop)
	return router
}

func RouterPrice(router chi.Router, priceController *v1.Price) chi.Router {
	router.With(md.Auth("ADMIN")).Post("/prices", priceController.AddPrice)
	router.With(md.Auth("ADMIN")).Put("/prices", priceController.EditPrice)
	router.With(md.Auth("USER")).Get("/prices", priceController.ListAllPrices)
	return router
}

func RouterUser(router chi.Router, usersController *v1.Users) chi.Router {
	router.With(md.Auth("ADMIN")).Post("/users", usersController.AddUser)
	router.With(md.Auth("ADMIN")).Put("/users", usersController.EditUser)
	router.With(md.Auth("ADMIN")).Put("/users/addrole", usersController.AddRole)
	router.With(md.Auth("ADMIN")).Put("/users/removerole", usersController.RemoveRole)
	return router
}

func RouterAuth(router chi.Router, authController *v1.Auth) chi.Router {
	router.Post("/auth", authController.Token)
	return router
}
