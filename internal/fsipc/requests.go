package fsipc

type PingRequest struct {
	Id string `json:"-"`
}

type GsNewSessionRequest struct {
	Id string `json:"-"`
}

type GsPlayerScoresRequest struct {
	Id            string  `json:"-"`
	Chart         string  `json:"chart" validate:"required"`
	ApiKeyPlayer1 *string `json:"api-key-player-1"`
	ApiKeyPlayer2 *string `json:"api-key-player-2"`
}

type GsPlayerLeaderboardsRequest struct {
	Id                    string  `json:"-"`
	Chart                 string  `json:"chart" validate:"required"`
	MaxLeaderboardResults *int    `json:"max-leaderboard-results"`
	ApiKeyPlayer1         *string `json:"api-key-player-1"`
	ApiKeyPlayer2         *string `json:"api-key-player-2"`
}

type gsScoreSubmitPlayerData struct {
	ApiKey      string `json:"api-key"`
	ProfileName string `json:"profile-name" validate:"required"`
	Score       int    `json:"score" validate:"min=0,max=10000"`
	Comment     string `json:"comment" validate:"required"`
	Rate        int    `json:"rate" validate:"min=0,max=300"`
}

type GsScoreSubmitRequest struct {
	Id                    string `json:"-"`
	Chart                 string `json:"chart" validate:"required"`
	MaxLeaderboardResults *int   `json:"max-leaderboard-results"`

	Player1 *gsScoreSubmitPlayerData `json:"player1"`
	Player2 *gsScoreSubmitPlayerData `json:"player2"`
}
