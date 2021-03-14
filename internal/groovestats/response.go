package groovestats

type scoreEntry struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
	Date  string `json:"date"`
}

type GetScoresResponse struct {
	Leaderboard []scoreEntry `json:"leaderboard"`

	RpgData struct {
		Leaderboard []scoreEntry `json:"leaderboard,omitempty"`
		RivalScores []scoreEntry `json:"rivalScores,omitempty"`
	} `json:"rpgData"`
}

type AutoSubmitScoreResponse struct {
	Result     string  `json:"result"`
	ScoreDelta *int    `json:"scoreDelta,omitempty"`
	RankingUrl *string `json:"rankingUrl,omitempty"`

	Profile struct {
		AvatarUrl string `json:"avatarUrl"`
	} `json:"profile"`
	Leaderboard []scoreEntry `json:"leaderboard"`

	RpgData struct {
		Name *string `json:"name,omitempty"`
		Url  *string `json:"url,omitempty"`

		Progress *struct {
			StatImprovements struct {
				Gold int `json:"gold"`
				Exp  int `json:"exp"`
				Tp   int `json:"tp"`
				Lp   int `json:"lp"`
				Jp   int `json:"jp"`
			} `json:"statImprovements"`

			SkillImprovements []string `json:"skillImprovements"`

			QuestsCompleted []struct {
				Title   string `json:"title"`
				Rewards []struct {
					Type            string `json:"type"`
					Description     string `json:"description"`
					SongDownloadUrl string `json:"songDownloadUrl"`
				} `json:"rewards"`
			} `json:"questsCompleted"`
		} `json:"progress,omitempty"`

		Leaderboard []scoreEntry `json:"leaderboard,omitempty"`
		RivalScores []scoreEntry `json:"rivalScores,omitempty"`
	} `json:"rpgData"`
}
