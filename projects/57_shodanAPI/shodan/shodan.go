package shodan

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
