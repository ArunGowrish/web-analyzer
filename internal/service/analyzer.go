package service

import (
	"errors"
	"log"
	"net/url"
	"strings"
	"sync"

	"golang.org/x/net/html"

	"github.com/ArunGowrish/web-analyzer/internal/client"
	"github.com/ArunGowrish/web-analyzer/internal/model"
	"github.com/ArunGowrish/web-analyzer/utils"
)

type AnalyzerService struct {
	HTTPClient client.HTTPClientInterface
}

func NewAnalyzerService(httpClient client.HTTPClientInterface) *AnalyzerService {
	return &AnalyzerService{
		HTTPClient: httpClient,
	}
}

// AnalyzeURL takes a URL string, validates it, fetches the page, parses the HTML,
// and returns an AnalysisResult containing the HTML version.
func (a *AnalyzerService) AnalyzeURL(requestUrl string) (*model.AnalysisResult, error) {
	log.Println("Analyze URL method invoked.")
	// Validate URL
	if msg := utils.IsUrlValid(requestUrl); msg != "" {
		return nil, errors.New(msg)
	}

	resp, err := a.HTTPClient.FetchResults(requestUrl)
	if err != nil {
		return nil, errors.New("Unable to fetch the webpage. Please check the URL or try again later.")
	}
	defer resp.Body.Close()

	// Parse HTML
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, errors.New("The webpage could not be processed. Please try a different URL.")
	}

	// Find the domain from the url
	baseURL, err := url.Parse(requestUrl)
	if err != nil {
		return nil, errors.New("The webpage could not be processed. Please try a different URL.")
	}

	htmlVersion := getHTMLVersion(doc)
	log.Println("HTML version of the web page - ", htmlVersion)

	title := getTitle(doc)
	log.Println("Title version of the web page - ", title)

	headingCounts := make(map[string]int)
	headings := getHeadingsAndCount(doc, headingCounts)
	log.Println("Headings of the web page - ", headings)

	link := &model.Link{}
	a.analyzeURL(doc, link, *baseURL)
	log.Println(
		"Internal links count from the web page - ", len(link.InternalLinks),
		"External links count from the web page - ", len(link.ExternalLinks),
		"In accessible links count from the web page - ", len(link.InAccessibleLinks),
	)

	isLoginForm := isPageContainsLoginForm(doc)
	log.Println("Is login form exist in web page - ", isLoginForm)

	result := &model.AnalysisResult{
		HTMLVersion: htmlVersion,
		Title:       title,
		Headings:    headings,
		Link:        *link,
		LoginForm:   isLoginForm,
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

// getHTMLVersion recursively traverses an HTML node tree to
// determine the Title of the document.
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

// getHeadingsAndCount recursively traverses an HTML node tree to
// determine the headings and counts.
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

// Analyze url method to check the url whether internal, external and
// it is accessible.
func (a *AnalyzerService) analyzeURL(n *html.Node, link *model.Link,
	baseURL url.URL) {

	// Extract all links
	a.extractLinks(n, link, baseURL)

	links := append(link.InternalLinks, link.ExternalLinks...)

	// Collect URLs
	var urls []string
	for _, l := range links {
		urls = append(urls, l)
	}

	// Run concurrent check
	results := a.checkLinksConcurrently(urls)

	// Append results
	for i, l := range links {
		if !results[l] {
			link.InAccessibleLinks = append(link.InAccessibleLinks, links[i])
		}
	}
}

// extractLinks recursively traverses an HTML node tree to
// determine the external links.
func (a *AnalyzerService) extractLinks(n *html.Node, link *model.Link,
	baseURL url.URL) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				href := attr.Val

				url, err := url.Parse(href)
				if err != nil {
					continue // skip -> no domain to find out internal or external link
				}

				if !isValidNavigationalLink(href) {
					continue // skip
				}

				isExternalLink := isExternalLink(href) && compareDomains(url.Host, baseURL.Host)

				// Treat not external links as internal links
				var finalUrl string
				if !isExternalLink {
					finalUrl = baseURL.ResolveReference(url).String()
					link.InternalLinks = append(link.InternalLinks, finalUrl)
				} else {
					finalUrl = href
					link.ExternalLinks = append(link.ExternalLinks, finalUrl)
				}
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		a.extractLinks(c, link, baseURL)
	}
}

func isExternalLink(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

func compareDomains(domainFromLink string, domainFromRequestUrl string) bool {
	return domainFromLink != "" && domainFromLink != domainFromRequestUrl
}

func isValidNavigationalLink(href string) bool {
	if href == "" {
		return false
	}

	hrefLower := strings.ToLower(href)

	if strings.HasPrefix(hrefLower, "javascript:") ||
		strings.HasPrefix(hrefLower, "#") ||
		strings.HasPrefix(hrefLower, "mailto:") ||
		strings.HasPrefix(hrefLower, "tel:") {
		return false
	}

	return true
}

// Check the accessibility of links concurrently by pool of goroutines.
func (a *AnalyzerService) checkLinksConcurrently(urls []string) map[string]bool {
	results := make(map[string]bool)
	var mu sync.Mutex

	channels := make(chan string, 100)
	var wg sync.WaitGroup
	verifiedAccessibility := make(map[string]bool)

	worker := func() {
		defer wg.Done()
		for url := range channels {
			ok := a.HTTPClient.IsLinkAccessible(url)

			mu.Lock()
			results[url] = ok
			mu.Unlock()
		}
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go worker()
	}
	for _, url := range urls {
		if !verifiedAccessibility[url] {
			verifiedAccessibility[url] = true
			channels <- url
		}
	}
	close(channels)

	wg.Wait()
	return results
}

// isPageContainsLoginForm traverses an HTML node tree to identify whether
// page has any Login form.
func isPageContainsLoginForm(n *html.Node) bool {
	if n.Type == html.ElementNode && n.Data == "form" {
		if isLoginForm(n) {
			return true
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if isPageContainsLoginForm(c) {
			return true
		}
	}
	return false
}

func isLoginForm(form *html.Node) bool {
	hasUsernameField, hasPassword := scanForLoginForm(form)
	return hasUsernameField && hasPassword
}

// Helper function to detect login form inside a <form> node
func scanForLoginForm(n *html.Node) (bool, bool) {
	hasUsernameField := false
	hasPassword := false

	if n.Type == html.ElementNode && n.Data == "input" {
		inputType := strings.ToLower(getAttr(n, "type"))
		switch inputType {
		case "text", "email":
			hasUsernameField = true
		case "password":
			hasPassword = true
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		username, password := scanForLoginForm(c)
		hasUsernameField = hasUsernameField || username
		hasPassword = hasPassword || password
	}
	return hasUsernameField, hasPassword
}

func getAttr(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}
