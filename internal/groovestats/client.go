package groovestats

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/archiveflax/gslauncher/internal/fsipc"
	"github.com/archiveflax/gslauncher/internal/settings"
)

type Client struct {
	client *http.Client

	allowScoreSubmit        bool
	allowPlayerScores       bool
	allowPlayerLeaderboards bool
	permanentError          bool
}

func NewClient() *Client {
	return &Client{
		client: &http.Client{Timeout: 15 * time.Second},

		allowScoreSubmit:        false,
		allowPlayerScores:       false,
		allowPlayerLeaderboards: false,
		permanentError:          false,
	}
}

func (client *Client) NewSession(request *fsipc.GsNewSessionRequest) (*NewSessionResponse, error) {
	if settings.Get().FakeGs {
		response, err := fakeNewSession()
		if err != nil {
			return nil, err
		}

		client.allowScoreSubmit = response.ServicesAllowed.ScoreSubmit
		client.allowPlayerScores = response.ServicesAllowed.PlayerScores
		client.allowPlayerLeaderboards = response.ServicesAllowed.PlayerLeaderboards

		return response, nil
	}

	params := url.Values{}
	params.Add("chartHashVersion", strconv.Itoa(request.ChartHashVersion))

	req, err := client.newGetRequest("/new-session.php", &params)
	if err != nil {
		return nil, err
	}

	var response NewSessionResponse
	err = client.doRequest(req, &response)
	if err != nil {
		return nil, err
	}

	client.allowScoreSubmit = response.ServicesAllowed.ScoreSubmit
	client.allowPlayerScores = response.ServicesAllowed.PlayerScores
	client.allowPlayerLeaderboards = response.ServicesAllowed.PlayerLeaderboards

	return &response, nil
}

func (client *Client) PlayerScores(request *fsipc.GsPlayerScoresRequest) (*PlayerScoresResponse, error) {
	if !client.allowPlayerScores {
		return nil, errors.New("not allowed to fetch player scores")
	}

	if settings.Get().FakeGs {
		return fakePlayerScores(request)
	}

	params := url.Values{}
	if request.Player1 != nil {
		params.Add("chartHashP1", request.Player1.ChartHash)
	}
	if request.Player2 != nil {
		params.Add("chartHashP2", request.Player2.ChartHash)
	}

	req, err := client.newGetRequest("/player-scores.php", &params)
	if err != nil {
		return nil, err
	}
	if request.Player1 != nil {
		req.Header.Add("x-api-key-player-1", request.Player1.ApiKey)
	}
	if request.Player2 != nil {
		req.Header.Add("x-api-key-player-2", request.Player2.ApiKey)
	}

	var response PlayerScoresResponse
	err = client.doRequest(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (client *Client) PlayerLeaderboards(request *fsipc.GsPlayerLeaderboardsRequest) (*PlayerLeaderboardsResponse, error) {
	if !client.allowPlayerLeaderboards {
		return nil, errors.New("not allowed to fetch player leaderboards")
	}

	if settings.Get().FakeGs {
		return fakePlayerLeaderboards(request)
	}

	params := url.Values{}
	if request.Player1 != nil {
		params.Add("chartHashP1", request.Player1.ChartHash)
	}
	if request.Player2 != nil {
		params.Add("chartHashP2", request.Player2.ChartHash)
	}
	if request.MaxLeaderboardResults != nil {
		params.Add("maxLeaderboardResults", strconv.Itoa(*request.MaxLeaderboardResults))
	}

	req, err := client.newGetRequest("/player-leaderboards.php", &params)
	if err != nil {
		return nil, err
	}
	if request.Player1 != nil {
		req.Header.Add("x-api-key-player-1", request.Player1.ApiKey)
	}
	if request.Player2 != nil {
		req.Header.Add("x-api-key-player-2", request.Player2.ApiKey)
	}

	var response PlayerLeaderboardsResponse
	err = client.doRequest(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (client *Client) ScoreSubmit(request *fsipc.GsScoreSubmitRequest) (*ScoreSubmitResponse, error) {
	if !client.allowScoreSubmit {
		return nil, errors.New("not allowed to submit scores")
	}

	if settings.Get().FakeGs {
		return fakeScoreSubmit(request)
	}

	params := url.Values{}
	if request.Player1 != nil {
		params.Add("chartHashP1", request.Player1.ChartHash)
	}
	if request.Player2 != nil {
		params.Add("chartHashP2", request.Player2.ChartHash)
	}
	if request.MaxLeaderboardResults != nil {
		params.Add("maxLeaderboardResults", strconv.Itoa(*request.MaxLeaderboardResults))
	}

	type scoreSubmitPlayerData struct {
		Score   int    `json:"score"`
		Comment string `json:"comment"`
		Rate    int    `json:"rate"`
	}

	data := struct {
		player1 *scoreSubmitPlayerData
		player2 *scoreSubmitPlayerData
	}{nil, nil}

	if request.Player1 != nil {
		data.player1 = &scoreSubmitPlayerData{
			Score:   request.Player1.Score,
			Comment: request.Player1.Comment,
			Rate:    request.Player1.Rate,
		}
	}

	if request.Player2 != nil {
		data.player2 = &scoreSubmitPlayerData{
			Score:   request.Player2.Score,
			Comment: request.Player2.Comment,
			Rate:    request.Player2.Rate,
		}
	}

	req, err := client.newPostRequest("/score-submit.php", &params, &data)
	if err != nil {
		return nil, err
	}

	if request.Player1 != nil {
		req.Header.Add("x-api-key-player-1", request.Player1.ApiKey)
	}

	if request.Player2 != nil {
		req.Header.Add("x-api-key-player-2", request.Player2.ApiKey)
	}

	var response ScoreSubmitResponse
	err = client.doRequest(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (client *Client) newGetRequest(path string, params *url.Values) (*http.Request, error) {
	url := settings.Get().GrooveStatsUrl + path
	if params != nil {
		url += "?" + params.Encode()
	}

	return http.NewRequest("GET", url, nil)
}

func (client *Client) newPostRequest(path string, params *url.Values, data interface{}) (*http.Request, error) {
	url := settings.Get().GrooveStatsUrl + path
	if params != nil {
		url += "?" + params.Encode()
	}

	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return http.NewRequest("POST", url, bytes.NewBuffer(body))
}

func (client *Client) doRequest(req *http.Request, response interface{}) error {
	if client.permanentError {
		return errors.New("request not sent due to protocol violation")
	}

	resp, err := client.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 && resp.StatusCode < 499 && resp.StatusCode != 429 {
		client.permanentError = true
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("status code %s", resp.Status)
	}

	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(response)
}
