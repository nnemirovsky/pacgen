package service

import (
	"context"
	"github.com/nnemirovsky/pacgen/internal/errs"
	"github.com/nnemirovsky/pacgen/internal/model"
	"github.com/rs/zerolog"
	"time"
)

type ProxyProfileRepository interface {
	GetAll(ctx context.Context) ([]model.ProxyProfile, error)
	GetByID(ctx context.Context, id int) (model.ProxyProfile, error)
	Create(ctx context.Context, profile *model.ProxyProfile) error
	Update(ctx context.Context, profile model.ProxyProfile) error
	Delete(ctx context.Context, id int) error
}

type ProxyProfileService struct {
	logger  zerolog.Logger
	repo    ProxyProfileRepository
	pacSrvc *PACService
}

func NewProxyProfileService(
	repo ProxyProfileRepository,
	pacSrvc *PACService,
	logger zerolog.Logger,
) *ProxyProfileService {
	return &ProxyProfileService{
		logger:  logger,
		repo:    repo,
		pacSrvc: pacSrvc,
	}
}

func (s *ProxyProfileService) GetAll(ctx context.Context) ([]model.ProxyProfile, error) {
	profiles, err := s.repo.GetAll(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("Error occurred while getting proxy profiles")
		return nil, errs.ServiceUnknownError
	}
	return profiles, nil
}

func (s *ProxyProfileService) GetByID(ctx context.Context, id int) (model.ProxyProfile, error) {
	profile, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if _, ok := err.(*errs.EntityNotFoundError); ok {
			s.logger.Debug().Err(err).Send()
			return profile, err
		}
		s.logger.Error().Err(err).Msg("Error occurred while getting proxy profile by id")
		return profile, errs.ServiceUnknownError
	}
	return profile, nil
}

func (s *ProxyProfileService) Create(ctx context.Context, profile *model.ProxyProfile) error {
	err := s.repo.Create(ctx, profile)
	if err != nil {
		if _, ok := err.(*errs.EntityAlreadyExistsError); ok {
			s.logger.Debug().Err(err).Send()
			return err
		}
		s.logger.Error().Err(err).Msg("Error occurred while creating proxy profile")
		return errs.ServiceUnknownError
	}

	s.logger.Debug().Int("profile-id", profile.ID).Msg("Proxy profile created")

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.pacSrvc.GeneratePACFile(ctx); err != nil {
			s.logger.Error().Err(err).Msg("Error occurred while generating pac file after creating profile")
			return
		}
		s.logger.Debug().Msg("Pac file generated after creating profile")
	}()

	return nil
}

func (s *ProxyProfileService) Update(ctx context.Context, profile model.ProxyProfile) error {
	err := s.repo.Update(ctx, profile)
	if err != nil {
		switch err.(type) {
		case *errs.EntityNotFoundError:
			s.logger.Debug().Err(err).Send()
			return err
		case *errs.EntityAlreadyExistsError:
			s.logger.Debug().Err(err).Send()
			return err
		default:
			s.logger.Error().Err(err).Msg("Error occurred while updating proxy profile")
			return errs.ServiceUnknownError
		}
	}

	s.logger.Debug().Int("profile-id", profile.ID).Msg("Proxy profile updated")

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.pacSrvc.GeneratePACFile(ctx); err != nil {
			s.logger.Error().Err(err).Msg("Error occurred while generating pac file after updating profile")
			return
		}
		s.logger.Debug().Msg("Pac file generated after updating profile")
	}()

	return nil
}

func (s *ProxyProfileService) Delete(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		switch err.(type) {
		case *errs.EntityNotFoundError:
			s.logger.Debug().Err(err).Send()
			return err
		case *errs.EntityStillReferencedError:
			s.logger.Debug().Err(err).Send()
			return err
		default:
			s.logger.Error().Err(err).Msg("Error occurred while deleting proxy profile")
			return errs.ServiceUnknownError
		}
	}

	s.logger.Debug().Int("profile-id", id).Msg("Proxy profile deleted")

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.pacSrvc.GeneratePACFile(ctx); err != nil {
			s.logger.Error().Err(err).Msg("Error occurred while generating pac file after deleting profile")
			return
		}
		s.logger.Debug().Msg("Pac file generated after deleting profile")
	}()

	return nil
}
