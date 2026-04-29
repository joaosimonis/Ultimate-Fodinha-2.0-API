package matches

import "time"

type Match struct {
	ID         string        `json:"id"`
	Mode       string        `json:"mode"`
	Status     string        `json:"status"`
	RoomCode   string        `json:"room_code"`
	MaxScore   *int          `json:"max_score,omitempty"`
	Lives      *int          `json:"lives,omitempty"`
	CreatedAt  time.Time     `json:"created_at"`
	FinishedAt *time.Time    `json:"finished_at,omitempty"`
	Players    []MatchPlayer `json:"players,omitempty"`
}

type MatchPlayer struct {
	ID         string `json:"id"`
	MatchID    string `json:"match_id"`
	PlayerID   string `json:"player_id"`
	TeamIndex  *int   `json:"team_index"`
	FinalScore int    `json:"final_score"`
	Placement  *int   `json:"placement"`
	IsWinner   bool   `json:"is_winner"`
}

type CreateMatchRequest struct {
	Mode     string                     `json:"mode" binding:"required"`
	RoomCode string                     `json:"room_code" binding:"required"`
	MaxScore *int                       `json:"max_score"`
	Lives    *int                       `json:"lives"`
	Players  []CreateMatchPlayerRequest `json:"players"`
}

type CreateMatchPlayerRequest struct {
	PlayerID  string `json:"player_id" binding:"required"`
	TeamIndex *int   `json:"team_index"`
}
