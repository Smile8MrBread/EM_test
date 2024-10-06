package main

import (
	"embed"
	"fmt"
	"github.com/Smile8MrBread/EM_test/app/internal/config"
	"github.com/Smile8MrBread/EM_test/app/internal/storage/postgres"
	"github.com/Smile8MrBread/EM_test/app/pkg/migrator"
)

//go:embed migrations/*.sql
var MigrationsFS embed.FS

func main() {
	cfg := config.MustLoad()

	m := migrator.MustGetNewMigrator(MigrationsFS, "migrations")

	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Storage.User, cfg.Storage.Password, cfg.Storage.Host, cfg.Storage.Port, cfg.Storage.DBName)
	conn, err := postgres.NewConn(url)
	if err != nil {
		panic(err)
	}

	err = m.ApplyMigrations(conn)
	if err != nil {
		panic(err)
	}
}
