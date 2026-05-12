package bootstrap

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joaosimonis/ultimate-fodinha-api/internal/config"
	"github.com/joaosimonis/ultimate-fodinha-api/internal/db"
	"github.com/joaosimonis/ultimate-fodinha-api/internal/http/middleware"
	"github.com/joaosimonis/ultimate-fodinha-api/internal/http/routes"
)

type App struct {
	Config config.Config
	DB     *pgxpool.Pool
	Router *gin.Engine
}

func BuildApp() (*App, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("carregar configurações: %w", err)
	}

	pool, err := db.Connect(cfg)
	if err != nil {
		return nil, fmt.Errorf("conectar ao banco: %w", err)
	}

	router := gin.Default()
	router.Use(middleware.CORS(cfg.AllowedOrigins))

	routes.SetupRoutes(router, pool, cfg)

	return &App{
		Config: cfg,
		DB:     pool,
		Router: router,
	}, nil
}

func (a *App) Run() error {
	return a.Router.Run(":" + a.Config.AppPort)
}

func (a *App) Close() {
	if a.DB != nil {
		a.DB.Close()
	}
}
