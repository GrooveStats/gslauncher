package groovestats

type NewSessionResponse struct {
	ActiveEvents []struct {
		Name      string `json:"name"`
		ShortName string `json:"shortName"`
		Url       string `json:"url"`
	} `json:"activeEvents"`

	ServicesAllowed struct {
		ScoreSubmit        bool `json:"scoreSubmit"`
		PlayerScores       bool `json:"playerScores"`
		PlayerLeaderboards bool `json:"playerLeaderboards"`
	} `json:"servicesAllowed"`
}

type leaderBoardEntry struct {
	Name       string  `json:"name"`
	MachineTag *string `json:"machineTag"`
	Score      int     `json:"score"`
	Date       string  `json:"date"`
	Rank       int     `json:"rank"`
	IsSelf     bool    `json:"isSelf"`
	IsRival    bool    `json:"isRival"`
}

type PlayerScoresResponse struct {
	Player1 *[]leaderBoardEntry `json:"player1"`
	Player2 *[]leaderBoardEntry `json:"player2"`
}

type scoreEntry struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
	Date  string `json:"date"`
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
