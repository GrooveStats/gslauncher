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

type SubmitScoreRequest struct {
	Id          string `json:"-"`
	ApiKey      string `json:"api-key" validate:"required"`
	ProfileName string `json:"profile-name" validate:"required""`
	Hash        string `json:"hash" validate:"required"`
	Rate        int    `json:"rate" validate:"min=5,max=200"`
	Score       int    `json:"score" validate:"min=0,max=10000"`
}
