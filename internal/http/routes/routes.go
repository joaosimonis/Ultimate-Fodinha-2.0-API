package routes

import (
	"net/http" // Pacote padrão para criar APIs REST

	"github.com/gin-gonic/gin"        // Framework para criar APIs REST
	"github.com/jackc/pgx/v5/pgxpool" // Pool de conexões com o PostgreSQL

	"github.com/joaosimonis/ultimate-fodinha-api/internal/http/handlers" // É o meu pacote com os handlers HTTP
	"github.com/joaosimonis/ultimate-fodinha-api/internal/matches"
	"github.com/joaosimonis/ultimate-fodinha-api/internal/players"
)

// Recebe o router do Gin e a pool de conexões com o Postgres
func SetupRoutes(r *gin.Engine, pool *pgxpool.Pool) {
	// 1. Cria os handlers
	playersRepo := players.NewRepository(pool)
	playersService := players.NewService(playersRepo)
	playersHandler := handlers.NewPlayerHandler(playersService)

	matchesRepo := matches.NewRepository(pool)
	matchesService := matches.NewService(matchesRepo)
	matchesHandler := handlers.NewMatchHandler(matchesService)

	rankingHandler := handlers.NewRankingHandler(pool)

	// 2. Grupo versionado -> basicamente para não repetirmos /api/v1 nas rotas que criarmos
	v1 := r.Group("/api/v1")
	{
		// 3. Rota de check da API (apenas para ver se está tudo certo)
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		// 4. Rotas da API, usando Group para melhor separação
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

		ranking := v1.Group("/ranking")
		{
			ranking.GET("/:mode", rankingHandler.GetByMode)
		}

	}
}
