package main

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"os"
	"pacgen/pkg/logutil"
	"strings"
)

var logger = logutil.Logger

func main() {
	action := getAction()

	m, err := migrate.New("file://migrations", "sqlite3://data/data.db")
	if err != nil {
		logger.Fatal().Err(err).Msg("Error occurred while trying to create new migrate instance")
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			logger.Fatal().Err(srcErr).Msg("Error occurred while closing source")
		}
		if dbErr != nil {
			logger.Fatal().Err(dbErr).Msg("Error occurred while closing database")
		}
	}()

	action(m)

	logger.Info().Msg("Migration completed successfully")
}

func getAction() func(m *migrate.Migrate) {
	if len(os.Args) < 2 {
		logger.Fatal().Msg("Specify action: 'up' or 'down'")
	}

	var action func(m *migrate.Migrate)
	switch a := strings.ToLower(strings.TrimSpace(os.Args[1])); a {
	case "up":
		action = up
	case "down":
		action = down
	default:
		logger.Fatal().Msgf("Invalid action '%s', should be either 'up' or 'down'", a)
	}
	return action
}

func up(m *migrate.Migrate) {
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Fatal().Err(err).Msg("Error occurred while migrating all the way up")
	}
}

func down(m *migrate.Migrate) {
	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		logger.Fatal().Err(err).Msg("Error occurred while migrating all the way down")
	}
}
