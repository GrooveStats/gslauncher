package groovestats

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/hashicorp/golang-lru"

	"github.com/GrooveStats/gslauncher/internal/fsipc"
	"github.com/GrooveStats/gslauncher/internal/settings"
	"github.com/GrooveStats/gslauncher/internal/stats"
)

type DisabledError struct {
	reason string
}

func (e *DisabledError) Error() string {
	return fmt.Sprintf("endpoint disabled: %s", e.reason)
}

type Client struct {
	getClient  *http.Client
	postClient *http.Client
	cache      *lru.Cache
	logger     *log.Logger

	allowScoreSubmit        bool
	allowPlayerScores       bool
	allowPlayerLeaderboards bool
	permanentError          bool
}

func NewClient() *Client {
	cache, _ := lru.New(512)
	logger := log.New(log.Writer(), "[GS] ", log.LstdFlags|log.Lmsgprefix)

	return &Client{
		getClient:  &http.Client{Timeout: 12 * time.Second},
		postClient: &http.Client{Timeout: 32 * time.Second},
		cache:      cache,
		logger:     logger,

		allowScoreSubmit:        false,
		allowPlayerScores:       false,
		allowPlayerLeaderboards: false,
		permanentError:          false,
	}
}

func (client *Client) NewSession(request *fsipc.GsNewSessionRequest) (*NewSessionResponse, error) {
	response, err := client.newSession(request)
	if err != nil {
		client.logger.Print(err)
	}
	return response, err
}

func (client *Client) newSession(request *fsipc.GsNewSessionRequest) (*NewSessionResponse, error) {
	stats.GsNewSessionCount++

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

	client.permanentError = false

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
	response, err := client.playerScores(request)
	if err != nil {
		client.logger.Print(err)
	}
	return response, err
}

func (client *Client) playerScores(request *fsipc.GsPlayerScoresRequest) (*PlayerScoresResponse, error) {
	cachedResponse := client.getPlayerScores(request)
	if cachedResponse != nil {
		stats.GsPlayerScoresCachedCount++
		return cachedResponse, nil
	}

	if !client.allowPlayerScores {
		return nil, &DisabledError{reason: "not allowed to fetch player scores"}
	}

	stats.GsPlayerScoresCount++

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

	client.addPlayerScores(request, &response)

	return &response, nil
}

func (client *Client) PlayerLeaderboards(request *fsipc.GsPlayerLeaderboardsRequest) (*PlayerLeaderboardsResponse, error) {
	response, err := client.playerLeaderboards(request)
	if err != nil {
		client.logger.Print(err)
	}
	return response, err
}

func (client *Client) playerLeaderboards(request *fsipc.GsPlayerLeaderboardsRequest) (*PlayerLeaderboardsResponse, error) {
	if !client.allowPlayerLeaderboards {
		return nil, &DisabledError{reason: "not allowed to fetch player leaderboards"}
	}

	stats.GsPlayerLeaderboardsCount++

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
	response, err := client.scoreSubmit(request)
	if err != nil {
		client.logger.Print(err)
	}
	return response, err
}

func (client *Client) scoreSubmit(request *fsipc.GsScoreSubmitRequest) (*ScoreSubmitResponse, error) {
	if !client.allowScoreSubmit {
		return nil, &DisabledError{reason: "not allowed to submit scores"}
	}

	stats.GsScoreSubmitCount++

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
		Score          int                   `json:"score"`
		Comment        string                `json:"comment,omitempty"`
		Rate           int                   `json:"rate"`
		JudgmentCounts *fsipc.JudgmentCounts `json:"judgmentCounts"`
		UsedCmod       *bool                 `json:"usedCmod"`
	}

	data := struct {
		Player1 *scoreSubmitPlayerData `json:"player1"`
		Player2 *scoreSubmitPlayerData `json:"player2"`
	}{nil, nil}

	if request.Player1 != nil {
		data.Player1 = &scoreSubmitPlayerData{
			Score:          request.Player1.Score,
			Comment:        request.Player1.Comment,
			Rate:           request.Player1.Rate,
			JudgmentCounts: request.Player1.JudgmentCounts,
			UsedCmod:       request.Player1.UsedCmod,
		}
	}

	if request.Player2 != nil {
		data.Player2 = &scoreSubmitPlayerData{
			Score:          request.Player2.Score,
			Comment:        request.Player2.Comment,
			Rate:           request.Player2.Rate,
			JudgmentCounts: request.Player1.JudgmentCounts,
			UsedCmod:       request.Player1.UsedCmod,
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

	client.removePlayerScores(request)

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
		return &DisabledError{reason: "request not sent due to protocol violation"}
	}

	var resp *http.Response
	var err error

	if req.Method == "GET" {
		resp, err = client.getClient.Do(req)
	} else {
		resp, err = client.postClient.Do(req)
	}
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	client.logger.Printf("%s %v %s", req.Method, req.URL, resp.Status)

	if resp.StatusCode >= 400 && resp.StatusCode < 499 && resp.StatusCode != 429 {
		client.permanentError = true
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Parse error response (if it actually is one)
	// XXX: This should only be done for 4xx and 5xx responses, but we have
	// to work around an issue in the API where error responses are
	// returned with status code 200
	errorResp := ErrorResponse{}
	json.Unmarshal(data, &errorResp)
	if errorResp.Error != "" || errorResp.Message != "" {
		if resp.StatusCode == 200 {
			client.permanentError = true
		}

		return fmt.Errorf("%s %v %d %v", req.Method, req.URL, resp.StatusCode, errorResp)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("status code %s", resp.Status)
	}

	return json.Unmarshal(data, response)
}
