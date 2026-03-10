package save

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"urlshort/internal/storage"
	"urlshort/lib/random"

	resp "urlshort/internal/http-server/handlers/api/response"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

const aliasLength = 6

type URLSaver interface {
	SaveURL(ctx context.Context, urlAlias string, longUrl string) error
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const action = "handels.url.save"
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
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		err = urlSaver.SaveURL(r.Context(), req.URL, alias)
		if errors.Is(err, storage.ErrUrlAlreadyExists) {
			log.Info("url already exists", slog.String("url", req.URL))
			render.Status(r, http.StatusConflict)
			render.JSON(w, r, resp.Error("url already exists"))
			return
		}

		if err != nil {
			log.Error("failed to save url", err.Error())
			render.JSON(w, r, resp.Error("failed to save url"))
			return
		}

		log.Info("url saved", slog.String("url", req.URL))
		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias:    alias,
		})
	}
}
