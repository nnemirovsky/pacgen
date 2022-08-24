package main

import (
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"pacgen/internal/repository"
	"pacgen/internal/service"
	"pacgen/pkg/logutil"
	"time"
)

func main() {
	logger := logutil.Logger

	db := sqlx.MustConnect("sqlite3", "./data/data.db")
	defer func() {
		if err := db.Close(); err != nil {
			logger.Fatal().Err(err).Send()
		}
	}()

	ruleRepo := repository.NewRuleRepository(db, logger)
	pacSrvc := service.NewPACService(ruleRepo, logger)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := pacSrvc.GeneratePACFile(ctx); err != nil {
		logger.Fatal().Err(err).Send()
	}

	logger.Info().Msg("PAC file generated successfully")
}
