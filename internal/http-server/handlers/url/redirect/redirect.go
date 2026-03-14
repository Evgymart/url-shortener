package redirect

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	resp "urlshort/internal/http-server/handlers/api/response"
	"urlshort/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type URLGetter interface {
	GetUrl(ctx context.Context, urlAlias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const action = "handels.url.redirect"
		log = log.With(
			slog.String("action", action),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParamFromCtx(r.Context(), "alias")
		url, err := urlGetter.GetUrl(r.Context(), alias)
		if errors.Is(err, storage.ErrUrlNotFound) {
			log.Error("Url not found", err.Error())
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, resp.Error(err.Error()))
			return
		}

		if err != nil {
			log.Error("Error getting url", err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("Error getting url"))
			return
		}

		log.Info("url redirected", slog.String("url", url))
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}
