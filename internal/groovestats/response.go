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

type playerLeaderboardsPlayerData struct {
	GsLeaderboard []leaderBoardEntry `json:"gsLeaderboard"`

	Rpg *struct {
		Name           string             `json:"name"`
		RpgLeaderboard []leaderBoardEntry `json:"rpgLeaderboard"`
	} `json:"rpg"`
}

type PlayerLeaderboardsResponse struct {
	Player1 *playerLeaderboardsPlayerData `json:"player1"`
	Player2 *playerLeaderboardsPlayerData `json:"player2"`
}

type scoreSubmitPlayerData struct {
	Result        string              `json:"result"`
	ScoreDelta    *int                `json:"scoreDelta,omitempty"`
	GsLeaderboard *[]leaderBoardEntry `json:"gsLeaderboard"`

	Rpg *struct {
		Name       string `json:"name,omitempty"`
		Result     string `json:"result"`
		ScoreDelta *int   `json:"scoreDelta,omitempty"`
		RateDelta  *int   `json:"rateDelta,omitempty"`

		Progress *struct {
			StatImprovements []struct {
				Name   string `json:"name"`
				Gained int    `json:"gained"`
			} `json:"statImprovements"`

			SkillImprovements []string `json:"skillImprovements"`

			QuestsCompleted []struct {
				Title   string `json:"title"`
				Rewards []struct {
					Type            string  `json:"type"`
					Description     string  `json:"description"`
					SongDownloadUrl *string `json:"songDownloadUrl"`
				} `json:"rewards"`
			} `json:"questsCompleted"`
		} `json:"progress"`

		RpgLeaderboard *[]leaderBoardEntry `json:"rpgLeaderboard"`
	} `json:"rpg"`
}

type ScoreSubmitResponse struct {
	Player1 *scoreSubmitPlayerData `json:"player1"`
	Player2 *scoreSubmitPlayerData `json:"player2"`
}
