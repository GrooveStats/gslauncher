package fsipc

type PingRequest struct {
	Id       string `json:"-"`
	Protocol int    `json:"protocol" validate:"required"`
}

type GsNewSessionRequest struct {
	Id string `json:"-"`
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

type gsScoreSubmitPlayerData struct {
	ApiKey      string `json:"apiKey" validate:"required"`
	ProfileName string `json:"profileName" validate:"required"`
	ChartHash   string `json:"chartHash" validate:"required"`
	Score       int    `json:"score" validate:"min=0,max=10000"`
	Comment     string `json:"comment" validate:"required"`
	Rate        int    `json:"rate" validate:"min=0,max=200"`
}

type GsScoreSubmitRequest struct {
	Id                    string `json:"-"`
	MaxLeaderboardResults *int   `json:"maxLeaderboardResults"`

	Player1 *gsScoreSubmitPlayerData `json:"player1"`
	Player2 *gsScoreSubmitPlayerData `json:"player2"`
}
