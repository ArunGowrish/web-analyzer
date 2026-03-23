package utils

import (
	"net/url"
	"strings"
)

func IsUrlValid(urlFromRequest string) string {
	trimmedUrl := strings.TrimSpace(urlFromRequest)
	if trimmedUrl == "" {
		return "Url cannot be empty"
	}
	parsedRequestUrl, err := url.ParseRequestURI(trimmedUrl)
	if err != nil || parsedRequestUrl.Scheme == "" || parsedRequestUrl.Host == "" {
		return "Invalid url"
	}
	return ""
}

func IdentifyHTMLVersion(doctype string) string {
	d := strings.ToLower(doctype)

	switch {
	case strings.Contains(d, "html 4.01") && strings.Contains(d, "strict"):
		return "HTML 4.01 Strict"

	case strings.Contains(d, "html 4.01") && strings.Contains(d, "transitional"):
		return "HTML 4.01 Transitional"

	case strings.Contains(d, "html 4.01") && strings.Contains(d, "frameset"):
		return "HTML 4.01 Frameset"

	case strings.Contains(d, "xhtml 1.0") && strings.Contains(d, "strict"):
		return "XHTML 1.0 Strict"

	case strings.Contains(d, "xhtml 1.0") && strings.Contains(d, "transitional"):
		return "XHTML 1.0 Transitional"

	case strings.Contains(d, "xhtml 1.0") && strings.Contains(d, "frameset"):
		return "XHTML 1.0 Frameset"

	case strings.Contains(d, "xhtml 1.1"):
		return "XHTML 1.1"

	default:
		return "HTML5"
	}
}
