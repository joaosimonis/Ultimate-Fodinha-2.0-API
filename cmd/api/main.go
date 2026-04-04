package main

import (
	"log"

	"github.com/jsimonis/ultimate-fodinha-api/internal/config"
	"github.com/jsimonis/ultimate-fodinha-api/internal/db"
	"github.com/jsimonis/ultimate-fodinha-api/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	pool := db.Connect(cfg)
	defer pool.Close()

	router := gin.Default()
	routes.SetupRoutes(router, pool)

	log.Printf("🚀 servidor rodando na porta %s", cfg.AppPort)
	log.Fatal(router.Run(":" + cfg.AppPort))
}