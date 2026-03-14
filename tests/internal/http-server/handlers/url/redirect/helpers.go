package redirect

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"

	"urlshort/internal/http-server/handlers/url/redirect"
	"urlshort/tests/internal/storage"

	"github.com/go-chi/chi/v5"
)

func SetupHandler(mockStorage *storage.MockURLSaver) http.HandlerFunc {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	return redirect.New(logger, mockStorage)
}

func CreateTestRequest(method, path, alias string) (*http.Request, error) {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("alias", alias)

	req := httptest.NewRequest(method, path, nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	return req, nil
}
