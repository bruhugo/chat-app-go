package main

import (
	"log"

	"github.com/grongoglongo/chatter-go/internal/db"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Panic(err)
	}

	db, err := db.ConnectAndMigrate()
	if err != nil {
		log.Panic("Error migrating or connecting to database: " + err.Error())
	}

	defer db.Close()
}
