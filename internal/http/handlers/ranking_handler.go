package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/joaosimonis/ultimate-fodinha-api/internal/http/response"
	"github.com/joaosimonis/ultimate-fodinha-api/internal/ranking"
)

type RankingService interface {
	Get(ctx context.Context, req ranking.GetRankingRequest) ([]ranking.Entry, error)
	GetTrucoTeams(ctx context.Context, req ranking.GetRankingRequest) ([]ranking.TeamEntry, error)
}

type RankingHandler struct {
	service RankingService
}

func NewRankingHandler(service RankingService) *RankingHandler {
	return &RankingHandler{service: service}
}

func (h *RankingHandler) Get(c *gin.Context) {
	mode := c.Query("mode")
	if mode == "" {
		mode = c.Param("mode")
	}

	entries, err := h.service.Get(c.Request.Context(), ranking.GetRankingRequest{
		Mode:     mode,
		TeamSize: c.Query("team_size"),
		Limit:    c.Query("limit"),
	})
	if err != nil {
		if isRankingValidationError(err) {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, "erro ao buscar ranking")
		return
	}

	response.Data(c, http.StatusOK, entries)
}

func (h *RankingHandler) GetTrucoTeams(c *gin.Context) {
	entries, err := h.service.GetTrucoTeams(c.Request.Context(), ranking.GetRankingRequest{
		TeamSize: c.Query("team_size"),
		Limit:    c.Query("limit"),
	})
	if err != nil {
		if isRankingValidationError(err) {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, "erro ao buscar ranking de duplas")
		return
	}

	response.Data(c, http.StatusOK, entries)
}

func isRankingValidationError(err error) bool {
	return errors.Is(err, ranking.ErrModeRequired) ||
		errors.Is(err, ranking.ErrInvalidMode) ||
		errors.Is(err, ranking.ErrTeamSizeRequired) ||
		errors.Is(err, ranking.ErrTeamSizeNotAllowed) ||
		errors.Is(err, ranking.ErrInvalidTeamSize) ||
		errors.Is(err, ranking.ErrInvalidLimit)
}
