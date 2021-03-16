package fsipc

type PingRequest struct {
	Id      string `json:"-"`
	Payload string `json:"payload" validate:"required"`
}

type GetScoresRequest struct {
	Id     string `json:"-"`
	ApiKey string `json:"api-key" validate:"required"`
	Hash   string `json:"hash" validate:"required"`
}

type SubmitScoreRequest struct {
	Id          string `json:"-"`
	ApiKey      string `json:"api-key" validate:"required"`
	ProfileName string `json:"profile-name" validate:"required""`
	Hash        string `json:"hash" validate:"required"`
	Rate        int    `json:"rate" validate:"min=5,max=200"`
	Score       int    `json:"score" validate:"min=0,max=10000"`
}
