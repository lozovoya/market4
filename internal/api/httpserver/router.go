package httpserver

import (
	"go.uber.org/zap"
	"market4/internal/api/httpserver/md"
	v1 "market4/internal/api/v1"
	"market4/internal/model"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(
	mux *chi.Mux,
	lg *zap.Logger,
	shopController *v1.Shop,
	categoryController *v1.Category,
	productController *v1.Product,
	priceController *v1.Price,
	usersController *v1.Users,
	authController *v1.Auth) chi.Mux {
	mux.Use(middleware.Logger)
	mux.Route("/api/v1", func(router chi.Router) {
		RouterShop(router, shopController, lg)
		RouterCategories(router, categoryController, lg)
		RouterProduct(router, productController, lg)
		RouterPrice(router, priceController, lg)
		RouterUser(router, usersController, lg)
		RouterAuth(router, authController)
	})

	lg.Info("new router is activated")
	return *mux
}

func RouterShop(router chi.Router, shopController *v1.Shop, lg *zap.Logger) chi.Router {
	router.With(md.Auth(model.USER, lg)).Get("/shops", shopController.ListAllShops)
	router.With(md.Auth(model.ADMIN, lg)).Post("/shops", shopController.AddShop)
	router.With(md.Auth(model.ADMIN, lg)).Put("/shops", shopController.EditShop)
	return router
}

func RouterCategories(router chi.Router, categoryController *v1.Category, lg *zap.Logger) chi.Router {
	router.With(md.Auth(model.USER, lg)).Get("/categories", categoryController.ListAllCategories)
	router.With(md.Auth(model.ADMIN, lg)).Post("/categories", categoryController.AddCategory)
	router.With(md.Auth(model.ADMIN, lg)).Put("/categories", categoryController.EditCategory)
	return router
}

func RouterProduct(router chi.Router, productController *v1.Product, lg *zap.Logger) chi.Router {
	router.With(md.Auth(model.ADMIN, lg)).Post("/products", productController.AddProduct)
	router.With(md.Auth(model.ADMIN, lg)).Put("/products", productController.EditProduct)
	router.With(md.Auth(model.USER, lg)).Get("/products", productController.ListAllProducts)
	router.With(md.Auth(model.USER, lg)).Get("/categories/{categoryID:.+}/products", productController.SearchProductsByCategory)
	router.With(md.Auth(model.USER, lg)).Get("/search/{product_name:.+}", productController.SearchProductByName)
	router.With(md.Auth(model.USER, lg)).Get("/shops/{shopID:.+}/products", productController.SearchActiveProductsOfShop)
	return router
}

func RouterPrice(router chi.Router, priceController *v1.Price, lg *zap.Logger) chi.Router {
	router.With(md.Auth(model.ADMIN, lg)).Post("/prices", priceController.AddPrice)
	router.With(md.Auth(model.ADMIN, lg)).Put("/prices", priceController.EditPrice)
	router.With(md.Auth(model.USER, lg)).Get("/prices", priceController.ListAllPrices)
	return router
}

func RouterUser(router chi.Router, usersController *v1.Users, lg *zap.Logger) chi.Router {
	router.With(md.Auth(model.ADMIN, lg)).Post("/users", usersController.AddUser)
	router.With(md.Auth(model.ADMIN, lg)).Put("/users", usersController.EditUser)
	router.With(md.Auth(model.ADMIN, lg)).Put("/users/addrole", usersController.AddRole)
	router.With(md.Auth(model.ADMIN, lg)).Put("/users/removerole", usersController.RemoveRole)
	return router
}

func RouterAuth(router chi.Router, authController *v1.Auth) chi.Router {
	router.Post("/auth", authController.Token)
	return router
}
