package matches

import "context"

const (
	ModeFodinha = "fodinha"
	ModeTruco   = "truco"
)

type RepositoryReader interface {
	FindLive(ctx context.Context) ([]Match, error)
	FindByID(ctx context.Context, id string) (Match, error)
	FindByRoomCode(ctx context.Context, roomCode string) (Match, error)
	Create(ctx context.Context, req CreateMatchRequest) (Match, error)
	Finish(ctx context.Context, id string) (Match, error)
	ExistsByRoomCode(ctx context.Context, roomCode string) (bool, error)
}

type Service struct {
	repo RepositoryReader
}

func NewService(repo RepositoryReader) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListLive(ctx context.Context) ([]Match, error) {
	matchesList, err := s.repo.FindLive(ctx)
	if err != nil {
		return nil, err
	}

	if matchesList == nil {
		return []Match{}, nil
	}

	return matchesList, nil
}

func (s *Service) GetByRoomCode(ctx context.Context, roomCode string) (Match, error) {
	return s.repo.FindByRoomCode(ctx, roomCode)
}

func (s *Service) GetByID(ctx context.Context, id string) (Match, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, req CreateMatchRequest) (Match, error) {
	if !isValidMode(req.Mode) {
		return Match{}, ErrInvalidMode
	}

	if err := validateMatchConfig(req); err != nil {
		return Match{}, err
	}

	exists, err := s.repo.ExistsByRoomCode(ctx, req.RoomCode)
	if err != nil {
		return Match{}, err
	}

	if exists {
		return Match{}, ErrRoomCodeTaken
	}

	return s.repo.Create(ctx, req)
}

func (s *Service) Finish(ctx context.Context, id string) (Match, error) {
	return s.repo.Finish(ctx, id)
}

func isValidMode(mode string) bool {
	return mode == ModeFodinha || mode == ModeTruco
}

func validateMatchConfig(req CreateMatchRequest) error {
	switch req.Mode {
	case ModeFodinha:
		return validateFodinhaConfig(req)
	case ModeTruco:
		return validateTrucoConfig(req)
	default:
		return ErrInvalidMode
	}
}

func validateFodinhaConfig(req CreateMatchRequest) error {
	if req.MaxScore != nil {
		return ErrMaxScoreNotAllowed
	}

	if req.Lives == nil {
		return ErrLivesRequired
	}

	if *req.Lives <= 0 {
		return ErrInvalidLives
	}

	if *req.Lives > 5 {
		return ErrLivesValueNotAllowed
	}

	return nil
}

func validateTrucoConfig(req CreateMatchRequest) error {
	if req.Lives != nil {
		return ErrLivesNotAllowed
	}

	return validateMaxScore(req.MaxScore)
}

func validateMaxScore(maxScore *int) error {
	if maxScore == nil {
		return ErrMaxScoreRequired
	}

	if *maxScore <= 0 {
		return ErrInvalidMaxScore
	}

	if *maxScore == 15 || *maxScore == 24 || *maxScore == 30 {
		return nil
	}

	return ErrMaxScoreValueNotAllowed
}
