package client

import (
	"net/http"
	"time"
)

type HTTPClient struct {
	Client *http.Client
}

func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *HTTPClient) IsLinkAccessible(url string) bool {
	// Checks the metadata like StatusCode.
	resp, err := c.Client.Head(url)

	if err != nil || resp.StatusCode >= 400 {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return false
		}
		req.Header.Set("User-Agent", "Mozilla/5.0")
		req.Header.Set("Accept", "text/html")

		resp, err = c.Client.Do(req)
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode < 400
	}

	defer resp.Body.Close()
	return resp.StatusCode < 400
}

func (c *HTTPClient) FetchResults(requestUrl string) (resp *http.Response, err error) {
	return c.Client.Get(requestUrl)
}
