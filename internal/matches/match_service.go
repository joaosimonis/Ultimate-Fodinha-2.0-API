package matches

import (
	"context"
	"crypto/rand"
	"fmt"
	"strings"
)

const (
	ModeFodinha        = "fodinha"
	ModeTruco          = "truco"
	minPlayers         = 2
	roomCodeLength     = 6
	roomCodeMaxTries   = 10
	roomCodeAlphabet   = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	StatusOngoing      = "ongoing"
	StatusFinished     = "finished"
	StatusLiveAlias    = "live"
	StatusOngoingAlias = "ongoing"
)

type RepositoryReader interface {
	FindAll(ctx context.Context, status string) ([]Match, error)
	FindByID(ctx context.Context, id string) (Match, error)
	FindByRoomCode(ctx context.Context, roomCode string) (Match, error)
	Create(ctx context.Context, req CreateMatchRequest, roomCode string) (Match, error)
	AddPlayers(ctx context.Context, id string, req AddMatchPlayersRequest) (Match, error)
	RemovePlayer(ctx context.Context, id string, playerID string) (Match, error)
	Finish(ctx context.Context, id string, results []FinishMatchPlayerResult) (Match, error)
	ExistsByRoomCode(ctx context.Context, roomCode string) (bool, error)
}

type Service struct {
	repo RepositoryReader
}

func NewService(repo RepositoryReader) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, status string) ([]Match, error) {
	matchesList, err := s.repo.FindAll(ctx, normalizeStatus(status))
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

	roomCode, err := s.generateUniqueRoomCode(ctx)
	if err != nil {
		return Match{}, err
	}

	return s.repo.Create(ctx, req, roomCode)
}

func (s *Service) AddPlayers(ctx context.Context, id string, req AddMatchPlayersRequest) (Match, error) {
	match, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return Match{}, err
	}

	if err := validatePlayers(req.Players); err != nil {
		return Match{}, err
	}

	if err := validatePlayersByMode(match.Mode, req.Players); err != nil {
		return Match{}, err
	}

	return s.repo.AddPlayers(ctx, id, req)
}

func (s *Service) RemovePlayer(ctx context.Context, id string, playerID string) (Match, error) {
	return s.repo.RemovePlayer(ctx, id, playerID)
}

func (s *Service) Finish(ctx context.Context, id string, req FinishMatchRequest) (Match, error) {
	match, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return Match{}, err
	}

	if match.Status == StatusFinished {
		return Match{}, ErrMatchAlreadyFinished
	}

	results, err := validateFinishResult(match, req)
	if err != nil {
		return Match{}, err
	}

	return s.repo.Finish(ctx, id, results)
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

func validatePlayers(players []CreateMatchPlayerRequest) error {
	if len(players) < minPlayers {
		return ErrPlayersRequired
	}

	seen := make(map[string]struct{}, len(players))
	for _, player := range players {
		playerID := strings.TrimSpace(player.PlayerID)
		if playerID == "" {
			return ErrPlayerIDRequired
		}

		if _, exists := seen[playerID]; exists {
			return ErrDuplicatePlayer
		}

		seen[playerID] = struct{}{}
	}

	return nil
}

func validatePlayersByMode(mode string, players []CreateMatchPlayerRequest) error {
	switch mode {
	case ModeTruco:
		return validateTrucoPlayers(players)
	default:
		return nil
	}
}

func validateTrucoPlayers(players []CreateMatchPlayerRequest) error {
	if len(players) != 4 && len(players) != 6 {
		return ErrTrucoPlayerCount
	}

	teamCounts := make(map[string]int, 2)

	for _, player := range players {
		teamName := strings.TrimSpace(player.TeamName)
		if teamName == "" {
			return ErrTrucoTeamRequired
		}

		teamCounts[teamName]++
	}

	if len(teamCounts) != 2 {
		return ErrTrucoTeamCount
	}

	expectedPlayersPerTeam := len(players) / 2
	for _, count := range teamCounts {
		if count != expectedPlayersPerTeam {
			return ErrTrucoTeamBalance
		}
	}

	return nil
}

func normalizeStatus(status string) string {
	switch status {
	case StatusLiveAlias:
		return StatusOngoing
	case StatusOngoingAlias:
		return StatusOngoing
	default:
		return status
	}
}

func (s *Service) generateUniqueRoomCode(ctx context.Context) (string, error) {
	for range roomCodeMaxTries {
		roomCode, err := generateRoomCode()
		if err != nil {
			return "", err
		}

		exists, err := s.repo.ExistsByRoomCode(ctx, roomCode)
		if err != nil {
			return "", err
		}

		if !exists {
			return roomCode, nil
		}
	}

	return "", ErrRoomCodeTaken
}

func generateRoomCode() (string, error) {
	bytes := make([]byte, roomCodeLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("gerar código da sala: %w", err)
	}

	for i := range bytes {
		bytes[i] = roomCodeAlphabet[int(bytes[i])%len(roomCodeAlphabet)]
	}

	return string(bytes), nil
}

func validateFinishResult(match Match, req FinishMatchRequest) ([]FinishMatchPlayerResult, error) {
	switch match.Mode {
	case ModeFodinha:
		return validateFodinhaFinish(match, req)
	case ModeTruco:
		return validateTrucoFinish(match, req)
	default:
		return nil, ErrInvalidMode
	}
}

func validateFinishPlayers(matchPlayers []MatchPlayer, resultPlayers []FinishMatchPlayerRequest) error {
	if len(matchPlayers) == 0 || len(matchPlayers) != len(resultPlayers) {
		return ErrFinishPlayersMismatch
	}

	expectedPlayers := make(map[string]struct{}, len(matchPlayers))
	for _, player := range matchPlayers {
		expectedPlayers[player.PlayerID] = struct{}{}
	}

	seen := make(map[string]struct{}, len(resultPlayers))
	for _, player := range resultPlayers {
		playerID := strings.TrimSpace(player.PlayerID)
		if playerID == "" {
			return ErrPlayerIDRequired
		}

		if _, exists := expectedPlayers[playerID]; !exists {
			return ErrFinishPlayersMismatch
		}

		if _, exists := seen[playerID]; exists {
			return ErrDuplicatePlayer
		}

		seen[playerID] = struct{}{}
	}

	return nil
}

func validateFodinhaFinish(match Match, req FinishMatchRequest) ([]FinishMatchPlayerResult, error) {
	if match.Lives == nil {
		return nil, ErrLivesRequired
	}

	if len(req.Teams) > 0 {
		return nil, ErrTeamsNotAllowed
	}

	if err := validateFinishPlayers(match.Players, req.Players); err != nil {
		return nil, err
	}

	placements := make(map[int]struct{}, len(req.Players))
	winners := 0
	results := make([]FinishMatchPlayerResult, 0, len(req.Players))
	for _, player := range req.Players {
		if player.RemainingLives == nil {
			return nil, ErrRemainingLivesRequired
		}

		if *player.RemainingLives < 0 || *player.RemainingLives > *match.Lives {
			return nil, ErrInvalidRemainingLives
		}

		if player.Placement == nil {
			return nil, ErrPlacementRequired
		}

		if *player.Placement < 1 || *player.Placement > len(req.Players) {
			return nil, ErrInvalidPlacement
		}

		if _, exists := placements[*player.Placement]; exists {
			return nil, ErrInvalidPlacement
		}
		placements[*player.Placement] = struct{}{}

		isWinner := *player.Placement == 1
		if isWinner {
			winners++
		}

		results = append(results, FinishMatchPlayerResult{
			PlayerID:       strings.TrimSpace(player.PlayerID),
			RemainingLives: player.RemainingLives,
			Placement:      *player.Placement,
			IsWinner:       isWinner,
		})
	}

	if winners != 1 {
		return nil, ErrInvalidWinnerCount
	}

	return results, nil
}

func validateTrucoFinish(match Match, req FinishMatchRequest) ([]FinishMatchPlayerResult, error) {
	if match.MaxScore == nil {
		return nil, ErrMaxScoreRequired
	}

	if len(req.Players) > 0 {
		return nil, ErrPlayersNotAllowed
	}

	if len(req.Teams) != 2 {
		return nil, ErrTeamsRequired
	}

	matchTeams := teamsFromMatchPlayers(match.Players)
	if len(matchTeams) != 2 {
		return nil, ErrTrucoTeamCount
	}

	scores := make(map[string]int, 2)
	for _, team := range req.Teams {
		teamName := strings.TrimSpace(team.TeamName)
		if teamName == "" {
			return nil, ErrTrucoTeamRequired
		}

		if _, exists := matchTeams[teamName]; !exists {
			return nil, ErrUnknownTeam
		}

		if _, exists := scores[teamName]; exists {
			return nil, ErrTrucoTeamCount
		}

		if team.FinalScore == nil {
			return nil, ErrTeamScoreRequired
		}

		if *team.FinalScore < 0 || *team.FinalScore > *match.MaxScore {
			return nil, ErrInvalidTeamScore
		}

		scores[teamName] = *team.FinalScore
	}

	if len(scores) != len(matchTeams) {
		return nil, ErrTeamsRequired
	}

	winnerTeam, err := inferTrucoWinner(scores, *match.MaxScore)
	if err != nil {
		return nil, err
	}

	results := make([]FinishMatchPlayerResult, 0, len(match.Players))
	for _, player := range match.Players {
		score := scores[player.TeamName]
		placement := 2
		isWinner := false
		if player.TeamName == winnerTeam {
			placement = 1
			isWinner = true
		}

		results = append(results, FinishMatchPlayerResult{
			PlayerID:   player.PlayerID,
			FinalScore: &score,
			Placement:  placement,
			IsWinner:   isWinner,
		})
	}

	return results, nil
}

func teamsFromMatchPlayers(players []MatchPlayer) map[string]struct{} {
	teams := make(map[string]struct{}, 2)
	for _, player := range players {
		teamName := strings.TrimSpace(player.TeamName)
		if teamName != "" {
			teams[teamName] = struct{}{}
		}
	}

	return teams
}

func inferTrucoWinner(scores map[string]int, maxScore int) (string, error) {
	winnerTeam := ""
	winnerScore := -1
	for teamName, score := range scores {
		if score == winnerScore {
			return "", ErrTrucoInvalidResult
		}

		if score > winnerScore {
			winnerTeam = teamName
			winnerScore = score
		}
	}

	if winnerScore != maxScore {
		return "", ErrTrucoInvalidResult
	}

	return winnerTeam, nil
}
