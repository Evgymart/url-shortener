package redirect

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"urlshort/internal/storage"
	storagemock "urlshort/tests/internal/storage"
)

func TestRedirectSuccess(t *testing.T) {
	mockStorage := storagemock.NewMockURLSaver()
	handler := SetupHandler(mockStorage)

	url := "https://example.com"
	alias := "test-alias"
	mockStorage.SavedURLs[alias] = url

	req, _ := CreateTestRequest(http.MethodGet, "/"+alias, alias)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusPermanentRedirect {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusPermanentRedirect)
	}

	resp, err := ParseRedirectResponse(rr.Result())
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.Status != "OK" {
		t.Errorf("expected status OK, got %s", resp.Status)
	}

	if resp.Url != url {
		t.Errorf("expected url %s, got %s", url, resp.Url)
	}
}

func TestRedirectNotFound(t *testing.T) {
	mockStorage := storagemock.NewMockURLSaver()
	handler := SetupHandler(mockStorage)

	alias := "non-existent-alias"
	req, _ := CreateTestRequest(http.MethodGet, "/"+alias, alias)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	resp, err := ParseRedirectResponse(rr.Result())
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.Status != "ERROR" {
		t.Errorf("expected status ERROR, got %s", resp.Status)
	}

	if resp.Error != storage.ErrUrlNotFound.Error() {
		t.Errorf("expected error %s, got %s", storage.ErrUrlNotFound.Error(), resp.Error)
	}
}

func TestRedirectEmptyAlias(t *testing.T) {
	mockStorage := storagemock.NewMockURLSaver()
	handler := SetupHandler(mockStorage)

	alias := ""
	req, _ := CreateTestRequest(http.MethodGet, "/", alias)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	resp, err := ParseRedirectResponse(rr.Result())
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.Status != "ERROR" {
		t.Errorf("expected status ERROR, got %s", resp.Status)
	}
}

func TestRedirectMultipleAliases(t *testing.T) {
	mockStorage := storagemock.NewMockURLSaver()
	handler := SetupHandler(mockStorage)

	testCases := []struct {
		alias string
		url   string
	}{
		{"google", "https://google.com"},
		{"github", "https://github.com"},
		{"example", "https://example.com"},
	}

	for _, tc := range testCases {
		mockStorage.SavedURLs[tc.alias] = tc.url
	}

	for _, tc := range testCases {
		t.Run(tc.alias, func(t *testing.T) {
			req, _ := CreateTestRequest(http.MethodGet, "/"+tc.alias, tc.alias)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusPermanentRedirect {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusPermanentRedirect)
			}

			resp, err := ParseRedirectResponse(rr.Result())
			if err != nil {
				t.Fatalf("failed to parse response: %v", err)
			}

			if resp.Url != tc.url {
				t.Errorf("expected url %s, got %s", tc.url, resp.Url)
			}
		})
	}
}
