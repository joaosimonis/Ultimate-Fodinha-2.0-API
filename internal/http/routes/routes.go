package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/joaosimonis/ultimate-fodinha-api/internal/http/handlers"
	"github.com/joaosimonis/ultimate-fodinha-api/internal/matches"
	"github.com/joaosimonis/ultimate-fodinha-api/internal/players"
	"github.com/joaosimonis/ultimate-fodinha-api/internal/ranking"
	"github.com/joaosimonis/ultimate-fodinha-api/internal/ws"
)

func SetupRoutes(r *gin.Engine, pool *pgxpool.Pool) {
	playersRepo := players.NewRepository(pool)
	playersService := players.NewService(playersRepo)
	playersHandler := handlers.NewPlayerHandler(playersService)

	matchesRepo := matches.NewRepository(pool)
	matchesService := matches.NewService(matchesRepo)
	matchesHandler := handlers.NewMatchHandler(matchesService)

	rankingRepo := ranking.NewRepository(pool)
	rankingService := ranking.NewService(rankingRepo)
	rankingHandler := handlers.NewRankingHandler(rankingService)

	wsHub := ws.NewHub()
	wsHandler := handlers.NewWebSocketHandler(wsHub, matchesService)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		players := v1.Group("/players")
		{
			players.GET("", playersHandler.ListPlayers)
			players.GET("/:id", playersHandler.GetPlayerByID)
			players.POST("", playersHandler.CreatePlayer)
			players.DELETE("/:id", playersHandler.DeletePlayer)
		}

		matches := v1.Group("/matches")
		{
			matches.GET("", matchesHandler.List)
			matches.GET("/:id", matchesHandler.GetByID)
			matches.POST("", matchesHandler.Create)
			matches.POST("/:id/players", matchesHandler.AddPlayers)
			matches.DELETE("/:id/players/:player_id", matchesHandler.RemovePlayer)
			matches.PATCH("/:id/finish", matchesHandler.Finish)
		}

		rooms := v1.Group("/rooms")
		{
			rooms.GET("/:room_code", matchesHandler.GetByRoomCode)
		}

		websocket := v1.Group("/ws")
		{
			websocket.GET("/rooms/:room_code", wsHandler.JoinRoom)
		}

		ranking := v1.Group("/ranking")
		{
			ranking.GET("", rankingHandler.Get)
			ranking.GET("/truco/teams", rankingHandler.GetTrucoTeams)
			ranking.GET("/:mode", rankingHandler.Get)
		}

	}
}
