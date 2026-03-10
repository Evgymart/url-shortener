package save

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"urlshort/internal/http-server/handlers/url/save"
	"urlshort/internal/storage"
)

func TestSaveURLSuccess(t *testing.T) {
	mockStorage := NewMockURLSaver()
	handler := SetupHandler(mockStorage)

	alias := "alias"
	reqBody := save.Request{
		URL:   "https://example.com",
		Alias: alias,
	}

	req, _ := CreateTestRequest(http.MethodPost, "/url", reqBody)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	resp, err := ParseSaveURLResponse(rr.Result())
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.Status != "OK" {
		t.Errorf("expected status OK, got %s", resp.Status)
	}

	if resp.Alias != alias {
		t.Errorf("expected alias %s, got %s", alias, resp.Alias)
	}
}

func TestAlreadyExistsError(t *testing.T) {
	mockStorage := NewMockURLSaver()
	handler := SetupHandler(mockStorage)
	url := "https://example.com"
	alias := "alias"
	mockStorage.SavedURLs[alias] = url

	reqBody := save.Request{
		URL:   url,
		Alias: alias,
	}

	req, _ := CreateTestRequest(http.MethodPost, "/url", reqBody)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusConflict)
	}

	resp, err := ParseSaveURLResponse(rr.Result())
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.Status != "ERROR" {
		t.Errorf("expected status Error, got %s", resp.Status)
	}

	if resp.Error != storage.ErrUrlAlreadyExists.Error() {
		t.Errorf("Expected error %s, got %s", storage.ErrUrlAlreadyExists.Error(), resp.Error)
	}
}
