package main

import (
	"context"
	"log"

	"github.com/jsimonis/ultimate-fodinha-api/internal/config"
	"github.com/jsimonis/ultimate-fodinha-api/internal/db"
	"github.com/jsimonis/ultimate-fodinha-api/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	conn := db.Connect(cfg)
	defer conn.Close(context.Background())

	router := gin.Default()

	routes.SetupRoutes(router)

	log.Printf("servidor rodando na porta %s", cfg.AppPort)
	log.Fatal(router.Run(":" + cfg.AppPort))
}