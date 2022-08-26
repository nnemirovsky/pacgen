package service

import (
	"context"
	"github.com/nnemirovsky/pacgen/internal/errs"
	"github.com/nnemirovsky/pacgen/internal/model"
	"github.com/rs/zerolog"
	"time"
)

type RuleService struct {
	logger  zerolog.Logger
	repo    RuleRepository
	pacSrvc pacService
}

func NewRuleService(repo RuleRepository, pacSrvc pacService, logger zerolog.Logger) *RuleService {
	return &RuleService{
		logger:  logger,
		repo:    repo,
		pacSrvc: pacSrvc,
	}
}

func (s *RuleService) GetAll(ctx context.Context) ([]model.Rule, error) {
	rules, err := s.repo.GetAll(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("Error occurred while getting rules")
		return nil, errs.ServiceUnknownError
	}
	return rules, nil
}

func (s *RuleService) GetAllWithProfiles(ctx context.Context) ([]model.Rule, error) {
	rules, err := s.repo.GetAllWithProfiles(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("Error occurred while getting rules")
		return nil, errs.ServiceUnknownError
	}
	return rules, nil
}

func (s *RuleService) GetByID(ctx context.Context, id int) (model.Rule, error) {
	rule, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if _, ok := err.(*errs.EntityNotFoundError); ok {
			s.logger.Debug().Err(err).Send()
			return rule, err
		}
		s.logger.Error().Err(err).Msg("Error occurred while getting rule by id")
		return rule, errs.ServiceUnknownError
	}
	return rule, nil
}

func (s *RuleService) Create(ctx context.Context, rule *model.Rule) error {
	err := s.repo.Create(ctx, rule)
	if err == errs.InvalidReferenceError {
		s.logger.Debug().Err(err).Send()
		return err
	}
	if err != nil {
		s.logger.Error().Err(err).Msg("Error occurred while creating rule")
		return errs.ServiceUnknownError
	}

	s.logger.Debug().Err(err).Int("rule-id", rule.ID).Msg("Rule created")

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.pacSrvc.GeneratePACFile(ctx); err != nil {
			s.logger.Error().Err(err).Msg("Error occurred while generating pac file after creating rule")
			return
		}
		s.logger.Debug().Msg("Pac file generated after creating rule")
	}()

	return nil
}

func (s *RuleService) Update(ctx context.Context, rule model.Rule) error {
	err := s.repo.Update(ctx, rule)
	if err == errs.InvalidReferenceError {
		s.logger.Debug().Err(err).Send()
		return err
	}
	if _, ok := err.(*errs.EntityNotFoundError); ok {
		s.logger.Debug().Err(err).Send()
		return err
	}
	if err != nil {
		s.logger.Error().Err(err).Msg("Error occurred while updating rule")
		return errs.ServiceUnknownError
	}

	s.logger.Debug().Int("rule-id", rule.ID).Msg("Rule updated")

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.pacSrvc.GeneratePACFile(ctx); err != nil {
			s.logger.Error().Err(err).Msg("Error occurred while generating pac file after updating rule")
			return
		}
		s.logger.Debug().Msg("Pac file generated after updating rule")
	}()

	return nil
}

func (s *RuleService) Delete(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if _, ok := err.(*errs.EntityNotFoundError); ok {
			s.logger.Debug().Err(err).Send()
			return err
		}
		s.logger.Error().Err(err).Msg("Error occurred while deleting rule")
		return errs.ServiceUnknownError
	}

	s.logger.Debug().Int("rule-id", id).Msg("Rule deleted")

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.pacSrvc.GeneratePACFile(ctx); err != nil {
			s.logger.Error().Err(err).Msg("Error occurred while generating pac file after deleting rule")
			return
		}
		s.logger.Debug().Msg("Pac file generated after deleting rule")
	}()

	return nil
}
