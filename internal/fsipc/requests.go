package fsipc

type PingRequest struct {
	Id       string `json:"-"`
	Protocol int    `json:"protocol" validate:"required"`
}

type GsNewSessionRequest struct {
	Id               string `json:"-"`
	ChartHashVersion int    `json:"chartHashVersion" validate:"required"`
}

type gsPlayerData struct {
	ChartHash string `json:"chartHash" validate:"required"`
	ApiKey    string `json:"apiKey" validate:"required"`
}

type GsPlayerScoresRequest struct {
	Id      string        `json:"-"`
	Player1 *gsPlayerData `json:"player1"`
	Player2 *gsPlayerData `json:"player2"`
}

type GsPlayerLeaderboardsRequest struct {
	Id                    string        `json:"-"`
	MaxLeaderboardResults *int          `json:"maxLeaderboardResults"`
	Player1               *gsPlayerData `json:"player1"`
	Player2               *gsPlayerData `json:"player2"`
}

type JudgmentCounts struct {
	FantasticPlus int `json:"fantasticPlus" validate:"min=0"`
	Fantastic     int `json:"fantastic" validate:"min=0"`
	Excellent     int `json:"excellent" validate:"min=0"`
	Great         int `json:"great" validate:"min=0"`
	Decent        int `json:"decent" validate:"min=0"`
	WayOff        int `json:"wayOff" validate:"min=0"`
	Miss          int `json:"miss" validate:"min=0"`
	TotalSteps    int `json:"totalSteps" validate:"min=0"`
	MinesHit      int `json:"minesHit" validate:"min=0"`
	TotalMines    int `json:"totalMines" validate:"min=0"`
	HoldsHeld     int `json:"holdsHeld" validate:"min=0"`
	TotalHolds    int `json:"totalHolds" validate:"min=0"`
	RollsHeld     int `json:"rollsHeld" validate:"min=0"`
	TotalRolls    int `json:"totalRolls" validate:"min=0"`
}

type gsScoreSubmitPlayerData struct {
	ApiKey         string          `json:"apiKey" validate:"required"`
	ProfileName    string          `json:"profileName"`
	ChartHash      string          `json:"chartHash" validate:"required"`
	Score          int             `json:"score" validate:"min=0,max=10000"`
	Comment        string          `json:"comment"`
	Rate           int             `json:"rate" validate:"min=0"`
	JudgmentCounts *JudgmentCounts `json:"judgmentCounts"`
	UsedCmod       *bool           `json:"usedCmod"`
}

type GsScoreSubmitRequest struct {
	Id                    string `json:"-"`
	MaxLeaderboardResults *int   `json:"maxLeaderboardResults"`

	Player1 *gsScoreSubmitPlayerData `json:"player1"`
	Player2 *gsScoreSubmitPlayerData `json:"player2"`
}
