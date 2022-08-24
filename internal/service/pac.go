package service

import (
	"context"
	"github.com/rs/zerolog"
	"io"
	"os"
	"pacgen/internal/model"
	"pacgen/pkg/gen"
)

type PACService struct {
	logger zerolog.Logger
	repo   RuleRepository
}

func NewPACService(repo RuleRepository, logger zerolog.Logger) *PACService {
	return &PACService{
		logger: logger,
		repo:   repo,
	}
}

func (s *PACService) GeneratePACFile(ctx context.Context) error {
	rules, err := s.repo.GetAllWithProfiles(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("Error occurred while getting rules to generate pac file")
		return err
	}

	if err = generatePACFile(rules); err != nil {
		s.logger.Error().Err(err).Msg("Error occurred while generating pac file")
		return err
	}

	return nil
}

func generatePAC(wr io.Writer, rules []model.Rule) error {
	conditions := make([]gen.Condition, 0)
	for _, rule := range rules {
		action := "DIRECT"
		if rule.ProxyProfile != nil {
			action = rule.ProxyProfile.Type.String() + " " + rule.ProxyProfile.Address
		}
		conditions = append(conditions, gen.Condition{Regex: rule.Regex, Action: action})
	}

	return gen.Generate(wr, conditions)
}

func generatePACFile(rules []model.Rule) error {
	file, err := os.OpenFile("data/proxy.pac", os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_SYNC, 0644)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	return generatePAC(file, rules)
}
