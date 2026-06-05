package matches

import (
	"errors"
	"testing"
)

func TestValidateTrucoFinishUsesTeamID(t *testing.T) {
	maxScore := 15
	match := trucoMatchForTest(maxScore)

	results, err := validateTrucoFinish(match, FinishMatchRequest{
		Teams: []FinishMatchTeamRequest{
			{
				TeamID:     "team-a",
				TeamName:   "Duda + Fabri + Grade",
				FinalScore: intPtr(15),
			},
			{
				TeamID:     "team-b",
				TeamName:   "João + Kevin + Knak",
				FinalScore: intPtr(0),
			},
		},
	})
	if err != nil {
		t.Fatalf("validateTrucoFinish returned error: %v", err)
	}

	assertTrucoResults(t, results, "player-a1", 15, true)
	assertTrucoResults(t, results, "player-b1", 0, false)
}

func TestValidateTrucoFinishFallsBackToTeamName(t *testing.T) {
	maxScore := 15
	match := trucoMatchForTest(maxScore)

	results, err := validateTrucoFinish(match, FinishMatchRequest{
		Teams: []FinishMatchTeamRequest{
			{
				TeamName:   "Fabri + Duda + Grade",
				FinalScore: intPtr(15),
			},
			{
				TeamName:   "João + Kevin + Knak",
				FinalScore: intPtr(0),
			},
		},
	})
	if err != nil {
		t.Fatalf("validateTrucoFinish returned error: %v", err)
	}

	assertTrucoResults(t, results, "player-a1", 15, true)
	assertTrucoResults(t, results, "player-b1", 0, false)
}

func TestValidateTrucoFinishRejectsUnknownTeamNameWithoutTeamID(t *testing.T) {
	maxScore := 15
	match := trucoMatchForTest(maxScore)

	_, err := validateTrucoFinish(match, FinishMatchRequest{
		Teams: []FinishMatchTeamRequest{
			{
				TeamName:   "Duda + Fabri + Grade",
				FinalScore: intPtr(15),
			},
			{
				TeamName:   "João + Kevin + Knak",
				FinalScore: intPtr(0),
			},
		},
	})
	if !errors.Is(err, ErrUnknownTeam) {
		t.Fatalf("expected ErrUnknownTeam, got %v", err)
	}
}

func trucoMatchForTest(maxScore int) Match {
	return Match{
		Mode:     ModeTruco,
		MaxScore: &maxScore,
		Teams: []MatchTeam{
			{
				ID:       "team-a",
				MatchID:  "match-1",
				TeamName: "Fabri + Duda + Grade",
			},
			{
				ID:       "team-b",
				MatchID:  "match-1",
				TeamName: "João + Kevin + Knak",
			},
		},
		Players: []MatchPlayer{
			{
				PlayerID: "player-a1",
				TeamID:   "team-a",
				TeamName: "Fabri + Duda + Grade",
			},
			{
				PlayerID: "player-a2",
				TeamID:   "team-a",
				TeamName: "Fabri + Duda + Grade",
			},
			{
				PlayerID: "player-a3",
				TeamID:   "team-a",
				TeamName: "Fabri + Duda + Grade",
			},
			{
				PlayerID: "player-b1",
				TeamID:   "team-b",
				TeamName: "João + Kevin + Knak",
			},
			{
				PlayerID: "player-b2",
				TeamID:   "team-b",
				TeamName: "João + Kevin + Knak",
			},
			{
				PlayerID: "player-b3",
				TeamID:   "team-b",
				TeamName: "João + Kevin + Knak",
			},
		},
	}
}

func assertTrucoResults(t *testing.T, results []FinishMatchPlayerResult, playerID string, score int, isWinner bool) {
	t.Helper()

	for _, result := range results {
		if result.PlayerID != playerID {
			continue
		}

		if result.FinalScore == nil || *result.FinalScore != score {
			t.Fatalf("expected score %d for %s, got %v", score, playerID, result.FinalScore)
		}

		if result.IsWinner != isWinner {
			t.Fatalf("expected winner %t for %s, got %t", isWinner, playerID, result.IsWinner)
		}

		return
	}

	t.Fatalf("result for player %s not found", playerID)
}

func intPtr(value int) *int {
	return &value
}
