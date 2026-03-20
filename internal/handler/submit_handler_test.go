package internal

import (
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ArunGowrish/web-analyzer/internal/model"
)

type MockAnalyzerService struct {
	Result *model.AnalysisResult
	Err    error
}

func (m *MockAnalyzerService) AnalyzeURL(url string) (*model.AnalysisResult, error) {
	return m.Result, m.Err
}

func TestHomeHandler(t *testing.T) {
	tmpl := template.Must(template.New("mock").Parse("mock template executed"))
	h := &Handler{Tmpl: tmpl}

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	h.HomeHandler(w, req)

	resp := w.Result()
	body := w.Body.String()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	if !strings.Contains(body, "mock template executed") {
		t.Errorf("expected template output, got %s", body)
	}
}

func TestSubmitHandler_ValidURL(t *testing.T) {
	tmpl := template.Must(template.New("mock").Parse("HTML Version: {{.HTMLVersion}}"))

	h := &Handler{
		Tmpl: tmpl,
		Analyzer: &MockAnalyzerService{
			Result: &model.AnalysisResult{
				HTMLVersion: "HTML5",
			},
			Err: nil,
		},
	}

	form := strings.NewReader("url=http://example.com")
	req := httptest.NewRequest("POST", "/submit", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	h.SubmitHandler(w, req)

	resp := w.Result()
	body := w.Body.String()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", resp.StatusCode)
	}
	if !strings.Contains(body, "HTML5") {
		t.Errorf("expected HTML version HTML5 in response, got %s", body)
	}
}

func TestSubmitHandler_ServiceError(t *testing.T) {
	tmpl := template.Must(template.New("mock").Parse("Error: {{.Error}}"))

	h := &Handler{
		Tmpl: tmpl,
		Analyzer: &MockAnalyzerService{
			Result: nil,
			Err:    fmt.Errorf("invalid URL"),
		},
	}

	form := strings.NewReader("url=http://example.com")
	req := httptest.NewRequest("POST", "/submit", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	h.SubmitHandler(w, req)

	body := w.Body.String()
	if !strings.Contains(body, "invalid URL") {
		t.Errorf("expected error message in response, got %s", body)
	}
}
