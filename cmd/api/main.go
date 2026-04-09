package main

import (
	"log"

	"github.com/joaosimonis/ultimate-fodinha-api/internal/config"
	"github.com/joaosimonis/ultimate-fodinha-api/internal/db"
	"github.com/joaosimonis/ultimate-fodinha-api/internal/http/routes"

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