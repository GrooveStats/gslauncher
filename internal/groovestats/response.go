package groovestats

type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

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

	ServicesResult string `json:"servicesResult"`
}

type leaderBoardEntry struct {
	Name       string  `json:"name"`
	MachineTag *string `json:"machineTag"`
	Score      int     `json:"score"`
	Date       string  `json:"date"`
	Rank       int     `json:"rank"`
	IsSelf     bool    `json:"isSelf"`
	IsRival    bool    `json:"isRival"`
	IsFail     bool    `json:"isFail"`
	Comments   *string `json:"comments"`
}

type playerScoresPlayerData struct {
	ChartHash     string              `json:"chartHash"`
	IsRanked      bool                `json:"isRanked"`
	GsLeaderboard *[]leaderBoardEntry `json:"gsLeaderboard"`
}

type PlayerScoresResponse struct {
	Player1 *playerScoresPlayerData `json:"player1"`
	Player2 *playerScoresPlayerData `json:"player2"`

	// added by the launcher
	Cached bool `json:"cached"`
}

type playerLeaderboardsPlayerData struct {
	ChartHash     string             `json:"chartHash"`
	IsRanked      bool               `json:"isRanked"`
	GsLeaderboard []leaderBoardEntry `json:"gsLeaderboard"`

	Rpg *struct {
		Name           string             `json:"name"`
		RpgLeaderboard []leaderBoardEntry `json:"rpgLeaderboard"`
	} `json:"rpg"`

	Itl *struct {
		Name           string             `json:"name"`
		ItlLeaderboard []leaderBoardEntry `json:"itlLeaderboard"`
	} `json:"itl"`
}

type PlayerLeaderboardsResponse struct {
	Player1 *playerLeaderboardsPlayerData `json:"player1"`
	Player2 *playerLeaderboardsPlayerData `json:"player2"`
}

type scoreSubmitPlayerData struct {
	ChartHash     string              `json:"chartHash"`
	IsRanked      bool                `json:"isRanked"`
	Result        string              `json:"result,omitempty"`
	ScoreDelta    *int                `json:"scoreDelta,omitempty"`
	GsLeaderboard *[]leaderBoardEntry `json:"gsLeaderboard"`

	Rpg *struct {
		Name       string `json:"name"`
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
					Type        string `json:"type"`
					Description string `json:"description"`
				} `json:"rewards"`
				SongDownloadUrl *string `json:"songDownloadUrl"`
			} `json:"questsCompleted"`
		} `json:"progress"`

		RpgLeaderboard *[]leaderBoardEntry `json:"rpgLeaderboard"`
	} `json:"rpg"`

	Itl *struct {
		Name           string              `json:"name"`
		ScoreDelta     *int                `json:"scoreDelta,omitempty"`
		ItlLeaderboard *[]leaderBoardEntry `json:"itlLeaderboard"`

		Progress *struct {
			StatImprovements []struct {
				Name    string `json:"name"`
				Current int    `json:"current"`
				Gained  int    `json:"gained"`
			} `json:"statImprovements"`

			QuestsCompleted []struct {
				Title   string `json:"title"`
				Rewards []struct {
					Type        string `json:"type"`
					Description string `json:"description"`
				} `json:"rewards"`
				SongDownloadUrl *string `json:"songDownloadUrl"`
			} `json:"questsCompleted"`
		} `json:"progress"`
	} `json:"itl"`
}

type ScoreSubmitResponse struct {
	Player1 *scoreSubmitPlayerData `json:"player1"`
	Player2 *scoreSubmitPlayerData `json:"player2"`
}
