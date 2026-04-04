package matches

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	pool *pgxpool.Pool
}

func NewHandler(pool *pgxpool.Pool) *Handler {
	return &Handler{pool: pool}
}

func (h *Handler) Create(c *gin.Context) {
	c.JSON(201, gin.H{"message": "criar partida"})
}

func (h *Handler) GetByRoomCode(c *gin.Context) {
	c.JSON(200, gin.H{"message": "buscar partida"})
}

func (h *Handler) Finish(c *gin.Context) {
	c.JSON(200, gin.H{"message": "finalizar partida"})
}

func (h *Handler) ListLive(c *gin.Context) {
	c.JSON(200, gin.H{"matches": []string{}})
}
