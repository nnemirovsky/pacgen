package service

import (
	"context"
	"github.com/nnemirovsky/pacgen/internal/model"
	"github.com/nnemirovsky/pacgen/pkg/gen"
	"github.com/rs/zerolog"
	"io"
	"os"
)

type PACService struct {
	logger   zerolog.Logger
	repo     RuleRepository
	filePath string
}

func NewPACService(repo RuleRepository, filePath string, logger zerolog.Logger) *PACService {
	return &PACService{
		logger:   logger,
		repo:     repo,
		filePath: filePath,
	}
}

func (s *PACService) GeneratePACFile(ctx context.Context) error {
	rules, err := s.repo.GetAllWithProfiles(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("Error occurred while getting rules to generate pac file")
		return err
	}

	if err = generatePACFile(rules, s.filePath); err != nil {
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

func generatePACFile(rules []model.Rule, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_SYNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	return generatePAC(file, rules)
}
