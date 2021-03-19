package groovestats

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/archiveflax/gslauncher/internal/settings"
)

type Client struct {
	client         *http.Client
	permanentError bool
}

func NewClient() *Client {
	return &Client{
		client:         &http.Client{Timeout: 15 * time.Second},
		permanentError: false,
	}
}

func (client *Client) NewSession() (*NewSessionResponse, error) {
	if settings.Get().FakeGroovestats {
		return fakeNewSession()
	}

	req, err := client.newGetRequest("/new-session.php", nil)
	if err != nil {
		return nil, err
	}

	var response NewSessionResponse
	err = client.doRequest(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (client *Client) AutoSubmitScore(apiKey, hash string, rate int, score int) (*AutoSubmitScoreResponse, error) {
	if settings.Get().FakeGroovestats {
		return fakeAutoSubmitScore(hash, rate, score)
	}

	data := struct {
		Hash  string `json:"hash"`
		Rate  int    `json:"rate"`
		Score int    `json:"score"`
	}{hash, rate, score}

	req, err := client.newPostRequest("/auto-submit-score.php", &data)
	if err != nil {
		return nil, err
	}
	req.Header.Add("x-api-key", apiKey)

	var response AutoSubmitScoreResponse
	err = client.doRequest(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (client *Client) GetScores(apiKey, hash string) (*GetScoresResponse, error) {
	if settings.Get().FakeGroovestats {
		return fakeGetScores(hash)
	}

	params := url.Values{}
	params.Add("h", hash)

	req, err := client.newGetRequest("/get-scores.php", &params)
	if err != nil {
		return nil, err
	}
	req.Header.Add("x-api-key", apiKey)

	var response GetScoresResponse
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

func (client *Client) newPostRequest(path string, data interface{}) (*http.Request, error) {
	url := settings.Get().GrooveStatsUrl + path

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
