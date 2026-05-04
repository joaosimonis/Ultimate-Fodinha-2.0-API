package ranking

type GetRankingRequest struct {
	Mode     string
	TeamSize string
	Limit    string
}

type Entry struct {
	Position     int      `json:"position"`
	PlayerID     string   `json:"player_id"`
	Name         string   `json:"name"`
	AvatarKey    string   `json:"avatar_key"`
	GamesPlayed  int      `json:"games_played"`
	Wins         *int     `json:"wins,omitempty"`
	Losses       *int     `json:"losses,omitempty"`
	WinRate      *float64 `json:"win_rate,omitempty"`
	GoldMedals   *int     `json:"gold_medals,omitempty"`
	SilverMedals *int     `json:"silver_medals,omitempty"`
	BronzeMedals *int     `json:"bronze_medals,omitempty"`
	Podiums      *int     `json:"podiums,omitempty"`
}
