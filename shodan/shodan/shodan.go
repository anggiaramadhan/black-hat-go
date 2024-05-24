package shodan

const baseUrl = "https://api.shodan.io"

type Client struct {
	apiKey string
}

func New(key string) *Client {
	return &Client{apiKey: key}
}
