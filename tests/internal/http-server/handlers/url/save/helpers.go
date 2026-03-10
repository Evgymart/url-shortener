package save

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"urlshort/internal/http-server/handlers/url/save"
)

func CreateTestRequest(method, path string, body interface{}) (*http.Request, error) {
	var bodyReader io.Reader

	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req := httptest.NewRequest(method, path, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func ParseSaveURLResponse(resp *http.Response) (save.Response, error) {
	var body save.Response
	err := json.NewDecoder(resp.Body).Decode(&body)
	return body, err
}

func SetupHandler(mockStorage *MockURLSaver) http.HandlerFunc {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	return save.New(logger, mockStorage)
}
