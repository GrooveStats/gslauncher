package groovestats

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
)

type Client struct {
	baseUrl string
	apiKey  string
	client  *http.Client
}

func NewClient(baseUrl string, apiKey string) *Client {
	return &Client{
		baseUrl: baseUrl,
		apiKey:  apiKey,
		client:  &http.Client{},
	}
}

func (client *Client) AutoSubmitScore(hash string, rate int, score int) (*AutoSubmitScoreResponse, error) {
	data := struct {
		H string `json:"h"`
		R int    `json:"r"`
		S int    `json:"s"`
	}{hash, rate, score}

	var response AutoSubmitScoreResponse

	err := client.jsonPost("/auto-submit-score.php", &data, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (client *Client) GetScores(hash string) (*GetScoresResponse, error) {
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

	resp, err := client.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(response)
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

	resp, err := client.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(response)
}
