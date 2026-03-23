package internal

import (
	"html/template"
	"log"
	"net/http"

	"github.com/ArunGowrish/web-analyzer/internal/model"
)

type Handler struct {
	Tmpl     *template.Template
	Analyzer serviceAnalyzerInterface
}

type serviceAnalyzerInterface interface {
	AnalyzeURL(url string) (*model.AnalysisResult, error)
}

// HomeHandler renders the home page template.
func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpl.Execute(w, nil)
	if err != nil {
		log.Println("Template execution error:", err)
		http.Error(w, "Something went wrong. Please try again.", http.StatusInternalServerError)
	}
}

// SubmitHandler processes the form submission containing a URL.
// Template execution errors are logged and returned as 500 Internal Server Error.
func (h *Handler) SubmitHandler(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	log.Println("Received URL:", url)

	result, err := h.Analyzer.AnalyzeURL(url)
	if err != nil {
		_ = h.Tmpl.Execute(w, map[string]string{
			"Error": err.Error(),
		})
		return
	}

	if err := h.Tmpl.Execute(w, map[string]interface{}{
		"HTMLVersion":              result.HTMLVersion,
		"Title":                    result.Title,
		"Headings":                 result.Headings,
		"ExternalLinksCount":       len(result.Link.ExternalLinks),
		"InternalLinksCount":       len(result.Link.InternalLinks),
		"InAssecessibleLinksCount": len(result.Link.InAccessibleLinks),
	}); err != nil {
		log.Println("Template execution error:", err)
		http.Error(w, "Something went wrong. Please try again.", http.StatusInternalServerError)
		return
	}
}
