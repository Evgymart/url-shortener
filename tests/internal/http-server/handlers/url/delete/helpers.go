package delete

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	resp "urlshort/internal/http-server/handlers/api/response"
	del "urlshort/internal/http-server/handlers/url/delete"
	"urlshort/tests/internal/storage"
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

func ParseDeleteURLResponse(response *http.Response) (resp.Response, error) {
	var body resp.Response
	err := json.NewDecoder(response.Body).Decode(&body)
	return body, err
}

func SetupHandler(mockStorage *storage.MockURLSaver) http.HandlerFunc {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	return del.New(logger, mockStorage)
}
