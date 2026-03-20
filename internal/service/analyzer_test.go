package service

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/ArunGowrish/web-analyzer/utils"
)

func TestAnalyzeURL_ValidHTML(t *testing.T) {
	htmlContent := `<!DOCTYPE html><html><head><title>Test</title></head><body></body></html>`

	service := &AnalyzerService{
		HTTPGet: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(htmlContent)),
				Header:     make(http.Header),
			}, nil
		},
	}

	result, err := service.AnalyzeURL("http://example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.HTMLVersion != utils.IdentifyHTMLVersion("html") {
		t.Errorf("expected HTML version %s, got %s", utils.IdentifyHTMLVersion("html"), result.HTMLVersion)
	}
}

func TestAnalyzeURL_InvalidURL(t *testing.T) {
	service := &AnalyzerService{}

	_, err := service.AnalyzeURL("invalid-url")
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}

func TestAnalyzeURL_HTTPError(t *testing.T) {
	service := &AnalyzerService{
		HTTPGet: func(url string) (*http.Response, error) {
			return nil, errors.New("network error")
		},
	}

	_, err := service.AnalyzeURL("http://example.com")
	if err == nil || !strings.Contains(err.Error(), "failed to fetch URL") {
		t.Fatalf("expected fetch error, got %v", err)
	}
}

func TestAnalyzeURL_InvalidHTML(t *testing.T) {
	service := &AnalyzerService{
		HTTPGet: func(url string) (*http.Response, error) {
			// Simulate invalid HTML which fails parsing
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader("<html></ht>")),
				Header:     make(http.Header),
			}, nil
		},
	}

	response, err := service.AnalyzeURL("http://example.com")
	HTMLVersion := response.HTMLVersion
	if err == nil && HTMLVersion != "Unknown" {
		t.Fatalf("expected HTML version Unknown, got %v", HTMLVersion)
	}
}

func TestAnalyzeURL_ValidTitle(t *testing.T) {
	service := &AnalyzerService{
		HTTPGet: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader("<html><head><title>Test title</title></head></html>")),
				Header:     make(http.Header),
			}, nil
		},
	}

	response, err := service.AnalyzeURL("http://example.com")
	Title := response.Title
	if err == nil && Title == "" {
		t.Fatalf("expected Title from parsed html, got %v", Title)
	}
}
