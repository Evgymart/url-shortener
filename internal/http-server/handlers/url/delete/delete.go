package delete

import (
	"context"
	"log/slog"
	"net/http"
	resp "urlshort/internal/http-server/handlers/api/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type URLDeleter interface {
	DeleteURL(ctx context.Context, urlAlias string) error
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const action = "handels.url.delete"
		log := log.With(
			slog.String("action", action),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParamFromCtx(r.Context(), "alias")
		err := urlDeleter.DeleteURL(r.Context(), alias)
		if err != nil {
			log.Error("failed to delete url", err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to delete url"))
			return
		}

		log.Info("url deleted", slog.String("alias", alias))
		render.JSON(w, r, resp.OK())
	}
}
