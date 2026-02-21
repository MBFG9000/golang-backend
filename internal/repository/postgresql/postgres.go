package postgresrepo

import (
	"context"
	"fmt"
	"log"

	"taskmanager/internal/config"
	"taskmanager/pkg/modules"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type Dialect struct {
	DB *sqlx.DB
}

func NewPostgres(ctx context.Context, cfg *modules.PostgreConfig) *Dialect {

	ConnString := config.GetConnURL(cfg)

	db := sqlx.MustConnect("postgres", ConnString)

	err := db.Ping()

	if err != nil {
		panic(err)
	}

	version := ""
	err = db.QueryRow("select version()").Scan(&version)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(version)

	AutoMigrate(cfg)

	return &Dialect{
		DB: db,
	}
}

func AutoMigrate(cfg *modules.PostgreConfig) {
	sourceURL := "file://database/migrations"
	databaseURL := config.GetConnURL(cfg)

	m, err := migrate.New(sourceURL, databaseURL)

	if err != nil {
		panic(err)
	}

	err = m.Up()

	if err != nil && err != migrate.ErrNoChange {
		panic(err)
	}
}
