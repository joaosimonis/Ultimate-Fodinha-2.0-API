package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/joaosimonis/ultimate-fodinha-api/internal/http/handlers"
)

func SetupRoutes(r *gin.Engine, pool *pgxpool.Pool) {
	playersHandler := handlers.NewPlayerHandler(pool)
	matchesHandler := handlers.NewMatchHandler(pool)
	rankingHandler := handlers.NewRankingHandler(pool)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		v1.GET("/players", playersHandler.List)
		v1.POST("/players", playersHandler.Create)
		v1.DELETE("/players/:id", playersHandler.Delete)

		v1.GET("/matches/live", matchesHandler.ListLive)
		v1.GET("/matches/:room_code", matchesHandler.GetByRoomCode)
		v1.POST("/matches", matchesHandler.Create)
		v1.PATCH("/matches/:id/finish", matchesHandler.Finish)

		v1.GET("/ranking/:mode", rankingHandler.GetByMode)
	}
}
