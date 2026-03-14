package delete

import (
	"net/http"
	"net/http/httptest"
	"testing"
	del "urlshort/internal/http-server/handlers/url/delete"
	storagemock "urlshort/tests/internal/storage"
)

func TestDeleteURLSuccess(t *testing.T) {
	mockStorage := storagemock.NewMockURLSaver()
	handler := SetupHandler(mockStorage)
	alias := "alias"
	url := "https://example.com"
	mockStorage.SavedURLs[alias] = url

	reqBody := del.Request{
		Alias: alias,
	}

	req, _ := CreateTestRequest(http.MethodDelete, "/url", reqBody)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	resp, err := ParseDeleteURLResponse(rr.Result())
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.Status != "OK" {
		t.Errorf("expected status OK, got %s", resp.Status)
	}

	if _, exists := mockStorage.SavedURLs[alias]; exists {
		t.Errorf("url was not deleted")
	}
}
