package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"market4/internal/api/httpserver"
	controllers "market4/internal/api/v1"
	"market4/internal/repository"
	"net"
	"net/http"
	"os"
)

const (
	defaultPort = "9999"
	defaultHost = "0.0.0.0"
	defaultDSN  = "postgres://app:pass@localhost:5432/marketdb"
)

func main() {
	port, ok := os.LookupEnv("APP_PORT")
	if !ok {
		port = defaultPort
	}

	host, ok := os.LookupEnv("APP_HOST")
	if !ok {
		host = defaultHost
	}

	dsn, ok := os.LookupEnv("DB")
	if !ok {
		dsn = defaultDSN
	}

	if err := execute(net.JoinHostPort(host, port), dsn); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func execute(addr string, dsn string) (err error) {

	shopCtx := context.Background()
	shopPool, err := pgxpool.Connect(shopCtx, dsn)
	shopRepo := repository.NewShopRepository(shopPool)
	shopController := controllers.NewShop(shopRepo)

	categoryCtx := context.Background()
	categoryPool, err := pgxpool.Connect(categoryCtx, dsn)
	categoryRepo := repository.NewCategoryRepository(categoryPool)
	categoryController := controllers.NewCategory(categoryRepo)

	productCtx := context.Background()
	productPool, err := pgxpool.Connect(productCtx, dsn)
	productRepo := repository.NewProductRepository(productPool, categoryRepo, shopRepo)
	productController := controllers.NewProduct(productRepo)

	priceCtx := context.Background()
	pricePool, err := pgxpool.Connect(priceCtx, dsn)
	priceRepo := repository.NewPriceRepository(pricePool, productRepo)
	priceController := controllers.NewPrice(priceRepo)

	router := httpserver.NewRouter(*chi.NewRouter(), shopController, categoryController, productController, priceController)

	server := http.Server{
		Addr:    addr,
		Handler: &router,
	}

	return server.ListenAndServe()
}
