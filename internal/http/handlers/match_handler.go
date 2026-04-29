package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/joaosimonis/ultimate-fodinha-api/internal/http/response"
	"github.com/joaosimonis/ultimate-fodinha-api/internal/matches"
)

type MatchService interface {
	ListLive(ctx context.Context) ([]matches.Match, error)
	GetByID(ctx context.Context, id string) (matches.Match, error)
	GetByRoomCode(ctx context.Context, roomCode string) (matches.Match, error)
	Create(ctx context.Context, req matches.CreateMatchRequest) (matches.Match, error)
	Finish(ctx context.Context, id string) (matches.Match, error)
}

type MatchHandler struct {
	service MatchService
}

func NewMatchHandler(service MatchService) *MatchHandler {
	return &MatchHandler{service: service}
}

func (h *MatchHandler) Create(c *gin.Context) {
	var req matches.CreateMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "dados inválidos")
		return
	}

	match, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, matches.ErrInvalidMode) {
			response.Error(c, http.StatusBadRequest, matches.ErrInvalidMode.Error())
			return
		}

		if errors.Is(err, matches.ErrInvalidMaxScore) {
			response.Error(c, http.StatusBadRequest, matches.ErrInvalidMaxScore.Error())
			return
		}

		if errors.Is(err, matches.ErrMaxScoreRequired) {
			response.Error(c, http.StatusBadRequest, matches.ErrMaxScoreRequired.Error())
			return
		}

		if errors.Is(err, matches.ErrMaxScoreNotAllowed) {
			response.Error(c, http.StatusBadRequest, matches.ErrMaxScoreNotAllowed.Error())
			return
		}

		if errors.Is(err, matches.ErrMaxScoreValueNotAllowed) {
			response.Error(c, http.StatusBadRequest, matches.ErrMaxScoreValueNotAllowed.Error())
			return
		}

		if errors.Is(err, matches.ErrLivesRequired) {
			response.Error(c, http.StatusBadRequest, matches.ErrLivesRequired.Error())
			return
		}

		if errors.Is(err, matches.ErrInvalidLives) {
			response.Error(c, http.StatusBadRequest, matches.ErrInvalidLives.Error())
			return
		}

		if errors.Is(err, matches.ErrLivesNotAllowed) {
			response.Error(c, http.StatusBadRequest, matches.ErrLivesNotAllowed.Error())
			return
		}

		if errors.Is(err, matches.ErrLivesValueNotAllowed) {
			response.Error(c, http.StatusBadRequest, matches.ErrLivesValueNotAllowed.Error())
			return
		}

		if errors.Is(err, matches.ErrRoomCodeTaken) {
			response.Error(c, http.StatusConflict, matches.ErrRoomCodeTaken.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, "erro ao criar partida")
		return
	}

	response.Data(c, http.StatusCreated, match)
}

func (h *MatchHandler) GetByRoomCode(c *gin.Context) {
	roomCode := c.Param("room_code")

	match, err := h.service.GetByRoomCode(c.Request.Context(), roomCode)
	if err != nil {
		if errors.Is(err, matches.ErrNotFound) {
			response.Error(c, http.StatusNotFound, matches.ErrNotFound.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, "erro ao buscar partida")
		return
	}

	response.Data(c, http.StatusOK, match)
}

func (h *MatchHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	match, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, matches.ErrNotFound) {
			response.Error(c, http.StatusNotFound, matches.ErrNotFound.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, "erro ao buscar partida")
		return
	}

	response.Data(c, http.StatusOK, match)
}

func (h *MatchHandler) Finish(c *gin.Context) {
	id := c.Param("id")

	match, err := h.service.Finish(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, matches.ErrNotFound) {
			response.Error(c, http.StatusNotFound, matches.ErrNotFound.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, "erro ao finalizar partida")
		return
	}

	response.Data(c, http.StatusOK, match)
}

func (h *MatchHandler) ListLive(c *gin.Context) {
	matchesList, err := h.service.ListLive(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "erro ao buscar partidas ao vivo")
		return
	}

	response.Data(c, http.StatusOK, matchesList)
}
