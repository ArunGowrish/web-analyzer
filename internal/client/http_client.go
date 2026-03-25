package client

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

type HTTPClient struct {
	Client *http.Client
}

type HTTPClientInterface interface {
	FetchResults(url string) (*http.Response, error)
	IsLinkAccessible(url string) bool
}

func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *HTTPClient) IsLinkAccessible(url string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	// HEAD request - Checks the metadata like StatusCode.
	req, err := http.NewRequestWithContext(ctx, "HEAD", url, nil)
	if err != nil {
		return false
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Println("Timeout when (HEAD): ", url)
		}
		return false
	}
	// GET request - Fallback to GET request if StatusCode represent fail.
	if resp.StatusCode >= 400 {
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

		if err != nil {
			return false
		}
		c.setCommonHeaders(req)

		resp, err = c.Client.Do(req)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				log.Println("Timeout when (GET): ", url)
			}
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode < 400
	}

	defer resp.Body.Close()
	return resp.StatusCode < 400
}

func (c *HTTPClient) FetchResults(requestUrl string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return nil, err
	}

	c.setCommonHeaders(req)
	response, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP error: Failed to fetch the webpage.")
	}
	return response, nil
}

func (c *HTTPClient) setCommonHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/120 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")
}
