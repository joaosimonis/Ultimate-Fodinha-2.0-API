package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/joaosimonis/ultimate-fodinha-api/internal/http/response"
	"github.com/joaosimonis/ultimate-fodinha-api/internal/players"

	"github.com/gin-gonic/gin"
)

type PlayerService interface {
	ListPlayers(ctx context.Context) ([]players.Player, error)
	GetPlayerByID(ctx context.Context, id string) (players.Player, error)
	CreatePlayer(ctx context.Context, req players.CreatePlayerRequest) (players.Player, error)
	DeletePlayer(ctx context.Context, id string) error
}

type PlayerHandler struct {
	service PlayerService
}

func NewPlayerHandler(service PlayerService) *PlayerHandler {
	return &PlayerHandler{
		service: service,
	}
}

func (h *PlayerHandler) ListPlayers(c *gin.Context) {
	playersList, err := h.service.ListPlayers(c.Request.Context())
	if err != nil {
		log.Printf("erro ao buscar jogadores: %v", err)
		response.Error(c, http.StatusInternalServerError, "erro ao buscar jogadores")
		return
	}

	response.Data(c, http.StatusOK, playersList)
}

func (h *PlayerHandler) GetPlayerByID(c *gin.Context) {
	id := c.Param("id")

	player, err := h.service.GetPlayerByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, players.ErrNotFound) {
			response.Error(c, http.StatusNotFound, players.ErrNotFound.Error())
			return
		}

		log.Printf("erro ao buscar jogador id=%s: %v", id, err)
		response.Error(c, http.StatusInternalServerError, "erro ao buscar jogador")
		return
	}

	response.Data(c, http.StatusOK, player)
}

func (h *PlayerHandler) CreatePlayer(c *gin.Context) {
	var req players.CreatePlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "dados inválidos")
		return
	}

	player, err := h.service.CreatePlayer(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, players.ErrNameTaken) {
			response.Error(c, http.StatusConflict, players.ErrNameTaken.Error())
			return
		}

		log.Printf("erro ao criar jogador name=%s: %v", req.Name, err)
		response.Error(c, http.StatusInternalServerError, "erro ao criar jogador")
		return
	}

	response.Data(c, http.StatusCreated, player)
}

func (h *PlayerHandler) DeletePlayer(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeletePlayer(c.Request.Context(), id); err != nil {
		if errors.Is(err, players.ErrNotFound) {
			response.Error(c, http.StatusNotFound, players.ErrNotFound.Error())
			return
		}

		log.Printf("erro ao remover jogador id=%s: %v", id, err)
		response.Error(c, http.StatusInternalServerError, "erro ao remover jogador")
		return
	}

	response.Message(c, http.StatusOK, "jogador removido com sucesso")
}