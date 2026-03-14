package delete

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	resp "urlshort/internal/http-server/handlers/api/response"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Alias string `json:"alias"`
}

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

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", err.Error())
			render.JSON(w, r, resp.Error(err.Error()))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))
		if err = validator.New().Struct(req); err != nil {
			log.Error("failed to validate request", err.Error())

			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)
			render.JSON(w, r, resp.ValidateErrors(validateErr))
			return
		}

		alias := req.Alias
		err = urlDeleter.DeleteURL(r.Context(), alias)
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
