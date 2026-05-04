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
	List(ctx context.Context, status string) ([]matches.Match, error)
	GetByID(ctx context.Context, id string) (matches.Match, error)
	GetByRoomCode(ctx context.Context, roomCode string) (matches.Match, error)
	Create(ctx context.Context, req matches.CreateMatchRequest) (matches.Match, error)
	AddPlayers(ctx context.Context, id string, req matches.AddMatchPlayersRequest) (matches.Match, error)
	RemovePlayer(ctx context.Context, id string, playerID string) (matches.Match, error)
	Finish(ctx context.Context, id string, req matches.FinishMatchRequest) (matches.Match, error)
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
		if errors.Is(err, matches.ErrRoomCodeTaken) {
			response.Error(c, http.StatusConflict, matches.ErrRoomCodeTaken.Error())
			return
		}

		if h.handleMatchValidationError(c, err) {
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

func (h *MatchHandler) AddPlayers(c *gin.Context) {
	id := c.Param("id")

	var req matches.AddMatchPlayersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "dados inválidos")
		return
	}

	match, err := h.service.AddPlayers(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, matches.ErrNotFound) {
			response.Error(c, http.StatusNotFound, matches.ErrNotFound.Error())
			return
		}

		if h.handleMatchValidationError(c, err) {
			return
		}

		response.Error(c, http.StatusInternalServerError, "erro ao adicionar jogadores")
		return
	}

	response.Data(c, http.StatusOK, match)
}

func (h *MatchHandler) RemovePlayer(c *gin.Context) {
	id := c.Param("id")
	playerID := c.Param("player_id")

	match, err := h.service.RemovePlayer(c.Request.Context(), id, playerID)
	if err != nil {
		if errors.Is(err, matches.ErrNotFound) {
			response.Error(c, http.StatusNotFound, matches.ErrNotFound.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, "erro ao remover jogador da partida")
		return
	}

	response.Data(c, http.StatusOK, match)
}

func (h *MatchHandler) Finish(c *gin.Context) {
	id := c.Param("id")

	var req matches.FinishMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "dados inválidos")
		return
	}

	match, err := h.service.Finish(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, matches.ErrNotFound) {
			response.Error(c, http.StatusNotFound, matches.ErrNotFound.Error())
			return
		}

		if h.handleMatchValidationError(c, err) {
			return
		}

		response.Error(c, http.StatusInternalServerError, "erro ao finalizar partida")
		return
	}

	response.Data(c, http.StatusOK, match)
}

func (h *MatchHandler) List(c *gin.Context) {
	status := c.Query("status")

	matchesList, err := h.service.List(c.Request.Context(), status)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "erro ao buscar partidas")
		return
	}

	response.Data(c, http.StatusOK, matchesList)
}

func (h *MatchHandler) handleMatchValidationError(c *gin.Context, err error) bool {
	validationErrors := []error{
		matches.ErrInvalidMode,
		matches.ErrInvalidMaxScore,
		matches.ErrMaxScoreRequired,
		matches.ErrMaxScoreNotAllowed,
		matches.ErrMaxScoreValueNotAllowed,
		matches.ErrLivesRequired,
		matches.ErrInvalidLives,
		matches.ErrLivesNotAllowed,
		matches.ErrLivesValueNotAllowed,
		matches.ErrPlayersRequired,
		matches.ErrPlayerIDRequired,
		matches.ErrDuplicatePlayer,
		matches.ErrTrucoPlayerCount,
		matches.ErrTrucoTeamRequired,
		matches.ErrTrucoTeamCount,
		matches.ErrTrucoTeamBalance,
		matches.ErrMatchAlreadyFinished,
		matches.ErrFinishPlayersMismatch,
		matches.ErrFinalScoreRequired,
		matches.ErrFinalScoreNotAllowed,
		matches.ErrInvalidFinalScore,
		matches.ErrRemainingLivesRequired,
		matches.ErrRemainingLivesNotAllowed,
		matches.ErrInvalidRemainingLives,
		matches.ErrPlacementRequired,
		matches.ErrInvalidPlacement,
		matches.ErrInvalidWinnerCount,
		matches.ErrTrucoInvalidResult,
		matches.ErrPlayersNotAllowed,
		matches.ErrTeamsRequired,
		matches.ErrTeamsNotAllowed,
		matches.ErrTeamScoreRequired,
		matches.ErrInvalidTeamScore,
		matches.ErrUnknownTeam,
		matches.ErrInvalidHostToken,
		matches.ErrInvalidLiveState,
	}

	for _, validationErr := range validationErrors {
		if errors.Is(err, validationErr) {
			response.Error(c, http.StatusBadRequest, validationErr.Error())
			return true
		}
	}

	return false
}
