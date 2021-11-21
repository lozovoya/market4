package main

import (
	"context"
	"go.uber.org/zap"
	"log"
	"market4/internal/api/auth"
	"market4/internal/api/httpserver"
	controllers "market4/internal/api/v1"
	cache2 "market4/internal/cache"
	"market4/internal/repository"
	"net"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	defaultPort     = "9999"
	defaultHost     = "0.0.0.0"
	defaultDSN      = "postgres://app:pass@marketdb:5432/marketdb"
	defaultCacheDSN = "redis://marketcache:6379/0"
	PRIVATEKEY      = "./keys/private.key"
	PUBLICKEY       = "./keys/public.key"
)

func main() {
	port, ok := os.LookupEnv("MARKET_PORT")
	if !ok {
		port = defaultPort
	}

	host, ok := os.LookupEnv("MARKET_HOST")
	if !ok {
		host = defaultHost
	}

	dsn, ok := os.LookupEnv("MARKET_DB")
	if !ok {
		dsn = defaultDSN
	}

	cacheDSN, ok := os.LookupEnv("MARKET_CACHE")
	if !ok {
		cacheDSN = defaultCacheDSN
	}

	privateJWTKey, ok := os.LookupEnv("privateJWTKey")
	if !ok {
		privateJWTKey = PRIVATEKEY
	}

	publicJWTKey, ok := os.LookupEnv("publicJWTKey")
	if !ok {
		publicJWTKey = PUBLICKEY
	}

	if err := execute(net.JoinHostPort(host, port), dsn, cacheDSN, privateJWTKey, publicJWTKey); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
func execute(addr, dsn, cacheDSN, privateJWTKey, publicJWTKey string) (err error) {
	lg := zap.NewExample()
	defer lg.Sync()

	cachePool := cache2.InitCache(cacheDSN)
	cache := cache2.NewRedisCache(cachePool)

	shopCtx := context.Background()
	shopPool, err := pgxpool.Connect(shopCtx, dsn)
	if err != nil {
		lg.Error("Execute", zap.Error(err))
		return err
	}
	shopRepo := repository.NewShopRepository(shopPool)
	shopController := controllers.NewShop(shopRepo, lg)

	categoryCtx := context.Background()
	categoryPool, err := pgxpool.Connect(categoryCtx, dsn)
	if err != nil {
		lg.Error("Execute", zap.Error(err))
		return err
	}
	categoryRepo := repository.NewCategoryRepository(categoryPool)
	categoryController := controllers.NewCategory(categoryRepo, lg)

	priceCtx := context.Background()
	pricePool, err := pgxpool.Connect(priceCtx, dsn)
	if err != nil {
		lg.Error("Execute", zap.Error(err))
		return err
	}
	priceRepo := repository.NewPriceRepository(pricePool)
	priceController := controllers.NewPrice(priceRepo, lg)

	productCtx := context.Background()
	productPool, err := pgxpool.Connect(productCtx, dsn)
	if err != nil {
		lg.Error("Execute", zap.Error(err))
		return err
	}
	productRepo := repository.NewProductRepository(productPool, categoryRepo, shopRepo, priceRepo)
	productController := controllers.NewProduct(productRepo, priceRepo, cache, lg)

	usersCtx := context.Background()
	usersPool, err := pgxpool.Connect(usersCtx, dsn)
	if err != nil {
		lg.Error("Execute", zap.Error(err))
		return err
	}
	usersRepo := repository.NewUsersRepo(usersPool)
	usersController := controllers.NewUser(usersRepo, lg)

	authService := auth.NewAuthService(privateJWTKey, publicJWTKey, usersRepo, lg)
	authController := controllers.NewAuth(*authService, usersRepo, lg)

	router := httpserver.NewRouter(chi.NewRouter(), lg,
		shopController,
		categoryController,
		productController,
		priceController,
		usersController,
		authController)

	server := http.Server{
		Addr:    addr,
		Handler: &router,
	}

	return server.ListenAndServe()
}
