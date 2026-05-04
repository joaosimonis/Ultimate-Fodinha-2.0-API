package bootstrap

import (
	"fmt" // Enviar erros formatados via fmt.Errorf

	"github.com/gin-gonic/gin"                                    // Framework para criar APIs REST
	"github.com/jackc/pgx/v5/pgxpool"                             // Pool de conexões com o PostgreSQL
	"github.com/joaosimonis/ultimate-fodinha-api/internal/config" // Pacote que carrega as variáveis de ambiente
	"github.com/joaosimonis/ultimate-fodinha-api/internal/db"     // Pacote que abre conexão com banco
	"github.com/joaosimonis/ultimate-fodinha-api/internal/http/middleware"
	"github.com/joaosimonis/ultimate-fodinha-api/internal/http/routes" // Pacote que registra as rotas da API
)

// Struct que representa a aplicação já montada
type App struct {
	Config config.Config // É a minha configuração com as variáveis de ambiente
	DB     *pgxpool.Pool // É o pool de conexões com o PostgreSQL
	Router *gin.Engine   // É o servidor HTTP configurado com Gin
}

// Função que inicializa a aplicação
func BuildApp() (*App, error) {
	// 1. Carrega config
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("carregar configurações: %w", err)
	}

	// 2. Conecta no banco
	pool, err := db.Connect(cfg)
	if err != nil {
		return nil, fmt.Errorf("conectar ao banco: %w", err)
	}

	// 3. Cria o router -> o Default cria um *gin.Engine com middlewares padrão, como logger e recovery
	// Resumindo, simplifica muito nossa vida
	router := gin.Default()
	router.Use(middleware.CORS(cfg.AllowedOrigins))

	// 4. Registra as rotas
	routes.SetupRoutes(router, pool)

	// 5. Devolve a aplicação pronta
	return &App{
		Config: cfg,
		DB:     pool,
		Router: router,
	}, nil
}

// Função com receiver, para subir o servidor HTTP
func (a *App) Run() error {
	return a.Router.Run(":" + a.Config.AppPort)
}

// Função com receiver, para fechar a aplicação
func (a *App) Close() {
	if a.DB != nil {
		a.DB.Close()
	}
}
