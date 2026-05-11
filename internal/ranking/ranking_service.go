package ranking

import (
	"context"
	"strconv"
	"strings"
)

const (
	ModeFodinha = "fodinha"
	ModeTruco   = "truco"

	defaultLimit = 10
	maxLimit     = 100
)

type RepositoryReader interface {
	FindFodinha(ctx context.Context, limit int) ([]Entry, error)
	FindTruco(ctx context.Context, teamSize int, limit int) ([]Entry, error)
	FindTrucoTeams(ctx context.Context, teamSize int, limit int) ([]TeamEntry, error)
}

type Service struct {
	repo RepositoryReader
}

func NewService(repo RepositoryReader) *Service {
	return &Service{repo: repo}
}

func (s *Service) Get(ctx context.Context, req GetRankingRequest) ([]Entry, error) {
	mode := strings.ToLower(strings.TrimSpace(req.Mode))
	limit, err := parseLimit(req.Limit)
	if err != nil {
		return nil, err
	}

	switch mode {
	case ModeFodinha:
		if strings.TrimSpace(req.TeamSize) != "" {
			return nil, ErrTeamSizeNotAllowed
		}
		return s.repo.FindFodinha(ctx, limit)
	case ModeTruco:
		teamSize, err := parseTeamSize(req.TeamSize)
		if err != nil {
			return nil, err
		}
		return s.repo.FindTruco(ctx, teamSize, limit)
	case "":
		return nil, ErrModeRequired
	default:
		return nil, ErrInvalidMode
	}
}

func (s *Service) GetTrucoTeams(ctx context.Context, req GetRankingRequest) ([]TeamEntry, error) {
	teamSize, err := parseTeamSize(req.TeamSize)
	if err != nil {
		return nil, err
	}

	limit, err := parseLimit(req.Limit)
	if err != nil {
		return nil, err
	}

	return s.repo.FindTrucoTeams(ctx, teamSize, limit)
}

func parseLimit(rawLimit string) (int, error) {
	rawLimit = strings.TrimSpace(rawLimit)
	if rawLimit == "" {
		return defaultLimit, nil
	}

	limit, err := strconv.Atoi(rawLimit)
	if err != nil || limit <= 0 {
		return 0, ErrInvalidLimit
	}

	if limit > maxLimit {
		return maxLimit, nil
	}

	return limit, nil
}

func parseTeamSize(rawTeamSize string) (int, error) {
	rawTeamSize = strings.TrimSpace(rawTeamSize)
	if rawTeamSize == "" {
		return 0, ErrTeamSizeRequired
	}

	teamSize, err := strconv.Atoi(rawTeamSize)
	if err != nil {
		return 0, ErrInvalidTeamSize
	}

	if teamSize != 2 && teamSize != 3 {
		return 0, ErrInvalidTeamSize
	}

	return teamSize, nil
}
