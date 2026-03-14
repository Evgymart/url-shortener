package redirect

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	resp "urlshort/internal/http-server/handlers/api/response"
	"urlshort/internal/storage"
	storagemock "urlshort/tests/internal/storage"
)

func TestRedirectSuccess(t *testing.T) {
	mockStorage := storagemock.NewMockURLSaver()
	handler := SetupHandler(mockStorage)

	alias := "google"
	url := "https://google.com"
	mockStorage.SavedURLs[alias] = url

	req, _ := CreateTestRequest(http.MethodGet, "/"+alias, alias)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusTemporaryRedirect {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusTemporaryRedirect)
	}

	location := rr.Header().Get("Location")
	if location != url {
		t.Errorf("expected location %s, got %s", url, location)
	}
}

func TestRedirectNotFound(t *testing.T) {
	mockStorage := storagemock.NewMockURLSaver()
	handler := SetupHandler(mockStorage)

	alias := "nonexistent"

	req, _ := CreateTestRequest(http.MethodGet, "/"+alias, alias)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	var response resp.Response
	err := json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if response.Status != resp.StatusError {
		t.Errorf("expected status %s, got %s", resp.StatusError, response.Status)
	}

	if response.Error != storage.ErrUrlNotFound.Error() {
		t.Errorf("expected error %s, got %s", storage.ErrUrlNotFound.Error(), response.Error)
	}
}

func TestRedirectMultiple(t *testing.T) {
	mockStorage := storagemock.NewMockURLSaver()
	handler := SetupHandler(mockStorage)

	testCases := []struct {
		name  string
		alias string
		url   string
	}{
		{"google", "google", "https://google.com"},
		{"github", "github", "https://github.com"},
		{"example", "example", "https://example.com"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockStorage.SavedURLs[tc.alias] = tc.url

			req, _ := CreateTestRequest(http.MethodGet, "/"+tc.alias, tc.alias)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusTemporaryRedirect {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusTemporaryRedirect)
			}

			location := rr.Header().Get("Location")
			if location != tc.url {
				t.Errorf("expected location %s, got %s", tc.url, location)
			}
		})
	}
}
