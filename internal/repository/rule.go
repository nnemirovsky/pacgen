package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
	"github.com/nnemirovsky/pacgen/internal/errs"
	"github.com/nnemirovsky/pacgen/internal/model"
	"github.com/rs/zerolog"
)

type RuleRepository struct {
	logger zerolog.Logger
	db     *sqlx.DB
}

func NewRuleRepository(db *sqlx.DB, logger zerolog.Logger) *RuleRepository {
	return &RuleRepository{
		logger: logger,
		db:     db,
	}
}

func (r *RuleRepository) GetAll(ctx context.Context) ([]model.Rule, error) {
	query := `SELECT id, regex, proxy_profile_id AS "proxy_profile.id" FROM rules`
	rules := make([]model.Rule, 0)
	if err := r.db.SelectContext(ctx, &rules, query); err != nil {
		r.logger.Error().Err(err).Msg("Error occurred while getting rules")
		return nil, errs.RepositoryUnknownError
	}
	return rules, nil
}

func (r *RuleRepository) GetAllWithProfiles(ctx context.Context) ([]model.Rule, error) {
	query := `SELECT r.id,
    				 r.regex,
    				 p.id AS "proxy_profile.id",
    				 p.name AS "proxy_profile.name",
    				 p.type AS "proxy_profile.type",
    				 p.address AS "proxy_profile.address"
			  FROM rules r
			  JOIN proxy_profiles p ON r.proxy_profile_id = p.id`

	rules := make([]model.Rule, 0)
	if err := r.db.SelectContext(ctx, &rules, query); err != nil {
		r.logger.Error().Err(err).Msg("Error occurred while getting rules")
		return nil, errs.RepositoryUnknownError
	}
	return rules, nil
}

func (r *RuleRepository) GetByID(ctx context.Context, id int) (model.Rule, error) {
	query := `SELECT r.id,
					 r.regex,
					 p.id AS "proxy_profile.id",
					 p.name AS "proxy_profile.name",
					 p.type AS "proxy_profile.type",
					 p.address AS "proxy_profile.address"
			  FROM rules r
			  JOIN proxy_profiles p ON r.proxy_profile_id = p.id
			  WHERE r.id = ?`

	var rule model.Rule
	if err := r.db.GetContext(ctx, &rule, query, id); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			err = &errs.EntityNotFoundError{Name: "rule", Key: "id", Value: id}
			r.logger.Debug().Err(err).Send()
		default:
			err = errs.RepositoryUnknownError
			r.logger.Error().Err(err).Msg("Error occurred while getting rule by id")
		}
		return model.Rule{}, err
	}
	return rule, nil
}

func (r *RuleRepository) Create(ctx context.Context, rule *model.Rule) error {
	cmd := `INSERT INTO rules (regex, proxy_profile_id) VALUES (:regex, :proxy_profile.id) RETURNING id`

	result, err := r.db.NamedExecContext(ctx, cmd, rule)
	if e, ok := err.(sqlite3.Error); ok && e.Code == sqlite3.ErrConstraint {
		err = errs.InvalidReferenceError
		r.logger.Debug().Err(err).Msg("Unknown proxy profile")
		return err
	}
	if err != nil {
		r.logger.Error().Err(err).Msg("Error occurred while creating rule")
		return errs.RepositoryUnknownError
	}

	id, err := result.LastInsertId()
	if err != nil {
		r.logger.Error().Err(err).Msg("Error occurred while retrieving created rule id")
		return errs.RepositoryUnknownError
	}
	rule.ID = int(id)
	return nil
}

func (r *RuleRepository) Update(ctx context.Context, rule model.Rule) error {
	cmd := `UPDATE rules SET regex = :regex, proxy_profile_id = :proxy_profile.id WHERE id = :id`

	result, err := r.db.NamedExecContext(ctx, cmd, rule)
	if e, ok := err.(sqlite3.Error); ok && e.Code == sqlite3.ErrConstraint {
		err = errs.InvalidReferenceError
		r.logger.Debug().Err(err).Msg("Unknown proxy profile")
		return err
	}
	if err != nil {
		r.logger.Error().Err(err).Msg("Error occurred while updating rule")
		return errs.RepositoryUnknownError
	}
	count, err := result.RowsAffected()
	if err != nil {
		r.logger.Error().Err(err).Msg("Error occurred while retrieving affected rows count after updating rule")
		return errs.RepositoryUnknownError
	}
	if count == 0 {
		err = &errs.EntityNotFoundError{Name: "rule", Key: "id", Value: rule.ID}
		r.logger.Debug().Err(err).Send()
		return err
	}
	return nil
}

func (r *RuleRepository) Delete(ctx context.Context, id int) error {
	cmd := `DELETE FROM rules WHERE id = ?`
	result, err := r.db.ExecContext(ctx, cmd, id)
	if err != nil {
		r.logger.Error().Err(err).Msg("Error occurred while deleting rule")
		return errs.RepositoryUnknownError
	}
	count, err := result.RowsAffected()
	if err != nil {
		r.logger.Error().Err(err).Msg("Error occurred while retrieving affected rows count after deleting rule")
		return errs.RepositoryUnknownError
	}
	if count == 0 {
		err = &errs.EntityNotFoundError{Name: "rule", Key: "id", Value: id}
		r.logger.Debug().Err(err).Send()
		return err
	}
	return nil
}
