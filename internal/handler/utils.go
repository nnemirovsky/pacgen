package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/nnemirovsky/pacgen/pkg/rest"
	"github.com/rs/zerolog"
	"net/http"
	"strconv"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func Render(w http.ResponseWriter, r *http.Request, v render.Renderer, logger zerolog.Logger) {
	if err := render.Render(w, r, v); err != nil {
		logger.Error().Err(err).Msg("Error occurred while trying to render response")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func getIDFromURL(w http.ResponseWriter, r *http.Request, logger zerolog.Logger) (id int, ok bool) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		logger.Debug().Err(err).Msg("Error occurred while converting id to int")
		Render(w, r, rest.BadRequestResponse("Invalid id format"), logger)
		return 0, false
	}
	return id, true
}

func getFromBodyAndValidate(w http.ResponseWriter, r *http.Request, logger zerolog.Logger, entity any) (ok bool) {
	if err := render.DecodeJSON(r.Body, entity); err != nil {
		logger.Debug().Err(err).Msg("Error occurred while decoding request body")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return false
	}

	if err := validate.Struct(entity); err != nil {
		logger.Debug().Err(err).Msg("Error occurred while validating request body")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return false
	}

	return true
}
