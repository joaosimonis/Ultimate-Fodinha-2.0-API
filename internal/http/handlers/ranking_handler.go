package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RankingHandler struct {
	pool *pgxpool.Pool
}

func NewRankingHandler(pool *pgxpool.Pool) *RankingHandler {
	return &RankingHandler{pool: pool}
}

func (h *RankingHandler) GetByMode(c *gin.Context) {
	c.JSON(200, gin.H{"ranking": []string{}})
}
