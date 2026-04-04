package players

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	repo *Repository
}

func NewHandler(pool *pgxpool.Pool) *Handler {
	return &Handler{repo: NewRepository(pool)}
}

func (h *Handler) List(c *gin.Context) {
	players, err := h.repo.FindAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao buscar jogadores"})
		return
	}

	if players == nil {
		players = []Player{}
	}

	c.JSON(http.StatusOK, gin.H{"players": players})
}

func (h *Handler) Create(c *gin.Context) {
	var req CreatePlayerRequest
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
		c.JSON(http.StatusConflict, gin.H{"error": ErrNameTaken.Error()})
		return
	}

	player, err := h.repo.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao criar jogador"})
		return
	}

	c.JSON(http.StatusCreated, player)
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": ErrNotFound.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao remover jogador"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "jogador removido com sucesso"})
}