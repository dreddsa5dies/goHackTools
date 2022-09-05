package shodan

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// BaseURL - ...
const BaseURL = "https://api.shodan.io"

// Client - ...
type Client struct {
	apiKey string
}

// New - ...
func New(apiKey string) *Client {
	return &Client{apiKey: apiKey}
}

// APIInfo - information of you account status
type APIInfo struct {
	QueryCredits int    `json:"query_credits"`
	ScanCredits  int    `json:"scan_credits"`
	Telnet       bool   `json:"telnet"`
	Plan         string `json:"plan"`
	HTTPS        bool   `json:"https"`
	Unlocked     bool   `json:"unlocked"`
}

// APIInfo - decoding respons from SHODAN.IO to APIInfo struct
func (s *Client) APIInfo() (*APIInfo, error) {
	res, err := http.Get(fmt.Sprintf("%s/api-info?key=%s", BaseURL, s.apiKey))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var ret APIInfo

	err = json.Unmarshal(body, &ret)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}
