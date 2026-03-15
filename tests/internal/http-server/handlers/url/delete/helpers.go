package delete

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	resp "urlshort/internal/http-server/handlers/api/response"
	del "urlshort/internal/http-server/handlers/url/delete"
	"urlshort/tests/internal/storage"

	"github.com/go-chi/chi/v5"
)

func CreateTestRequest(method, path, alias string) (*http.Request, error) {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("alias", alias)

	req := httptest.NewRequest(method, path, nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
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
