package service

import (
	"errors"
	"net/http"

	"golang.org/x/net/html"

	"github.com/ArunGowrish/web-analyzer/internal/model"
	"github.com/ArunGowrish/web-analyzer/utils"
)

type AnalyzerService struct {
	HTTPGet func(url string) (*http.Response, error)
}

// AnalyzeURL takes a URL string, validates it, fetches the page, parses the HTML,
// and returns an AnalysisResult containing the HTML version.
func (s *AnalyzerService) AnalyzeURL(url string) (*model.AnalysisResult, error) {
	// Validate URL
	if msg := utils.IsUrlValid(url); msg != "" {
		return nil, errors.New(msg)
	}

	httpGetFunc := s.HTTPGet
	if httpGetFunc == nil {
		httpGetFunc = http.Get
	}

	// Fetch the URL
	resp, err := httpGetFunc(url)
	if err != nil {
		return nil, errors.New("failed to fetch URL")
	}
	defer resp.Body.Close()

	// Parse HTML
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, errors.New("failed to parse HTML")
	}

	HeadingCounts := make(map[string]int)

	// Extract values
	result := &model.AnalysisResult{
		HTMLVersion: getHTMLVersion(doc),
		Title:       getTitle(doc),
		Headings:    getHeadingsAndCount(doc, HeadingCounts),
	}

	return result, nil
}

// getHTMLVersion recursively traverses an HTML node tree to identify the DOCTYPE
// and determine the HTML version.
func getHTMLVersion(n *html.Node) string {
	if n.Type == html.DoctypeNode {
		return n.Data
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if v := getHTMLVersion(c); v != "" {
			return utils.IdentifyHTMLVersion(v)
		}
	}
	return "unknown"
}

// getHTMLVersion recursively traverses an HTML node tree to identify the DOCTYPE
// and determine the Title of the document.
func getTitle(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
		return n.FirstChild.Data
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if title := getTitle(c); title != "" {
			return title
		}
	}
	return ""
}

// getHeadingsAndCount recursively traverses an HTML node tree to identify the DOCTYPE
// and determine the headings and counts.
func getHeadingsAndCount(n *html.Node, counts map[string]int) map[string]int {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "h1", "h2", "h3", "h4", "h5", "h6":
			counts[n.Data]++
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		getHeadingsAndCount(c, counts)
	}
	return counts
}
