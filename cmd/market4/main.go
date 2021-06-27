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

	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, dsn)

	storage := repository.NewStorage(pool)
	cont := controllers.NewMarketController(*storage)
	router := httpserver.NewRouter(*chi.NewRouter(), cont)

	server := http.Server{
		Addr:    addr,
		Handler: &router,
	}

	return server.ListenAndServe()
}
