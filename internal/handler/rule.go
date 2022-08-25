package handler

import (
	"context"
	"fmt"
	"github.com/go-chi/render"
	"github.com/nnemirovsky/pacgen/internal/errs"
	"github.com/nnemirovsky/pacgen/internal/model"
	"github.com/nnemirovsky/pacgen/pkg/rest"
	"github.com/rs/zerolog"
	"net/http"
)

type RuleService interface {
	GetAll(ctx context.Context) ([]model.Rule, error)
	GetAllWithProfiles(ctx context.Context) ([]model.Rule, error)
	GetByID(ctx context.Context, id int) (model.Rule, error)
	Create(ctx context.Context, profile *model.Rule) error
	Update(ctx context.Context, profile model.Rule) error
	Delete(ctx context.Context, id int) error
}

type RuleHandler struct {
	logger  zerolog.Logger
	service RuleService
}

func NewRuleHandler(service RuleService, logger zerolog.Logger) *RuleHandler {
	return &RuleHandler{
		logger:  logger,
		service: service,
	}
}

func (h *RuleHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	rules, err := h.service.GetAll(r.Context())
	if err != nil {
		h.logger.Error().Err(err).Msg("Error occurred while getting all rules")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ruleEntities := make([]RuleR, 0)
	for _, rule := range rules {
		ruleR := RuleR{}
		ruleR.FromModel(rule)
		ruleEntities = append(ruleEntities, ruleR)
	}

	render.JSON(w, r, ruleEntities)
	w.WriteHeader(http.StatusOK)
}

func (h *RuleHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, ok := getIDFromURL(w, r, h.logger)
	if !ok {
		return
	}

	rule, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if _, ok := err.(*errs.EntityNotFoundError); ok {
			h.logger.Debug().Err(err).Send()
			Render(w, r, rest.NotFoundResponse(err.Error()), h.logger)
			return
		}
		h.logger.Error().Err(err).Msg("Error occurred while getting rule by id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ruleR := RuleR{}
	ruleR.FromModel(rule)

	render.JSON(w, r, ruleR)
	w.WriteHeader(http.StatusOK)
}

func (h *RuleHandler) Create(w http.ResponseWriter, r *http.Request) {
	rule := RuleCU{}
	if ok := getFromBodyAndValidate(w, r, h.logger, &rule); !ok {
		return
	}

	ruleModel, err := rule.ToModel()
	if err != nil {
		h.logger.Error().Err(err).Msg("Error occurred while converting rule entity to corresponding model")
	}

	err = h.service.Create(r.Context(), &ruleModel)
	if err == errs.InvalidReferenceError {
		h.logger.Debug().Err(err).Send()
		Render(w, r, rest.ConflictResponse(err.Error()), h.logger)
		return
	}
	if err != nil {
		h.logger.Error().Err(err).Msg("Error occurred while creating rule")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rest.Created(w, fmt.Sprintf("%s/%v", r.URL, ruleModel.ID))
}

func (h *RuleHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := getIDFromURL(w, r, h.logger)
	if !ok {
		return
	}

	rule := RuleCU{}
	if ok := getFromBodyAndValidate(w, r, h.logger, &rule); !ok {
		return
	}

	ruleModel, err := rule.ToModel()
	if err != nil {
		h.logger.Error().Err(err).Msg("Error occurred while converting rule entity to corresponding model")
	}
	ruleModel.ID = id

	err = h.service.Update(r.Context(), ruleModel)
	if err == errs.InvalidReferenceError {
		h.logger.Debug().Err(err).Send()
		Render(w, r, rest.ConflictResponse(err.Error()), h.logger)
		return
	}
	if err != nil {
		if _, ok := err.(*errs.EntityNotFoundError); ok {
			h.logger.Debug().Err(err).Send()
			Render(w, r, rest.NotFoundResponse(err.Error()), h.logger)
			return
		}
		h.logger.Error().Err(err).Msg("Error occurred while updating rule")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.NoContent(w, r)
}

func (h *RuleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := getIDFromURL(w, r, h.logger)
	if !ok {
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if _, ok := err.(*errs.EntityNotFoundError); ok {
			h.logger.Debug().Err(err).Send()
			Render(w, r, rest.NotFoundResponse(err.Error()), h.logger)
			return
		}
		h.logger.Error().Err(err).Msg("Error occurred while deleting rule")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.NoContent(w, r)
}
