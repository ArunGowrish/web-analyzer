package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchResults_Success(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	client := NewHTTPClient()

	resp, err := client.FetchResults(mockServer.URL)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestIsLinkAccessible_HeadSuccess(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodHead {
			w.WriteHeader(http.StatusOK)
			return
		}
	}))
	defer mockServer.Close()

	client := NewHTTPClient()

	isLinkAccessible := client.IsLinkAccessible(mockServer.URL)
	if !isLinkAccessible {
		t.Errorf("expected true, got false")
	}
}

func TestIsLinkAccessible_FallbackToGet(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodHead {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			return
		}
	}))
	defer mockServer.Close()

	client := NewHTTPClient()

	isLinkAccessible := client.IsLinkAccessible(mockServer.URL)
	if !isLinkAccessible {
		t.Errorf("expected true from GET fallback, got false")
	}
}

func TestIsLinkAccessible_Failure(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	client := NewHTTPClient()

	isLinkAccessible := client.IsLinkAccessible(mockServer.URL)
	if isLinkAccessible {
		t.Errorf("expected false, got true")
	}
}
