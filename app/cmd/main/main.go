package main

import (
	"fmt"
	"github.com/Smile8MrBread/EM_test/app/internal/config"
	"github.com/Smile8MrBread/EM_test/app/internal/services"
	"github.com/Smile8MrBread/EM_test/app/internal/storage/postgres"
	"github.com/Smile8MrBread/EM_test/app/internal/transport/rest"
	"github.com/Smile8MrBread/EM_test/app/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"runtime"
)

func main() {
	fmt.Println(runtime.GOROOT())

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)
	log.Info("Lib starting...")

	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Storage.User, cfg.Storage.Password, cfg.Storage.Host, cfg.Storage.Port, cfg.Storage.DBName)
	conn, err := postgres.NewConn(url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	db := postgres.New(conn)

	lib := services.New(log, db)
	log.Info("Lib started")
	rest.StartServer(r, lib)
}
