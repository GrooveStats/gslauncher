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
	baseUrl        string
	apiKey         string
	client         *http.Client
	permanentError bool
}

func NewClient(apiKey string) *Client {
	baseUrl := settings.Get().GrooveStatsUrl

	return &Client{
		baseUrl:        baseUrl,
		apiKey:         apiKey,
		client:         &http.Client{Timeout: 15 * time.Second},
		permanentError: false,
	}
}

func (client *Client) AutoSubmitScore(hash string, rate int, score int) (*AutoSubmitScoreResponse, error) {
	if settings.Get().FakeGroovestats {
		return fakeAutoSubmitScore(hash, rate, score)
	}

	data := struct {
		Hash  string `json:"hash"`
		Rate  int    `json:"rate"`
		Score int    `json:"score"`
	}{hash, rate, score}

	var response AutoSubmitScoreResponse

	err := client.jsonPost("/auto-submit-score.php", &data, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (client *Client) GetScores(hash string) (*GetScoresResponse, error) {
	if settings.Get().FakeGroovestats {
		return fakeGetScores(hash)
	}

	params := url.Values{}
	params.Add("h", hash)

	var response GetScoresResponse

	err := client.jsonGet("/get-scores.php", &params, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (client *Client) jsonGet(path string, params *url.Values, response interface{}) error {
	url := client.baseUrl + path
	if params != nil {
		url += "?" + params.Encode()
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("x-api-key", client.apiKey)

	return client.doRequest(req, response)
}

func (client *Client) jsonPost(path string, data interface{}, response interface{}) error {
	url := client.baseUrl + path

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Add("x-api-key", client.apiKey)

	return client.doRequest(req, response)
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
