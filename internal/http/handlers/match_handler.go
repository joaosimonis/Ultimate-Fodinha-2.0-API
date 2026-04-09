package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MatchHandler struct {
	pool *pgxpool.Pool
}

func NewMatchHandler(pool *pgxpool.Pool) *MatchHandler {
	return &MatchHandler{pool: pool}
}

func (h *MatchHandler) Create(c *gin.Context) {
	c.JSON(201, gin.H{"message": "criar partida"})
}

func (h *MatchHandler) GetByRoomCode(c *gin.Context) {
	c.JSON(200, gin.H{"message": "buscar partida"})
}

func (h *MatchHandler) Finish(c *gin.Context) {
	c.JSON(200, gin.H{"message": "finalizar partida"})
}

func (h *MatchHandler) ListLive(c *gin.Context) {
	c.JSON(200, gin.H{"matches": []string{}})
}
