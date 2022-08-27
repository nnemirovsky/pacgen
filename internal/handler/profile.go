package handler

import (
	"fmt"
	"github.com/go-chi/render"
	"github.com/nnemirovsky/pacgen/internal/errs"
	"github.com/nnemirovsky/pacgen/pkg/rest"
	"github.com/rs/zerolog"
	"net/http"
)

type ProxyProfileHandler struct {
	logger  zerolog.Logger
	service ProxyProfileService
}

func NewProxyProfileHandler(service ProxyProfileService, logger zerolog.Logger) *ProxyProfileHandler {
	return &ProxyProfileHandler{
		logger:  logger,
		service: service,
	}
}

func (h *ProxyProfileHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	profiles, err := h.service.GetAll(r.Context())
	if err != nil {
		h.logger.Error().Err(err).Msg("Error occurred while getting all proxy profiles")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	profileEntities := make([]ProxyProfileR, 0)
	for _, profile := range profiles {
		profileR := ProxyProfileR{}
		profileR.FromModel(profile)
		profileEntities = append(profileEntities, profileR)
	}

	render.JSON(w, r, profileEntities)
	w.WriteHeader(http.StatusOK)
}

func (h *ProxyProfileHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, ok := getIDFromURL(w, r, h.logger)
	if !ok {
		return
	}

	profile, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if _, ok := err.(*errs.EntityNotFoundError); ok {
			h.logger.Debug().Err(err).Send()
			Render(w, r, rest.NotFoundResponse(err.Error()), h.logger)
			return
		}
		h.logger.Error().Err(err).Msg("Error occurred while getting proxy profile by id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	profileR := ProxyProfileR{}
	profileR.FromModel(profile)

	render.JSON(w, r, profileR)
	w.WriteHeader(http.StatusOK)
}

func (h *ProxyProfileHandler) Create(w http.ResponseWriter, r *http.Request) {
	profileCU := ProxyProfileCU{}
	if ok := getFromBodyAndValidate(w, r, h.logger, &profileCU); !ok {
		return
	}

	profileModel, err := profileCU.ToModel()
	if err != nil {
		h.logger.Debug().Err(err).Msg("Error occurred while converting profile entity to corresponding model")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if err := h.service.Create(r.Context(), &profileModel); err != nil {
		if _, ok := err.(*errs.EntityAlreadyExistsError); ok {
			h.logger.Debug().Err(err).Send()
			Render(w, r, rest.ConflictResponse(err.Error()), h.logger)
			return
		}
		h.logger.Error().Err(err).Msg("Error occurred while creating profile")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rest.Created(w, fmt.Sprintf("%s/%v", r.URL, profileModel.ID))
}

func (h *ProxyProfileHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := getIDFromURL(w, r, h.logger)
	if !ok {
		return
	}

	profile := ProxyProfileCU{}
	if ok := getFromBodyAndValidate(w, r, h.logger, &profile); !ok {
		return
	}

	profileModel, err := profile.ToModel()
	if err != nil {
		h.logger.Debug().Err(err).Msg("Error occurred while converting profile entity to corresponding model")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	profileModel.ID = id

	if err := h.service.Update(r.Context(), profileModel); err != nil {
		switch err.(type) {
		case *errs.EntityAlreadyExistsError:
			h.logger.Debug().Err(err).Send()
			Render(w, r, rest.ConflictResponse(err.Error()), h.logger)
			return
		case *errs.EntityNotFoundError:
			h.logger.Debug().Err(err).Send()
			Render(w, r, rest.NotFoundResponse(err.Error()), h.logger)
			return
		default:
			h.logger.Error().Err(err).Msg("Error occurred while updating profile")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	render.NoContent(w, r)
}

func (h *ProxyProfileHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := getIDFromURL(w, r, h.logger)
	if !ok {
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		switch err.(type) {
		case *errs.EntityNotFoundError:
			h.logger.Debug().Err(err).Send()
			Render(w, r, rest.NotFoundResponse(err.Error()), h.logger)
			return
		case *errs.EntityStillReferencedError:
			h.logger.Debug().Err(err).Send()
			Render(w, r, rest.ConflictResponse(err.Error()), h.logger)
			return
		default:
			h.logger.Error().Err(err).Msg("Error occurred while deleting profile")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	render.NoContent(w, r)
}
