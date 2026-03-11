package redirect

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"urlshort/internal/http-server/handlers/url/redirect"
	"urlshort/tests/internal/storage"

	"github.com/go-chi/chi/v5"
)

func CreateTestRequest(method, path, alias string) (*http.Request, error) {
	req := httptest.NewRequest(method, path, nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("alias", alias)

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	return req, nil
}

func ParseRedirectResponse(resp *http.Response) (redirect.Response, error) {
	var body redirect.Response
	err := json.NewDecoder(resp.Body).Decode(&body)
	return body, err
}

func SetupHandler(mockStorage *storage.MockURLSaver) http.HandlerFunc {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	return redirect.New(logger, mockStorage)
}
