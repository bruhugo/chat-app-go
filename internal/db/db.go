package db

import (
	"database/sql"
	"embed"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grongoglongo/chatter-go/internal/config"

	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var embedMigrationsFS embed.FS

func ConnectAndMigrate() (*sql.DB, error) {

	user := config.EnvConfig.DbUser
	password := config.EnvConfig.DbPassword
	host := config.EnvConfig.DbHost
	port := config.EnvConfig.DbPort
	database := config.EnvConfig.DbDatabase

	dbString := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + database + "?parseTime=true"

	log.Print("db string is: " + dbString)
	log.Print("Starting database migration...")

	// database migration
	src, err := iofs.New(embedMigrationsFS, "migrations")
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithSourceInstance("iofs", src, "mysql://"+dbString)

	if err != nil {
		return nil, err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, err
	}

	log.Print("Database migration concluded. Starting connection...")

	// database connection
	db, err := sql.Open("mysql", dbString)

	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(time.Minute * 3)

	if err = db.Ping(); err != nil {
		return nil, err
	}

	log.Print("Database connected.")

	return db, nil
}
