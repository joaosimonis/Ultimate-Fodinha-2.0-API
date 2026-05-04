package matches

import (
	"encoding/json"
	"time"
)

type Match struct {
	ID         string          `json:"id"`
	Mode       string          `json:"mode"`
	Status     string          `json:"status"`
	RoomCode   string          `json:"room_code"`
	HostToken  *string         `json:"host_token,omitempty"`
	MaxScore   *int            `json:"max_score,omitempty"`
	Lives      *int            `json:"lives,omitempty"`
	LiveState  json.RawMessage `json:"live_state,omitempty"`
	CreatedAt  time.Time       `json:"created_at"`
	FinishedAt *time.Time      `json:"finished_at,omitempty"`
	Players    []MatchPlayer   `json:"players,omitempty"`
}

type MatchPlayer struct {
	ID             string `json:"id"`
	MatchID        string `json:"match_id"`
	PlayerID       string `json:"player_id"`
	TeamName       string `json:"team_name,omitempty"`
	FinalScore     *int   `json:"final_score,omitempty"`
	RemainingLives *int   `json:"remaining_lives,omitempty"`
	Placement      *int   `json:"placement"`
	IsWinner       bool   `json:"is_winner"`
}

type CreateMatchRequest struct {
	Mode     string `json:"mode" binding:"required"`
	MaxScore *int   `json:"max_score"`
	Lives    *int   `json:"lives"`
}

type AddMatchPlayersRequest struct {
	Players []CreateMatchPlayerRequest `json:"players" binding:"required"`
}

type CreateMatchPlayerRequest struct {
	PlayerID string `json:"player_id" binding:"required"`
	TeamName string `json:"team_name"`
}

type FinishMatchRequest struct {
	Players []FinishMatchPlayerRequest `json:"players"`
	Teams   []FinishMatchTeamRequest   `json:"teams"`
}

type FinishMatchPlayerRequest struct {
	PlayerID       string `json:"player_id" binding:"required"`
	RemainingLives *int   `json:"remaining_lives"`
	Placement      *int   `json:"placement"`
}

type FinishMatchTeamRequest struct {
	TeamName   string `json:"team_name" binding:"required"`
	FinalScore *int   `json:"final_score"`
}

type FinishMatchPlayerResult struct {
	PlayerID       string
	FinalScore     *int
	RemainingLives *int
	Placement      int
	IsWinner       bool
}
