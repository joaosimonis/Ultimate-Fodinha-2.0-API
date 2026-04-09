package handlers

import (
	"errors"
	"net/http"

	"github.com/joaosimonis/ultimate-fodinha-api/internal/players"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PlayerHandler struct {
	repo *players.Repository
}

func NewPlayerHandler(pool *pgxpool.Pool) *PlayerHandler {
	return &PlayerHandler{
		repo: players.NewRepository(pool),
	}
}


func (h *PlayerHandler) List(c *gin.Context) {
	playersList, err := h.repo.FindAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao buscar jogadores"})
		return
	}

	if playersList == nil {
		playersList = []players.Player{}
	}

	c.JSON(http.StatusOK, gin.H{"players": playersList})
}

func (h *PlayerHandler) Create(c *gin.Context) {
	var req players.CreatePlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos"})
		return
	}

	exists, err := h.repo.ExistsByName(c.Request.Context(), req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao verificar jogador"})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": players.ErrNameTaken.Error()})
		return
	}

	player, err := h.repo.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao criar jogador"})
		return
	}

	c.JSON(http.StatusCreated, player)
}

func (h *PlayerHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, players.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": players.ErrNotFound.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao remover jogador"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "jogador removido com sucesso"})
}