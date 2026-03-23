package mocks

import (
	"io"
	"net/http"
	"strings"
)

// HTTPClientMock is a mock for AnalyzerService's HTTPClient
type HTTPClientMock struct {
	FetchResultsFunc     func(url string) (*http.Response, error)
	IsLinkAccessibleFunc func(url string) bool
}

func (m *HTTPClientMock) FetchResults(url string) (*http.Response, error) {
	if m.FetchResultsFunc != nil {
		return m.FetchResultsFunc(url)
	}
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader("<html></html>")),
		Header:     make(http.Header),
	}, nil
}

func (m *HTTPClientMock) IsLinkAccessible(url string) bool {
	if m.IsLinkAccessibleFunc != nil {
		return m.IsLinkAccessibleFunc(url)
	}
	return true
}
