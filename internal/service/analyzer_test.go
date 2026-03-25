package service

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/ArunGowrish/web-analyzer/internal/mocks"
	"github.com/ArunGowrish/web-analyzer/utils"
)

func TestAnalyzeURL_ValidHTML(t *testing.T) {
	htmlContent := `<!DOCTYPE html><html><head><title>Test</title></head><body></body></html>`

	// Mock HTTP client
	mockClient := &mocks.HTTPClientMock{
		FetchResultsFunc: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(htmlContent)),
				Header:     make(http.Header),
			}, nil
		},
		IsLinkAccessibleFunc: func(url string) bool {
			return true
		},
	}

	service := NewAnalyzerService(mockClient)

	result, err := service.AnalyzeURL("http://mock.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.HTMLVersion != utils.IdentifyHTMLVersion("html") {
		t.Errorf("expected HTML version %s, got %s", utils.IdentifyHTMLVersion("html"), result.HTMLVersion)
	}
}

func TestAnalyzeURL_InvalidURL(t *testing.T) {
	mockClient := &mocks.HTTPClientMock{}
	service := NewAnalyzerService(mockClient)

	_, err := service.AnalyzeURL("invalid-url")
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}

func TestAnalyzeURL_HTTPError(t *testing.T) {
	mockClient := &mocks.HTTPClientMock{
		FetchResultsFunc: func(url string) (*http.Response, error) {
			return nil, errors.New("network error")
		},
	}

	service := NewAnalyzerService(mockClient)

	_, err := service.AnalyzeURL("http://mock.com")
	if err == nil || !strings.Contains(err.Error(), "Unable to fetch the webpage. Please check the URL or try again later.") {
		t.Fatalf("expected fetch error, got %v", err)
	}
}

func TestAnalyzeURL_InvalidHTML(t *testing.T) {
	mockClient := &mocks.HTTPClientMock{
		FetchResultsFunc: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader("<html></ht>")),
				Header:     make(http.Header),
			}, nil
		},
	}

	service := NewAnalyzerService(mockClient)

	response, err := service.AnalyzeURL("http://mock.com")
	htmlVersion := response.HTMLVersion
	if err == nil && htmlVersion != "HTML5" {
		t.Fatalf("expected HTML version HTML5, got %v", htmlVersion)
	}
}

func TestAnalyzeURL_ValidTitle(t *testing.T) {
	mockClient := &mocks.HTTPClientMock{
		FetchResultsFunc: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader("<html><head><title>Test title</title></head></html>")),
				Header:     make(http.Header),
			}, nil
		},
		IsLinkAccessibleFunc: func(url string) bool {
			return true
		},
	}

	service := NewAnalyzerService(mockClient)
	response, err := service.AnalyzeURL("http://mock.com")
	title := response.Title
	if err == nil && title == "" {
		t.Fatalf("expected Title from parsed html, got %v", title)
	}
}

func TestAnalyzeURL_ExtractHeadingsAndCounts(t *testing.T) {
	htmlContent := `<!DOCTYPE html>
	<html><body>
	<h1></h1>
	<h2></h2><h2></h2>
	<h3></h3><h3></h3><h3></h3>	<h3></h3>
	<h4></h4><h4></h4><h4></h4><h4></h4>
	<h5></h5><h5></h5><h5></h5>
	<h6></h6><h6></h6>
	</body></html>`

	mockClient := &mocks.HTTPClientMock{
		FetchResultsFunc: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(htmlContent)),
				Header:     make(http.Header),
			}, nil
		},
		IsLinkAccessibleFunc: func(url string) bool { return true },
	}

	service := NewAnalyzerService(mockClient)
	response, _ := service.AnalyzeURL("http://mock.com")

	mockHeadingsCountMap := map[string]int{
		"h1": 1,
		"h2": 2,
		"h3": 4,
		"h4": 4,
		"h5": 3,
		"h6": 2,
	}
	for k, v := range response.Headings {
		headingsCountMapValue := mockHeadingsCountMap[k]
		if v != headingsCountMapValue {
			t.Fatalf("expected headings and count from parsed html %s: %v, got %s: %v", k, v, k, headingsCountMapValue)
		}
	}
}

func TestAnalyzeURL_ValidateNavigationLink(t *testing.T) {
	htmlContent := `<!DOCTYPE html>
	<html><body>
		<a href="javascript:"></a>
		<a href="#"></a>
	</body></html>`

	mockClient := &mocks.HTTPClientMock{
		FetchResultsFunc: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(htmlContent)),
				Header:     make(http.Header),
			}, nil
		},
		IsLinkAccessibleFunc: func(url string) bool {
			return true
		},
	}

	service := NewAnalyzerService(mockClient)
	response, err := service.AnalyzeURL("http://mock.com")
	link := response.Link
	internalLinks := len(link.InternalLinks)
	externalLinks := len(link.ExternalLinks)
	if err == nil && internalLinks != 0 && externalLinks != 0 {
		t.Fatalf("expected internal and external links count from url %d , got %v", 0, internalLinks+externalLinks)
	}
}

func TestAnalyzeURL_ValidateInternalLinksCount(t *testing.T) {
	htmlContent := `<!DOCTYPE html>
	<html><body>
		<a href="http://mock.com/about"></a>
		<a href="/inventory"></a>
	</body></html>`

	mockClient := &mocks.HTTPClientMock{
		FetchResultsFunc: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(htmlContent)),
				Header:     make(http.Header),
			}, nil
		},
		IsLinkAccessibleFunc: func(url string) bool {
			return true
		},
	}

	service := NewAnalyzerService(mockClient)
	response, err := service.AnalyzeURL("http://mock.com")
	link := response.Link
	linkCount := len(link.InternalLinks)
	if err == nil && len(link.InternalLinks) != 2 {
		t.Fatalf("expected internal links count from url %d , got %v", 2, linkCount)
	}
}

func TestAnalyzeURL_ValidateExternalLinksCount(t *testing.T) {
	htmlContent := `<!DOCTYPE html>
	<html><body>
		<a href="http://external.com/data"></a>
		<a href="http://external.com/address"></a>
	</body></html>`

	mockClient := &mocks.HTTPClientMock{
		FetchResultsFunc: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(htmlContent)),
				Header:     make(http.Header),
			}, nil
		},
		IsLinkAccessibleFunc: func(url string) bool {
			return true
		},
	}

	service := NewAnalyzerService(mockClient)
	response, err := service.AnalyzeURL("http://mock.com")
	link := response.Link
	linkCount := len(link.ExternalLinks)
	if err == nil && len(link.ExternalLinks) != 2 {
		t.Fatalf("expected internal links count from url %d , got %v", 2, linkCount)
	}
}

func TestAnalyzeURL_ValidateInAccesibleLinksCount(t *testing.T) {
	htmlContent := `<!DOCTYPE html>
	<html><body>
		<a href="http://xxx.com/"></a>
		<a href="http://yyy.com/"></a>
	</body></html>`

	mockClient := &mocks.HTTPClientMock{
		FetchResultsFunc: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(htmlContent)),
				Header:     make(http.Header),
			}, nil
		},
		IsLinkAccessibleFunc: func(url string) bool {
			return false
		},
	}

	service := NewAnalyzerService(mockClient)
	response, err := service.AnalyzeURL("http://mock.com")
	link := response.Link
	linkCount := len(link.InAccessibleLinks)
	if err == nil && linkCount != 2 {
		t.Fatalf("expected in accessible links count from url %d , got %v", 2, linkCount)
	}
}

func TestAnalyzeURL_LoginFormNotExist(t *testing.T) {
	htmlContent := `<!DOCTYPE html>
	<html><body>
		<form>
			<input type="text" name="username">
		</form>
	</body></html>`

	mockClient := &mocks.HTTPClientMock{
		FetchResultsFunc: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(htmlContent)),
				Header:     make(http.Header),
			}, nil
		},
		IsLinkAccessibleFunc: func(url string) bool {
			return true
		},
	}

	service := NewAnalyzerService(mockClient)
	response, err := service.AnalyzeURL("http://mock.com")
	isLoginFormExist := response.LoginForm
	if err == nil && isLoginFormExist {
		t.Fatalf("expected login form not exist in the webpage, got %v", isLoginFormExist)
	}
}

func TestAnalyzeURL_LoginFormExist(t *testing.T) {
	htmlContent := `<!DOCTYPE html>
	<html><body>
		<form>
			<input type="text" name="username">
			<input type="password" name="password">
		</form>
	</body></html>`

	mockClient := &mocks.HTTPClientMock{
		FetchResultsFunc: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(htmlContent)),
				Header:     make(http.Header),
			}, nil
		},
		IsLinkAccessibleFunc: func(url string) bool {
			return true
		},
	}

	service := NewAnalyzerService(mockClient)
	response, err := service.AnalyzeURL("http://mock.com")
	isLoginFormExist := response.LoginForm
	if err == nil && !isLoginFormExist {
		t.Fatalf("expected a login form in the webpage, got %v", isLoginFormExist)
	}
}
