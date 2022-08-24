package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	"pacgen/internal/errs"
	"pacgen/internal/model"
)

type ProxyProfileRepository struct {
	logger zerolog.Logger
	db     *sqlx.DB
}

func NewProxyProfileRepository(db *sqlx.DB, logger zerolog.Logger) *ProxyProfileRepository {
	return &ProxyProfileRepository{
		logger: logger,
		db:     db,
	}
}

func (r *ProxyProfileRepository) GetAll(ctx context.Context) ([]model.ProxyProfile, error) {
	query := `SELECT id, name, type, address FROM proxy_profiles`
	profiles := make([]model.ProxyProfile, 0)
	if err := r.db.SelectContext(ctx, &profiles, query); err != nil {
		r.logger.Error().Err(err).Msg("Error occurred while getting proxy profiles")
		return nil, errs.RepositoryUnknownError
	}
	return profiles, nil
}

func (r *ProxyProfileRepository) GetByID(ctx context.Context, id int) (model.ProxyProfile, error) {
	query := `SELECT id, name, type, address FROM proxy_profiles WHERE id = ?`
	var profile model.ProxyProfile
	if err := r.db.GetContext(ctx, &profile, query, id); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			err = &errs.EntityNotFoundError{Name: "proxy profile", Key: "id", Value: id}
			r.logger.Debug().Err(err).Send()
		default:
			err = errs.RepositoryUnknownError
			r.logger.Error().Err(err).Msg("Error occurred while getting proxy profile by id")
		}
		return model.ProxyProfile{}, err
	}
	return profile, nil
}

func (r *ProxyProfileRepository) Create(ctx context.Context, profile *model.ProxyProfile) error {
	cmd := `INSERT INTO proxy_profiles (name, type, address) VALUES (:name, :type, :address)`
	result, err := r.db.NamedExecContext(ctx, cmd, profile)
	if err != nil {
		if s, ok := err.(sqlite3.Error); ok && s.Code == sqlite3.ErrConstraint {
			err := &errs.EntityAlreadyExistsError{Name: "proxy profile", Key: "name", Value: profile.Name}
			r.logger.Debug().Err(err).Send()
			return err
		}
		r.logger.Error().Err(err).Msg("Error occurred while creating proxy profile")
		return errs.RepositoryUnknownError
	}
	id, err := result.LastInsertId()
	if err != nil {
		r.logger.Error().Err(err).Msg("Error occurred while retrieving created profile id")
		return errs.RepositoryUnknownError
	}
	profile.ID = int(id)
	return nil
}

func (r *ProxyProfileRepository) Update(ctx context.Context, profile model.ProxyProfile) error {
	cmd := `UPDATE proxy_profiles SET name = :name, type = :type, address = :address WHERE id = :id`
	result, err := r.db.NamedExecContext(ctx, cmd, profile)
	if err != nil {
		if s, ok := err.(sqlite3.Error); ok && s.Code == sqlite3.ErrConstraint {
			err := &errs.EntityAlreadyExistsError{Name: "proxy profile", Key: "name", Value: profile.Name}
			r.logger.Debug().Err(err).Send()
			return err
		}
		r.logger.Error().Err(err).Msg("Error occurred while updating proxy profile")
		return errs.RepositoryUnknownError
	}
	count, err := result.RowsAffected()
	if err != nil {
		r.logger.Error().Err(err).Msg("Error occurred while retrieving affected rows count after updating proxy profile")
		return errs.RepositoryUnknownError
	}
	if count == 0 {
		err = &errs.EntityNotFoundError{Name: "proxy profile", Key: "id", Value: profile.ID}
		r.logger.Print(err)
		return err
	}
	return nil
}

func (r *ProxyProfileRepository) Delete(ctx context.Context, id int) error {
	cmd := `DELETE FROM proxy_profiles WHERE id = ?`
	result, err := r.db.ExecContext(ctx, cmd, id)
	if err != nil {
		if s, ok := err.(sqlite3.Error); ok && s.Code == sqlite3.ErrConstraint {
			err := &errs.EntityReferencedError{Name: "proxy profile", Key: "id", Value: id}
			r.logger.Debug().Err(err).Send()
			return err
		}
		r.logger.Error().Err(err).Msg("Error occurred while deleting proxy profile")
		return errs.RepositoryUnknownError
	}
	count, err := result.RowsAffected()
	if err != nil {
		r.logger.Error().Err(err).Msg("Error occurred while retrieving affected rows count after deleting proxy profile")
		return errs.RepositoryUnknownError
	}
	if count == 0 {
		err = &errs.EntityNotFoundError{Name: "proxy profile", Key: "id", Value: id}
		r.logger.Print(err)
		return err
	}
	return nil
}
