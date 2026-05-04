package players

import "context"

type PlayerRepository interface {
	FindAll(ctx context.Context) ([]Player, error)
	FindByID(ctx context.Context, id string) (Player, error)
	Create(ctx context.Context, req CreatePlayerRequest) (Player, error)
	Delete(ctx context.Context, id string) error
	ExistsByName(ctx context.Context, name string) (bool, error)
}

type Service struct {
	repo PlayerRepository
}

func NewService(repo PlayerRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListPlayers(ctx context.Context) ([]Player, error) {
	playersList, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	if playersList == nil {
		return []Player{}, nil
	}

	return playersList, nil
}

func (s *Service) GetPlayerByID(ctx context.Context, id string) (Player, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) CreatePlayer(ctx context.Context, req CreatePlayerRequest) (Player, error) {
	exists, err := s.repo.ExistsByName(ctx, req.Name)
	if err != nil {
		return Player{}, err
	}

	if exists {
		return Player{}, ErrNameTaken
	}

	return s.repo.Create(ctx, req)
}

func (s *Service) DeletePlayer(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
