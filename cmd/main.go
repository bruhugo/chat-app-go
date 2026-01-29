package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/grongoglongo/chatter-go/internal/config"
	"github.com/grongoglongo/chatter-go/internal/db"
	"github.com/grongoglongo/chatter-go/internal/routes"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Panic(err)
	}

	db, err := db.ConnectAndMigrate()
	if err != nil {
		log.Panic("Error migrating or connecting to database: " + err.Error())
	}

	defer func() {
		db.Close()
		log.Print("Database disconnected.")
	}()

	router := gin.Default()

	routes.ApplyRoutes(router, db)

	port := config.Port
	if config.Port == "" {
		log.Print("Port not provided in env variables (port), hence using default 8080")
		port = "8080"
	}

	router.Run(":" + port)
}
