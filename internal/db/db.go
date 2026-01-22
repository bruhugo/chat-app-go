package db

import (
	"database/sql"
	"embed"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var embedMigrationsFS embed.FS

func ConnectAndMigrate() (*sql.DB, error) {

	user := os.Getenv("db_user")
	password := os.Getenv("db_password")
	host := os.Getenv("db_host")
	port := os.Getenv("db_port")
	database := os.Getenv("db_database")

	dbString := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + database

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
